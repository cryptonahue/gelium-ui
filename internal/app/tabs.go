package app

import (
	"html/template"
	"net/http"
)

// Tabs are real navigation links: each tab is an <a href> to its own
// page/section and the selected tab is marked by the handler from the current
// URL (?tab= / ?sub=). There is no role="tablist", no roving focus and no
// component JavaScript: the native link keyboard contract (Tab to move,
// Enter to activate) already satisfies accessibility for this pattern, and the
// roadmap only allows the tablist ARIA pattern when its full keyboard contract
// is genuinely resolved. See the platform-first audit in the delivery report.

type tabsTabView struct {
	Href      string
	Label     string
	Icon      template.HTML
	Active    bool
	Stacked   bool
	AriaLabel string
}

type tabsBarView struct {
	Variant   string
	AriaLabel string
	Heading   string
	Tabs      []tabsTabView
	Panel     template.HTML
}

type tabsDemo struct {
	Bars []tabsBarView
}

// Trusted inline SVG glyphs for the Tabs documentation. Decorative: the
// wrapping span carries aria-hidden and the visible label (or the aria-label
// on icon-only tabs) supplies the accessible name.
const (
	tabImageSVG  template.HTML = `<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" focusable="false"><path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z"></path></svg>`                                                                                                                                                                                                                                                                                                                                // #nosec G203 -- trusted, internal decorative glyph.
	tabVideoSVG  template.HTML = `<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" focusable="false"><path d="M17 10.5V7c0-.55-.45-1-1-1H4c-.55 0-1 .45-1 1v10c0 .55.45 1 1 1h12c.55 0 1-.45 1-1v-3.5l4 4v-11l-4 4z"></path></svg>`                                                                                                                                                                                                                                                                                                                                                   // #nosec G203 -- trusted, internal decorative glyph.
	tabMusicSVG  template.HTML = `<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" focusable="false"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"></path></svg>`                                                                                                                                                                                                                                                                                                                                                                   // #nosec G203 -- trusted, internal decorative glyph.
	tabFlightSVG template.HTML = `<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" focusable="false"><path d="M21 16v-2l-8-5V3.5c0-.83-.67-1.5-1.5-1.5S10 2.67 10 3.5V9l-8 5v2l8-2.5V19l-2 1.5V22l3.5-1 3.5 1v-1.5L13 19v-5.5l8 2.5z"></path></svg>`                                                                                                                                                                                                                                                                                                                                  // #nosec G203 -- trusted, internal decorative glyph.
	tabHotelSVG  template.HTML = `<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" focusable="false"><path d="M7 13c1.66 0 3-1.34 3-3S8.66 7 7 7s-3 1.34-3 3 1.34 3 3 3zm12-6h-8v7H3V5H1v15h2v-3h18v3h2v-9c0-2.21-1.79-4-4-4z"></path></svg>`                                                                                                                                                                                                                                                                                                                                         // #nosec G203 -- trusted, internal decorative glyph.
	tabHikingSVG template.HTML = `<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" focusable="false"><path d="M9.83 8.79c1.02-.62 1.67-1.72 1.67-2.98C11.5 4.29 10.21 3 8.62 3S5.75 4.29 5.75 5.81c0 1.26.65 2.36 1.67 2.98L2.5 20.48h2.4l2.09-5.04h3.26l2.09 5.04h2.4L9.83 8.79zm1.29 3.03h-5l1.79-4.32L8.62 7.5c-.19 0-.38-.04-.57-.11-.66-.25-1.12-.78-1.28-1.42h1.05c.1.17.24.3.42.38.31.14.64.2.98.2.34 0 .67-.07.98-.2.18-.08.32-.21.42-.38h1.05c-.16.64-.62 1.17-1.28 1.42-.19.07-.38.11-.57.11l-.71 2.88zm5.34-3.58l.92 5.04h2.4L18.62 3h-2.1l-1.7 5.24h2.64z"></path></svg>` // #nosec G203 -- trusted, internal decorative glyph.
)

func (s *server) tabsDocs(w http.ResponseWriter, r *http.Request) {
	demo := newTabsDemo(r.URL.Query().Get("tab"), r.URL.Query().Get("sub"))
	s.renderMarkdownPage(w, r, pageView{
		Title:    "Tabs",
		TabsDemo: demo,
	}, "content/tabs.md")
}

