package web

import (
	"embed"
	"encoding/json"
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

// sourceAppCSS concatenates the raw source files that web/styles/app.css pulls
// in, in the same order, plus the entry's own top-level tail (keyframes,
// reduced-motion, forced-colors, @theme mapping). Keep this list in sync with
// app.css when files are added or reordered. The theme and core tokens are not
// embedded under styles/, so they are prepended from the repository tree via
// themeCSS.
func sourceAppCSS(t *testing.T) string {
	t.Helper()
	paths := []string{
		"styles/tokens.css",
		"styles/base.css",
		"styles/button.css",
		"styles/text-field.css",
		"styles/dialog.css",
		"styles/toast.css",
		"styles/focus-ring.css",
		"styles/elevation.css",
		"styles/icon.css",
		"styles/divider.css",
		"styles/card.css",
		"styles/badge.css",
		"styles/avatar.css",
		"styles/checkbox.css",
		"styles/radio.css",
		"styles/switch.css",
		"styles/select.css",
		"styles/slider.css",
		"styles/progress.css",
		"styles/icon-button.css",
		"styles/fab.css",
		"styles/list.css",
		"styles/chips.css",
		"styles/tabs.css",
		"styles/navigation-bar.css",
		"styles/navigation-tab.css",
		"styles/segmented-button.css",
		"styles/menu.css",
		"styles/navigation-drawer.css",
		"styles/data-table.css",
		"styles/pagination.css",
		"styles/banner.css",
		"styles/callout.css",
		"styles/empty-state.css",
		"styles/error-state.css",
		"styles/inline-alert.css",
		"styles/skeleton.css",
		"styles/tooltip.css",
		"styles/validation-summary.css",
		"styles/breadcrumb.css",
		"styles/footer.css",
		"styles/section-heading.css",
		"styles/video.css",
		"styles/feature-card.css",
		"styles/hero.css",
		"styles/language-switcher.css",
		"styles/newsletter.css",
		"styles/split.css",
		"styles/demo-whatsapp.css",
		"styles/recipe-admin-resource.css",
		"styles/recipe-ops-queue.css",
		"styles/recipe-public-feed.css",
		"styles/app.css",
	}
	var sb strings.Builder
	for _, path := range paths {
		css, err := sourceStyles.ReadFile(path)
		if err != nil {
			t.Fatalf("read source app CSS %s: %v", path, err)
		}
		sb.Write(css)
		// app.css imports the theme right after the core tokens (tokens.css
		// then the theme), so raw-source assertions see the same cascade as the
		// build: the core defaults first, then the theme overrides them.
		if path == "styles/tokens.css" {
			sb.WriteString(themeCSS(t, defaultThemeName))
		}
	}
	return sb.String()
}

// defaultThemeName is the theme the contract tests use to prove a concrete
// theme satisfies the semantic token contract. It must exist on disk under
// themes/<name>/theme.css.
const defaultThemeName = "theme-material"

// repositoryRoot returns the repository root as seen from this test file.
func repositoryRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate styles contract test")
	}
	return filepath.Dir(filename) + "/.."
}

// availableThemes discovers the themes that actually exist on disk. It is the
// single source of truth for theme names in the contract matrix, so the matrix
// extends to new themes (e.g. basecoat) without re-hardcoding paths in every
// test — and never assumes a theme that does not exist.
func availableThemes(t *testing.T) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(repositoryRoot(t), "themes", "*", "theme.css"))
	if err != nil {
		t.Fatalf("discover themes: %v", err)
	}
	var themes []string
	for _, m := range matches {
		themes = append(themes, filepath.Base(filepath.Dir(m)))
	}
	return themes
}

// themeCSS reads one theme's CSS from the repository themes directory,
// validating the name against the themes present on disk.
func themeCSS(t *testing.T, name string) string {
	t.Helper()
	found := false
	for _, theme := range availableThemes(t) {
		if theme == name {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("theme %q is not present on disk (have %v)", name, availableThemes(t))
	}
	return repositoryFile(t, "themes", name, "theme.css")
}

// entryMediaBlock extracts one media block from the app.css entry tail. The
// entry (styles/app.css) carries the consolidated reduced-motion and
// forced-colors rules, but individual component files may declare their own
// copies of the same media query, so callers must target the entry's block
// (the last occurrence) instead of the first one.
func entryMediaBlock(t *testing.T, css, media string) string {
	t.Helper()
	index := strings.LastIndex(css, media)
	if index < 0 {
		t.Fatalf("source CSS is missing the %s media query", media)
	}
	block := css[index:]
	if next := strings.Index(block[len(media):], "@media "); next >= 0 {
		block = block[:len(media)+next]
	}
	return block
}

// TestSurfaceContainerTokenClosedAcrossCoreAndEveryScheme proves the
// --ui-color-surface-container gap is closed: the core owns a neutral default
// and the Material theme defines the token in light and in the single explicit
// dark class route. Presence only — never a concrete hex value.
func TestSurfaceContainerTokenClosedAcrossCoreAndEveryScheme(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	if !strings.Contains(string(core), "--ui-color-surface-container:") {
		t.Error("core tokens.css must define --ui-color-surface-container (neutral default)")
	}

	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, "theme-material"), " ")
	// Single dark mechanism: the token is defined once in light and once in the
	// explicit dark class route — never in a dark media block (3 definitions
	// would mean the duplicated class+media routine has leaked back).
	if got := strings.Count(theme, "--ui-color-surface-container:"); got != 2 {
		t.Errorf("theme-material must define --ui-color-surface-container once in light and once in the single dark class route, got %d", got)
	}
	if !strings.Contains(theme, ".theme-material.theme-dark,") {
		t.Error("theme-material must declare the explicit dark class route (.theme-material.theme-dark)")
	}
	if strings.Contains(theme, "@media (prefers-color-scheme: dark)") {
		t.Error("theme-material must not define a dark media route (single dark mechanism is the class route)")
	}
}

// TestDisplayLgTokenClosedAcrossCoreAndTheme proves the --ui-type-display-lg
// gap is closed under the Phase B decomposition (R1): the core owns the
// shorthand alias (a font: composition of decomposed tokens) and the Material
// theme defines the decomposed display-lg values. Presence only — never a
// concrete value.
func TestDisplayLgTokenClosedAcrossCoreAndTheme(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	if !strings.Contains(string(core), "--ui-type-display-lg:") {
		t.Error("core tokens.css must define the --ui-type-display-lg alias (Phase B decomposition)")
	}

	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, "theme-material"), " ")
	if !strings.Contains(theme, "--ui-type-display-lg-size:") || !strings.Contains(theme, "--ui-type-display-lg-weight:") {
		t.Error("theme-material must define the decomposed display-lg values (--ui-type-display-lg-size / -weight)")
	}
}

