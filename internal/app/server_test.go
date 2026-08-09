package app

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	webassets "loomui/web"
)

func openingTagWithID(t *testing.T, body, element, id string) string {
	t.Helper()
	pattern := `<` + element + `\b[^>]*\bid="` + regexp.QuoteMeta(id) + `"[^>]*>`
	tag := regexp.MustCompile(pattern).FindString(body)
	if tag == "" {
		t.Fatalf("body is missing <%s> with id %q", element, id)
	}
	return tag
}

func renderButton(t *testing.T, view buttonView) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/button.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "button", view); err != nil {
		t.Fatalf("execute button template: %v", err)
	}
	return rendered.String()
}

func renderTextField(t *testing.T, view textFieldView) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/text-field.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "text-field", view); err != nil {
		t.Fatalf("execute text field template: %v", err)
	}
	return rendered.String()
}

func TestTextFieldRendersMaterialControlAnatomy(t *testing.T) {
	rendered := renderTextField(t, textFieldView{
		ID: "account-email", Label: "Email", Value: "ada@example.com", Variant: "outlined",
	})

	control := regexp.MustCompile(`<span class="ui-text-field-control">(?s:.*?)</span>`).FindString(rendered)
	if control == "" {
		t.Fatal("text field must group its floating label and native input in a Material control container")
	}
	labelIndex := strings.Index(control, `<label for="account-email">Email</label>`)
	inputIndex := strings.Index(control, `<input`)
	if labelIndex < 0 || inputIndex < 0 || labelIndex > inputIndex {
		t.Errorf("control = %q, want associated floating label before native input", control)
	}
}

func TestTextFieldErrorRendersTrustedDecorativeTrailingIcon(t *testing.T) {
	for _, tt := range []struct {
		name     string
		textarea bool
		native   string
	}{
		{name: "input", native: "<input"},
		{name: "textarea", textarea: true, native: "<textarea"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rendered := renderTextField(t, textFieldView{
				ID: "invalid-field", Label: "Value", Variant: "outlined", Textarea: tt.textarea, Error: "Invalid value",
			})
			control := regexp.MustCompile(`<span class="ui-text-field-control">(?s:.*?)</span>`).FindString(rendered)
			nativeIndex := strings.Index(control, tt.native)
			iconIndex := strings.Index(control, `<svg class="ui-text-field-error-icon" aria-hidden="true" focusable="false"`)
			if nativeIndex < 0 || iconIndex < nativeIndex {
				t.Errorf("error control = %q, want an aria-hidden, unfocusable decorative trusted SVG after the native %s", control, tt.name)
			}
			if strings.Contains(control, "&lt;svg") {
				t.Error("fixed internal error icon must render as trusted SVG markup")
			}
		})
	}
}

func TestTextFieldNativeControlsExposeEmptyStateForFloatingLabel(t *testing.T) {
	for _, tt := range []struct {
		name     string
		textarea bool
		element  string
	}{
		{name: "input", element: "input"},
		{name: "textarea", textarea: true, element: "textarea"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rendered := renderTextField(t, textFieldView{
				ID: "profile-field", Label: "Profile", Variant: "filled", Textarea: tt.textarea,
			})
			control := openingTagWithID(t, rendered, tt.element, "profile-field")
			if !strings.Contains(control, `placeholder=" "`) {
				t.Errorf("%s = %q, want an invisible placeholder hook for the CSS empty/populated state", tt.element, control)
			}
		})
	}
}

func TestButtonDeclarativeInvokerAttributesOnlyRenderOnActiveNativeButtons(t *testing.T) {
	active := renderButton(t, buttonView{
		Label: "Open dialog", Variant: "text", Command: "show-modal", CommandFor: "confirm-dialog", Autofocus: true, Value: "confirm",
	})
	for _, contract := range []string{`class="ui-button ui-button-text"`, `type="button"`, `command="show-modal"`, `commandfor="confirm-dialog"`, `autofocus`, `value="confirm"`} {
		if !strings.Contains(active, contract) {
			t.Errorf("active button = %q, want %s", active, contract)
		}
	}

	views := []buttonView{
		{Label: "Link", Variant: "text", Href: "/docs", Command: "show-modal", CommandFor: "confirm-dialog", Autofocus: true, Value: "confirm"},
		{Label: "Disabled", Variant: "text", Disabled: true, Command: "show-modal", CommandFor: "confirm-dialog", Autofocus: true, Value: "confirm"},
		{Label: "Loading", Variant: "text", Loading: true, Command: "show-modal", CommandFor: "confirm-dialog", Autofocus: true, Value: "confirm"},
	}
	for _, view := range views {
		rendered := renderButton(t, view)
		for _, forbidden := range []string{" command=", " commandfor=", " autofocus", " value="} {
			if strings.Contains(rendered, forbidden) {
				t.Errorf("inactive button %q = %q, must omit %s", view.Label, rendered, forbidden)
			}
		}
	}
	for _, view := range []buttonView{
		{Label: "Unpaired command", Variant: "text", Command: "show-modal"},
		{Label: "Unpaired target", Variant: "text", CommandFor: "confirm-dialog"},
	} {
		rendered := renderButton(t, view)
		if strings.Contains(rendered, " command=") || strings.Contains(rendered, " commandfor=") {
			t.Errorf("unpaired invoker attributes for %q = %q, must omit both", view.Label, rendered)
		}
	}
}

