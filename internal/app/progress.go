package app

import (
	"net/http"
)

type progressDemo struct{}

func (s *server) progressDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:        "Progress",
		ProgressDemo: &progressDemo{},
	}, "content/progress.md")
}
