package lib

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The formal theme matrix (Phase C) verifies the theme-agnostic contract:
//
//	theme × component × scheme        — every discovered theme defines each
//	                                    component's token family in light and
//	                                    in the single explicit dark class route.
//	component × state                 — every component CSS covers its
//	                                    documented states with tokens, never
//	                                    literals.
//
// The suite discovers themes by glob (availableThemes), so adding a new theme
// under lib/themes/<name>.css (e.g. theme-basecoat, Phase I) extends the
// matrix without editing tests — the matrix never assumes a theme that does
// not exist, and it never asserts a concrete value, only token presence.
//
// Dark coverage has one documented mode, mirroring the theme contract:
//
//   - direct:   the family is re-declared in the dark class route (at least
//               one token of the family). This is the Material behaviour for
//               field, dialog, toast, card, switch, slider, progress, fab,
//               select and the semantic color family.
//   - derived:  the family lives in light and references --ui-color-* semantic
//               tokens; dark legibility comes from those semantic tokens being
//               redefined in the dark class route. This is the Material
//               behaviour for badge, checkbox, radio and divider.
//
// A new theme satisfies the matrix by (1) defining every family below in light
// and (2) covering dark either directly or through redefined semantic colors —
// the checklist in docs/gelium-ui-theme-verification.md. The dark media route
// (@media prefers-color-scheme) is no longer part of the contract: dark is
// served by the single explicit class route (Phase A).

// stateExpectation documents one state a component CSS must cover with tokens.
type stateExpectation struct {
	state    string // documented state: hover, focus, pressed, disabled, selected, error, loading, empty (or variant/tone)
	selector string // selector the component CSS must contain
	token    string // var(--ui-*) token the state must be driven by (never a literal)
}

// themeMatrixComponent is one row of the formal matrix.
type themeMatrixComponent struct {
	component     string             // documented component name
	cssFile       string             // source file in web/styles
	family        string             // token family prefix, e.g. "--ui-fab-"
	derivedColors []string           // empty = direct dark coverage; non-empty = family derives from these --ui-color-* tokens
	states        []stateExpectation // the states the component CSS must cover with tokens
}

