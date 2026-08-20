package lib

import (
	"regexp"
	"strings"
	"testing"
)

// Layer 1 (mobile foundations) contracts: touch-target and container-cap
// tokens, the documented breakpoint ladder, and the step-based type scale
// finding. Assertions are source-level (web/styles via the shared embed) plus
// compiled-bundle presence checks (npm run build regenerates
// web/static/app.css, committed with the change).

// TestMobileTokensDefinedInCore proves the mobile foundation tokens exist in
// the core contract: --ui-touch-target (interactive hit-area floor) and
// --ui-container-max (consumer layout cap).

// Layer 2 of the mobile UX audit — horizontal overflow containment.
//
// The user rule is absolute: horizontal overflow must NEVER be masked by
// hiding it on the document root (body { overflow-x: hidden } hides content
// loss instead of fixing it). Containment instead uses the verified pattern:
// (a) wide surfaces scroll INTERNALLY inside an overflow-x: auto container,
// (b) flex/grid children get min-width: 0 so intrinsic min widths can shrink,
// (c) nowrap / min-width:*rem is only legal inside a scroll container or a
// wrapping row — never on the page. Every test below pins one contract in
// source CSS (plus the demo template where a wrapper element is required).

// findCSSRule returns the first rule whose selector list contains want
// (comma-split, trimmed), reporting whether it was found.
func findCSSRule(t *testing.T, rules []cssRule, want string) (cssRule, bool) {
	t.Helper()
	for _, r := range rules {
		for _, sel := range strings.Split(r.selector, ",") {
			if strings.TrimSpace(sel) == want {
				return r, true
			}
		}
	}
	return cssRule{}, false
}

// maxWidth48remBlocks returns every @media (max-width: 48rem) { ... } block in
// css, brace-balanced so multi-rule blocks are captured whole (a naive regex
// only matches the first inner rule). Like reducedMotionBlocks, the scan
// starts at the media block's own opening brace.
func maxWidth48remBlocks(css string) []string {
	var out []string
	for _, m := range regexp.MustCompile(`(?s)@media\s*\(max-width:\s*48rem\)\s*\{`).FindAllStringIndex(css, -1) {
		depth := 0
		for i := m[1] - 1; i < len(css); i++ {
			switch css[i] {
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					out = append(out, css[m[0]:i+1])
					i = len(css)
				}
			}
		}
	}
	return out
}

// TestComponentPreviewContainsPageOverflow pins the Layer-2 safety net: demo
// previews scroll their own overflow (never the page) and their flex children
// get min-width: 0 so intrinsic min widths (wide tables, chat demos, chip
// rows) can shrink instead of pushing the viewport. The preview margin stays
// symmetric (margin: 2rem 0).

func TestBreakpointLadderDocumented(t *testing.T) {
	css, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	const marker = "Breakpoint ladder"
	text := string(css)
	i := strings.Index(text, marker)
	if i < 0 {
		t.Fatalf("tokens.css must document the breakpoint ladder with a %q comment marker", marker)
	}
	ladder := text[i:]
	for _, rung := range []string{"320px", "480px", "48rem", "64rem"} {
		if !strings.Contains(ladder, rung) {
			t.Errorf("breakpoint ladder comment must include rung %q", rung)
		}
	}
}

// TestMobileTokensReachCompiledBundle proves npm run build regenerates
// web/static/app.css carrying the mobile foundation tokens alongside their
// source definitions.

func TestCoreTokensStayStepBased(t *testing.T) {
	css, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	if strings.Contains(string(css), "clamp(") {
		t.Error("core tokens.css must stay step-based (GOV.UK/USWDS finding): no clamp() fluid steps")
	}
}

// TestBreakpointLadderDocumented proves the four-rung breakpoint ladder is
// documented in a tokens.css comment (GOV.UK/USWDS reference ladder). The
// marker anchors the scan so the rungs must come from the comment, not from
// token values.

