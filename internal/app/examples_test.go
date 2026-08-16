package app

import (
	"html/template"
	"strings"
	"testing"

	webassets "geliumui/site/web"
)

// pilotComponentSlugs is the pilot set for the Base UI "Examples + code +
// API reference" pattern: the three component pages that render the
// server-driven "## Examples" and "## API reference" sections. Every other
// component page must stay unchanged (no sections, no headings).
var pilotComponentSlugs = []string{"button", "text-field", "select"}

// pilotPartialInvocations maps each pilot slug to the real Go template
// partial invocations its Examples code blocks must contain. These are the
// actual partials that render the live demos (button, text-field,
// validation-form, select-demo, select-menu-demo) — never invented names.
var pilotPartialInvocations = map[string][]string{
	"button":     {`{{template "button"`},
	"text-field": {`{{template "text-field"`, `{{template "validation-form"`},
	"select":     {`{{template "select-demo"`, `{{template "select-menu-demo"`},
}

// pilotAPIRefProps maps each pilot slug to real fields of its view structs
// that the API reference table must document: buttonView fields for button,
// textFieldView fields for text-field, selectMenuDemo/selectMenuOption
// fields for select.
var pilotAPIRefProps = map[string][]string{
	"button":     {"Label", "Variant", "Href", "IconSVG", "Command", "CommandFor", "Value", "Disabled", "Loading", "Submit", "Autofocus"},
	"text-field": {"ID", "Label", "Name", "Value", "Variant", "Helper", "MessageRole", "Error", "Disabled", "Textarea", "Autofocus"},
	"select":     {"Options", "Error", "Value", "Label", "Selected"},
}

// examplesSection returns the rendered body between the "## Examples" and
// "## API reference" headings ("" when Examples is absent).
func examplesSection(body string) string {
	start := strings.Index(body, ">Examples</h2>")
	if start < 0 {
		return ""
	}
	rest := body[start:]
	if end := strings.Index(rest, ">API reference</h2>"); end >= 0 {
		return rest[:end]
	}
	return rest
}

// TestPilotPagesRenderExamplesSections is the positive pilot contract: each
// pilot page renders a "## Examples" heading with at least 2 example blocks,
// every block carries a live demo (a component-preview rendered by the real
// partial) AND a code block naming that same partial, and the examples
// section sits at the end of the article before the API reference table.
func TestPilotPagesRenderExamplesSections(t *testing.T) {
	for _, slug := range pilotComponentSlugs {
		t.Run(slug, func(t *testing.T) {
			body := getOKBody(t, "/components/"+slug)

			section := examplesSection(body)
			if section == "" {
				t.Fatalf("%s page is missing the '## Examples' section", slug)
			}
			if got := strings.Count(section, `class="example-block"`); got < 2 {
				t.Errorf("%s Examples section has %d example blocks, want at least 2", slug, got)
			}
			// Every example block must pair a live demo (real partial
			// output inside a component-preview) with a code block.
			for _, block := range strings.Split(section, `class="example-block"`)[1:] {
				if !strings.Contains(block, `<div class="component-preview"`) {
					t.Errorf("%s has an example block without a live demo (component-preview)", slug)
				}
				if !strings.Contains(block, "<pre><code>") {
					t.Errorf("%s has an example block without a code block", slug)
				}
			}
			// The code blocks name the real partials that render the demos.
			for _, want := range pilotPartialInvocations[slug] {
				if !strings.Contains(body, want) {
					t.Errorf("%s Examples code blocks are missing the real partial invocation %q", slug, want)
				}
			}
			// Base UI order: Examples and API reference come AFTER the
			// intro but BEFORE the Guidance block — show the component,
			// then the rules (owner decision). Guidance (when present)
			// must appear after both sections.
			ex := strings.Index(body, ">Examples</h2>")
			api := strings.Index(body, ">API reference</h2>")
			if ex < 0 || api < 0 {
				t.Fatalf("%s must render Examples and API reference sections", slug)
			}
			if ex > api {
				t.Errorf("%s must place '## Examples' before '## API reference'", slug)
			}
			if g := strings.Index(body, ">Guidance</h2>"); g >= 0 && g < ex {
				t.Errorf("%s must place '## Guidance' after '## Examples' (show the component, then the rules)", slug)
			}
			if strings.Index(body, "</article>") < api {
				t.Errorf("%s must place '## API reference' inside the article", slug)
			}
		})
	}
}

