package app

import (
	"bytes"
	"net/http"
	"strings"
)

// newsletterView is the server-driven view model for the NEWSLETTER public
// content pattern (Phase F): a zero-JS subscription form. The server owns the
// whole contract — POST + 422 (X-Gelium-Validation) with an inline alert on an
// invalid email, persistent success view on a valid one — so the aside never
// mutates state client-side. Action is the POST target; Error and Success are
// mutually exclusive, and Success replaces the form entirely.
type newsletterView struct {
	ID          string
	Title       string
	Description string
	Action      string
	EmailID     string
	EmailLabel  string
	EmailName   string
	EmailValue  string
	Error       *inlineAlertView
	Success     string
	Submit      buttonView
}

// defaultNewsletter returns the example in its initial state: the subscribe
// form ready for a no-JS POST, with HTMX as an optional enhancement that swaps
// only the aside fragment.
func defaultNewsletter() newsletterView {
	return newsletterView{
		ID:          "newsletter",
		Title:       "Stay in the loop",
		Description: "Product updates and Gelium UI releases. No spam, unsubscribe anytime.",
		Action:      "/examples/newsletter",
		EmailID:     "newsletter-email",
		EmailLabel:  "Email address",
		EmailName:   "email",
		Submit:      buttonView{Label: "Subscribe", Variant: "primary", Submit: true},
	}
}

// newsletterExampleMarkdown is the copy of the GET example page. It lives in
// code (not under content/) because the route is a noindex /examples/*
// contract demonstration, not a component docs page.
const newsletterExampleMarkdown = `# Newsletter example

The newsletter is a **zero-JS** subscription form: the server owns the whole
contract, so there is nothing to script client-side.

- **POST + 422** — an invalid email is rejected with status 422 and the
  ` + "`X-Gelium-Validation: true`" + ` header, re-rendering the aside with an
  inline alert and the submitted value preserved.
- **POST → 200 success** — a valid email replaces the form with a persistent
  ` + "`role=\"status\"`" + ` confirmation.

Try it below: submit an invalid address to exercise the 422 contract, then a
valid one for the success view. No subscription is performed — this page is a
contract demonstration only (` + "`noindex`" + `).

> HTMX is an optional enhancement: the form posts to the same URL and only the
> aside fragment is swapped.
`

func (s *server) newsletterExample(w http.ResponseWriter, r *http.Request) {
	newsletter := defaultNewsletter()
	s.renderMarkdown(w, r, pageView{
		Title:      "Newsletter",
		Newsletter: &newsletter,
	}, newsletterExampleMarkdown, "/examples/newsletter")
}

// newsletterSubscribe completes the no-JS server round-trip for the newsletter
// example, mirroring text_field.go: an invalid email is rejected with a 422
// carrying the validation header (HX branch) and an inline alert; a valid one
// swaps the form for a persistent success view. Without JavaScript the whole
// page is re-rendered; with HTMX only the aside fragment is swapped.
func (s *server) newsletterSubscribe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	value := strings.TrimSpace(r.FormValue("email"))
	data := defaultNewsletter()
	data.EmailValue = value

	isHX := strings.EqualFold(r.Header.Get("HX-Request"), "true")
	status := http.StatusOK
	if !validEmail(value) {
		status = http.StatusUnprocessableEntity
		data.Error = &inlineAlertView{
			Tone:  "error",
			Title: "Invalid email address",
			Body:  "Enter a valid email address to subscribe.",
		}
	} else {
		data.Success = "You're subscribed — watch your inbox for a confirmation email."
	}

	if !isHX {
		s.renderMarkdownStatus(w, r, pageView{
			Title:      "Newsletter",
			Newsletter: &data,
		}, newsletterExampleMarkdown, "/examples/newsletter", status)
		return
	}

	var rendered bytes.Buffer
	if err := s.templates.ExecuteTemplate(&rendered, "newsletter", data); err != nil {
		http.Error(w, "newsletter unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if status == http.StatusUnprocessableEntity {
		w.Header().Set("X-Gelium-Validation", "true")
	}
	w.WriteHeader(status)
	_, _ = w.Write(rendered.Bytes())
}

// validEmail is a deliberately minimal example validator: an address must have
// a non-empty local part, an @, a dotted domain and a non-empty domain. It
// proves the 422 contract without pretending to be a full RFC validator.
func validEmail(s string) bool {
	i := strings.LastIndexByte(s, '@')
	if i <= 0 || i == len(s)-1 {
		return false
	}
	return strings.Contains(s[i+1:], ".")
}

// newsletterDocs renders the Newsletter component docs page (Phase F public
// content pattern): a zero-JS subscription form whose whole contract is
// server-driven — POST + 422 with an inline alert on an invalid email, a
// persistent success view on a valid one. The page documents the aside
// pattern; the live example lives at /examples/newsletter.
func (s *server) newsletterDocs(w http.ResponseWriter, r *http.Request) {
	s.renderMarkdownPage(w, r, pageView{
		Title: "Newsletter",
	}, "content/newsletter.md")
}