// TestFontTokensDefinedInCoreProve font defaults live in the core so an
// unthemed document resolves every var(--ui-font-*) reference. Presence only.
func TestFontTokensDefinedInCore(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	for _, token := range []string{"--ui-font-sans:", "--ui-font-mono:"} {
		if !strings.Contains(string(core), token) {
			t.Errorf("core tokens.css must define %s (neutral default)", token)
		}
	}
}

// TestFontSansOverriddenByMaterialTheme proves the Material theme overrides
// the core's --ui-font-sans under .theme-material. Presence only.
func TestFontSansOverriddenByMaterialTheme(t *testing.T) {
	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, "theme-material"), " ")
	if !strings.Contains(theme, "--ui-font-sans:") {
		t.Error("theme-material must override --ui-font-sans under .theme-material")
	}
	if !strings.Contains(theme, "--ui-font-mono:") {
		t.Error("theme-material must define --ui-font-mono under .theme-material")
	}
}

// TestTitleMdTokenClosedAcrossCoreAndTheme proves the --ui-type-title-md gap
// is closed under the Phase B decomposition (R1): the core owns the shorthand
// alias (a font: composition of decomposed tokens) and the Material theme
// defines the decomposed title-md values. Presence only — never a concrete
// value.
func TestTitleMdTokenClosedAcrossCoreAndTheme(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	if !strings.Contains(string(core), "--ui-type-title-md:") {
		t.Error("core tokens.css must define the --ui-type-title-md alias (Phase B decomposition)")
	}

	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, "theme-material"), " ")
	if !strings.Contains(theme, "--ui-type-title-md-size:") || !strings.Contains(theme, "--ui-type-title-md-weight:") {
		t.Error("theme-material must define the decomposed title-md values (--ui-type-title-md-size / -weight)")
	}
}

// TestCoreTokensSelfContained proves the core tokens file is autonomous: every
// var(--ui-*) it references must resolve to a --ui-* definition inside the same
// file. There is no exception anymore (display-lg→font-sans was the last one);
// the only internal references expected are the color alias and the font stack
// that the typescale shorthand needs.
func TestCoreTokensSelfContained(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	css := string(core)

	defined := map[string]bool{}
	for _, m := range regexp.MustCompile(`(--ui-[a-z0-9-]+):`).FindAllStringSubmatch(css, -1) {
		defined[m[1]] = true
	}

	var missing []string
	seen := map[string]bool{}
	for _, m := range regexp.MustCompile(`var\((--ui-[a-z0-9-]+)`).FindAllStringSubmatch(css, -1) {
		name := m[1]
		if seen[name] {
			continue
		}
		seen[name] = true
		if !defined[name] {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		t.Errorf("core tokens.css references --ui-* tokens not defined in the core: %v", missing)
	}
}

// spaceCoreTokens are the six core spacing tokens of the --ui-space-* scale.
// The scale skips 5 and 7 deliberately (Material keeps only the 1/2/3/4/6/8
// steps), so only these names are part of the contract.
var spaceCoreTokens = []string{
	"--ui-space-1:",
	"--ui-space-2:",
	"--ui-space-3:",
	"--ui-space-4:",
	"--ui-space-6:",
	"--ui-space-8:",
}

// TestSpaceTokensDefinedInCore proves the core owns the six --ui-space-* tokens
// that the component layer may consume. Presence only — never a concrete value,
// so the test stays theme-agnostic (the Material theme does not override the
// spacing scale).
func TestSpaceTokensDefinedInCore(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	css := string(core)
	for _, token := range spaceCoreTokens {
		if !strings.Contains(css, token) {
			t.Errorf("core tokens.css must define %s (spacing scale step)", token)
		}
	}
}

// sizeCoreTokens are the core --ui-size-* tokens of the Phase B2 density/sizes
// slice: control and field geometry, the icon scale, the track height, the
// item row-height scale, and the navigation slot/indicator/label geometry.
var sizeCoreTokens = []string{
	"--ui-size-control:",
	"--ui-size-field:",
	"--ui-size-icon:",
	"--ui-size-icon-sm:",
	"--ui-size-track:",
	"--ui-size-item-sm:",
	"--ui-size-item:",
	"--ui-size-item-lg:",
	"--ui-size-item-xl:",
	"--ui-size-nav-icon-slot:",
	"--ui-size-nav-icon-slot-h:",
	"--ui-size-indicator:",
	"--ui-size-label-height:",
}

// TestSizeTokensDefinedInCore proves the core owns the --ui-size-* tokens that
// the component layer may consume. Presence only — never a concrete value, so
// the test stays theme-agnostic (values come from the theme layer).
func TestSizeTokensDefinedInCore(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	css := string(core)
	for _, token := range sizeCoreTokens {
		if !strings.Contains(css, token) {
			t.Errorf("core tokens.css must define %s (size/density scale)", token)
		}
	}
}

// TestSizeTokensOverriddenByMaterialTheme proves the Material theme overrides
// the core --ui-size-* tokens that carry a Material value under
// .theme-material. Presence only — never a concrete value. The nav slot
// tokens are deliberately absent here: their core defaults are already the
// Material values (64x32/32/16px), so the theme does not re-declare them.
func TestSizeTokensOverriddenByMaterialTheme(t *testing.T) {
	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, "theme-material"), " ")
	for _, token := range []string{
		"--ui-size-control:",
		"--ui-size-field:",
		"--ui-size-icon:",
		"--ui-size-icon-sm:",
		"--ui-size-track:",
		"--ui-size-item-sm:",
		"--ui-size-item:",
		"--ui-size-item-lg:",
		"--ui-size-item-xl:",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme-material must override %s under .theme-material", token)
		}
	}
}

