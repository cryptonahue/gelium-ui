package app

import (
	"net/http"
)

type checkboxDemo struct{}

func (s *server) checkboxDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title:        "Checkbox",
		CheckboxDemo: &checkboxDemo{},
	}, "content/checkbox.md")
}
