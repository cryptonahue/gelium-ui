package app

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

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

// toastIcons is a closed set of trusted, internal, decorative SVG markers keyed by
// toast type. Each is aria-hidden and unfocusable by contract; the visible Message
// is what assistive technology announces. Never place user input in these strings.
var toastIcons = map[string]template.HTML{ // #nosec G203 -- trusted, internal decorative SVG markup.
	"info":    `<svg class="ui-icon ui-toast-icon ui-toast-icon-info" aria-hidden="true" focusable="false" viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm1 15h-2v-6h2v6Zm0-8h-2V7h2v2Z"></path></svg>`,
	"success": `<svg class="ui-icon ui-toast-icon ui-toast-icon-success" aria-hidden="true" focusable="false" viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm-2 15-5-5 1.41-1.41L10 14.17l6.59-6.58L18 9l-8 8Z"></path></svg>`,
	"warning": `<svg class="ui-icon ui-toast-icon ui-toast-icon-warning" aria-hidden="true" focusable="false" viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M1 21h22L12 2 1 21Zm12-3h-2v-2h2v2Zm0-4h-2v-4h2v4Z"></path></svg>`,
	"error":   `<svg class="ui-icon ui-toast-icon ui-toast-icon-error" aria-hidden="true" focusable="false" viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm1 15h-2v-2h2v2Zm0-4h-2V7h2v6Z"></path></svg>`,
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

func (s *server) toastDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, toastPage(defaultToastDemo()), "content/toast.md")
}

// toastTriggerJSON encodes the wire contract of the HX-Trigger response header:
//
//	{"gelium:toast":{"type":"success","message":"Saved"}}
//
// The message is JSON-escaped so it can never break out of the header or body.
func toastTriggerJSON(typ, message string) (string, error) {
	payload := struct {
		Toast struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"gelium:toast"`
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
// returns only the form fragment and an HX-Trigger that raises the gelium:toast event,
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
		s.renderMarkdownPageStatus(w, r, toastPage(demo), "content/toast.md", status)
		return
	}

	if status == http.StatusUnprocessableEntity {
		w.Header().Set("X-Gelium-Validation", "true")
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
