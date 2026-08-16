package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestRadioPrimitiveCSSMapsNativeControlAndStates(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-radio {`,
		`gap: var(--ui-space-2);`,
		`input[type="radio"] {`,
		`appearance: none;`,
		`width: var(--ui-radio-size);`,
		`height: var(--ui-radio-size);`,
		`border-radius: var(--ui-radio-radius);`,
		`border: var(--ui-radio-outline-width) solid var(--ui-radio-outline);`,
		`input:checked + .ui-radio-mark {`,
		`background: var(--ui-radio-checked);`,
		`input:disabled {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing radio contract %q", contract)
		}
	}

	stateLayers := []string{
		`.ui-radio:hover {`,
		`:focus-visible`,
		`:active:not(:disabled)`,
	}
	for _, sel := range stateLayers {
		if !strings.Contains(css, sel) {
			t.Errorf("source CSS is missing radio state selector %q", sel)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	for _, contract := range []string{
		`border-color: CanvasText;`,
		`background: CanvasText;`,
		`background: GrayText;`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("radio must stay operable in forced colors; missing %q", contract)
		}
	}
}

func TestRadioThemeDefinesPublicUIFamily(t *testing.T) {
	theme := themeCSS(t, "theme-material")
	for _, token := range []string{
		"--ui-radio-size:",
		"--ui-radio-radius:",
		"--ui-radio-outline-width:",
		"--ui-radio-outline:",
		"--ui-radio-hover-outline:",
		"--ui-radio-checked:",
		"--ui-radio-disabled:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing radio token %q", token)
		}
	}
}

func TestRadioReducedMotionDisablesTransitions(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "radio.css"), " ")
	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("radio.css must include a reduced-motion media query")
	}
	if !strings.Contains(css, `.ui-radio input[type="radio"], .ui-radio-mark::after { transition: none; }`) {
		t.Error("radio reduced-motion must disable the control and dot transitions")
	}
	if !strings.Contains(css, `.ui-radio:active:not(:disabled) input[type="radio"] { transform: none; }`) {
		t.Error("radio reduced-motion must drop the active scale transform (G11)")
	}
}

func TestEmbeddedCompiledCSSIncludesRadioContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-radio`,
		`var(--ui-radio-size)`,
		`var(--ui-radio-outline)`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded radio CSS is missing %q", contract)
		}
	}
}
