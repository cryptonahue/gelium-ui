package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestBadgePrimitiveCSSMapsSizesAndNotColorOnly(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-badge {`,
		`width: var(--ui-badge-size, 6px);`,
		`height: var(--ui-badge-size, 6px);`,
		`background: var(--ui-badge-container);`,
		`border-radius: var(--ui-radius-full);`,
		`.ui-badge-large {`,
		`min-width: var(--ui-badge-large-size, 16px);`,
		`height: var(--ui-badge-large-size, 16px);`,
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
	css := compiledAppCSS(t)
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

// TestBadgeToneVariantsReuseSemanticTokens proves the tone variants reuse the
// closed semantic color tokens (danger/success/warning/info + containers) so
// every tone follows light/dark/forced-colors and never a raw literal.
func TestBadgeToneVariantsReuseSemanticTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-badge--error { background: var(--ui-color-danger); color: var(--ui-color-danger-fg); }`,
		`.ui-badge--success { background: var(--ui-color-success); color: var(--ui-color-success-fg); }`,
		`.ui-badge--warning { background: var(--ui-color-warning-container); color: var(--ui-color-warning-fg); }`,
		`.ui-badge--info { background: var(--ui-color-info); color: var(--ui-color-info-fg); }`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing badge tone contract %q", contract)
		}
	}
}

// TestBadgeInfoToneTokenDefinedInCoreAndTheme proves the --ui-color-info-fg
// on-color the info tone needs is closed across the core (neutral default) and
// the Material theme (light + the single dark class route).
func TestBadgeInfoToneTokenDefinedInCoreAndTheme(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	if !strings.Contains(string(core), "--ui-color-info-fg:") {
		t.Error("core tokens.css must define --ui-color-info-fg (badge info tone on-color)")
	}

	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, defaultThemeName), " ")
	if n := strings.Count(theme, "--ui-color-info-fg:"); n != 2 {
		t.Errorf("theme-material must define --ui-color-info-fg once in light and once in the single dark class route, got %d definitions", n)
	}
}
