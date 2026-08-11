package app

import (
	"net/http"
)

type switchDemo struct{}

func (s *server) switchDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title:      "Switch",
		SwitchDemo: &switchDemo{},
	}, "content/switch.md")
}