func TestButtonLoadingAccessibleNameDerivesFromLabel(t *testing.T) {
	tests := []struct {
		name string
		href string
	}{
		{name: "button"},
		{name: "link", href: "/must-not-activate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := renderButton(t, buttonView{
				Label:   "Delete account",
				Variant: "primary",
				Href:    tt.href,
				Loading: true,
			})

			if !strings.Contains(rendered, `<span class="sr-only">Loading Delete account</span>`) {
				t.Errorf("loading %s = %q, want accessible name derived from its label", tt.name, rendered)
			}
			if strings.Contains(rendered, "Saving changes") {
				t.Errorf("loading %s = %q, must not contain a hardcoded accessible name", tt.name, rendered)
			}
		})
	}
}

func TestButtonHrefDisabledHasNoActivationPath(t *testing.T) {
	rendered := renderButton(t, buttonView{
		Label:    "Unavailable destination",
		Variant:  "primary",
		Href:     "/must-not-activate",
		Disabled: true,
	})

	for _, attribute := range []string{`role="link"`, `aria-disabled="true"`, `tabindex="-1"`} {
		if !strings.Contains(rendered, attribute) {
			t.Errorf("inactive link = %q, want attribute %s", rendered, attribute)
		}
	}
	if strings.Contains(rendered, "href=") {
		t.Errorf("inactive link = %q, must not contain any href attribute", rendered)
	}
}

func TestButtonHrefLoadingHasNoActivationPathAndAnnouncesState(t *testing.T) {
	rendered := renderButton(t, buttonView{
		Label:   "Save changes",
		Variant: "primary",
		Href:    "/must-not-activate",
		Loading: true,
	})

	for _, contract := range []string{
		`role="link"`,
		`aria-disabled="true"`,
		`tabindex="-1"`,
		`aria-busy="true"`,
		`<span class="sr-only">Loading Save changes</span>`,
	} {
		if !strings.Contains(rendered, contract) {
			t.Errorf("loading link = %q, want %s", rendered, contract)
		}
	}
	if strings.Contains(rendered, "href=") {
		t.Errorf("loading link = %q, must not contain any href attribute", rendered)
	}
}

func TestButtonActiveHrefRemainsNavigable(t *testing.T) {
	rendered := renderButton(t, buttonView{
		Label:   "Read documentation",
		Variant: "outline",
		Href:    "/docs",
	})

	if !strings.Contains(rendered, `href="/docs"`) {
		t.Errorf("active link = %q, want destination href", rendered)
	}
	for _, inactiveAttribute := range []string{`aria-disabled="true"`, `tabindex="-1"`} {
		if strings.Contains(rendered, inactiveAttribute) {
			t.Errorf("active link = %q, must not contain %s", rendered, inactiveAttribute)
		}
	}
}

func TestButtonRendersEachTrustedInlineSVGIconSlotUnescaped(t *testing.T) {
	icons := []template.HTML{
		template.HTML(`<svg aria-hidden="true" focusable="false" viewBox="0 0 10 10"><path data-icon="save" d="M1 1h8v8H1z"></path></svg>`),        // #nosec G203 -- test fixture models trusted internal markup.
		template.HTML(`<svg aria-hidden="true" focusable="false" viewBox="0 0 12 12"><circle data-icon="add" cx="6" cy="6" r="5"></circle></svg>`), // #nosec G203 -- test fixture models trusted internal markup.
	}

	for index, icon := range icons {
		rendered := renderButton(t, buttonView{
			Label:   "Action",
			Variant: "primary",
			IconSVG: icon,
		})
		if !strings.Contains(rendered, string(icon)) {
			t.Errorf("icon case %d = %q, want trusted SVG rendered unescaped", index, rendered)
		}
		if strings.Contains(rendered, "&lt;svg") {
			t.Errorf("icon case %d = %q, SVG slot was escaped", index, rendered)
		}
	}
}

func TestHealthzReturnsPlainTextOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	res := httptest.NewRecorder()

	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if got := res.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q, want plain UTF-8 text", got)
	}
	if got := res.Body.String(); got != "ok\n" {
		t.Errorf("body = %q, want %q", got, "ok\\n")
	}
}

func TestHomeRendersMarkdownInsideDogfoodedLayout(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Loom UI</h1>`,
		`<main`,
		`class="ui-button ui-button-primary"`,
		`href="/components/button"`,
		`src="/static/htmx.min.js?v=0.4.0"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("home does not contain contract %q", contract)
		}
	}
}

func TestLayoutCacheBustsEmbeddedAssetsAcrossExeUpgrades(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	body := res.Body.String()

	for _, asset := range []string{
		`href="/static/app.css?v=0.4.0"`,
		`src="/static/htmx.min.js?v=0.4.0"`,
		`src="/static/app.js?v=0.4.0"`,
	} {
		if !strings.Contains(body, asset) {
			t.Errorf("layout must cache-bust upgraded embedded asset %s", asset)
		}
	}
}

func TestButtonDocsRenderEveryVariantAndAccessibleState(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/components/button", nil)
	res := httptest.NewRecorder()

	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Button</h1>`,
		`class="ui-button ui-button-primary"`,
		`class="ui-button ui-button-secondary"`,
		`class="ui-button ui-button-outline"`,
		`disabled`,
		`aria-disabled="true"`,
		`aria-busy="true"`,
		`<span class="sr-only">Loading Save changes</span>`,
		`<svg aria-hidden="true" focusable="false"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("button docs do not contain contract %q", contract)
		}
	}
}

