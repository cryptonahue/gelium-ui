package app

import (
	"net/http"
)

type radioDemo struct{}

func (s *server) radioDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title:     "Radio",
		RadioDemo: &radioDemo{},
	}, "content/radio.md")
}
