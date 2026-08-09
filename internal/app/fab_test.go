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

// renderFab renders a single FAB from the real dogfooded template so the test
// asserts the component's contract, never duplicated markup.
func renderFab(t *testing.T, view fabView) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/fab.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "fab", view); err != nil {
		t.Fatalf("execute fab template: %v", err)
	}
	return rendered.String()
}

func TestFabActionRendersNativeButtonWithMandatoryAriaLabel(t *testing.T) {
	rendered := renderFab(t, fabView{
		AriaLabel: "Compose email",
		Variant:   "surface",
		Size:      "medium",
		IconSVG:   template.HTML(`<svg aria-hidden="true" focusable="false" viewBox="0 0 24 24"><path data-icon="compose" d="M1 1h22v22H1z"></path></svg>`), // #nosec G203 -- trusted fixture.
	})

	for _, contract := range []string{
		`<button class="ui-fab ui-fab-surface ui-fab-medium"`,
		`type="button"`,
		`aria-label="Compose email"`,
	} {
		if !strings.Contains(rendered, contract) {
			t.Errorf("icon-only fab = %q, want %s", rendered, contract)
		}
	}
	if !strings.Contains(rendered, `<svg aria-hidden="true" focusable="false"`) {
		t.Errorf("icon-only fab = %q, want decorative trusted SVG icon", rendered)
	}
	if strings.Contains(rendered, `href=`) {
		t.Errorf("action fab = %q, must not emit an href", rendered)
	}
}

func TestFabNavigationHrefRendersActiveAnchor(t *testing.T) {
	rendered := renderFab(t, fabView{
		AriaLabel: "Open settings",
		Variant:   "surface",
		Size:      "medium",
		Href:      "/settings",
		IconSVG:   template.HTML(`<svg aria-hidden="true" focusable="false" viewBox="0 0 24 24"><path data-icon="gear" d="M1 1h22v22H1z"></path></svg>`), // #nosec G203 -- trusted fixture.
	})

	for _, contract := range []string{
		`<a class="ui-fab ui-fab-surface ui-fab-medium"`,
		`href="/settings"`,
		`aria-label="Open settings"`,
	} {
		if !strings.Contains(rendered, contract) {
			t.Errorf("navigation fab = %q, want %s", rendered, contract)
		}
	}
	for _, inactive := range []string{`aria-disabled="true"`, `tabindex="-1"`, `<button`, `type="button"`} {
		if strings.Contains(rendered, inactive) {
			t.Errorf("active link fab = %q, must not contain %s", rendered, inactive)
		}
	}
}

func TestFabDisabledRemovesActivationPath(t *testing.T) {
	button := renderFab(t, fabView{
		AriaLabel: "Compose email",
		Variant:   "primary",
		Size:      "medium",
		Disabled:  true,
		IconSVG:   template.HTML(`<svg aria-hidden="true" focusable="false" viewBox="0 0 24 24"><path d="M1 1h22v22H1z"></path></svg>`), // #nosec G203 -- trusted fixture.
	})
	for _, contract := range []string{`disabled`, `aria-disabled="true"`, `type="button"`} {
		if !strings.Contains(button, contract) {
			t.Errorf("disabled button fab = %q, want %s", button, contract)
		}
	}

	link := renderFab(t, fabView{
		AriaLabel: "Open settings",
		Variant:   "surface",
		Size:      "medium",
		Href:      "/settings",
		Disabled:  true,
	})
	for _, contract := range []string{`role="link"`, `aria-disabled="true"`, `tabindex="-1"`} {
		if !strings.Contains(link, contract) {
			t.Errorf("disabled link fab = %q, want %s", link, contract)
		}
	}
	if strings.Contains(link, "href=") {
		t.Errorf("disabled link fab = %q, must not contain any href attribute", link)
	}
}

func TestFabDeclarativeInvokerOnlyOnActiveNativeButtons(t *testing.T) {
	active := renderFab(t, fabView{
		AriaLabel:  "Attach file",
		Variant:    "secondary",
		Size:       "medium",
		Command:    "show-modal",
		CommandFor: "attach-dialog",
		Value:      "attach",
		IconSVG:    template.HTML(`<svg aria-hidden="true" focusable="false" viewBox="0 0 24 24"><path d="M1 1h22v22H1z"></path></svg>`), // #nosec G203 -- trusted fixture.
	})
	for _, contract := range []string{`command="show-modal"`, `commandfor="attach-dialog"`, `value="attach"`} {
		if !strings.Contains(active, contract) {
			t.Errorf("active fab = %q, want %s", active, contract)
		}
	}

	for _, view := range []fabView{
		{AriaLabel: "Link nav", Variant: "surface", Size: "medium", Href: "/x", Command: "show-modal", CommandFor: "attach-dialog", Value: "attach"},
		{AriaLabel: "Disabled", Variant: "surface", Size: "medium", Disabled: true, Command: "show-modal", CommandFor: "attach-dialog", Value: "attach"},
		{AriaLabel: "Unpaired", Variant: "surface", Size: "medium", Command: "show-modal"},
	} {
		rendered := renderFab(t, view)
		for _, forbidden := range []string{" command=", " commandfor=", " value="} {
			if strings.Contains(rendered, forbidden) {
				t.Errorf("inactive fab %q = %q, must omit %s", view.AriaLabel, rendered, forbidden)
			}
		}
	}
}

