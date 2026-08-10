package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestSegmentedButtonPrimitiveCSSMapsMaterialAnatomy(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "segmented-button.css"), " ")

	for _, contract := range []string{
		`.ui-segmented-button-set {`,
		`height: var(--ui-segmented-button-container-height);`,
		`border: var(--ui-border-width-1) var(--ui-border-style-solid) var(--ui-segmented-button-outline);`,
		`border-radius: var(--ui-segmented-button-container-radius);`,
		`--ui-segmented-button-icon-size: var(--ui-size-icon-sm);`,
		`.ui-segmented-button:first-of-type {`,
		`border-start-start-radius: var(--ui-segmented-button-container-radius);`,
		`.ui-segmented-button:not(:first-of-type) {`,
		`border-inline-start: var(--ui-border-width-1) var(--ui-border-style-solid) var(--ui-segmented-button-outline);`,
		`.ui-segmented-button-checkmark-path {`,
		`stroke-width: 2px;`,
		`stroke-dasharray: 29.7833385;`,
		`.ui-segmented-button-label {`,
		`font-size: .875rem;`,
		`font-weight: 500;`,
		`letter-spacing: .00625rem;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source segmented-button CSS is missing contract %q", contract)
		}
	}
}

// TestSegmentedButtonSelectionComesFromNativeChecked guards the platform-first
// decision that replaced upstream's <button aria-pressed> + JS selection with
// native radio/checkbox semantics: every selected-state style derives from the
// native :checked pseudo-class, and every state layer uses the shared
// --ui-state-* opacities.
func TestSegmentedButtonSelectionComesFromNativeChecked(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "segmented-button.css"), " ")

	for _, contract := range []string{
		`.ui-segmented-button:has(input:checked) {`,
		`background: var(--ui-segmented-button-selected-container);`,
		`.ui-segmented-button:has(input:checked) .ui-segmented-button-graphic {`,
		`.ui-segmented-button:has(input:checked) .ui-segmented-button-checkmark-path {`,
		`.ui-segmented-button:has(input:checked) .ui-segmented-button-icon { display: none;`,
		`.ui-segmented-button--icon-only:has(input:checked) .ui-segmented-button-icon {`,
		`.ui-segmented-button:hover::before {`,
		`.ui-segmented-button:focus-visible::before,`,
		`.ui-segmented-button:has(input:focus-visible)::before {`,
		`.ui-segmented-button:active::before {`,
		`.ui-segmented-button:has(input:disabled) {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source segmented-button CSS is missing selection/state contract %q", contract)
		}
	}
}

// TestSegmentedButtonDemoClassVocabularyIsClosed is the TDD regression for the
// closed vocabulary rule: the demo/preview scaffolding classes in
// segmented-button.html MUST match the CSS selectors exactly, and the ui- prefix
// is reserved for the component primitives (.ui-segmented-button*), never for
// the demo grid.
func TestSegmentedButtonDemoClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "web", "templates", "segmented-button.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "segmented-button.css"), " ")

	for _, cls := range []string{
		"segmented-button-demo-grid",
		"segmented-button-demo-group",
		"segmented-button-demo-heading",
		"segmented-button-demo-form",
		"segmented-button-demo-submit",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("segmented-button.html is missing demo class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("segmented-button.css is missing demo selector .%s", cls)
		}
	}

	for _, cls := range []string{
		"ui-segmented-button-set",
		"ui-segmented-button-legend",
		"ui-segmented-button",
		"ui-segmented-button--icon-only",
		"ui-segmented-button--action",
		"ui-segmented-button-graphic",
		"ui-segmented-button-checkmark",
		"ui-segmented-button-checkmark-path",
		"ui-segmented-button-icon",
		"ui-segmented-button-label",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("segmented-button.html is missing primitive class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("segmented-button.css is missing primitive selector .%s", cls)
		}
	}

	if strings.Contains(tmpl, "ui-segmented-button-demo") {
		t.Error("segmented-button.html must not ui-prefix the demo scaffolding classes")
	}
	if strings.Contains(css, ".ui-segmented-button-demo") {
		t.Error("segmented-button.css must not define .ui-segmented-button-demo selectors")
	}
}

func TestSegmentedButtonReducedMotionAndForcedColorsWired(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "segmented-button.css"), " ")

	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("segmented-button.css must include a reduced-motion media query")
	}
	if !strings.Contains(css, `transition: none;`) {
		t.Error("segmented-button reduced-motion must disable transitions")
	}
	if !strings.Contains(css, `@media (forced-colors: active)`) {
		t.Error("segmented-button.css must include a forced-colors media query")
	}
	if !strings.Contains(css, `.ui-segmented-button:has(input:checked) { background: Highlight; color: HighlightText;`) {
		t.Error("segmented-button selected fill must switch to Highlight/HighlightText in forced colors")
	}
	if !strings.Contains(css, `.ui-segmented-button:has(input:disabled) { background: Canvas; color: GrayText;`) {
		t.Error("segmented-button disabled segments must switch to Canvas/GrayText in forced colors")
	}
	if !strings.Contains(css, `.ui-segmented-button:has(input:focus-visible) { outline-color: Highlight;`) {
		t.Error("segmented-button focus ring must switch to Highlight in forced colors")
	}
}

func TestEmbeddedCompiledCSSIncludesSegmentedButtonContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-segmented-button-set`,
		`.ui-segmented-button`,
		`.ui-segmented-button--icon-only`,
		`.ui-segmented-button--action`,
		`.ui-segmented-button-graphic`,
		`.ui-segmented-button-checkmark-path`,
		`.ui-segmented-button-icon`,
		`.ui-segmented-button-label`,
		`.segmented-button-demo-grid`,
		`.segmented-button-demo-group`,
		`.segmented-button-demo-heading`,
		`.segmented-button-demo-form`,
		`.segmented-button-demo-submit`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled segmented-button CSS is missing %q", contract)
		}
	}
}
