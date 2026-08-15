package app

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Ops Queue screen recipe (Phase G): a server-rendered work queue. Each row is
// a List two-line item composed from the Avatar, Badge tone and Button
// primitives; the order is operational (server-side status rank, then FIFO by
// received time — position is state, never presentation), every transition is
// a real POST+303 with a persistent success banner, and the refresh reuses the
// Data table refresh contract (HX-Request bifurcation + loom:toast). The recipe
// consumes the Avatar, badge tone and standalone pagination primitives
// introduced alongside it; the queue itself introduces no new primitives.

// recipeOpsQueuePageSize controls the server-side pagination slice.
const recipeOpsQueuePageSize = 5

// recipeQueueStatuses and recipeQueueKinds are the closed filter vocabularies;
// unknown values sanitize to "" (all). Status drives both the operational order
// and the badge tone; kind is a plain category filter.
var recipeQueueStatuses = []string{"pending", "in_progress", "done", "blocked"}

var recipeQueueStatusLabels = map[string]string{
	"pending":     "Pending",
	"in_progress": "In progress",
	"done":        "Done",
	"blocked":     "Blocked",
}

var recipeQueueKinds = []string{"message", "order", "support", "billing"}

var recipeQueueKindLabels = map[string]string{
	"message": "Message",
	"order":   "Order",
	"support": "Support",
	"billing": "Billing",
}

// recipeQueueItem is one unit of work in the queue.
type recipeQueueItem struct {
	ID          string
	Subject     string
	Requester   string
	Kind        string
	Status      string
	Assignee    string
	ReceivedAt  time.Time
	SLADeadline time.Time
}

// recipeQueueStore is the in-memory mock store behind the recipe. A single
// mutex guards the slice and the flash banner so handlers mutate state without
// racing; the dataset lives only in server memory (like the WhatsApp demo
// store) — there is no persistence layer in this recipe.
type recipeQueueStore struct {
	mu     sync.Mutex
	seq    int
	items  []recipeQueueItem
	banner *bannerView
}

// queueDemoStore is the shared demo store for the recipe routes.
var queueDemoStore = newRecipeQueueStore()

// resetRecipeQueueStore restores the demo store to its seed state. Tests use it
// to stay deterministic regardless of execution order.
func resetRecipeQueueStore() {
	queueDemoStore = newRecipeQueueStore()
}

// newRecipeQueueStore seeds the mock dataset. Received times and SLA deadlines
// are relative to the caller's clock so the demo always looks fresh; the tone
// derivation below only flips across generous boundaries, keeping renders
// deterministic within any realistic run.
func newRecipeQueueStore() *recipeQueueStore {
	now := time.Now()
	return &recipeQueueStore{
		items: []recipeQueueItem{
			{ID: "billing-invoice-77", Subject: "Invoice reconciliation for Q2", Requester: "Carla M.", Kind: "billing", Status: "pending", Assignee: "", ReceivedAt: now.Add(-3 * time.Hour), SLADeadline: now.Add(-10 * time.Minute)},
			{ID: "billing-credit-12", Subject: "Credit note for cancelled plan", Requester: "Bob T.", Kind: "billing", Status: "pending", Assignee: "Alicia R.", ReceivedAt: now.Add(-30 * time.Minute), SLADeadline: now.Add(3 * time.Hour)},
			{ID: "order-1042", Subject: "Order 1042 — expedited shipping", Requester: "Alicia R.", Kind: "order", Status: "pending", Assignee: "", ReceivedAt: now.Add(-5 * time.Minute), SLADeadline: now.Add(30 * time.Minute)},
			{ID: "order-1051", Subject: "Order 1051 — damaged box on arrival", Requester: "Dev Ops", Kind: "order", Status: "in_progress", Assignee: "Bob T.", ReceivedAt: now.Add(-6 * time.Hour), SLADeadline: now.Add(3 * time.Hour)},
			{ID: "support-ticket-221", Subject: "Refund request for batch purchase", Requester: "Bob T.", Kind: "support", Status: "in_progress", Assignee: "Carla M.", ReceivedAt: now.Add(-20 * time.Minute), SLADeadline: now.Add(3 * time.Hour)},
			{ID: "message-probe", Subject: "Payment webhook verification", Requester: "Dev Ops", Kind: "message", Status: "blocked", Assignee: "Alicia R.", ReceivedAt: now.Add(-1 * time.Hour), SLADeadline: now.Add(-10 * time.Minute)},
			{ID: "support-ticket-230", Subject: "Login issue after migration", Requester: "Alicia R.", Kind: "support", Status: "done", Assignee: "Carla M.", ReceivedAt: now.Add(-26 * time.Hour), SLADeadline: now.Add(-4 * time.Hour)},
			{ID: "message-announce", Subject: "SLA policy announcement", Requester: "Ops Lead", Kind: "message", Status: "done", Assignee: "", ReceivedAt: now.Add(-30 * time.Hour), SLADeadline: now.Add(-4 * time.Hour)},
		},
	}
}