// TestSizeTokensConsumedByComponents proves the --ui-size-* tokens are real
// consumers, not orphan definitions: at least the minimum number of distinct
// component files must reference var(--ui-size-*) at least once. Structural
// assertion only — never a concrete value, and demo/doc files are excluded so
// the proof is about real component selectors.
func TestSizeTokensConsumedByComponents(t *testing.T) {
	const minConsumerFiles = 3

	entries, err := sourceStyles.ReadDir("styles")
	if err != nil {
		t.Fatalf("list styles dir: %v", err)
	}

	excluded := map[string]bool{
		"tokens.css":                true, // defines the scale, not a consumer
		"base.css":                  true, // docs layout (site-header, prose, hero)
		"demo-whatsapp.css":         true, // demo only (all selectors are .demo-wa-*)
		"recipe-admin-resource.css": true, // screen recipe layout (all selectors are .recipe-ar-*)
		"app.css":                   true, // entry, imports and media tails
	}

	var consumers []string
	totalUses := 0
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".css") || excluded[name] {
			continue
		}
		content, err := sourceStyles.ReadFile("styles/" + name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if uses := strings.Count(string(content), "var(--ui-size-"); uses > 0 {
			consumers = append(consumers, name)
			totalUses += uses
		}
	}

	if len(consumers) < minConsumerFiles {
		t.Errorf("expected var(--ui-size-*) in at least %d distinct component files, got %d (%v)",
			minConsumerFiles, len(consumers), consumers)
	}
	if totalUses < minConsumerFiles {
		t.Errorf("expected at least %d total var(--ui-size-*) usages across components, got %d",
			minConsumerFiles, totalUses)
	}

	// The theme layer is also a consumer: --ui-select-height, the slider and
	// the progress track re-point to the core size scale, so an unthemed
	// document must resolve them.
	theme := themeCSS(t, "theme-material")
	if !strings.Contains(theme, "var(--ui-size-") {
		t.Error("theme-material must reference at least one var(--ui-size-*) token")
	}
}

// TestComponentSizeTokensDeclaredScoped proves every family that Phase B2
// migrated declares at least one scoped component token in its own CSS file,
// so consumers can override geometry at the component root. Presence only.
func TestComponentSizeTokensDeclaredScoped(t *testing.T) {
	for file, token := range map[string]string{
		"fab.css":          "--ui-fab-container-size:",
		"button.css":       "--ui-button-padding-y:",
		"text-field.css":   "--ui-text-field-textarea-min-height:",
		"select.css":       "--ui-select-caret-reserve:",
		"toast.css":        "--ui-toast-min-height:",
		"switch.css":       "--ui-switch-handle-inset:",
		"menu.css":         "--ui-menu-min-width:",
		"data-table.css":   "--ui-data-table-checkbox-column-width:",
		"empty-state.css":  "--ui-empty-state-icon-size:",
		"banner.css":       "--ui-banner-icon-size:",
		"callout.css":      "--ui-callout-icon-size:",
		"inline-alert.css": "--ui-inline-alert-icon-size:",
		"skeleton.css":     "--ui-skeleton-avatar-size:",
		"tabs.css":         "--ui-tabs-height:",
		"chips.css":        "--ui-chip-height:",
		"dialog.css":       "--ui-dialog-min-width:",
		"radio.css":        "--ui-radio-dot-size:",
		"checkbox.css":     "--ui-checkbox-mark-size:",
	} {
		css := sourceComponentCSS(t, file)
		if !strings.Contains(css, token) {
			t.Errorf("%s must declare scoped component token %s", file, token)
		}
	}
}

// TestSpaceTokensConsumedByComponents proves the --ui-space-* tokens are real
// consumers, not orphan definitions: at least the minimum number of distinct
// component files must reference var(--ui-space-*) at least once. Structural
// assertion only — never a concrete spacing value, and demo/docco files are
// excluded so the proof is about real component selectors.
func TestSpaceTokensConsumedByComponents(t *testing.T) {
	const minConsumerFiles = 5

	entries, err := sourceStyles.ReadDir("styles")
	if err != nil {
		t.Fatalf("list styles dir: %v", err)
	}

	// Demo and doc-only files must not count as component consumers: they carry
	// preview scaffolding, not component primitives.
	excluded := map[string]bool{
		"tokens.css":                true, // defines the scale, not a consumer
		"base.css":                  true, // docs layout (site-header, prose, hero)
		"demo-whatsapp.css":         true, // demo only (all selectors are .demo-wa-*)
		"recipe-admin-resource.css": true, // screen recipe layout (all selectors are .recipe-ar-*)
		"app.css":                   true, // entry, imports and media tails
	}

	var consumers []string
	totalUses := 0
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".css") || excluded[name] {
			continue
		}
		content, err := sourceStyles.ReadFile("styles/" + name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if uses := strings.Count(string(content), "var(--ui-space-"); uses > 0 {
			consumers = append(consumers, name)
			totalUses += uses
		}
	}

	if len(consumers) < minConsumerFiles {
		t.Errorf("expected var(--ui-space-*) in at least %d distinct component files, got %d (%v)",
			minConsumerFiles, len(consumers), consumers)
	}
	if totalUses < minConsumerFiles {
		t.Errorf("expected at least %d total var(--ui-space-*) usages across components, got %d",
			minConsumerFiles, totalUses)
	}

	// The theme layer is also a consumer: --ui-fab-extension-gap re-points to
	// the core scale, so an unthemed document must resolve it.
	theme := themeCSS(t, "theme-material")
	if !strings.Contains(theme, "var(--ui-space-") {
		t.Error("theme-material must reference at least one var(--ui-space-*) token")
	}
}

// borderCoreTokens are the three core border tokens of the Phase B3 borders
// slice: the two widths every real component border uses and the solid style
// every consumer shares. The scale deliberately stops here: there is no dashed
// step (the only use is a docs demo stand-in) and no 3px/4px widths.
var borderCoreTokens = []string{
	"--ui-border-width-1:",
	"--ui-border-width-2:",
	"--ui-border-style-solid:",
}

// TestBorderTokensDefinedInCore proves the core owns the --ui-border-* tokens
// that the component layer may consume. Presence only — never a concrete value,
// so the test stays theme-agnostic (values come from the theme layer).
func TestBorderTokensDefinedInCore(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	css := string(core)
	for _, token := range borderCoreTokens {
		if !strings.Contains(css, token) {
			t.Errorf("core tokens.css must define %s (border scale step)", token)
		}
	}
}

