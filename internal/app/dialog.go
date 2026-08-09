package app

import (
	"net/http"
)

type dialogView struct {
	Trigger buttonView
	Cancel  buttonView
	Confirm buttonView
}

func (s *server) dialogDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title: "Dialog",
		Dialog: &dialogView{
			Trigger: buttonView{Label: "Open confirmation dialog", Variant: "primary", Command: "show-modal", CommandFor: "confirm-dialog"},
			Cancel:  buttonView{Label: "Cancel", Variant: "text", Autofocus: true, Command: "request-close", CommandFor: "confirm-dialog"},
			Confirm: buttonView{Label: "Confirm", Variant: "text", Command: "close", CommandFor: "confirm-dialog", Value: "confirm"},
		},
	}, "content/dialog.md")
}
