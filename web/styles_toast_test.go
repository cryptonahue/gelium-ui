package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestToastReleaseDocsAndPackageVersionStayCoherent(t *testing.T) {
	packageJSON := repositoryFile(t, "package.json")
	readme := repositoryFile(t, "README.md")
	docs := repositoryFile(t, "web", "content", "toast.md")
	if !strings.Contains(packageJSON, `"version": "0.4.0"`) {
		t.Error("package version must identify the Toast release as 0.4.0")
	}
	for _, contract := range []string{"/components/toast", "v0.4.0", "loom:toast", "no-JS"} {
		if !strings.Contains(readme, contract) {
			t.Errorf("README is missing Toast release contract %q", contract)
		}
	}
	for _, contract := range []string{"role=\"alert\"", "role=\"status\"", "aria-live", "HX-Trigger", "prefers-reduced-motion", "forced-colors"} {
		if !strings.Contains(docs, contract) {
			t.Errorf("Toast documentation is missing %q", contract)
		}
	}
}

func TestMaterialThemeDefinesToastTokensInEveryColorScheme(t *testing.T) {
	theme := regexp.MustCompile(`\s+`).ReplaceAllString(repositoryFile(t, "themes", "theme-material", "theme.css"), " ")
	count := func(token string) int {
		return len(regexp.MustCompile(`(?i)`+regexp.QuoteMeta(token)).FindAllStringIndex(theme, -1))
	}
	if got := count(`--ui-toast-container:`); got != 3 {
		t.Errorf("toast container theme declarations = %d, want 3 (light, explicit dark, media dark)", got)
	}
	for _, token := range []string{
		`--ui-toast-container: #322f35;`,
		`--ui-toast-fg: #f3edf7;`,
		`--ui-toast-radius: 4px;`,
		`--ui-toast-action: #d0bcff;`,
		`--ui-toast-container: #ece6f0;`,
		`--ui-toast-fg: #1d1b20;`,
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
		`.ui-toast { display: flex;`,
		`min-height: 3rem;`,
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
	reducedIndex := strings.Index(css, "@media (prefers-reduced-motion: reduce)")
	if reducedIndex < 0 {
		t.Fatal("source CSS is missing the reduced-motion media query")
	}
	reduced := css[reducedIndex:]
	if next := strings.Index(reduced[1:], "@media "); next >= 0 {
		reduced = reduced[:next+1]
	}
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
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
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
