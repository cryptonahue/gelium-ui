package app

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"geliumui/lib"
)

// docsNavLink is one destination in the docs shell sidebar / footer export.
// Path is the canonical route (identity, active match, footer). Href is what
// the sidebar renders — Path plus an optional ?theme=<slug> so navigating the
// IA does not drop the document-root theme selection.
type docsNavLink struct {
	Path    string
	Href    string
	Label   string
	Current bool // exact path match against the active request path
}

// docsNavGroup is a labeled block of docs nav links (IA section).
type docsNavGroup struct {
	Title string
	Links []docsNavLink
}

// docsNavView is the shell chrome model: grouped IA plus honest topbar slots.
// Version aligns with the static asset query (?v=0.5.0). SearchIndex is the
// JSON index of every nav destination ({title, href, group}) that search.js
// filters client-side; ThemeSlug/Scheme let the search form preserve the
// current chrome on its 0-JS GET fallback.
type docsNavView struct {
	Groups      []docsNavGroup
	Version     string
	SearchIndex template.JS // JSON [{title,href,group}] of the whole nav model
	ThemeSlug   string      // current ?theme= slug (search form hidden input)
	Scheme      string      // current ?scheme= value (search form hidden input)
	// ChromeQuery is the allowlisted ?theme=&scheme= suffix for topbar links
	// that are NOT sidebar destinations (brand, blog, changelog). Without it
	// boosted navigation from those links silently drops the selected
	// theme/scheme and the server renders the default.
	ChromeQuery string
}

// docsShellVersion is the static version badge shown in the docs topbar.
// Derived from the single cache-busting version constant so the badge, the
// npm package and the asset cache-buster can never drift (REQ-6.2).
const docsShellVersion = lib.AssetsVersion

// usesDocsShell reports whether path should render the two-pane docs chrome.
// Home, recipes, and demos stay on the legacy site-header layout.
func usesDocsShell(path string) bool {
	return path == "/docs" ||
		strings.HasPrefix(path, "/docs/") ||
		strings.HasPrefix(path, "/components/")
}

// themeSlugFromClass maps a resolved theme class (theme-basecoat) to the
// public ?theme= slug (basecoat). Empty when the class is unknown.
func themeSlugFromClass(class string) string {
	for _, t := range availableThemes {
		if t.Class == class {
			return t.Slug
		}
	}
	return ""
}

// docsNavFor builds the Scalar-style sidebar IA for activePath.
// Exact path match sets Current; empty activePath marks nothing current
// (used when flattening the same model into the site footer).
// themeSlug/scheme, when allowlisted, are appended on every sidebar Href so
// in-shell navigation keeps direction and light/dark selection.
func docsNavFor(activePath, themeSlug, scheme string) docsNavView {
	link := func(path, label string) docsNavLink {
		return docsNavLink{
			Path:    path,
			Href:    chromeHref(path, themeSlug, scheme),
			Label:   label,
			Current: activePath != "" && path == activePath,
		}
	}

	groups := make([]docsNavGroup, 0, 5+len(docsSections))
	groups = append(groups, docsNavGroup{
		Title: "Getting started",
		Links: []docsNavLink{link("/docs", "Documentation")},
	})
	// Handbook is the concept layer: it sits at position 2, right after
	// Getting started and BEFORE the component reference sections, so
	// onboarding precedes lookup (Information architecture handbook page).
	handbook := make([]docsNavLink, 0, len(handbookNavLinks))
	for _, l := range handbookNavLinks {
		handbook = append(handbook, link(l.Path, l.Label))
	}
	groups = append(groups, docsNavGroup{Title: "Handbook", Links: handbook})
	for _, section := range docsSections {
		links := make([]docsNavLink, 0, len(section.Links))
		for _, l := range section.Links {
			links = append(links, link(l.Path, l.Label))
		}
		groups = append(groups, docsNavGroup{Title: section.Title, Links: links})
	}
	groups = append(groups,
		docsNavGroup{
			Title: "Patterns",
			Links: []docsNavLink{link("/docs/patterns", "Patterns")},
		},
		docsNavGroup{
			Title: "Recipes",
			Links: []docsNavLink{
				link("/recipes/admin-resource", "Admin Resource"),
				link("/recipes/ops-queue", "Ops Queue"),
				link("/recipes/public-feed", "Public Feed"),
			},
		},
	)

	return docsNavView{
		Groups:      groups,
		Version:     docsShellVersion,
		SearchIndex: buildSearchIndex(groups),
		ThemeSlug:   themeSlug,
		Scheme:      scheme,
		ChromeQuery: chromeQuery(themeSlug, scheme),
	}
}