// snapshot returns a copy of every item so handlers never hold a live reference
// into the store's slice.
func (s *recipeQueueStore) snapshot() []recipeQueueItem {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]recipeQueueItem, len(s.items))
	copy(out, s.items)
	return out
}

func (s *recipeQueueStore) get(id string) (recipeQueueItem, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, it := range s.items {
		if it.ID == id {
			return it, true
		}
	}
	return recipeQueueItem{}, false
}

func (s *recipeQueueStore) updateStatus(id, status string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].Status = status
			return true
		}
	}
	return false
}

func (s *recipeQueueStore) delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, it := range s.items {
		if it.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}
	return false
}

// setBanner stores a persistent flash banner for the next list render (the
// POST+303 contract: the mutation redirects, the following GET shows it and
// consumes it).
func (s *recipeQueueStore) setBanner(b bannerView) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.banner = &b
}

func (s *recipeQueueStore) takeBanner() *bannerView {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.banner == nil {
		return nil
	}
	b := s.banner
	s.banner = nil
	return b
}

// ----- view models -----

// recipeOpsQueueView is the list screen: app shell + banner slot + filter form
// + the queue panel (two-line List rows composed from Avatar + Badge tone +
// Button) + standalone pagination + the remote refresh form.
type recipeOpsQueueView struct {
	Meta          metaView
	ThemeClass    string
	Title         string
	Description   string
	FilterAction  string
	StatusValue   string
	KindValue     string
	StatusOptions []recipeQueueOption
	KindOptions   []recipeQueueOption
	Items         []recipeQueueItemView
	Caption       string
	Pagination    *paginationView
	EmptyState    emptyStateView
	Banner        *bannerView
	Refreshed     bool
	RefreshToast  *toastView
}

// recipeQueueOption is one closed-vocabulary filter option.
type recipeQueueOption struct {
	Value string
	Label string
}

// recipeQueueItemView is one rendered queue row with its per-row actions.
type recipeQueueItemView struct {
	ID            string
	Subject       string
	Requester     string
	Avatar        avatarView
	KindLabel     string
	StatusLabel   string
	Tone          string
	Assignee      string
	ReceivedText  string
	SLAText       string
	AdvanceAction string
	DequeueAction string
	CanAdvance    bool
}

// ----- handlers -----

func (s *server) recipeOpsQueueList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	view := newRecipeOpsQueueView(q.Get("status"), q.Get("kind"), q.Get("page"), queueDemoStore.takeBanner())
	applyRequestTheme(r, view)

	if strings.EqualFold(r.Header.Get("HX-Request"), "true") {
		s.renderRecipeTemplate(w, http.StatusOK, "recipe-ops-queue-panel", view)
		return
	}
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-ops-queue", view)
}

// recipeOpsQueueAdvance moves an item to its next operational state and
// redirects (POST+303) with a persistent success banner. Done is terminal; a
// re-advance of a done item is a no-op reported with an informational banner.
func (s *server) recipeOpsQueueAdvance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, ok := queueDemoStore.get(id)
	if !ok {
		s.recipeQueueNotFound(w, "Item not found", "The queue item you are trying to advance does not exist or has been removed.")
		return
	}
	next := recipeQueueAdvanceStatus(item.Status)
	if next == item.Status {
		queueDemoStore.setBanner(recipeQueueInfoBanner("Already completed", fmt.Sprintf("%q is already done.", item.Subject)))
	} else {
		queueDemoStore.updateStatus(id, next)
		queueDemoStore.setBanner(recipeQueueSuccessBanner("Queue advanced", fmt.Sprintf("%q was marked as %s.", item.Subject, strings.ToLower(recipeQueueStatusLabels[next]))))
	}
	http.Redirect(w, r, "/recipes/ops-queue", http.StatusSeeOther)
}

