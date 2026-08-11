package app

import (
	"net/http"
)

type cardDemoCard struct {
	Title string
	Body  string
	Href  string
}

type cardDemo struct {
	Static cardDemoCard
	Link   cardDemoCard
	Action cardDemoCard
}

func (s *server) cardDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Card",
		CardDemo: &cardDemo{
			Static: cardDemoCard{
				Title: "Elevated",
				Body:  "A static article that groups related content.",
			},
			Link: cardDemoCard{
				Title: "Outlined link",
				Body:  "Navigate with a real anchor underneath.",
				Href:  "/components/elevation",
			},
			Action: cardDemoCard{
				Title: "Filled action",
				Body:  "Act with a real button underneath.",
			},
		},
	}, "content/card.md")
}