func TestDialogDocsRouteDogfoodsNativeDeclarativeDialog(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/dialog", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<title>Dialog · Loom UI</title>`, `<h1>Dialog</h1>`, `href="/components/dialog"`,
		`id="confirm-dialog-title"`, `id="confirm-dialog-description"`,
		`command="show-modal" commandfor="confirm-dialog"`,
		`command="request-close" commandfor="confirm-dialog"`,
		`command="close" commandfor="confirm-dialog" value="confirm"`,
		`class="ui-button ui-button-text"`, `autofocus`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("dialog docs are missing %q", contract)
		}
	}
	dialog := openingTagWithID(t, body, "dialog", "confirm-dialog")
	for _, attribute := range []string{`closedby="any"`, `aria-labelledby="confirm-dialog-title"`, `aria-describedby="confirm-dialog-description"`} {
		if !strings.Contains(dialog, attribute) {
			t.Errorf("dialog = %q, want %s", dialog, attribute)
		}
	}
	for _, forbidden := range []string{" open", " role=", " aria-modal=", " tabindex="} {
		if strings.Contains(dialog, forbidden) {
			t.Errorf("dialog = %q, must omit redundant attribute %s", dialog, forbidden)
		}
	}
	for _, id := range []string{"confirm-dialog", "confirm-dialog-title", "confirm-dialog-description"} {
		if got := strings.Count(body, `id="`+id+`"`); got != 1 {
			t.Errorf("id %q occurs %d times, want exactly once", id, got)
		}
	}
}

func TestDialogDocsExplainProgressiveBrowserCompatibility(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/dialog", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}

	body := res.Body.String()
	for _, contract := range []string{
		"supporting browsers",
		"Baseline Low",
		"no component JavaScript fallback",
		"server-rendered fallback or adapter",
		"request-close",
		"newer than the invoker commands",
		"not Baseline",
		"Chromium-only",
		"instant or asymmetric",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("dialog docs are missing compatibility contract %q", contract)
		}
	}
}

func TestDialogDocsRouteKeepsMethodAndUnknownRouteSemantics(t *testing.T) {
	for _, tt := range []struct {
		method, path string
		want         int
	}{
		{http.MethodPost, "/components/dialog", http.StatusMethodNotAllowed},
		{http.MethodGet, "/components/dialog/missing", http.StatusNotFound},
	} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(tt.method, tt.path, nil))
		if res.Code != tt.want {
			t.Errorf("%s %s status = %d, want %d", tt.method, tt.path, res.Code, tt.want)
		}
	}
}

