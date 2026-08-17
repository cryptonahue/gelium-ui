package app

import (
	"net/http"
)

func (s *server) heroDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{Title: "Hero"}, "content/hero.md")
}