// searchIndexEntry is one destination in the client-side docs search index.
type searchIndexEntry struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Group string `json:"group"`
}

// buildSearchIndex flattens the nav model into the JSON index search.js
// filters: every sidebar destination with its group label. Href carries the
// same chrome query as the sidebar links, so search navigation preserves the
// selected theme and light/dark scheme.
func buildSearchIndex(groups []docsNavGroup) template.JS {
	entries := make([]searchIndexEntry, 0, 64)
	for _, g := range groups {
		for _, l := range g.Links {
			entries = append(entries, searchIndexEntry{Title: l.Label, Href: l.Href, Group: g.Title})
		}
	}
	b, err := json.Marshal(entries)
	if err != nil {
		return template.JS("[]")
	}
	return template.JS(b)
}

// orderedDocsNav is the flat, ordered IA shared by the sidebar AND the
// previous/next pagination: Getting started, Handbook, component sections,
// Patterns, Recipes. One model, two renderings — they can never drift.
func orderedDocsNav() []navLink {
	links := make([]navLink, 0, 48)
	links = append(links, navLink{Path: "/docs", Label: "Documentation"})
	links = append(links, handbookNavLinks...)
	for _, section := range docsSections {
		links = append(links, section.Links...)
	}
	links = append(links, navLink{Path: "/docs/patterns", Label: "Patterns"})
	links = append(links,
		navLink{Path: "/recipes/admin-resource", Label: "Admin Resource"},
		navLink{Path: "/recipes/ops-queue", Label: "Ops Queue"},
		navLink{Path: "/recipes/public-feed", Label: "Public Feed"},
	)
	return links
}

// prevNextFor returns the previous/next destinations around activePath in the
// flat IA order (GOV.UK pattern). Both hrefs carry the allowlisted chrome
// query so pagination never silently resets theme/scheme. nil on the first
// page (no previous) or the last page (no next); nil entirely when activePath
// is not part of the ordered IA.
func prevNextFor(activePath, themeSlug, scheme string) *prevNextView {
	ordered := orderedDocsNav()
	for i, l := range ordered {
		if l.Path != activePath {
			continue
		}
		pn := &prevNextView{}
		if i > 0 {
			pn.Prev = &prevNextLink{Href: chromeHref(ordered[i-1].Path, themeSlug, scheme), Label: ordered[i-1].Label}
		}
		if i < len(ordered)-1 {
			pn.Next = &prevNextLink{Href: chromeHref(ordered[i+1].Path, themeSlug, scheme), Label: ordered[i+1].Label}
		}
		return pn
	}
	return nil
}

// handbookNavLinks is the Gelium Handbook IA: the concept pages that explain
// how the library works, ordered first-to-read first. It drives the sidebar
// group (position 2, before the component reference), the /docs hub section,
// and the sitemap so the destinations can never drift.
var handbookNavLinks = []navLink{
	{Path: "/docs/information-architecture", Label: "Information architecture"},
	{Path: "/docs/screens", Label: "Screens"},
	{Path: "/docs/journeys", Label: "Journeys"},
	{Path: "/docs/data-display", Label: "Data display"},
	{Path: "/docs/feedback", Label: "Feedback"},
	{Path: "/docs/density", Label: "Density"},
	{Path: "/docs/motion", Label: "Motion"},
	{Path: "/docs/ui-definition-of-done", Label: "UI definition of done"},
	{Path: "/docs/choose-the-right-control", Label: "Choose the right control"},
	{Path: "/docs/forms", Label: "Forms"},
	{Path: "/docs/compare", Label: "Why Gelium"},
	{Path: "/docs/performance", Label: "Performance"},
	{Path: "/docs/responsive", Label: "Responsive"},
	{Path: "/docs/themes", Label: "Themes"},
	{Path: "/docs/tokens", Label: "Tokens"},
	{Path: "/docs/server-contracts", Label: "Server contracts"},
	{Path: "/docs/accessibility", Label: "Accessibility"},
	{Path: "/docs/principles", Label: "Design principles"},
	{Path: "/docs/browser-support", Label: "Browser support"},
	{Path: "/docs/content-style", Label: "Content style"},
	{Path: "/docs/acknowledgments", Label: "Acknowledgments"},
	{Path: "/docs/contributing", Label: "Contributing"},
	{Path: "/docs/changelog", Label: "Changelog"},
	{Path: "/docs/roadmap", Label: "Roadmap"},
}

