package app

import (
	"net/http"
)

// videoDocs renders the Video component page. Video is a responsive container
// (not a content component): .ui-video wraps a native <video> with controls,
// lazy loading, captions, and a fallback — zero component JavaScript.
func (s *server) videoDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Video",
	}, "content/video.md")
}