func TestHeroTitleSmallScreenStep(t *testing.T) {
	css, err := sourceStyles.ReadFile("styles/hero.css")
	if err != nil {
		t.Fatalf("read hero.css: %v", err)
	}
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(string(css), " ")
	const want = "@media (max-width: 47.99rem) { .ui-hero-title { font: var(--ui-type-display-sm); letter-spacing: var(--ui-type-display-sm-letter-spacing); } }"
	if !strings.Contains(compact, want) {
		t.Errorf("hero.css must step the hero title to the small display class under 48rem: %s", want)
	}
}

func TestMobileTokensDefinedInCore(t *testing.T) {
	css, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	for _, token := range []string{"--ui-touch-target:", "--ui-container-max:"} {
		if !strings.Contains(string(css), token) {
			t.Errorf("core tokens.css must define %s (mobile foundation)", token)
		}
	}
}

// TestTouchTargetPinnedDefaultsAndThemeOwnership pins the ONE documented
// choice: the core floor is 44px (clears the GOV.UK 40px minimum), the
// Material theme overrides it to 48px (Material 3 48x48, matching USWDS's
// 48px menu button), and frozen theme-basecoat stays untouched so it
// inherits the core floor.

func TestMobileTokensReachCompiledBundle(t *testing.T) {
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(compiledAppCSS(t), " ")
	for _, token := range []string{"--ui-touch-target:", "--ui-container-max:"} {
		if !strings.Contains(compact, token) {
			t.Errorf("compiled static/app.css must carry %s (regenerate with npm run build)", token)
		}
	}
}

// TestReducedMotionGlobalRuleInBaseCSS proves the global reduced-motion
// safety net lives in base.css: under prefers-reduced-motion: reduce every
// animation and transition is killed with !important, so motion that slips
// past the per-component neutralizations cannot play for users who asked for
// less motion. The universal selector keeps the rule cheap and future-proof.

func TestMobileTokensWiredIntoComponents(t *testing.T) {
	for file, needle := range map[string]string{
		"styles/button.css":                "min-height: max(var(--ui-touch-target), var(--ui-button-min-height, var(--ui-touch-target)));",
		"styles/icon-button.css":           "width: var(--ui-touch-target); height: var(--ui-touch-target);",
		"styles/recipe-admin-resource.css": "max-width: var(--ui-container-max);",
	} {
		css := sourceComponentCSS(t, file[len("styles/"):])
		compact := regexp.MustCompile(`\s+`).ReplaceAllString(css, " ")
		if !strings.Contains(compact, needle) {
			t.Errorf("%s must consume the mobile foundation token via %q", file, needle)
		}
	}
}

// TestCoreTokensStayStepBased pins the GOV.UK/USWDS scale finding: the type
// and spacing scales stay step-based, so tokens.css must never introduce
// clamp() fluid steps.

func TestReducedMotionGlobalRuleInBaseCSS(t *testing.T) {
	css := sourceComponentCSS(t, "base.css")
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(css, " ")
	const want = "@media (prefers-reduced-motion: reduce) { * { animation: none !important; transition: none !important; } }"
	if !strings.Contains(compact, want) {
		t.Errorf("base.css must append the global reduced-motion rule: %s", want)
	}
}

// TestReducedMotionGlobalRuleReachesCompiledBundle proves npm run build
// carries the global safety net into web/static/app.css. The minifier strips
// whitespace and may reorder the killed properties, so the compiled rule is
// matched order-independently on the universal selector.

func TestReducedMotionGlobalRuleReachesCompiledBundle(t *testing.T) {
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(compiledAppCSS(t), " ")
	// RE2 has no lookaheads; the two regexes match the same universal rule
	// whichever property order the minifier produced.
	transition := regexp.MustCompile(`\*\{[^}]*transition:none!important[^}]*\}`)
	animation := regexp.MustCompile(`\*\{[^}]*animation:none!important[^}]*\}`)
	if !transition.MatchString(compact) || !animation.MatchString(compact) {
		t.Error("compiled static/app.css must carry the global *{transition/animation:none!important} rule (regenerate with npm run build)")
	}
}

