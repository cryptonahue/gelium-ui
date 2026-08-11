package app

// paginationView is the server-driven view model for the standalone
// "pagination" partial: the reusable extraction of the Data table footer nav.
// Current and Total are the numeric state; PageLinks carries every numbered
// page with the current one marked; PrevHref/NextHref are real GET links and
// HasPrev/HasNext switch the boundary links to aria-disabled spans. A nil
// *paginationView renders nothing, so recipes can omit pagination on short
// sets.
type paginationView struct {
	Label     string
	PageLinks []paginationPageLink
	PrevHref  string
	NextHref  string
	HasPrev   bool
	HasNext   bool
	Current   int
	Total     int
}

// paginationPageLink is one numbered page: a real link, or the current page
// rendered as a span with aria-current="page".
type paginationPageLink struct {
	Num     int
	Href    string
	Current bool
}

// newPaginationView builds the view from a page-counting href function, so any
// recipe paginates over its own query vocabulary (the same contract as
// dataTableHref). The current page is clamped into [1, total]; when there is a
// single page the returned view is non-nil but a recipe should render it only
// when total > 1.
func newPaginationView(current, total int, href func(int) string) *paginationView {
	if total < 1 {
		total = 1
	}
	if current < 1 {
		current = 1
	}
	if current > total {
		current = total
	}
	links := make([]paginationPageLink, 0, total)
	for n := 1; n <= total; n++ {
		links = append(links, paginationPageLink{Num: n, Href: href(n), Current: n == current})
	}
	return &paginationView{
		PageLinks: links,
		PrevHref:  href(current - 1),
		NextHref:  href(current + 1),
		HasPrev:   current > 1,
		HasNext:   current < total,
		Current:   current,
		Total:     total,
	}
}