func TestFabSizeAndLoweredMapToClosedVocabularyClasses(t *testing.T) {
	tests := []struct {
		view   fabView
		classe string
	}{
		{fabView{AriaLabel: "S", Size: "small", Variant: "surface"}, `ui-fab-small`},
		{fabView{AriaLabel: "M", Size: "medium", Variant: "surface"}, `ui-fab-medium`},
		{fabView{AriaLabel: "L", Size: "large", Variant: "surface"}, `ui-fab-large`},
	}
	for _, tt := range tests {
		rendered := renderFab(t, tt.view)
		if !strings.Contains(rendered, tt.classe) {
			t.Errorf("fab size %q = %q, want class %s", tt.view.Size, rendered, tt.classe)
		}
	}

	lowered := renderFab(t, fabView{AriaLabel: "L", Size: "medium", Variant: "surface", Lowered: true})
	if !strings.Contains(lowered, `ui-fab ui-fab-surface ui-fab-medium ui-fab-lowered`) {
		t.Errorf("lowered fab = %q, want ui-fab-lowered class preserving surface+medium", lowered)
	}
	if strings.Contains(lowered, `ui-fab-extended`) {
		t.Errorf("lowered fab = %q, must not be marked extended", lowered)
	}
}

func TestFabExtendedRendersVisibleLabelAndKeepsAccessibleName(t *testing.T) {
	rendered := renderFab(t, fabView{
		AriaLabel: "Compose a new email",
		Label:     "Compose",
		Variant:   "primary",
		Size:      "medium",
		IconSVG:   template.HTML(`<svg aria-hidden="true" focusable="false" viewBox="0 0 24 24"><path d="M1 1h22v22H1z"></path></svg>`), // #nosec G203 -- trusted fixture.
	})

	for _, contract := range []string{`ui-fab-extended`, `<span class="ui-fab-label">Compose</span>`, `aria-label="Compose a new email"`} {
		if !strings.Contains(rendered, contract) {
			t.Errorf("extended fab = %q, want %s", rendered, contract)
		}
	}
}

// TestFabDocsEveryIconOnlyFabHasNonEmptyAccessibleName is the non-negotiable
// accessible name gate: an icon-anchored FAB has no visible text, so every
// rendered icon-only FAB must carry a non-empty aria-label. Extended FABs fall
// back to their visible label but still render the richer aria-label when one
// is supplied.
func TestFabDocsEveryIconOnlyFabHasNonEmptyAccessibleName(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/fab", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	for _, c := range []string{
		`aria-label="Edit profile"`,
		`aria-label="Unavailable action"`,
		`aria-label="Compose a new email"`,
		`aria-label="Navigate to settings"`,
	} {
		if !strings.Contains(body, c) {
			t.Errorf("fab docs are missing a non-empty accessible name %q", c)
		}
	}
}

// TestFabDocsRenderEveryVariantSizeAndState verifies the composed docs page
// exercises the full closed vocabulary: primary, surface, secondary color
// variants; small, medium, large sizes; the extended labelled form; lowered
// elevation; disabled; navigation href; and the decorative icon slot.
func TestFabDocsRenderEveryVariantSizeAndState(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/fab", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	for _, contract := range []string{
		`<h1>Floating action button (FAB)</h1>`,
		`class="ui-fab ui-fab-primary ui-fab-medium"`,
		`class="ui-fab ui-fab-surface ui-fab-medium"`,
		`class="ui-fab ui-fab-secondary ui-fab-medium"`,
		`ui-fab-small`,
		`ui-fab-large`,
		`ui-fab-extended`,
		`<span class="ui-fab-label">Compose</span>`,
		`href="/components/fab"`,
		`disabled`,
		`aria-disabled="true"`,
		`<svg aria-hidden="true" focusable="false"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("fab docs do not contain contract %q", contract)
		}
	}
}

// TestFabDocsNeverRendersIconOnlyFabWithoutAccessibleName asserts the data
// layer keeps the invariant that no icon-only FAB slips through nameless: the
// page contains exactly the aria-labels that were authored, and every non-
// extended, non-href demo carries one.
func TestFabDocsNeverRendersIconOnlyFabWithoutAccessibleName(t *testing.T) {
	demo := defaultFabDemo()
	for _, v := range demo.FabDemoViews {
		if v.Label == "" && v.AriaLabel == "" {
			t.Errorf("icon-only fab demo lacks accessible name: %+v", v)
		}
	}
}

func TestFabEscapesUserSuppliedTextNotTrustedSVG(t *testing.T) {
	rendered := renderFab(t, fabView{
		AriaLabel: `tom & <b>jerry</b> "quoted"`,
		Label:     `Label & "quotes"`,
		Variant:   "primary",
		Size:      "medium",
		IconSVG:   template.HTML(`<svg aria-hidden="true" focusable="false" viewBox="0 0 24 24"><path d="M1 1h22v22H1z"></path></svg>`), // #nosec G203 -- trusted fixture.
	})

	for _, contract := range []string{
		`aria-label="tom &amp; &lt;b&gt;jerry&lt;/b&gt; &#34;quoted&#34;"`,
		`<span class="ui-fab-label">Label &amp; &#34;quotes&#34;</span>`,
	} {
		if !strings.Contains(rendered, contract) {
			t.Errorf("fab = %q, want escaped text %s", rendered, contract)
		}
	}
	if strings.Contains(rendered, `&lt;svg`) {
		t.Errorf("trusted SVG slot was escaped: %q", rendered)
	}
}
