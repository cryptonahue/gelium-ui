package app

import (
	"net/http"
)

type listDemo struct{}

func (s *server) listDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title:    "List",
		ListDemo: &listDemo{},
	}, "content/list.md")
}
