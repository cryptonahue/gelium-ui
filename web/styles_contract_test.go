package web

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

//go:embed styles/*.css
var sourceStyles embed.FS

func sourceAppCSS(t *testing.T) string {
	t.Helper()
	css, err := sourceStyles.ReadFile("styles/app.css")
	if err != nil {
		t.Fatalf("read source app CSS: %v", err)
	}
	return string(css)
}

func repositoryFile(t *testing.T, path ...string) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate styles contract test")
	}
	parts := append([]string{filepath.Dir(filename), ".."}, path...)
	content, err := os.ReadFile(filepath.Join(parts...))
	if err != nil {
		t.Fatalf("read repository file: %v", err)
	}
	return string(content)
}

func TestButtonInteractiveStateSelectorsExcludeAriaDisabled(t *testing.T) {
	css := sourceAppCSS(t)
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(css, " ")

	selectors := []string{
		`.ui-button:hover:not(:disabled):not([aria-disabled="true"])`,
		`.ui-button:active:not(:disabled):not([aria-disabled="true"])`,
	}
	for _, selector := range selectors {
		if !strings.Contains(compact, selector+" {") {
			t.Errorf("source CSS is missing inactive-safe selector %q", selector)
		}
	}

	disabledSelector := `.ui-button:disabled, .ui-button[aria-disabled="true"]`
	disabledIndex := strings.Index(compact, disabledSelector+" {")
	if disabledIndex < 0 {
		t.Fatalf("source CSS is missing neutral disabled rule %q", disabledSelector)
	}
	for _, selector := range selectors {
		if stateIndex := strings.Index(compact, selector+" {"); stateIndex >= disabledIndex {
			t.Errorf("neutral disabled rule must follow interactive state selector %q", selector)
		}
	}
}

func TestTextButtonUsesPrimaryStateLayersWithoutInactiveInteraction(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for _, contract := range []string{
		`.ui-button-text { background: transparent; color: var(--ui-color-primary);`,
		`.ui-button-text:hover:not(:disabled):not([aria-disabled="true"]) { box-shadow: inset 0 0 0 999px color-mix(in srgb, var(--ui-color-primary) calc(var(--ui-state-hover-opacity) * 100%), transparent);`,
		`.ui-button-text:active:not(:disabled):not([aria-disabled="true"]) { box-shadow: inset 0 0 0 999px color-mix(in srgb, var(--ui-color-primary) calc(var(--ui-state-pressed-opacity) * 100%), transparent);`,
		`.ui-button:focus-visible { outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring); outline-offset: var(--ui-focus-offset);`,
		`.ui-button:disabled, .ui-button[aria-disabled="true"] { opacity: var(--ui-state-disabled-opacity); cursor: not-allowed; box-shadow: none;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("text button state contract is missing %q", contract)
		}
	}
}

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

func TestTextFieldSourceCSSUsesThemeTokensAndAccessibleStates(t *testing.T) {
	css := sourceAppCSS(t)
	for _, contract := range []string{
		`.ui-text-field-outlined`,
		`.ui-text-field-filled`,
		`.ui-text-field input:hover:not(:disabled)`,
		`.ui-text-field textarea:hover:not(:disabled)`,
		`.ui-text-field input:focus-visible`,
		`.ui-text-field textarea:focus-visible`,
		`.ui-text-field-error`,
		`.ui-text-field :disabled`,
		`var(--ui-field-container)`,
		`var(--ui-field-border)`,
		`var(--ui-field-error)`,
		`@media (prefers-reduced-motion: reduce)`,
		`@media (forced-colors: active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source text-field CSS is missing %q", contract)
		}
	}
}

