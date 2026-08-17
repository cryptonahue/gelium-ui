package app

import (
	"net/http"
)

func (s *server) emptyStateDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Empty state",
	}, "content/empty-state.md")
}