// TestPilotPagesRenderAPIRefTables is the positive API contract: each pilot
// page renders a "## API reference" heading with a semantic table whose rows
// document the component's real view-struct fields.
func TestPilotPagesRenderAPIRefTables(t *testing.T) {
	for _, slug := range pilotComponentSlugs {
		t.Run(slug, func(t *testing.T) {
			body := getOKBody(t, "/components/"+slug)

			if !strings.Contains(body, ">API reference</h2>") {
				t.Fatalf("%s page is missing the '## API reference' section", slug)
			}
			if !strings.Contains(body, `class="api-ref-table"`) {
				t.Errorf("%s API reference must render a semantic table", slug)
			}
			if !strings.Contains(body, "<th scope=\"col\">Prop</th>") ||
				!strings.Contains(body, "<th scope=\"col\">Type</th>") ||
				!strings.Contains(body, "<th scope=\"col\">Default</th>") ||
				!strings.Contains(body, "<th scope=\"col\">Description</th>") {
				t.Errorf("%s API reference table must carry Prop/Type/Default/Description headers", slug)
			}
			for _, prop := range pilotAPIRefProps[slug] {
				if !strings.Contains(body, "<code>"+prop+"</code>") {
					t.Errorf("%s API reference table is missing the real prop %q", slug, prop)
				}
			}
		})
	}
}

// TestNonPilotPagesOmitExamplesAndAPIRef is the negative pilot contract: the
// pilot sections render on the three pilot pages only. Every other component
// page keeps its current anatomy with no Examples or API reference markup.
func TestNonPilotPagesOmitExamplesAndAPIRef(t *testing.T) {
	for _, slug := range []string{"dialog", "toast", "card"} {
		t.Run(slug, func(t *testing.T) {
			body := getOKBody(t, "/components/"+slug)
			for _, gone := range []string{
				">Examples</h2>",
				">API reference</h2>",
				`class="example-block"`,
				`class="api-ref-table"`,
			} {
				if strings.Contains(body, gone) {
					t.Errorf("%s must not render pilot markup %q", slug, gone)
				}
			}
		})
	}
}

// TestExampleDescriptionsStayUnder25Words extends the copy contract (no
// sentence over 25 words) to the pilot Examples data: example descriptions
// are server-side Go strings, so the length contract is enforced on the
// registry directly.
func TestExampleDescriptionsStayUnder25Words(t *testing.T) {
	for _, slug := range pilotComponentSlugs {
		page, ok := pilotPages[slug]
		if !ok {
			t.Fatalf("pilot registry is missing %s", slug)
		}
		if got := len(page.Examples); got < 2 {
			t.Fatalf("%s must define at least 2 examples, got %d", slug, got)
		}
		for _, ex := range page.Examples {
			if n := len(strings.Fields(ex.Desc)); n > 25 {
				t.Errorf("%s example %q description is %d words (> 25): %q", slug, ex.Name, n, ex.Desc)
			}
		}
	}
}

// TestExampleCodeMatchesRenderedPartial is the no-drift guard: every code
// block shown on a pilot page must invoke the exact partial that renders the
// live demo next to it, and that partial must exist in the template bundle.
// If a partial is renamed or an example's code drifts, this test fails.
func TestExampleCodeMatchesRenderedPartial(t *testing.T) {
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/*.html"))
	for _, slug := range pilotComponentSlugs {
		for _, ex := range pilotPages[slug].Examples {
			if tmpl.Lookup(ex.Partial) == nil {
				t.Errorf("%s example %q names partial %q which is not a defined template", slug, ex.Name, ex.Partial)
			}
			want := `{{template "` + ex.Partial + `"`
			if !strings.Contains(string(ex.Code), want) {
				t.Errorf("%s example %q code %q does not invoke the rendered partial %q", slug, ex.Name, ex.Code, ex.Partial)
			}
		}
	}
}
