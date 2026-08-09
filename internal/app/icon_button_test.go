package app

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	webassets "loomui/web"
)

func renderIconButton(t *testing.T, view iconButtonView) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/icon-button.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "icon-button", view); err != nil {
		t.Fatalf("execute icon-button template: %v", err)
	}
	return rendered.String()
}

func TestIconButtonRendersNativeActionButtonWithAccessibleName(t *testing.T) {
	rendered := renderIconButton(t, iconButtonView{
		Label:   "Add to favorites",
		Variant: "standard",
		IconSVG: template.HTML(`<svg aria-hidden="true" focusable="false" viewBox="0 0 24 24"><path d="M12 21l-1.5-1.4C5 15.2 2 12.4 2 9 2 6.2 4.2 4 7 4c1.7 0 3.4.8 5 2.2C13.6 4.8 15.3 4 17 4c2.8 0 5 2.2 5 5 0 3.4-3 6.2-8.5 10.6L12 21z"/></svg>`),
	})

	for _, contract := range []string{
		`<button type="button" class="ui-icon-button ui-icon-button-standard"`,
		`aria-label="Add to favorites"`,
		`<svg aria-hidden="true" focusable="false"`,
	} {
		if !strings.Contains(rendered, contract) {
			t.Errorf("standard icon button = %q, want %s", rendered, contract)
		}
	}
	if strings.Contains(rendered, "href=") {
		t.Errorf("action icon button = %q, must not render a link", rendered)
	}
}

func TestIconButtonRequiresAccessibleName(t *testing.T) {
	named := renderIconButton(t, iconButtonView{
		Label:        "Add to favorites",
		Variant:      "standard",
		IconSVG:      template.HTML(`<svg aria-hidden="true" focusable="false" viewBox="0 0 24 24"><path d="M0 0h24v24H0z"/></svg>`),
		VisibleLabel: true,
	})
	if !strings.Contains(named, `<span class="ui-icon-button-label">Add to favorites</span>`) {
		t.Errorf("standard icon button = %q, visible label must carry the accessible name", named)
	}

	iconOnly := renderIconButton(t, iconButtonView{
		Label:   "Add to favorites",
		Variant: "filled",
		IconSVG: template.HTML(`<svg aria-hidden="true" focusable="false" viewBox="0 0 24 24"><path d="M0 0h24v24H0z"/></svg>`),
	})
	if !strings.Contains(iconOnly, `aria-label="Add to favorites"`) {
		t.Errorf("icon-only icon button = %q, must carry a non-empty aria-label", iconOnly)
	}

	unlabelled := renderIconButton(t, iconButtonView{
		Variant: "standard",
		IconSVG: template.HTML(`<svg aria-hidden="true" focusable="false" viewBox="0 0 24 24"><path d="M0 0h24v24H0z"/></svg>`),
	})
	if strings.Contains(unlabelled, "aria-label=\"\"") {
		t.Errorf("icon button = %q, aria-label must never be empty", unlabelled)
	}
}

func TestIconButtonDeclarativeInvokerAttributesOnlyOnActiveNativeButtons(t *testing.T) {
	active := renderIconButton(t, iconButtonView{
		Label: "Open menu", Variant: "standard", IconSVG: iconButtonFavoriteSVG,
		Command: "show-menu", CommandFor: "example-menu", Value: "menu",
	})
	for _, contract := range []string{`type="button"`, `command="show-menu"`, `commandfor="example-menu"`, `value="menu"`} {
		if !strings.Contains(active, contract) {
			t.Errorf("active icon button = %q, want %s", active, contract)
		}
	}

	forbidden := []string{" command=", " commandfor=", " value="}
	for _, view := range []iconButtonView{
		{Label: "Link", Variant: "standard", IconSVG: iconButtonFavoriteSVG, Href: "/docs", Command: "show-menu", CommandFor: "example-menu", Value: "menu"},
		{Label: "Disabled", Variant: "standard", IconSVG: iconButtonFavoriteSVG, Disabled: true, Command: "show-menu", CommandFor: "example-menu", Value: "menu"},
	} {
		rendered := renderIconButton(t, view)
		for _, attr := range forbidden {
			if strings.Contains(rendered, attr) {
				t.Errorf("inactive icon button %q = %q, must omit %s", view.Label, rendered, attr)
			}
		}
	}
}