func TestTextFieldDocsRouteDogfoodsNativeAccessibleVariantsAndStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/text-field", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Text field</h1>`,
		`href="/components/text-field"`,
		`ui-text-field-outlined`,
		`ui-text-field-filled`,
		`<textarea`,
		`disabled`,
		`type="submit"`,
		`<section class="component-preview" aria-label="Text field validation example" id="validation-example">`,
		`action="/examples/text-field/validate#validation-example"`,
		`hx-post="/examples/text-field/validate"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("text field docs do not contain contract %q", contract)
		}
	}

	for id, label := range map[string]string{
		"text-normal":            "Name",
		"text-helper":            "Email",
		"text-error":             "Username",
		"text-disabled":          "Account ID",
		"text-disabled-outlined": "Server path",
		"text-disabled-textarea": "Changelog",
		"text-textarea":          "Biography",
	} {
		if !strings.Contains(body, `<label for="`+id+`">`+label+`</label>`) {
			t.Errorf("field %q is missing its associated label", id)
		}
	}

	for _, id := range []string{"text-disabled", "text-disabled-outlined"} {
		input := openingTagWithID(t, body, "input", id)
		if !strings.Contains(input, "disabled") {
			t.Errorf("dogfooded disabled input %q is missing the native disabled attribute", id)
		}
	}
	textarea := openingTagWithID(t, body, "textarea", "text-disabled-textarea")
	if !strings.Contains(textarea, "disabled") {
		t.Error("dogfooded disabled textarea is missing the native disabled attribute")
	}

	helperInput := openingTagWithID(t, body, "input", "text-helper")
	if !strings.Contains(helperInput, `aria-describedby="text-helper-help"`) || !strings.Contains(body, `id="text-helper-help"`) {
		t.Error("helper input must reference its unique helper ID")
	}
	errorInput := openingTagWithID(t, body, "input", "text-error")
	for _, attribute := range []string{`aria-invalid="true"`, `aria-describedby="text-error-error"`} {
		if !strings.Contains(errorInput, attribute) {
			t.Errorf("error input is missing %s", attribute)
		}
	}
	if !strings.Contains(body, `id="text-error-error"`) || !strings.Contains(body, `role="alert"`) || !strings.Contains(body, `Error:`) {
		t.Error("error state must be associated and visibly signalled beyond color")
	}
	if strings.Contains(errorInput, `autofocus`) {
		t.Error("the static Username error demo must not steal focus on GET")
	}
}

func TestTextFieldDocsExplainHTTP422AndHTMXSwapContract(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/text-field", nil))
	body := res.Body.String()
	for _, contract := range []string{"HTTP 422", "X-Loom-Validation", "htmx:beforeSwap", "shouldSwap", "isError", "outerHTML", "without JavaScript", "complete documentation page"} {
		if !strings.Contains(body, contract) {
			t.Errorf("text field docs are missing 422/HTMX explanation %q", contract)
		}
	}
}

func TestTextFieldValidationWithHXRejectsWhitespaceWithCompleteAccessibleFormFragment(t *testing.T) {
	form := strings.NewReader("name=+++++")
	req := httptest.NewRequest(http.MethodPost, "/examples/text-field/validate", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	if got := res.Header().Get("X-Loom-Validation"); got != "true" {
		t.Errorf("X-Loom-Validation = %q, want %q", got, "true")
	}
	body := res.Body.String()
	if strings.Contains(body, `<!doctype html>`) || strings.Contains(body, `<title>`) {
		t.Error("HX 422 response must be a form fragment, not a complete document")
	}
	if !strings.Contains(body, `<form`) || !strings.Contains(body, `Name is required`) {
		t.Error("422 response must return the complete form and required error")
	}
	input := openingTagWithID(t, body, "input", "validation-name")
	for _, contract := range []string{`name="name"`, `value="     "`, `aria-invalid="true"`, `aria-describedby="validation-name-error"`} {
		if !strings.Contains(input, contract) {
			t.Errorf("422 input is missing %s", contract)
		}
	}
	if !strings.Contains(body, `id="validation-name-error"`) || !strings.Contains(body, `role="alert"`) {
		t.Error("422 error must be associated and announced")
	}
	if strings.Contains(input, `autofocus`) {
		t.Error("HX 422 fragment must not add autofocus")
	}
}

func TestTextFieldValidationWithoutHXRejectsWhitespaceInCompleteDocumentationPage(t *testing.T) {
	form := strings.NewReader("name=+++++")
	req := httptest.NewRequest(http.MethodPost, "/examples/text-field/validate", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	if got := res.Header().Get("X-Loom-Validation"); got != "" {
		t.Errorf("non-HX X-Loom-Validation = %q, want no fragment-only validation header", got)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<!doctype html>`,
		`<title>Text field · Loom UI</title>`,
		`<nav aria-label="Primary">`,
		`<article class="prose"><h1>Text field</h1>`,
		`aria-label="Text field examples"`,
		`aria-label="Text field validation example"`,
		`Name is required`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("full-page 422 response is missing %q", contract)
		}
	}
	input := openingTagWithID(t, body, "input", "validation-name")
	for _, contract := range []string{`value="     "`, `aria-invalid="true"`, `aria-describedby="validation-name-error"`, `autofocus`} {
		if !strings.Contains(input, contract) {
			t.Errorf("full-page 422 input is missing %s", contract)
		}
	}
	demoInput := openingTagWithID(t, body, "input", "text-error")
	if strings.Contains(demoInput, `autofocus`) {
		t.Error("full-page 422 must focus only the submitted validation field, not the static Username demo")
	}
}

