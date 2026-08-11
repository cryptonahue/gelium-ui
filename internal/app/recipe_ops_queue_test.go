package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func getRecipe(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, path, nil))
	return res
}

func postRecipe(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)
	return res
}

// TestRecipeOpsQueueListRendersQueueWithAvatarToneBadgeAndActions proves the
// queue screen composes the primitives: two-line List rows with a decorative
// Avatar, a tone Badge carrying the status label, real POST actions per row and
// server-side pagination.
func TestRecipeOpsQueueListRendersQueueWithAvatarToneBadgeAndActions(t *testing.T) {
	resetRecipeQueueStore()

	res := getRecipe(t, "/recipes/ops-queue")
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1 class="recipe-oq-title">Work queue</h1>`,
		`<ul class="ui-list recipe-oq-list">`,
		`<li class="ui-list-item ui-list-item--two-line recipe-oq-item">`,
		`class="ui-avatar ui-avatar--sm"`,
		`aria-hidden="true"`,
		`class="ui-avatar-initials">CM</span>`,
		`class="ui-badge ui-badge-large ui-badge--error">Pending</span>`,
		`class="ui-badge ui-badge-large ui-badge--warning">Pending</span>`,
		`class="ui-badge ui-badge-large ui-badge--info">In progress</span>`,
		`method="post" action="/recipes/ops-queue/order-1042/advance"`,
		`method="post" action="/recipes/ops-queue/order-1042/dequeue"`,
		`id="queue-panel"`,
		`id="loom-toast-region"`,
		`<span class="ui-pagination-page ui-pagination-page--current" aria-current="page">1</span>`,
		`8 items · page 1 of 2`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("queue list is missing %q", contract)
		}
	}

	// Operational order: pending first (FIFO), then in_progress. The pending
	// invoice received earliest leads the page.
	if idxInvoice, idxCredit := strings.Index(body, "Invoice reconciliation"), strings.Index(body, "Credit note"); idxInvoice < 0 || idxCredit < 0 || idxInvoice > idxCredit {
		t.Errorf("pending FIFO order must lead with the earliest received invoice (invoice=%d credit=%d)", idxInvoice, idxCredit)
	}
}

// TestRecipeOpsQueueFiltersSanitizeAndPaginate proves the closed status/kind
// vocabularies filter the list, unknown values sanitize to "all", and page=2
// shows the next slice.
func TestRecipeOpsQueueFiltersSanitizeAndPaginate(t *testing.T) {
	resetRecipeQueueStore()

	res := getRecipe(t, "/recipes/ops-queue?status=done")
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`Login issue after migration`,
		`SLA policy announcement`,
		`2 items · page 1 of 1`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("done filter is missing %q", contract)
		}
	}
	if strings.Contains(body, "Invoice reconciliation") {
		t.Error("done filter must not include a pending item")
	}
	if strings.Contains(body, "ui-pagination") {
		t.Error("a single page of results must not render pagination")
	}

	res = getRecipe(t, "/recipes/ops-queue?status=bogus&kind=bogus")
	body = res.Body.String()
	if !strings.Contains(body, `8 items · page 1 of 2`) {
		t.Errorf("unknown vocabularies must sanitize to all, got caption %q", body)
	}

	res = getRecipe(t, "/recipes/ops-queue?page=2")
	body = res.Body.String()
	for _, contract := range []string{
		`<span class="ui-pagination-page ui-pagination-page--current" aria-current="page">2</span>`,
		`Payment webhook verification`,
		`SLA policy announcement`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("page 2 is missing %q", contract)
		}
	}
	if strings.Contains(body, "Invoice reconciliation") {
		t.Error("page 2 must not contain a row from page 1")
	}
}

// TestRecipeOpsQueueListHXFragment proves the HX-Request bifurcation: HTMX swaps
// only the #queue-panel fragment, not the whole page.
func TestRecipeOpsQueueListHXFragment(t *testing.T) {
	resetRecipeQueueStore()

	req := httptest.NewRequest(http.MethodGet, "/recipes/ops-queue?status=pending", nil)
	req.Header.Set("HX-Request", "true")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Contains(body, "<html") || strings.Contains(body, "<!doctype html>") {
		t.Error("HX request must return the queue panel fragment, not a full document")
	}
	for _, contract := range []string{`id="queue-panel"`, `Invoice reconciliation`} {
		if !strings.Contains(body, contract) {
			t.Errorf("queue panel fragment is missing %q", contract)
		}
	}
}

// TestRecipeOpsQueueAdvance303 proves the advance mutation follows POST+303 and
// the following render shows the persistent success banner and the new status.
func TestRecipeOpsQueueAdvance303(t *testing.T) {
	resetRecipeQueueStore()

	res := postRecipe(t, "/recipes/ops-queue/order-1042/advance")
	if res.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusSeeOther)
	}
	if loc := res.Header().Get("Location"); loc != "/recipes/ops-queue" {
		t.Errorf("Location = %q, want /recipes/ops-queue", loc)
	}

	res = getRecipe(t, "/recipes/ops-queue")
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`class="ui-banner ui-banner--success"`,
		`>Queue advanced</p>`,
		`was marked as in progress.`,
		`class="ui-badge ui-badge-large ui-badge--warning">In progress</span>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("advance result is missing %q", contract)
		}
	}

	// The flash banner is consumed after one render.
	res = getRecipe(t, "/recipes/ops-queue")
	if strings.Contains(res.Body.String(), "Queue advanced") {
		t.Error("the flash success banner must be consumed after one render")
	}
}