func TestTextFieldSourceCSSBuildsMaterialContainerAndFloatingLabel(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for _, contract := range []string{
		`.ui-text-field-control { position: relative; display: grid; height: 3.5rem; box-sizing: border-box;`,
		`.ui-text-field-control > label { position: absolute;`,
		`.ui-text-field-control:focus-within > label`,
		`.ui-text-field-control:has(input:not(:placeholder-shown)) > label`,
		`.ui-text-field-control:has(textarea:not(:placeholder-shown)) > label`,
		`.ui-text-field-control:has(textarea:placeholder-shown):not(:focus-within) > label { top: 1rem; transform: none;`,
		`.ui-text-field-outlined .ui-text-field-control { border: 1px solid var(--ui-field-border);`,
		`.ui-text-field-filled .ui-text-field-control { border-bottom: 1px solid var(--ui-field-border);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("Material text-field anatomy is missing %q", contract)
		}
	}
}

func TestTextFieldLabelRemainsClickableAndSignalsTextInteraction(t *testing.T) {
	css := sourceAppCSS(t)
	labelRule := regexp.MustCompile(`(?s)\.ui-text-field-control\s*>\s*label\s*\{([^}]*)\}`).FindStringSubmatch(css)
	if labelRule == nil {
		t.Fatal("source CSS is missing the base text-field label rule")
	}
	if strings.Contains(labelRule[1], "pointer-events: none") {
		t.Error("text-field label must receive pointer events so its for attribute can focus the native control")
	}
	if !regexp.MustCompile(`(?:^|\s)cursor:\s*text\s*;`).MatchString(labelRule[1]) {
		t.Error("clickable text-field label must use cursor: text")
	}
}

func TestTextFieldSingleLineControlKeepsStableExternalHeight(t *testing.T) {
	css := sourceAppCSS(t)
	controlRule := regexp.MustCompile(`(?s)\.ui-text-field-control\s*\{([^}]*)\}`).FindStringSubmatch(css)
	if controlRule == nil {
		t.Fatal("source CSS is missing the base text-field control rule")
	}
	for _, contract := range []string{`(?:^|\s)height:\s*3\.5rem\s*;`, `box-sizing:\s*border-box\s*;`} {
		if !regexp.MustCompile(contract).MatchString(controlRule[1]) {
			t.Errorf("base text-field control must match %q", contract)
		}
	}
	if strings.Contains(controlRule[1], "min-height:") {
		t.Error("base text-field control must use a fixed border-box height, not min-height-only sizing")
	}

	for _, pattern := range []string{
		`(?s)\.ui-text-field input\s*\{[^}]*height:\s*100%\s*;[^}]*box-sizing:\s*border-box\s*;`,
		`(?s)\.ui-text-field-control:has\(textarea\)\s*\{[^}]*height:\s*auto\s*;[^}]*min-height:\s*7rem\s*;`,
		`(?s)\.ui-text-field textarea\s*\{[^}]*min-height:\s*7rem\s*;[^}]*resize:\s*vertical\s*;`,
	} {
		if !regexp.MustCompile(pattern).MatchString(css) {
			t.Errorf("stable text-field sizing contract is missing pattern %q", pattern)
		}
	}
}

func TestTextFieldRestingLabelKeepsBodyLargeLineHeight(t *testing.T) {
	css := sourceAppCSS(t)
	labelRule := regexp.MustCompile(`(?s)\.ui-text-field-control\s*>\s*label\s*\{([^}]*)\}`).FindStringSubmatch(css)
	if labelRule == nil {
		t.Fatal("source CSS is missing the base text-field label rule")
	}
	if !strings.Contains(labelRule[1], "font: var(--ui-type-body-lg);") {
		t.Error("resting text-field label must use --ui-type-body-lg")
	}
	if regexp.MustCompile(`line-height:\s*1\s*;`).MatchString(labelRule[1]) {
		t.Error("resting text-field label must not override the 16/24 body-large shorthand with line-height: 1")
	}
}

func TestTextFieldPreviewUsesAlignedTwoColumnGrid(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for _, contract := range []string{
		`.text-field-preview { display: grid;`,
		`grid-template-columns: repeat(2, minmax(0, 1fr));`,
		`align-items: start;`,
		`gap: 2rem;`,
		`.text-field-preview .ui-text-field { width: 100%;`,
		`@media (max-width: 48rem)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("text-field preview spacing contract is missing %q", contract)
		}
	}
}

func TestTextFieldForcedColorsStylesIntegratedControlWithoutNativeDoubleBorder(t *testing.T) {
	css := sourceAppCSS(t)
	forcedColorsIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedColorsIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forcedColorsCSS := css[forcedColorsIndex:]
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(forcedColorsCSS, " ")
	for _, contract := range []string{
		`.ui-text-field-outlined .ui-text-field-control, .ui-text-field-filled .ui-text-field-control { border-color: CanvasText; forced-color-adjust: auto;`,
		`.ui-text-field-outlined .ui-text-field-control:focus-within { border-color: Highlight; box-shadow: inset 0 0 0 2px Highlight;`,
		`.ui-text-field-filled .ui-text-field-control:focus-within { border-bottom-color: Highlight; box-shadow: inset 0 -2px 0 Highlight;`,
		`.ui-text-field-error .ui-text-field-control { border-color: Mark;`,
		`.ui-text-field-error.ui-text-field-outlined .ui-text-field-control:focus-within { border-color: Mark; box-shadow: inset 0 0 0 2px Mark;`,
		`.ui-text-field-error.ui-text-field-filled .ui-text-field-control:focus-within { border-bottom-color: Mark; box-shadow: inset 0 -2px 0 Mark;`,
		`.ui-text-field-disabled .ui-text-field-control { border-color: GrayText; box-shadow: none;`,
		`.ui-text-field input, .ui-text-field textarea { color: FieldText; forced-color-adjust: auto;`,
	} {
		if !strings.Contains(compact, contract) {
			t.Errorf("forced-colors text-field contract is missing %q", contract)
		}
	}
	inputBorder := regexp.MustCompile(`(?s)\.ui-text-field (?:input|textarea)[^{]*\{[^}]*\bborder(?:-[a-z]+)?:`)
	if inputBorder.MatchString(forcedColorsCSS) {
		t.Error("forced colors must not add a second border to the native input or textarea")
	}
}

