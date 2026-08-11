package app

import (
	"net/http"
)

type sliderDemo struct{}

func (s *server) sliderDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title:      "Slider",
		SliderDemo: &sliderDemo{},
	}, "content/slider.md")
}