var themeMatrixComponents = []themeMatrixComponent{
	{
		component: "button", cssFile: "button.css", family: "--ui-color-",
		states: []stateExpectation{
			{state: "hover", selector: `.ui-button:hover:not(:disabled):not([aria-disabled="true"])`, token: "--ui-state-hover-opacity"},
			{state: "focus", selector: `.ui-button:focus-visible`, token: "--ui-color-focus-ring"},
			{state: "pressed", selector: `.ui-button:active:not(:disabled):not([aria-disabled="true"])`, token: "--ui-state-pressed-opacity"},
			{state: "disabled", selector: `.ui-button:disabled`, token: "--ui-state-disabled-opacity"},
			{state: "loading", selector: `.ui-button-spinner`, token: "--ui-button-spinner-size"},
		},
	},
	{
		component: "text-field", cssFile: "text-field.css", family: "--ui-field-",
		states: []stateExpectation{
			{state: "hover", selector: `.ui-text-field-control:hover:not(:focus-within):has(input:not(:disabled))`, token: "--ui-field-border-hover"},
			{state: "focus", selector: `.ui-text-field-control:focus-within`, token: "--ui-color-primary"},
			{state: "error", selector: `.ui-text-field-error .ui-text-field-message`, token: "--ui-field-error"},
			{state: "disabled", selector: `.ui-text-field-disabled .ui-text-field-control > label`, token: "--ui-state-disabled-opacity"},
			{state: "empty", selector: `.ui-text-field-control:has(input:not(:placeholder-shown)) > label`, token: "--ui-type-label-sm"},
		},
	},
	{
		component: "dialog", cssFile: "dialog.css", family: "--ui-dialog-",
		states: []stateExpectation{
			{state: "open", selector: `.ui-dialog[open]`, token: "--ui-motion-long"},
			{state: "backdrop", selector: `.ui-dialog::backdrop`, token: "--ui-dialog-scrim"},
		},
	},
	{
		component: "toast", cssFile: "toast.css", family: "--ui-toast-",
		states: []stateExpectation{
			{state: "loading", selector: `.ui-toast-show`, token: "--ui-motion-short"},
			{state: "focus", selector: `.ui-toast-action:focus-visible`, token: "--ui-color-focus-ring"},
			{state: "error", selector: `.ui-toast-icon-error`, token: "--ui-toast-icon-error"},
		},
	},
	{
		component: "card", cssFile: "card.css", family: "--ui-card-",
		states: []stateExpectation{
			{state: "focus", selector: `.ui-card:focus-visible`, token: "--ui-color-focus-ring"},
			{state: "variant-elevated", selector: `.ui-card-elevated`, token: "--ui-card-container-elevated"},
			{state: "variant-filled", selector: `.ui-card-filled`, token: "--ui-card-container-filled"},
			{state: "variant-outlined", selector: `.ui-card-outlined`, token: "--ui-card-outline-color"},
		},
	},
	{
		component: "badge", cssFile: "badge.css", family: "--ui-badge-",
		derivedColors: []string{"--ui-color-danger", "--ui-color-danger-fg"},
		states: []stateExpectation{
			{state: "tone-error", selector: `.ui-badge--error`, token: "--ui-color-danger"},
			{state: "tone-success", selector: `.ui-badge--success`, token: "--ui-color-success"},
			{state: "tone-warning", selector: `.ui-badge--warning`, token: "--ui-color-warning-container"},
			{state: "tone-info", selector: `.ui-badge--info`, token: "--ui-color-info"},
		},
	},
	{
		component: "checkbox", cssFile: "checkbox.css", family: "--ui-checkbox-",
		derivedColors: []string{"--ui-color-primary", "--ui-color-primary-fg", "--ui-color-fg-muted", "--ui-color-fg", "--ui-color-danger"},
		states: []stateExpectation{
			{state: "hover", selector: `.ui-checkbox:hover`, token: "--ui-checkbox-hover-outline"},
			{state: "focus", selector: `.ui-checkbox input:focus-visible`, token: "--ui-color-focus-ring"},
			{state: "selected", selector: `.ui-checkbox input:checked`, token: "--ui-checkbox-container"},
			{state: "error", selector: `.ui-checkbox input[type="checkbox"][aria-invalid="true"]`, token: "--ui-checkbox-error"},
			{state: "disabled", selector: `.ui-checkbox input:disabled`, token: "--ui-state-disabled-opacity"},
		},
	},
	{
		component: "radio", cssFile: "radio.css", family: "--ui-radio-",
		derivedColors: []string{"--ui-color-primary", "--ui-color-fg-muted", "--ui-color-fg", "--ui-color-danger"},
		states: []stateExpectation{
			{state: "hover", selector: `.ui-radio:hover`, token: "--ui-radio-hover-outline"},
			{state: "focus", selector: `.ui-radio input:focus-visible`, token: "--ui-color-focus-ring"},
			{state: "selected", selector: `.ui-radio input:checked`, token: "--ui-radio-checked"},
			{state: "error", selector: `.ui-radio input[type="radio"][aria-invalid="true"]`, token: "--ui-color-danger"},
			{state: "disabled", selector: `.ui-radio input:disabled`, token: "--ui-state-disabled-opacity"},
		},
	},
	{
		component: "switch", cssFile: "switch.css", family: "--ui-switch-",
		states: []stateExpectation{
			{state: "hover", selector: `.ui-switch:hover`, token: "--ui-color-fg"},
			{state: "focus", selector: `.ui-switch input:focus-visible`, token: "--ui-color-focus-ring"},
			{state: "selected", selector: `.ui-switch input:checked + .ui-switch-track`, token: "--ui-switch-track"},
			{state: "pressed", selector: `.ui-switch input:active:not(:disabled) ~ .ui-switch-handle`, token: "--ui-switch-handle-pressed-size"},
			{state: "disabled", selector: `.ui-switch input:disabled`, token: "--ui-switch-disabled-track-opacity"},
		},
	},
	{
		component: "slider", cssFile: "slider.css", family: "--ui-slider-",
		states: []stateExpectation{
			{state: "focus", selector: `.ui-slider input[type="range"]:focus-visible`, token: "--ui-color-focus-ring"},
			{state: "pressed", selector: `.ui-slider input[type="range"]:active:not(:disabled)`, token: "--ui-slider-handle-pressed-size"},
			{state: "disabled", selector: `.ui-slider input[type="range"]:disabled`, token: "--ui-slider-disabled"},
		},
	},
	{
		component: "progress", cssFile: "progress.css", family: "--ui-progress-",
		states: []stateExpectation{
			{state: "determinate", selector: `.ui-progress progress`, token: "--ui-progress-indicator"},
		},
	},
	{
		component: "fab", cssFile: "fab.css", family: "--ui-fab-",
		states: []stateExpectation{
			{state: "hover", selector: `.ui-fab:hover:not(:disabled):not([aria-disabled="true"])`, token: "--ui-state-hover-opacity"},
			{state: "focus", selector: `.ui-fab:focus-visible:not(:disabled):not([aria-disabled="true"])`, token: "--ui-color-focus-ring"},
			{state: "pressed", selector: `.ui-fab:active:not(:disabled):not([aria-disabled="true"])`, token: "--ui-state-pressed-opacity"},
			{state: "disabled", selector: `.ui-fab:disabled`, token: "--ui-state-disabled-opacity"},
		},
	},
	{
		component: "select", cssFile: "select.css", family: "--ui-select-",
		states: []stateExpectation{
			{state: "hover", selector: `.ui-select select:hover:not(:disabled):not([aria-invalid="true"])`, token: "--ui-select-hover-outline"},
			{state: "focus", selector: `.ui-select select:focus-visible`, token: "--ui-color-focus-ring"},
			{state: "error", selector: `.ui-select select[aria-invalid="true"]`, token: "--ui-select-error"},
			{state: "disabled", selector: `.ui-select select:disabled`, token: "--ui-select-disabled-opacity"},
			{state: "empty", selector: `.ui-select:has(select option:checked:not([value=""])) .ui-select-label`, token: "--ui-type-body-sm"},
		},
	},
	{
		component: "divider", cssFile: "divider.css", family: "--ui-divider-",
		derivedColors: []string{"--ui-color-border"},
		// Divider is decorative: it owns no interactive states. Its only
		// contract is the family + the derived semantic color above.
		states: []stateExpectation{
			{state: "base", selector: `.ui-divider`, token: "--ui-divider-color"},
		},
	},
}

