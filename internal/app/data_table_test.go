package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDataTableDocsRouteDogfoodsNativeTableSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/data-table", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("data table docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Data table</h1>`,
		`class="ui-data-table"`,
		`<table`,
		`<caption`,
		`<thead`,
		`<tbody`,
		`<th scope="col"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("data table docs are missing %q", contract)
		}
	}
}

func TestDataTableDocsRouteUsesRealGETLinksForSortAndPagination(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/data-table", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("data table docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`class="ui-data-table-sort ui-data-table-sort--active" href="?sort=name&amp;dir=desc"`,
		`class="ui-data-table-sort" href="?sort=status&amp;dir=asc"`,
		`class="ui-data-table-sort" href="?sort=date&amp;dir=asc"`,
		`<form class="data-table-demo-filter" method="get" action="/components/data-table#data-table-demo"`,
		`name="q"`,
		`class="ui-data-table-page" href="?sort=name&amp;dir=asc&amp;page=2"`,
		`aria-sort="ascending"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("data table docs are missing %q", contract)
		}
	}
}

func TestDataTableDocsRouteSortsRowsServerSide(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/data-table?sort=status&dir=desc", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `aria-sort="descending"`) {
		t.Error("active sort column must carry aria-sort=descending")
	}
	if !strings.Contains(body, `class="ui-data-table-sort ui-data-table-sort--active" href="?sort=status&amp;dir=asc"`) {
		t.Error("active sort column link must toggle to asc")
	}
	// status desc: Pending > Done > Active alphabetically (P > D > A). The
	// first page holds the four Pending rows before the first Done row.
	firstPending := strings.Index(body, ">Pending</td>")
	firstDone := strings.Index(body, ">Done</td>")
	if firstPending == -1 || firstDone == -1 {
		t.Fatal("status cells must render with closed vocabulary values")
	}
	if !(firstPending < firstDone) {
		t.Errorf("status desc must order Pending before Done, got %d %d", firstPending, firstDone)
	}
	if !strings.Contains(body, ">Lambda mailer</td>") {
		t.Error("status desc must surface a Pending row on the first page")
	}
}

func TestDataTableDocsRoutePaginatesRowsServerSide(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/data-table?page=2", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `aria-current="page">2</span>`) {
		t.Error("current page must be marked with aria-current")
	}
	if !strings.Contains(body, `class="ui-data-table-page" href="?sort=name&amp;dir=asc&amp;page=1"`) {
		t.Error("page 2 must link back to page 1")
	}
	if !strings.Contains(body, `class="ui-data-table-page" href="?sort=name&amp;dir=asc&amp;page=3"`) {
		t.Error("page 2 must link forward to page 3")
	}
	if !strings.Contains(body, ">Kappa layout</td>") {
		t.Error("page 2 must render rows from the second page")
	}
}

func TestDataTableDocsRouteFiltersRowsServerSide(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/data-table?q=delta", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, ">Delta docs</td>") {
		t.Error("filter q=delta must keep the Delta docs row")
	}
	if strings.Contains(body, ">Alpha release</td>") {
		t.Error("filter q=delta must drop the Alpha release row")
	}
	if !strings.Contains(body, `value="delta"`) {
		t.Error("filter input must preserve the submitted query")
	}
}

func TestDataTableDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/data-table", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST data table status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestDataTableHXRequestSwapsOnlyTablePanel(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/components/data-table?sort=name&dir=desc", nil)
	req.Header.Set("HX-Request", "true")
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("HX status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `<div class="ui-data-table" id="data-table-panel">`) {
		t.Error("HX response must be the table panel fragment")
	}
	if strings.Contains(body, "<html") {
		t.Error("HX response must not contain a full document")
	}
	if strings.Contains(body, `<form class="data-table-demo-filter"`) {
		t.Error("HX response must not include the surrounding filter form")
	}
}

func TestDataTableNoJSSortLinkIsAFullPageGETFallback(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/data-table?sort=name&dir=desc", nil))
	body := res.Body.String()
	if !strings.Contains(body, "<html") {
		t.Error("no-JS sort must return a full document")
	}
	if !strings.Contains(body, `aria-sort="descending"`) {
		t.Error("full-page sort must mark the active column with aria-sort")
	}
}

func TestDataTableSelectionUsesNativeCheckboxesInARealForm(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/data-table", nil))
	body := res.Body.String()
	for _, contract := range []string{
		`<form class="data-table-demo-select" method="get" action="/components/data-table#data-table-demo">`,
		`<input type="checkbox" name="selection" value="all"`,
		`<input type="checkbox" name="selection" value="alpha"`,
		`class="data-table-demo-submit"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("selection demo is missing %q", contract)
		}
	}
}

func TestDataTableSelectionRoundTripMarksRowsChecked(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/data-table?selection=alpha&selection=beta", nil))
	body := res.Body.String()
	if !strings.Contains(body, `<input type="checkbox" name="selection" value="alpha" checked`) {
		t.Error("selected row must render checked server-side")
	}
	if !strings.Contains(body, `<input type="checkbox" name="selection" value="beta" checked`) {
		t.Error("second selected row must render checked server-side")
	}
	if !strings.Contains(body, "2 rows selected.") {
		t.Error("selection round-trip must show a status notice")
	}
	if !strings.Contains(body, `<tr class="ui-data-table-row">`) {
		t.Error("selected row markup must stay a native table row")
	}
}

