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
	// Handbook criteria split into Core / System / Meta so the sidebar stays
	// scannable (USWDS-style: don't dump one giant concept list). All three
	// sit after Getting started and BEFORE component reference sections.
	for _, section := range handbookSections {
		links := make([]docsNavLink, 0, len(section.Links))
		for _, l := range section.Links {
			links = append(links, link(l.Path, l.Label))
		}
		groups = append(groups, docsNavGroup{Title: section.Title, Links: links})
	}
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
				link("/recipes/rich-article", "Rich Article"),
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

// docsNavForNavigation preserves the source vocabulary of the active selection:
// legacy theme/scheme links stay legacy, while recipe fields stay canonical only
// inside recipe/gallery contexts selected by navigationSelectionFor.
func docsNavForNavigation(activePath string, selection navigationSelection) docsNavView {
	nav := docsNavFor(activePath, selection.themeSlug, selection.scheme)
	if selection.canonical {
		for gi := range nav.Groups {
			for li := range nav.Groups[gi].Links {
				link := &nav.Groups[gi].Links[li]
				link.Href = selection.href(link.Path)
			}
		}
		nav.SearchIndex = buildSearchIndex(nav.Groups)
	}
	nav.ChromeQuery = selection.query()
	return nav
}

// docsNavForSelection remains available to focused callers that explicitly
// request canonical recipe links.
func docsNavForSelection(activePath string, selection documentSelection, execution accordionExecution) docsNavView {
	return docsNavForNavigation(activePath, navigationSelection{canonical: true, selection: selection, execution: execution})
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
// previous/next pagination: Getting started, Core/System/Meta handbook,
// component sections, Patterns, Recipes. One model, two renderings.
func orderedDocsNav() []navLink {
	links := make([]navLink, 0, 48)
	links = append(links, navLink{Path: "/docs", Label: "Documentation"})
	links = append(links, handbookNavLinks()...)
	for _, section := range docsSections {
		links = append(links, section.Links...)
	}
	links = append(links, navLink{Path: "/docs/patterns", Label: "Patterns"})
	links = append(links,
		navLink{Path: "/recipes/admin-resource", Label: "Admin Resource"},
		navLink{Path: "/recipes/ops-queue", Label: "Ops Queue"},
		navLink{Path: "/recipes/public-feed", Label: "Public Feed"},
		navLink{Path: "/recipes/rich-article", Label: "Rich Article"},
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

// prevNextForNavigation follows the same source-aware rule as docsNavForNavigation.
func prevNextForNavigation(activePath string, selection navigationSelection) *prevNextView {
	ordered := orderedDocsNav()
	for i, link := range ordered {
		if link.Path != activePath {
			continue
		}
		out := &prevNextView{}
		if i > 0 {
			out.Prev = &prevNextLink{Href: selection.href(ordered[i-1].Path), Label: ordered[i-1].Label}
		}
		if i < len(ordered)-1 {
			out.Next = &prevNextLink{Href: selection.href(ordered[i+1].Path), Label: ordered[i+1].Label}
		}
		return out
	}
	return nil
}

// prevNextForSelection remains available to focused callers that explicitly
// request canonical recipe links.
func prevNextForSelection(activePath string, selection documentSelection, execution accordionExecution) *prevNextView {
	return prevNextForNavigation(activePath, navigationSelection{canonical: true, selection: selection, execution: execution})
}

// handbookSection is one scannable handbook tier in the docs sidebar.
type handbookSection struct {
	Title string
	Links []navLink
}

// handbookSections splits concept docs into Core (how to design screens),
// System (tokens/themes/platform), and Meta (project). Order within and
// across sections is first-to-read; flattening feeds prev/next + sitemap.
var handbookSections = []handbookSection{
	{
		Title: "Core",
		Links: []navLink{
			{Path: "/docs/information-architecture", Label: "Information architecture"},
			{Path: "/docs/screens", Label: "Screens"},
			{Path: "/docs/journeys", Label: "Journeys"},
			{Path: "/docs/data-display", Label: "Data display"},
			{Path: "/docs/feedback", Label: "Feedback"},
			{Path: "/docs/choose-the-right-control", Label: "Choose the right control"},
			{Path: "/docs/forms", Label: "Forms"},
			{Path: "/docs/agent-workflow", Label: "Agent workflow"},
			{Path: "/docs/ui-definition-of-done", Label: "UI definition of done"},
			{Path: "/docs/density", Label: "Density"},
			{Path: "/docs/motion", Label: "Motion"},
			{Path: "/docs/compare", Label: "Why Gelium"},
		},
	},
	{
		Title: "System",
		Links: []navLink{
			{Path: "/docs/themes", Label: "Themes"},
			{Path: "/docs/themes/gallery", Label: "Theme Gallery"},
			{Path: "/docs/tokens", Label: "Tokens"},
			{Path: "/docs/typography", Label: "Typography"},
			{Path: "/docs/spacing", Label: "Spacing"},
			{Path: "/docs/colors", Label: "Colors"},
			{Path: "/docs/server-contracts", Label: "Server contracts"},
			{Path: "/docs/accessibility", Label: "Accessibility"},
			{Path: "/docs/principles", Label: "Design principles"},
			{Path: "/docs/performance", Label: "Performance"},
			{Path: "/docs/responsive", Label: "Responsive"},
			{Path: "/docs/media", Label: "Media"},
			{Path: "/docs/browser-support", Label: "Browser support"},
			{Path: "/docs/content-style", Label: "Content style"},
		},
	},
	{
		Title: "Meta",
		Links: []navLink{
			{Path: "/docs/acknowledgments", Label: "Acknowledgments"},
			{Path: "/docs/contributing", Label: "Contributing"},
			{Path: "/docs/changelog", Label: "Changelog"},
			{Path: "/docs/roadmap", Label: "Roadmap"},
		},
	},
}

// handbookNavLinks flattens handbookSections in sidebar order for prev/next,
// sitemap, and any caller that needs a single ordered concept list.
func handbookNavLinks() []navLink {
	n := 0
	for _, s := range handbookSections {
		n += len(s.Links)
	}
	out := make([]navLink, 0, n)
	for _, s := range handbookSections {
		out = append(out, s.Links...)
	}
	return out
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
			{Path: "/components/banner", Label: "Banner"},
			{Path: "/components/inline-alert", Label: "Inline alert"},
			{Path: "/components/callout", Label: "Callout"},
			{Path: "/components/skeleton", Label: "Skeleton"},
			{Path: "/components/empty-state", Label: "Empty state"},
			{Path: "/components/error-state", Label: "Error state"},
			{Path: "/components/validation-summary", Label: "Validation summary"},
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
	{
		Title: "Composition & content",
		Intro: "Page-level primitives: layout, brand, and editorial content.",
		Links: []navLink{
			{Path: "/components/hero", Label: "Hero"},
			{Path: "/components/avatar", Label: "Avatar"},
			{Path: "/components/breadcrumb", Label: "Breadcrumb"},
			{Path: "/components/footer", Label: "Footer"},
			{Path: "/components/pagination", Label: "Pagination"},
			{Path: "/components/section-heading", Label: "Section heading"},
			{Path: "/components/feature-card", Label: "Feature card"},
			{Path: "/components/split", Label: "Split"},
			{Path: "/components/image", Label: "Image"},
			{Path: "/components/media", Label: "Media"},
			{Path: "/components/video", Label: "Video"},
			{Path: "/components/newsletter", Label: "Newsletter"},
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

// docsIndex is the GET /docs handler. Orientation first (Start here, grouped
// like the sidebar tiers), then recipes, then a deep-dive essay from
// content/index.md. Sidebar owns the full Core/System/Meta + component catalog.
func (s *server) docsIndex(w http.ResponseWriter, r *http.Request) {
	var md string
	md += "# Documentation\n\n"
	md += "This hub **orients** you. The **sidebar** is the map — grouped as **Core** (how to design screens), **System** (tokens, themes, platform), **Meta** (project), then components, Patterns, and Recipes. This page does not repeat that catalog.\n\n"
	if q := strings.TrimSpace(r.URL.Query().Get("q")); q != "" {
		md += searchResultsMarkdown(q)
	}
	md += "## Start here\n\n"
	md += "### Install and agents\n\n"
	md += "- [npm `gelium-ui`](https://www.npmjs.com/package/gelium-ui) — CSS, themes, templates, optional JS for product apps.\n"
	md += "- [Agent brief](/llms.txt) — machine-readable project summary.\n"
	md += "- [Agent UX pack](/llms-ux.txt) — SURFACE / FEED / DATA / JOURNEY / WF tables.\n"
	md += "- [Agent workflow](/docs/agent-workflow) — shape → build → audit → polish (ethos-safe).\n"
	md += "- [UI definition of done](/docs/ui-definition-of-done) — ship checklist before accepting UI.\n"
	md += "- [AEO](/docs/aeo) — answer-first content and discovery aids are not ranking guarantees.\n"
	md += "- [SEO](/docs/seo) — metadata ownership, indexability, sitemap, and social coherence.\n"
	md += "### Core (screen criteria)\n\n"
	md += "- [Screens](/docs/screens) — screen types, hierarchy, nav, surface modes.\n"
	md += "- [Journeys](/docs/journeys) — multi-step flows and post-submit landings.\n"
	md += "- [Data display](/docs/data-display) — table vs list vs cards (`DATA-*`).\n"
	md += "- [Feedback](/docs/feedback) — toast vs summary vs banner vs empty (`FEED-*`).\n"
	md += "- [Forms](/docs/forms) — labels, `inputmode`/`type`, validate after interaction.\n"
	md += "- [Choose the right control](/docs/choose-the-right-control) — radio vs select vs chips…\n"
	md += "- [Why Gelium](/docs/compare) — when to use this stack, and explicit no-gos.\n"
	md += "- [Patterns](/docs/patterns) — domain skeletons (forum, catalog, admin).\n\n"
	md += "### System (design system)\n\n"
	md += "- [Themes](/docs/themes) — class-based direction; dark is a class route.\n"
	md += "- [Tokens](/docs/tokens) — `--ui-*` ownership (core vs theme vs component).\n"
	md += "- [Typography](/docs/typography) — live type roles, measure, weight, and line-height.\n"
	md += "- [Spacing](/docs/spacing) — the `--ui-space-*` scale and compositions.\n"
	md += "- [Colors](/docs/colors) — semantic roles, focus, status, and light/dark examples.\n"
	md += "- [Server contracts](/docs/server-contracts) — GET, POST+303, 422, toast wire.\n"
	md += "- [Accessibility](/docs/accessibility) — platform defaults and contracts.\n"
	md += "- [Performance](/docs/performance) — ~50KB JS stance; CSS is large by design.\n"
	md += "- [Responsive](/docs/responsive) — viewports, not device names.\n"
	md += "- [Media](/docs/media) — accessible images, audio, video, transcripts, and safe embeds.\n"
	md += "- [Density](/docs/density) — comfortable / cozy / compact + shell.\n"
	md += "- [Motion](/docs/motion) — when to animate; reduced-motion.\n"
	md += "- [Content style](/docs/content-style) — voice + H1–H3 / list / table grammar.\n\n"
	md += "## Use the sidebar\n\n"
	md += "Open **Core**, **System**, or **Meta** for the full concept list; component groups (Foundation, Actions, …) for the reference. Mobile: same IA in the docs menu. Come back here for orientation, not for every link.\n\n"
	md += "## Try a full screen\n\n"
	md += "Runnable compositions (Recipes in the sidebar):\n\n"
	md += "- [Admin Resource](/recipes/admin-resource) — data table + form + dialog on the server contract.\n"
	md += "- [Ops Queue](/recipes/ops-queue) — work queue with POST+303 transitions.\n"
	md += "- [Public/Social Feed](/recipes/public-feed) — feed with reactions and loading states.\n"
	md += "- [WhatsApp manager](/demo/whatsapp) — larger composed demo.\n\n"
	if source, err := fs.ReadFile(s.assets, docsRootLeadSource); err == nil {
		if lead := stripDocsRootH1(string(source)); strings.TrimSpace(lead) != "" {
			md += "## Deep dive\n\n"
			md += lead + "\n"
		}
	}
	s.renderMarkdown(w, r, pageView{Title: "Documentation"}, md, "/docs")
}

// searchResultsMarkdown honors the 0-JS GET /docs?q= search fallback
// server-side. It filters the SAME nav model the client-side index is built
// from (docsNavFor -> docsSections + handbookSections), so the no-JS path can
// never agree on a different page set than the JS path. Titles and paths are
// matched case-insensitively; results render as sidebar group + link.
func searchResultsMarkdown(q string) string {
	lower := strings.ToLower(q)
	nav := docsNavFor("/docs", "", "")
	type hit struct {
		label string
		href  string
		group string
	}
	hits := make([]hit, 0, 8)
	seen := make(map[string]bool, 64)
	for _, g := range nav.Groups {
		for _, l := range g.Links {
			if seen[l.Path] {
				continue
			}
			seen[l.Path] = true
			if strings.Contains(strings.ToLower(l.Label), lower) || strings.Contains(strings.ToLower(l.Path), strings.ToLower(q)) {
				hits = append(hits, hit{label: l.Label, href: l.Href, group: g.Title})
			}
		}
	}
	if len(hits) == 0 {
		return "## No matches for “" + q + "”\n\nNothing in the sidebar matches. Try a broader term, or browse the **Core**, **System**, and **Meta** groups; every component, handbook page, and recipe lives in the sidebar.\n\n"
	}
	var b strings.Builder
	b.WriteString("## Search results for “" + q + "”\n\n")
	for _, h := range hits {
		b.WriteString("- [" + h.label + "](" + h.href + ") — " + h.group + ".\n")
	}
	return b.String()
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

// docsAgentWorkflow is GET /docs/agent-workflow — ethos-safe agent passes + surface modes.
func (s *server) docsAgentWorkflow(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Agent workflow"}, "content/handbook-agent-workflow.md", "/docs/agent-workflow")
}

// docsTemplateProduct is GET /docs/templates/product — consumer PRODUCT.md template.
func (s *server) docsTemplateProduct(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "PRODUCT.md template"}, "content/templates/product.md", "/docs/templates/product")
}

// docsTemplateDesign is GET /docs/templates/design — consumer DESIGN.md template.
func (s *server) docsTemplateDesign(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "DESIGN.md template"}, "content/templates/design.md", "/docs/templates/design")
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

// docsTypography, docsSpacing, and docsColors are the human-readable
// foundation showcases. Their specimens stay in Markdown/HTML so the server
// renders the same token-driven markup without a client runtime.
func (s *server) docsTypography(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Typography"}, "content/handbook-typography.md", "/docs/typography")
}

func (s *server) docsSpacing(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Spacing"}, "content/handbook-spacing.md", "/docs/spacing")
}

func (s *server) docsColors(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Colors"}, "content/handbook-colors.md", "/docs/colors")
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

// docsMedia is GET /docs/media — server-rendered media contracts and guidance.
func (s *server) docsMedia(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Media"}, "content/handbook-media.md", "/docs/media")
}

// docsBrowserSupport is GET /docs/browser-support — the Browser support
// handbook page: the consolidated Baseline API table (Popover, anchor
// positioning, Invoker Commands, dialog closedby) and the "what always
// works" contract (native semantics, 0-JS, AA contrast, forced-colors).
func (s *server) docsBrowserSupport(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Browser support"}, "content/handbook-browser-support.md", "/docs/browser-support")
}

func (s *server) docsSEO(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "SEO"}, "content/handbook-seo.md", "/docs/seo")
}

func (s *server) docsAEO(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "AEO"}, "content/handbook-aeo.md", "/docs/aeo")
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
