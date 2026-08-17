package app

import (
	"net/http"
)

// mediaDocs renders the Media component page: audio, transcript, and embed
// figures styled through media.css with native controls, transcripts, and an
// allowlist + consent boundary for third-party embeds.
func (s *server) mediaDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Media",
	}, "content/media.md")
}