package app

import (
	"html/template"
	"net/http"
)

type iconDemo struct {
	DecorativeSVG template.HTML
	NamedSVG      template.HTML
	Existing      []template.HTML
}

const saveIconSVG template.HTML = `<svg aria-hidden="true" focusable="false" class="ui-icon" viewBox="0 0 24 24" width="18" height="18" fill="currentColor"><path d="M17 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V7l-4-4Zm-5 16a3 3 0 1 1 0-6 3 3 0 0 1 0 6Zm3-10H5V5h10v4Z"></path></svg>` // #nosec G203 -- trusted, internal decorative icon markup.

// iconDocsSVG set: trusted, internal inline SVG constants for the Icon primitive
// documentation. Every decorative SVG must be aria-hidden and unfocusable; an
// icon that carries meaning must pair with visible text and never rely on an
// aria-label alone. Never place user input in these strings.
const decorativeIconSVG template.HTML = `<svg class="ui-icon" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M19 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2Zm-7 14H6v-2h6v2Zm3-4H6v-2h9v2Zm3-5-1.41 1.41L19.17 8l-1.59 1.58L16.16 8.16 19 5.33l3.41-3.41L20 0l5 5-5.83 5.83Z"></path></svg>` // #nosec G203 -- trusted, internal decorative glyph.

// namedIconSVG is the trusted icon whose meaning is carried by the visible text
// printed beside it in the documentation. The glyph itself is aria-hidden; the
// visible text supplies the accessible name.
const namedIconSVG template.HTML = `<svg class="ui-icon" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm1 15h-2v-6h2v6Zm0-8h-2V7h2v2Z"></path></svg>` // #nosec G203 -- trusted, internal meaningful icon glyph.

func (s *server) iconDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title: "Icon",
		IconDemo: &iconDemo{
			DecorativeSVG: decorativeIconSVG,
			NamedSVG:      namedIconSVG,
			Existing: []template.HTML{
				saveIconSVG,
				toastIcons["info"],
				toastIcons["success"],
				toastIcons["warning"],
				toastIcons["error"],
			},
		},
	}, "content/icon.md")
}