func TestTextFieldValidationWithHXAcceptsValueAndReturnsAccessibleSuccessFormFragment(t *testing.T) {
	form := strings.NewReader("name=Ada+Lovelace")
	req := httptest.NewRequest(http.MethodPost, "/examples/text-field/validate", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Contains(body, `<!doctype html>`) || strings.Contains(body, `<title>`) {
		t.Error("HX 200 response must be a form fragment, not a complete document")
	}
	input := openingTagWithID(t, body, "input", "validation-name")
	for _, contract := range []string{`name="name"`, `value="Ada Lovelace"`, `aria-describedby="validation-name-help"`} {
		if !strings.Contains(input, contract) {
			t.Errorf("200 input is missing %s", contract)
		}
	}
	if !strings.Contains(body, `id="validation-name-help"`) || !strings.Contains(body, `role="status"`) || !strings.Contains(body, `Name accepted`) {
		t.Error("200 form must include an associated accessible success message")
	}
	if strings.Contains(input, `aria-invalid="true"`) {
		t.Error("valid input must not be marked invalid")
	}
}

func TestTextFieldValidationWithoutHXAcceptsValueInCompleteDocumentationPage(t *testing.T) {
	form := strings.NewReader("name=Ada+Lovelace")
	req := httptest.NewRequest(http.MethodPost, "/examples/text-field/validate", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<!doctype html>`,
		`<title>Text field · Loom UI</title>`,
		`<nav aria-label="Primary">`,
		`<article class="prose"><h1>Text field</h1>`,
		`aria-label="Text field examples"`,
		`aria-label="Text field validation example"`,
		`role="status"`,
		`Name accepted`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("full-page 200 response is missing %q", contract)
		}
	}
	input := openingTagWithID(t, body, "input", "validation-name")
	for _, contract := range []string{`value="Ada Lovelace"`, `aria-describedby="validation-name-help"`} {
		if !strings.Contains(input, contract) {
			t.Errorf("full-page 200 input is missing %s", contract)
		}
	}
	if strings.Contains(input, `aria-invalid="true"`) {
		t.Error("valid full-page input must not be marked invalid")
	}
	if strings.Contains(input, `autofocus`) {
		t.Error("valid full-page input must not autofocus")
	}
}

func TestOnlyTrueHXRequestHeaderReceivesValidationFragment(t *testing.T) {
	tests := []struct {
		name         string
		header       string
		wantDocument bool
	}{
		{name: "missing", wantDocument: true},
		{name: "false", header: "false", wantDocument: true},
		{name: "surrounded true", header: " true ", wantDocument: true},
		{name: "true", header: "true"},
		{name: "case insensitive true", header: "TRUE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/examples/text-field/validate", strings.NewReader("name=Ada"))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if tt.header != "" {
				req.Header.Set("HX-Request", tt.header)
			}
			res := httptest.NewRecorder()
			New().ServeHTTP(res, req)

			gotDocument := strings.Contains(res.Body.String(), `<!doctype html>`)
			if gotDocument != tt.wantDocument {
				t.Errorf("HX-Request %q document = %t, want %t", tt.header, gotDocument, tt.wantDocument)
			}
		})
	}
}

func TestTextFieldValidationRouteUsesModernMuxMethodSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/examples/text-field/validate", nil))
	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET validation status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}

	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/examples/text-field/missing", nil))
	if res.Code != http.StatusNotFound {
		t.Fatalf("unknown route status = %d, want %d", res.Code, http.StatusNotFound)
	}
}

func TestTextFieldFormAndLocalScriptImplementHTMX422SwapContract(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/text-field", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	form := regexp.MustCompile(`<form\b[^>]*>`).FindString(body)
	if form == "" {
		t.Fatal("docs are missing validation form")
	}
	for _, attribute := range []string{
		`novalidate`,
		`hx-post="/examples/text-field/validate"`,
		`hx-target="this"`,
		`hx-swap="outerHTML"`,
	} {
		if !strings.Contains(form, attribute) {
			t.Errorf("validation form is missing %s", attribute)
		}
	}
	htmxIndex := strings.Index(body, `src="/static/htmx.min.js?v=0.4.0"`)
	appIndex := strings.Index(body, `src="/static/app.js?v=0.4.0"`)
	if htmxIndex < 0 || appIndex < 0 || appIndex < htmxIndex {
		t.Error("local app.js must load after local HTMX")
	}

	asset := httptest.NewRecorder()
	New().ServeHTTP(asset, httptest.NewRequest(http.MethodGet, "/static/app.js", nil))
	if asset.Code != http.StatusOK {
		t.Fatalf("app.js status = %d, want %d", asset.Code, http.StatusOK)
	}
	if got := asset.Header().Get("Content-Type"); got != "text/javascript; charset=utf-8" {
		t.Errorf("app.js Content-Type = %q", got)
	}
	js := asset.Body.String()
	for _, contract := range []string{
		`htmx:beforeSwap`,
		`shouldSwap = true`,
		`isError = false`,
	} {
		if !strings.Contains(js, contract) {
			t.Errorf("app.js is missing 422 hook contract %q", contract)
		}
	}
	validation422 := regexp.MustCompile(`status\s*===\s*422\s*&&\s*event\.detail\.xhr\.getResponseHeader\("X-Loom-Validation"\)\s*===\s*"true"`)
	if !validation422.MatchString(js) {
		t.Error("app.js must only swap a 422 when X-Loom-Validation is true")
	}
}

func TestStaticBuildArtifactsAreServedFromEmbeddedFilesystem(t *testing.T) {
	tests := []struct {
		path        string
		contentType string
		contract    string
	}{
		{path: "/static/app.css", contentType: "text/css; charset=utf-8", contract: ".ui-button"},
		{path: "/static/htmx.min.js", contentType: "text/javascript; charset=utf-8", contract: "htmx"},
		{path: "/static/app.js", contentType: "text/javascript; charset=utf-8", contract: "X-Loom-Validation"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			res := httptest.NewRecorder()
			New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if res.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
			}
			if got := res.Header().Get("Content-Type"); got != tt.contentType {
				t.Errorf("Content-Type = %q, want %q", got, tt.contentType)
			}
			if got := res.Header().Get("Cache-Control"); got != "no-cache" {
				t.Errorf("Cache-Control = %q, want revalidation with no-cache", got)
			}
			if !strings.Contains(res.Body.String(), tt.contract) {
				t.Errorf("asset does not contain build contract %q", tt.contract)
			}
		})
	}
}

func TestMaterialDarkThemeKeepsFilledFieldDistinctFromSurface(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/static/app.css", nil))
	css := res.Body.String()

	if got := strings.Count(css, "--ui-color-surface:#211f26"); got != 2 {
		t.Fatalf("compiled dark theme surface declarations = %d, want 2", got)
	}
	if got := strings.Count(css, "--ui-field-container:#36343b"); got != 2 {
		t.Errorf("compiled dark theme filled container declarations = %d, want 2", got)
	}
	if strings.Contains(css, "--ui-field-container:#211f26") {
		t.Error("dark filled field container must differ from the #211f26 surface")
	}
}

func TestMaterialThemeDefinesTextFieldTypescaleTokens(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/static/app.css", nil))
	css := res.Body.String()

	for _, token := range []string{
		`--ui-type-body-lg:400 1rem/1.5rem var(--ui-font-sans)`,
		`--ui-type-body-sm:400 .75rem/1rem var(--ui-font-sans)`,
		`--ui-type-label-sm:500 .75rem/1rem var(--ui-font-sans)`,
	} {
		if !strings.Contains(css, token) {
			t.Errorf("compiled Material theme is missing text-field typescale definition %q", token)
		}
	}
}

func TestMaterialThemeExposesSemanticFoundationContracts(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/static/app.css", nil))
	css := res.Body.String()

	for _, token := range []string{
		"--ui-color-primary:",
		"--ui-type-display-sm:",
		"--ui-radius-full:",
		"--ui-shadow-5:",
		"--ui-state-dragged-opacity:",
		"--ui-focus-thickness:3px",
		"--ui-focus-offset:2px",
		"--ui-field-container:",
		"--ui-field-border:",
		"--ui-field-error:",
		`.theme-material.dark`,
	} {
		if !strings.Contains(css, token) {
			t.Errorf("compiled Material theme is missing %q", token)
		}
	}
}
func renderToast(t *testing.T, view toastView) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/toast.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "toast", view); err != nil {
		t.Fatalf("execute toast template: %v", err)
	}
	return rendered.String()
}

func TestToastRendersMaterialSnackbarAnatomy(t *testing.T) {
	rendered := renderToast(t, toastView{ID: "toast-1", Type: "success", Role: "status", Message: "Saved", IconSVG: toastIcons["success"], Dismiss: true})
	for _, contract := range []string{
		`<div class="ui-toast ui-toast-success" id="toast-1" role="status">`,
		`aria-hidden="true"`,
		`focusable="false"`,
		`<span class="ui-toast-message">Saved</span>`,
		`<button class="ui-toast-action" type="button" data-loom-toast-dismiss aria-label="Dismiss notification">Dismiss</button>`,
	} {
		if !strings.Contains(rendered, contract) {
			t.Errorf("toast is missing %q", contract)
		}
	}
	if strings.Contains(rendered, "&lt;svg") {
		t.Error("fixed internal icon must render as trusted SVG markup, not escaped text")
	}
}

func TestToastRoleIsAssertiveOnlyForErrors(t *testing.T) {
	for _, tt := range []struct{ typ, role string }{
		{"info", "status"}, {"success", "status"}, {"warning", "status"}, {"error", "alert"},
	} {
		rendered := renderToast(t, newToast(tt.typ, "x", "m"))
		role := regexp.MustCompile(`role="[^"]+"`).FindString(rendered)
		if role != `role="`+tt.role+`"` {
			t.Errorf("toast type %q role = %q, want role=%q", tt.typ, role, tt.role)
		}
	}
}

