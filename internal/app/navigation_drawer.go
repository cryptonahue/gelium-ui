package app

import (
	"html/template"
	"net/http"
)

// navigationDrawerDestination is one link in a navigation drawer. The
// destination is a real <a href> so activation and tab order are native; Active
// is a real boolean derived server-side from the current page (aria-current),
// never from JavaScript. Unlike the navigation bar, the drawer does NOT swap an
// active/inactive icon pair (Material keeps a single glyph per drawer
// destination and recolors it), so there is exactly one Glyph slot. Glyph is a
// trusted, internal decorative SVG. An optional badge reuses the existing
// .ui-badge primitive.
type navigationDrawerDestination struct {
	Href       string
	Label      string
	Active     bool
	HasBadge   bool
	BadgeDot   bool
	BadgeValue string
	Glyph      template.HTML
}

// navigationDrawerDemoView is one drawer instance. Modal selects the roadmap's
// modal variant (a native <dialog> with an invoker-command trigger); otherwise
// the demo renders the permanent standard variant as a <nav> embedded in the
// layout. ID is the dialog id used by the trigger's commandfor; Trigger is the
// modal trigger button. Heading labels the demo section only.
type navigationDrawerDemoView struct {
	Heading      string
	Label        string
	Modal        bool
	ID           string
	Trigger      *buttonView
	Destinations []navigationDrawerDestination
}

// navigationDrawerDemo is the view model for the Navigation drawer
// documentation preview. It carries trusted, internal inline SVG glyphs and the
// demo drawers so the template stays free of raw HTML strings and
// caller-controlled classes.
type navigationDrawerDemo struct {
	Drawers []navigationDrawerDemoView
}

// drawerGlyph constants: trusted, internal inline SVG glyphs (Material Icons,
// 24x24). Every glyph is decorative: aria-hidden and unfocusable, colored by
// currentColor, and the destination's visible label supplies the accessible
// name. Never place user input in these strings.
const (
	// #nosec G203 -- trusted, internal decorative glyph.
	drawerHomeSVG template.HTML = `<svg class="ui-navigation-drawer-glyph-svg" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	drawerInboxSVG template.HTML = `<svg class="ui-navigation-drawer-glyph-svg" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M19 3H4.99c-1.11 0-1.98.9-1.98 2L3 19c0 1.1.88 2 1.99 2H19c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 12h-4c0 1.66-1.35 3-3 3s-3-1.34-3-3H4.99V5H19v10z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	drawerMailSVG template.HTML = `<svg class="ui-navigation-drawer-glyph-svg" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	drawerChatSVG template.HTML = `<svg class="ui-navigation-drawer-glyph-svg" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M20 2H4c-1.1 0-1.99.9-1.99 2L2 22l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zM6 9h12v2H6V9zm8 5H6v-2h8v2zm4-6H6V6h12v2z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	drawerMenuSVG template.HTML = `<svg class="ui-navigation-drawer-glyph-svg" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	drawerBellSVG template.HTML = `<svg class="ui-navigation-drawer-glyph-svg" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M12 22c1.1 0 2-.9 2-2h-4c0 1.1.9 2 2 2zm6-6v-5c0-3.07-1.63-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68C7.64 5.36 6 7.92 6 11v5l-2 2v1h16v-1l-2-2z"></path></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	drawerSettingsSVG template.HTML = `<svg class="ui-navigation-drawer-glyph-svg" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z"></path></svg>`
)

// defaultNavigationDrawerDestinations is the shared destination set. Every
// destination is a real link to a Loom component page so clicking it navigates
// normally (no JavaScript); the current page is marked active server-side by
// markActive.
func defaultNavigationDrawerDestinations() []navigationDrawerDestination {
	return []navigationDrawerDestination{
		{Label: "Home", Href: "/components/button", Glyph: drawerHomeSVG},
		{Label: "Inbox", Href: "/components/text-field", Glyph: drawerInboxSVG},
		{Label: "Mail", Href: "/components/list", Glyph: drawerMailSVG},
		{Label: "Chat", Href: "/components/dialog", Glyph: drawerChatSVG},
		{Label: "Navigation drawer", Href: "/components/navigation-drawer", Glyph: drawerMenuSVG},
		{Label: "Settings", Href: "/components/select", Glyph: drawerSettingsSVG},
	}
}

// markActive sets Active on the destination whose href is the current page, so
// the selected state is always server-derived from the request path.
func markDrawerActive(destinations []navigationDrawerDestination, currentPath string) []navigationDrawerDestination {
	for i := range destinations {
		destinations[i].Active = destinations[i].Href == currentPath
	}
	return destinations
}

func defaultNavigationDrawerDemo() navigationDrawerDemo {
	standard := defaultNavigationDrawerDestinations()

	badged := defaultNavigationDrawerDestinations()
	badged[1].HasBadge = true
	badged[1].BadgeDot = true
	badged[2].HasBadge = true
	badged[2].BadgeValue = "3"

	modal := defaultNavigationDrawerDestinations()

	return navigationDrawerDemo{
		Drawers: []navigationDrawerDemoView{
			{
				Heading:      "Standard",
				Label:        "Primary",
				Destinations: standard,
			},
			{
				Heading:      "Standard with badges",
				Label:        "Primary with badges",
				Destinations: badged,
			},
			{
				Heading: "Modal",
				Label:   "Primary modal",
				Modal:   true,
				ID:      "navigation-drawer-modal",
				Trigger: &buttonView{
					Label:      "Open navigation drawer",
					Variant:    "primary",
					Command:    "show-modal",
					CommandFor: "navigation-drawer-modal",
				},
				Destinations: modal,
			},
		},
	}
}

func (s *server) navigationDrawerDocs(w http.ResponseWriter, r *http.Request) {
	demo := defaultNavigationDrawerDemo()
	for i := range demo.Drawers {
		demo.Drawers[i].Destinations = markDrawerActive(demo.Drawers[i].Destinations, r.URL.Path)
	}
	s.renderMarkdownPage(w, pageView{
		Title:                "Navigation drawer",
		NavigationDrawerDemo: &demo,
	}, "content/navigation-drawer.md")
}
