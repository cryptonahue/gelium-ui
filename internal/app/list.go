package app

import (
	"net/http"
)

type listDemo struct{}

func (s *server) listDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:    "List",
		ListDemo: &listDemo{},
	}, "content/list.md")
}
