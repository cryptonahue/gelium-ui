package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestSelectMenuPrimitiveCSSMapsMaterialMenuSurfaceAndItems(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-select-menu {`,
		`min-width: var(--ui-select-menu-min-width);`,
		`border-radius: var(--ui-select-menu-radius);`,
		`background: var(--ui-select-menu-container);`,
		`box-shadow: var(--ui-select-menu-elevation);`,
		`.ui-select-menu-item {`,
		`min-height: var(--ui-select-menu-item-height);`,
		`font: var(--ui-type-label-lg);`,
		`.ui-select-menu-item:hover {`,
		`.ui-select-menu-item[aria-selected="true"] { background: var(--ui-select-menu-item-selected);`,
		`.ui-select-menu-divider {`,
		`height: 1px;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing select-menu contract %q", contract)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	for _, contract := range []string{
		`.ui-select-menu { border: 1px solid CanvasText;`,
		`.ui-select-menu-item[aria-selected="true"] { background: Highlight;`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("select menu must stay distinguishable in forced colors; missing %q", contract)
		}
	}

	reduced := entryMediaBlock(t, css, "@media (prefers-reduced-motion: reduce)")
	if !strings.Contains(reduced, ".ui-select-menu") {
		t.Error("reduced-motion CSS must disable select-menu transitions")
	}
}

func TestSelectMenuThemeDefinesPublicUIFamily(t *testing.T) {
	theme := themeCSS(t, "theme-material")
	for _, token := range []string{
		"--ui-select-menu-container:",
		"--ui-select-menu-radius:",
		"--ui-select-menu-elevation:",
		"--ui-select-menu-min-width:",
		"--ui-select-menu-item-height:",
		"--ui-select-menu-item-fg:",
		"--ui-select-menu-item-selected:",
		"--ui-select-menu-divider:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing select-menu token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesSelectMenuContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-select-menu`,
		`var(--ui-select-menu-container)`,
		`var(--ui-select-menu-item-height)`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled select-menu CSS is missing %q", contract)
		}
	}
}