// TestTextFieldDisabledPrecedenceOverError declares the explicit state
// precedence contract from the roadmap: when a consumer passes both error and
// disabled classes, disabled must win. The disabled border/background rules must
// follow the error rules in the source and the combined selectors must resolve
// to the disabled palettes rather than the error palette.
func TestTextFieldDisabledPrecedenceOverError(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	errorIndex := strings.LastIndex(css, `border-color: var(--ui-field-error);`)
	disabledIndex := strings.Index(css, `.ui-text-field-disabled .ui-text-field-control { box-shadow: none;`)
	if errorIndex < 0 || disabledIndex < 0 {
		t.Fatal("source CSS is missing the error/disabled ordering baseline")
	}
	if disabledIndex < errorIndex {
		t.Error("disabled rules must be declared after error rules so disabled wins on equal specificity")
	}

	for _, contract := range []string{
		`.ui-text-field-error.ui-text-field-disabled .ui-text-field-control { border-color: color-mix(in srgb, var(--ui-color-fg) 12%, transparent); box-shadow: none;`,
		`.ui-text-field-error.ui-text-field-disabled .ui-text-field-control > label { color: var(--ui-field-label);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("disabled precedence contract is missing combined rule %q", contract)
		}
	}
}

func TestTextFieldDisabledStateAppliesMaterialOpacityPerPart(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for _, contract := range []string{
		`.ui-text-field-disabled .ui-text-field-control > label, .ui-text-field-disabled input, .ui-text-field-disabled textarea, .ui-text-field-disabled .ui-text-field-message, .ui-text-field-disabled .ui-text-field-error-icon { opacity: var(--ui-state-disabled-opacity);`,
		`.ui-text-field-disabled .ui-text-field-control { box-shadow: none;`,
		`.ui-text-field-disabled.ui-text-field-outlined .ui-text-field-control { border-color: color-mix(in srgb, var(--ui-color-fg) 12%, transparent);`,
		`.ui-text-field-disabled.ui-text-field-filled .ui-text-field-control { background: color-mix(in srgb, var(--ui-color-fg) 4%, var(--ui-color-surface)); border-bottom-color: color-mix(in srgb, var(--ui-color-fg) 38%, transparent);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("disabled Material text-field contract is missing %q", contract)
		}
	}
	if strings.Contains(css, `.ui-text-field :disabled { opacity:`) {
		t.Error("disabled opacity must cover every component part, not only the native control")
	}
}

func TestTextFieldTrailingSlotAndFloatingLabelUseLogicalRTLGeometry(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for _, contract := range []string{
		`.ui-text-field-control > label { position: absolute; z-index: 1; top: 50%; inset-inline-start: 1rem;`,
		`.ui-text-field-error .ui-text-field-control > label { max-width: calc(100% - 3.5rem);`,
		`.ui-text-field-error-icon { position: absolute; top: 50%; inset-inline-end: 1rem; width: 1.5rem; height: 1.5rem;`,
		`.ui-text-field-error input, .ui-text-field-error textarea { padding-inline-end: 3.5rem;`,
		`.ui-text-field-control:dir(rtl) > label { transform-origin: right top;`,
		`.ui-text-field-outlined .ui-text-field-control:dir(rtl):focus-within > label`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("logical trailing-slot/RTL contract is missing %q", contract)
		}
	}
	for _, physical := range []string{"left: 1rem;", "right: 1rem;", "padding-right: 3.5rem;"} {
		if strings.Contains(css, physical) {
			t.Errorf("text-field CSS must not use physical trailing geometry %q", physical)
		}
	}
}

func TestTextFieldHoverDoesNotOverrideFocusIndicator(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for _, control := range []string{"input", "textarea"} {
		selector := `.ui-text-field-control:hover:not(:focus-within):has(` + control + `:not(:disabled))`
		if !strings.Contains(css, selector) {
			t.Errorf("hover selector for %s must exclude focused fields with %q", control, selector)
		}
	}

	interactiveHover := regexp.MustCompile(`\.ui-text-field-control:hover[^,{]*:has\((?:input|textarea):not\(:disabled\)\)`)
	for _, selector := range interactiveHover.FindAllString(css, -1) {
		if !strings.Contains(selector, ":not(:focus-within)") {
			t.Errorf("interactive text-field hover selector must not apply during focus: %q", selector)
		}
	}
}

func TestTextFieldErrorFocusAndHoverKeepErrorIndicatorForBothVariants(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for _, variant := range []string{"outlined", "filled"} {
		selectors := []string{
			`.ui-text-field-error.ui-text-field-` + variant + ` .ui-text-field-control:focus-within`,
			`.ui-text-field-error.ui-text-field-` + variant + ` .ui-text-field-control:hover:not(:focus-within):has(input:not(:disabled))`,
			`.ui-text-field-error.ui-text-field-` + variant + ` .ui-text-field-control:hover:not(:focus-within):has(textarea:not(:disabled))`,
		}
		for _, selector := range selectors {
			pattern := regexp.MustCompile(regexp.QuoteMeta(selector) + `(?:\s*,[^{}]*)?\s*\{[^}]*border-color:\s*var\(--ui-field-error\)\s*;?`)
			if !pattern.MatchString(css) {
				t.Errorf("error %s field selector %q must keep its error indicator via border-color", variant, selector)
			}
		}
	}

	for _, test := range []struct {
		variant string
		shadow  string
	}{
		{"outlined", `box-shadow: inset 0 0 0 2px var(--ui-field-error);`},
		{"filled", `box-shadow: inset 0 -2px 0 var(--ui-field-error);`},
	} {
		selector := `.ui-text-field-error.ui-text-field-` + test.variant + ` .ui-text-field-control:focus-within`
		pattern := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(selector) + `(?:\s*,[^{}]*)?\s*\{([^}]*)\}`)
		rule := pattern.FindStringSubmatch(sourceAppCSS(t))
		if rule == nil || !strings.Contains(rule[1], test.shadow) {
			t.Errorf("focused error %s field must replace the primary overlay with %q", test.variant, test.shadow)
		}
	}

	errorRuleIndex := strings.LastIndex(css, `border-color: var(--ui-field-error)`)
	for _, normalState := range []string{
		`.ui-text-field-filled .ui-text-field-control:focus-within`,
		`.ui-text-field-control:hover:not(:focus-within):has(input:not(:disabled))`,
	} {
		if stateIndex := strings.Index(css, normalState); stateIndex < 0 || errorRuleIndex < stateIndex {
			t.Errorf("error-state indicator override must follow normal state rule %q", normalState)
		}
	}
}

func TestTextFieldSupportingMessageUsesMaterialTypographyAndSpacing(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for _, contract := range []string{
		`.ui-text-field { display: grid; width: min(100%, 24rem); gap: 0;`,
		`.ui-text-field-message { margin: 0; padding: .25rem 1rem 0; color: var(--ui-color-fg-muted); font: var(--ui-type-body-sm);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("supporting message Material contract is missing %q", contract)
		}
	}
}

func TestTextFieldFocusUsesInsetOverlayWithoutChangingBorderGeometry(t *testing.T) {
	css := sourceAppCSS(t)
	for _, test := range []struct {
		name     string
		selector string
		color    string
		shadow   string
	}{
		{"outlined", `.ui-text-field-outlined .ui-text-field-control:focus-within`, `border-color: var(--ui-color-primary);`, `box-shadow: inset 0 0 0 2px var(--ui-color-primary);`},
		{"filled", `.ui-text-field-filled .ui-text-field-control:focus-within`, `border-bottom-color: var(--ui-color-primary);`, `box-shadow: inset 0 -2px 0 var(--ui-color-primary);`},
	} {
		t.Run(test.name, func(t *testing.T) {
			pattern := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(test.selector) + `\s*\{([^}]*)\}`)
			rule := pattern.FindStringSubmatch(css)
			if rule == nil {
				t.Fatalf("source CSS is missing focus rule %q", test.selector)
			}
			for _, contract := range []string{test.color, test.shadow} {
				if !strings.Contains(rule[1], contract) {
					t.Errorf("%s focus rule is missing %q", test.name, contract)
				}
			}
			if regexp.MustCompile(`border(?:-[a-z]+)?-width\s*:`).MatchString(rule[1]) {
				t.Errorf("%s focus must preserve the base 1px border geometry, got %q", test.name, rule[1])
			}
		})
	}

	compact := regexp.MustCompile(`\s+`).ReplaceAllString(css, " ")
	for _, contract := range []string{
		`.ui-text-field-control:focus-within > label { color: var(--ui-color-primary);`,
		`.ui-text-field-error .ui-text-field-control:focus-within > label { color: var(--ui-field-error);`,
	} {
		if !strings.Contains(compact, contract) {
			t.Errorf("integrated Material focus contract is missing %q", contract)
		}
	}
	if strings.Contains(compact, `.ui-text-field-control:has(:focus-visible) { outline:`) {
		t.Error("text field focus must not add an exterior focus outline")
	}
	if strings.Contains(compact, `.ui-text-field-control:has(input:not(:placeholder-shown)) > label, .ui-text-field-control:has(textarea:not(:placeholder-shown)) > label { color: var(--ui-color-primary);`) {
		t.Error("an unfocused populated field label must retain its normal label color")
	}
}

func TestReducedMotionDisablesTextFieldControlAndLabelTransitions(t *testing.T) {
	css := sourceAppCSS(t)
	reducedMotionIndex := strings.Index(css, "@media (prefers-reduced-motion: reduce)")
	if reducedMotionIndex < 0 {
		t.Fatal("source CSS is missing the reduced-motion media query")
	}
	reducedMotionCSS := css[reducedMotionIndex:]
	if nextMedia := strings.Index(reducedMotionCSS[1:], "@media "); nextMedia >= 0 {
		reducedMotionCSS = reducedMotionCSS[:nextMedia+1]
	}
	for _, selector := range []string{`.ui-text-field-control`, `.ui-text-field-control > label`} {
		pattern := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(selector) + `\s*\{[^}]*transition:\s*none\s*;?[^}]*\}`)
		if !pattern.MatchString(reducedMotionCSS) {
			t.Errorf("reduced-motion CSS must disable transitions for %s", selector)
		}
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

func TestReducedMotionDisablesButtonSpinnerAnimation(t *testing.T) {
	css := sourceAppCSS(t)
	reducedMotionIndex := strings.Index(css, "@media (prefers-reduced-motion: reduce)")
	if reducedMotionIndex < 0 {
		t.Fatal("source CSS is missing the reduced-motion media query")
	}

	reducedMotionCSS := css[reducedMotionIndex:]
	if nextMedia := strings.Index(reducedMotionCSS[1:], "@media "); nextMedia >= 0 {
		reducedMotionCSS = reducedMotionCSS[:nextMedia+1]
	}
	spinnerAnimationNone := regexp.MustCompile(`(?s)\.ui-button-spinner\s*\{[^}]*animation:\s*none\s*;?[^}]*\}`)
	if !spinnerAnimationNone.MatchString(reducedMotionCSS) {
		t.Error("reduced-motion CSS must disable the spinner with animation: none")
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

func TestDividerPrimitiveCSSMapsTokensAndInsets(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-divider {`,
		`height: var(--ui-divider-thickness);`,
		`background: var(--ui-divider-color);`,
		`.ui-divider-inset { padding-inline: 1rem;`,
		`.ui-divider-inset-start { padding-inline-start: 1rem;`,
		`.ui-divider-inset-end { padding-inline-end: 1rem;`,
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
	theme := repositoryFile(t, "themes", "theme-material", "theme.css")
	for _, token := range []string{
		"--ui-divider-color:",
		"--ui-divider-thickness:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing divider token %q", token)
		}
	}
}