// recipeOpsQueueDequeue removes an item from the queue entirely (POST+303 +
// persistent success banner).
func (s *server) recipeOpsQueueDequeue(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, ok := queueDemoStore.get(id)
	if !ok {
		s.recipeQueueNotFound(w, "Item not found", "The queue item you are trying to remove does not exist or has already been removed.")
		return
	}
	queueDemoStore.delete(id)
	queueDemoStore.setBanner(recipeQueueSuccessBanner("Item removed", fmt.Sprintf("%q was removed from the queue.", item.Subject)))
	http.Redirect(w, r, "/recipes/ops-queue", http.StatusSeeOther)
}

// recipeOpsQueueRefresh completes a remote refresh for the queue. With HTMX it
// returns the refresh form fragment and an HX-Trigger raising loom:toast;
// without JavaScript it re-renders the full page with the queue refreshed and a
// persistent inline toast, exactly like the Data table refresh demo.
func (s *server) recipeOpsQueueRefresh(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	isHX := strings.EqualFold(r.Header.Get("HX-Request"), "true")

	view := newRecipeOpsQueueView("", "", "", queueDemoStore.takeBanner())
	applyRequestTheme(r, view)
	view.Refreshed = true

	if isHX {
		trigger, err := toastTriggerJSON("success", "Queue refreshed.")
		if err != nil {
			http.Error(w, "refresh unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Trigger", trigger)
		s.renderRecipeTemplate(w, http.StatusOK, "recipe-ops-queue-refresh-form", view)
		return
	}

	toast := newToast("success", "recipe-oq-refresh-result", "Queue refreshed.")
	view.RefreshToast = &toast
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-ops-queue", view)
}

func (s *server) recipeQueueNotFound(w http.ResponseWriter, title, body string) {
	s.renderErrorPage(w, http.StatusNotFound, title, body, true, "/recipes/ops-queue", "Back to the queue", "/recipes/ops-queue")
}

// ----- view builders -----

// newRecipeOpsQueueView validates the request against the closed vocabularies
// (status/kind/page), filters, applies the operational sort (status rank then
// FIFO) and paginates the slice into a view model.
func newRecipeOpsQueueView(statusParam, kindParam, page string, banner *bannerView) *recipeOpsQueueView {
	statusValue := recipeClosedValue(statusParam, recipeQueueStatuses)
	kindValue := recipeClosedValue(kindParam, recipeQueueKinds)
	pageNum := recipeParsePage(page)

	now := time.Now()
	snapshot := queueDemoStore.snapshot()
	items := make([]recipeQueueItem, 0, len(snapshot))
	for _, it := range snapshot {
		if statusValue != "" && it.Status != statusValue {
			continue
		}
		if kindValue != "" && it.Kind != kindValue {
			continue
		}
		items = append(items, it)
	}

	sort.SliceStable(items, func(i, j int) bool {
		ri, rj := recipeQueueRank(items[i].Status), recipeQueueRank(items[j].Status)
		if ri != rj {
			return ri < rj
		}
		return items[i].ReceivedAt.Before(items[j].ReceivedAt)
	})

	total := len(items)
	totalPages := (total + recipeOpsQueuePageSize - 1) / recipeOpsQueuePageSize
	if totalPages < 1 {
		totalPages = 1
	}
	if pageNum > totalPages {
		pageNum = totalPages
	}

	start := (pageNum - 1) * recipeOpsQueuePageSize
	end := start + recipeOpsQueuePageSize
	if end > total {
		end = total
	}

	rows := make([]recipeQueueItemView, 0, end-start)
	for _, it := range items[start:end] {
		rows = append(rows, recipeQueueRow(it, now))
	}

	var pagination *paginationView
	if totalPages > 1 {
		pagination = newPaginationView(pageNum, totalPages, func(n int) string {
			return recipeOpsQueueHref(statusValue, kindValue, n)
		})
	}

	meta := screenRecipeMeta(
		"Work queue · Ops Queue recipe",
		"Attend the work queue: triage items by status and SLA, advance them to the next state and remove them, server-side with the Ops Queue screen recipe.",
		"/recipes/ops-queue",
	)

	return &recipeOpsQueueView{
		Meta:          meta,
		ThemeClass:    themeClass(""),
		Title:         "Work queue",
		Description:   "The Ops Queue screen recipe: a server-rendered work queue composed from List, Avatar, badge tones, Button, Toast, Empty state and Banner primitives. Order is operational state, not presentation.",
		FilterAction:  "/recipes/ops-queue",
		StatusValue:   statusValue,
		KindValue:     kindValue,
		StatusOptions: recipeQueueOptions(recipeQueueStatuses, recipeQueueStatusLabels),
		KindOptions:   recipeQueueOptions(recipeQueueKinds, recipeQueueKindLabels),
		Items:         rows,
		Caption:       fmt.Sprintf("%d items · page %d of %d", total, pageNum, totalPages),
		Pagination:    pagination,
		EmptyState:    recipeQueueEmptyState(statusValue, kindValue),
		Banner:        banner,
	}
}

// recipeQueueRow renders one item: a decorative sm Avatar (paired with the
// visible requester name), the two-line List text, a tone Badge carrying the
// status label and the per-row POST actions.
func recipeQueueRow(it recipeQueueItem, now time.Time) recipeQueueItemView {
	next := recipeQueueAdvanceStatus(it.Status)
	return recipeQueueItemView{
		ID:            it.ID,
		Subject:       it.Subject,
		Requester:     it.Requester,
		Avatar:        avatarView{Initials: recipeInitials(it.Requester), Decorative: true, Size: "sm"},
		KindLabel:     recipeQueueKindLabels[it.Kind],
		StatusLabel:   recipeQueueStatusLabels[it.Status],
		Tone:          recipeQueueItemTone(it, now),
		Assignee:      it.Assignee,
		ReceivedText:  recipeRelativeTime(it.ReceivedAt, now),
		SLAText:       recipeQueueSLAText(it, now),
		AdvanceAction: "/recipes/ops-queue/" + it.ID + "/advance",
		DequeueAction: "/recipes/ops-queue/" + it.ID + "/dequeue",
		CanAdvance:    next != it.Status,
	}
}

// recipeQueueRank is the operational priority of a status: pending items are
// attended first, then in progress, blocked, and finally done.
func recipeQueueRank(status string) int {
	switch status {
	case "pending":
		return 0
	case "in_progress":
		return 1
	case "blocked":
		return 2
	case "done":
		return 3
	}
	return 9
}

// recipeQueueAdvanceStatus is the closed state machine: pending → in_progress →
// done; blocked re-enters in_progress; done is terminal.
func recipeQueueAdvanceStatus(status string) string {
	switch status {
	case "pending":
		return "in_progress"
	case "in_progress":
		return "done"
	case "blocked":
		return "in_progress"
	}
	return status
}

// recipeQueueItemTone derives the Badge tone server-side: done is success,
// blocked is error, and pending/in_progress items follow the SLA deadline
// (overdue → error, near the deadline → warning, otherwise → info). The status
// label is always visible next to the tone — never color-only.
func recipeQueueItemTone(it recipeQueueItem, now time.Time) string {
	switch it.Status {
	case "done":
		return "success"
	case "blocked":
		return "error"
	}
	if now.After(it.SLADeadline) {
		return "error"
	}
	if now.Add(time.Hour).After(it.SLADeadline) {
		return "warning"
	}
	return "info"
}

// recipeQueueSLAText reports the deadline state as a visible label.
func recipeQueueSLAText(it recipeQueueItem, now time.Time) string {
	if it.Status == "done" {
		return "Completed"
	}
	d := it.SLADeadline.Sub(now)
	if d < 0 {
		return "SLA overdue"
	}
	return "SLA in " + recipeHumanDuration(d)
}

// recipeRelativeTime formats a timestamp relative to now ("just now", "5m ago",
// "3h ago", "2d ago") — the terseness contract of a scan surface.
func recipeRelativeTime(t, now time.Time) string {
	d := now.Sub(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d/time.Minute))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d/time.Hour))
	default:
		return fmt.Sprintf("%dd ago", int(d/(24*time.Hour)))
	}
}

