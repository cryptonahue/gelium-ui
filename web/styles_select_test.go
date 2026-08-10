package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestSelectPrimitiveCSSMapsNativeVariantsAndStates(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-select {`,
		`gap: 0;`,
		`.ui-select select {`,
		`appearance: none;`,
		`height: var(--ui-select-height);`,
		`border-radius: var(--ui-select-radius);`,
		`.ui-select-filled select { background: var(--ui-select-container-filled); border: var(--ui-border-width-1) var(--ui-border-style-solid) transparent; border-bottom: var(--ui-border-width-1) var(--ui-border-style-solid) var(--ui-select-outline);`,
		`.ui-select-outlined select { background: transparent; border: var(--ui-border-width-1) var(--ui-border-style-solid) var(--ui-select-outline);`,
		`.ui-select select:focus-visible { outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring); outline-offset: var(--ui-focus-offset);`,
		`.ui-select select:disabled { cursor: not-allowed;`,
		`.ui-select select[aria-invalid="true"] { border-color: var(--ui-select-error);`,
		`.ui-select-caret`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing select contract %q", contract)
		}
	}

	labelFloating := []string{
		`.ui-select:focus-within .ui-select-label`,
		`.ui-select:has(select option:checked:not([value=""])) .ui-select-label`,
	}
	for _, sel := range labelFloating {
		if !strings.Contains(css, sel) {
			t.Errorf("source CSS is missing floating-label selector %q", sel)
		}
	}
	if strings.Contains(css, `.ui-select select { appearance: none; opacity: 0;`) {
		t.Error("select must keep the native select visible as the control surface, not hide it")
	}
	if !strings.Contains(css, `.ui-select select:not(:has(option:checked:not([value=""]))) { color: transparent;`) {
		t.Error("source CSS must hide the placeholder option text while empty so the resting label is the only prompt")
	}
	if !strings.Contains(css, `.ui-select select option:checked { background: var(--ui-select-list-bg); color: var(--ui-select-list-fg);`) {
		t.Error("source CSS must style the checked option row with the list palette so the browser popup does not tint the selected placeholder")
	}
	labelRule := regexp.MustCompile(`(?s)\.ui-select-label\s*\{([^}]*)\}`).FindStringSubmatch(css)
	if labelRule == nil {
		t.Fatal("source CSS is missing the base select label rule")
	}
	if strings.Contains(labelRule[1], "pointer-events: none") {
		t.Error("select label must receive pointer events so its for attribute can focus the native select")
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	for _, contract := range []string{
		`.ui-select select { border-color: CanvasText;`,
		`.ui-select select:disabled { border-color: GrayText;`,
		`.ui-select select[aria-invalid="true"] { border-color: Mark;`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("select must stay distinguishable in forced colors; missing %q", contract)
		}
	}
}

func TestSelectThemeDefinesPublicUIFamily(t *testing.T) {
	theme := themeCSS(t, "theme-material")
	for _, token := range []string{
		"--ui-select-height:",
		"--ui-select-radius:",
		"--ui-select-radius-top:",
		"--ui-select-caret:",
		"--ui-select-fg:",
		"--ui-select-label:",
		"--ui-select-outline:",
		"--ui-select-container-filled:",
		"--ui-select-hover-outline:",
		"--ui-select-focus:",
		"--ui-select-error:",
		"--ui-select-disabled-opacity:",
		"--ui-select-list-bg:",
		"--ui-select-list-fg:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing select token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesSelectContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-select`,
		`var(--ui-select-height)`,
		`var(--ui-select-outline)`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled select CSS is missing %q", contract)
		}
	}
}
