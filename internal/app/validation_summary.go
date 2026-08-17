package app

import (
	"net/http"
)

func (s *server) validationSummaryDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Validation summary",
	}, "content/validation-summary.md")
}