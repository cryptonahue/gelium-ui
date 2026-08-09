package app

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
