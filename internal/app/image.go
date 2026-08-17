package app

import (
	"net/http"
)

// imageDocs renders the Image component page. Image is a semantic media
// pattern (figures with alt text, intrinsic dimensions, and responsive
// sources) that styles through media.css — the image templates have no
// dedicated styles file.
func (s *server) imageDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Image",
	}, "content/image.md")
}