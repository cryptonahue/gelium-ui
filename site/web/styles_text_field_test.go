package web

import (
	"regexp"
	"strings"
	"testing"
)

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
		`.ui-text-field-control { position: relative; display: grid; height: var(--ui-size-field); box-sizing: border-box;`,
		`.ui-text-field-control > label { position: absolute;`,
		`.ui-text-field-control:focus-within > label`,
		`.ui-text-field-control:has(input:not(:placeholder-shown)) > label`,
		`.ui-text-field-control:has(textarea:not(:placeholder-shown)) > label`,
		`.ui-text-field-control:has(textarea:placeholder-shown):not(:focus-within) > label { top: 1rem; transform: none;`,
		`.ui-text-field-outlined .ui-text-field-control { border: var(--ui-border-width-1) var(--ui-border-style-solid) var(--ui-field-border);`,
		`.ui-text-field-filled .ui-text-field-control { border-bottom: var(--ui-border-width-1) var(--ui-border-style-solid) var(--ui-field-border);`,
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
	for _, contract := range []string{`(?:^|\s)height:\s*var\(--ui-size-field\)\s*;`, `box-sizing:\s*border-box\s*;`} {
		if !regexp.MustCompile(contract).MatchString(controlRule[1]) {
			t.Errorf("base text-field control must match %q", contract)
		}
	}
	if strings.Contains(controlRule[1], "min-height:") {
		t.Error("base text-field control must use a fixed border-box height, not min-height-only sizing")
	}

	for _, pattern := range []string{
		`(?s)\.ui-text-field input\s*\{[^}]*height:\s*100%\s*;[^}]*box-sizing:\s*border-box\s*;`,
		`(?s)\.ui-text-field-control:has\(textarea\)\s*\{[^}]*height:\s*auto\s*;[^}]*min-height:\s*var\(--ui-text-field-textarea-min-height\)\s*;`,
		`(?s)\.ui-text-field textarea\s*\{[^}]*min-height:\s*var\(--ui-text-field-textarea-min-height\)\s*;[^}]*resize:\s*vertical\s*;`,
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
		`.ui-text-field-error-icon { position: absolute; top: 50%; inset-inline-end: 1rem; width: var(--ui-text-field-error-icon-size); height: var(--ui-text-field-error-icon-size);`,
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
		`.ui-text-field { --ui-text-field-textarea-min-height: 7rem; --ui-text-field-error-icon-size: var(--ui-size-icon); display: grid; width: min(100%, 24rem); gap: 0;`,
		`.ui-text-field-message { margin: 0; padding: var(--ui-space-1) var(--ui-space-4) 0; color: var(--ui-color-fg-muted); font: var(--ui-type-body-sm);`,
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
	reducedMotionCSS := entryMediaBlock(t, css, "@media (prefers-reduced-motion: reduce)")
	for _, selector := range []string{`.ui-text-field-control`, `.ui-text-field-control > label`} {
		pattern := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(selector) + `\s*\{[^}]*transition:\s*none\s*;?[^}]*\}`)
		if !pattern.MatchString(reducedMotionCSS) {
			t.Errorf("reduced-motion CSS must disable transitions for %s", selector)
		}
	}
}
