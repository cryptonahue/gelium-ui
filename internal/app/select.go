package app

import (
	"bytes"
	"net/http"
	"strings"
)

type selectDemo struct{}

// selectMenuOption is one option of the server-driven Select menu demo.
// Selected marks the currently chosen option so the docs page can render the
// native <select>'s selected option.
type selectMenuOption struct {
	Value    string
	Label    string
	Selected bool
}

// selectMenuDemo is the view model for the server-driven Select menu demo: a
// real native <select> whose value is posted back to the server. The control is
// the Select component's own field surface, so the form works in every browser
// with zero component JavaScript and no Invoker Commands dependency.
type selectMenuDemo struct {
	Options []selectMenuOption
	Error   string
}

// selectMenuOptions is the closed vocabulary of the server-driven demo.
// Unknown values must be rejected with a 422 so the docs page can teach the
// validation contract without inventing arbitrary options.
var selectMenuOptions = []selectMenuOption{
	{Value: "standard", Label: "Standard"},
	{Value: "priority", Label: "Priority"},
	{Value: "enterprise", Label: "Enterprise"},
}

// defaultSelectMenuDemo returns the demo in its initial state: Priority is the
// server-side default selection, matching the native <select> demo where the
// second field is populated.
func defaultSelectMenuDemo() selectMenuDemo {
	return selectMenuDemo{
		Options: []selectMenuOption{
			{Value: "standard", Label: "Standard", Selected: false},
			{Value: "priority", Label: "Priority", Selected: true},
			{Value: "enterprise", Label: "Enterprise", Selected: false},
		},
	}
}

func (s *server) selectDocs(w http.ResponseWriter, r *http.Request) {
	demo := defaultSelectMenuDemo()
	s.renderMarkdownPage(w, r, pageView{
		Title:          "Select",
		SelectDemo:     &selectDemo{},
		SelectMenuDemo: &demo,
	}, "content/select.md")
}

// selectMenu completes the no-JS server round-trip for the Select menu demo:
// the native <select> posts its value and the server validates it against the
// closed vocabulary, then re-renders. Without JavaScript the whole
// documentation page is re-rendered; with HTMX only the form fragment is
// swapped. Unknown values are rejected with a 422 and carry the validation
// header plus a visible inline error.
func (s *server) selectMenu(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	isHX := strings.EqualFold(r.Header.Get("HX-Request"), "true")
	demo := defaultSelectMenuDemo()

	value := r.FormValue("value")
	status := http.StatusOK
	found := false
	for i := range demo.Options {
		demo.Options[i].Selected = false
		if demo.Options[i].Value == value {
			demo.Options[i].Selected = true
			found = true
		}
	}

	if !found {
		status = http.StatusUnprocessableEntity
		demo.Error = "Select a valid option"
	}

	if !isHX {
		s.renderMarkdownPageStatus(w, r, pageView{
			Title:          "Select",
			SelectDemo:     &selectDemo{},
			SelectMenuDemo: &demo,
		}, "content/select.md", status)
		return
	}

	// HX request: return only the form fragment so htmx swaps the updated
	// native select field. A 422 carries the validation header plus the
	// visible error; a success returns the field with the new selection.
	if status == http.StatusUnprocessableEntity {
		w.Header().Set("X-Gelium-Validation", "true")
	}

	var rendered bytes.Buffer
	if err := s.templates.ExecuteTemplate(&rendered, "select-menu-demo", demo); err != nil {
		http.Error(w, "select menu unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(rendered.Bytes())
}
