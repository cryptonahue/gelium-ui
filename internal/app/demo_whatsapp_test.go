package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// WhatsApp demo: main app route, 24h window visibility, expired-window
// composer, send/template/typing/read round trips, and the admin section.

func TestWhatsAppDemoRouteRendersConversationsAndWindow(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/demo/whatsapp", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("demo whatsapp status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`LoomChat`,
		`Ana Souza`,
		`Carlos Lima`,
		`María Fernanda`,
		`demo-wa-conv`,
		`demo-wa-window`,
		`restantes`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("demo whatsapp is missing %q", contract)
		}
	}
	// The expired conversation must carry the expired tone.
	if !strings.Contains(body, `demo-wa-window--expired`) {
		t.Error("demo whatsapp must show an expired-window chip")
	}
}

func TestWhatsAppDemoActiveChatRendersThreadAndComposer(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/demo/whatsapp?c=ana", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("demo whatsapp chat status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`demo-wa-bubble--in`,
		`demo-wa-bubble--out`,
		`demo-wa-composer`,
		`demo-wa-chat-head-window`,
		`Ventana de servicio de 24 h`,
		`name="conversation" value="ana"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("demo whatsapp chat is missing %q", contract)
		}
	}
}

func TestWhatsAppExpiredConversationBlocksComposer(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/demo/whatsapp?c=maria", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("demo whatsapp expired status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	// Expired: no free-text composer, template-only notice present.
	if strings.Contains(body, `demo-wa-composer-input`) {
		t.Error("expired conversation must not render the free-text composer")
	}
	for _, contract := range []string{
		`demo-wa-expired`,
		`La ventana de servicio de 24 h venció`,
		`demo-wa-template-send`,
		`name="template"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("expired conversation is missing %q", contract)
		}
	}
}

func TestWhatsAppSearchFiltersConversations(t *testing.T) {
	form := url.Values{}
	form.Set("q", "Carlos")
	req := httptest.NewRequest(http.MethodGet, "/demo/whatsapp?"+form.Encode(), nil)
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("demo whatsapp search status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `Carlos Lima`) {
		t.Error("search for Carlos must include Carlos Lima")
	}
	if strings.Contains(body, `Ana Souza`) {
		t.Error("search for Carlos must exclude Ana Souza")
	}
}

func TestWhatsAppSendAppendsMessageAndRedirects(t *testing.T) {
	form := url.Values{}
	form.Set("conversation", "ana")
	form.Set("message", "Hola de nuevo!")
	req := httptest.NewRequest(http.MethodPost, "/demo/whatsapp/send", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusSeeOther {
		t.Fatalf("send status = %d, want %d", res.Code, http.StatusSeeOther)
	}
	// The message must now appear in the chat.
	res2 := httptest.NewRecorder()
	New().ServeHTTP(res2, httptest.NewRequest(http.MethodGet, "/demo/whatsapp?c=ana", nil))
	if !strings.Contains(res2.Body.String(), "Hola de nuevo!") {
		t.Error("sent message must appear in the chat thread")
	}
}

func TestWhatsAppSendTemplateAppendsAndRedirects(t *testing.T) {
	form := url.Values{}
	form.Set("conversation", "maria")
	form.Set("template", "boas_vindas")
	req := httptest.NewRequest(http.MethodPost, "/demo/whatsapp/send-template", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusSeeOther {
		t.Fatalf("send-template status = %d, want %d", res.Code, http.StatusSeeOther)
	}
}

func TestWhatsAppReadMarksConversationRead(t *testing.T) {
	form := url.Values{}
	form.Set("conversation", "carlos")
	req := httptest.NewRequest(http.MethodPost, "/demo/whatsapp/read", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusSeeOther {
		t.Fatalf("read status = %d, want %d", res.Code, http.StatusSeeOther)
	}
}

func TestWhatsAppAdminRendersTechnicalSurface(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/demo/whatsapp/admin", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("demo whatsapp admin status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`Modo administrador`,
		`GREEN`,
		`YELLOW`,
		`demo-wa-admin-table`,
		`sk_live_`,
		`cloud.datafyapi.com.br`,
		`HMAC-SHA256`,
		`x-datafy`,
		`<main class="demo-wa-admin">`,
		`aria-current="location"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("admin section is missing %q", contract)
		}
	}
}

// TestWhatsAppAdminLandmarksAndLabels closes gaps G7 and G9: the admin demo
// exposes a main landmark, marks the active tab with aria-current (not
// class-only), and labels every emoji-only action link so the emoji glyph is
// never the accessible name.
func TestWhatsAppAdminLandmarksAndLabels(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/demo/whatsapp/admin", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("demo whatsapp admin status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<main class="demo-wa-admin">`,
		`demo-wa-admin-tab--active" href="#numbers" aria-current="location">Números`,
		`aria-label="Configurar número">⚙`,
		`aria-label="Eliminar número">🗑`,
		`aria-label="Editar código QR">✎`,
		`aria-label="Eliminar código QR">🗑`,
		`aria-label="Ver token">👁`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("admin demo is missing %q", contract)
		}
	}
}

// TestWhatsAppDemoDropsRedundantListRole closes gap G10: a native <ul> must not
// restate its own semantics with role="list".
func TestWhatsAppDemoDropsRedundantListRole(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/demo/whatsapp", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("demo whatsapp status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `<ul class="demo-wa-conversations">`) {
		t.Error("conversations list must render as a plain native ul")
	}
	if strings.Contains(body, `role="list"`) {
		t.Error("demo must not carry the redundant role=list on a native ul (G10)")
	}
}

