package app

import (
	"bytes"
	"net/http"
	"net/url"
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

// selectPopupExperiment is a deliberately small, inline disclosure experiment
// shown only for the optional HTMX execution profile. It is not a custom
// listbox: native details/summary owns disclosure while native submit buttons
// carry the closed option vocabulary to the server.
type selectPopupExperiment struct {
	Options        []selectMenuOption
	SelectedLabel  string
	Error          string
	SubmittedValue string
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

func defaultSelectPopupExperiment() selectPopupExperiment {
	demo, _ := selectPopupExperimentForValue("priority")
	return demo
}

// selectPopupExperimentForValue resolves only the documented option vocabulary.
// An unknown submitted value leaves the safe default selected while preserving
// the raw value for the server-rendered validation message.
func selectPopupExperimentForValue(value string) (selectPopupExperiment, bool) {
	selected := "priority"
	found := false
	for _, option := range selectMenuOptions {
		if option.Value == value {
			selected = value
			found = true
			break
		}
	}

	demo := selectPopupExperiment{Options: make([]selectMenuOption, len(selectMenuOptions))}
	for i, option := range selectMenuOptions {
		option.Selected = option.Value == selected
		if option.Selected {
			demo.SelectedLabel = option.Label
		}
		demo.Options[i] = option
	}
	return demo, found
}

// selectExamples returns the pilot Examples for the Select page with the
// server-driven menu example bound to the current demo, so the example's
// live form always reflects the latest selection. The Select partials
// hardcode their element ids, so the Examples section OWNS the demos on the
// pilot page: the standalone SelectDemo/SelectMenuDemo preview sections stay
// nil to avoid duplicate ids on one page.
func selectExamples(demo selectMenuDemo) ([]exampleBlock, *apiRefView) {
	examples, apiRef := examplesFor("select")
	out := make([]exampleBlock, len(examples))
	copy(out, examples)
	for i := range out {
		if out[i].Partial == "select-menu-demo" {
			out[i].Views = []any{demo}
		}
	}
	return out, apiRef
}

func (s *server) selectDocs(w http.ResponseWriter, r *http.Request) {
	demo := defaultSelectMenuDemo()
	examples, apiRef := selectExamples(demo)
	page := pageView{
		Title:    "Select",
		Examples: examples,
		APIRef:   apiRef,
	}
	if accordionExecutionFromRequest(r) == accordionExecutionHTMX {
		popup := defaultSelectPopupExperiment()
		if value := r.URL.Query().Get("select_popup_value"); value != "" {
			if resolved, ok := selectPopupExperimentForValue(value); ok {
				popup = resolved
			}
		}
		page.SelectPopup = &popup
	}
	s.renderMarkdownPage(w, r, page, "content/select.md")
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
		examples, apiRef := selectExamples(demo)
		s.renderMarkdownPageStatus(w, r, pageView{
			Title:    "Select",
			Examples: examples,
			APIRef:   apiRef,
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

// selectPopup handles the optional inline disclosure experiment. It is a real
// progressive enhancement boundary: regular form posts redirect to a canonical
// GET URL, while HTMX posts receive only the replacement section. HTMX 4 swaps
// 422 responses by default, so the validation response is a complete section.
func (s *server) selectPopup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Vary", "HX-Request")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	// HTMX 4 sends the literal request marker `HX-Request: true`. Treat every
	// other value as a regular form post so Vary: HX-Request cleanly separates
	// the full-page and fragment representations.
	isHX := r.Header.Get("HX-Request") == "true"
	value := r.FormValue("popup_value")
	demo, found := selectPopupExperimentForValue(value)
	status := http.StatusOK
	if !found {
		status = http.StatusUnprocessableEntity
		demo.Error = "Select a valid plan"
		demo.SubmittedValue = value
	}

	if !isHX {
		if found {
			state := url.Values{}
			state.Set("execution", string(accordionExecutionHTMX))
			state.Set("select_popup_value", value)
			http.Redirect(w, r, "/components/select?"+state.Encode()+"#select-popup-experiment", http.StatusSeeOther)
			return
		}
		examples, apiRef := selectExamples(defaultSelectMenuDemo())
		s.renderMarkdownPageStatus(w, r, pageView{
			Title:       "Select",
			Examples:    examples,
			APIRef:      apiRef,
			SelectPopup: &demo,
		}, "content/select.md", status)
		return
	}

	// This header remains compatible with the site's shared validation hook;
	// HTMX 4 itself swaps this full section on HTTP 422 by default.
	if status == http.StatusUnprocessableEntity {
		w.Header().Set("X-Gelium-Validation", "true")
	}
	var rendered bytes.Buffer
	if err := s.templates.ExecuteTemplate(&rendered, "select-popup-experiment", demo); err != nil {
		http.Error(w, "select popup experiment unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(rendered.Bytes())
}
