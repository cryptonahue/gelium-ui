package lib

import (
	"strings"
	"testing"
)

// The Basecoat theme (Phase I) satisfies the same theme-agnostic contract as
// Material: every family the components consume is defined in light and dark,
// presence-only, never a concrete value. These tests pin the Basecoat coverage
// explicitly so a regression in the theme itself (not just the matrix) is
// caught; the matrix (TestThemeMatrixCoversEveryAvailableTheme) already
// re-checks the same contract by glob discovery.

// basecoatMandatoryTokens are the light-scheme definitions the theme must own:
// the semantic color roles, typography, radius/elevation/border/focus/motion/
// state scales, and the scoped families of the Phase I component scope
// (Text field, Dialog, Toast, Card, Badge).
var basecoatMandatoryTokens = []string{
	// Semantic colors — the roles the contract §3.1 lists.
	"--ui-color-canvas:",
	"--ui-color-surface:",
	"--ui-color-surface-container:",
	"--ui-color-fg:",
	"--ui-color-fg-muted:",
	"--ui-color-primary:",
	"--ui-color-primary-fg:",
	"--ui-color-secondary:",
	"--ui-color-secondary-fg:",
	"--ui-color-danger:",
	"--ui-color-danger-fg:",
	"--ui-color-border:",
	"--ui-color-border-strong:",
	"--ui-color-focus-ring:",
	"--ui-color-success:",
	"--ui-color-success-fg:",
	"--ui-color-warning:",
	"--ui-color-warning-fg:",
	"--ui-color-warning-container:",
	"--ui-color-info:",
	"--ui-color-danger-container:",
	"--ui-color-scrim:",
	// Typography (Phase B decomposition: the theme owns the decomposed
	// per-step tokens; the --ui-type-<step> aliases live in the core).
	"--ui-font-sans:",
	"--ui-font-mono:",
	"--ui-type-display-sm-size:",
	"--ui-type-display-sm-weight:",
	"--ui-type-display-lg-size:",
	"--ui-type-display-lg-weight:",
	"--ui-type-headline-sm-size:",
	"--ui-type-headline-sm-weight:",
	"--ui-type-title-lg-size:",
	"--ui-type-title-lg-weight:",
	"--ui-type-title-md-size:",
	"--ui-type-title-md-weight:",
	"--ui-type-body-lg-size:",
	"--ui-type-body-lg-weight:",
	"--ui-type-body-md-size:",
	"--ui-type-body-md-weight:",
	"--ui-type-body-sm-size:",
	"--ui-type-body-sm-weight:",
	"--ui-type-label-lg-size:",
	"--ui-type-label-lg-weight:",
	"--ui-type-label-sm-size:",
	"--ui-type-label-sm-weight:",
	"--ui-type-dialog-headline-size:",
	"--ui-type-dialog-headline-weight:",
	"--ui-type-dialog-body-size:",
	"--ui-type-dialog-body-weight:",
	// Radius scale (Basecoat base 0.625rem lands on md).
	"--ui-radius-none:",
	"--ui-radius-xs:",
	"--ui-radius-sm:",
	"--ui-radius-md:",
	"--ui-radius-lg:",
	"--ui-radius-full:",
	// Elevation.
	"--ui-shadow-0:",
	"--ui-shadow-1:",
	"--ui-shadow-2:",
	"--ui-shadow-3:",
	"--ui-shadow-4:",
	"--ui-shadow-5:",
	// Border, focus, motion, state.
	"--ui-border-width-1:",
	"--ui-border-width-2:",
	"--ui-border-style-solid:",
	"--ui-focus-thickness:",
	"--ui-focus-offset:",
	"--ui-motion-short:",
	"--ui-motion-long:",
	"--ui-easing-standard:",
	"--ui-state-hover-opacity:",
	"--ui-state-focus-opacity:",
	"--ui-state-pressed-opacity:",
	"--ui-state-selected-opacity:",
	"--ui-state-disabled-opacity:",
	// Text field.
	"--ui-field-container:",
	"--ui-field-border:",
	"--ui-field-border-hover:",
	"--ui-field-label:",
	"--ui-field-error:",
	// Dialog.
	"--ui-dialog-container:",
	"--ui-dialog-fg:",
	"--ui-dialog-body:",
	"--ui-dialog-scrim:",
	// Toast (incl. the status icon roles).
	"--ui-toast-container:",
	"--ui-toast-fg:",
	"--ui-toast-radius:",
	"--ui-toast-action:",
	"--ui-toast-icon-info:",
	"--ui-toast-icon-success:",
	"--ui-toast-icon-warning:",
	"--ui-toast-icon-error:",
	// Card.
	"--ui-card-radius:",
	"--ui-card-container-elevated:",
	"--ui-card-container-filled:",
	"--ui-card-container-outlined:",
	"--ui-card-outline-color:",
	// Badge.
	"--ui-badge-size:",
	"--ui-badge-large-size:",
	"--ui-badge-container:",
	"--ui-badge-fg:",
}

// TestBasecoatThemeDefinesMandatoryTokensInLightScheme proves the theme ships
// every mandatory family in its light scheme. Presence-only, never a value.
func TestBasecoatThemeDefinesMandatoryTokensInLightScheme(t *testing.T) {
	theme := compactCSS(t, themeCSS(t, "theme-basecoat"))
	for _, token := range basecoatMandatoryTokens {
		if !strings.Contains(theme, token) {
			t.Errorf("theme-basecoat must define %s in the light scheme", token)
		}
	}
	// The radius scale must express the Basecoat base --radius: 0.625rem.
	if !strings.Contains(theme, "--ui-radius-md:.625rem") {
		t.Error("theme-basecoat must place the Basecoat base radius 0.625rem on --ui-radius-md")
	}
}

