package web

import (
	"regexp"
	"strings"
	"testing"
)

// TestSelectMenuDemoUsesNativeSelectFieldSurface proves the server-driven
// Select menu demo dogfoods the Select component's own field surface (the
// native <select>) instead of the old M3 dialog menu: the dead trigger and the
// menu surface must not ship, and the field's validation styling must apply.
func TestSelectMenuDemoUsesNativeSelectFieldSurface(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-select {`,
		`.ui-select select {`,
		`appearance: none;`,
		`height: var(--ui-select-height);`,
		`border-radius: var(--ui-select-radius);`,
		`.ui-select-filled select { background: var(--ui-select-container-filled); border: var(--ui-border-width-1) var(--ui-border-style-solid) transparent; border-bottom: var(--ui-border-width-1) var(--ui-border-style-solid) var(--ui-select-outline);`,
		`.ui-select select[aria-invalid="true"] { border-color: var(--ui-select-error);`,
		`.ui-select-error {`,
		`color: var(--ui-select-error);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing native select field contract %q", contract)
		}
	}

	for _, gone := range []string{".ui-select-menu", ".m3-select-trigger", "closedby"} {
		if strings.Contains(css, gone) {
			t.Errorf("dead M3 menu markup/CSS %q must not ship once the native select is the base control", gone)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	for _, contract := range []string{
		`.ui-select select { border-color: CanvasText;`,
		`.ui-select-error { color: Mark;`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("the select menu demo field must stay distinguishable in forced colors; missing %q", contract)
		}
	}

	reduced := entryMediaBlock(t, css, "@media (prefers-reduced-motion: reduce)")
	if !strings.Contains(reduced, ".ui-select") {
		t.Error("reduced-motion CSS must keep disabling the select field transitions")
	}
	if strings.Contains(reduced, ".ui-select-menu") {
		t.Error("reduced-motion CSS must not reference the removed M3 menu")
	}
}

func TestSelectMenuThemeDoesNotDefineDeadFamily(t *testing.T) {
	// The M3 menu was removed in the G1 fix (native <select> is the base
	// control). The select-menu token family is dead: the theme must not
	// carry it, otherwise the dead tokens linger without consumers.
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
		if strings.Contains(theme, token) {
			t.Errorf("theme must not define dead select-menu token %q after the G1 native-select fix", token)
		}
	}
}

func TestEmbeddedCompiledCSSExcludesDeadSelectMenu(t *testing.T) {
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
	for _, gone := range []string{".ui-select-menu", ".m3-select-trigger"} {
		if strings.Contains(css, gone) {
			t.Errorf("embedded compiled CSS must not contain the removed M3 menu %q", gone)
		}
	}
}
