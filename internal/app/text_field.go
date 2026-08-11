package app

import (
	"bytes"
	"net/http"
	"strings"
)

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

func (s *server) textFieldDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, textFieldPage(defaultValidationForm()), "content/text-field.md")
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
		s.renderMarkdownPageStatus(w, r, textFieldPage(data), "content/text-field.md", status)
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