func TestSanitizeToastTypeFallsBackToClosedVocabulary(t *testing.T) {
	for _, in := range []string{"info", "success", "warning", "error"} {
		if got := sanitizeToastType(in); got != in {
			t.Errorf("sanitizeToastType(%q) = %q, want %q", in, got, in)
		}
	}
	for _, in := range []string{"", "fatal", "SUCCESS", "danger", "success alert"} {
		if got := sanitizeToastType(in); got != "info" {
			t.Errorf("sanitizeToastType(%q) = %q, want info", in, got)
		}
	}
}

func TestToastDocsRenderVariantsAndLiveRegion(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/toast", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`id="loom-toast-region"`,
		`aria-live="polite"`,
		`class="ui-toast ui-toast-success"`,
		`class="ui-toast ui-toast-error"`,
		`name="type"`,
		`hx-post="/examples/toast/demo"`,
		`/components/toast`,
		`?v=0.4.0`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("toast docs page is missing %q", contract)
		}
	}
}

func TestToastDemoNoJSRendersPersistentInlineFeedback(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/toast/demo", strings.NewReader("message=Record+updated&type=success"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	New().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `<div class="ui-toast ui-toast-success" id="toast-demo-result" role="status">`) {
		t.Error("no-JS toast demo must render a persistent inline toast")
	}
	if !strings.Contains(body, "Record updated") {
		t.Error("inline toast must carry the submitted message")
	}
	if got := res.Header().Get("HX-Trigger"); got != "" {
		t.Errorf("no-JS response must not rely on HX-Trigger, got %q", got)
	}
}

