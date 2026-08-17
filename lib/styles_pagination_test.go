package lib

import (
	"regexp"
	"strings"
	"testing"
)

// TestPaginationPrimitiveCSSMapsPagesAndBoundaries proves the standalone
// pagination paints pages with the scoped tokens and the full radius, with a
// distinct active and disabled state.
func TestPaginationPrimitiveCSSMapsPagesAndBoundaries(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "pagination.css"), " ")

	for _, contract := range []string{
		`.ui-pagination {`,
		`--ui-pagination-page-color: var(--ui-color-fg-muted);`,
		`--ui-pagination-active-color: var(--ui-color-primary);`,
		`.ui-pagination-page {`,
		`border-radius: var(--ui-radius-full);`,
		`font: var(--ui-type-label-sm);`,
		`.ui-pagination-page--current {`,
		`color: var(--ui-pagination-active-color);`,
		`.ui-pagination-page--disabled {`,
		`opacity: var(--ui-state-disabled-opacity);`,
		`pointer-events: none;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("pagination.css is missing contract %q", contract)
		}
	}

	if !strings.Contains(css, "var(--ui-space-") {
		t.Error("pagination.css must consume the spacing scale for the page gap/padding")
	}
	if !strings.Contains(css, "var(--ui-focus-thickness)") {
		t.Error("pagination.css must paint the focus ring from the focus tokens")
	}
}
