package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestCardPrimitiveCSSMapsVariantsToThemeTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-card {`,
		`border-radius: var(--ui-card-radius);`,
		`.ui-card-elevated { background: var(--ui-card-container-elevated); box-shadow: var(--ui-shadow-1);`,
		`.ui-card-filled { background: var(--ui-card-container-filled);`,
		`.ui-card-outlined { background: var(--ui-card-container-outlined); border: 1px solid var(--ui-card-outline-color);`,
		`.ui-card:focus-visible { outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring); outline-offset: var(--ui-focus-offset);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing card contract %q", contract)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	if !strings.Contains(forced, ".ui-card { border: 1px solid CanvasText;") {
		t.Error("card must keep a visible boundary in forced colors")
	}
}

func TestCardThemeDefinesPublicUIFamily(t *testing.T) {
	theme := repositoryFile(t, "themes", "theme-material", "theme.css")
	for _, token := range []string{
		"--ui-card-radius:",
		"--ui-card-container-elevated:",
		"--ui-card-container-filled:",
		"--ui-card-container-outlined:",
		"--ui-card-outline-color:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing card token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesCardContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-card{`,
		`.ui-card-elevated{`,
		`.ui-card-filled{`,
		`.ui-card-outlined{`,
		`@media (forced-colors:active)`,
		`var(--ui-card-container-filled)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled card CSS is missing %q", contract)
		}
	}
}
