package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestCheckboxPrimitiveCSSMapsNativeControlAndStates(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-checkbox {`,
		`gap: .5rem;`,
		`input[type="checkbox"] {`,
		`appearance: none;`,
		`width: var(--ui-checkbox-size);`,
		`height: var(--ui-checkbox-size);`,
		`border-radius: var(--ui-checkbox-radius);`,
		`border: var(--ui-checkbox-outline-width) solid var(--ui-checkbox-outline);`,
		`input:checked + .ui-checkbox-mark {`,
		`input:checked ~ .ui-checkbox-label {`,
		`background: var(--ui-checkbox-container);`,
		`input:disabled {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing checkbox contract %q", contract)
		}
	}

	stateLayers := []string{
		`.ui-checkbox:hover {`,
		`:focus-visible`,
		`:active:not(:disabled)`,
	}
	for _, sel := range stateLayers {
		if !strings.Contains(css, sel) {
			t.Errorf("source CSS is missing checkbox state selector %q", sel)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	for _, contract := range []string{
		`border-color: CanvasText;`,
		`border-color: Canvas;`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("checkbox must stay operable in forced colors; missing %q", contract)
		}
	}
}

func TestCheckboxThemeDefinesPublicUIFamily(t *testing.T) {
	theme := repositoryFile(t, "themes", "theme-material", "theme.css")
	for _, token := range []string{
		"--ui-checkbox-size:",
		"--ui-checkbox-radius:",
		"--ui-checkbox-outline-width:",
		"--ui-checkbox-outline:",
		"--ui-checkbox-hover-outline:",
		"--ui-checkbox-container:",
		"--ui-checkbox-icon:",
		"--ui-checkbox-error:",
		"--ui-checkbox-checked-disabled-container:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing checkbox token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesCheckboxContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-checkbox`,
		`var(--ui-checkbox-size)`,
		`var(--ui-checkbox-outline)`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled checkbox CSS is missing %q", contract)
		}
	}
}