func recipeHumanDuration(d time.Duration) string {
	switch {
	case d >= 24*time.Hour:
		return fmt.Sprintf("%dd", int(d/(24*time.Hour)))
	case d >= time.Hour:
		return fmt.Sprintf("%dh", int(d/time.Hour))
	case d >= time.Minute:
		return fmt.Sprintf("%dm", int(d/time.Minute))
	}
	return "<1m"
}

// recipeQueueEmptyState offers a real CTA: clear the filters when filtering
// hides everything, otherwise surface the completed work.
func recipeQueueEmptyState(statusValue, kindValue string) emptyStateView {
	if statusValue != "" || kindValue != "" {
		return emptyStateView{
			Title:    "No matching items",
			Body:     "No queue items match the selected filters. Clear them to see the full queue.",
			CTA:      true,
			CTAHref:  "/recipes/ops-queue",
			CTALabel: "Clear filters",
		}
	}
	return emptyStateView{
		Title:    "Queue is clear",
		Body:     "Every item is done. Review completed work or wait for new items.",
		CTA:      true,
		CTAHref:  "/recipes/ops-queue?status=done",
		CTALabel: "View completed",
	}
}

// recipeOpsQueueHref builds a real GET link preserving the current filters and
// page (parameter order is stable: status, kind, page).
func recipeOpsQueueHref(status, kind string, page int) string {
	var b strings.Builder
	b.WriteByte('?')
	write := func(key, value string) {
		if b.Len() > 1 {
			b.WriteByte('&')
		}
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteString(url.QueryEscape(value))
	}
	if status != "" {
		write("status", status)
	}
	if kind != "" {
		write("kind", kind)
	}
	if page >= 1 {
		write("page", strconv.Itoa(page))
	}
	return b.String()
}

