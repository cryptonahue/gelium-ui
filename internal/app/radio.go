package app

import (
	"net/http"
)

type radioDemo struct{}

func (s *server) radioDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:     "Radio",
		RadioDemo: &radioDemo{},
	}, "content/radio.md")
}
