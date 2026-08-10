package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestMaterialThemeDefinesDialogSemanticsInEveryColorScheme(t *testing.T) {
	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, "theme-material"), " ")
	for _, contract := range []string{
		`--ui-dialog-radius:`,
		`--ui-type-dialog-headline:`,
		`--ui-type-dialog-body:`,
		`--ui-dialog-container:`,
		`--ui-dialog-fg:`,
		`--ui-dialog-body:`,
		`--ui-dialog-scrim:`,
		`.theme-material.theme-dark,`,
		`@media (prefers-color-scheme: dark)`,
	} {
		if !strings.Contains(theme, contract) {
			t.Errorf("Material dialog theme contract is missing %q", contract)
		}
	}
	// Every surface/foreground/scrim token must be defined once per scheme:
	// light, explicit dark, and media dark. The contract is token presence
	// across schemes, not a concrete hex value.
	for _, token := range []string{
		"--ui-dialog-container:",
		"--ui-dialog-fg:",
		"--ui-dialog-body:",
		"--ui-dialog-scrim:",
	} {
		if got := strings.Count(theme, token); got != 3 {
			t.Errorf("dark dialog semantics must be defined for explicit and preferred dark schemes: %s appears %d times, want 3", token, got)
		}
	}
}

func TestDialogSourceCSSImplementsMaterialGeometryStatesAndProgressiveMotion(t *testing.T) {
	css := sourceAppCSS(t)
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(css, " ")
	for _, contract := range []string{
		`.ui-dialog {`, `min-width: var(--ui-dialog-min-width);`, `min-height: var(--ui-dialog-min-height);`,
		`max-width: min(var(--ui-dialog-max-width), calc(100% - 48px));`, `max-height: min(var(--ui-dialog-max-width), calc(100% - 48px));`,
		`width: fit-content;`, `height: fit-content;`, `margin: auto;`, `border-radius: var(--ui-dialog-radius);`,
		`background: var(--ui-dialog-container);`, `color: var(--ui-dialog-fg);`,
		`.ui-dialog-headline { margin: 0; padding: var(--ui-space-6) var(--ui-space-6) 0; font: var(--ui-type-dialog-headline);`,
		`.ui-dialog-content { padding: var(--ui-space-6); color: var(--ui-dialog-body); font: var(--ui-type-dialog-body);`,
		`.ui-dialog-actions { display: flex; flex-wrap: nowrap; justify-content: flex-end; gap: var(--ui-space-2); padding: var(--ui-space-4) var(--ui-space-6) var(--ui-space-6);`,
		`.ui-dialog::backdrop { background: var(--ui-dialog-scrim);`,
		`transition: opacity var(--ui-motion-short)`, `transition-behavior: allow-discrete;`, `overlay var(--ui-motion-short)`, `display var(--ui-motion-short)`,
		`.ui-dialog[open] {`, `translate: 0;`, `scale: 1;`, `opacity: 1;`,
		`@starting-style`, `translate: 0 -50px;`, `scale: .35;`, `var(--ui-motion-long)`,
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
