package app

import (
	"net/http"
)

func (s *server) breadcrumbDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{Title: "Breadcrumb"}, "content/breadcrumb.md")
}