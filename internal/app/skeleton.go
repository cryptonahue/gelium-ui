package app

import (
	"net/http"
)

func (s *server) skeletonDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Skeleton",
	}, "content/skeleton.md")
}