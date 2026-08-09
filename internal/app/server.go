package app

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/yuin/goldmark"

	webassets "loomui/web"
)

type buttonView struct {
	Label      string
	Variant    string
	Href       string
	IconSVG    template.HTML
	Command    string
	CommandFor string
	Value      string
	Disabled   bool
	Loading    bool
	Submit     bool
	Autofocus  bool
}

type textFieldView struct {
	ID          string
	Label       string
	Name        string
	Value       string
	Variant     string
	Helper      string
	MessageRole string
	Error       string
	Disabled    bool
	Textarea    bool
	Autofocus   bool
}

type validationFormView struct {
	Field  textFieldView
	Submit buttonView
}

type dialogView struct {
	Trigger buttonView
	Cancel  buttonView
	Confirm buttonView
}

type toastView struct {
	ID      string
	Type    string // closed vocabulary: info | success | warning | error
	Role    string // status (polite) or alert (assertive for errors)
	Message string
	IconSVG template.HTML
	Dismiss bool
}

type toastTypeOption struct {
	Value string
	Label string
}

type toastDemoView struct {
	Field  textFieldView
	Types  []toastTypeOption
	Type   string
	Submit buttonView
	Toast  *toastView // persistent inline feedback shown when JavaScript is unavailable
}

type elevationDemo struct {
	Levels []elevationDemoLevel
}

type elevationDemoLevel struct {
	Level int
}

type focusRingDemo struct{}

type iconDemo struct {
	DecorativeSVG template.HTML
	NamedSVG      template.HTML
	Existing      []template.HTML
}

type dividerDemo struct{}

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

type badgeDemo struct{}

type checkboxDemo struct{}

type radioDemo struct{}

type switchDemo struct{}

type selectDemo struct{}

const saveIconSVG template.HTML = `<svg aria-hidden="true" focusable="false" class="ui-icon" viewBox="0 0 24 24" width="18" height="18" fill="currentColor"><path d="M17 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V7l-4-4Zm-5 16a3 3 0 1 1 0-6 3 3 0 0 1 0 6Zm3-10H5V5h10v4Z"></path></svg>` // #nosec G203 -- trusted, internal decorative icon markup.

// toastIcons is a closed set of trusted, internal, decorative SVG markers keyed by
// toast type. Each is aria-hidden and unfocusable by contract; the visible Message
// is what assistive technology announces. Never place user input in these strings.
var toastIcons = map[string]template.HTML{ // #nosec G203 -- trusted, internal decorative SVG markup.
	"info":    `<svg class="ui-icon ui-toast-icon ui-toast-icon-info" aria-hidden="true" focusable="false" viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm1 15h-2v-6h2v6Zm0-8h-2V7h2v2Z"></path></svg>`,
	"success": `<svg class="ui-icon ui-toast-icon ui-toast-icon-success" aria-hidden="true" focusable="false" viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm-2 15-5-5 1.41-1.41L10 14.17l6.59-6.58L18 9l-8 8Z"></path></svg>`,
	"warning": `<svg class="ui-icon ui-toast-icon ui-toast-icon-warning" aria-hidden="true" focusable="false" viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M1 21h22L12 2 1 21Zm12-3h-2v-2h2v2Zm0-4h-2v-4h2v4Z"></path></svg>`,
	"error":   `<svg class="ui-icon ui-toast-icon ui-toast-icon-error" aria-hidden="true" focusable="false" viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm1 15h-2v-2h2v2Zm0-4h-2V7h2v6Z"></path></svg>`,
}

// iconDocsSVG set: trusted, internal inline SVG constants for the Icon primitive
// documentation. Every decorative SVG must be aria-hidden and unfocusable; an
// icon that carries meaning must pair with visible text and never rely on an
// aria-label alone. Never place user input in these strings.
const decorativeIconSVG template.HTML = `<svg class="ui-icon" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M19 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2Zm-7 14H6v-2h6v2Zm3-4H6v-2h9v2Zm3-5-1.41 1.41L19.17 8l-1.59 1.58L16.16 8.16 19 5.33l3.41-3.41L20 0l5 5-5.83 5.83Z"></path></svg>` // #nosec G203 -- trusted, internal decorative glyph.