func TestDataTableRefreshNoJSRendersPersistentToastAndProgress(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/data-table/refresh", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	New().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `class="ui-progress data-table-demo-progress"`) {
		t.Error("no-JS refresh must render the .ui-progress primitive")
	}
	if !strings.Contains(body, `<div class="ui-toast ui-toast-success" id="data-table-refresh-result"`) {
		t.Error("no-JS refresh must render a persistent inline toast")
	}
	if got := res.Header().Get("HX-Trigger"); got != "" {
		t.Errorf("no-JS refresh must not rely on HX-Trigger, got %q", got)
	}
}

func TestDataTableRefreshHXReturnsFragmentAndToastTrigger(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/data-table/refresh", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	New().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if got := res.Header().Get("HX-Trigger"); !strings.Contains(got, "loom:toast") {
		t.Errorf("HX refresh must raise loom:toast, got %q", got)
	}
	body := res.Body.String()
	if !strings.Contains(body, `<form class="data-table-demo-refresh"`) {
		t.Error("HX response must be the refresh form fragment")
	}
	if strings.Contains(body, "<html") {
		t.Error("HX response must not contain a full document")
	}
	if strings.Contains(body, `id="data-table-refresh-result"`) {
		t.Error("HX refresh must rely on the live region, not prerender an inline toast")
	}
}

func TestDataTableDocsRouteSanitizesClosedVocabularies(t *testing.T) {
	for _, u := range []string{
		"/components/data-table?sort=bogus&dir=sideways",
		"/components/data-table?page=0",
		"/components/data-table?page=999",
		"/components/data-table?sort=%22%3E%3Cscript%3Ealert(1)%3C%2Fscript%3E",
	} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, u, nil))
		if res.Code != http.StatusOK {
			t.Fatalf("status = %d for %s, want %d", res.Code, u, http.StatusOK)
		}
		if strings.Contains(res.Body.String(), "bogus") {
			t.Errorf("invalid sort key must not survive sanitization for %s", u)
		}
		if strings.Contains(res.Body.String(), "sideways") {
			t.Errorf("invalid dir must not survive sanitization for %s", u)
		}
		if strings.Contains(res.Body.String(), "<script>") {
			t.Errorf("injected sort key must be escaped for %s", u)
		}
	}
}

func TestDataTableDocsRouteEscapesFilterQuery(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, `/components/data-table?q=<script>alert(1)</script>`, nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Contains(body, `<script>alert(1)`) {
		t.Error("filter query must be HTML-escaped in the input value and links")
	}
	if !strings.Contains(body, `&lt;script&gt;`) {
		t.Error("filter query must render escaped in the input value")
	}
}

func TestDataTablePageLinksPreserveSortAndFilter(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/data-table?q=a&sort=date&dir=desc&page=1", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `class="ui-data-table-page" href="?q=a&amp;sort=date&amp;dir=desc&amp;page=2"`) {
		t.Error("page 2 link must preserve the active filter and sort")
	}
	if !strings.Contains(body, `class="ui-data-table-sort" href="?q=a&amp;sort=name&amp;dir=asc"`) {
		t.Error("switching sort must preserve the active filter")
	}
}

func TestDataTableEmptyStateRendersMessageAndCTA(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/data-table?q=zzz", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"0 rows · page 1 of 1",
		`class="ui-empty-state ui-empty-state--compact"`,
		`colspan="4"`,
		`<a class="ui-button" href="?sort=name&amp;dir=asc">Clear search</a>`,
		">No results</p>",
		">Try adjusting your search or filters.</p>",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("empty state is missing %q", contract)
		}
	}
	if strings.Contains(body, `<tr class="ui-data-table-row">`) {
		t.Error("empty state must not render any data rows")
	}
}

func TestDataTableEmptyStateSelectAllDisabledOrAbsent(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/data-table?q=zzz", nil))
	body := res.Body.String()
	if strings.Contains(body, `<input type="checkbox" name="selection" value="all"`) {
		t.Error("select-all checkbox must not render with zero rows")
	}
	if strings.Contains(body, `class="data-table-demo-submit"`) {
		t.Error("submit selection button must not render with zero rows")
	}
}

func TestDataTableEmptyStateHXFragmentIncludesEmptyRow(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/components/data-table?q=zzz", nil)
	req.Header.Set("HX-Request", "true")
	New().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Contains(body, "<html") {
		t.Error("HX empty state response must be a fragment, not a full document")
	}
	for _, contract := range []string{
		`<div class="ui-data-table" id="data-table-panel">`,
		`class="ui-empty-state ui-empty-state--compact"`,
		`colspan="4"`,
		">No results</p>",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("HX empty state fragment is missing %q", contract)
		}
	}
}

func TestDataTableEmptyStateAllSelectionCountsZero(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/data-table?q=zzz&selection=all", nil))
	body := res.Body.String()
	if strings.Contains(body, "All rows selected.") {
		t.Error("empty state must not claim all rows are selected")
	}
	if strings.Contains(body, `<input type="checkbox" name="selection" value="all" checked`) {
		t.Error("select-all must not be checked with zero rows even after selection=all")
	}
}

func TestDataTableEmptyStateEscapesQuery(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, `/components/data-table?q=<script>alert(1)</script>`, nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Contains(body, `<script>alert(1)`) {
		t.Error("injected query must be HTML-escaped in the empty state")
	}
	if !strings.Contains(body, `class="ui-empty-state ui-empty-state--compact"`) {
		t.Error("an unmatched query must still render the empty state")
	}
}
