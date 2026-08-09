package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestMaterialThemeDefinesDialogSemanticsInEveryColorScheme(t *testing.T) {
	theme := regexp.MustCompile(`\s+`).ReplaceAllString(repositoryFile(t, "themes", "theme-material", "theme.css"), " ")
	for _, contract := range []string{
		`--ui-dialog-radius: 28px;`,
		`--ui-type-dialog-headline: 400 1.5rem/2rem var(--ui-font-sans);`,
		`--ui-type-dialog-body: 400 .875rem/1.25rem var(--ui-font-sans);`,
		`--ui-dialog-container: #ece6f0;`,
		`--ui-dialog-fg: #1d1b20;`,
		`--ui-dialog-body: #49454f;`,
		`--ui-dialog-scrim: rgb(0 0 0 / .32);`,
		`.theme-material.theme-dark,`,
		`--ui-dialog-container: #2b2930;`,
		`--ui-dialog-fg: #e6e0e9;`,
		`--ui-dialog-body: #cac4d0;`,
		`@media (prefers-color-scheme: dark)`,
	} {
		if !strings.Contains(theme, contract) {
			t.Errorf("Material dialog theme contract is missing %q", contract)
		}
	}
	if strings.Count(theme, `--ui-dialog-container: #2b2930;`) != 2 {
		t.Error("dark dialog semantics must be defined for explicit and preferred dark schemes")
	}
}

func TestDialogSourceCSSImplementsMaterialGeometryStatesAndProgressiveMotion(t *testing.T) {
	css := sourceAppCSS(t)
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(css, " ")
	for _, contract := range []string{
		`.ui-dialog {`, `min-width: 280px;`, `min-height: 140px;`,
		`max-width: min(560px, calc(100% - 48px));`, `max-height: min(560px, calc(100% - 48px));`,
		`width: fit-content;`, `height: fit-content;`, `margin: auto;`, `border-radius: var(--ui-dialog-radius);`,
		`background: var(--ui-dialog-container);`, `color: var(--ui-dialog-fg);`,
		`.ui-dialog-headline { margin: 0; padding: 24px 24px 0; font: var(--ui-type-dialog-headline);`,
		`.ui-dialog-content { padding: 24px; color: var(--ui-dialog-body); font: var(--ui-type-dialog-body);`,
		`.ui-dialog-actions { display: flex; flex-wrap: nowrap; justify-content: flex-end; gap: 8px; padding: 16px 24px 24px;`,
		`.ui-dialog::backdrop { background: var(--ui-dialog-scrim);`,
		`transition: opacity 150ms`, `transition-behavior: allow-discrete;`, `overlay 150ms`, `display 150ms`,
		`.ui-dialog[open] {`, `translate: 0;`, `scale: 1;`, `opacity: 1;`,
		`@starting-style`, `translate: 0 -50px;`, `scale: .35;`, `500ms`,
		`@media (prefers-reduced-motion: reduce)`, `transition: none;`,
		`@media (forced-colors: active)`, `border: 2px solid WindowText;`,
	} {
		if !strings.Contains(compact, contract) {
			t.Errorf("dialog CSS is missing %q", contract)
		}
	}
	dialogRule := regexp.MustCompile(`(?s)\.ui-dialog\s*\{([^}]*)\}`).FindStringSubmatch(css)
	if dialogRule == nil {
		t.Fatal("source CSS is missing .ui-dialog rule")
	}
	if strings.Contains(dialogRule[1], "box-shadow") {
		t.Error("dialog container must not add box-shadow/elevation")
	}
}

func TestEmbeddedCompiledCSSIncludesDialogContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}

	css := string(compiled)
	for _, contract := range []string{
		`.ui-dialog{`,
		`.ui-dialog::backdrop{`,
		`@starting-style`,
		`@media (prefers-reduced-motion:reduce)`,
		`@media (forced-colors:active)`,
		`var(--ui-dialog-container)`,
		`var(--ui-dialog-scrim)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled dialog CSS is missing %q", contract)
		}
	}
}
