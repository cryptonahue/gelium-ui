package web

import (
	"regexp"
	"strings"
	"testing"
)

// sourceComponentCSS reads a single source CSS file from the same embed that
// styles_contract_test.go exposes, avoiding any mutation of the shared helper.
func sourceComponentCSS(t *testing.T, name string) string {
	t.Helper()
	css, err := sourceStyles.ReadFile("styles/" + name)
	if err != nil {
		t.Fatalf("read source component CSS %s: %v", name, err)
	}
	return string(css)
}

func TestFabPrimitiveCSSMapsMaterialAnatomy(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "fab.css"), " ")

	for _, contract := range []string{
		`.ui-fab {`,
		`position: relative;`,
		`display: inline-flex;`,
		`.ui-fab-medium {`,
		`width: 56px;`,
		`height: 56px;`,
		`box-shadow: var(--ui-shadow-3);`,
		`.ui-fab-small {`,
		`width: 40px;`,
		`height: 40px;`,
		`.ui-fab-large {`,
		`width: 96px;`,
		`height: 96px;`,
		`.ui-fab-extended {`,
		`height: 48px;`,
		`font: var(--ui-type-label-lg);`,
		`.ui-fab-primary { background: var(--ui-fab-primary-container); color: var(--ui-fab-primary-fg); }`,
		`.ui-fab-surface { background: var(--ui-fab-surface-container); color: var(--ui-fab-surface-fg); }`,
		`.ui-fab-secondary { background: var(--ui-fab-secondary-container); color: var(--ui-fab-secondary-fg); }`,
		`.ui-fab:hover:not(:disabled):not([aria-disabled="true"])`,
		`.ui-fab:active:not(:disabled):not([aria-disabled="true"])`,
		`.ui-fab:focus-visible`,
		`outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring);`,
		`.ui-fab:disabled, .ui-fab[aria-disabled="true"]`,
		`opacity: var(--ui-state-disabled-opacity);`,
		`currentColor`,
		`var(--ui-state-hover-opacity)`,
		`var(--ui-state-pressed-opacity)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source FAB CSS is missing contract %q", contract)
		}
	}
}

func TestFabThemeDefinesPublicUITokens(t *testing.T) {
	theme := regexp.MustCompile(`\s+`).ReplaceAllString(repositoryFile(t, "themes", "theme-material", "theme.css"), " ")
	light := theme
	for _, token := range []string{
		"--ui-fab-primary-container:",
		"--ui-fab-primary-fg:",
		"--ui-fab-surface-container:",
		"--ui-fab-surface-fg:",
		"--ui-fab-secondary-container:",
		"--ui-fab-secondary-fg:",
		"--ui-fab-container-shape:",
		"--ui-fab-container-shape-small:",
		"--ui-fab-container-shape-large:",
		"--ui-fab-extended-shape:",
		"--ui-fab-icon-size:",
		"--ui-fab-icon-size-large:",
		"--ui-fab-icon-size-extended:",
		"--ui-fab-extension-gap:",
	} {
		if !strings.Contains(light, token) {
			t.Errorf("theme is missing FAB token %q", token)
		}
	}

	// The dark scheme must remap the container pair so the FAB stays legible.
	// Both the explicit dark override and the media query repeat the pairs.
	count := strings.Count(theme, "--ui-fab-primary-container: #4f378b") +
		strings.Count(theme, "--ui-fab-surface-container: #36343b") +
		strings.Count(theme, "--ui-fab-secondary-container: #4a4458")
	if count < 3 {
		t.Errorf("theme dark scheme must define the three FAB container colors, found dark markers %d", count)
	}
}

func TestFabReducedMotionAndForcedColorsWiredInEntry(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Fatal("entry CSS is missing the reduced-motion media query")
	}
	if !strings.Contains(css, `.ui-fab { transition: none; }`) {
		t.Error("FAB must drop transitions under reduced motion")
	}
	if !strings.Contains(css, `@media (forced-colors: active)`) {
		t.Fatal("entry CSS is missing the forced-colors media query")
	}
	if !strings.Contains(css, `.ui-fab { border: 1px solid CanvasText; }`) {
		t.Error("FAB must keep a visible boundary in forced colors")
	}
}

func TestEmbeddedCompiledCSSIncludesFabContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-fab{`,
		`.ui-fab-medium{`,
		`.ui-fab-small{`,
		`.ui-fab-large{`,
		`.ui-fab-extended{`,
		`var(--ui-fab-primary-container)`,
		`var(--ui-fab-surface-container)`,
		`var(--ui-fab-secondary-container)`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled FAB CSS is missing %q", contract)
		}
	}
}