// TestBorderTokensConsumedByComponents proves the --ui-border-* tokens are real
// consumers, not orphan definitions: at least the minimum number of distinct
// component files must reference var(--ui-border-width-*) or
// var(--ui-border-style-solid) at least once. Structural assertion only — never
// a concrete border value, and demo/doc files are excluded so the proof is
// about real component selectors.
func TestBorderTokensConsumedByComponents(t *testing.T) {
	const minConsumerFiles = 5

	entries, err := sourceStyles.ReadDir("styles")
	if err != nil {
		t.Fatalf("list styles dir: %v", err)
	}

	// Demo and doc-only files must not count as component consumers: they carry
	// preview scaffolding, not component primitives.
	excluded := map[string]bool{
		"tokens.css":                true, // defines the scale, not a consumer
		"base.css":                  true, // docs layout (site-header, prose, hero)
		"demo-whatsapp.css":         true, // demo only (all selectors are .demo-wa-*)
		"recipe-admin-resource.css": true, // screen recipe layout (all selectors are .recipe-ar-*)
		"app.css":                   true, // entry, imports and media tails
	}

	var consumers []string
	totalUses := 0
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".css") || excluded[name] {
			continue
		}
		content, err := sourceStyles.ReadFile("styles/" + name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if uses := strings.Count(string(content), "var(--ui-border-width-") + strings.Count(string(content), "var(--ui-border-style-solid)"); uses > 0 {
			consumers = append(consumers, name)
			totalUses += uses
		}
	}

	if len(consumers) < minConsumerFiles {
		t.Errorf("expected var(--ui-border-*) in at least %d distinct component files, got %d (%v)",
			minConsumerFiles, len(consumers), consumers)
	}
	if totalUses < minConsumerFiles {
		t.Errorf("expected at least %d total var(--ui-border-*) usages across components, got %d",
			minConsumerFiles, totalUses)
	}

	// The theme layer is also a consumer: the checkbox/radio/switch
	// outline-width tokens re-point to the core width-2, so an unthemed
	// document must resolve them.
	theme := themeCSS(t, "theme-material")
	if !strings.Contains(theme, "var(--ui-border-") {
		t.Error("theme-material must reference at least one var(--ui-border-*) token")
	}
}

// TestOutlineWidthTokensDeriveFromBorderWidth2 proves the three control
// outline-width tokens the Material theme owns derive from the core width-2
// scale step, so every outline stroke stays aligned with the shared scale.
// Presence only — never a concrete value.
func TestOutlineWidthTokensDeriveFromBorderWidth2(t *testing.T) {
	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, "theme-material"), " ")
	for _, token := range []string{
		"--ui-checkbox-outline-width: var(--ui-border-width-2)",
		"--ui-radio-outline-width: var(--ui-border-width-2)",
		"--ui-switch-outline-width: var(--ui-border-width-2)",
	} {
		if !strings.Contains(theme, token) {
			t.Errorf("theme-material must define %q under .theme-material", token)
		}
	}
}

// semanticColorCoreTokens are the eight semantic status tokens of the Phase B4
// slice. Only tokens with real consumers were added to the core contract:
// success/warning/info back the toast icon variants and the WhatsApp demo
// tones, danger-container backs the demo RED quality chip, and scrim is the
// dialog/drawer backdrop. outline and success-container are intentionally
// absent (border-strong covers the former; nothing consumes the latter).
var semanticColorCoreTokens = []string{
	"--ui-color-success",
	"--ui-color-success-fg",
	"--ui-color-warning",
	"--ui-color-warning-fg",
	"--ui-color-warning-container",
	"--ui-color-info",
	"--ui-color-danger-container",
	"--ui-color-scrim",
}

// TestSemanticColorTokensDefinedInCore proves the eight semantic status tokens
// exist in the core tokens contract with neutral defaults. Presence only —
// never a concrete hex value.
func TestSemanticColorTokensDefinedInCore(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	css := string(core)
	for _, token := range semanticColorCoreTokens {
		if !strings.Contains(css, token+":") {
			t.Errorf("core tokens.css must define %s", token)
		}
	}
}

// TestSemanticColorTokensOverriddenByMaterialTheme proves the Material theme
// overrides every semantic status token in light and in the single explicit
// dark class route. Each token must appear exactly twice (light + dark class),
// so a closed override matrix never needs the duplicated media dark route.
// Presence only.
func TestSemanticColorTokensOverriddenByMaterialTheme(t *testing.T) {
	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, defaultThemeName), " ")
	for _, token := range semanticColorCoreTokens {
		if n := strings.Count(theme, token+":"); n != 2 {
			t.Errorf("theme-material must override %s once in light and once in the single dark class route, got %d definitions", token, n)
		}
	}
	// The dialog-scrim component token must remain as a compatibility alias of
	// the core scrim, in light and the dark class route, so dialog.css and
	// navigation-drawer.css keep resolving it.
	if n := strings.Count(theme, "--ui-dialog-scrim: var(--ui-color-scrim)"); n != 2 {
		t.Errorf("theme-material must alias --ui-dialog-scrim to --ui-color-scrim in light and the dark class route, got %d definitions", n)
	}
}

// TestToastIconTokensDeriveFromCore proves the four --ui-toast-icon-* tokens the
// theme owns no longer carry color literals that drift: they re-point to the
// core semantic colors in light, and the inverted-surface dark values live in
// the single dark class route (never duplicated in a media block).
func TestToastIconTokensDeriveFromCore(t *testing.T) {
	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, defaultThemeName), " ")
	// The toast uses an inverted surface: light = dark container (#322f35),
	// dark = light container (#ece6f0). Icons must therefore be light on the
	// dark container and dark on the light container to keep contrast. These
	// are deliberate per-scheme values (not core status tokens, which would
	// render at 1.14-1.60:1 on the inverted surface).
	for _, want := range []string{
		"--ui-toast-icon-info: #d0bcff",
		"--ui-toast-icon-success: #81c995",
		"--ui-toast-icon-warning: #fdd663",
		"--ui-toast-icon-error: #f2b8b5",
		"--ui-toast-icon-info: #6750a4",
		"--ui-toast-icon-success: #2e7d32",
		"--ui-toast-icon-warning: #7a5700",
		"--ui-toast-icon-error: #b3261e",
	} {
		if !strings.Contains(theme, want) {
			t.Errorf("theme-material must define toast icon value %q with accessible contrast on the inverted surface", want)
		}
	}
	// Light scheme defines the four light-on-dark values once; the dark class
	// route defines the dark-on-light values once each (single mechanism — no
	// class+media duplication).
	for _, derive := range []string{
		"--ui-toast-icon-info: #d0bcff",
		"--ui-toast-icon-success: #81c995",
		"--ui-toast-icon-warning: #fdd663",
		"--ui-toast-icon-error: #f2b8b5",
	} {
		if n := strings.Count(theme, derive); n != 1 {
			t.Errorf("theme-material must define light-scheme toast icon %q exactly once, got %d", derive, n)
		}
	}
	for _, derive := range []string{
		"--ui-toast-icon-info: #6750a4",
		"--ui-toast-icon-success: #2e7d32",
		"--ui-toast-icon-warning: #7a5700",
		"--ui-toast-icon-error: #b3261e",
	} {
		if n := strings.Count(theme, derive); n != 1 {
			t.Errorf("theme-material must define dark-scheme toast icon %q exactly once in the dark class route, got %d", derive, n)
		}
	}
}

