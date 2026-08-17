package app

import (
	"net/http"
)

// avatarDocs is the Avatar component page handler. The companion view model
// (avatarView) lives in avatar.go; this file only serves the docs content so
// the recipe primitive and its documentation stay one concern apart.
func (s *server) avatarDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{Title: "Avatar"}, "content/avatar.md")
}