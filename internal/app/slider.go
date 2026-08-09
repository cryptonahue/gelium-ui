package app

import (
	"net/http"
)

type sliderDemo struct{}

func (s *server) sliderDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:      "Slider",
		SliderDemo: &sliderDemo{},
	}, "content/slider.md")
}
