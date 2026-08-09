package app

import (
	"net/http"
)

type checkboxDemo struct{}

func (s *server) checkboxDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:        "Checkbox",
		CheckboxDemo: &checkboxDemo{},
	}, "content/checkbox.md")
}
