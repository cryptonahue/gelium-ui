package app

import (
	"net/http"
)

func (s *server) errorStateDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Error state",
	}, "content/error-state.md")
}