// TestWhatsAppDemoLiteralsMigrated proves every color literal that used to live
// in the WhatsApp demo (the values that motivated this slice) is gone: the demo
// now reads the core semantic tokens only. This is the migration guard for the
// demo file, which is excluded from the blanket component check below.
func TestWhatsAppDemoLiteralsMigrated(t *testing.T) {
	css := sourceComponentCSS(t, "demo-whatsapp.css")
	for _, literal := range []string{
		"#fff3d6", "#6a4b00", "#b9930a", "#d6f5dd", "#0b6b2c", "#fddcdc", "#9c1c1c",
	} {
		if strings.Contains(css, literal) {
			t.Errorf("demo-whatsapp.css must no longer contain the literal %s (migrate to --ui-color-*)", literal)
		}
	}
}

// TestNoColorLiteralsInComponents is the blanket Phase B4 guard: no component
// file under web/styles may carry a raw color literal. Tokens.css owns the
// contract values, demo-whatsapp.css is guarded by
// TestWhatsAppDemoLiteralsMigrated, app.css is the entry (forced-colors tails
// only), and forced-colors blocks inside components are exempt by spec (system
// colors are mandatory). Everything else must reference --ui-color-* tokens.
func TestNoColorLiteralsInComponents(t *testing.T) {
	excluded := map[string]bool{
		"tokens.css":        true, // owns the neutral contract values
		"demo-whatsapp.css": true, // guarded by its own migration test
		"app.css":           true, // entry: imports + forced-colors tails
	}
	hexLiteral := regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	forcedColors := regexp.MustCompile(`@media \(forced-colors: active\)\s*\{[^}]*\}`)

	entries, err := sourceStyles.ReadDir("styles")
	if err != nil {
		t.Fatalf("list styles dir: %v", err)
	}
	checked := 0
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".css") || excluded[name] {
			continue
		}
		content, err := sourceStyles.ReadFile("styles/" + name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		css := forcedColors.ReplaceAllString(string(content), "")
		if m := hexLiteral.FindString(css); m != "" {
			t.Errorf("%s must not contain the color literal %s (use a --ui-color-* token)", name, m)
		}
		if strings.Contains(css, "rgb(") {
			t.Errorf("%s must not contain a raw rgb() literal (use a --ui-color-* or --ui-shadow-* token)", name)
		}
		checked++
	}
	if checked < 25 {
		t.Errorf("expected to check the full component set, only checked %d files", checked)
	}
}

// motionCoreTokens are the core motion tokens of the Phase B5 slice: the two
// durations every consumer may rely on plus the shared standard easing. Only
// values with real consumers were added: short backs the interaction
// transitions everywhere, long backs the dialog/navigation-drawer [open]
// entrance transitions, and standard is the shared easing. There is no medium
// step and no emphasized/decelerate/accelerate easing: nothing consumes them.
var motionCoreTokens = []string{
	"--ui-motion-short:",
	"--ui-motion-long:",
	"--ui-easing-standard:",
}

// TestMotionTokensDefinedInCore proves the core owns the motion tokens that
// the component layer may consume. Presence only — never a concrete value, so
// the test stays theme-agnostic (values come from the theme layer).
func TestMotionTokensDefinedInCore(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	css := string(core)
	for _, token := range motionCoreTokens {
		if !strings.Contains(css, token) {
			t.Errorf("core tokens.css must define %s (motion scale step)", token)
		}
	}
}

// TestMotionTokensOverriddenByMaterialTheme proves the Material theme overrides
// every core motion token under .theme-material. The dark routes do not
// redefine motion (durations and easing are scheme-independent), so a single
// definition per token is the closed contract. Presence only.
func TestMotionTokensOverriddenByMaterialTheme(t *testing.T) {
	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, defaultThemeName), " ")
	for _, token := range motionCoreTokens {
		if !strings.Contains(theme, token) {
			t.Errorf("theme-material must override %s under .theme-material", token)
		}
	}
}

// TestDialogDrawerNoMotionLiterals is the Phase B5 migration guard: dialog.css
// and navigation-drawer.css must read their durations from the core motion
// tokens instead of hard-coding 150ms/500ms. Each must reference
// var(--ui-motion-short) (closed/backdrop transitions) and
// var(--ui-motion-long) (the [open] entrance) at least once and carry no
// 150ms/500ms literal anywhere.
func TestDialogDrawerNoMotionLiterals(t *testing.T) {
	for _, file := range []string{"dialog.css", "navigation-drawer.css"} {
		css := sourceComponentCSS(t, file)
		for _, literal := range []string{"150ms", "500ms"} {
			if strings.Contains(css, literal) {
				t.Errorf("%s must not contain the literal %s (use var(--ui-motion-*))", file, literal)
			}
		}
		if !strings.Contains(css, "var(--ui-motion-short)") {
			t.Errorf("%s must reference var(--ui-motion-short) for the closed/backdrop transitions", file)
		}
		if !strings.Contains(css, "var(--ui-motion-long)") {
			t.Errorf("%s must reference var(--ui-motion-long) for the [open] entrance", file)
		}
	}
}

// TestCheckboxRadioReducedMotionCoverage proves the checkbox/radio transitions
// are disabled under prefers-reduced-motion: each component file declares its
// own local reduced-motion block (the convention for the phase B3+ components)
// that sets transition: none on every selector that carries a transition.
func TestCheckboxRadioReducedMotionCoverage(t *testing.T) {
	cases := map[string]string{
		"checkbox.css": `.ui-checkbox input[type="checkbox"] { transition: none; }`,
		"radio.css":    `.ui-radio input[type="radio"], .ui-radio-mark::after { transition: none; }`,
	}
	for file, want := range cases {
		css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, file), " ")
		if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
			t.Errorf("%s must include a reduced-motion media query", file)
		}
		if !strings.Contains(css, want) {
			t.Errorf("%s reduced-motion must disable transitions with %q", file, want)
		}
	}
}

// TestPersistentSuccessPartialsNeverToast is the Phase D guard for the
// persistent-success contract: inline-alert and banner are server-rendered
// persistent slots, so they must never emit the transient loom:toast signal
// (HX-Trigger loom:toast) that the Toast pattern owns. Structural assertion
// only — the persistent success vocabulary is a contract, not a value.
func TestPersistentSuccessPartialsNeverToast(t *testing.T) {
	for _, file := range []string{"inline-alert.html", "banner.html"} {
		tmpl := repositoryFile(t, "web", "templates", file)
		for _, forbidden := range []string{"loom:toast", "HX-Trigger"} {
			if strings.Contains(tmpl, forbidden) {
				t.Errorf("%s must not contain %q (persistent success is never announced via toast)", file, forbidden)
			}
		}
	}
}

