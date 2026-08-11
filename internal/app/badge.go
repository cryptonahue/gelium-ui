package app

import (
	"net/http"
)

type badgeDemo struct{}

func (s *server) badgeDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title:     "Badge",
		BadgeDemo: &badgeDemo{},
	}, "content/badge.md")
}
