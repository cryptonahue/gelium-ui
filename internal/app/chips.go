package app

import (
	"html/template"
	"net/http"
)

// removeableChip is one removable input chip in the docs example. Removed is a
// real boolean so the template can render a physically absent control without
// introducing a conditional string or caller-controlled class.
type removeableChip struct {
	Value   string
	Label   string
	Removed bool
}

// chipsDemo is the view model for the Chips documentation preview. It carries
// the trusted, internal inline SVG markers used by the demo scaffolding and the
// removable input chips (each with a server-defined removal state) so the
// template stays free of raw HTML strings and caller-controlled classes.
type chipsDemo struct {
	CalendarSVG  template.HTML
	RemoveSVG    template.HTML
	Removable    []removeableChip
	RemoveNotice string
}

// chipsCalendarSVG is the leading assist-chip icon (decorative; aria-hidden).
const chipsCalendarSVG template.HTML = `<svg class="ui-icon ui-chip-icon" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M19 4h-1V2h-2v2H8V2H6v2H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2Zm0 16H5V9h14v11ZM7 11h2v2H7v-2Zm4 0h2v2h-2v-2Zm4 0h2v2h-2v-2Zm-8 4h2v2H7v-2Zm4 0h2v2h-2v-2Zm4 0h2v2h-2v-2Z"></path></svg>` // #nosec G203 -- trusted, internal decorative glyph.

// chipsRemoveSVG is the trailing remove icon used by the input-chip remove
// action. It is decorative and aria-hidden; the remove button's accessible name
// comes from its aria-label, never from the glyph.
const chipsRemoveSVG template.HTML = `<svg class="ui-icon" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M19 6.41 17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12 19 6.41Z"></path></svg>` // #nosec G203 -- trusted, internal decorative glyph.

func defaultRemovableChips() []removeableChip {
	return []removeableChip{
		{Value: "star-wars", Label: "Star Wars"},
		{Value: "star-trek", Label: "Star Trek"},
	}
}

func defaultChipsDemo() chipsDemo {
	return chipsDemo{
		CalendarSVG: chipsCalendarSVG,
		RemoveSVG:   chipsRemoveSVG,
		Removable:   defaultRemovableChips(),
	}
}

func chipsPage(demo chipsDemo) pageView {
	return pageView{
		Title:     "Chips",
		ChipsDemo: &demo,
	}
}

func (s *server) chipsDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, chipsPage(defaultChipsDemo()), "content/chips.md")
}

// chipsRemoveDemo completes a no-JS server round-trip for the removable input
// chip: the form posts the chip value and the server re-renders the page with
// that chip removed and a persistent inline notice. This is the platform-first
// answer to chip removal; no component JavaScript is required.
func (s *server) chipsRemoveDemo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	demo := defaultChipsDemo()
	removed := r.FormValue("remove")
	if removed != "" {
		known := false
		for i := range demo.Removable {
			if demo.Removable[i].Value == removed {
				demo.Removable[i].Removed = true
				demo.RemoveNotice = "Removed " + demo.Removable[i].Label + "."
				known = true
			}
		}
		if !known {
			demo.RemoveNotice = "Nothing to remove."
		}
	}

	s.renderMarkdownPageStatus(w, chipsPage(demo), "content/chips.md", http.StatusOK)
}
