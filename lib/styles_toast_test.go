package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestReleaseDocsAndPackageVersionStayCoherent(t *testing.T) {
	packageJSON := repositoryFile(t, "package.json")
	readme := repositoryFile(t, "README.md")
	toastDocs := repositoryFile(t, "site", "web", "content", "toast.md")
	if !strings.Contains(packageJSON, `"version": "0.5.0"`) {
		t.Error("package version must identify the 0.5.0 release")
	}
	for _, contract := range []string{"v0.5.0", "HTMX 4", "on-this-page rail", "prev/next pagination", "/components/toast", "gelium:toast", "sin JS"} {
		if !strings.Contains(readme, contract) {
			t.Errorf("README is missing release contract %q", contract)
		}
	}
	for _, contract := range []string{"role=\"alert\"", "role=\"status\"", "aria-live", "HX-Trigger", "prefers-reduced-motion", "forced-colors"} {
		if !strings.Contains(toastDocs, contract) {
			t.Errorf("Toast documentation is missing %q", contract)
		}
	}
}

func TestMaterialThemeDefinesToastTokensInEveryColorScheme(t *testing.T) {
	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, "theme-material"), " ")
	count := func(token string) int {
		return len(regexp.MustCompile(`(?i)`+regexp.QuoteMeta(token)).FindAllStringIndex(theme, -1))
	}
	if got := count(`--ui-toast-container:`); got != 2 {
		t.Errorf("toast container theme declarations = %d, want 2 (light + single dark class route)", got)
	}
	// The contract is the token family, never a concrete hex value. Color
	// tokens must be defined once per scheme (light + single dark class
	// route); shape tokens are scheme-independent and just need to exist.
	for _, token := range []string{
		`--ui-toast-container:`,
		`--ui-toast-fg:`,
		`--ui-toast-action:`,
	} {
		if got := count(token); got != 2 {
			t.Errorf("toast token %s must be defined once in light and once in the single dark class route, got %d", token, got)
		}
	}
	for _, token := range []string{
		`--ui-toast-radius:`,
		`--ui-toast-icon-info:`,
		`--ui-toast-icon-success:`,
		`--ui-toast-icon-warning:`,
		`--ui-toast-icon-error:`,
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("Material toast theme is missing %q", token)
		}
	}
}

func TestToastSourceCSSImplementsSnackbarAnatomyAndAccessibleStates(t *testing.T) {
	css := sourceAppCSS(t)
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(css, " ")
	for _, contract := range []string{
		`.ui-toast-region { position: fixed;`,
		`.ui-toast { --ui-toast-min-height: 3rem; --ui-toast-padding: .875rem 1rem; display: flex;`,
		`min-height: var(--ui-toast-min-height);`,
		`max-width: min(100%, 26rem);`,
		`border-radius: var(--ui-toast-radius);`,
		`background: var(--ui-toast-container);`,
		`color: var(--ui-toast-fg);`,
		`box-shadow: var(--ui-shadow-3);`,
		`.ui-toast-action:focus-visible { outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring); outline-offset: var(--ui-focus-offset);`,
		`.ui-toast-icon-info { color: var(--ui-toast-icon-info);`,
		`.ui-toast-icon-error { color: var(--ui-toast-icon-error);`,
	} {
		if !strings.Contains(compact, contract) {
			t.Errorf("toast CSS is missing %q", contract)
		}
	}
}

func TestToastReducedMotionDisablesRegionTransition(t *testing.T) {
	css := sourceAppCSS(t)
	reduced := entryMediaBlock(t, css, "@media (prefers-reduced-motion: reduce)")
	pattern := regexp.MustCompile(`(?s)\.ui-toast-region \.ui-toast\s*\{[^}]*transition:\s*none\s*;?[^}]*\}`)
	if !pattern.MatchString(reduced) {
		t.Error("reduced-motion CSS must disable the toast enter/exit transition")
	}
}

func TestToastForcedColorsProvidesBordersBeyondColor(t *testing.T) {
	css := sourceAppCSS(t)
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(css, " ")
	for _, contract := range []string{
		`@media (forced-colors: active)`,
		`.ui-toast { border: 1px solid CanvasText; forced-color-adjust: auto;`,
		`.ui-toast-action { color: Highlight;`,
	} {
		if !strings.Contains(compact, contract) {
			t.Errorf("toast forced-colors CSS is missing %q", contract)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesToastContracts(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-toast-region{`,
		`.ui-toast-region .ui-toast.ui-toast-show{`,
		`@media (prefers-reduced-motion:reduce)`,
		`@media (forced-colors:active)`,
		`var(--ui-toast-container)`,
		`var(--ui-toast-action)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled toast CSS is missing %q", contract)
		}
	}
}