// TestHeroTitleSmallScreenStep pins the 320px typography evidence: at 320px
// the hero h1 (display-lg, 56px) produces lines of 364px/322px inside a
// 254px visible hero clip region (the unbreakable 15-char word
// "Server-rendered" is ~364px at 56px), so the 56px display step silently
// clips on the smallest supported width. The step-scale fix (GOV.UK: the
// large 48px step drops to the small 32px class) is a media step under
// 48rem that falls the hero title back to the existing small display class
// (display-sm, 36px — its longest line then measures ~234px and fits).
// Body steps are untouched: body copy stays at 16px everywhere.

func TestTouchTargetPinnedDefaultsAndThemeOwnership(t *testing.T) {
	core, err := sourceStyles.ReadFile("styles/tokens.css")
	if err != nil {
		t.Fatalf("read core tokens: %v", err)
	}
	compact := regexp.MustCompile(`\s+`).ReplaceAllString(string(core), " ")
	if !strings.Contains(compact, "--ui-touch-target: 44px;") {
		t.Error("core --ui-touch-target default must be pinned at 44px (>=44px floor, GOV.UK 40px minimum)")
	}
	for _, ref := range []string{"M3", "USWDS", "GOV.UK"} {
		if !strings.Contains(string(core), ref) {
			t.Errorf("core tokens.css must document the touch-target references (%s)", ref)
		}
	}

	material := regexp.MustCompile(`\s+`).ReplaceAllString(themeCSS(t, "theme-material"), " ")
	if !strings.Contains(material, "--ui-touch-target: 48px;") {
		t.Error("theme-material must override --ui-touch-target at 48px (Material 3 48x48 touch target)")
	}

	basecoat := themeCSS(t, "theme-basecoat")
	if strings.Contains(basecoat, "--ui-touch-target") {
		t.Error("theme-basecoat is frozen and must not override --ui-touch-target (inherits the core 44px floor)")
	}
}

// TestMobileTokensWiredIntoComponents proves the foundation tokens have real
// consumers, not orphan definitions: button and icon-button reach the
// hit-area floor via the touch-target token, and the recipe consumer layout
// resolves the container cap.

func TestChipsRowsWrapInsteadOfPushing(t *testing.T) {
	rules := cssRules(t, sourceComponentCSS(t, "chips.css"))

	row, ok := findCSSRule(t, rules, ".chips-demo-row")
	if !ok || !strings.Contains(row.body, "flex-wrap: wrap") {
		t.Error(".chips-demo-row must wrap (flex-wrap: wrap) so chips never push the page")
	}
	for _, sel := range []string{".chips-demo-grid", ".chips-demo-group"} {
		r, ok := findCSSRule(t, rules, sel)
		if !ok {
			t.Errorf("chips.css must define %s", sel)
			continue
		}
		if !strings.Contains(r.body, "min-width: 0") {
			t.Errorf("%s must declare min-width: 0 so the preview can shrink it", sel)
		}
	}
	chip, ok := findCSSRule(t, rules, ".ui-chip")
	if !ok || !strings.Contains(chip.body, "white-space: nowrap") {
		t.Error(".ui-chip must keep white-space: nowrap (the row wraps whole chips, labels never wrap inside a chip)")
	}
}

// TestWhatsAppDemoContainedOnNarrowScreens pins the demo-whatsapp containment:
// the 16rem sidebar minimum relaxes under 48rem and the two panes stack, the
// appbar search shrinks instead of pushing the brand/actions off-page, and the
// chat-head window block yields to its container.

