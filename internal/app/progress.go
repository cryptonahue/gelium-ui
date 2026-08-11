package app

import (
	"net/http"
)

type progressDemo struct{}

func (s *server) progressDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title:        "Progress",
		ProgressDemo: &progressDemo{},
	}, "content/progress.md")
}