// TestRecipeOpsQueueAdvanceTerminalIsInfoBanner proves re-advancing a done item
// is a no-op reported with an informational banner, not another mutation.
func TestRecipeOpsQueueAdvanceTerminalIsInfoBanner(t *testing.T) {
	resetRecipeQueueStore()

	res := postRecipe(t, "/recipes/ops-queue/support-ticket-230/advance")
	if res.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusSeeOther)
	}
	body := getRecipe(t, "/recipes/ops-queue").Body.String()
	for _, contract := range []string{`class="ui-banner ui-banner--info"`, `>Already completed</p>`} {
		if !strings.Contains(body, contract) {
			t.Errorf("terminal advance is missing %q", contract)
		}
	}
}

// TestRecipeOpsQueueDequeue303 proves the dequeue mutation follows POST+303 and
// removes the row.
func TestRecipeOpsQueueDequeue303(t *testing.T) {
	resetRecipeQueueStore()

	res := postRecipe(t, "/recipes/ops-queue/billing-invoice-77/dequeue")
	if res.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusSeeOther)
	}
	body := getRecipe(t, "/recipes/ops-queue").Body.String()
	for _, contract := range []string{`class="ui-banner ui-banner--success"`, `>Item removed</p>`} {
		if !strings.Contains(body, contract) {
			t.Errorf("dequeue result is missing %q", contract)
		}
	}
	if strings.Contains(body, `ui-list-item-headline">Invoice reconciliation for Q2`) {
		t.Error("dequeue must remove the item from the queue")
	}
}

// TestRecipeOpsQueueNotFound proves unknown item ids render the shared
// error-state page with the real 404 status on both mutations.
func TestRecipeOpsQueueNotFound(t *testing.T) {
	resetRecipeQueueStore()

	for _, path := range []string{"/recipes/ops-queue/nope/advance", "/recipes/ops-queue/nope/dequeue"} {
		res := postRecipe(t, path)
		if res.Code != http.StatusNotFound {
			t.Errorf("POST %s status = %d, want %d", path, res.Code, http.StatusNotFound)
		}
		if body := res.Body.String(); !strings.Contains(body, `class="ui-error-state"`) {
			t.Errorf("POST %s must render the error-state page", path)
		}
	}
}

// TestRecipeOpsQueueEmptyState proves the empty state offers a real CTA: clear
// the filters when filtering hides everything, otherwise surface completed work.
func TestRecipeOpsQueueEmptyState(t *testing.T) {
	resetRecipeQueueStore()

	res := getRecipe(t, "/recipes/ops-queue?status=blocked&kind=billing")
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{`No matching items`, `Clear filters`, `href="/recipes/ops-queue"`} {
		if !strings.Contains(body, contract) {
			t.Errorf("filtered empty state is missing %q", contract)
		}
	}

	resetRecipeQueueStore()
	for _, it := range queueDemoStore.snapshot() {
		queueDemoStore.delete(it.ID)
	}
	body = getRecipe(t, "/recipes/ops-queue").Body.String()
	for _, contract := range []string{`Queue is clear`, `href="/recipes/ops-queue?status=done"`, `View completed`} {
		if !strings.Contains(body, contract) {
			t.Errorf("empty queue state is missing %q", contract)
		}
	}
}

