package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestDividerPrimitiveCSSMapsTokensAndInsets(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-divider {`,
		`height: var(--ui-divider-thickness);`,
		`background: var(--ui-divider-color);`,
		`.ui-divider-inset { padding-inline: var(--ui-space-4);`,
		`.ui-divider-inset-start { padding-inline-start: var(--ui-space-4);`,
		`.ui-divider-inset-end { padding-inline-end: var(--ui-space-4);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing divider contract %q", contract)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	if !strings.Contains(forced, ".ui-divider { background: CanvasText;") {
		t.Error("divider must keep a visible line in forced colors")
	}
}

func TestDividerThemeDefinesPublicUIPair(t *testing.T) {
	theme := themeCSS(t, "theme-material")
	for _, token := range []string{
		"--ui-divider-color:",
		"--ui-divider-thickness:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing divider token %q", token)
		}
	}
}