// namedIconSVG is the trusted icon whose meaning is carried by the visible text
// printed beside it in the documentation. The glyph itself is aria-hidden; the
// visible text supplies the accessible name.
const namedIconSVG template.HTML = `<svg class="ui-icon" aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm1 15h-2v-6h2v6Zm0-8h-2V7h2v2Z"></path></svg>` // #nosec G203 -- trusted, internal meaningful icon glyph.

type pageView struct {
	Title          string
	Content        template.HTML
	CTA            *buttonView
	Buttons        []buttonView
	TextFields     []textFieldView
	ValidationForm *validationFormView
	Dialog         *dialogView
	Toasts         []toastView
	ToastDemo      *toastDemoView
	ElevationDemo  *elevationDemo
	FocusRingDemo  *focusRingDemo
	IconDemo       *iconDemo
	DividerDemo    *dividerDemo
	CardDemo       *cardDemo
	BadgeDemo      *badgeDemo
	CheckboxDemo   *checkboxDemo
	RadioDemo      *radioDemo
	SwitchDemo     *switchDemo
	SelectDemo     *selectDemo
}

type server struct {
	templates *template.Template
	markdown  goldmark.Markdown
	assets    fs.FS
}

// New builds the Loom UI documentation HTTP handler from embedded assets.
func New() http.Handler {
	templates := template.Must(template.ParseFS(webassets.Assets, "templates/*.html"))
	s := &server{
		templates: templates,
		markdown:  goldmark.New(),
		assets:    webassets.Assets,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /{$}", s.home)
	mux.HandleFunc("GET /components/button", s.buttonDocs)
	mux.HandleFunc("GET /components/text-field", s.textFieldDocs)
	mux.HandleFunc("GET /components/dialog", s.dialogDocs)
	mux.HandleFunc("GET /components/toast", s.toastDocs)
	mux.HandleFunc("GET /components/elevation", s.elevationDocs)
	mux.HandleFunc("GET /components/focus-ring", s.focusRingDocs)
	mux.HandleFunc("GET /components/icon", s.iconDocs)
	mux.HandleFunc("GET /components/divider", s.dividerDocs)
	mux.HandleFunc("GET /components/card", s.cardDocs)
	mux.HandleFunc("GET /components/badge", s.badgeDocs)
	mux.HandleFunc("GET /components/checkbox", s.checkboxDocs)
	mux.HandleFunc("GET /components/radio", s.radioDocs)
	mux.HandleFunc("GET /components/switch", s.switchDocs)
	mux.HandleFunc("GET /components/select", s.selectDocs)
	mux.HandleFunc("POST /examples/text-field/validate", s.validateTextField)
	mux.HandleFunc("POST /examples/toast/demo", s.toastDemo)
	mux.HandleFunc("GET /static/{name}", s.staticAsset)
	return mux
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

func (s *server) staticAsset(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	contentTypes := map[string]string{
		"app.css":     "text/css; charset=utf-8",
		"htmx.min.js": "text/javascript; charset=utf-8",
		"app.js":      "text/javascript; charset=utf-8",
	}
	contentType, ok := contentTypes[name]
	if !ok {
		http.NotFound(w, r)
		return
	}
	asset, err := fs.ReadFile(s.assets, "static/"+name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(asset)
}

func (s *server) home(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title: "Loom UI",
		CTA: &buttonView{
			Label:   "Explore Button",
			Variant: "primary",
			Href:    "/components/button",
		},
	}, "content/index.md")
}

func (s *server) buttonDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title: "Button",
		Buttons: []buttonView{
			{Label: "Save changes", Variant: "primary", IconSVG: saveIconSVG},
			{Label: "Continue", Variant: "secondary"},
			{Label: "Learn more", Variant: "outline"},
			{Label: "Unavailable", Variant: "primary", Disabled: true},
			{Label: "Save changes", Variant: "primary", Loading: true},
		},
	}, "content/button.md")
}

func (s *server) textFieldDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, textFieldPage(defaultValidationForm()), "content/text-field.md")
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

