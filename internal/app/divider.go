package app

import (
	"net/http"
)

type dividerDemo struct{}

func (s *server) dividerDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title:       "Divider",
		DividerDemo: &dividerDemo{},
	}, "content/divider.md")
}
