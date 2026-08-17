package app

import (
	"net/http"
)

func (s *server) bannerDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Banner",
	}, "content/banner.md")
}