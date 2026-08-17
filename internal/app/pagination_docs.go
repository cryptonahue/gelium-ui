package app

import (
	"net/http"
)

// paginationDocs is the Pagination component page handler. The companion view
// model (paginationView) lives in pagination.go; this file only serves the
// docs content so the reusable partial and its documentation stay one concern
// apart.
func (s *server) paginationDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{Title: "Pagination"}, "content/pagination.md")
}