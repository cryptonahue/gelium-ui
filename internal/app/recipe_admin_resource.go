package app

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Admin Resource screen recipe (Phase G): a server-rendered resource manager
// that composes existing primitives into a full screen. The list reuses the
// Data table server-driven pattern (sort/filter/page/selection via stable GET
// params), the create/edit forms reuse the Text field + Select components with
// the 422 validation contract, the delete confirmation is a Dialog page
// variant rendered as a real native <dialog open>, and every mutation follows
// POST+303 with a persistent success banner (never a toast). The recipe
// introduces no new primitives — only wiring.

// recipeAdminPageSize controls the server-side pagination slice of the recipe.
const recipeAdminPageSize = 5

// recipeResource is one row in the Admin Resource dataset.
type recipeResource struct {
	ID     string
	Name   string
	Status string
	Date   string
	Owner  string
}

// recipeResourceStore is the in-memory mock store behind the recipe. A single
// mutex guards the slice and the flash banner, so handlers mutate state without
// racing; the dataset lives only in server memory (like the WhatsApp demo
// store) — there is no persistence layer in this recipe.
type recipeResourceStore struct {
	mu     sync.Mutex
	seq    int
	items  []recipeResource
	banner *bannerView
}

// resourceDemoStore is the shared demo store for the recipe routes.
var resourceDemoStore = newRecipeResourceStore()

// resetRecipeResourceStore restores the demo store to its seed state. Tests use
// it to stay deterministic regardless of execution order; a real deployment
// would reset only at process start.
func resetRecipeResourceStore() {
	resourceDemoStore = newRecipeResourceStore()
}

// newRecipeResourceStore seeds the mock dataset. Dates are ISO-8601 so string
// ordering equals chronological ordering (reuses dataTableSortKeys).
func newRecipeResourceStore() *recipeResourceStore {
	return &recipeResourceStore{
		items: []recipeResource{
			{ID: "alpha", Name: "Alpha release", Status: "Active", Date: "2026-01-12", Owner: "Alicia R."},
			{ID: "beta", Name: "Beta rollout", Status: "Pending", Date: "2026-02-03", Owner: "Bob T."},
			{ID: "gamma", Name: "Gamma refactor", Status: "Done", Date: "2026-03-15", Owner: "Carla M."},
			{ID: "delta", Name: "Delta docs", Status: "Active", Date: "2026-04-28", Owner: "Alicia R."},
			{ID: "epsilon", Name: "Epsilon migration", Status: "Pending", Date: "2026-05-06", Owner: "Dev Ops"},
			{ID: "zeta", Name: "Zeta dashboard", Status: "Done", Date: "2026-06-02", Owner: "Bob T."},
			{ID: "eta", Name: "Eta audit", Status: "Active", Date: "2026-07-19", Owner: "Carla M."},
			{ID: "theta", Name: "Theta backup", Status: "Pending", Date: "2026-08-07", Owner: "Dev Ops"},
			{ID: "iota", Name: "Iota index", Status: "Done", Date: "2026-09-11", Owner: "Alicia R."},
			{ID: "kappa", Name: "Kappa layout", Status: "Active", Date: "2026-10-23", Owner: "Bob T."},
			{ID: "lambda", Name: "Lambda mailer", Status: "Pending", Date: "2026-11-05", Owner: "Carla M."},
			{ID: "mu", Name: "Mu navigation", Status: "Done", Date: "2026-12-14", Owner: "Dev Ops"},
		},
	}
}

// snapshot returns a copy of every item so handlers never hold a live reference
// into the store's slice.
func (s *recipeResourceStore) snapshot() []recipeResource {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]recipeResource, len(s.items))
	copy(out, s.items)
	return out
}

func (s *recipeResourceStore) get(id string) (recipeResource, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, it := range s.items {
		if it.ID == id {
			return it, true
		}
	}
	return recipeResource{}, false
}

