package app

import (
	"net/http"
)

type focusRingDemo struct{}

func (s *server) focusRingDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title:         "Focus ring",
		FocusRingDemo: &focusRingDemo{},
	}, "content/focus-ring.md")
}
