package app

import (
	"html/template"
	"net/http"
)

// navigationBarDestination is one link in a navigation bar. The destination is
// a real <a href> (the roadmap's "navegación real con links"); Active is a real
// boolean derived server-side from the current page, never from JavaScript.
// InactiveSVG/ActiveSVG carry the trusted internal decorative glyphs; the two
// slots mirror Material's active/inactive icon pair and the active one is
// swapped in with CSS only.
type navigationBarDestination struct {
	Href        string
	Label       string
	Active      bool
	HasBadge    bool
	BadgeDot    bool
	BadgeValue  string
	InactiveSVG template.HTML
	ActiveSVG   template.HTML
}

// navigationBarDemoView is one <nav> instance: its accessible label, optional
// "hide inactive labels" modifier (the upstream hideInactiveLabels option), and
// the ordered destinations. Heading labels the demo section only.
type navigationBarDemoView struct {
	Heading            string
	Label              string
	HideInactiveLabels bool
	Destinations       []navigationBarDestination
}

// navigationBarDemo is the view model for the Navigation bar documentation
// preview. It carries trusted, internal inline SVG glyphs and the demo nav bars
// so the template stays free of raw HTML strings and caller-controlled classes.
type navigationBarDemo struct {
	Bars []navigationBarDemoView
}

// navBarGlyphs: trusted, internal inline SVG glyphs (Material Icons, 24x24).
// Every glyph is decorative: aria-hidden and unfocusable, colored by
// currentColor, and the destination's visible label supplies the accessible
// name. Never place user input in these strings.
const (
	// #nosec G203 -- trusted, internal decorative glyph.
	navBarHomeSVG template.HTML = `<svg class="ui-nav-bar-glyph" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	navBarListSVG template.HTML = `<svg class="ui-nav-bar-glyph" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M3 13h2v-2H3v2zm0 4h2v-2H3v2zm0-8h2V7H3v2zm4 4h14v-2H7v2zm0 4h14v-2H7v2zM7 7v2h14V7H7z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	navBarMenuSVG template.HTML = `<svg class="ui-nav-bar-glyph" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	navBarAppsSVG template.HTML = `<svg class="ui-nav-bar-glyph ui-nav-bar-glyph--active" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M4 8h4V4H4v4zm6 12h4v-4h-4v4zm-6 0h4v-4H4v4zm0-6h4v-4H4v4zm6 0h4v-4h-4v4zm6-10v4h4V4h-4zm-6 4h4V4h-4v4zm6 6h4v-4h-4v4zm0 6h4v-4h-4v4z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	navBarLabelSVG template.HTML = `<svg class="ui-nav-bar-glyph" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M17.63 5.84C17.27 5.33 16.67 5 16 5L5 5.01C3.9 5.01 3 5.9 3 7v10c0 1.1.9 1.99 2 1.99L16 19c.67 0 1.27-.33 1.63-.84L22 12l-4.37-6.16z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	navBarBellSVG template.HTML = `<svg class="ui-nav-bar-glyph" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M12 22c1.1 0 2-.9 2-2h-4c0 1.1.9 2 2 2zm6-6v-5c0-3.07-1.63-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68C7.64 5.36 6 7.92 6 11v5l-2 2v1h16v-1l-2-2z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	navBarMailSVG template.HTML = `<svg class="ui-nav-bar-glyph" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z"></path></svg>`
)

// defaultNavigationBarDestinations is the demo app's shared destination set.
// Every destination is a real link to a Loom component page so clicking it
// navigates normally (no JavaScript); the current page is marked active
// server-side by markActive.
func defaultNavigationBarDestinations() []navigationBarDestination {
	return []navigationBarDestination{
		{Label: "Home", Href: "/components/button", InactiveSVG: navBarHomeSVG, ActiveSVG: navBarHomeSVG},
		{Label: "List", Href: "/components/list", InactiveSVG: navBarListSVG, ActiveSVG: navBarListSVG},
		{Label: "Navigation bar", Href: "/components/navigation-bar", InactiveSVG: navBarMenuSVG, ActiveSVG: navBarAppsSVG},
		{Label: "Chips", Href: "/components/chips", InactiveSVG: navBarLabelSVG, ActiveSVG: navBarLabelSVG},
		{Label: "Badge", Href: "/components/badge", InactiveSVG: navBarBellSVG, ActiveSVG: navBarBellSVG},
	}
}

// markActive sets Active on the destination whose href is the current page, so
// the selected state is always server-derived from the request path.
func markActive(destinations []navigationBarDestination, currentPath string) []navigationBarDestination {
	for i := range destinations {
		destinations[i].Active = destinations[i].Href == currentPath
	}
	return destinations
}

func defaultNavigationBarDemo() navigationBarDemo {
	primary := defaultNavigationBarDestinations()

	badged := defaultNavigationBarDestinations()
	badged[0].HasBadge = true
	badged[0].BadgeDot = true
	badged[3].HasBadge = true
	badged[3].BadgeValue = "3"
	badged[4].HasBadge = true
	badged[4].BadgeValue = "12"

	return navigationBarDemo{
		Bars: []navigationBarDemoView{
			{
				Heading:      "Standard",
				Label:        "Primary",
				Destinations: primary,
			},
			{
				Heading:            "Hide inactive labels",
				Label:              "Primary compact",
				HideInactiveLabels: true,
				Destinations:       primary,
			},
			{
				Heading:      "Badges",
				Label:        "Primary with badges",
				Destinations: badged,
			},
		},
	}
}

func (s *server) navigationBarDocs(w http.ResponseWriter, r *http.Request) {
	demo := defaultNavigationBarDemo()
	for i := range demo.Bars {
		demo.Bars[i].Destinations = markActive(demo.Bars[i].Destinations, r.URL.Path)
	}
	s.renderMarkdownPage(w, pageView{
		Title:             "Navigation bar",
		NavigationBarDemo: &demo,
	}, "content/navigation-bar.md")
}
