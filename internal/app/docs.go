package app

import (
	"net/http"
	"net/url"
	"strings"
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
// Version aligns with the static asset query (?v=0.4.0). SearchDisabled is
// always true in this change — placeholder only, no live corpus search.
type docsNavView struct {
	Groups         []docsNavGroup
	Version        string
	SearchDisabled bool
}

// docsShellVersion is the static version badge shown in the docs topbar.
const docsShellVersion = "0.4.0"

// usesDocsShell reports whether path should render the two-pane docs chrome.
// Home, recipes, and demos stay on the legacy site-header layout.
func usesDocsShell(path string) bool {
	return path == "/docs" ||
		strings.HasPrefix(path, "/docs/") ||
		strings.HasPrefix(path, "/components/")
}

// docsNavHref builds a sidebar href. When themeSlug is a catalog slug, the
// link carries only ?theme=<slug> so chrome navigation preserves visual
// direction without re-emitting unrelated query state.
func docsNavHref(path, themeSlug string) string {
	if themeSlug == "" {
		return path
	}
	if _, ok := themeBySlugOrClass(themeSlug); !ok {
		return path
	}
	// Prefer the canonical short slug from the catalog (not theme-*).
	slug := themeSlug
	if t, ok := themeBySlugOrClass(themeSlug); ok {
		slug = t.Slug
	}
	q := url.Values{}
	q.Set("theme", slug)
	return path + "?" + q.Encode()
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
// themeSlug, when non-empty and allowlisted, is appended as ?theme= on every
// sidebar Href so in-shell navigation keeps the selected direction.
func docsNavFor(activePath, themeSlug string) docsNavView {
	link := func(path, label string) docsNavLink {
		return docsNavLink{
			Path:    path,
			Href:    docsNavHref(path, themeSlug),
			Label:   label,
			Current: activePath != "" && path == activePath,
		}
	}

	groups := make([]docsNavGroup, 0, 4+len(docsSections))
	groups = append(groups, docsNavGroup{
		Title: "Getting started",
		Links: []docsNavLink{link("/docs", "Documentation")},
	})
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
		docsNavGroup{
			Title: "Themes",
			Links: []docsNavLink{link("/docs/themes", "Themes")},
		},
	)

	return docsNavView{
		Groups:         groups,
		Version:        docsShellVersion,
		SearchDisabled: true,
	}
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

// docsIndex is the GET /docs handler. It renders a Markdown page whose content
// is generated from the same componentRoutes registry that drives the nav, so
// the index can never drift from the actual library.
func (s *server) docsIndex(w http.ResponseWriter, r *http.Request) {
	var md string
	md += "# Documentation\n\n"
	md += "The Gelium UI component library, organized by category. Every page is dogfooded: it renders the real component it documents.\n\n"
	for _, section := range docsSections {
		md += "## " + section.Title + "\n\n"
		md += section.Intro + "\n\n"
		for _, link := range section.Links {
			md += "- [" + link.Label + "](" + link.Path + ")\n"
		}
		md += "\n"
	}
	md += "## Demos\n\n"
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

// docsThemes is GET /docs/themes — thin stub explaining visual directions and
// the 0-JS ?theme= switcher hosted in the docs topbar.
func (s *server) docsThemes(w http.ResponseWriter, r *http.Request) {
	md := `# Themes

Gelium UI ships multiple visual directions on one markup surface. Selection is document-root and zero-JS: append ` + "`?theme=<slug>`" + ` to any docs or component URL, or use the Theme control in the docs topbar.

## Directions

| Direction | Query |
|-----------|-------|
| Material (default) | ` + "`?theme=material`" + ` |
| Basecoat | ` + "`?theme=basecoat`" + ` |

Only allowlisted slugs apply. Unknown values keep the default direction.

## What stays the same

Theme switching must not change URLs, landmarks, or SEO metadata — only the root class and cascade tokens.
`
	s.renderMarkdown(w, r, pageView{Title: "Themes"}, md, "/docs/themes")
}
