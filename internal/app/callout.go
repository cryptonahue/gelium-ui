package app

import (
	"net/http"
)

func (s *server) calloutDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Callout",
	}, "content/callout.md")
}