func textFieldPage(validationForm validationFormView) pageView {
	return pageView{
		Title: "Text field",
		TextFields: []textFieldView{
			{ID: "text-normal", Label: "Name", Variant: "outlined"},
			{ID: "text-helper", Label: "Email", Variant: "filled", Helper: "We'll only use this for account updates."},
			{ID: "text-error", Label: "Username", Variant: "outlined", Value: "?", Error: "Use letters and numbers only."},
			{ID: "text-disabled", Label: "Account ID", Variant: "filled", Value: "ACCT-1042", Disabled: true},
			{ID: "text-disabled-outlined", Label: "Server path", Variant: "outlined", Disabled: true},
			{ID: "text-disabled-textarea", Label: "Changelog", Variant: "filled", Textarea: true, Value: "Locked notes.", Disabled: true},
			{ID: "text-textarea", Label: "Biography", Variant: "outlined", Textarea: true, Helper: "Tell people a little about yourself."},
		},
		ValidationForm: &validationForm,
	}
}

func defaultValidationForm() validationFormView {
	return validationFormView{
		Field:  textFieldView{ID: "validation-name", Label: "Name", Name: "name", Variant: "outlined", Helper: "Enter your name."},
		Submit: buttonView{Label: "Validate name", Variant: "primary", Submit: true},
	}
}

func (s *server) validateTextField(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	value := r.FormValue("name")
	field := textFieldView{ID: "validation-name", Label: "Name", Name: "name", Value: value, Variant: "outlined"}
	isHX := strings.EqualFold(r.Header.Get("HX-Request"), "true")
	status := http.StatusOK
	if strings.TrimSpace(value) == "" {
		field.Error = "Name is required"
		field.Autofocus = !isHX
		status = http.StatusUnprocessableEntity
	} else {
		field.Helper = "Name accepted"
		field.MessageRole = "status"
	}

	data := defaultValidationForm()
	data.Field = field
	if !isHX {
		s.renderMarkdownPageStatus(w, textFieldPage(data), "content/text-field.md", status)
		return
	}

	var rendered bytes.Buffer
	if err := s.templates.ExecuteTemplate(&rendered, "validation-form", data); err != nil {
		http.Error(w, "validation unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if status == http.StatusUnprocessableEntity {
		w.Header().Set("X-Loom-Validation", "true")
	}
	w.WriteHeader(status)
	_, _ = w.Write(rendered.Bytes())
}

func (s *server) renderMarkdownPage(w http.ResponseWriter, data pageView, contentPath string) {
	s.renderMarkdownPageStatus(w, data, contentPath, http.StatusOK)
}

func (s *server) renderMarkdownPageStatus(w http.ResponseWriter, data pageView, contentPath string, status int) {
	source, err := fs.ReadFile(s.assets, contentPath)
	if err != nil {
		http.Error(w, "documentation unavailable", http.StatusInternalServerError)
		return
	}

	var rendered bytes.Buffer
	if err := s.markdown.Convert(source, &rendered); err != nil {
		http.Error(w, "documentation unavailable", http.StatusInternalServerError)
		return
	}

	var page bytes.Buffer
	data.Content = template.HTML(rendered.String()) // #nosec G203 -- source is trusted, embedded Markdown.
	if err := s.templates.ExecuteTemplate(&page, "layout", data); err != nil {
		http.Error(w, "documentation unavailable", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(page.Bytes())
}

// toastTypes is the closed vocabulary of a Toast notification. Unknown labels
// fall back to info rather than introducing arbitrary styling or semantics.
var toastTypes = []string{"info", "success", "warning", "error"}

func sanitizeToastType(t string) string {
	for _, v := range toastTypes {
		if t == v {
			return v
		}
	}
	return "info"
}

func toastRole(t string) string {
	if t == "error" {
		return "alert"
	}
	return "status"
}

func newToast(t, id, message string) toastView {
	if _, ok := toastIcons[t]; !ok {
		t = "info"
	}
	return toastView{
		ID:      id,
		Type:    t,
		Role:    toastRole(t),
		Message: message,
		IconSVG: toastIcons[t],
		Dismiss: true,
	}
}

func toastPage(demo toastDemoView) pageView {
	return pageView{
		Title: "Toast",
		Toasts: []toastView{
			newToast("info", "toast-static-info", "This is an informational notification."),
			newToast("success", "toast-static-success", "Your changes were saved."),
			newToast("warning", "toast-static-warning", "Your session is about to expire."),
			newToast("error", "toast-static-error", "Something went wrong. Try again."),
		},
		ToastDemo: &demo,
	}
}

func defaultToastDemo() toastDemoView {
	return toastDemoView{
		Field: textFieldView{ID: "toast-message", Label: "Message", Name: "message", Value: "Record updated", Variant: "outlined"},
		Types: []toastTypeOption{
			{Value: "info", Label: "Info"},
			{Value: "success", Label: "Success"},
			{Value: "warning", Label: "Warning"},
			{Value: "error", Label: "Error"},
		},
		Type:   "success",
		Submit: buttonView{Label: "Show toast", Variant: "primary", Submit: true},
	}
}

func (s *server) toastDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, toastPage(defaultToastDemo()), "content/toast.md")
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

func (s *server) focusRingDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:         "Focus ring",
		FocusRingDemo: &focusRingDemo{},
	}, "content/focus-ring.md")
}