// create inserts a new resource with a stable slug ID and returns the stored
// item. Callers validate the fields first; the store only persists.
func (s *recipeResourceStore) create(name, status, date, owner string) recipeResource {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	base := recipeResourceSlug(name)
	id := base
	for _, it := range s.items {
		if it.ID == id {
			id = fmt.Sprintf("%s-%d", base, s.seq)
		}
	}
	item := recipeResource{ID: id, Name: name, Status: status, Date: date, Owner: owner}
	s.items = append(s.items, item)
	return item
}

func (s *recipeResourceStore) update(id, name, status, date, owner string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].Name = name
			s.items[i].Status = status
			s.items[i].Date = date
			s.items[i].Owner = owner
			return true
		}
	}
	return false
}

func (s *recipeResourceStore) delete(id string) bool {
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
func (s *recipeResourceStore) setBanner(b bannerView) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.banner = &b
}

func (s *recipeResourceStore) takeBanner() *bannerView {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.banner == nil {
		return nil
	}
	b := s.banner
	s.banner = nil
	return b
}

// recipeResourceSlug builds a stable, deep-linkable ID from a display name.
func recipeResourceSlug(name string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case r == ' ' || r == '-' || r == '_':
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	if slug := strings.Trim(b.String(), "-"); slug != "" {
		return slug
	}
	return "project"
}

// ----- view models -----

// recipeAdminResourceView is the list screen: app shell + banner slot + filter
// form + the Data table panel (selection form + panel + pagination) + the
// remote refresh form. The columns, page links and empty state reuse the Data
// table view models.
type recipeAdminResourceView struct {
	Meta            metaView
	ThemeClass      string
	DataTheme       string
	Title           string
	Description     string
	NewButton       buttonView
	FilterAction    string
	SearchField     textFieldView
	Query           string
	Sort            string
	Dir             string
	Page            int
	Total           int
	Pages           int
	Caption         string
	Columns         []dataTableColumn
	Rows            []recipeResourceRowView
	PageLinks       []dataTablePageView
	PrevHref        string
	NextHref        string
	HasPrev         bool
	HasNext         bool
	SelectedCount   int
	SelectionNotice string
	Colspan         int
	EmptyState      emptyStateView
	Banner          *bannerView
	Refreshed       bool
	RefreshToast    *toastView
}

// recipeResourceRowView is one rendered table row with per-row actions.
type recipeResourceRowView struct {
	ID         string
	Name       string
	Status     string
	Date       string
	Owner      string
	Selected   bool
	EditHref   string
	DeleteHref string
}

// recipeAdminResourceFormView is the create/edit form screen. Validation
// failures re-render this view with the inline alert + validation summary set
// and per-field errors, preserving the submitted values.
type recipeAdminResourceFormView struct {
	Meta          metaView
	ThemeClass    string
	DataTheme     string
	Heading       string
	Intro         string
	Action        string
	CancelHref    string
	SubmitLabel   string
	NameField     textFieldView
	StatusValue   string
	StatusOptions []recipeStatusOption
	StatusError   string
	DateField     textFieldView
	OwnerField    textFieldView
	InlineAlert   *inlineAlertView
	Validation    *validationSummaryView
}

// recipeStatusOption is one closed-vocabulary status option (reuses
// dataTableStatuses).
type recipeStatusOption struct {
	Value string
	Label string
}

// recipeAdminResourceConfirmView is the delete confirmation screen: a Dialog
// page variant (real native <dialog open>) with the destructive action as a
// real POST form.
type recipeAdminResourceConfirmView struct {
	Meta        metaView
	ThemeClass  string
	DataTheme   string
	Item        recipeResource
	DeleteHref  string
	DeleteLabel string
	ListHref    string
}

// validationSummaryView is the production view model for the shared
// "validation-summary" partial. The items are real links to each field error
// anchor, so a no-JS form failure navigates straight to the broken field.
type validationSummaryView struct {
	Title        string
	HeadingLevel int
	Items        []validationSummaryItemView
}

type validationSummaryItemView struct {
	Href    string
	Message string
}

// ----- handlers -----