// docsSections groups the component library into logical categories for the
// /docs index page. Each entry links to the dogfooded component page.
var docsSections = []struct {
	Title string
	Intro string
	Links []navLink
}{
	{
		Title: "Foundation",
		Intro: "The primitives every component builds on.",
		Links: []navLink{
			{Path: "/components/elevation", Label: "Elevation"},
			{Path: "/components/focus-ring", Label: "Focus ring"},
			{Path: "/components/icon", Label: "Icon"},
			{Path: "/components/divider", Label: "Divider"},
		},
	},
	{
		Title: "Actions",
		Intro: "Buttons and single-action controls.",
		Links: []navLink{
			{Path: "/components/button", Label: "Button"},
			{Path: "/components/icon-button", Label: "Icon button"},
			{Path: "/components/fab", Label: "FAB"},
			{Path: "/components/chips", Label: "Chips"},
			{Path: "/components/segmented-button", Label: "Segmented buttons"},
			{Path: "/components/menu", Label: "Menu"},
		},
	},
	{
		Title: "Input",
		Intro: "Native form controls with Material styling.",
		Links: []navLink{
			{Path: "/components/text-field", Label: "Text field"},
			{Path: "/components/checkbox", Label: "Checkbox"},
			{Path: "/components/radio", Label: "Radio"},
			{Path: "/components/switch", Label: "Switch"},
			{Path: "/components/select", Label: "Select"},
			{Path: "/components/slider", Label: "Slider"},
			{Path: "/components/list", Label: "List"},
		},
	},
	{
		Title: "Feedback & status",
		Intro: "Communicate state and progress.",
		Links: []navLink{
			{Path: "/components/dialog", Label: "Dialog"},
			{Path: "/components/toast", Label: "Toast"},
			{Path: "/components/progress", Label: "Progress"},
			{Path: "/components/badge", Label: "Badge"},
			{Path: "/components/card", Label: "Card"},
			{Path: "/components/tooltip", Label: "Tooltip"},
		},
	},
	{
		Title: "Navigation",
		Intro: "Real links and server-derived active state.",
		Links: []navLink{
			{Path: "/components/tabs", Label: "Tabs"},
			{Path: "/components/navigation-bar", Label: "Navigation bar"},
			{Path: "/components/navigation-tab", Label: "Navigation tab"},
			{Path: "/components/navigation-drawer", Label: "Navigation drawer"},
		},
	},
	{
		Title: "Data",
		Intro: "Server-side tables.",
		Links: []navLink{
			{Path: "/components/data-table", Label: "Data table"},
		},
	},
}

// docsRootLeadSource is the embedded docs root (web/content/index.md) that
// leads the generated /docs hub. It is embedded in the bundle so the hub can
// explain what Gelium UI is without duplicating prose in Go.
const docsRootLeadSource = "content/index.md"

// stripDocsRootH1 removes a single leading H1 line so an embedded docs root
// can lead the /docs hub without creating a second heading level 1 on the page.
// Anything after the H1 — including further headings — is preserved verbatim;
// blank lines left by the stripped heading are collapsed.
func stripDocsRootH1(source string) string {
	if strings.HasPrefix(source, "# ") {
		if i := strings.IndexByte(source, '\n'); i >= 0 {
			return strings.TrimLeft(source[i+1:], "\n")
		}
		return ""
	}
	return source
}