func TestToastDemoNoJSEmptyMessageRejects422WithoutToast(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/toast/demo", strings.NewReader("message=&type=error"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	New().ServeHTTP(res, req)
	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	body := res.Body.String()
	if !strings.Contains(body, "Message is required") {
		t.Error("empty message must render an inline validation error")
	}
	if strings.Contains(body, `id="toast-demo-result"`) {
		t.Error("validation failures must not be announced as toast notifications")
	}
}
func TestToastDemoHXReturnsTriggerAndFragment(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/toast/demo", strings.NewReader("message=Saved&type=success"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	New().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	trigger := res.Header().Get("HX-Trigger")
	if trigger == "" {
		t.Fatal("HTMX toast demo must set HX-Trigger")
	}
	var parsed struct {
		Toast struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"loom:toast"`
	}
	if err := json.Unmarshal([]byte(trigger), &parsed); err != nil {
		t.Fatalf("HX-Trigger is not valid JSON: %v", err)
	}
	if parsed.Toast.Type != "success" || parsed.Toast.Message != "Saved" {
		t.Errorf("HX-Trigger = %q, want type success and message Saved", trigger)
	}
	body := res.Body.String()
	if !strings.Contains(body, `<form class="toast-demo-form"`) {
		t.Error("HTMX response must be the toast demo form fragment, not a full page")
	}
	if strings.Contains(body, "<html") {
		t.Error("HTMX response must not contain a full document")
	}
	if strings.Contains(body, `id="toast-demo-result"`) {
		t.Error("HTMX demo must rely on the live region, not prerender an inline toast")
	}
}

func TestToastDemoHXEmptyMessageRejects422WithoutTrigger(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/toast/demo", strings.NewReader("message=&type=success"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	New().ServeHTTP(res, req)
	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	if got := res.Header().Get("X-Loom-Validation"); got != "true" {
		t.Errorf("X-Loom-Validation = %q, want true", got)
	}
	if got := res.Header().Get("HX-Trigger"); got != "" {
		t.Errorf("validation failure must not raise loom:toast, got %q", got)
	}
	body := res.Body.String()
	if !strings.Contains(body, "Message is required") {
		t.Error("fragment must render the inline validation error")
	}
	if !strings.Contains(body, `<form class="toast-demo-form"`) {
		t.Error("HTMX 422 must return the toast demo form fragment")
	}
}

func TestToastTriggerJSONEscapesMessage(t *testing.T) {
	trigger, err := toastTriggerJSON("success", `He said "hi" <b> & "bye"`)
	if err != nil {
		t.Fatalf("toastTriggerJSON: %v", err)
	}
	if strings.Contains(trigger, "<b>") {
		t.Error("toast message HTML must be escaped inside the JSON header")
	}
	var parsed struct {
		Toast struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"loom:toast"`
	}
	if err := json.Unmarshal([]byte(trigger), &parsed); err != nil {
		t.Fatalf("escaped header must remain valid JSON: %v", err)
	}
	if parsed.Toast.Message != `He said "hi" <b> & "bye"` {
		t.Errorf("round-trip message = %q, want original", parsed.Toast.Message)
	}
}

// TestTextFieldDisabledStateKeepsBlockedAttributeAndPrecedesError locks the
// explicit state precedence contract demanded by the roadmap: disabled wins
// over error. A disabled field must carry the real disabled attribute, must
// never advertise aria-invalid or an error alert, and must not expose an error
// icon even when the view model carries both flags.
func TestTextFieldDisabledStateKeepsBlockedAttributeAndPrecedesError(t *testing.T) {
	rendered := renderTextField(t, textFieldView{ID: "d", Label: "Account", Name: "account", Value: "ACCT-1", Variant: "filled", Disabled: true, Error: "Should never surface"})

	input := openingTagWithID(t, rendered, "input", "d")
	if !strings.Contains(input, "disabled") {
		t.Error("disabled field must carry the native disabled attribute")
	}
	for _, forbidden := range []string{`aria-invalid`, `aria-describedby`, `role="alert"`, `<p class="ui-text-field-message"`} {
		if strings.Contains(rendered, forbidden) {
			t.Errorf("disabled state must precede error state, but rendered %q", forbidden)
		}
	}
	if strings.Contains(rendered, "ui-text-field-error") {
		t.Error("disabled field must not expose the error state class")
	}
	if strings.Contains(rendered, "invalid") {
		t.Error("disabled field must not render the error message text")
	}
}

// TestTextFieldDocsDisabledFieldCarriesBlockingAttribute isolates the dogfooded
// disabled preview: the exact <input> with id "text-disabled" must carry the
// native disabled attribute and must not be announced as invalid.
func TestTextFieldDocsDisabledFieldCarriesBlockedAttribute(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/text-field", nil))
	body := res.Body.String()
	disabledInput := openingTagWithID(t, body, "input", "text-disabled")
	if !strings.Contains(disabledInput, "disabled") {
		t.Error("dogfooded disabled input is missing the native disabled attribute")
	}
	if strings.Contains(disabledInput, "aria-invalid") {
		t.Error("dogfooded disabled input must not be announced as invalid")
	}
}

func TestElevationDocsRouteDogfoodsTokenMappedUtilityLevels(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/elevation", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("elevation docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Elevation</h1>`,
		`href="/components/elevation"`,
		`aria-label="Elevation example"`,
		`class="ui-elevation-0"`,
		`class="ui-elevation-1"`,
		`class="ui-elevation-2"`,
		`class="ui-elevation-3"`,
		`class="ui-elevation-4"`,
		`class="ui-elevation-5"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("elevation docs are missing %q", contract)
		}
	}
}

func TestElevationDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/elevation", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST elevation status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestFocusRingDocsRouteDogfoodsSharedFocusVisibleContract(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/focus-ring", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("focus ring docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Focus ring</h1>`,
		`href="/components/focus-ring"`,
		`aria-label="Focus ring example"`,
		`class="focus-demo-grid"`,
		`class="focus-demo-link"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("focus ring docs are missing %q", contract)
		}
	}
}

func TestFocusRingDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/focus-ring", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST focus ring status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestIconDocsRouteDogfoodsTrustedSVGContracts(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/icon", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("icon docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Icon</h1>`,
		`href="/components/icon"`,
		`aria-label="Icon examples"`,
		`class="ui-icon"`,
		`aria-hidden="true"`,
		`focusable="false"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("icon docs are missing %q", contract)
		}
	}
}

func TestIconDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/icon", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST icon status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestDividerDocsRouteDogfoodsNativeHRSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/divider", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("divider docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Divider</h1>`,
		`href="/components/divider"`,
		`aria-label="Divider examples"`,
		`class="ui-divider"`,
		`class="ui-divider ui-divider-inset"`,
		`class="ui-divider ui-divider-inset-start"`,
		`class="ui-divider ui-divider-inset-end"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("divider docs are missing %q", contract)
		}
	}
}

func TestDividerDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/divider", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST divider status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestCardDocsRouteDogfoodsSemanticRootNodes(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/card", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("card docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Card</h1>`,
		`href="/components/card"`,
		`aria-label="Card examples"`,
		`<article class="ui-card ui-card-elevated"`,
		`<a class="ui-card ui-card-outlined"`,
		`<button class="ui-card ui-card-filled"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("card docs are missing %q", contract)
		}
	}
}

func TestCardDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/card", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST card status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestBadgeDocsRouteDogfoodsDotAndCountStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/badge", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("badge docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Badge</h1>`,
		`href="/components/badge"`,
		`aria-label="Badge examples"`,
		`class="ui-badge"`,
		`class="ui-badge ui-badge-large"`,
		`>3</span>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("badge docs are missing %q", contract)
		}
	}
}

func TestBadgeDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/badge", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST badge status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestCheckboxDocsRouteDogfoodsStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/checkbox", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("checkbox docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Checkbox</h1>`,
		`href="/components/checkbox"`,
		`aria-label="Checkbox examples"`,
		`type="checkbox"`,
		`class="ui-checkbox"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("checkbox docs are missing %q", contract)
		}
	}
}

func TestCheckboxDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/checkbox", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST checkbox status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestRadioDocsRouteDogfoodsStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/radio", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("radio docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Radio</h1>`,
		`href="/components/radio"`,
		`aria-label="Radio examples"`,
		`type="radio"`,
		`class="ui-radio"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("radio docs are missing %q", contract)
		}
	}
}

func TestRadioDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/radio", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST radio status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestSwitchDocsRouteDogfoodsStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/switch", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("switch docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Switch</h1>`,
		`href="/components/switch"`,
		`aria-label="Switch examples"`,
		`type="checkbox"`,
		`class="ui-switch"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("switch docs are missing %q", contract)
		}
	}
}

func TestSwitchDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/switch", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST switch status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestSelectDocsRouteDogfoodsNativeVariantsAndStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/select", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("select docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Select</h1>`,
		`href="/components/select"`,
		`aria-label="Select examples"`,
		`<select`,
		`ui-select-filled`,
		`ui-select-outlined`,
		`<option`,
		`disabled`,
		`aria-invalid="true"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("select docs are missing %q", contract)
		}
	}
}

func TestSelectDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/select", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST select status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestSelectMenuDocsRouteDogfoodsServerDrivenMenu(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/select", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("select docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`class="ui-select-menu"`,
		`aria-label="Select menu example"`,
		`command="show-modal"`,
		`commandfor="select-menu"`,
		`hx-post="/examples/select/menu"`,
		`hx-target="this"`,
		`hx-swap="outerHTML"`,
		`command="request-close"`,
		`class="ui-select-menu-item"`,
		`aria-selected="true"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("select menu docs are missing %q", contract)
		}
	}
}

func TestSelectMenuChangeNoHSetsValueAndRendersSessionState(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/select/menu", strings.NewReader("value=priority&id=select-menu"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `<!doctype html>`) || !strings.Contains(body, `<title>Select · Loom UI</title>`) {
		t.Error("no-HX response must be a complete documentation page")
	}
	if !strings.Contains(body, `value="priority"`) {
		t.Error("selected value must be reflected in the session state")
	}
	if !strings.Contains(body, `aria-selected="true"`) {
		t.Error("selected menu item must carry aria-selected")
	}
}

func TestSelectMenuChangeHXReturnsFragmentReflectingSelection(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/select/menu", strings.NewReader("value=priority&id=select-menu"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Contains(body, `<!doctype html>`) || strings.Contains(body, `<title>`) {
		t.Error("HX response must be a fragment, not a complete document")
	}
	if !strings.Contains(body, `class="ui-select m3-select"`) {
		t.Errorf("fragment must return the M3 select wrapper, got %q", body)
	}
	if !strings.Contains(body, `value="priority"`) {
		t.Error("fragment must persist the newly selected value on the hidden input")
	}
	if strings.Contains(body, `aria-selected="true"`) {
		t.Error("closed-menu fragment must not contain the open menu list")
	}
}

func TestSelectMenuChangeRejectsUnknownValueWith422(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/select/menu", strings.NewReader("value=not-a-plan&id=select-menu"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	New().ServeHTTP(res, req)

	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	if got := res.Header().Get("X-Loom-Validation"); got != "true" {
		t.Errorf("X-Loom-Validation = %q, want true", got)
	}
	body := res.Body.String()
	if !strings.Contains(body, "Select a valid option") {
		t.Error("422 fragment must carry a visible validation error")
	}
}

func TestAppJSWiresHTMXToastTriggerEvent(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/static/app.js", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("app.js status = %d, want %d", res.Code, http.StatusOK)
	}
	js := res.Body.String()
	for _, contract := range []string{
		`addEventListener("loom:toast"`,
		`#loom-toast-region`,
		`ui-toast-show`,
		`Dismiss notification`,
	} {
		if !strings.Contains(js, contract) {
			t.Errorf("app.js is missing toast enhancement contract %q", contract)
		}
	}
}
