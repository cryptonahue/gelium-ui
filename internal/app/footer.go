package app

import (
	"net/http"
)

func (s *server) footerDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{Title: "Footer"}, "content/footer.md")
}