// docsIndex is the GET /docs handler. It is an orientation hub, not a second
// sidebar: the embedded docs root (content/index.md) explains what Gelium is
// and how to start, then a short "Start here" list points at high-value
// handbook pages. The full Handbook and component catalog live in the docs
// shell sidebar (and footer) — this body must not re-list every destination.
func (s *server) docsIndex(w http.ResponseWriter, r *http.Request) {
	var md string
	if source, err := fs.ReadFile(s.assets, docsRootLeadSource); err == nil {
		if lead := stripDocsRootH1(string(source)); strings.TrimSpace(lead) != "" {
			md += lead + "\n\n"
		}
	}
	md += "# Documentation\n\n"
	md += "This hub orients you. The **sidebar** is the map: Handbook concepts and the full component reference stay there so this page does not repeat them.\n\n"
	md += "## Start here\n\n"
	md += "High-value entry points (not the full catalog):\n\n"
	md += "- [Screens](/docs/screens) — screen types, hierarchy, nav patterns (sourced criteria).\n"
	md += "- [Journeys](/docs/journeys) — multi-step flows and where to land after submit.\n"
	md += "- [Data display](/docs/data-display) — table vs list vs cards (DATA-*).\n"
	md += "- [Feedback](/docs/feedback) — toast vs summary vs banner vs empty (decision matrix).\n"
	md += "- [Patterns](/docs/patterns) — domain skeletons (forum, catalog, admin).\n"
	md += "- [UI definition of done](/docs/ui-definition-of-done) — ship checklist for humans and agents.\n"
	md += "- [Why Gelium](/docs/compare) — when to use Gelium vs React/headless kits, and explicit no-gos.\n"
	md += "- [Forms contract](/docs/forms) — labels, `inputmode`/`type`, autocomplete, validate after interaction.\n"
	md += "- [Themes](/docs/themes) — class-based direction; dark is a class route, not media-only.\n"
	md += "- [Tokens](/docs/tokens) — `--ui-*` ownership (core vs theme vs component).\n"
	md += "- [Performance](/docs/performance) — ~50KB JS stance; CSS is the largest asset by design.\n"
	md += "- [Responsive](/docs/responsive) — design for screen sizes, not devices.\n"
	md += "- [npm `gelium-ui`](https://www.npmjs.com/package/gelium-ui) — install the package consumers get.\n"
	md += "- [Agent brief](/llms.txt) — machine-readable project summary.\n"
	md += "- [Agent UX pack](/llms-ux.txt) — screen/feedback decision tables for LLMs.\n\n"
	md += "## Use the sidebar\n\n"
	md += "Every Handbook page and every component is already linked in the left nav (desktop) and the mobile docs menu. Jump from there for reference; come back here for the story and the start list.\n\n"
	md += "## Try a full screen\n\n"
	md += "Composed recipes and demos (also in the nav under Recipes / outside the handbook):\n\n"
	md += "- [Admin Resource](/recipes/admin-resource) — data table + form + dialog on the server contract.\n"
	md += "- [Ops Queue](/recipes/ops-queue) — work queue with POST+303 transitions.\n"
	md += "- [Public/Social Feed](/recipes/public-feed) — feed with reactions and loading states.\n"
	md += "- [WhatsApp manager](/demo/whatsapp) — larger composed demo.\n"
	s.renderMarkdown(w, r, pageView{Title: "Documentation"}, md, "/docs")
}

// docsPatterns is GET /docs/patterns — domain skeletons and recipe affinity.
func (s *server) docsPatterns(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Patterns"}, "content/handbook-patterns.md", "/docs/patterns")
}

// docsInformationArchitecture is GET /docs/information-architecture — the
// Information architecture handbook page: the concept-before-reference
// hierarchy rule, the criteria for adding a group or page, and the agent
// prompt that lets LLMs evaluate or improve the docs IA.
func (s *server) docsInformationArchitecture(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Information architecture"}, "content/handbook-information-architecture.md", "/docs/information-architecture")
}

// docsChooseTheRightControl is GET /docs/choose-the-right-control — the
// control-selection handbook page: the decision table and rules of thumb
// for picking the right input component per situation.
func (s *server) docsChooseTheRightControl(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Choose the right control"}, "content/handbook-choose-the-right-control.md", "/docs/choose-the-right-control")
}

// docsScreens is GET /docs/screens — screen types, hierarchy, nav patterns,
// and the build checklist adapted from GOV.UK / USWDS / M3 / NNG sources.
func (s *server) docsScreens(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Screens"}, "content/handbook-screens.md", "/docs/screens")
}

// docsFeedback is GET /docs/feedback — decision matrix for toast vs
// validation-summary vs banner vs empty/error/skeleton (sourced criteria).
func (s *server) docsFeedback(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Feedback"}, "content/handbook-feedback.md", "/docs/feedback")
}

// docsJourneys is GET /docs/journeys — multi-step shapes and post-submit landings.
func (s *server) docsJourneys(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Journeys"}, "content/handbook-journeys.md", "/docs/journeys")
}

// docsDataDisplay is GET /docs/data-display — table vs list vs cards (DATA-*).
func (s *server) docsDataDisplay(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Data display"}, "content/handbook-data-display.md", "/docs/data-display")
}

// docsDensity is GET /docs/density — comfortable/cozy/compact and app shell.
func (s *server) docsDensity(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Density"}, "content/handbook-density.md", "/docs/density")
}

// docsMotion is GET /docs/motion — when to animate; reduced-motion policy.
func (s *server) docsMotion(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Motion"}, "content/handbook-motion.md", "/docs/motion")
}

// docsUIDefinitionOfDone is GET /docs/ui-definition-of-done — ship checklist.
func (s *server) docsUIDefinitionOfDone(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "UI definition of done"}, "content/handbook-ui-definition-of-done.md", "/docs/ui-definition-of-done")
}