func (s *server) iconDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title: "Icon",
		IconDemo: &iconDemo{
			DecorativeSVG: decorativeIconSVG,
			NamedSVG:      namedIconSVG,
			Existing: []template.HTML{
				saveIconSVG,
				toastIcons["info"],
				toastIcons["success"],
				toastIcons["warning"],
				toastIcons["error"],
			},
		},
	}, "content/icon.md")
}

func (s *server) dividerDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:       "Divider",
		DividerDemo: &dividerDemo{},
	}, "content/divider.md")
}

func (s *server) cardDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
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

func (s *server) badgeDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:     "Badge",
		BadgeDemo: &badgeDemo{},
	}, "content/badge.md")
}

func (s *server) checkboxDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:        "Checkbox",
		CheckboxDemo: &checkboxDemo{},
	}, "content/checkbox.md")
}

func (s *server) radioDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:     "Radio",
		RadioDemo: &radioDemo{},
	}, "content/radio.md")
}

func (s *server) switchDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:      "Switch",
		SwitchDemo: &switchDemo{},
	}, "content/switch.md")
}

func (s *server) selectDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:      "Select",
		SelectDemo: &selectDemo{},
	}, "content/select.md")
}

// toastTriggerJSON encodes the wire contract of the HX-Trigger response header:
//
//	{"loom:toast":{"type":"success","message":"Saved"}}
//
// The message is JSON-escaped so it can never break out of the header or body.
func toastTriggerJSON(typ, message string) (string, error) {
	payload := struct {
		Toast struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"loom:toast"`
	}{}
	payload.Toast.Type = typ
	payload.Toast.Message = message
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// toastDemo completes a server-driven action. Without JavaScript it re-renders the
// full documentation page with a persistent inline toast (no-JS flow); with HTMX it
// returns only the form fragment and an HX-Trigger that raises the loom:toast event,
// which the local enhancement layer displays as an auto-dismissing toast in the
// aria-live region. Validation failures are never reported as toasts.
func (s *server) toastDemo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	message := r.FormValue("message")
	typ := sanitizeToastType(r.FormValue("type"))
	isHX := strings.EqualFold(r.Header.Get("HX-Request"), "true")

	status := http.StatusOK
	var inline *toastView
	field := textFieldView{ID: "toast-message", Label: "Message", Name: "message", Value: message, Variant: "outlined"}
	if strings.TrimSpace(message) == "" {
		field.Error = "Message is required"
		field.Autofocus = !isHX
		status = http.StatusUnprocessableEntity
	} else {
		field.Helper = "Toast queued"
		field.MessageRole = "status"
		if isHX {
			trigger, err := toastTriggerJSON(typ, message)
			if err != nil {
				http.Error(w, "toast unavailable", http.StatusInternalServerError)
				return
			}
			w.Header().Set("HX-Trigger", trigger)
		} else {
			toast := newToast(typ, "toast-demo-result", message)
			inline = &toast
		}
	}

	demo := defaultToastDemo()
	demo.Field = field
	demo.Type = typ
	demo.Toast = inline

	if !isHX {
		s.renderMarkdownPageStatus(w, toastPage(demo), "content/toast.md", status)
		return
	}

	if status == http.StatusUnprocessableEntity {
		w.Header().Set("X-Loom-Validation", "true")
	}

	var rendered bytes.Buffer
	if err := s.templates.ExecuteTemplate(&rendered, "toast-demo-form", demo); err != nil {
		http.Error(w, "toast unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(rendered.Bytes())
}
