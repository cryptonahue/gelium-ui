package app

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// Data table: a server-rendered Material 3 data table reimplemented over native
// HTML table semantics (<table>/<thead>/<tbody>/<th scope>/<caption>). Sort,
// filter and pagination are real GET requests handled by the server: column
// headers are <a href> links carrying ?sort=&dir=, the filter is a real GET
// form carrying ?q=, and page numbers are real GET links carrying ?page=. The
// no-JS flow is a normal full-page GET; the same links optionally carry
// hx-get so HTMX swaps only the table panel. Row selection uses native
// checkboxes in a real form, so multi-select submits without JavaScript.

// dataTableSortKeys is the closed vocabulary of sortable columns.
var dataTableSortKeys = []string{"name", "status", "date"}

// dataTablePageSize controls the server-side pagination slice of the demo.
const dataTablePageSize = 5

// dataTableStatuses is the closed vocabulary of the status column.
var dataTableStatuses = []string{"Active", "Pending", "Done"}

// dataTableItem is one row in the server-side demo dataset.
type dataTableItem struct {
	ID     string
	Name   string
	Status string
	Date   string
}

// dataTableItems is the small closed-vocabulary demo dataset. Dates are
// ISO-8601 so string ordering equals chronological ordering.
var dataTableItems = []dataTableItem{
	{ID: "alpha", Name: "Alpha release", Status: "Active", Date: "2026-01-12"},
	{ID: "beta", Name: "Beta rollout", Status: "Pending", Date: "2026-02-03"},
	{ID: "gamma", Name: "Gamma refactor", Status: "Done", Date: "2026-03-15"},
	{ID: "delta", Name: "Delta docs", Status: "Active", Date: "2026-04-28"},
	{ID: "epsilon", Name: "Epsilon migration", Status: "Pending", Date: "2026-05-06"},
	{ID: "zeta", Name: "Zeta dashboard", Status: "Done", Date: "2026-06-02"},
	{ID: "eta", Name: "Eta audit", Status: "Active", Date: "2026-07-19"},
	{ID: "theta", Name: "Theta backup", Status: "Pending", Date: "2026-08-07"},
	{ID: "iota", Name: "Iota index", Status: "Done", Date: "2026-09-11"},
	{ID: "kappa", Name: "Kappa layout", Status: "Active", Date: "2026-10-23"},
	{ID: "lambda", Name: "Lambda mailer", Status: "Pending", Date: "2026-11-05"},
	{ID: "mu", Name: "Mu navigation", Status: "Done", Date: "2026-12-14"},
}

// dataTableColumn is one sortable column header.
type dataTableColumn struct {
	Key      string
	Label    string
	Href     string
	AriaSort string // "" | "ascending" | "descending"
	Active   bool
	SortIcon template.HTML
}

// dataTableRowView is one rendered table row.
type dataTableRowView struct {
	ID       string
	Name     string
	Status   string
	Date     string
	Selected bool
}

// dataTablePageView is one pagination page link.
type dataTablePageView struct {
	Num     int
	Href    string
	Current bool
}

// emptyStateView is the view data for the shared "empty-state" partial. Title
// and Body carry the message, Icon may carry a trusted inline SVG glyph
// (template.HTML, decorative), and CTA is a real link for GET navigation —
// never a div/span control. Compact switches the primitive to start alignment
// for inline contexts like a table row.
type emptyStateView struct {
	Title    string
	Body     string
	Icon     template.HTML
	CTA      bool
	CTAHref  string
	CTALabel string
	Compact  bool
}

// dataTableDemo is the view model for the Data table documentation preview.
type dataTableDemo struct {
	Query           string
	Sort            string
	Dir             string
	Page            int
	Total           int
	Pages           int
	Caption         string
	Columns         []dataTableColumn
	Rows            []dataTableRowView
	PageLinks       []dataTablePageView
	PrevHref        string
	NextHref        string
	HasPrev         bool
	HasNext         bool
	SelectedCount   int
	SelectionNotice string
	Colspan         int
	EmptyState      emptyStateView
	Refreshed       bool
	RefreshToast    *toastView
}