func TestElevationPrimitiveCSSMapsLevelsToShadowTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for level := 0; level <= 5; level++ {
		contract := fmt.Sprintf(".ui-elevation-%d { box-shadow: var(--ui-shadow-%d);", level, level)
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing %q", contract)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	if !strings.Contains(forced, ".ui-elevation-0, .ui-elevation-1, .ui-elevation-2, .ui-elevation-3, .ui-elevation-4, .ui-elevation-5 { box-shadow: none;") {
		t.Error("elevation is visual-only and must not keep decorative shadows in forced colors")
	}

	reducedIndex := strings.Index(css, "@media (prefers-reduced-motion: reduce)")
	if reducedIndex < 0 {
		t.Fatal("source CSS is missing the reduced-motion media query")
	}
	reduced := css[reducedIndex:]
	if nextMedia := strings.Index(reduced[1:], "@media "); nextMedia >= 0 {
		reduced = reduced[:nextMedia+1]
	}
	if !strings.Contains(reduced, ".ui-elevation-0, .ui-elevation-1, .ui-elevation-2, .ui-elevation-3, .ui-elevation-4, .ui-elevation-5 { transition: none;") {
		t.Error("reduced-motion CSS must disable the elevation box-shadow transition")
	}
}

func TestEmbeddedCompiledCSSIncludesElevationContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-elevation-0{box-shadow:var(--ui-shadow-0)}`,
		`.ui-elevation-1{box-shadow:var(--ui-shadow-1)}`,
		`.ui-elevation-5{box-shadow:var(--ui-shadow-5)}`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled elevation CSS is missing %q", contract)
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

func TestCardPrimitiveCSSMapsVariantsToThemeTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-card {`,
		`border-radius: var(--ui-card-radius);`,
		`.ui-card-elevated { background: var(--ui-card-container-elevated); box-shadow: var(--ui-shadow-1);`,
		`.ui-card-filled { background: var(--ui-card-container-filled);`,
		`.ui-card-outlined { background: var(--ui-card-container-outlined); border: 1px solid var(--ui-card-outline-color);`,
		`.ui-card:focus-visible { outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring); outline-offset: var(--ui-focus-offset);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing card contract %q", contract)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	if !strings.Contains(forced, ".ui-card { border: 1px solid CanvasText;") {
		t.Error("card must keep a visible boundary in forced colors")
	}
}

func TestCardThemeDefinesPublicUIFamily(t *testing.T) {
	theme := repositoryFile(t, "themes", "theme-material", "theme.css")
	for _, token := range []string{
		"--ui-card-radius:",
		"--ui-card-container-elevated:",
		"--ui-card-container-filled:",
		"--ui-card-container-outlined:",
		"--ui-card-outline-color:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing card token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesCardContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-card{`,
		`.ui-card-elevated{`,
		`.ui-card-filled{`,
		`.ui-card-outlined{`,
		`@media (forced-colors:active)`,
		`var(--ui-card-container-filled)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled card CSS is missing %q", contract)
		}
	}
}

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
	theme := repositoryFile(t, "themes", "theme-material", "theme.css")
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

func TestCheckboxPrimitiveCSSMapsNativeControlAndStates(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-checkbox {`,
		`gap: .5rem;`,
		`input[type="checkbox"] {`,
		`appearance: none;`,
		`width: var(--ui-checkbox-size);`,
		`height: var(--ui-checkbox-size);`,
		`border-radius: var(--ui-checkbox-radius);`,
		`border: var(--ui-checkbox-outline-width) solid var(--ui-checkbox-outline);`,
		`input:checked + .ui-checkbox-mark {`,
		`input:checked ~ .ui-checkbox-label {`,
		`background: var(--ui-checkbox-container);`,
		`input:disabled {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing checkbox contract %q", contract)
		}
	}

	stateLayers := []string{
		`.ui-checkbox:hover {`,
		`:focus-visible`,
		`:active:not(:disabled)`,
	}
	for _, sel := range stateLayers {
		if !strings.Contains(css, sel) {
			t.Errorf("source CSS is missing checkbox state selector %q", sel)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	for _, contract := range []string{
		`border-color: CanvasText;`,
		`border-color: Canvas;`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("checkbox must stay operable in forced colors; missing %q", contract)
		}
	}
}

func TestCheckboxThemeDefinesPublicUIFamily(t *testing.T) {
	theme := repositoryFile(t, "themes", "theme-material", "theme.css")
	for _, token := range []string{
		"--ui-checkbox-size:",
		"--ui-checkbox-radius:",
		"--ui-checkbox-outline-width:",
		"--ui-checkbox-outline:",
		"--ui-checkbox-hover-outline:",
		"--ui-checkbox-container:",
		"--ui-checkbox-icon:",
		"--ui-checkbox-error:",
		"--ui-checkbox-checked-disabled-container:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing checkbox token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesCheckboxContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-checkbox`,
		`var(--ui-checkbox-size)`,
		`var(--ui-checkbox-outline)`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled checkbox CSS is missing %q", contract)
		}
	}
}

