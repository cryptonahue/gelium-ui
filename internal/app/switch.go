package app

import (
	"net/http"
)

type switchDemo struct{}

func (s *server) switchDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:      "Switch",
		SwitchDemo: &switchDemo{},
	}, "content/switch.md")
}