// Trusted, internal, decorative sort glyphs. Each is aria-hidden and
// unfocusable by contract; the visible header label supplies the accessible
// name, and the active column additionally carries aria-sort.
const (
	// #nosec G203 -- trusted, internal decorative glyph.
	dataTableSortAscSVG template.HTML = `<svg class="ui-data-table-sort-icon" aria-hidden="true" focusable="false" viewBox="0 0 24 24" width="18" height="18" fill="currentColor"><path d="M7.41 15.41 12 10.83l4.59 4.58L18 14l-6-6-6 6z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	dataTableSortDescSVG template.HTML = `<svg class="ui-data-table-sort-icon" aria-hidden="true" focusable="false" viewBox="0 0 24 24" width="18" height="18" fill="currentColor"><path d="M7.41 8.59 12 13.17l4.59-4.58L18 10l-6 6-6-6z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	dataTableSortIdleSVG template.HTML = `<svg class="ui-data-table-sort-icon" aria-hidden="true" focusable="false" viewBox="0 0 24 24" width="18" height="18" fill="currentColor"><path d="M7.41 8.59 12 13.17l4.59-4.58L18 10l-6 6-6-6z"></path></svg>`
)

func (s *server) dataTableDocs(w http.ResponseWriter, r *http.Request) {
	selection := r.URL.Query()["selection"]
	demo := newDataTableDemo(r.URL.Query().Get("q"), r.URL.Query().Get("sort"), r.URL.Query().Get("dir"), r.URL.Query().Get("page"), selection)

	if strings.EqualFold(r.Header.Get("HX-Request"), "true") {
		var out strings.Builder
		if err := s.templates.ExecuteTemplate(&out, "data-table-panel", demo); err != nil {
			http.Error(w, "table unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(out.String()))
		return
	}

	s.renderMarkdownPage(w, pageView{
		Title:         "Data table",
		DataTableDemo: demo,
	}, "content/data-table.md")
}

// newDataTableDemo validates the request against closed vocabularies and builds
// the filtered, sorted, paginated view model. Unknown sort keys, directions and
// page numbers fall back to safe defaults; the submitted selection values are
// checked against the dataset and mark the matching rows checked server-side.
func newDataTableDemo(q, sortParam, dir, page string, selection []string) *dataTableDemo {
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

	items := make([]dataTableItem, 0, len(dataTableItems))
	for _, it := range dataTableItems {
		if query != "" && !strings.Contains(strings.ToLower(it.Name), strings.ToLower(query)) && !strings.Contains(strings.ToLower(it.Status), strings.ToLower(query)) {
			continue
		}
		items = append(items, it)
	}

	sort.SliceStable(items, func(i, j int) bool {
		a, b := dataTableField(items[i], sortKey), dataTableField(items[j], sortKey)
		if direction == "desc" {
			return a > b
		}
		return a < b
	})

	total := len(items)
	totalPages := (total + dataTablePageSize - 1) / dataTablePageSize
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

	start := (pageNum - 1) * dataTablePageSize
	end := start + dataTablePageSize
	if end > total {
		end = total
	}

	rows := make([]dataTableRowView, 0, end-start)
	for _, it := range items[start:end] {
		rows = append(rows, dataTableRowView{
			ID:       it.ID,
			Name:     it.Name,
			Status:   it.Status,
			Date:     it.Date,
			Selected: selectAll || selectedSet[it.ID],
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

	// The empty row spans every column plus the leading checkbox column, so
	// the message never collapses the table grid.
	colspan := 1 + len(dataTableColumns(query, sortKey, direction))
	var emptyState emptyStateView
	if total == 0 {
		emptyState = emptyStateView{
			Title:    "No results",
			Body:     "Try adjusting your search or filters.",
			CTA:      true,
			CTAHref:  dataTableHref("", sortKey, direction, 0),
			CTALabel: "Clear search",
			Compact:  true,
		}
	}

	pageLinks := make([]dataTablePageView, 0, totalPages)
	for n := 1; n <= totalPages; n++ {
		pageLinks = append(pageLinks, dataTablePageView{
			Num:     n,
			Href:    dataTableHref(query, sortKey, direction, n),
			Current: n == pageNum,
		})
	}

	demo := &dataTableDemo{
		Query:           query,
		Sort:            sortKey,
		Dir:             direction,
		Page:            pageNum,
		Total:           total,
		Pages:           totalPages,
		Caption:         fmt.Sprintf("%d rows · page %d of %d", total, pageNum, totalPages),
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
		EmptyState:      emptyState,
		Refreshed:       false,
	}
	return demo
}

// dataTableField returns the sortable field for a row.
func dataTableField(it dataTableItem, key string) string {
	switch key {
	case "status":
		return it.Status
	case "date":
		return it.Date
	default:
		return it.Name
	}
}

// dataTableColumns builds the three sortable columns. The active column carries
// aria-sort and toggles direction; every other column links to its ascending
// sort. Hrefs preserve the current query so filtering survives a sort.
func dataTableColumns(query, sortKey, direction string) []dataTableColumn {
	labels := map[string]string{"name": "Name", "status": "Status", "date": "Date"}
	cols := make([]dataTableColumn, 0, len(dataTableSortKeys))
	for _, k := range dataTableSortKeys {
		col := dataTableColumn{Key: k, Label: labels[k]}
		if k == sortKey {
			col.Active = true
			col.SortIcon = dataTableSortDescSVG
			col.AriaSort = "descending"
			if direction == "asc" {
				col.AriaSort = "ascending"
				col.SortIcon = dataTableSortAscSVG
			}
			col.Href = dataTableHref(query, k, toggleDataTableDir(direction), 0)
		} else {
			col.SortIcon = dataTableSortIdleSVG
			col.Href = dataTableHref(query, k, "asc", 0)
		}
		cols = append(cols, col)
	}
	return cols
}

func toggleDataTableDir(dir string) string {
	if dir == "asc" {
		return "desc"
	}
	return "asc"
}

// dataTableHref builds a real GET sort/page link preserving the current query.
// The parameter order is stable (q, sort, dir, page) so links are predictable
// and testable; q is percent-encoded by url.QueryEscape.
func dataTableHref(q, sort, dir string, page int) string {
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
	if q != "" {
		write("q", q)
	}
	write("sort", sort)
	write("dir", dir)
	if page >= 1 {
		write("page", strconv.Itoa(page))
	}
	return b.String()
}

// dataTableSelectionNotice summarizes a submitted selection for the no-JS form
// round-trip, mirroring the List selection notice contract.
func dataTableSelectionNotice(selection []string) string {
	selected := 0
	all := false
	for _, v := range selection {
		if v == "all" {
			all = true
		}
		selected++
	}
	if len(selection) == 0 {
		return ""
	}
	if all {
		return "All rows selected."
	}
	if selected == 1 {
		return "1 row selected."
	}
	return fmt.Sprintf("%d rows selected.", selected)
}

// dataTableRefreshDemo completes a remote refresh operation for the docs demo.
// Without JavaScript it re-renders the full documentation page with the table
// refreshed and a persistent inline toast plus a determinate progress bar; with
// HTMX it returns only the refresh form fragment and an HX-Trigger that raises
// the loom:toast event, which the local enhancement layer displays in the
// aria-live region. The progress element reuses the existing .ui-progress
// primitive; the toast reuses .ui-toast. This mirrors the Toast demo flow.
func (s *server) dataTableRefreshDemo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	isHX := strings.EqualFold(r.Header.Get("HX-Request"), "true")

	demo := newDataTableDemo("", "", "", "", nil)
	demo.Refreshed = true

	if isHX {
		trigger, err := toastTriggerJSON("success", "Data refreshed.")
		if err != nil {
			http.Error(w, "refresh unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Trigger", trigger)

		var rendered strings.Builder
		if err := s.templates.ExecuteTemplate(&rendered, "data-table-refresh-form", demo); err != nil {
			http.Error(w, "refresh unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(rendered.String()))
		return
	}

	toast := newToast("success", "data-table-refresh-result", "Data refreshed.")
	demo.RefreshToast = &toast
	s.renderMarkdownPageStatus(w, pageView{
		Title:         "Data table",
		DataTableDemo: demo,
	}, "content/data-table.md", http.StatusOK)
}
