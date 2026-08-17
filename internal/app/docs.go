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
	{Path: "/docs/choose-the-right-control", Label: "Choose the right control"},
	{Path: "/docs/compare", Label: "Why Gelium"},
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

// docsIndex is the GET /docs handler. It renders a Markdown page whose content
// is generated from the same componentRoutes registry that drives the nav, so
// the index can never drift from the actual library. The embedded docs root
// (web/content/index.md) leads the page so the hub explains what Gelium UI is
// and how to use it before cataloging the library.
func (s *server) docsIndex(w http.ResponseWriter, r *http.Request) {
	var md string
	if source, err := fs.ReadFile(s.assets, docsRootLeadSource); err == nil {
		if lead := stripDocsRootH1(string(source)); strings.TrimSpace(lead) != "" {
			md += lead + "\n\n"
		}
	}
	md += "# Documentation\n\n"
	md += "The Gelium UI documentation, ordered concept before reference: the Handbook explains how the library works, then the component reference is organized by category. Every page is dogfooded: it renders the real component it documents.\n\n"
	md += "## Handbook\n\n"
	md += "How Gelium UI works: information architecture, themes, tokens, server contracts, accessibility, and the principles behind the library.\n\n"
	for _, link := range handbookNavLinks {
		md += "- [" + link.Label + "](" + link.Path + ")\n"
	}
	md += "\n"
	for _, section := range docsSections {
		md += "## " + section.Title + "\n\n"
		md += section.Intro + "\n\n"
		for _, link := range section.Links {
			md += "- [" + link.Label + "](" + link.Path + ")\n"
		}
		md += "\n"
	}
	md += "\n## Demos\n\n"
	md += "- [WhatsApp manager](/demo/whatsapp) — a complete chat application built with the library.\n"
	md += "\n## Screen recipes\n\n"
	md += "Phase G screen recipes — full screens composed from the library primitives on the canonical server contract.\n\n"
	md += "- [Admin Resource](/recipes/admin-resource) — a server-rendered resource manager (Data table + form + dialog + banner).\n"
	md += "- [Ops Queue](/recipes/ops-queue) — a work queue with avatar, tone badges and POST+303 transitions.\n"
	md += "- [Public/Social Feed](/recipes/public-feed) — a reverse-chronological activity feed with views, reactions and loading states.\n"
	s.renderMarkdown(w, r, pageView{Title: "Documentation"}, md, "/docs")
}

// docsPatterns is GET /docs/patterns — thin stub so the sidebar IA has a real
// href while pattern registry content grows later.
func (s *server) docsPatterns(w http.ResponseWriter, r *http.Request) {
	md := `# Patterns

Composition patterns for Gelium UI (Phase F–G). This page is a stub destination for the docs shell IA.

## Screen recipes

Full screens composed from library primitives live under Recipes:

- [Admin Resource](/recipes/admin-resource)
- [Ops Queue](/recipes/ops-queue)
- [Public/Social Feed](/recipes/public-feed)

## Component patterns

See the [Documentation](/docs) index for foundation, actions, input, feedback, navigation, and data primitives.
`
	s.renderMarkdown(w, r, pageView{Title: "Patterns"}, md, "/docs/patterns")
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

// docsCompare is GET /docs/compare — Why Gelium: honest comparison vs
// React/headless kits (Radix, shadcn, Base UI), payload orders of magnitude,
// when to use and explicit no-gos.
func (s *server) docsCompare(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPageAt(w, r, pageView{Title: "Why Gelium"}, "content/handbook-compare.md", "/docs/compare")
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
