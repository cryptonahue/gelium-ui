package app

import (
	"net/http"
)

type dividerDemo struct{}

func (s *server) dividerDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:       "Divider",
		DividerDemo: &dividerDemo{},
	}, "content/divider.md")
}
