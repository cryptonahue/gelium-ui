package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// postRecipeAdminResource posts a create/update form and returns the recorder.
func postRecipeAdminResource(t *testing.T, path string, form url.Values) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)
	return res
}

// TestRecipeAdminResourceListRendersDataTableWithFilterSortPage proves the list
// screen reuses the Data table server-driven pattern: stable GET params for
// filter/sort/page, aria-sort on the active column, aria-current on the current
// page and server-side pagination.
func TestRecipeAdminResourceListRendersDataTableWithFilterSortPage(t *testing.T) {
	resetRecipeResourceStore()

	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1 class="recipe-ar-title">Projects</h1>`,
		`id="recipe-ar-q" type="search"`,
		`<table class="ui-data-table-table">`,
		`<caption class="ui-data-table-caption">12 projects · page 1 of 3</caption>`,
		`aria-sort="ascending"`,
		`<th scope="col" class="ui-data-table-cell">Actions</th>`,
		`href="/recipes/admin-resource/alpha/edit"`,
		`href="/recipes/admin-resource/alpha/delete"`,
		`name="selection" value="all" aria-label="Select all rows"`,
		`<span class="ui-data-table-page ui-data-table-page--current" aria-current="page">1</span>`,
		`id="resource-panel"`,
		`id="gelium-toast-region"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("list does not contain contract %q", contract)
		}
	}

	// Filter: q matches name, status and owner.
	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource?q=alpha", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("filter status = %d, want %d", res.Code, http.StatusOK)
	}
	body = res.Body.String()
	for _, contract := range []string{`Alpha release`, `1 projects · page 1 of 1`} {
		if !strings.Contains(body, contract) {
			t.Errorf("filtered list does not contain contract %q", contract)
		}
	}
	if strings.Contains(body, "Beta rollout") {
		t.Error("filtered list must not contain a non-matching row")
	}

	// Status filter: the closed GET vocabulary filters server-side and remains
	// selected in the native control while sort links preserve the filter.
	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource?status=Pending", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status filter status = %d, want %d", res.Code, http.StatusOK)
	}
	body = res.Body.String()
	for _, contract := range []string{
		`name="status"`,
		`<option value="Pending" selected>Pending</option>`,
		`4 projects · page 1 of 1`,
		`Beta rollout`,
		`href="?dir=desc&amp;sort=name&amp;status=Pending"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("status filter does not contain contract %q", contract)
		}
	}
	if strings.Contains(body, "Alpha release") {
		t.Error("status filter must drop a non-matching row")
	}

	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource?q=alpha&status=Pending", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("combined empty filter status = %d, want %d", res.Code, http.StatusOK)
	}
	body = res.Body.String()
	for _, contract := range []string{`No results`, `Clear filters`, `href="?dir=asc&amp;sort=name"`} {
		if !strings.Contains(body, contract) {
			t.Errorf("combined empty filter does not contain contract %q", contract)
		}
	}

	// Sort by date descending: the newest project leads the first page.
	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource?sort=date&dir=desc", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("sort status = %d, want %d", res.Code, http.StatusOK)
	}
	body = res.Body.String()
	if !strings.Contains(body, `aria-sort="descending"`) {
		t.Error("date desc sort must mark the active column with aria-sort=descending")
	}
	if idxMu, idxLambda := strings.Index(body, "Mu navigation"), strings.Index(body, "Lambda mailer"); idxMu < 0 || idxLambda < 0 || idxMu > idxLambda {
		t.Errorf("date desc sort must lead with Mu navigation before Lambda mailer (Mu=%d Lambda=%d)", idxMu, idxLambda)
	}

	// Page 2 with the name ascending order shows rows 6-10.
	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource?sort=name&dir=asc&page=2", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("page status = %d, want %d", res.Code, http.StatusOK)
	}
	body = res.Body.String()
	for _, contract := range []string{
		`12 projects · page 2 of 3`,
		`<span class="ui-data-table-page ui-data-table-page--current" aria-current="page">2</span>`,
		`Gamma refactor`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("page 2 does not contain contract %q", contract)
		}
	}
	if strings.Contains(body, "Alpha release") {
		t.Error("page 2 must not contain a row from page 1")
	}
}

// TestRecipeAdminResourceListHXFragment proves the HX-Request bifurcation: HTMX
// swaps only the #resource-panel fragment, not the whole page.
func TestRecipeAdminResourceListHXFragment(t *testing.T) {
	resetRecipeResourceStore()

	req := httptest.NewRequest(http.MethodGet, "/recipes/admin-resource?q=delta", nil)
	req.Header.Set("HX-Request", "true")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Contains(body, "<html") || strings.Contains(body, "<!doctype html>") {
		t.Error("HX request must return the panel fragment, not a full document")
	}
	for _, contract := range []string{`id="resource-panel"`, `Delta docs`, `hx-target="#resource-panel"`} {
		if !strings.Contains(body, contract) {
			t.Errorf("panel fragment does not contain contract %q", contract)
		}
	}
}

func TestRecipeAdminResourceDetailRendersReadOnlyRecord(t *testing.T) {
	resetRecipeResourceStore()
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource/alpha", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<article class="recipe-ar-detail">`,
		`<h1 class="recipe-ar-title">Alpha release</h1>`,
		`<dl class="recipe-ar-detail-fields">`,
		`<dt>Name</dt><dd>Alpha release</dd>`,
		`<dt>Status</dt><dd><span class="ui-badge ui-badge-large ui-badge--success">Active</span>`,
		`href="/recipes/admin-resource/alpha/edit"`,
		`href="/recipes/admin-resource/alpha/delete"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("detail does not contain contract %q", contract)
		}
	}
}

func TestRecipeAdminResourceDetailNotFoundAndGETOnly(t *testing.T) {
	resetRecipeResourceStore()
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource/nope", nil))
	if res.Code != http.StatusNotFound || !strings.Contains(res.Body.String(), `class="ui-error-state"`) {
		t.Fatalf("unknown detail status/body = %d/%q, want 404 error state", res.Code, res.Body.String())
	}
	res = postRecipeAdminResource(t, "/recipes/admin-resource/alpha", url.Values{})
	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST detail status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

// TestRecipeAdminResourceEmptyState proves the table empty row uses the shared
// empty-state primitive: a search with no matches offers a real "Clear search"
// CTA and the empty dataset offers "New project".
func TestRecipeAdminResourceEmptyState(t *testing.T) {
	resetRecipeResourceStore()

	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource?q=zzzznothing", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{`No results`, `Try adjusting the filters.`, `Clear search`} {
		if !strings.Contains(body, contract) {
			t.Errorf("empty state does not contain contract %q", contract)
		}
	}

	resetRecipeResourceStore()
	for _, it := range resourceDemoStore.snapshot() {
		resourceDemoStore.delete(it.ID)
	}
	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource", nil))
	body = res.Body.String()
	for _, contract := range []string{`No projects yet`, `Create your first project to get started.`, `href="/recipes/admin-resource/new"`} {
		if !strings.Contains(body, contract) {
			t.Errorf("empty dataset state does not contain contract %q", contract)
		}
	}
}

// TestRecipeAdminResourceSelection proves selection round-trips as stable GET
// params: the matching rows render checked server-side and the notice reports
// the count, with the select-all box unchecked when not everything is selected.
func TestRecipeAdminResourceSelection(t *testing.T) {
	resetRecipeResourceStore()

	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource?selection=alpha&selection=beta", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`value="alpha" checked`,
		`value="beta" checked`,
		`2 rows selected.`,
		`href="?dir=desc&amp;selection=alpha&amp;selection=beta&amp;sort=name"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("selection state does not contain contract %q", contract)
		}
	}

	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource?selection=alpha&selection=alpha&selection=unknown", nil))
	body = res.Body.String()
	if !strings.Contains(body, `1 row selected.`) {
		t.Error("selection notice must deduplicate valid IDs and ignore unknown IDs")
	}

	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource?selection=all", nil))
	body = res.Body.String()
	if !strings.Contains(body, `href="?dir=asc&amp;page=2&amp;selection=all&amp;sort=name"`) {
		t.Error("select-all scope must survive pagination links")
	}
}