// splitThemeSchemes splits one theme's CSS into its light block, its explicit
// dark-class block, and (while the legacy media route still exists) its dark
// media block. The contract is the single dark mechanism: only the class route
// may carry dark values, so darkMedia is expected to be empty; the helper
// tolerates a leftover media block so themes can be inspected mid-migration,
// and the no-media assertion lives in the theme-specific contract tests.
// Presence-only: the blocks are never parsed for values.
func splitThemeSchemes(t *testing.T, theme string) (light, darkClass, darkMedia string) {
	t.Helper()
	css := themeCSS(t, theme)

	darkClassStart := strings.Index(css, ".theme-dark")
	if darkClassStart < 0 {
		t.Errorf("theme %q must declare an explicit dark class route (.theme-dark / .dark / [data-theme=\"dark\"])", theme)
		return "", "", ""
	}

	light = css[:darkClassStart]
	darkClass = css[darkClassStart:]
	if mediaStart := strings.Index(css, "@media (prefers-color-scheme: dark)"); mediaStart >= 0 {
		darkClass = css[darkClassStart:mediaStart]
		darkMedia = css[mediaStart:]
	}
	return light, darkClass, darkMedia
}

// hasFamilyDefinition reports whether a CSS block declares at least one token
// of the given family (--ui-<family>-<token>:), as a definition, not a reference.
func hasFamilyDefinition(block, family string) bool {
	re := regexp.MustCompile(regexp.QuoteMeta(family) + `[a-z0-9-]*\s*:`)
	return re.MatchString(block)
}

// TestThemeMatrixCoversEveryAvailableTheme proves the matrix is theme-agnostic
// by construction: it iterates the themes the glob discovers (never a hardcoded
// path), runs the full component matrix for each, and documents the contract a
// new theme must satisfy. Adding lib/themes/theme-basecoat.css later extends
// the matrix with zero test edits.
func TestThemeMatrixCoversEveryAvailableTheme(t *testing.T) {
	themes := availableThemes(t)
	if len(themes) == 0 {
		t.Fatal("the glob lib/themes/*.css must discover at least one theme")
	}
	for _, theme := range themes {
		theme := theme
		t.Run(theme, func(t *testing.T) {
			t.Parallel()
			light, darkClass, _ := splitThemeSchemes(t, theme)

			for _, comp := range themeMatrixComponents {
				comp := comp
				t.Run(comp.component, func(t *testing.T) {
					t.Parallel()
					if !hasFamilyDefinition(light, comp.family) {
						t.Errorf("%s: light scheme must define the %s token family", comp.component, comp.family)
					}
					if len(comp.derivedColors) == 0 {
						// direct dark coverage: the family is re-declared in
						// the single dark class route.
						if !hasFamilyDefinition(darkClass, comp.family) {
							t.Errorf("%s: dark class route must redefine a %s token", comp.component, comp.family)
						}
						return
					}
					// derived dark coverage: the family lives in light and
					// references semantic colors that the dark class route
					// redefines.
					for _, color := range comp.derivedColors {
						if !strings.Contains(light, "var("+color+")") {
							t.Errorf("%s: light scheme must derive %s from %s", comp.component, comp.family, color)
						}
						if !strings.Contains(darkClass, color+":") {
							t.Errorf("%s: dark class route must redefine %s (derived legibility)", comp.component, color)
						}
					}
				})
			}
		})
	}
}