func TestRadioPrimitiveCSSMapsNativeControlAndStates(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-radio {`,
		`gap: .5rem;`,
		`input[type="radio"] {`,
		`appearance: none;`,
		`width: var(--ui-radio-size);`,
		`height: var(--ui-radio-size);`,
		`border-radius: var(--ui-radio-radius);`,
		`border: var(--ui-radio-outline-width) solid var(--ui-radio-outline);`,
		`input:checked + .ui-radio-mark {`,
		`background: var(--ui-radio-checked);`,
		`input:disabled {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing radio contract %q", contract)
		}
	}

	stateLayers := []string{
		`.ui-radio:hover {`,
		`:focus-visible`,
		`:active:not(:disabled)`,
	}
	for _, sel := range stateLayers {
		if !strings.Contains(css, sel) {
			t.Errorf("source CSS is missing radio state selector %q", sel)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	for _, contract := range []string{
		`border-color: CanvasText;`,
		`background: CanvasText;`,
		`background: GrayText;`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("radio must stay operable in forced colors; missing %q", contract)
		}
	}
}

func TestRadioThemeDefinesPublicUIFamily(t *testing.T) {
	theme := repositoryFile(t, "themes", "theme-material", "theme.css")
	for _, token := range []string{
		"--ui-radio-size:",
		"--ui-radio-radius:",
		"--ui-radio-outline-width:",
		"--ui-radio-outline:",
		"--ui-radio-hover-outline:",
		"--ui-radio-checked:",
		"--ui-radio-disabled:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing radio token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesRadioContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-radio`,
		`var(--ui-radio-size)`,
		`var(--ui-radio-outline)`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded radio CSS is missing %q", contract)
		}
	}
}

func TestSwitchPrimitiveCSSMapsNativeControlAndStates(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-switch {`,
		`gap: .5rem;`,
		`input[type="checkbox"] {`,
		`appearance: none;`,
		`width: var(--ui-switch-width);`,
		`height: var(--ui-switch-height);`,
		`border-radius: var(--ui-switch-radius);`,
		`border: var(--ui-switch-outline-width) solid var(--ui-switch-track-outline);`,
		`input:checked + .ui-switch-track {`,
		`background: var(--ui-switch-track);`,
		`input:checked ~ .ui-switch-handle {`,
		`translateX`,
		`input:disabled {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing switch contract %q", contract)
		}
	}

	stateLayers := []string{
		`.ui-switch:hover {`,
		`:focus-visible`,
		`:active:not(:disabled)`,
	}
	for _, sel := range stateLayers {
		if !strings.Contains(css, sel) {
			t.Errorf("source CSS is missing switch state selector %q", sel)
		}
	}
	if strings.Contains(css, `input[type="checkbox"] { appearance: none; opacity: 0;`) {
		t.Error("switch must keep the native checkbox visually integrated as the track, not hide it with opacity")
	}
	if strings.Contains(css, `.ui-switch input[type="checkbox"]:disabled { opacity: 0.38;`) {
		t.Error("switch disabled opacity must follow Material track 0.12 / handle 0.38 split, not dim the whole input")
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	for _, contract := range []string{
		`.ui-switch input:checked + .ui-switch-track { background: ButtonText;`,
		`.ui-switch input:checked ~ .ui-switch-handle { background: ButtonText;`,
		`.ui-switch input:disabled + .ui-switch-track { border-color: GrayText;`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("switch must stay operable in forced colors; missing %q", contract)
		}
	}
}

func TestSwitchThemeDefinesPublicUIFamily(t *testing.T) {
	theme := repositoryFile(t, "themes", "theme-material", "theme.css")
	for _, token := range []string{
		"--ui-switch-width:",
		"--ui-switch-height:",
		"--ui-switch-radius:",
		"--ui-switch-outline-width:",
		"--ui-switch-track:",
		"--ui-switch-track-unselected:",
		"--ui-switch-track-outline:",
		"--ui-switch-handle:",
		"--ui-switch-handle-selected:",
		"--ui-switch-handle-size:",
		"--ui-switch-handle-selected-size:",
		"--ui-switch-handle-pressed-size:",
		"--ui-switch-disabled-track-opacity:",
		"--ui-switch-disabled-handle-opacity:",
		"--ui-switch-disabled-handle:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing switch token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesSwitchContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-switch`,
		`var(--ui-switch-width)`,
		`var(--ui-switch-track-outline)`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled switch CSS is missing %q", contract)
		}
	}
}

func TestSelectPrimitiveCSSMapsNativeVariantsAndStates(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-select {`,
		`gap: 0;`,
		`.ui-select select {`,
		`appearance: none;`,
		`height: var(--ui-select-height);`,
		`border-radius: var(--ui-select-radius);`,
		`.ui-select-filled select { background: var(--ui-select-container-filled); border: 1px solid transparent; border-bottom: 1px solid var(--ui-select-outline);`,
		`.ui-select-outlined select { background: transparent; border: 1px solid var(--ui-select-outline);`,
		`.ui-select select:focus-visible { outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring); outline-offset: var(--ui-focus-offset);`,
		`.ui-select select:disabled { cursor: not-allowed;`,
		`.ui-select select[aria-invalid="true"] { border-color: var(--ui-select-error);`,
		`.ui-select-caret`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing select contract %q", contract)
		}
	}

	labelFloating := []string{
		`.ui-select:focus-within .ui-select-label`,
		`.ui-select:has(select option:checked:not([value=""])) .ui-select-label`,
	}
	for _, sel := range labelFloating {
		if !strings.Contains(css, sel) {
			t.Errorf("source CSS is missing floating-label selector %q", sel)
		}
	}
	if strings.Contains(css, `.ui-select select { appearance: none; opacity: 0;`) {
		t.Error("select must keep the native select visible as the control surface, not hide it")
	}
	if !strings.Contains(css, `.ui-select select:not(:has(option:checked:not([value=""]))) { color: transparent;`) {
		t.Error("source CSS must hide the placeholder option text while empty so the resting label is the only prompt")
	}
	if !strings.Contains(css, `.ui-select select option:checked { background: var(--ui-select-list-bg); color: var(--ui-select-list-fg);`) {
		t.Error("source CSS must style the checked option row with the list palette so the browser popup does not tint the selected placeholder")
	}
	labelRule := regexp.MustCompile(`(?s)\.ui-select-label\s*\{([^}]*)\}`).FindStringSubmatch(css)
	if labelRule == nil {
		t.Fatal("source CSS is missing the base select label rule")
	}
	if strings.Contains(labelRule[1], "pointer-events: none") {
		t.Error("select label must receive pointer events so its for attribute can focus the native select")
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	for _, contract := range []string{
		`.ui-select select { border-color: CanvasText;`,
		`.ui-select select:disabled { border-color: GrayText;`,
		`.ui-select select[aria-invalid="true"] { border-color: Mark;`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("select must stay distinguishable in forced colors; missing %q", contract)
		}
	}
}

func TestSelectThemeDefinesPublicUIFamily(t *testing.T) {
	theme := repositoryFile(t, "themes", "theme-material", "theme.css")
	for _, token := range []string{
		"--ui-select-height:",
		"--ui-select-radius:",
		"--ui-select-radius-top:",
		"--ui-select-caret:",
		"--ui-select-fg:",
		"--ui-select-label:",
		"--ui-select-outline:",
		"--ui-select-container-filled:",
		"--ui-select-hover-outline:",
		"--ui-select-focus:",
		"--ui-select-error:",
		"--ui-select-disabled-opacity:",
		"--ui-select-list-bg:",
		"--ui-select-list-fg:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing select token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesSelectContracts(t *testing.T) {
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
}

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

	reducedIndex := strings.Index(css, "@media (prefers-reduced-motion: reduce)")
	if reducedIndex < 0 {
		t.Fatal("source CSS is missing the reduced-motion media query")
	}
	reduced := css[reducedIndex:]
	if nextMedia := strings.Index(reduced[1:], "@media "); nextMedia >= 0 {
		reduced = reduced[:nextMedia+1]
	}
	if !strings.Contains(reduced, ".ui-select-menu") {
		t.Error("reduced-motion CSS must disable select-menu transitions")
	}
}

func TestSelectMenuThemeDefinesPublicUIFamily(t *testing.T) {
	theme := repositoryFile(t, "themes", "theme-material", "theme.css")
	for _, token := range []string{
		"--ui-select-menu-container:",
		"--ui-select-menu-radius:",
		"--ui-select-menu-elevation:",
		"--ui-select-menu-min-width:",
		"--ui-select-menu-item-height:",
		"--ui-select-menu-item-fg:",
		"--ui-select-menu-item-icon:",
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

func TestSliderPrimitiveCSSMapsNativeRangeAndStates(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-slider {`,
		`input[type="range"] {`,
		`appearance: none;`,
		`::-webkit-slider-runnable-track {`,
		`var(--ui-slider-fill`,
		`::-webkit-slider-thumb {`,
		`::-moz-range-track {`,
		`::-moz-range-progress {`,
		`::-moz-range-thumb {`,
		`var(--ui-slider-handle-pressed-size)`,
		`input[type="range"]:disabled {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing slider contract %q", contract)
		}
	}

	for _, sel := range []string{`:focus-visible`, `:active:not(:disabled)`} {
		if !strings.Contains(css, sel) {
			t.Errorf("source CSS is missing slider state selector %q", sel)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	for _, contract := range []string{
		`::-webkit-slider-thumb`,
		`::-moz-range-thumb`,
		`::-webkit-slider-runnable-track`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("slider must stay distinguishable in forced colors; missing %q", contract)
		}
	}
}

func TestSliderThemeDefinesPublicUIFamily(t *testing.T) {
	theme := repositoryFile(t, "themes", "theme-material", "theme.css")
	for _, token := range []string{
		"--ui-slider-track-height:",
		"--ui-slider-track-radius:",
		"--ui-slider-handle-size:",
		"--ui-slider-handle-pressed-size:",
		"--ui-slider-active:",
		"--ui-slider-inactive:",
		"--ui-slider-handle:",
		"--ui-slider-handle-elevation:",
		"--ui-slider-disabled:",
		"--ui-slider-disabled-opacity:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing slider token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesSliderContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-slider`,
		`-webkit-slider-runnable-track`,
		`-moz-range-progress`,
		`-moz-range-thumb`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled slider CSS is missing %q", contract)
		}
	}
}

func TestProgressPrimitiveCSSMapsNativeProgressAndStates(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-progress {`,
		`progress {`,
		`appearance: none;`,
		`height: var(--ui-progress-track-height);`,
		`border-radius: var(--ui-progress-radius);`,
		`::-webkit-progress-bar {`,
		`::-webkit-progress-value {`,
		`var(--ui-progress-indicator)`,
		`::-moz-progress-bar {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source CSS is missing progress contract %q", contract)
		}
	}

	forcedIndex := strings.Index(css, "@media (forced-colors: active)")
	if forcedIndex < 0 {
		t.Fatal("source CSS is missing the forced-colors media query")
	}
	forced := css[forcedIndex:]
	for _, contract := range []string{
		`::-webkit-progress-value`,
		`::-moz-progress-bar`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("progress must stay distinguishable in forced colors; missing %q", contract)
		}
	}
}

func TestProgressThemeDefinesPublicUIFamily(t *testing.T) {
	theme := repositoryFile(t, "themes", "theme-material", "theme.css")
	for _, token := range []string{
		"--ui-progress-track-height:",
		"--ui-progress-radius:",
		"--ui-progress-track:",
		"--ui-progress-indicator:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme is missing progress token %q", token)
		}
	}
}

func TestEmbeddedCompiledCSSIncludesProgressContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-progress`,
		`-webkit-progress-value`,
		`-moz-progress-bar`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled progress CSS is missing %q", contract)
		}
	}
}