func TestComponentPreviewContainsPageOverflow(t *testing.T) {
	rules := cssRules(t, sourceComponentCSS(t, "docs-chrome.css"))

	preview, ok := findCSSRule(t, rules, ".component-preview")
	if !ok {
		t.Fatal("base.css must define .component-preview")
	}
	for _, want := range []string{
		"display: flex;",
		"flex-wrap: wrap;",
		"overflow-x: auto;",
		"margin: 2rem 0;",
	} {
		if !strings.Contains(preview.body, want) {
			t.Errorf(".component-preview must keep %q", want)
		}
	}

	children, ok := findCSSRule(t, rules, ".component-preview > *")
	if !ok {
		t.Fatal("base.css must give .component-preview children min-width: 0")
	}
	if !strings.Contains(children.body, "min-width: 0") {
		t.Error(".component-preview > * must declare min-width: 0 so intrinsic min widths shrink inside the preview")
	}
}

// TestNoBodyOverflowXHiddenMasking is the user's hard rule: overflow must
// never be masked by hiding it on the document root. No styles file may
// declare overflow-x: hidden (or the overflow shorthand) on a body/html
// selector.

func TestDataTableWideTablesScrollInside(t *testing.T) {
	css := sourceComponentCSS(t, "data-table.css")
	rules := cssRules(t, css)

	scroll, ok := findCSSRule(t, rules, ".ui-data-table-scroll")
	if !ok {
		t.Fatal("data-table.css must define .ui-data-table-scroll")
	}
	if !strings.Contains(scroll.body, "overflow-x: auto") {
		t.Error(".ui-data-table-scroll must scroll horizontally (overflow-x: auto)")
	}

	// The wrapper must sit between the bordered container and the <table> so
	// the table's overflow lands on the wrapper's scrollbar, never on the
	// page. Wrapping outside the container would be clipped by the container's
	// overflow: hidden.
	tmpl := repositoryFile(t, "lib", "templates", "data-table.html")
	if !regexp.MustCompile(`<div class="ui-data-table-scroll">\s*<table class="ui-data-table-table">`).MatchString(tmpl) {
		t.Error("data-table.html must open .ui-data-table-scroll directly around the <table>")
	}
	if !strings.Contains(tmpl, `<div class="ui-data-table" id="data-table-panel">`) {
		t.Error("data-table.html must keep the .ui-data-table panel as the hx-target")
	}

	// min-width: 8rem items live in flex-wrap rows, so they wrap instead of
	// pushing.
	for _, sel := range []string{".data-table-demo-filter", ".data-table-demo-refresh-row"} {
		r, ok := findCSSRule(t, rules, sel)
		if !ok || !strings.Contains(r.body, "flex-wrap: wrap") {
			t.Errorf("%s must wrap (flex-wrap: wrap) so min-width: 8rem items stay on-page", sel)
		}
	}

	// The nowrap cells stay legal because they self-clip (ellipsis) inside the
	// scroll wrapper, and the backdrop of the whole chain is the component
	// preview safety net.
	label, ok := findCSSRule(t, rules, "tbody .ui-data-table-cell--label")
	if !ok {
		t.Fatal("data-table.css must style tbody .ui-data-table-cell--label")
	}
	for _, want := range []string{"overflow: hidden;", "text-overflow: ellipsis;", "white-space: nowrap;"} {
		if !strings.Contains(label.body, want) {
			t.Errorf("tbody .ui-data-table-cell--label must keep %q (self-clipping nowrap cell)", want)
		}
	}
}

// TestChipsRowsWrapInsteadOfPushing pins the chips containment: chips keep
// their nowrap pill labels because the demo row wraps whole chips (flex-wrap),
// and the grid/group chain gets min-width: 0 so the wrapping can actually
// happen when the preview shrinks. No scroll container is needed: the
// component's design is wrapping, not scrolling.

func TestExistingProseScrollPatternsStayPinned(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "docs-chrome.css"), " ")
	for _, contract := range []string{
		`.prose pre { margin: 1rem 0; padding: 1rem 1.25rem; overflow-x: auto;`,
		`@media (max-width: 48rem) { .prose table { display: block; overflow-x: auto; } }`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("base.css must keep established prose scroll pattern %q", contract)
		}
	}
}

