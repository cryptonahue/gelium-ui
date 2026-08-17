package app

import (
	"net/http"
)

// sectionHeadingDocs renders the Section heading component page. Section
// heading is a typographic utility, not a full component (Phase F decision):
// .ui-section-heading always renders an h2 — never h1 — because the page owns
// a single h1 (P2). The utility needs no view model and no JavaScript.
func (s *server) sectionHeadingDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Section heading",
	}, "content/section-heading.md")
}