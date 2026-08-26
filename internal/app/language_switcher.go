package app

import (
	"net/http"
)

// languageSwitcherDocs renders the Language switcher component page. The
// switcher is a GET navigation form: changing language navigates with
// ?lang=<code> (the server answers with a 303), the submit stays visible,
// and no script ever auto-submits — so it works with zero JS.
func (s *server) languageSwitcherDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Language switcher",
	}, "content/language-switcher.md")
}
