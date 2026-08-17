package app

import (
	"net/http"
)

func (s *server) inlineAlertDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Inline alert",
	}, "content/inline-alert.md")
}