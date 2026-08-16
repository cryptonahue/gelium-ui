package web

import (
	"regexp"
	"strings"
	"testing"
)

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
func TestComponentPreviewContainsPageOverflow(t *testing.T) {
	rules := cssRules(t, sourceComponentCSS(t, "base.css"))

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
	tmpl := repositoryFile(t, "web", "templates", "data-table.html")
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
func TestExistingProseScrollPatternsStayPinned(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "base.css"), " ")
	for _, contract := range []string{
		`.prose pre { margin: 1rem 0; padding: 1rem 1.25rem; overflow-x: auto;`,
		`@media (max-width: 48rem) { .prose table { display: block; overflow-x: auto; } }`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("base.css must keep established prose scroll pattern %q", contract)
		}
	}
}
