package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func renderAccordionTest(t *testing.T, view *accordionView) string {
	t.Helper()
	tmpl := parseTestTemplates(t, "templates/*.html")
	var out bytes.Buffer
	if err := tmpl.ExecuteTemplate(&out, "accordion", view); err != nil {
		t.Fatalf("render accordion: %v", err)
	}
	return out.String()
}

func TestAccordionNativeContract(t *testing.T) {
	view := prepareAccordion(&accordionView{
		ID: "faq", Heading: "Questions", Label: "Fallback", Multiple: true,
		Items: []accordionItem{
			{Value: "first", Heading: "First", Body: "<em>escaped</em>", Open: true},
			{Value: "second item", Heading: "Second", Body: "Second body"},
		},
	})
	body := renderAccordionTest(t, view)
	for _, want := range []string{
		`<section class="ui-accordion ui-accordion--behavior-native ui-accordion--execution-native" data-behavior="native" data-execution="native" id="faq" aria-labelledby="faq-heading">`,
		`<h2 class="ui-accordion-heading" id="faq-heading">Questions</h2>`,
		`<details class="ui-accordion-item" open data-value="first">`,
		`<summary class="ui-accordion-trigger" id="faq-trigger-first" aria-controls="faq-panel-first">`,
		`<section class="ui-accordion-panel" id="faq-panel-first" aria-labelledby="faq-trigger-first">`,
		`data-value="second item"`,
		`id="faq-trigger-second-item"`,
		`<span class="ui-accordion-icon-plus">+</span>`,
		`&lt;em&gt;escaped&lt;/em&gt;`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("accordion output missing %q\n%s", want, body)
		}
	}
	for _, forbidden := range []string{`role="button"`, `aria-expanded=`, `disabled`} {
		if strings.Contains(body, forbidden) {
			t.Errorf("native accordion must not emit %q", forbidden)
		}
	}

}

func TestAccordionLabelFallbackAndExclusiveName(t *testing.T) {
	view := prepareAccordion(&accordionView{
		ID: "settings", Label: "Settings", Multiple: false, Name: "settings-sections",
		Items: []accordionItem{{Value: "profile", Heading: "Profile"}, {Value: "security", Heading: "Security"}},
	})
	body := renderAccordionTest(t, view)
	if !strings.Contains(body, `class="ui-accordion `) || !strings.Contains(body, `id="settings" aria-label="Settings"`) {
		t.Error("Label must remain the accessible-name fallback when Heading is absent")
	}
	if got := strings.Count(body, ` name="settings-sections"`); got != 2 {
		t.Fatalf("exclusive mode must name every details item, got %d", got)
	}
	if strings.Contains(body, `data-multiple`) {
		t.Error("accordion must not invent a data-multiple behavior contract")
	}
	multi := renderAccordionTest(t, prepareAccordion(&accordionView{ID: "multi", Multiple: true, Name: "ignored", Items: []accordionItem{{Value: "one"}}}))
	if strings.Contains(multi, ` name="ignored"`) {
		t.Error("multiple mode must not emit exclusive details name")
	}
}

func TestAccordionProfilesValidateAndPreserveNativeFallback(t *testing.T) {
	for _, tc := range []struct {
		name, rawBehavior, rawExecution, wantBehavior, wantExecution string
	}{
		{"valid", "baseui", "htmx", "baseui", "htmx"},
		{"invalid", "<script>", "fetch", "native", "native"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			view := prepareAccordion(&accordionView{ID: "profile", Behavior: accordionBehavior(tc.rawBehavior), Execution: accordionExecution(tc.rawExecution), Items: []accordionItem{{Value: "one"}}})
			body := renderAccordionTest(t, view)
			if !strings.Contains(body, `data-behavior="`+tc.wantBehavior+`"`) || !strings.Contains(body, `data-execution="`+tc.wantExecution+`"`) {
				t.Fatalf("profiles not normalized safely: %s", body)
			}
			if strings.Contains(body, "hx-") {
				t.Error("accordion profile must not emit arbitrary or fake hx attributes")
			}
		})
	}
}

func TestAccordionBehaviorPolicies(t *testing.T) {
	for _, behavior := range []accordionBehavior{accordionBehaviorBasecoat, accordionBehaviorMaterial} {
		view := prepareAccordion(&accordionView{ID: "policy", Behavior: behavior, Items: []accordionItem{
			{Value: "first", Open: true}, {Value: "second", Open: true},
		}})
		if view.Multiple || view.Name != "policy-group" {
			t.Fatalf("%s must default to a named exclusive group: multiple=%v name=%q", behavior, view.Multiple, view.Name)
		}
		if !view.Items[0].Open || view.Items[1].Open {
			t.Fatalf("%s must retain only the first initial open item", behavior)
		}
		body := renderAccordionTest(t, view)
		if got := strings.Count(body, ` name="policy-group"`); got != 2 {
			t.Fatalf("%s must name every details item, got %d", behavior, got)
		}
		if !strings.Contains(body, `ui-accordion-icon-chevron`) {
			t.Fatalf("%s must use the reference-style chevron icon", behavior)
		}
	}
	baseUI := prepareAccordion(&accordionView{ID: "baseui", Behavior: accordionBehaviorBaseUI, Multiple: true, MultipleSet: true, Items: []accordionItem{{Open: true}, {Open: true}}})
	if !baseUI.Multiple || baseUI.Name != "" || !baseUI.Items[1].Open {
		t.Fatal("baseui must preserve multiple-open state")
	}
	materialMultiple := prepareAccordion(&accordionView{ID: "material-multiple", Behavior: accordionBehaviorMaterial, Multiple: true, MultipleSet: true, Items: []accordionItem{{Open: true}, {Open: true}}})
	if !materialMultiple.Multiple || materialMultiple.Name != "" || !materialMultiple.Items[1].Open {
		t.Fatal("explicit multiple=true must override the material default")
	}
}

func TestAccordionHTMXRetainsNativeFallbackAndBehaviorIsIndependentOfSkin(t *testing.T) {
	htmx := renderAccordionTest(t, prepareAccordion(&accordionView{ID: "htmx", Behavior: accordionBehaviorMaterial, Execution: accordionExecutionHTMX, Items: []accordionItem{{Value: "one", Open: true}}}))
	for _, marker := range []string{`data-behavior="material"`, `data-execution="htmx"`, `<details`, `<summary`, `open`} {
		if !strings.Contains(htmx, marker) {
			t.Errorf("HTMX mode missing native fallback marker %q", marker)
		}
	}
	nativeSkin := renderAccordionTest(t, prepareAccordion(&accordionView{ID: "skin", Behavior: accordionBehaviorMaterial, Execution: accordionExecutionNative, Items: []accordionItem{{Value: "one"}}}))
	if !strings.Contains(nativeSkin, `data-behavior="material"`) || strings.Contains(nativeSkin, `data-behavior="baseui"`) {
		t.Error("behavior must be explicit and independent from the visual skin")
	}
}
func TestAccordionDocsAndRoute(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/accordion", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, want := range []string{`class="ui-accordion `, `data-value="native"`, `accordion-demo`} {
		if !strings.Contains(body, want) && want != "faq" {
			t.Errorf("accordion route missing %q", want)
		}
	}
}
