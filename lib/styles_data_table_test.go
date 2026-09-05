package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestDataTablePrimitiveCSSMapsMaterialAnatomy(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "data-table.css"), " ")

	for _, contract := range []string{
		`.ui-data-table {`,
		`background: var(--ui-data-table-container-color);`,
		`border: var(--ui-border-width-1) var(--ui-border-style-solid) var(--ui-data-table-outline-color);`,
		`border-radius: var(--ui-radius-sm);`,
		`height: var(--ui-data-table-header-height);`,
		`height: var(--ui-data-table-row-height);`,
		`padding: 0 var(--ui-data-table-cell-padding);`,
		`width: var(--ui-data-table-checkbox-column-width);`,
		`.ui-data-table-sort-icon {`,
		`width: var(--ui-data-table-sort-icon-size);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source data-table CSS is missing contract %q", contract)
		}
	}
}

// TestDataTableRowStateComesFromNativeChecked guards the platform-first
// decision that replaced upstream's JS-driven selection with native checkbox
// semantics: the selected-row styles derive from :has(input:checked), and the
// hover/focus/active state layer uses the shared --ui-state-* opacities.
func TestDataTableRowStateComesFromNativeChecked(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "data-table.css"), " ")

	for _, contract := range []string{
		`.ui-data-table-row:hover .ui-data-table-cell::before {`,
		`opacity: var(--ui-data-table-hover-opacity);`,
		`.ui-data-table-row:focus-within .ui-data-table-cell::before {`,
		`.ui-data-table-row:active .ui-data-table-cell::before {`,
		`.ui-data-table-row:has(input:checked) .ui-data-table-cell {`,
		`background: var(--ui-data-table-selected-container-color);`,
		`.ui-data-table-checkbox input[type="checkbox"]:focus-visible {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source data-table CSS is missing state/selection contract %q", contract)
		}
	}
}

// TestDataTableDemoClassVocabularyIsClosed is the TDD regression for the closed
// vocabulary rule: the demo/preview scaffolding classes in data-table.html MUST
// match the CSS selectors exactly, and the ui- prefix is reserved for the
// component primitives (.ui-data-table*), never for the demo grid.
func TestDataTableDemoClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "lib", "templates", "data-table.html")
	layout := repositoryFile(t, "site", "web", "templates", "layout.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "data-table.css"), " ")

	for _, cls := range []string{
		"data-table-demo-grid",
		"data-table-demo-group",
		"data-table-demo-heading",
		"data-table-demo-filter",
		"data-table-demo-filter-label",
		"data-table-demo-filter-input",
		"data-table-demo-filter-submit",
		"data-table-demo-select",
		"data-table-demo-actions",
		"data-table-demo-submit",
		"data-table-demo-notice",
		"data-table-demo-refresh",
		"data-table-demo-refresh-row",
		"data-table-demo-refresh-button",
		"data-table-demo-refresh-hint",
		"data-table-demo-progress",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("data-table.html is missing demo class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("data-table.css is missing demo selector .%s", cls)
		}
	}

	// The preview modifier class wraps the demo in layout.html.
	if !strings.Contains(layout, "data-table-preview") {
		t.Error("layout.html is missing data-table-preview on the demo section")
	}
	if !strings.Contains(css, ".data-table-preview") {
		t.Error("data-table.css is missing demo selector .data-table-preview")
	}

	for _, cls := range []string{
		"ui-data-table",
		"ui-data-table-table",
		"ui-data-table-caption",
		"ui-data-table-cell",
		"ui-data-table-cell--checkbox",
		"ui-data-table-cell--label",
		"ui-data-table-cell--sortable",
		"ui-data-table-checkbox",
		"ui-data-table-row",
		"ui-data-table-sort",
		"ui-data-table-sort--active",
		"ui-data-table-pagination",
		"ui-data-table-page",
		"ui-data-table-page--current",
		"ui-data-table-page--disabled",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("data-table.html is missing primitive class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("data-table.css is missing primitive selector .%s", cls)
		}
	}
	if !strings.Contains(css, "min-width: var(--ui-touch-target);") || !strings.Contains(css, "min-height: var(--ui-touch-target);") {
		t.Error("data-table pagination pages must meet the shared touch-target token")
	}

	// The sort glyph and its class live in the trusted internal SVG constant in
	// data_table.go; the template references it through {{.SortIcon}}.
	goSrc := repositoryFile(t, "internal", "app", "data_table.go")
	if !strings.Contains(goSrc, "ui-data-table-sort-icon") {
		t.Error("data_table.go must define the ui-data-table-sort-icon SVG class")
	}
	if !strings.Contains(css, ".ui-data-table-sort-icon") {
		t.Error("data-table.css is missing primitive selector .ui-data-table-sort-icon")
	}

	if strings.Contains(tmpl, "ui-data-table-demo") {
		t.Error("data-table.html must not ui-prefix the demo scaffolding classes")
	}
	if strings.Contains(css, ".ui-data-table-demo") {
		t.Error("data-table.css must not define .ui-data-table-demo selectors")
	}
}

func TestDataTableReducedMotionAndForcedColorsWired(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "data-table.css"), " ")

	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("data-table.css must include a reduced-motion media query")
	}
	if !strings.Contains(css, `transition: none;`) {
		t.Error("data-table reduced-motion must disable transitions")
	}
	if !strings.Contains(css, `@media (forced-colors: active)`) {
		t.Error("data-table.css must include a forced-colors media query")
	}
	if !strings.Contains(css, `.ui-data-table-row:has(input:checked) .ui-data-table-cell { background: Highlight; color: HighlightText;`) {
		t.Error("data-table selected rows must switch to Highlight/HighlightText in forced colors")
	}
	if !strings.Contains(css, `.ui-data-table-sort { color: LinkText;`) {
		t.Error("data-table sort links must keep LinkText in forced colors")
	}
	if !strings.Contains(css, `.ui-data-table-sort:focus-visible { outline-color: Highlight;`) {
		t.Error("data-table sort focus ring must switch to Highlight in forced colors")
	}
}

// TestDataTableRowDividerAvoidsFirstLastChild guards the prior lesson: a
// container whose first child may not be the row itself must never rely on
// :first-child/:last-child for dividers. The divider uses adjacent sibling
// rows so the header row and caption stay clean.
func TestDataTableRowDividerAvoidsFirstLastChild(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "data-table.css"), " ")

	if !strings.Contains(css, `tbody .ui-data-table-row + .ui-data-table-row .ui-data-table-cell {`) {
		t.Error("data-table row divider must use adjacent-sibling rows")
	}
	if strings.Contains(css, `:last-child`) {
		t.Error("data-table.css must not rely on :last-child for dividers")
	}
	if strings.Contains(css, `.ui-data-table-row:last-child`) {
		t.Error("data-table.css must not style the last row via :last-child")
	}
}

func TestEmbeddedCompiledCSSIncludesDataTableContracts(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-data-table`,
		`.ui-data-table-cell`,
		`.ui-data-table-sort`,
		`.ui-data-table-row:has(input:checked)`,
		`.ui-data-table-pagination`,
		`.data-table-demo-grid`,
		`.data-table-demo-refresh`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled data-table CSS is missing %q", contract)
		}
	}
}
