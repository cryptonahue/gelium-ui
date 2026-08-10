package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestBadgePrimitiveCSSMapsSizesAndNotColorOnly(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-badge {`,
		`width: var(--ui-badge-size);`,
		`height: var(--ui-badge-size);`,
		`background: var(--ui-badge-container);`,
		`border-radius: var(--ui-radius-full);`,
		`.ui-badge-large {`,
		`min-width: var(--ui-badge-large-size);`,
		`height: var(--ui-badge-large-size);`,
		`color: var(--ui-badge-fg);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing badge contract %q", contract)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	if !strings.Contains(forced, ".ui-badge { border: 1px solid CanvasText;") {
		t.Error("badge must keep a visible boundary in forced colors")
	}
}

func TestBadgeThemeDefinesPublicUIPair(t *testing.T) {
	theme := themeCSS(t, "theme-material")
	for _, token := range []string{
		"--ui-badge-size:",
		"--ui-badge-large-size:",
		"--ui-badge-container:",
		"--ui-badge-fg:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing badge token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesBadgeContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-badge{`,
		`.ui-badge-large{`,
		`@media (forced-colors:active)`,
		`var(--ui-badge-container)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled badge CSS is missing %q", contract)
		}
	}
}
