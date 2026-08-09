package app

import (
	"net/http"
)

type badgeDemo struct{}

func (s *server) badgeDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:     "Badge",
		BadgeDemo: &badgeDemo{},
	}, "content/badge.md")
}
