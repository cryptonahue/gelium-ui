package app

import (
	"fmt"
	"strings"
	"testing"
)

// TestPaginationBuildsLinksAndCurrentMarker proves the standalone partial renders
// real links for every page plus the current page as an aria-current span.
func TestPaginationBuildsLinksAndCurrentMarker(t *testing.T) {
	view := newPaginationView(2, 3, func(n int) string { return fmt.Sprintf("/feed?page=%d", n) })
	if view.Current != 2 || view.Total != 3 {
		t.Fatalf("Current/Total = %d/%d, want 2/3", view.Current, view.Total)
	}
	if !view.HasPrev || !view.HasNext {
		t.Error("a middle page must have prev and next links")
	}

	body := renderPartial(t, "pagination", view)
	for _, contract := range []string{
		`<nav class="ui-pagination" aria-label="Pagination">`,
		`href="/feed?page=1"`,
		`href="/feed?page=3"`,
		`<span class="ui-pagination-page ui-pagination-page--current" aria-current="page">2</span>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("pagination is missing %q", contract)
		}
	}
}

// TestPaginationBoundariesUseDisabledSpans proves the first/last page switch the
// boundary link to an aria-disabled span instead of a dead link.
func TestPaginationBoundariesUseDisabledSpans(t *testing.T) {
	view := newPaginationView(1, 2, func(n int) string { return fmt.Sprintf("/q?page=%d", n) })
	if view.HasPrev {
		t.Error("the first page must not expose a previous link")
	}
	if !view.HasNext {
		t.Error("the first page of two must expose a next link")
	}
	body := renderPartial(t, "pagination", view)
	for _, contract := range []string{
		`<span class="ui-pagination-page ui-pagination-page--disabled" aria-disabled="true">Previous</span>`,
		`<a class="ui-pagination-page" href="/q?page=2">Next</a>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("first page is missing %q", contract)
		}
	}

	view = newPaginationView(2, 2, func(n int) string { return fmt.Sprintf("/q?page=%d", n) })
	if view.HasNext {
		t.Error("the last page must not expose a next link")
	}
	body = renderPartial(t, "pagination", view)
	if !strings.Contains(body, `<span class="ui-pagination-page ui-pagination-page--disabled" aria-disabled="true">Next</span>`) {
		t.Error("the last page must render a disabled Next span")
	}
}

// TestPaginationClampsPage proves out-of-range pages clamp into [1, total] and
// the label can be customized.
func TestPaginationClampsPage(t *testing.T) {
	view := newPaginationView(99, 3, func(n int) string { return fmt.Sprintf("?p=%d", n) })
	if view.Current != 3 {
		t.Errorf("page 99 must clamp to the last page, got %d", view.Current)
	}
	view = newPaginationView(0, 3, func(n int) string { return fmt.Sprintf("?p=%d", n) })
	if view.Current != 1 {
		t.Errorf("page 0 must clamp to the first page, got %d", view.Current)
	}
	view.Label = "Queue pages"
	if body := renderPartial(t, "pagination", view); !strings.Contains(body, `aria-label="Queue pages"`) {
		t.Error("pagination must honor the custom aria-label")
	}
}