func TestNoBodyOverflowXHiddenMasking(t *testing.T) {
	files, err := sourceStyles.ReadDir("styles")
	if err != nil {
		t.Fatalf("read styles dir: %v", err)
	}
	for _, f := range files {
		name := f.Name()
		if !strings.HasSuffix(name, ".css") {
			continue
		}
		raw, err := sourceStyles.ReadFile("styles/" + name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		for _, rule := range cssRules(t, string(raw)) {
			for _, sel := range strings.Split(rule.selector, ",") {
				sel = strings.TrimSpace(sel)
				if sel != "body" && sel != "html" && !strings.HasPrefix(sel, "body.") && !strings.HasPrefix(sel, "html.") {
					continue
				}
				if strings.Contains(rule.body, "overflow-x: hidden") || strings.Contains(rule.body, "overflow: hidden") {
					t.Errorf("%s: selector %q must not hide page overflow (user rule: contain, never mask)", name, sel)
				}
			}
		}
	}
}

// TestDataTableWideTablesScrollInside pins the data-table containment: the
// nowrap label cells keep their ellipsis, but the table itself is wrapped in
// .ui-data-table-scroll (overflow-x: auto) inside the demo template, so wide
// tables scroll internally instead of pushing the page.

func TestWhatsAppDemoContainedOnNarrowScreens(t *testing.T) {
	css := sourceComponentCSS(t, "demo-whatsapp.css")

	blocks := maxWidth48remBlocks(css)
	if len(blocks) == 0 {
		t.Fatal("demo-whatsapp.css must include an @media (max-width: 48rem) containment block")
	}
	var mediaRules []cssRule
	for _, b := range blocks {
		mediaRules = append(mediaRules, cssRules(t, b)...)
	}

	layout, ok := findCSSRule(t, mediaRules, ".demo-wa-layout")
	if !ok || !strings.Contains(layout.body, "flex-wrap: wrap") {
		t.Error("mobile block: .demo-wa-layout must wrap so sidebar and chat stack")
	}
	sidebar, ok := findCSSRule(t, mediaRules, ".demo-wa-sidebar")
	if !ok {
		t.Error("mobile block: .demo-wa-sidebar must relax its 16rem minimum")
	} else {
		if !strings.Contains(sidebar.body, "min-width: 0") {
			t.Error("mobile block: .demo-wa-sidebar must declare min-width: 0 (16rem minimum relaxed)")
		}
		if !strings.Contains(sidebar.body, "flex: 1 1 100%") {
			t.Error("mobile block: .demo-wa-sidebar must take its own full line (flex: 1 1 100%)")
		}
	}
	chat, ok := findCSSRule(t, mediaRules, ".demo-wa-chat")
	if !ok || !strings.Contains(chat.body, "flex: 1 1 100%") {
		t.Error("mobile block: .demo-wa-chat must take its own full line (flex: 1 1 100%)")
	}
	headWindow, ok := findCSSRule(t, mediaRules, ".demo-wa-chat-head-window")
	if !ok || !strings.Contains(headWindow.body, "min-width: 0") {
		t.Error("mobile block: .demo-wa-chat-head-window must declare min-width: 0")
	}
	windowBar, ok := findCSSRule(t, mediaRules, ".demo-wa-window-bar")
	if !ok || !strings.Contains(windowBar.body, "min(12rem, 100%)") {
		t.Error("mobile block: .demo-wa-window-bar must yield to the container (width: min(12rem, 100%))")
	}

	// Appbar: the search field may shrink to zero rather than push anything
	// off-page; declared on the base rule (harmless on desktop).
	search, ok := findCSSRule(t, cssRules(t, css), ".demo-wa-search")
	if !ok || !strings.Contains(search.body, "min-width: 0") {
		t.Error(".demo-wa-search must declare min-width: 0 so the appbar never pushes the page")
	}
}

// TestExistingProseScrollPatternsStayPinned guards the established good
// patterns Layer 2 builds on: code blocks and markdown tables scroll inside
// their own containers on narrow screens. Regression protection — never break
// these while containing new surfaces.
