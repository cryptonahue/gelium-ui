package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestChipsPrimitiveCSSMapsMaterialAnatomy(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "chips.css"), " ")

	for _, contract := range []string{
		`.ui-chip {`,
		`height: var(--ui-chip-height);`,
		`border: var(--ui-border-width-1) var(--ui-border-style-solid) var(--ui-color-border-strong);`,
		`border-radius: var(--ui-radius-sm);`,
		`font: var(--ui-type-label-lg);`,
		`.ui-chip:hover:not(:disabled):not([aria-disabled="true"])`,
		`.ui-chip:active:not(:disabled):not([aria-disabled="true"])`,
		`.ui-chip:focus-visible`,
		`outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring);`,
		`.ui-chip:disabled,`,
		`opacity: var(--ui-state-disabled-opacity);`,
		`var(--ui-state-hover-opacity)`,
		`var(--ui-state-pressed-opacity)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source chips CSS is missing contract %q", contract)
		}
	}
}

func TestChipsFilterUsesNativeCheckboxSemantics(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "chips.css"), " ")

	for _, contract := range []string{
		`.ui-chip-filter input[type="checkbox"] {`,
		`appearance: none;`,
		`.ui-chip-filter input:checked {`,
		`background: var(--ui-color-secondary);`,
		`color: var(--ui-color-secondary-fg);`,
		`.ui-chip-filter input:checked ~ .ui-chip-selected-icon`,
		`opacity: 1;`,
		`.ui-chip-filter input:checked ~ .ui-chip-label`,
		`color: var(--ui-color-secondary-fg);`,
		`.ui-chip-remove:focus-visible`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source chips CSS is missing filter/remove contract %q", contract)
		}
	}
}

// TestChipsDemoGridSelectorsMatchTemplate is the TDD regression for the closed
// vocabulary rule: the demo/preview scaffolding classes in chips.html MUST match
// the CSS selectors exactly, and the ui- prefix is reserved for the component
// primitives (.ui-chip*), never for the demo grid.
func TestChipsDemoGridClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "lib", "templates", "chips.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "chips.css"), " ")

	for _, cls := range []string{
		"chips-demo-grid",
		"chips-demo-group",
		"chips-demo-row",
		"chips-demo-heading",
		"chips-demo-notice",
		"chips-demo-row",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("chips.html is missing demo class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("chips.css is missing demo selector .%s", cls)
		}
	}

	for _, cls := range []string{
		"ui-chip",
		"ui-chip-assist",
		"ui-chip-filter",
		"ui-chip-suggestion",
		"ui-chip-input",
		"ui-chip-remove",
		"ui-chip-label",
		"ui-chip-icon",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("chips.html is missing primitive class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("chips.css is missing primitive selector .%s", cls)
		}
	}
}

func TestChipsReducedMotionAndForcedColorsWired(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "chips.css"), " ")

	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("chips.css must include a reduced-motion media query")
	}
	if !strings.Contains(css, `.ui-chip,`) {
		t.Error("chips reduced-motion must disable .ui-chip transitions")
	}
	if !strings.Contains(css, `@media (forced-colors: active)`) {
		t.Error("chips.css must include a forced-colors media query")
	}
	if !strings.Contains(css, `.ui-chip { border: 1px solid CanvasText;`) {
		t.Error("chips must keep a visible boundary in forced colors")
	}
}

func TestEmbeddedCompiledCSSIncludesChipsContracts(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-chip{`,
		`.ui-chip-assist`,
		`.ui-chip-filter`,
		`.ui-chip-suggestion`,
		`.ui-chip-input`,
		`.ui-chip-remove`,
		`.chips-demo-grid`,
		`.chips-demo-group`,
		`.chips-demo-row`,
		`.chips-demo-heading`,
		`.chips-demo-notice`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled chips CSS is missing %q", contract)
		}
	}
}