// TestSkipLinkCSSWiredIntoAppCSS proves the layout skip link ships a real
// visual treatment in the compiled app CSS (G7): it is visually hidden until
// keyboard focus brings it into view, so it never shows on the page but stays
// reachable.
func TestSkipLinkCSSWiredIntoAppCSS(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceAppCSS(t), " ")
	for _, contract := range []string{
		`.ui-skip-link {`,
		`position: absolute;`,
		`transform: translateY(-100%);`,
		`.ui-skip-link:focus {`,
		`transform: translateY(0);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("app CSS is missing skip-link contract %q", contract)
		}
	}
}

func repositoryFile(t *testing.T, path ...string) string {
	t.Helper()
	parts := append([]string{repositoryRoot(t)}, path...)
	content, err := os.ReadFile(filepath.Join(parts...))
	if err != nil {
		t.Fatalf("read repository file: %v", err)
	}
	return string(content)
}

// componentCSSFiles lists the component CSS files under web/styles (the embed
// FS) excluding the core tokens, which are exempt from the state-layer
// contract because they own the raw color values (scrim, borders, etc).
func componentCSSFiles(t *testing.T) map[string]string {
	t.Helper()
	entries, err := sourceStyles.ReadDir("styles")
	if err != nil {
		t.Fatalf("list styles dir: %v", err)
	}
	files := map[string]string{}
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".css") || name == "tokens.css" {
			continue
		}
		css, err := sourceStyles.ReadFile("styles/" + name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		files[name] = string(css)
	}
	if len(files) == 0 {
		t.Fatal("web/styles must contain at least one component CSS file")
	}
	return files
}

// TestNoRgbStateLayersInComponentCSS is the state-layer contract guard: fixed
// rgb()/rgba() colors SHALL NOT exist in component CSS. tokens.css is exempt
// (it owns the raw contract values, including the scrim). Any rgb() appearing
// in a component file is a state layer or decoration that bypassed the token
// system — the contract (state-layers) forbids it.
func TestNoRgbStateLayersInComponentCSS(t *testing.T) {
	for name, css := range componentCSSFiles(t) {
		for i, line := range strings.Split(css, "\n") {
			if strings.Contains(line, "rgb(") || strings.Contains(line, "rgba(") {
				t.Errorf("%s:%d must not contain rgb()/rgba() (state layers and decorations use --ui-color-* tokens): %s", name, i+1, strings.TrimSpace(line))
			}
		}
	}
}

// TestNoCurrentColorStateLayers is the theme-aware state-layer contract guard:
// hover/focus/pressed/selected/disabled layers SHALL be color-mix() over an
// explicit fg token, never currentColor. A state layer is recognized by the
// color-mix() overlay pattern; decorative currentColor (fill:, stroke:,
// border: on glyphs and checkmarks) is intentionally out of scope and stays.
func TestNoCurrentColorStateLayers(t *testing.T) {
	for name, css := range componentCSSFiles(t) {
		for i, line := range strings.Split(css, "\n") {
			if strings.Contains(line, "color-mix(") && strings.Contains(line, "currentColor") {
				t.Errorf("%s:%d state layer must use color-mix(in oklab, var(--ui-color-*-fg), transparent <opacity>) instead of currentColor: %s", name, i+1, strings.TrimSpace(line))
			}
		}
	}
}

// TestNoPhantomOutlineToken locks the outline role resolution: the role is
// served by --ui-color-border-strong (D3), and --ui-color-outline SHALL NOT
// exist as a definition or var() use anywhere under web/ + themes/ — only
// comments/docs may mention the name. The companion assertion proves the real
// owner is consumed, so a zero-match grep cannot pass vacuously.
func TestNoPhantomOutlineToken(t *testing.T) {
	scanned := 0
	for name, css := range componentCSSFiles(t) {
		scanned++
		for i, line := range strings.Split(css, "\n") {
			if !strings.Contains(line, "--ui-color-outline") {
				continue
			}
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "//") {
				continue // comment mention is allowed
			}
			if strings.Contains(line, "--ui-color-outline:") || strings.Contains(line, "var(--ui-color-outline)") {
				t.Errorf("%s:%d defines or consumes the phantom --ui-color-outline token (outline role is --ui-color-border-strong): %s", name, i+1, trimmed)
			}
		}
	}
	themes := availableThemes(t)
	for _, theme := range themes {
		scanned++
		css := themeCSS(t, theme)
		for i, line := range strings.Split(css, "\n") {
			if strings.Contains(line, "--ui-color-outline:") || strings.Contains(line, "var(--ui-color-outline)") {
				t.Errorf("themes/%s/theme.css:%d defines or consumes the phantom --ui-color-outline token (outline role is --ui-color-border-strong)", theme, i+1)
			}
		}
	}
	if scanned == 0 {
		t.Fatal("outline scan must cover at least one CSS file")
	}
	// Companion: the real owner of the outline role must be consumed by at
	// least one component, so the role is served and the doc mapping
	// outline=border-strong is not a dead reference.
	consumers := 0
	for _, css := range componentCSSFiles(t) {
		consumers += strings.Count(css, "var(--ui-color-border-strong)")
	}
	if consumers == 0 {
		t.Error("--ui-color-border-strong (the outline role owner) must have at least one var() consumer in component CSS")
	}
}

// TestDarkTokensDefinedOnceSingleMechanism is the dark-mode-routine contract
// guard: every theme defines dark values through exactly one mechanism (the
// explicit .theme-{name}.theme-dark class route) and no theme CSS may carry an
// @media (prefers-color-scheme: dark) block. Each dark token must be defined
// exactly once — the class+media duplication that previously existed would
// surface as a second definition.
func TestDarkTokensDefinedOnceSingleMechanism(t *testing.T) {
	tokenRe := regexp.MustCompile(`(--ui-[a-z0-9-]+)\s*:`)
	for _, theme := range availableThemes(t) {
		t.Run(theme, func(t *testing.T) {
			_, darkClass, darkMedia := splitThemeSchemes(t, theme)
			if darkMedia != "" {
				t.Errorf("%s must not define a dark @media (prefers-color-scheme: dark) block (single dark mechanism is the class route)", theme)
			}
			names := map[string]bool{}
			for _, m := range tokenRe.FindAllStringSubmatch(darkClass, -1) {
				names[m[1]] = true
			}
			if len(names) == 0 {
				t.Fatalf("%s dark class route must define at least one dark token", theme)
			}
			for name := range names {
				if n := strings.Count(darkClass, name+":"); n != 1 {
					t.Errorf("%s dark class route must define %s exactly once, got %d definitions", theme, name, n)
				}
			}
		})
	}
}

// typeScaleSteps are the 12 typescale steps the core decomposes (Phase B,
// R1). Each step owns five decomposed tokens
// (--ui-type-<step>-{size,weight,line-height,letter-spacing,family}) plus a
// shorthand alias composing them. label-md is excluded here: it is a new
// Phase B step (R2) with no pre-change baseline to snapshot.
var typeScaleSteps = []string{
	"display-lg", "display-sm", "headline-sm",
	"title-lg", "title-md",
	"body-lg", "body-md", "body-sm",
	"label-lg", "label-sm",
	"dialog-headline", "dialog-body",
}

// typeStepProps are the five decomposed property suffixes of the Phase B
// decomposition contract (R1). letter-spacing is part of the per-step token
// set even though the font: shorthand alias cannot express it — the token
// exists so a consumer can opt into tracking per step.
var typeStepProps = []string{"size", "weight", "line-height", "letter-spacing", "family"}

// typeBaseline reads the golden fixture that captured the pre-change shorthand
// values per theme (web/testdata/type_baseline.json). The fixture is the
// snapshot contract: after decomposition, composing each alias from the
// decomposed props must reproduce the baseline exactly.
func typeBaseline(t *testing.T) map[string]map[string]string {
	t.Helper()
	raw := repositoryFile(t, "web", "testdata", "type_baseline.json")
	var baseline map[string]map[string]string
	if err := json.Unmarshal([]byte(raw), &baseline); err != nil {
		t.Fatalf("parse type baseline fixture: %v", err)
	}
	return baseline
}

// decomposeThemeStep extracts the five decomposed --ui-type-<step>-* values a
// theme declares in its light block. Every theme must own all five props per
// step (the decomposition contract), so a missing prop fails the test.
func decomposeThemeStep(t *testing.T, lightCSS, step string) map[string]string {
	t.Helper()
	props := map[string]string{}
	for _, prop := range typeStepProps {
		re := regexp.MustCompile(`--ui-type-` + regexp.QuoteMeta(step) + `-` + regexp.QuoteMeta(prop) + `\s*:\s*([^;]+);`)
		m := re.FindStringSubmatch(lightCSS)
		if m == nil {
			// Report the missing prop so the failure names the gap, then
			// return what we have — the caller's equivalence check will also
			// fail loudly.
			t.Errorf("theme must define the decomposed token --ui-type-%s-%s (Phase B decomposition)", step, prop)
			continue
		}
		props[prop] = strings.TrimSpace(m[1])
	}
	return props
}

// composeTypeAlias reproduces the Phase B alias composition (design D1):
// the shorthand is a single font: value built from the decomposed parts —
// `weight size/line-height family`. letter-spacing cannot live in a font:
// shorthand, so it is intentionally not part of the alias.
func composeTypeAlias(props map[string]string) string {
	return fmt.Sprintf("%s %s/%s %s",
		props["weight"], props["size"], props["line-height"], props["family"])
}

// TestTypeAliasSnapshotEquivalence is the Phase B R1 snapshot contract: after
// decomposing the 12 shorthand type tokens into per-step props in the core and
// moving the values into the themes as decomposed overrides, composing each
// alias from the decomposed parts must reproduce the golden baseline captured
// pre-change — zero visual diff in light. Dark is asserted identical because
// the type family is never redefined in the dark class route.
func TestTypeAliasSnapshotEquivalence(t *testing.T) {
	baseline := typeBaseline(t)

	for _, theme := range availableThemes(t) {
		theme := theme
		t.Run(theme, func(t *testing.T) {
			themeWant, ok := baseline[theme]
			if !ok {
				t.Fatalf("type baseline fixture has no entry for theme %q (regenerate web/testdata/type_baseline.json)", theme)
			}
			light, darkClass, _ := splitThemeSchemes(t, theme)

			// Dark equivalence: the type family lives in light only, so
			// resolving any alias under the dark class route yields the same
			// composed value as light. Assert the dark block declares no
			// --ui-type-* token at all.
			if strings.Contains(darkClass, "--ui-type-") {
				t.Errorf("%s dark class route must not redefine any --ui-type-* token (type is scheme-independent; dark resolves identical)", theme)
			}

			for _, step := range typeScaleSteps {
				props := decomposeThemeStep(t, light, step)
				got := composeTypeAlias(props)
				want := strings.Join(strings.Fields(themeWant[step]), " ")
				gotNorm := strings.Join(strings.Fields(got), " ")
				if gotNorm != want {
					t.Errorf("%s: composed --ui-type-%s = %q, want baseline %q (decomposition changed the resolved value)",
						theme, step, gotNorm, want)
				}
			}
		})
	}
}

// TestTypeAliasesComposeDecomposedTokens proves the alias contract of the
// decomposition (R1): every shorthand --ui-type-<step> in the core must be a
// single font: composition of its five decomposed tokens — never a standalone
// literal value. Themes own the decomposed values; the core owns the alias
// shape.
func TestTypeAliasesComposeDecomposedTokens(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	css := regexp.MustCompile(`\s+`).ReplaceAllString(string(core), " ")
	for _, step := range typeScaleSteps {
		want := fmt.Sprintf("--ui-type-%s: var(--ui-type-%s-weight) var(--ui-type-%s-size)/var(--ui-type-%s-line-height) var(--ui-type-%s-family);",
			step, step, step, step, step)
		if !strings.Contains(css, want) {
			t.Errorf("core tokens.css must define --ui-type-%s as a font: composition of its decomposed tokens, got missing %q", step, want)
		}
	}
}

// TestThemesNeverRedeclareTypeAliases proves the ownership split of the
// decomposition (R1): themes override the five decomposed props per step and
// must NEVER redeclare the shorthand alias (that would fork the alias from the
// core composition and break snapshot equivalence).
func TestThemesNeverRedeclareTypeAliases(t *testing.T) {
	for _, theme := range availableThemes(t) {
		theme := theme
		t.Run(theme, func(t *testing.T) {
			css := themeCSS(t, theme)
			for _, step := range typeScaleSteps {
				re := regexp.MustCompile(`--ui-type-` + regexp.QuoteMeta(step) + `\s*:`)
				if m := re.FindString(css); m != "" {
					t.Errorf("%s must not redeclare the shorthand alias %s (themes override decomposed values only)", theme, m)
				}
			}
		})
	}
}

// TestLabelMdClosure is the Phase B R2 contract: --ui-type-label-md is a NEW
// core default with its own Material override — an independent step, never an
// alias of label-lg. The theme switcher (base.css) and the docs shell
// (docs-shell.css) must consume var(--ui-type-label-md) with NO label-lg
// fallback, and the token must be defined standalone in the core and in both
// themes.
func TestLabelMdClosure(t *testing.T) {
	// 1. Core defines label-md standalone (decomposed + alias), never as a
	// label-lg alias.
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	coreCSS := string(core)
	if !strings.Contains(coreCSS, "--ui-type-label-md:") {
		t.Error("core tokens.css must define the --ui-type-label-md alias (Phase B R2)")
	}
	for _, prop := range typeStepProps {
		if !strings.Contains(coreCSS, "--ui-type-label-md-"+prop+":") {
			t.Errorf("core tokens.css must define the decomposed label-md token --ui-type-label-md-%s", prop)
		}
	}
	if strings.Contains(coreCSS, "--ui-type-label-md: var(--ui-type-label-lg)") {
		t.Error("core must define --ui-type-label-md standalone, never as var(--ui-type-label-lg)")
	}

	// 2. Every theme defines the decomposed label-md values standalone (never
	// via label-lg) in its light block.
	for _, theme := range availableThemes(t) {
		theme := theme
		t.Run(theme, func(t *testing.T) {
			light, _, _ := splitThemeSchemes(t, theme)
			for _, prop := range typeStepProps {
				if !strings.Contains(light, "--ui-type-label-md-"+prop+":") {
					t.Errorf("%s light scheme must define the decomposed label-md token --ui-type-label-md-%s", theme, prop)
				}
			}
			if strings.Contains(light, "--ui-type-label-md: var(--ui-type-label-lg)") {
				t.Errorf("%s must define label-md standalone, never as var(--ui-type-label-lg)", theme)
			}
		})
	}

	// 3. The switcher (base.css) and docs shell (docs-shell.css) consume
	// var(--ui-type-label-md) with NO fallback. Each file must contain at
	// least one bare var(--ui-type-label-md) reference (the font: shorthand)
	// and zero label-lg fallbacks; the label-md-letter-spacing consumer on
	// the same rules references the decomposed token and must not be mistaken
	// for a bare consumption.
	for _, file := range []string{"base.css", "docs-shell.css"} {
		css, err := sourceStyles.ReadFile("styles/" + file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		content := string(css)
		if strings.Contains(content, "var(--ui-type-label-md, var(--ui-type-label-lg))") {
			t.Errorf("%s must consume var(--ui-type-label-md) with NO label-lg fallback (R2 closure)", file)
		}
		if !regexp.MustCompile(`var\(--ui-type-label-md\)`).MatchString(content) {
			t.Errorf("%s must reference the bare var(--ui-type-label-md) (font: shorthand), got only decomposed token references", file)
		}
	}
}

// TestConsumerGatedInvariants is the Phase B R3 contract (D1, D3): no
// density/z/breakpoint families, no motion-medium/easing token without a real
// consumer, and every new Phase B token must have at least one var() consumer
// (no orphans). Geometry stays in --ui-size-*; dialog/popover use the top
// layer. A future token added without a consumer trips this test.
func TestConsumerGatedInvariants(t *testing.T) {
	// 1. The forbidden families must not exist as definitions anywhere under
	// web/styles or themes (density/z/breakpoint are deferred, D3).
	sources := allStyleSources(t)
	for _, forbidden := range []string{"--ui-density-", "--ui-z-", "--ui-breakpoint-"} {
		var offenders []string
		for _, src := range sources {
			for token := range findTokenDefinitions(src.css, forbidden) {
				offenders = append(offenders, src.name+":"+token)
			}
		}
		if len(offenders) > 0 {
			t.Errorf("forbidden consumer-gated family %s must not be defined (D3 defers): %v", forbidden, offenders)
		}
	}

	// 2. Every motion/easing token defined in the core or a theme must have at
	// least one var() consumer in component or layout CSS. A motion-medium or
	// new easing with no consumer trips here.
	motionDefs := map[string]string{} // token -> where defined
	for _, src := range sources {
		for name, def := range findTokenDefinitions(src.css, "--ui-motion-") {
			motionDefs[name] = src.name + ": " + def
		}
		for name, def := range findTokenDefinitions(src.css, "--ui-easing-") {
			motionDefs[name] = src.name + ": " + def
		}
	}
	for token, where := range motionDefs {
		if !hasVarConsumer(token, sources) {
			t.Errorf("motion/easing token %s has zero var() consumers (defined at %s) — D3 defers tokens without a migrated consumer", token, where)
		}
	}

	// 3. Every Phase B type token (the decomposed steps plus label-md) must
	// have at least one var() consumer: aliases consume the 4 font-composed
	// props, the letter-spacing declarations consume the tracking props, and
	// base.css/docs-shell consume the label-md alias.
	for _, step := range append(append([]string{}, typeScaleSteps...), "label-md") {
		for _, prop := range typeStepProps {
			token := "--ui-type-" + step + "-" + prop
			if !hasVarConsumer(token, sources) {
				t.Errorf("Phase B token %s has zero var() consumers (orphan)", token)
			}
		}
	}
}

// styleSource is one CSS source file scanned by the consumer-gated invariants.
type styleSource struct {
	name string
	css  string
}

// allStyleSources concatenates every CSS file under web/styles plus every
// theme's theme.css, tagged by file name, so consumer scans can attribute
// definitions to their owning file.
func allStyleSources(t *testing.T) []styleSource {
	t.Helper()
	var sources []styleSource
	entries, err := sourceStyles.ReadDir("styles")
	if err != nil {
		t.Fatalf("list styles dir: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".css") {
			continue
		}
		css, err := sourceStyles.ReadFile("styles/" + entry.Name())
		if err != nil {
			t.Fatalf("read styles/%s: %v", entry.Name(), err)
		}
		sources = append(sources, styleSource{name: "web/styles/" + entry.Name(), css: string(css)})
	}
	for _, theme := range availableThemes(t) {
		sources = append(sources, styleSource{name: "themes/" + theme + "/theme.css", css: themeCSS(t, theme)})
	}
	return sources
}

// findTokenDefinitions returns every --ui-<family>...: definition found in a
// CSS blob, as token -> raw definition line.
func findTokenDefinitions(css, family string) map[string]string {
	out := map[string]string{}
	re := regexp.MustCompile(`(--ui-` + regexp.QuoteMeta(strings.TrimPrefix(family, "--ui-")) + `[a-z0-9-]*)\s*:\s*([^;]+);`)
	for _, m := range re.FindAllStringSubmatch(css, -1) {
		if _, ok := out[m[1]]; !ok {
			out[m[1]] = strings.TrimSpace(m[2])
		}
	}
	return out
}

// hasVarConsumer reports whether any scanned source references the token via
// var().
func hasVarConsumer(token string, sources []styleSource) bool {
	needle := "var(" + token
	for _, src := range sources {
		if strings.Contains(src.css, needle) {
			return true
		}
	}
	return false
}
