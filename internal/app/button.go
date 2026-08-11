package app

import (
	"html/template"
	"net/http"
)

type buttonView struct {
	Label      string
	Variant    string
	Href       string
	IconSVG    template.HTML
	Command    string
	CommandFor string
	Value      string
	Disabled   bool
	Loading    bool
	Submit     bool
	Autofocus  bool
}

func (s *server) buttonDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Button",
		Buttons: []buttonView{
			{Label: "Save changes", Variant: "primary", IconSVG: saveIconSVG},
			{Label: "Continue", Variant: "secondary"},
			{Label: "Learn more", Variant: "outline"},
			{Label: "Unavailable", Variant: "primary", Disabled: true},
			{Label: "Save changes", Variant: "primary", Loading: true},
		},
	}, "content/button.md")
}