func TestWhatsAppTypingRouteResponds(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/demo/whatsapp/typing?c=ana", nil))

	if res.Code != http.StatusSeeOther {
		t.Fatalf("typing status = %d, want %d", res.Code, http.StatusSeeOther)
	}
}

// TestWhatsAppDemoIsSpanishAndNoIndexed closes gap G2 (lang on Spanish demos)
// and the SEO robots rule: the demo is es and never indexed.
func TestWhatsAppDemoIsSpanishAndNoIndexed(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/demo/whatsapp", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("demo whatsapp status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<html lang="es" class="theme-material">`,
		`<meta name="robots" content="noindex, nofollow">`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("demo whatsapp is missing %q", contract)
		}
	}
}

// TestWhatsAppAdminIsSpanishAndNoIndexed is the admin section counterpart of
// the Spanish/noindex contract.
func TestWhatsAppAdminIsSpanishAndNoIndexed(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/demo/whatsapp/admin", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("demo whatsapp admin status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<html lang="es" class="theme-material">`,
		`<meta name="robots" content="noindex, nofollow">`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("demo whatsapp admin is missing %q", contract)
		}
	}
}

// TestWhatsAppAdminWebhookSaveRedirects closes gap G3: the admin webhook form
// must POST to a real handler that persists and redirects (POST+303), instead
// of 405'ing.
func TestWhatsAppAdminWebhookSaveRedirects(t *testing.T) {
	form := url.Values{}
	form.Set("webhook_url", "https://example.com/hook")
	form.Set("webhook_secret", "whsec_test")
	req := httptest.NewRequest(http.MethodPost, "/demo/whatsapp/admin/webhook", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusSeeOther {
		t.Fatalf("webhook save status = %d, want %d", res.Code, http.StatusSeeOther)
	}
	if loc := res.Header().Get("Location"); loc != "/demo/whatsapp/admin" {
		t.Errorf("redirect location = %q, want /demo/whatsapp/admin", loc)
	}
	res2 := httptest.NewRecorder()
	New().ServeHTTP(res2, httptest.NewRequest(http.MethodGet, "/demo/whatsapp/admin", nil))
	if !strings.Contains(res2.Body.String(), "https://example.com/hook") {
		t.Error("saved webhook URL must render back in the admin form")
	}
}

// TestWhatsAppAdminWebhookSaveGETIs405 keeps the POST-only semantics: a GET to
// the webhook save route stays a 405.
func TestWhatsAppAdminWebhookSaveGETIs405(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/demo/whatsapp/admin/webhook", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET webhook save status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}
