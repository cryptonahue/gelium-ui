package app

import (
	"net/http"
)

// featureCardDocs renders the Feature card component page. Feature card is a
// composition (Card + media + CTA link), not a primitive: the wrapper reuses
// the .ui-card surface, title, body and action slots, so every visual signal
// comes from Card itself.
func (s *server) featureCardDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Feature card",
	}, "content/feature-card.md")
}