// TestRecipeAdminResourceNewForm proves the create form is a full server-rendered
// page with the shared form components and a real Cancel link.
func TestRecipeAdminResourceNewForm(t *testing.T) {
	resetRecipeResourceStore()

	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource/new", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1 class="recipe-ar-title">New project</h1>`,
		`method="post" action="/recipes/admin-resource" novalidate`,
		`id="recipe-ar-name"`,
		`id="recipe-ar-status"`,
		`id="recipe-ar-date"`,
		`id="recipe-ar-owner"`,
		`>Create project</button>`,
		`href="/recipes/admin-resource"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("new form does not contain contract %q", contract)
		}
	}
}

// TestRecipeAdminResourceCreateValidation422 proves the create failure uses the
// 422 contract: the X-Gelium-Validation header, the inline alert, the validation
// summary linking to each field error and the preserved submitted values.
func TestRecipeAdminResourceCreateValidation422(t *testing.T) {
	resetRecipeResourceStore()

	res := postRecipeAdminResource(t, "/recipes/admin-resource", url.Values{
		"name":   {"  "},
		"status": {"On hold"},
		"date":   {"not-a-date"},
		"owner":  {"Alicia R."},
	})
	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	if res.Header().Get("X-Gelium-Validation") != "true" {
		t.Errorf("X-Gelium-Validation = %q, want true", res.Header().Get("X-Gelium-Validation"))
	}
	body := res.Body.String()
	for _, contract := range []string{
		`class="ui-validation-summary" role="alert"`,
		`href="#recipe-ar-name-error">Enter the project name.</a>`,
		`href="#recipe-ar-status-error">Choose a status from the list.</a>`,
		`href="#recipe-ar-date-error">Date must use the YYYY-MM-DD format.</a>`,
		`class="ui-inline-alert ui-inline-alert--error"`,
		`id="recipe-ar-name-error" role="alert"`,
		`id="recipe-ar-status-error" role="alert"`,
		`id="recipe-ar-date-error" role="alert"`,
		`aria-invalid="true" aria-describedby="recipe-ar-status-error"`,
		`value="Alicia R."`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("422 form does not contain contract %q", contract)
		}
	}
}

// TestRecipeAdminResourceCreateSuccess303 proves the create success follows
// POST+303 and the following list render shows the persistent success banner
// (never a toast), which is consumed on the next render.
func TestRecipeAdminResourceCreateSuccess303(t *testing.T) {
	resetRecipeResourceStore()

	res := postRecipeAdminResource(t, "/recipes/admin-resource", url.Values{
		"name":   {"Zulu launch"},
		"status": {"Active"},
		"date":   {"2026-12-20"},
		"owner":  {"Alicia R."},
	})
	if res.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusSeeOther)
	}
	if loc := res.Header().Get("Location"); loc != "/recipes/admin-resource" {
		t.Errorf("Location = %q, want /recipes/admin-resource", loc)
	}

	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`class="ui-banner ui-banner--success"`,
		`class="ui-banner-title">Project created</p>`,
		`Zulu launch`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("create result does not contain contract %q", contract)
		}
	}

	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource", nil))
	if strings.Contains(res.Body.String(), "Project created") {
		t.Error("the flash success banner must be consumed after one render")
	}
}

// TestRecipeAdminResourceEditForm proves the edit form is pre-populated with the
// stored resource and posts to the per-item edit action.
func TestRecipeAdminResourceEditForm(t *testing.T) {
	resetRecipeResourceStore()

	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource/alpha/edit", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1 class="recipe-ar-title">Edit Alpha release</h1>`,
		`method="post" action="/recipes/admin-resource/alpha/edit" novalidate`,
		`value="Alpha release"`,
		`<option value="Active" selected>Active</option>`,
		`value="2026-01-12"`,
		`value="Alicia R."`,
		`>Save changes</button>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("edit form does not contain contract %q", contract)
		}
	}
}

// TestRecipeAdminResourceEditValidation422 proves the edit failure re-renders
// the form with the 422 contract while preserving the submitted values.
func TestRecipeAdminResourceEditValidation422(t *testing.T) {
	resetRecipeResourceStore()

	res := postRecipeAdminResource(t, "/recipes/admin-resource/alpha/edit", url.Values{
		"name":   {"  "},
		"status": {"Active"},
		"date":   {"2026-01-12"},
		"owner":  {"Alicia R."},
	})
	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	if res.Header().Get("X-Gelium-Validation") != "true" {
		t.Errorf("X-Gelium-Validation = %q, want true", res.Header().Get("X-Gelium-Validation"))
	}
	body := res.Body.String()
	for _, contract := range []string{
		`class="ui-validation-summary" role="alert"`,
		`href="#recipe-ar-name-error">Enter the project name.</a>`,
		`value="Alicia R."`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("edit 422 form does not contain contract %q", contract)
		}
	}
}

// TestRecipeAdminResourceEditSuccess303 proves the edit success redirects and
// shows the persistent "Project updated" banner on the next render.
func TestRecipeAdminResourceEditSuccess303(t *testing.T) {
	resetRecipeResourceStore()

	res := postRecipeAdminResource(t, "/recipes/admin-resource/alpha/edit", url.Values{
		"name":   {"Alpha refresh"},
		"status": {"Pending"},
		"date":   {"2026-01-12"},
		"owner":  {"Alicia R."},
	})
	if res.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusSeeOther)
	}

	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource?q=alpha", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{`Project updated`, `Alpha refresh`} {
		if !strings.Contains(body, contract) {
			t.Errorf("edit result does not contain contract %q", contract)
		}
	}
}

// TestRecipeAdminResourceDeleteConfirm proves the delete confirmation is a
// Dialog page variant: a real native <dialog open> with the destructive action
// as a real POST form and a Cancel link.
func TestRecipeAdminResourceDeleteConfirm(t *testing.T) {
	resetRecipeResourceStore()

	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource/alpha/delete", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<dialog class="ui-dialog" open`,
		`id="recipe-ar-confirm-title">Delete Alpha release?</h2>`,
		`You are about to permanently delete Alpha release (owned by Alicia R.). This action cannot be undone.`,
		`method="post" action="/recipes/admin-resource/alpha/delete"`,
		`>Delete Alpha release</button>`,
		`href="/recipes/admin-resource"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("delete confirm does not contain contract %q", contract)
		}
	}
}

// TestRecipeAdminResourceDeleteSuccess303 proves the delete mutation follows
// POST+303, shows the persistent success banner and removes the row.
func TestRecipeAdminResourceDeleteSuccess303(t *testing.T) {
	resetRecipeResourceStore()

	res := postRecipeAdminResource(t, "/recipes/admin-resource/alpha/delete", url.Values{})
	if res.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusSeeOther)
	}

	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource?q=alpha", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{`class="ui-banner ui-banner--success"`, `Project deleted`, `No results`} {
		if !strings.Contains(body, contract) {
			t.Errorf("delete result does not contain contract %q", contract)
		}
	}
}

func TestRecipeAdminResourceBulkDeleteConfirm(t *testing.T) {
	resetRecipeResourceStore()
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource/bulk-delete?selection=alpha&selection=beta", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`id="recipe-ar-bulk-confirm-title">Delete 2 projects?`,
		`action="/recipes/admin-resource/bulk-delete"`,
		`name="selection" value="alpha"`,
		`name="selection" value="beta"`,
		`This action cannot be undone.`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("bulk confirmation does not contain contract %q", contract)
		}
	}
}

func TestRecipeAdminResourceBulkDeleteSuccess303(t *testing.T) {
	resetRecipeResourceStore()
	res := postRecipeAdminResource(t, "/recipes/admin-resource/bulk-delete", url.Values{
		"selection": {"alpha", "beta"},
	})
	if res.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusSeeOther)
	}
	if got := res.Header().Get("Location"); got != "/recipes/admin-resource" {
		t.Fatalf("Location = %q, want /recipes/admin-resource", got)
	}
	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource", nil))
	body := res.Body.String()
	for _, contract := range []string{`2 projects deleted`, `class="ui-banner ui-banner--success"`} {
		if !strings.Contains(body, contract) {
			t.Errorf("bulk delete result does not contain contract %q", contract)
		}
	}
	if strings.Contains(body, "Alpha release") || strings.Contains(body, "Beta rollout") {
		t.Error("bulk-deleted projects must not remain in the list")
	}
}

// TestRecipeAdminResourceNotFound proves unknown resource IDs render the shared
// error-state page with the real 404 status on both edit and delete.
func TestRecipeAdminResourceAuthorizationIsConsumerOwnedAndRevalidated(t *testing.T) {
	resetRecipeResourceStore()
	denyBeta := func(_ *http.Request, action string, item *recipeResource) bool {
		if action != recipeAdminDeleteAction {
			return true
		}
		return item == nil || item.ID != "beta"
	}
	h := newWithRecipeAdminAuthorizer(denyBeta)

	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, "/recipes/admin-resource/alpha/delete") {
		t.Error("authorized row should render its delete action")
	}
	if strings.Contains(body, "/recipes/admin-resource/beta/delete") {
		t.Error("unauthorized row must not render a delete action")
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource/beta/delete", nil))
	if res.Code != http.StatusForbidden {
		t.Fatalf("unauthorized confirmation status = %d, want %d", res.Code, http.StatusForbidden)
	}

	form := url.Values{"selection": {"alpha", "beta"}}
	req := httptest.NewRequest(http.MethodPost, "/recipes/admin-resource/bulk-delete", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res = httptest.NewRecorder()
	h.ServeHTTP(res, req)
	if res.Code != http.StatusSeeOther {
		t.Fatalf("mixed bulk delete status = %d, want %d", res.Code, http.StatusSeeOther)
	}
	if _, ok := resourceDemoStore.get("alpha"); ok {
		t.Error("authorized selected row should be deleted")
	}
	if _, ok := resourceDemoStore.get("beta"); !ok {
		t.Error("unauthorized selected row must survive")
	}
	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource", nil))
	if !strings.Contains(res.Body.String(), "not authorized and were left unchanged") {
		t.Error("mixed bulk delete should report the unauthorized selection")
	}
}

func TestRecipeAdminResourceNotFound(t *testing.T) {
	resetRecipeResourceStore()

	for _, path := range []string{"/recipes/admin-resource/nope/edit", "/recipes/admin-resource/nope/delete"} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, path, nil))
		if res.Code != http.StatusNotFound {
			t.Errorf("GET %s status = %d, want %d", path, res.Code, http.StatusNotFound)
		}
		if body := res.Body.String(); !strings.Contains(body, `class="ui-error-state"`) {
			t.Errorf("GET %s must render the error-state page", path)
		}
	}
}

// TestRecipeAdminResourceRefreshPostOnly proves the refresh action is POST-only
// (a GET answers 405 with Allow: POST), that the no-JS refresh re-renders the
// list with a persistent inline toast + progress, and that the HTMX refresh
// returns the fragment plus an HX-Trigger gelium:toast.
func TestRecipeAdminResourceRefreshPostOnly(t *testing.T) {
	resetRecipeResourceStore()

	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource/refresh", nil))
	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET refresh status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
	if got := res.Header().Get("Allow"); got != "POST" {
		t.Errorf("Allow = %q, want POST", got)
	}

	res = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/recipes/admin-resource/refresh", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	New().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("POST refresh status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`Projects refreshed.`,
		`<progress value="100" max="100"></progress>`,
		`class="ui-toast ui-toast-success"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("no-JS refresh does not contain contract %q", contract)
		}
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/recipes/admin-resource/refresh", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	New().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("HX refresh status = %d, want %d", res.Code, http.StatusOK)
	}
	if trigger := res.Header().Get("HX-Trigger"); !strings.Contains(trigger, "gelium:toast") || !strings.Contains(trigger, "Projects refreshed.") {
		t.Errorf("HX-Trigger = %q, want gelium:toast payload", trigger)
	}
	body = res.Body.String()
	if strings.Contains(body, "<html") {
		t.Error("HX refresh must return the refresh form fragment, not a full document")
	}
	if !strings.Contains(body, `class="recipe-ar-refresh"`) {
		t.Error("HX refresh must return the refresh form fragment")
	}
}

