package app

import (
	"html/template"
	"net/http"
)

// navigationTabDestination is one Navigation tab: the individual destination of
// a navigation bar, reimplemented as a real <a href> link. The roadmap's
// "link semántico, no tab falso" contract applies: there is no role="tab", no
// role="tablist" and no roving focus, so the active tab is a real boolean
// derived server-side from the current page (aria-current="page"), never from
// JavaScript. Active is derived from the current page; the tab composes the
// same destination contract as .ui-nav-bar so it drops into the delivered
// navigation bar. InactiveSVG/ActiveSVG carry the trusted internal decorative
// glyphs; when both slots carry the same glyph, the template renders a single
// glyph and skips the redundant second copy (there is nothing to swap).
type navigationTabDestination struct {
	Href        string
	Label       string
	Active      bool
	HasBadge    bool
	BadgeDot    bool
	BadgeValue  string
	InactiveSVG template.HTML
	ActiveSVG   template.HTML
}

// HasDistinctActiveGlyph reports whether the tab swaps between two different
// glyphs. When inactive and active are the same trusted SVG, the template
// renders one copy only — the CSS swap would otherwise show the icon twice.
func (t navigationTabDestination) HasDistinctActiveGlyph() bool {
	return t.InactiveSVG != t.ActiveSVG
}

// navigationTabRow is one navigation-tab demo row: a standalone strip of tabs
// or an in-bar composition. InBar wraps the destinations in the delivered
// .ui-nav-bar contract so the demo proves the tab composes into the bar;
// HideInactiveLabels mirrors the upstream per-tab hideInactiveLabel option.
type navigationTabRow struct {
	Heading            string
	AriaLabel          string
	InBar              bool
	HideInactiveLabels bool
	Destinations       []navigationTabDestination
}

// navigationTabDemo is the view model for the Navigation tab documentation
// preview. It carries trusted, internal inline SVG glyphs and the demo rows so
// the template stays free of raw HTML strings and caller-controlled classes.
type navigationTabDemo struct {
	Rows []navigationTabRow
}

// navTabGlyphs: trusted, internal inline SVG glyphs (Material Icons, 24x24).
// Every glyph is decorative: aria-hidden and unfocusable, colored by
// currentColor, and the destination's visible label supplies the accessible
// name. Never place user input in these strings.
const (
	// #nosec G203 -- trusted, internal decorative glyph.
	navTabHomeSVG template.HTML = `<svg class="ui-nav-tab-glyph" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	navTabListSVG template.HTML = `<svg class="ui-nav-tab-glyph" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M3 13h2v-2H3v2zm0 4h2v-2H3v2zm0-8h2V7H3v2zm4 4h14v-2H7v2zm0 4h14v-2H7v2zM7 7v2h14V7H7z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	navTabMenuSVG template.HTML = `<svg class="ui-nav-tab-glyph" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	navTabAppsSVG template.HTML = `<svg class="ui-nav-tab-glyph ui-nav-tab-glyph--active" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M4 8h4V4H4v4zm6 12h4v-4h-4v4zm-6 0h4v-4H4v4zm0-6h4v-4H4v4zm6 0h4v-4h-4v4zm6-10v4h4V4h-4zm-6 4h4V4h-4v4zm6 6h4v-4h-4v4zm0 6h4v-4h-4v4z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	navTabLabelSVG template.HTML = `<svg class="ui-nav-tab-glyph" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M17.63 5.84C17.27 5.33 16.67 5 16 5L5 5.01C3.9 5.01 3 5.9 3 7v10c0 1.1.9 1.99 2 1.99L16 19c.67 0 1.27-.33 1.63-.84L22 12l-4.37-6.16z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	navTabBellSVG template.HTML = `<svg class="ui-nav-tab-glyph" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M12 22c1.1 0 2-.9 2-2h-4c0 1.1.9 2 2 2zm6-6v-5c0-3.07-1.63-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68C7.64 5.36 6 7.92 6 11v5l-2 2v1h16v-1l-2-2z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	navTabMailSVG template.HTML = `<svg class="ui-nav-tab-glyph" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z"></path></svg>`
)

// defaultNavigationTabDestinations is the demo app's shared destination set.
// Every destination is a real link to a Loom component page so clicking it
// navigates normally (no JavaScript); the current page is marked active
// server-side by markNavigationTabActive.
func defaultNavigationTabDestinations() []navigationTabDestination {
	return []navigationTabDestination{
		{Label: "Home", Href: "/components/button", InactiveSVG: navTabHomeSVG, ActiveSVG: navTabHomeSVG},
		{Label: "List", Href: "/components/list", InactiveSVG: navTabListSVG, ActiveSVG: navTabListSVG},
		{Label: "Navigation tab", Href: "/components/navigation-tab", InactiveSVG: navTabMenuSVG, ActiveSVG: navTabAppsSVG},
		{Label: "Badge", Href: "/components/badge", InactiveSVG: navTabBellSVG, ActiveSVG: navTabBellSVG},
	}
}

// markNavigationTabActive sets Active on the tab whose href is the current
// page, so the selected state is always server-derived from the request path.
func markNavigationTabActive(destinations []navigationTabDestination, currentPath string) []navigationTabDestination {
	for i := range destinations {
		destinations[i].Active = destinations[i].Href == currentPath
	}
	return destinations
}

func defaultNavigationTabDemo() navigationTabDemo {
	standard := defaultNavigationTabDestinations()

	badged := defaultNavigationTabDestinations()
	badged[0].HasBadge = true
	badged[0].BadgeDot = true
	badged[1].HasBadge = true
	badged[1].BadgeValue = "3"
	badged[3].HasBadge = true
	badged[3].BadgeValue = "12"

	return navigationTabDemo{
		Rows: []navigationTabRow{
			{
				Heading:      "Standard",
				AriaLabel:    "Navigation tab examples",
				Destinations: standard,
			},
			{
				Heading:            "Hide inactive labels",
				AriaLabel:          "Navigation tab compact",
				HideInactiveLabels: true,
				Destinations:       standard,
			},
			{
				Heading:      "Badges",
				AriaLabel:    "Navigation tabs with badges",
				Destinations: badged,
			},
			{
				Heading:      "Inside the navigation bar",
				AriaLabel:    "Navigation bar composition",
				InBar:        true,
				Destinations: standard,
			},
		},
	}
}

func (s *server) navigationTabDocs(w http.ResponseWriter, r *http.Request) {
	demo := defaultNavigationTabDemo()
	for i := range demo.Rows {
		demo.Rows[i].Destinations = markNavigationTabActive(demo.Rows[i].Destinations, r.URL.Path)
	}
	s.renderMarkdownPage(w, pageView{
		Title:             "Navigation tab",
		NavigationTabDemo: &demo,
	}, "content/navigation-tab.md")
}
