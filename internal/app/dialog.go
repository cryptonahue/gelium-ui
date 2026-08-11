package app

import (
	"net/http"
)

// dialogView is the view model for the base page-variant preview: the trigger
// is a real link to the server-rendered confirmation page, so no overlay
// markup and no inert control ever ships.
type dialogView struct {
	Trigger buttonView
}

// dialogConfirmView is the view model for the server-rendered confirmation
// page. Confirm is a real form POST that redirects back; Cancel is a link back.
type dialogConfirmView struct {
	Headline    string
	Description string
	Cancel      buttonView
	Confirm     buttonView
}

func (s *server) dialogDocs(w http.ResponseWriter, r *http.Request) {
	data := pageView{
		Title: "Dialog",
		Dialog: &dialogView{
			Trigger: buttonView{Label: "Open confirmation dialog", Variant: "primary", Href: "/components/dialog/confirm"},
		},
	}
	if r.URL.Query().Get("confirmed") == "1" {
		data.InlineAlert = &inlineAlertView{
			Tone: "success",
			Body: "Action confirmed.",
		}
	}
	s.renderMarkdownPage(w, data, "content/dialog.md")
}

// dialogConfirm renders the page variant of the Dialog: the same headline and
// description the modal variant shows, but inline as normal page content, with
// Confirm as a real form POST and Cancel as a link back. This is the G1
// fallback: it works in every browser with zero component JavaScript.
func (s *server) dialogConfirm(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title: "Dialog",
		DialogConfirm: &dialogConfirmView{
			Headline:    "Confirm action",
			Description: "This action will apply the selected changes. Do you want to continue?",
			Cancel:      buttonView{Label: "Cancel", Variant: "text", Href: "/components/dialog"},
			Confirm:     buttonView{Label: "Confirm", Variant: "text", Submit: true, Value: "confirm"},
		},
	}, "content/dialog.md")
}

// dialogConfirmPost completes the no-JS confirmation round-trip: the form POSTs
// and the server answers with a 303 back to the docs page, which renders the
// result via the persistent inline-alert slot.
func (s *server) dialogConfirmPost(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/components/dialog?confirmed=1", http.StatusSeeOther)
}
