package app

import (
	"net/http"
)

// splitDocs renders the Split component page. Split is a two-column
// editorial composition (media + body) that stacks on narrow screens and
// mirrors automatically in right-to-left documents.
func (s *server) splitDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Split",
	}, "content/split.md")
}