package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewsletterExampleGETRendersZeroJSFormNoindex(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/examples/newsletter", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<!doctype html>`,
		`<h1>Newsletter example</h1>`,
		`<meta name="robots" content="noindex, nofollow">`,
		`<aside class="ui-newsletter" aria-labelledby="newsletter-title">`,
		`<h2 id="newsletter-title" class="ui-newsletter-title">Stay in the loop</h2>`,
		`<form class="ui-newsletter-form" method="post" action="/examples/newsletter"`,
		`hx-post="/examples/newsletter"`,
		`id="newsletter-email"`,
		`name="email"`,
		`type="email"`,
		`required`,
		`<button class="ui-button ui-button-primary" type="submit"><span>Subscribe</span></button>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("newsletter example GET is missing contract %q", contract)
		}
	}
	if strings.Contains(body, "ui-newsletter-success") {
		t.Error("initial GET must render the subscribe form, not the success view")
	}
}

func TestNewsletterPOSTInvalidEmailWithoutHXRejectsWith422FullPage(t *testing.T) {
	form := strings.NewReader("email=not-an-email")
	req := httptest.NewRequest(http.MethodPost, "/examples/newsletter", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	// The fragment-only validation header belongs to the HX branch only.
	if got := res.Header().Get("X-Gelium-Validation"); got != "" {
		t.Errorf("non-HX X-Gelium-Validation = %q, want no fragment-only validation header", got)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<!doctype html>`,
		`<aside class="ui-newsletter"`,
		`ui-inline-alert ui-inline-alert--error`,
		`role="alert"`,
		`Invalid email address`,
		`value="not-an-email"`,
		`aria-invalid="true"`,
		`aria-describedby="newsletter-error"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("full-page 422 newsletter is missing %q", contract)
		}
	}
	input := openingTagWithID(t, body, "input", "newsletter-email")
	if !strings.Contains(input, `type="email"`) || !strings.Contains(input, `name="email"`) {
		t.Error("422 must re-render the real email input preserving its name/type")
	}
}

func TestNewsletterPOSTValidEmailWithoutHXRendersSuccessView(t *testing.T) {
	form := strings.NewReader("email=ada%40example.com")
	req := httptest.NewRequest(http.MethodPost, "/examples/newsletter", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<!doctype html>`,
		`<p class="ui-newsletter-success" role="status">You&#39;re subscribed`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("full-page 200 newsletter is missing %q", contract)
		}
	}
	if strings.Contains(body, `name="email"`) {
		t.Error("success view must replace the form with the confirmation, not keep the email input")
	}
}

func TestNewsletterPOSTInvalidEmailWithHXReturnsFragmentAndHeader(t *testing.T) {
	form := strings.NewReader("email=+++")
	req := httptest.NewRequest(http.MethodPost, "/examples/newsletter", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	if got := res.Header().Get("X-Gelium-Validation"); got != "true" {
		t.Errorf("X-Gelium-Validation = %q, want %q", got, "true")
	}
	body := res.Body.String()
	if strings.Contains(body, `<!doctype html>`) || strings.Contains(body, `<title>`) {
		t.Error("HX 422 response must be the aside fragment, not a complete document")
	}
	for _, contract := range []string{
		`<aside class="ui-newsletter" aria-labelledby="newsletter-title">`,
		`ui-inline-alert ui-inline-alert--error`,
		`role="alert"`,
		`Invalid email address`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("HX 422 aside fragment is missing %q", contract)
		}
	}
}

func TestNewsletterPOSTValidEmailWithHXReturnsSuccessFragment(t *testing.T) {
	form := strings.NewReader("email=grace%40example.com")
	req := httptest.NewRequest(http.MethodPost, "/examples/newsletter", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Contains(body, `<!doctype html>`) || strings.Contains(body, `<title>`) {
		t.Error("HX 200 response must be the aside fragment, not a complete document")
	}
	if !strings.Contains(body, `<p class="ui-newsletter-success" role="status">`) {
		t.Error("HX success fragment must render the persistent confirmation with role=status")
	}
}

func TestNewsletterExampleIsNotAComponentRoute(t *testing.T) {
	// The example lives under /examples/* (noindex) and must not appear in the
	// primary navigation or the components docs.
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	body := res.Body.String()
	if strings.Contains(body, `href="/examples/newsletter"`) {
		t.Error("home nav must not link the /examples/newsletter route")
	}
}