func (s *server) recipeAdminResourceList(w http.ResponseWriter, r *http.Request) {
	selection := r.URL.Query()["selection"]
	view := newRecipeAdminResourceView(
		r.URL.Query().Get("q"), r.URL.Query().Get("sort"), r.URL.Query().Get("dir"), r.URL.Query().Get("page"),
		selection, resourceDemoStore.takeBanner(),
	)
	applyRequestChrome(r, view)

	if strings.EqualFold(r.Header.Get("HX-Request"), "true") {
		s.renderRecipeTemplate(w, http.StatusOK, "recipe-admin-resource-panel", view)
		return
	}
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-admin-resource-list", view)
}

func (s *server) recipeAdminResourceNew(w http.ResponseWriter, r *http.Request) {
	view := newRecipeAdminResourceFormView("create", "", "", "Pending", "", "", nil)
	applyRequestChrome(r, view)
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-admin-resource-form", view)
}

func (s *server) recipeAdminResourceCreate(w http.ResponseWriter, r *http.Request) {
	name, status, date, owner := s.parseRecipeResourceForm(r)
	if errs := validateRecipeResourceForm(name, status, date); len(errs) > 0 {
		s.recipeAdminResourceFormInvalid(w, r, "create", "", name, status, date, owner, errs)
		return
	}
	item := resourceDemoStore.create(name, status, date, owner)
	resourceDemoStore.setBanner(recipeAdminSuccessBanner("Project created", fmt.Sprintf("%q was added to the projects list.", item.Name)))
	http.Redirect(w, r, "/recipes/admin-resource", http.StatusSeeOther)
}

func (s *server) recipeAdminResourceEdit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, ok := resourceDemoStore.get(id)
	if !ok {
		s.recipeAdminResourceNotFound(w, "Project not found", "The project you are trying to edit does not exist or has been deleted.")
		return
	}
	view := newRecipeAdminResourceFormView("edit", id, item.Name, item.Status, item.Date, item.Owner, nil)
	applyRequestChrome(r, view)
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-admin-resource-form", view)
}

func (s *server) recipeAdminResourceUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, ok := resourceDemoStore.get(id); !ok {
		s.recipeAdminResourceNotFound(w, "Project not found", "The project you are trying to edit does not exist or has been deleted.")
		return
	}
	name, status, date, owner := s.parseRecipeResourceForm(r)
	if errs := validateRecipeResourceForm(name, status, date); len(errs) > 0 {
		s.recipeAdminResourceFormInvalid(w, r, "edit", id, name, status, date, owner, errs)
		return
	}
	resourceDemoStore.update(id, name, status, date, owner)
	resourceDemoStore.setBanner(recipeAdminSuccessBanner("Project updated", fmt.Sprintf("%q was saved.", name)))
	http.Redirect(w, r, "/recipes/admin-resource", http.StatusSeeOther)
}

func (s *server) recipeAdminResourceDeleteConfirm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, ok := resourceDemoStore.get(id)
	if !ok {
		s.recipeAdminResourceNotFound(w, "Project not found", "The project you are trying to delete does not exist or has already been removed.")
		return
	}
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-admin-resource-confirm", recipeAdminResourceConfirmView{
		Meta:        recipeAdminMeta("Delete "+item.Name+" · Admin Resource recipe", "Confirm the deletion of a project in the Admin Resource screen recipe.", "/recipes/admin-resource/"+id+"/delete"),
		ThemeClass:  themeClass(""),
		Item:        item,
		DeleteHref:  "/recipes/admin-resource/" + id + "/delete",
		DeleteLabel: "Delete " + item.Name,
		ListHref:    "/recipes/admin-resource",
	})
}

func (s *server) recipeAdminResourceDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, ok := resourceDemoStore.get(id)
	if !ok {
		s.recipeAdminResourceNotFound(w, "Project not found", "The project you are trying to delete does not exist or has already been removed.")
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	resourceDemoStore.delete(id)
	resourceDemoStore.setBanner(recipeAdminSuccessBanner("Project deleted", fmt.Sprintf("%q was removed from the projects list.", item.Name)))
	http.Redirect(w, r, "/recipes/admin-resource", http.StatusSeeOther)
}

