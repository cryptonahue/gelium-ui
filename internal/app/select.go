package app

import (
	"bytes"
	"net/http"
	"strings"
)

type selectDemo struct{}

// selectMenuOption is one server-driven menu item in the Select menu demo.
// Selected marks the currently chosen option so the docs page can render
// aria-selected; the HX fragment renders Open=false and omits the marks so a
// closed menu never leaks the "open list" structure into the response.
type selectMenuOption struct {
	Value    string
	Label    string
	Selected bool
}

// selectMenuDemo is the view model for the progressive menu demo: a native
// <dialog> opened by command/commandfor, whose items are submit buttons that
// post the chosen value back to the server. Without command support the menu
// stays closed and the native <select> field remains the fallback control.
type selectMenuDemo struct {
	ID      string
	Value   string
	Label   string
	Options []selectMenuOption
	Error   string
	Open    bool
}

// selectMenuOptions is the closed vocabulary of the server-driven menu demo.
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
		ID:    "select-menu",
		Value: "priority",
		Label: "Priority",
		Open:  true,
		Options: []selectMenuOption{
			{Value: "standard", Label: "Standard", Selected: false},
			{Value: "priority", Label: "Priority", Selected: true},
			{Value: "enterprise", Label: "Enterprise", Selected: false},
		},
	}
}

func (s *server) selectDocs(w http.ResponseWriter, _ *http.Request) {
	demo := defaultSelectMenuDemo()
	s.renderMarkdownPage(w, pageView{
		Title:          "Select",
		SelectDemo:     &selectDemo{},
		SelectMenuDemo: &demo,
	}, "content/select.md")
}

// selectFromClosedSet picks the option whose value matches, returning a demo
// whose selection label follows the chosen option. Unknown values keep the
// previous state so the field never shows a fabricated option.
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
			demo.Value = value
			demo.Label = demo.Options[i].Label
			found = true
		}
	}

	if !found {
		status = http.StatusUnprocessableEntity
		demo.Error = "Select a valid option"
	}

	if !isHX {
		s.renderMarkdownPageStatus(w, pageView{
			Title:          "Select",
			SelectDemo:     &selectDemo{},
			SelectMenuDemo: &demo,
		}, "content/select.md", status)
		return
	}

	// HX request: return only the form fragment so htmx swaps the closed menu
	// field. A 422 carries the validation header plus the visible error; a
	// success returns the updated closed menu without the open-list marks.
	demo.Open = false
	if status == http.StatusUnprocessableEntity {
		w.Header().Set("X-Loom-Validation", "true")
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