// TestThemeMatrixStateCoverageIsTokenDriven proves every documented state of
// every component is covered by a var(--ui-*) token in its own source CSS —
// never by a literal. This is the component × state half of the matrix. It
// reads component source (theme-agnostic), so no theme dependency.
func TestThemeMatrixStateCoverageIsTokenDriven(t *testing.T) {
	for _, comp := range themeMatrixComponents {
		comp := comp
		t.Run(comp.component, func(t *testing.T) {
			css := sourceComponentCSS(t, comp.cssFile)
			for _, st := range comp.states {
				st := st
				t.Run(st.state, func(t *testing.T) {
					if !strings.Contains(css, st.selector) {
						t.Errorf("%s.css must contain selector %q for state %q", comp.cssFile, st.selector, st.state)
					}
					if !strings.Contains(css, "var("+st.token+")") {
						t.Errorf("%s.css must drive state %q with token %s (no literals)", comp.cssFile, st.state, st.token)
					}
				})
			}
		})
	}
}

// TestSourceAppCSSKeepsCoreBeforeThemeCascade proves the sourceAppCSS helper
// mirrors the app.css import order: core tokens (tokens.css) must come before
// the theme so raw-source assertions see the same cascade as the build — the
// theme overrides the core defaults, never the other way around. Presence-only.
func TestSourceAppCSSKeepsCoreBeforeThemeCascade(t *testing.T) {
	css := sourceAppCSS(t)
	coreMarker := "--ui-color-canvas:"
	themeMarker := ".theme-material"
	coreIdx, themeIdx := strings.Index(css, coreMarker), strings.Index(css, themeMarker)
	if coreIdx < 0 {
		t.Fatalf("sourceAppCSS is missing the core marker %q", coreMarker)
	}
	if themeIdx < 0 {
		t.Fatalf("sourceAppCSS is missing the theme marker %q", themeMarker)
	}
	if themeIdx < coreIdx {
		t.Errorf("sourceAppCSS must concatenate core tokens before the theme (app.css imports tokens.css then the theme)")
	}
}

// TestThemeMatrixDefaultThemeIsPartOfDiscovery pins the contract matrix to the
// concrete theme the repository ships today: theme-material must be present on
// disk and the matrix must cover it, so regressions in the only theme cannot
// hide behind a matrix that silently ran over an empty set.
func TestThemeMatrixDefaultThemeIsPartOfDiscovery(t *testing.T) {
	for _, theme := range availableThemes(t) {
		if theme == defaultThemeName {
			return
		}
	}
	t.Fatalf("the contract matrix must cover %q, the repository's only theme", defaultThemeName)
}

// TestThemeMatrixLabelMdDefinedPerTheme is the Phase B R2 extension of the
// formal matrix: every theme must define the label-md type family in its light
// scheme (the decomposed --ui-type-label-md-* tokens). label-md is the R2
// closure step — a real step, never an alias of label-lg — so the matrix
// covers it like any other family the components consume.
func TestThemeMatrixLabelMdDefinedPerTheme(t *testing.T) {
	for _, theme := range availableThemes(t) {
		theme := theme
		t.Run(theme, func(t *testing.T) {
			light, _, _ := splitThemeSchemes(t, theme)
			for _, prop := range typeStepProps {
				token := "--ui-type-label-md-" + prop + ":"
				if !strings.Contains(light, token) {
					t.Errorf("%s light scheme must define the label-md decomposed token %s (R2 closure)", theme, token)
				}
			}
			// Closure: the theme must never define label-md as a label-lg
			// alias (that would defeat the R2 split).
			if strings.Contains(light, "--ui-type-label-md: var(--ui-type-label-lg)") {
				t.Errorf("%s must define label-md standalone, never as var(--ui-type-label-lg)", theme)
			}
		})
	}
}

// TestThemeFilesLiveUnderLibThemes pins the flat package layout: every
// discovered theme is a single file at lib/themes/<name>.css — not the
// retired nested themes/<name>/theme.css tree.
func TestThemeFilesLiveUnderLibThemes(t *testing.T) {
	themes := availableThemes(t)
	if len(themes) == 0 {
		t.Fatal("lib/themes/*.css must discover at least one theme")
	}
	root := repositoryRoot(t)
	for _, name := range themes {
		flat := filepath.Join(root, "lib", "themes", name+".css")
		if _, err := os.Stat(flat); err != nil {
			t.Errorf("theme %s must live at lib/themes/%s.css: %v", name, name, err)
		}
		nested := filepath.Join(root, "themes", name, "theme.css")
		if _, err := os.Stat(nested); err == nil {
			t.Errorf("retired nested path themes/%s/theme.css must not exist (use lib/themes/%s.css)", name, name)
		}
		nestedLib := filepath.Join(root, "lib", "themes", name, "theme.css")
		if _, err := os.Stat(nestedLib); err == nil {
			t.Errorf("retired nested path lib/themes/%s/theme.css must not exist (use flat lib/themes/%s.css)", name, name)
		}
	}
}
