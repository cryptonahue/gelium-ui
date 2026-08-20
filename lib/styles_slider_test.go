package lib

import (
	"regexp"
	"strings"
	"testing"
)

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
		`.ui-slider--ticks {`,
		`var(--ui-slider-tick-interval,`,
		`repeating-linear-gradient(to right, var(--ui-slider-tick)`,
		`.ui-slider--value-label {`,
		`attr(data-value)`,
		`left: var(--ui-slider-fill, 0%);`,
		`.ui-slider--value-label:has(input[type="range"]:focus-visible)::after`,
		`.ui-slider--value-label:has(input[type="range"]:active:not(:disabled))::after`,
		`.ui-slider--ticks:has(input:disabled)::before { opacity: 0; }`,
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
		`repeating-linear-gradient(to right, Canvas`,
		`.ui-slider--value-label::after { background: CanvasText;`,
	} {
		if !strings.Contains(forced, contract) {
			t.Errorf("slider must stay distinguishable in forced colors; missing %q", contract)
		}
	}
}

func TestSliderThemeDefinesPublicUIFamily(t *testing.T) {
	theme := themeCSS(t, "theme-material")
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
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-slider`,
		`-webkit-slider-runnable-track`,
		`-moz-range-progress`,
		`-moz-range-thumb`,
		`.ui-slider--ticks`,
		`.ui-slider--value-label`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled slider CSS is missing %q", contract)
		}
	}
}
