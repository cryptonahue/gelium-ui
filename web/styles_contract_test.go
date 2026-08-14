package web

import (
	"embed"
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
// gap is closed: the core owns the neutral shorthand default and the Material
// theme defines the Material value. Presence only — never a concrete value.
func TestDisplayLgTokenClosedAcrossCoreAndTheme(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	if !strings.Contains(string(core), "--ui-type-display-lg:") {
		t.Error("core tokens.css must define --ui-type-display-lg (neutral default)")
	}

	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, "theme-material"), " ")
	if !strings.Contains(theme, "--ui-type-display-lg:") {
		t.Error("theme-material must define --ui-type-display-lg")
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
// is closed: the core owns the neutral shorthand default (consumed by the
// WhatsApp demo admin subtitle) and the Material theme defines its override.
// Presence only — never a concrete value.
func TestTitleMdTokenClosedAcrossCoreAndTheme(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	if !strings.Contains(string(core), "--ui-type-title-md:") {
		t.Error("core tokens.css must define --ui-type-title-md (neutral default)")
	}

	theme := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, "theme-material"), " ")
	if !strings.Contains(theme, "--ui-type-title-md:") {
		t.Error("theme-material must define --ui-type-title-md")
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