// docsForms is GET /docs/forms — the Forms contract handbook page: labels,
// type/inputmode pairing, autocomplete, validation timing, and native-first
// rules that apply before any component page.
func (s *server) docsForms(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Forms"}, "content/handbook-forms.md", "/docs/forms")
}

// docsCompare is GET /docs/compare — Why Gelium: honest comparison vs
// React/headless kits (Radix, shadcn, Base UI), payload orders of magnitude,
// when to use and explicit no-gos.
func (s *server) docsCompare(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Why Gelium"}, "content/handbook-compare.md", "/docs/compare")
}

// docsPerformance is GET /docs/performance — product stance on payload:
// JS as progressive enhancement, large CSS (tokens+themes) by design,
// how to measure static assets and npm pack size.
func (s *server) docsPerformance(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Performance"}, "content/handbook-performance.md", "/docs/performance")
}

// docsResponsive is GET /docs/responsive — design for viewports and content
// reflow, not device names; containment without overflow-x:hidden masking;
// --ui-touch-target / --ui-container-max and prose 65ch.
func (s *server) docsResponsive(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Responsive"}, "content/handbook-responsive.md", "/docs/responsive")
}

// docsThemes is GET /docs/themes — the Themes handbook page: how themes work
// over one markup surface, the Material default and Basecoat direction, and
// the 0-JS ?theme=/class selection routes hosted in the docs topbar.
func (s *server) docsThemes(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Themes"}, "content/handbook-themes.md", "/docs/themes")
}

// docsTokens is GET /docs/tokens — the Tokens handbook page: the --ui-* token
// vocabulary, core families, theme-owned values, and naming conventions.
func (s *server) docsTokens(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Tokens"}, "content/handbook-tokens.md", "/docs/tokens")
}

// docsServerContracts is GET /docs/server-contracts — the server contract
// handbook page: GET+query state, POST+303 mutations, 422 validation, and
// HX-Trigger toast feedback.
func (s *server) docsServerContracts(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Server contracts"}, "content/handbook-server-contracts.md", "/docs/server-contracts")
}

// docsAccessibility is GET /docs/accessibility — the accessibility handbook
// page: native semantics, accessible names, states, focus, live regions,
// reduced motion, forced colors, contrast, and RTL.
func (s *server) docsAccessibility(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Accessibility"}, "content/handbook-accessibility.md", "/docs/accessibility")
}

// docsPrinciples is GET /docs/principles — the Design Principles page: the
// four foundation principles with what/why/example each, plus how tests
// enforce them.
func (s *server) docsPrinciples(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Design principles"}, "content/principles.md", "/docs/principles")
}

// docsBrowserSupport is GET /docs/browser-support — the Browser support
// handbook page: the consolidated Baseline API table (Popover, anchor
// positioning, Invoker Commands, dialog closedby) and the "what always
// works" contract (native semantics, 0-JS, AA contrast, forced-colors).
func (s *server) docsBrowserSupport(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Browser support"}, "content/handbook-browser-support.md", "/docs/browser-support")
}

// docsContributing is GET /docs/contributing — the Contributing page: where
// the project lives, the development setup, the gates every contribution
// must pass, commit conventions, and the PR workflow.
func (s *server) docsContributing(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Contributing"}, "content/handbook-contributing.md", "/docs/contributing")
}

// docsContentStyle is GET /docs/content-style — the Content style handbook
// page: the one voice every Gelium UI string follows (errors, toasts, empty
// states, banners, docs), with the action-pattern rules and the copy contract
// tests that enforce them.
func (s *server) docsContentStyle(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Content style"}, "content/handbook-content-style.md", "/docs/content-style")
}

// docsAcknowledgments is GET /docs/acknowledgments — the Acknowledgments page:
// an honest record of the design systems and component libraries that
// inspired Gelium UI, what was taken from each, and how it was adapted to the
// Gelium model (server-rendered, 0-JS, token-driven).
func (s *server) docsAcknowledgments(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Acknowledgments"}, "content/handbook-acknowledgments.md", "/docs/acknowledgments")
}

// docsChangelog is GET /docs/changelog — the Changelog page: the full project
// process documented per version, mirrored from the repository CHANGELOG.md.
func (s *server) docsChangelog(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Changelog"}, "content/handbook-changelog.md", "/docs/changelog")
}

// docsRoadmap is GET /docs/roadmap — the public Roadmap page: phases A-J
// shipped with contract tests, docs/DX status, and the post-A-J next list,
// mirrored from the internal system roadmap.
func (s *server) docsRoadmap(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Roadmap"}, "content/handbook-roadmap.md", "/docs/roadmap")
}
