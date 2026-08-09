package app

import (
	"html/template"
	"net/http"
)

type iconButtonView struct {
	Label             string
	Variant           string
	IconSVG           template.HTML
	SelectedIcon      template.HTML
	Href              string
	Command           string
	CommandFor        string
	Value             string
	Disabled          bool
	Toggle            bool
	Selected          bool
	AriaLabelSelected string
}

type iconButtonDemo struct {
	Buttons []iconButtonView
}

func (s *server) iconButtonDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title: "Icon button",
		IconButtonDemo: &iconButtonDemo{
			Buttons: []iconButtonView{
				{Label: "Add to favorites", Variant: "standard", IconSVG: iconButtonFavoriteSVG},
				{Label: "Add to favorites", Variant: "filled", IconSVG: iconButtonFavoriteSVG},
				{Label: "Add to favorites", Variant: "filled-tonal", IconSVG: iconButtonFavoriteSVG},
				{Label: "Add to favorites", Variant: "outlined", IconSVG: iconButtonFavoriteSVG},
				{Label: "Add to favorites", Variant: "filled", IconSVG: iconButtonFavoriteSVG, Disabled: true},
				{Label: "Navigate to home", Variant: "standard", IconSVG: iconButtonHomeSVG, Href: "/"},
				{Label: "Add to favorites", Variant: "filled", IconSVG: iconButtonFavoriteSVG, Toggle: true},
				{Label: "Add to favorites", Variant: "filled", IconSVG: iconButtonFavoriteSVG, SelectedIcon: iconButtonCheckSVG, Toggle: true, Selected: true, AriaLabelSelected: "Remove from favorites"},
			},
		},
	}, "content/icon-button.md")
}

// Trusted, internal inline SVG glyphs for the Icon button documentation. Every
// decorative glyph is aria-hidden and unfocusable; the visible text or
// aria-label supplies the accessible name. Never place user input in these
// strings.
const (
	// #nosec G203 -- trusted, internal decorative glyph.
	iconButtonFavoriteSVG template.HTML = `<svg class="ui-icon" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	iconButtonCheckSVG template.HTML = `<svg class="ui-icon" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg>`
	// #nosec G203 -- trusted, internal decorative glyph.
	iconButtonHomeSVG template.HTML = `<svg class="ui-icon" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z"/></svg>`
)
