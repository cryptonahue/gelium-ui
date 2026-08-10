package web

import (
	"regexp"
	"strings"
	"testing"
)

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

func TestReducedMotionDisablesButtonSpinnerAnimation(t *testing.T) {
	css := sourceAppCSS(t)
	reducedMotionCSS := entryMediaBlock(t, css, "@media (prefers-reduced-motion: reduce)")
	spinnerAnimationNone := regexp.MustCompile(`(?s)\.ui-button-spinner\s*\{[^}]*animation:\s*none\s*;?[^}]*\}`)
	if !spinnerAnimationNone.MatchString(reducedMotionCSS) {
		t.Error("reduced-motion CSS must disable the spinner with animation: none")
	}
}