// recipeAdminResourceRefresh completes a remote refresh for the recipe list.
// With HTMX it returns the refresh form fragment and an HX-Trigger raising
// gelium:toast (transient result); without JavaScript it re-renders the full
// page with the list refreshed and a persistent inline toast, exactly like the
// Data table refresh demo.
func (s *server) recipeAdminResourceRefresh(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	isHX := strings.EqualFold(r.Header.Get("HX-Request"), "true")

	view := newRecipeAdminResourceView("", "", "", "", nil, resourceDemoStore.takeBanner())
	applyRequestChrome(r, view)
	view.Refreshed = true

	if isHX {
		trigger, err := toastTriggerJSON("success", "Projects refreshed.")
		if err != nil {
			http.Error(w, "refresh unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Trigger", trigger)
		s.renderRecipeTemplate(w, http.StatusOK, "recipe-admin-resource-refresh-form", view)
		return
	}

	toast := newToast("success", "recipe-ar-refresh-result", "Projects refreshed.")
	view.RefreshToast = &toast
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-admin-resource-list", view)
}

// recipeAdminResourceFormInvalid re-renders a create/edit form with the 422
// validation contract: the X-Gelium-Validation header, the inline alert, the
// validation summary (real links to each field error) and per-field errors.
func (s *server) recipeAdminResourceFormInvalid(w http.ResponseWriter, r *http.Request, mode, id, name, status, date, owner string, errs []recipeFieldError) {
	view := newRecipeAdminResourceFormView(mode, id, name, status, date, owner, errs)
	applyRequestChrome(r, view)
	w.Header().Set("X-Gelium-Validation", "true")
	s.renderRecipeTemplate(w, http.StatusUnprocessableEntity, "recipe-admin-resource-form", view)
}

func (s *server) recipeAdminResourceNotFound(w http.ResponseWriter, title, body string) {
	s.renderErrorPage(w, http.StatusNotFound, title, body, true, "/recipes/admin-resource", "Back to projects", "/recipes/admin-resource")
}

// parseRecipeResourceForm reads and trims the shared form fields.
func (s *server) parseRecipeResourceForm(r *http.Request) (name, status, date, owner string) {
	if err := r.ParseForm(); err != nil {
		return "", "", "", ""
	}
	return strings.TrimSpace(r.FormValue("name")),
		r.FormValue("status"),
		strings.TrimSpace(r.FormValue("date")),
		strings.TrimSpace(r.FormValue("owner"))
}

// recipeFieldError is one field-level validation failure. Href anchors the
// error message element the Text field / Select partials render (the
// validation summary links to it).
type recipeFieldError struct {
	Field   string
	Message string
	Href    string
}

// validateRecipeResourceForm validates the shared resource fields against the
// closed vocabularies and formats. Name and date are required (date must be
// ISO-8601 YYYY-MM-DD); status must be one of dataTableStatuses. Validation
// failures are never reported as toasts — they use the 422 contract.
func validateRecipeResourceForm(name, status, date string) []recipeFieldError {
	var errs []recipeFieldError
	if name == "" {
		errs = append(errs, recipeFieldError{Field: "name", Message: "Enter the project name.", Href: "#recipe-ar-name-error"})
	}
	validStatus := false
	for _, st := range dataTableStatuses {
		if status == st {
			validStatus = true
			break
		}
	}
	if !validStatus {
		errs = append(errs, recipeFieldError{Field: "status", Message: "Choose a status from the list.", Href: "#recipe-ar-status-error"})
	}
	if date == "" {
		errs = append(errs, recipeFieldError{Field: "date", Message: "Enter the project date.", Href: "#recipe-ar-date-error"})
	} else if _, err := time.Parse("2006-01-02", date); err != nil {
		errs = append(errs, recipeFieldError{Field: "date", Message: "Date must use the YYYY-MM-DD format.", Href: "#recipe-ar-date-error"})
	}
	return errs
}

// ----- view builders -----

func recipeAdminMeta(title, description, routePath string) metaView {
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

func recipeAdminSuccessBanner(title, body string) bannerView {
	return bannerView{Tone: "success", Title: title, Body: body}
}

// newRecipeAdminResourceView validates the request against the closed Data
// table vocabularies (reusing dataTableSortKeys and the column/href builders)
// and builds the filtered, sorted, paginated list view from the store.
func newRecipeAdminResourceView(q, sortParam, dir, page string, selection []string, banner *bannerView) *recipeAdminResourceView {
	query := strings.TrimSpace(q)
	sortKey := "name"
	for _, k := range dataTableSortKeys {
		if sortParam == k {
			sortKey = k
		}
	}
	direction := "asc"
	if dir == "desc" {
		direction = "desc"
	}
	pageNum := 1
	if n, err := strconv.Atoi(page); err == nil && n >= 1 {
		pageNum = n
	}

	snapshot := resourceDemoStore.snapshot()
	items := make([]recipeResource, 0, len(snapshot))
	for _, it := range snapshot {
		if query != "" && !strings.Contains(strings.ToLower(it.Name), strings.ToLower(query)) &&
			!strings.Contains(strings.ToLower(it.Status), strings.ToLower(query)) &&
			!strings.Contains(strings.ToLower(it.Owner), strings.ToLower(query)) {
			continue
		}
		items = append(items, it)
	}

	sort.SliceStable(items, func(i, j int) bool {
		a, b := recipeAdminField(items[i], sortKey), recipeAdminField(items[j], sortKey)
		if direction == "desc" {
			return a > b
		}
		return a < b
	})

	total := len(items)
	totalPages := (total + recipeAdminPageSize - 1) / recipeAdminPageSize
	if totalPages < 1 {
		totalPages = 1
	}
	if pageNum > totalPages {
		pageNum = totalPages
	}

	selectedSet := make(map[string]bool)
	selectAll := false
	for _, v := range selection {
		if v == "all" {
			selectAll = true
		}
		selectedSet[v] = true
	}

	start := (pageNum - 1) * recipeAdminPageSize
	end := start + recipeAdminPageSize
	if end > total {
		end = total
	}

	rows := make([]recipeResourceRowView, 0, end-start)
	for _, it := range items[start:end] {
		rows = append(rows, recipeResourceRowView{
			ID:         it.ID,
			Name:       it.Name,
			Status:     it.Status,
			Date:       it.Date,
			Owner:      it.Owner,
			Selected:   selectAll || selectedSet[it.ID],
			EditHref:   "/recipes/admin-resource/" + it.ID + "/edit",
			DeleteHref: "/recipes/admin-resource/" + it.ID + "/delete",
		})
	}

	selectedCount := 0
	if len(selection) > 0 {
		if selectAll {
			selectedCount = total
		} else {
			for _, it := range items {
				if selectedSet[it.ID] {
					selectedCount++
				}
			}
		}
	}

	colspan := 1 + len(dataTableColumns(query, sortKey, direction)) + 1
	empty := recipeAdminEmptyState(query, sortKey, direction)

	pageLinks := make([]dataTablePageView, 0, totalPages)
	for n := 1; n <= totalPages; n++ {
		pageLinks = append(pageLinks, dataTablePageView{
			Num:     n,
			Href:    dataTableHref(query, sortKey, direction, n),
			Current: n == pageNum,
		})
	}

	meta := recipeAdminMeta(
		"Projects · Admin Resource recipe",
		"Manage the projects set: filter, sort, paginate, create, edit and delete server-side with the Admin Resource screen recipe.",
		"/recipes/admin-resource",
	)

	return &recipeAdminResourceView{
		Meta:            meta,
		ThemeClass:      themeClass(""),
		Title:           "Projects",
		Description:     "The Admin Resource screen recipe: a server-rendered resource manager composed from the Data table, form, dialog, banner and toast primitives.",
		NewButton:       buttonView{Label: "New project", Variant: "primary", Href: "/recipes/admin-resource/new"},
		FilterAction:    "/recipes/admin-resource",
		SearchField:     textFieldView{ID: "recipe-ar-q", Label: "Filter", Name: "q", Value: query, Variant: "outlined", Helper: "Filter by name, status or owner."},
		Query:           query,
		Sort:            sortKey,
		Dir:             direction,
		Page:            pageNum,
		Total:           total,
		Pages:           totalPages,
		Caption:         fmt.Sprintf("%d projects · page %d of %d", total, pageNum, totalPages),
		Columns:         dataTableColumns(query, sortKey, direction),
		Rows:            rows,
		PageLinks:       pageLinks,
		PrevHref:        dataTableHref(query, sortKey, direction, pageNum-1),
		NextHref:        dataTableHref(query, sortKey, direction, pageNum+1),
		HasPrev:         pageNum > 1,
		HasNext:         pageNum < totalPages,
		SelectedCount:   selectedCount,
		SelectionNotice: dataTableSelectionNotice(selection),
		Colspan:         colspan,
		EmptyState:      empty,
		Banner:          banner,
	}
}

func recipeAdminField(it recipeResource, key string) string {
	switch key {
	case "status":
		return it.Status
	case "date":
		return it.Date
	default:
		return it.Name
	}
}

// recipeAdminEmptyState builds the table empty row: with a search it invites
// clearing the filters; with no results at all it points to creating the first
// resource.
func recipeAdminEmptyState(query, sortKey, direction string) emptyStateView {
	if query != "" {
		return emptyStateView{
			Title:    "No results",
			Body:     "No projects match your search. Try adjusting the filters.",
			CTA:      true,
			CTAHref:  dataTableHref("", sortKey, direction, 0),
			CTALabel: "Clear search",
			Compact:  true,
		}
	}
	return emptyStateView{
		Title:    "No projects yet",
		Body:     "Create your first project to get started.",
		CTA:      true,
		CTAHref:  "/recipes/admin-resource/new",
		CTALabel: "New project",
		Compact:  true,
	}
}

// newRecipeAdminResourceFormView builds the create/edit form. mode is "create"
// or "edit"; status defaults to the first closed vocabulary value on create.
// errs, when non-empty, produce the inline alert + validation summary + field
// errors that re-render the form with submitted values preserved.
func newRecipeAdminResourceFormView(mode, id, name, status, date, owner string, errs []recipeFieldError) recipeAdminResourceFormView {
	isEdit := mode == "edit"
	action := "/recipes/admin-resource"
	routePath := "/recipes/admin-resource/new"
	heading := "New project"
	intro := "Create a project and add it to the list. Fields marked with an asterisk are required."
	submit := "Create project"
	title := "New project · Admin Resource recipe"
	description := "Create a project in the Admin Resource screen recipe."
	if isEdit {
		action = "/recipes/admin-resource/" + id + "/edit"
		routePath = action
		heading = "Edit " + name
		intro = "Update the project details and save your changes."
		submit = "Save changes"
		title = "Edit " + name + " · Admin Resource recipe"
		description = "Update a project in the Admin Resource screen recipe."
	}
	nameErr, statusErr, dateErr := "", "", ""
	statusInvalid := false
	if len(errs) > 0 {
		for _, e := range errs {
			switch e.Field {
			case "name":
				nameErr = e.Message
			case "status":
				statusErr = e.Message
				statusInvalid = true
			case "date":
				dateErr = e.Message
			}
		}
	}
	if status == "" && !isEdit && !statusInvalid {
		status = dataTableStatuses[0]
	}

	statusOptions := make([]recipeStatusOption, 0, len(dataTableStatuses))
	for _, st := range dataTableStatuses {
		statusOptions = append(statusOptions, recipeStatusOption{Value: st, Label: st})
	}

	return recipeAdminResourceFormView{
		Meta:          recipeAdminMeta(title, description, routePath),
		ThemeClass:    themeClass(""),
		Heading:       heading,
		Intro:         intro,
		Action:        action,
		CancelHref:    "/recipes/admin-resource",
		SubmitLabel:   submit,
		NameField:     textFieldView{ID: "recipe-ar-name", Label: "Name", Name: "name", Value: name, Variant: "outlined", Helper: "Display name of the project.", Error: nameErr},
		StatusValue:   status,
		StatusOptions: statusOptions,
		StatusError:   statusErr,
		DateField:     textFieldView{ID: "recipe-ar-date", Label: "Date", Name: "date", Value: date, Variant: "outlined", Helper: "Use the YYYY-MM-DD format, e.g. 2026-08-10.", Error: dateErr},
		OwnerField:    textFieldView{ID: "recipe-ar-owner", Label: "Owner", Name: "owner", Value: owner, Variant: "outlined", Helper: "Optional. Who is responsible for this project."},
		InlineAlert:   recipeAdminFormInlineAlert(errs),
		Validation:    recipeAdminFormValidation(errs),
	}
}

func recipeAdminFormInlineAlert(errs []recipeFieldError) *inlineAlertView {
	if len(errs) == 0 {
		return nil
	}
	return &inlineAlertView{
		Tone:  "error",
		Title: "The project could not be saved",
		Body:  "Check the highlighted fields below and try again.",
	}
}

func recipeAdminFormValidation(errs []recipeFieldError) *validationSummaryView {
	if len(errs) == 0 {
		return nil
	}
	items := make([]validationSummaryItemView, 0, len(errs))
	for _, e := range errs {
		items = append(items, validationSummaryItemView{Href: e.Href, Message: e.Message})
	}
	return &validationSummaryView{
		Title:        "Please fix the following issues",
		HeadingLevel: 2,
		Items:        items,
	}
}

// renderRecipeTemplate executes one recipe template into the response with the
// given status.
// applyRequestChrome applies the document-root theme AND scheme selection
// (?theme=/?scheme= query, Phase H) to a recipe view. Recipe templates are
// standalone (not the docs layout), so the query middleware alone cannot
// reach them; handlers call this right after building their view. Dark
// appends theme-dark + data-theme="dark"; light emits data-theme="light".
func applyRequestChrome(r *http.Request, view interface{}) {
	theme := themeFromRequest(r)
	scheme := schemeFromRequest(r)
	switch v := view.(type) {
	case *recipeAdminResourceView:
		applyChromeToView(theme, &v.ThemeClass, &v.DataTheme, scheme)
	case *recipeAdminResourceFormView:
		applyChromeToView(theme, &v.ThemeClass, &v.DataTheme, scheme)
	case *recipeAdminResourceConfirmView:
		applyChromeToView(theme, &v.ThemeClass, &v.DataTheme, scheme)
	case *recipeOpsQueueView:
		applyChromeToView(theme, &v.ThemeClass, &v.DataTheme, scheme)
	case *recipeFeedView:
		applyChromeToView(theme, &v.ThemeClass, &v.DataTheme, scheme)
	}
}

// applyChromeToView mirrors applyDocumentRootScheme for the standalone recipe
// views: a valid ?theme= overrides ThemeClass, and the scheme mutates the
// document-root classes/data-theme exactly like the docs layout.
func applyChromeToView(theme string, themeClass *string, dataTheme *string, scheme string) {
	if theme != "" {
		*themeClass = theme
	}
	switch normalizeScheme(scheme) {
	case "dark":
		if !strings.Contains(*themeClass, "theme-dark") {
			*themeClass = strings.TrimSpace(*themeClass + " theme-dark")
		}
		*dataTheme = "dark"
	case "light":
		*dataTheme = "light"
	default:
		*dataTheme = ""
	}
}

func (s *server) renderRecipeTemplate(w http.ResponseWriter, status int, name string, data interface{}) {
	var page bytes.Buffer
	if err := s.templates.ExecuteTemplate(&page, name, data); err != nil {
		http.Error(w, "recipe unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(page.Bytes())
}