// TestRecipeOpsQueueRefreshPostOnly proves the refresh action is POST-only (GET
// answers 405 with Allow: POST), the no-JS refresh re-renders the list with a
// persistent inline toast + progress, and the HTMX refresh returns the fragment
// plus an HX-Trigger loom:toast.
func TestRecipeOpsQueueRefreshPostOnly(t *testing.T) {
	resetRecipeQueueStore()

	res := getRecipe(t, "/recipes/ops-queue/refresh")
	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET refresh status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
	if got := res.Header().Get("Allow"); got != "POST" {
		t.Errorf("Allow = %q, want POST", got)
	}

	res = postRecipe(t, "/recipes/ops-queue/refresh")
	if res.Code != http.StatusOK {
		t.Fatalf("POST refresh status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`Queue refreshed.`,
		`<progress value="100" max="100"></progress>`,
		`class="ui-toast ui-toast-success"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("no-JS refresh is missing %q", contract)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/recipes/ops-queue/refresh", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	res = httptest.NewRecorder()
	New().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("HX refresh status = %d, want %d", res.Code, http.StatusOK)
	}
	if trigger := res.Header().Get("HX-Trigger"); !strings.Contains(trigger, "loom:toast") || !strings.Contains(trigger, "Queue refreshed.") {
		t.Errorf("HX-Trigger = %q, want loom:toast payload", trigger)
	}
	body = res.Body.String()
	if strings.Contains(body, "<html") {
		t.Error("HX refresh must return the refresh form fragment, not a full document")
	}
	if !strings.Contains(body, `class="recipe-oq-refresh"`) {
		t.Error("HX refresh must return the refresh form fragment")
	}
}

// TestRecipeOpsQueueNoindexProvesRecipeSurfacesAreNeverIndexed proves the queue
// route emits robots noindex, nofollow plus a canonical and a per-route
// description in the recipe's own document head.
func TestRecipeOpsQueueNoindexProvesRecipeSurfacesAreNeverIndexed(t *testing.T) {
	resetRecipeQueueStore()

	res := getRecipe(t, "/recipes/ops-queue")
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`lang="en"`,
		`<meta name="robots" content="noindex, nofollow">`,
		`<link rel="canonical" href="https://gelium-ui.example/recipes/ops-queue">`,
		`<title>Work queue · Ops Queue recipe</title>`,
		`<meta name="description"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("queue page is missing %q", contract)
		}
	}
}

// TestRecipeQueueToneDerivation proves the server-side tone derivation for the
// closed statuses and SLA deadlines.
func TestRecipeQueueToneDerivation(t *testing.T) {
	now := time.Date(2026, 8, 11, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		item recipeQueueItem
		want string
	}{
		{recipeQueueItem{Status: "pending", SLADeadline: now.Add(-10 * time.Minute)}, "error"},
		{recipeQueueItem{Status: "pending", SLADeadline: now.Add(30 * time.Minute)}, "warning"},
		{recipeQueueItem{Status: "pending", SLADeadline: now.Add(3 * time.Hour)}, "info"},
		{recipeQueueItem{Status: "in_progress", SLADeadline: now.Add(3 * time.Hour)}, "info"},
		{recipeQueueItem{Status: "in_progress", SLADeadline: now.Add(30 * time.Minute)}, "warning"},
		{recipeQueueItem{Status: "blocked", SLADeadline: now.Add(3 * time.Hour)}, "error"},
		{recipeQueueItem{Status: "done", SLADeadline: now.Add(-4 * time.Hour)}, "success"},
	}
	for _, tc := range cases {
		if got := recipeQueueItemTone(tc.item, now); got != tc.want {
			t.Errorf("tone for %+v = %q, want %q", tc.item, got, tc.want)
		}
	}
}