// newTabsDemo builds the dogfooded Tabs documentation demo. Selection is
// derived from the current page: the handler validates the query parameters
// against a closed vocabulary and marks the matching tabs server-side.
func newTabsDemo(tab, sub string) *tabsDemo {
	primary := "photos"
	for _, k := range []string{"photos", "videos", "music"} {
		if tab == k {
			primary = k
		}
	}
	secondary := "travel"
	for _, k := range []string{"travel", "hotel", "activities"} {
		if sub == k {
			secondary = k
		}
	}

	icon := func(k string) template.HTML {
		switch k {
		case "photos":
			return tabImageSVG
		case "videos":
			return tabVideoSVG
		default:
			return tabMusicSVG
		}
	}
	subIcon := func(k string) template.HTML {
		switch k {
		case "travel":
			return tabFlightSVG
		case "hotel":
			return tabHotelSVG
		default:
			return tabHikingSVG
		}
	}

	primaryTab := func(k string) tabsTabView {
		return tabsTabView{
			Href:    "/components/tabs?tab=" + k,
			Label:   capitalize(k),
			Icon:    icon(k),
			Active:  primary == k,
			Stacked: true,
		}
	}
	secondaryTab := func(k string) tabsTabView {
		return tabsTabView{
			Href:   "/components/tabs?sub=" + k,
			Label:  capitalize(k),
			Icon:   subIcon(k),
			Active: secondary == k,
		}
	}

	primaryKeys := []string{"photos", "videos", "music"}
	secondaryKeys := []string{"travel", "hotel", "activities"}

	return &tabsDemo{Bars: []tabsBarView{
		{
			Variant:   "primary",
			AriaLabel: "Primary tabs example",
			Heading:   "Primary · icon + label",
			Tabs:      mapKeys(primaryKeys, primaryTab),
			Panel:     template.HTML("<p>Server-rendered panel for the <strong>" + capitalize(primary) + "</strong> tab. The selected tab is marked by the handler from the current URL — no JavaScript.</p>"), // #nosec G203 -- trusted, internal static markup.
		},
		{
			Variant:   "primary",
			AriaLabel: "Primary tabs label only",
			Heading:   "Primary · label only",
			Tabs:      mapKeys(primaryKeys, labelOnlyTab(primaryTab)),
		},
		{
			Variant:   "primary",
			AriaLabel: "Primary tabs icon only",
			Heading:   "Primary · icon only",
			Tabs:      mapKeys(primaryKeys, iconOnlyTab(primaryTab)),
		},
		{
			Variant:   "secondary",
			AriaLabel: "Secondary tabs example",
			Heading:   "Secondary · icon + label",
			Tabs:      mapKeys(secondaryKeys, secondaryTab),
			Panel:     template.HTML("<p>Server-rendered panel for the <strong>" + capitalize(secondary) + "</strong> tab. The selected tab is marked by the handler from the current URL — no JavaScript.</p>"), // #nosec G203 -- trusted, internal static markup.
		},
		{
			Variant:   "secondary",
			AriaLabel: "Secondary tabs label only",
			Heading:   "Secondary · label only",
			Tabs:      mapKeys(secondaryKeys, labelOnlyTab(secondaryTab)),
		},
		{
			Variant:   "secondary",
			AriaLabel: "Secondary tabs icon only",
			Heading:   "Secondary · icon only",
			Tabs:      mapKeys(secondaryKeys, iconOnlyTab(secondaryTab)),
		},
	}}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	b[0] = b[0] - 'a' + 'A'
	return string(b)
}

func mapKeys(keys []string, build func(string) tabsTabView) []tabsTabView {
	out := make([]tabsTabView, 0, len(keys))
	for _, k := range keys {
		out = append(out, build(k))
	}
	return out
}

// labelOnlyTab strips the icon so the tab shows text only.
func labelOnlyTab(build func(string) tabsTabView) func(string) tabsTabView {
	return func(k string) tabsTabView {
		t := build(k)
		t.Icon = ""
		t.Stacked = false
		return t
	}
}

// iconOnlyTab keeps the icon and moves the accessible name to an aria-label,
// matching the upstream guidance for icon-only tabs.
func iconOnlyTab(build func(string) tabsTabView) func(string) tabsTabView {
	return func(k string) tabsTabView {
		t := build(k)
		t.Label = ""
		t.AriaLabel = capitalize(k)
		t.Stacked = false
		return t
	}
}