func TestIconButtonDisabledRemovesActivationPath(t *testing.T) {
	button := renderIconButton(t, iconButtonView{
		Label: "Add to favorites", Variant: "filled", IconSVG: iconButtonFavoriteSVG, Disabled: true,
	})
	if !strings.Contains(button, ` disabled aria-disabled="true"`) {
		t.Errorf("disabled icon button = %q, want disabled + aria-disabled", button)
	}

	link := renderIconButton(t, iconButtonView{
		Label: "Navigate", Variant: "standard", IconSVG: iconButtonFavoriteSVG, Href: "/must-not-activate", Disabled: true,
	})
	for _, contract := range []string{`role="link"`, `aria-disabled="true"`, `tabindex="-1"`} {
		if !strings.Contains(link, contract) {
			t.Errorf("disabled link icon button = %q, want %s", link, contract)
		}
	}
	if strings.Contains(link, "href=") {
		t.Errorf("disabled link icon button = %q, must not contain any href", link)
	}
}

func TestIconButtonToggleReflectsPressedAndSwapsIcon(t *testing.T) {
	unselected := renderIconButton(t, iconButtonView{
		Label: "Add to favorites", Variant: "filled", IconSVG: iconButtonFavoriteSVG, Toggle: true,
	})
	if !strings.Contains(unselected, `aria-pressed="false"`) {
		t.Errorf("unselected toggle = %q, want aria-pressed=false", unselected)
	}
	if !strings.Contains(unselected, string(iconButtonFavoriteSVG)) {
		t.Errorf("unselected toggle = %q, want unselected icon", unselected)
	}

	selected := renderIconButton(t, iconButtonView{
		Label: "Add to favorites", Variant: "filled", IconSVG: iconButtonFavoriteSVG,
		SelectedIcon: iconButtonCheckSVG, Toggle: true, Selected: true, AriaLabelSelected: "Remove from favorites",
	})
	if !strings.Contains(selected, `aria-pressed="true"`) {
		t.Errorf("selected toggle = %q, want aria-pressed=true", selected)
	}
	if !strings.Contains(selected, string(iconButtonCheckSVG)) {
		t.Errorf("selected toggle = %q, want selected icon", selected)
	}
	if !strings.Contains(selected, `aria-label="Remove from favorites"`) {
		t.Errorf("selected toggle = %q, want selected aria-label", selected)
	}
}

func TestIconButtonActiveHrefRemainsNavigable(t *testing.T) {
	rendered := renderIconButton(t, iconButtonView{
		Label: "Navigate home", Variant: "standard", IconSVG: iconButtonHomeSVG, Href: "/",
	})
	if !strings.Contains(rendered, `href="/"`) {
		t.Errorf("active link icon button = %q, want href", rendered)
	}
	for _, inactive := range []string{`aria-disabled="true"`, `tabindex="-1"`} {
		if strings.Contains(rendered, inactive) {
			t.Errorf("active link icon button = %q, must not contain %s", rendered, inactive)
		}
	}
}

func TestIconButtonDocsRenderEveryVariantAndState(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/components/icon-button", nil)
	res := httptest.NewRecorder()

	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Icon button</h1>`,
		`class="ui-icon-button ui-icon-button-standard"`,
		`class="ui-icon-button ui-icon-button-filled"`,
		`class="ui-icon-button ui-icon-button-filled-tonal"`,
		`class="ui-icon-button ui-icon-button-outlined"`,
		`disabled aria-disabled="true"`,
		`href="/"`,
		`aria-pressed="false"`,
		`aria-pressed="true"`,
		`aria-label="Remove from favorites"`,
		`<svg class="ui-icon" aria-hidden="true" focusable="false"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("icon button docs do not contain contract %q", contract)
		}
	}
}
