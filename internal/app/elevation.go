package app

import (
	"net/http"
)

type elevationDemo struct {
	Levels []elevationDemoLevel
}

type elevationDemoLevel struct {
	Level int
}

func elevationLevels() []elevationDemoLevel {
	return []elevationDemoLevel{
		{Level: 0}, {Level: 1}, {Level: 2},
		{Level: 3}, {Level: 4}, {Level: 5},
	}
}

func (s *server) elevationDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title: "Elevation",
		ElevationDemo: &elevationDemo{
			Levels: elevationLevels(),
		},
	}, "content/elevation.md")
}
