package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestSwitchPrimitiveCSSMapsNativeControlAndStates(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")

	for _, contract := range []string{
		`.ui-switch {`,
		`gap: var(--ui-space-2);`,
		`input[type="checkbox"] {`,
		`appearance: none;`,
		`width: var(--ui-switch-width, 52px);`,
		`height: var(--ui-switch-height, 32px);`,
		`border-radius: var(--ui-switch-radius, var(--ui-radius-full));`,
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
	theme := themeCSS(t, "theme-material")
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
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-switch`,
		`var(--ui-switch-width`,
		`var(--ui-switch-track-outline)`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled switch CSS is missing %q", contract)
		}
	}
}
