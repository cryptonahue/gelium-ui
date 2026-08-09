package app

import (
	"net/http"
)

type focusRingDemo struct{}

func (s *server) focusRingDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:         "Focus ring",
		FocusRingDemo: &focusRingDemo{},
	}, "content/focus-ring.md")
}
