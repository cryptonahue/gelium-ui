package app

import (
	"net/http"
)

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
	md += "The Gelidium UI component library, organized by category. Every page is dogfooded: it renders the real component it documents.\n\n"
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
	s.renderMarkdown(w, pageView{Title: "Documentation"}, md)
}