// TestRecipeAdminResourceNoindexProvesRecipeSurfacesAreNeverIndexed proves every
// recipe route emits robots noindex, nofollow plus a canonical and a per-route
// description, in the recipe's own document head.
func TestRecipeAdminResourceNoindexProvesRecipeSurfacesAreNeverIndexed(t *testing.T) {
	resetRecipeResourceStore()

	cases := []struct {
		path  string
		title string
	}{
		{path: "/recipes/admin-resource", title: "Projects · Admin Resource recipe"},
		{path: "/recipes/admin-resource/new", title: "New project · Admin Resource recipe"},
		{path: "/recipes/admin-resource/alpha/edit", title: "Edit Alpha release · Admin Resource recipe"},
		{path: "/recipes/admin-resource/alpha/delete", title: "Delete Alpha release · Admin Resource recipe"},
	}
	for _, tc := range cases {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if res.Code != http.StatusOK {
			t.Fatalf("GET %s status = %d, want %d", tc.path, res.Code, http.StatusOK)
		}
		body := res.Body.String()
		for _, contract := range []string{
			`lang="en"`,
			`<meta name="robots" content="noindex, nofollow">`,
			`<link rel="canonical" href="https://gelium-ui.example` + tc.path + `">`,
			`<title>` + tc.title + `</title>`,
			`<meta name="description"`,
		} {
			if !strings.Contains(body, contract) {
				t.Errorf("GET %s does not contain contract %q", tc.path, contract)
			}
		}
	}
}