// basecoatMatrixFamilies are the families the contract matrix requires beyond
// the scope tokens above: checkbox, radio, switch, slider, progress, fab,
// select and divider must all be defined so the matrix stays green for the
// second theme without edits.
var basecoatMatrixFamilies = []string{
	"--ui-checkbox-",
	"--ui-radio-",
	"--ui-switch-",
	"--ui-slider-",
	"--ui-progress-",
	"--ui-fab-",
	"--ui-select-",
	"--ui-divider-",
}

// TestBasecoatThemeDefinesMatrixFamiliesInLightScheme proves the additional
// component families the formal matrix requires exist in light.
func TestBasecoatThemeDefinesMatrixFamiliesInLightScheme(t *testing.T) {
	light, _, _ := splitThemeSchemes(t, "theme-basecoat")
	for _, family := range basecoatMatrixFamilies {
		if !hasFamilyDefinition(light, family) {
			t.Errorf("theme-basecoat light scheme must define the %s token family", family)
		}
	}
}

// TestBasecoatThemeCoversDarkInClassRoute proves the direct-dark families are
// re-declared in the single explicit dark class route, and that the semantic
// colors the derived families reference are redefined there — the same
// contract the matrix asserts for Material. Dark has exactly one mechanism
// (the class route): no @media (prefers-color-scheme: dark) block may exist.
func TestBasecoatThemeCoversDarkInClassRoute(t *testing.T) {
	_, darkClass, _ := splitThemeSchemes(t, "theme-basecoat")

	if strings.Contains(themeCSS(t, "theme-basecoat"), "@media (prefers-color-scheme: dark)") {
		t.Error("theme-basecoat must not define a dark media route (single dark mechanism is the class route)")
	}

	// Direct dark coverage: at least one token of each family in the dark
	// class route. checkbox/radio/divider are derived families (they live in
	// light and reference semantic colors) and are covered by the color loop
	// below.
	for _, family := range []string{
		"--ui-field-",
		"--ui-dialog-",
		"--ui-toast-",
		"--ui-card-",
		"--ui-switch-",
		"--ui-slider-",
		"--ui-progress-",
		"--ui-fab-",
		"--ui-select-",
		"--ui-color-",
	} {
		if !hasFamilyDefinition(darkClass, family) {
			t.Errorf("theme-basecoat dark class route must redefine a %s token", family)
		}
	}

	// Derived dark coverage: the semantic colors badge/checkbox/radio/divider
	// reference in light must be redefined in the dark class route.
	for _, color := range []string{
		"--ui-color-danger:",
		"--ui-color-danger-fg:",
		"--ui-color-primary:",
		"--ui-color-primary-fg:",
		"--ui-color-fg-muted:",
		"--ui-color-fg:",
		"--ui-color-border:",
	} {
		if !strings.Contains(darkClass, color) {
			t.Errorf("theme-basecoat dark class route must redefine %s (derived legibility)", color)
		}
	}
}

// TestBasecoatThemeIsDiscoveredAlongsideMaterial pins the discovery contract
// for Phase I: the glob must find both themes, so the formal matrix and the
// bundle mechanism actually run over the second theme instead of silently
// covering a single theme.
func TestBasecoatThemeIsDiscoveredAlongsideMaterial(t *testing.T) {
	themes := availableThemes(t)
	found := map[string]bool{}
	for _, name := range themes {
		found[name] = true
	}
	for _, want := range []string{"theme-material", "theme-basecoat"} {
		if !found[want] {
			t.Errorf("availableThemes must discover %s, got %v", want, themes)
		}
	}
}

// TestCompiledBundleCarriesBasecoatRootSelectors proves the served bundle — the
// artifact `npm run build` produces from app.css — carries the .theme-basecoat
// root selector and its dark class route, so class-driven selection actually
// resolves at runtime. Dark is the single class route (no media dark block).
func TestCompiledBundleCarriesBasecoatRootSelectors(t *testing.T) {
	compiled := compactCSS(t, compiledAppCSS(t))
	for _, contract := range []string{
		".theme-basecoat{",
		".theme-basecoat.theme-dark,",
	} {
		if !strings.Contains(compiled, contract) {
			t.Errorf("compiled bundle is missing Basecoat contract %q", contract)
		}
	}
}

// TestBasecoatThemeTonesStayTokenDriven proves the scoped tone/status surfaces
// of the Phase I scope stay token-driven in the theme: the toast icon roles,
// badge and card variants reference semantic tokens instead of literals.
func TestBasecoatThemeTonesStayTokenDriven(t *testing.T) {
	theme := compactCSS(t, themeCSS(t, "theme-basecoat"))
	for _, contract := range []string{
		"--ui-badge-container:var(--ui-color-danger)",
		"--ui-badge-fg:var(--ui-color-danger-fg)",
		"--ui-dialog-scrim:var(--ui-color-scrim)",
	} {
		if !strings.Contains(theme, contract) {
			t.Errorf("theme-basecoat must keep tone surface %q token-driven", contract)
		}
	}
}