// ----- shared recipe helpers -----

// recipeClosedValue sanitizes a request value against a closed vocabulary,
// returning "" for anything unknown (the "all" default).
func recipeClosedValue(v string, vocab []string) string {
	for _, k := range vocab {
		if v == k {
			return k
		}
	}
	return ""
}

func recipeParsePage(page string) int {
	if n, err := strconv.Atoi(page); err == nil && n >= 1 {
		return n
	}
	return 1
}

func recipeQueueOptions(vocab []string, labels map[string]string) []recipeQueueOption {
	out := make([]recipeQueueOption, 0, len(vocab))
	for _, v := range vocab {
		out = append(out, recipeQueueOption{Value: v, Label: labels[v]})
	}
	return out
}

func recipeQueueSuccessBanner(title, body string) bannerView {
	return bannerView{Tone: "success", Title: title, Body: body}
}

func recipeQueueInfoBanner(title, body string) bannerView {
	return bannerView{Tone: "info", Title: title, Body: body}
}

// recipeInitials derives the avatar initials from the first letters of the
// first two name words ("Alicia R." → "AR", "Dev Ops" → "DO").
func recipeInitials(name string) string {
	fields := strings.Fields(name)
	if len(fields) == 0 {
		return "?"
	}
	var b strings.Builder
	n := len(fields)
	if n > 2 {
		n = 2
	}
	for _, f := range fields[:n] {
		if f == "" {
			continue
		}
		b.WriteString(strings.ToUpper(f[:1]))
	}
	if b.Len() == 0 {
		return "?"
	}
	return b.String()
}

// screenRecipeMeta builds the server-driven metadata for a recipe surface.
// Recipe routes live under /recipes/* and are never indexed (noindex, nofollow
// — demo surfaces, exactly like /demo/*).
func screenRecipeMeta(title, description, routePath string) metaView {
	return metaView{
		Title:         title,
		Description:   description,
		Lang:          "en",
		Robots:        "noindex, nofollow",
		Canonical:     siteBaseURL + routePath,
		OGType:        "website",
		OGTitle:       title + " · Gelium UI",
		OGDescription: description,
	}
}
