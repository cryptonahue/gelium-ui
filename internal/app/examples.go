package app

import (
	"bytes"
	"html/template"
)

// exampleBlock is one demo+code pair in a component page's "## Examples"
// section (Base UI pattern pilot). Demo is the live render of the real
// component partial; Code is the actual Go template invocation that produced
// it — both derive from the same Partial name, so the shown code can never
// drift from the demo. Views holds one view per partial execution.
type exampleBlock struct {
	Name    string        // "### " heading
	Desc    string        // one-line description of what the example shows
	Code    template.HTML // the real partial invocation, shown in the code block
	Partial string        // partial executed for the live demo (same name Code names)
	Views   []any         // one view per partial execution
	Demo    template.HTML // live demo markup rendered from Partial+Views
}

// apiRefRow is one row of the "## API reference" table: a real field of the
// component's view structs, its Go type, its zero-value default, and what it
// does. The registry is hand-maintained but verified against the actual
// structs by the pilot tests, so a renamed or dropped field fails the suite.
type apiRefRow struct {
	Prop        string
	Type        string
	Default     string
	Description string
}

// apiRefView is the "## API reference" section of one pilot component: a
// semantic table sourced from the component's real view structs.
type apiRefView struct {
	Component string
	Rows      []apiRefRow
}

// pilotPage bundles the Examples + API reference sections of one pilot page.
type pilotPage struct {
	Examples []exampleBlock
	APIRef   *apiRefView
}

// pilotPages is the pilot registry for the Base UI "Examples + code + API
// reference" pattern. Only the three pilot components (button, text-field,
// select) render the sections; every other page leaves Examples/APIRef nil
// and the layout emits nothing. The example views reuse the real component
// view structs (buttonView, textFieldView, validationFormView, selectDemo,
// selectMenuDemo) and the Code strings name the real partials — the same
// names passed to ExecuteTemplate when the demos are rendered.
var pilotPages = map[string]pilotPage{
	"button": {
		Examples: []exampleBlock{
			{
				Name:    "Variants",
				Desc:    "Primary, secondary, and outline buttons carry the page's action hierarchy.",
				Code:    template.HTML(`{{template "button" .}}`),
				Partial: "button",
				Views: []any{
					buttonView{Label: "Save changes", Variant: "primary", IconSVG: saveIconSVG},
					buttonView{Label: "Continue", Variant: "secondary"},
					buttonView{Label: "Learn more", Variant: "outline"},
				},
			},
			{
				Name:    "Disabled and loading",
				Desc:    "Disabled and loading buttons stay native disabled controls and never rely on color alone.",
				Code:    template.HTML(`{{template "button" .}}`),
				Partial: "button",
				Views: []any{
					buttonView{Label: "Unavailable", Variant: "primary", Disabled: true},
					buttonView{Label: "Save changes", Variant: "primary", Loading: true},
				},
			},
			{
				Name:    "Link button",
				Desc:    "Passing Href renders the button as a real link for navigation.",
				Code:    template.HTML(`{{template "button" .}}`),
				Partial: "button",
				Views: []any{
					buttonView{Label: "Open the docs", Variant: "outline", Href: "/docs"},
				},
			},
		},
		APIRef: &apiRefView{
			Component: "Button",
			Rows: []apiRefRow{
				{Prop: "Label", Type: "string", Default: `""`, Description: "Accessible action text rendered inside the button."},
				{Prop: "Variant", Type: "string", Default: `""`, Description: "Emphasis level: primary, secondary, outline, or text."},
				{Prop: "Href", Type: "string", Default: `""`, Description: "When set, renders a link-shaped anchor instead of a button."},
				{Prop: "IconSVG", Type: "template.HTML", Default: `""`, Description: "Decorative inline SVG; must be aria-hidden and focusable=false."},
				{Prop: "Command", Type: "string", Default: `""`, Description: "Invoker Command name for the button."},
				{Prop: "CommandFor", Type: "string", Default: `""`, Description: "Invoker Command target element id."},
				{Prop: "Value", Type: "string", Default: `""`, Description: "Submitted value when the button is inside a form."},
				{Prop: "Disabled", Type: "bool", Default: "false", Description: "Native disabled attribute; leaves the tab order."},
				{Prop: "Loading", Type: "bool", Default: "false", Description: "Native disabled plus aria-busy and the accessible name Loading {Label}."},
				{Prop: "Submit", Type: "bool", Default: "false", Description: "Renders type=submit for form submission."},
				{Prop: "Autofocus", Type: "bool", Default: "false", Description: "Native autofocus attribute."},
			},
		},
	},
	"text-field": {
		Examples: []exampleBlock{
			{
				Name:    "Outlined and filled",
				Desc:    "Outlined and filled fields cover the two visual variants of the component.",
				Code:    template.HTML(`{{template "text-field" .}}`),
				Partial: "text-field",
				Views: []any{
					textFieldView{ID: "ex-name", Label: "Name", Variant: "outlined", Required: true, MaxLength: 40},
					textFieldView{ID: "ex-email", Label: "Email", Variant: "filled", Type: "email", Autocomplete: "email", InputMode: "email"},
				},
			},
			{
				Name:    "Native input surface",
				Desc:    "Type, constraints, and autofill hints pass through as native attributes.",
				Code:    template.HTML(`{{template "text-field" .}}`),
				Partial: "text-field",
				Views: []any{
					textFieldView{ID: "ex-surface", Label: "Email address", Variant: "filled", Type: "email", Autocomplete: "email", InputMode: "email", Required: true, MinLength: 5, MaxLength: 254, Helper: "We check the format on submit."},
					textFieldView{ID: "ex-surface-search", Label: "Search", Variant: "outlined", Type: "search", InputMode: "search", Placeholder: "Search posts…"},
				},
			},
			{
				Name:    "Helper and error messages",
				Desc:    "Helper text and visible errors are associated through aria-describedby, never color alone.",
				Code:    template.HTML(`{{template "text-field" .}}`),
				Partial: "text-field",
				Views: []any{
					textFieldView{ID: "ex-helper", Label: "Email", Variant: "filled", Helper: "We'll only use this for account updates."},
					textFieldView{ID: "ex-error", Label: "Username", Variant: "outlined", Value: "?", Error: "Use letters and numbers only."},
				},
			},
			{
				Name:    "Server-side validation",
				Desc:    "The validation form posts to the server and swaps the fragment on HTTP 422.",
				Code:    template.HTML(`{{template "validation-form" .}}`),
				Partial: "validation-form",
				// The demo uses a distinct field id so the docs page never
				// carries two elements with id="validation-name" (the live
				// preview below already owns that id). The partial is the
				// same one the code block names.
				Views: []any{validationFormView{
					Field:  textFieldView{ID: "example-validation-name", Label: "Name", Name: "name", Variant: "outlined", Helper: "Enter your name."},
					Submit: buttonView{Label: "Validate name", Variant: "primary", Submit: true},
				}},
			},
		},
		APIRef: &apiRefView{
			Component: "Text field",
			Rows: []apiRefRow{
				{Prop: "ID", Type: "string", Default: `""`, Description: "Element id paired with the label's for attribute."},
				{Prop: "Label", Type: "string", Default: `""`, Description: "Visible label text, explicitly associated via for and id."},
				{Prop: "Name", Type: "string", Default: `""`, Description: "Form submission name."},
				{Prop: "Value", Type: "string", Default: `""`, Description: "Initial input or textarea value."},
				{Prop: "Variant", Type: "string", Default: `""`, Description: "Field style: outlined or filled."},
				{Prop: "Helper", Type: "string", Default: `""`, Description: "Helper message referenced by aria-describedby."},
				{Prop: "MessageRole", Type: "string", Default: `""`, Description: "ARIA role for the helper message, for example status."},
				{Prop: "Error", Type: "string", Default: `""`, Description: "Error message; renders aria-invalid and a visible alert."},
				{Prop: "Disabled", Type: "bool", Default: "false", Description: "Native disabled attribute; takes precedence over error."},
				{Prop: "Textarea", Type: "bool", Default: "false", Description: "Renders a multi-line textarea instead of an input."},
				{Prop: "Autofocus", Type: "bool", Default: "false", Description: "Native autofocus attribute."},
				{Prop: "Type", Type: "string", Default: `"text"`, Description: "Native input type: email, url, tel, number, password, search, or text. Inputs only."},
				{Prop: "Required", Type: "bool", Default: "false", Description: "Native required attribute; the value must be submitted."},
				{Prop: "MaxLength", Type: "int", Default: "0", Description: "Native maxlength in characters; 0 omits the attribute."},
				{Prop: "MinLength", Type: "int", Default: "0", Description: "Native minlength in characters; 0 omits the attribute."},
				{Prop: "Pattern", Type: "string", Default: `""`, Description: "Native pattern regex for inputs; textareas do not support it."},
				{Prop: "Autocomplete", Type: "string", Default: `""`, Description: "Standard autocomplete token (email, current-password, …)."},
				{Prop: "InputMode", Type: "string", Default: `""`, Description: "Native inputmode hint for the on-screen keyboard."},
				{Prop: "ReadOnly", Type: "bool", Default: "false", Description: "Native readonly attribute; focusable and submitted, not editable."},
				{Prop: "Placeholder", Type: "string", Default: `""`, Description: "Visible hint text; keeps the floating label floated and never replaces it."},
			},
		},
	},
	"select": {
		Examples: []exampleBlock{
			{
				Name:    "Filled and outlined fields",
				Desc:    "The two field variants share the native select as the focusable control.",
				Code:    template.HTML(`{{template "select-demo" .}}`),
				Partial: "select-demo",
				Views:   []any{&selectDemo{}},
			},
			{
				Name:    "Server-driven menu",
				Desc:    "The menu posts its value and the server validates it against a closed vocabulary.",
				Code:    template.HTML(`{{template "select-menu-demo" .}}`),
				Partial: "select-menu-demo",
				Views:   []any{defaultSelectMenuDemo()},
			},
		},
		APIRef: &apiRefView{
			Component: "Select",
			Rows: []apiRefRow{
				{Prop: "Options", Type: "[]selectMenuOption", Default: "standard, priority, enterprise", Description: "Options rendered as native option elements; the closed vocabulary of the server-driven demo."},
				{Prop: "Error", Type: "string", Default: `""`, Description: "Server-side validation message; renders aria-invalid when set."},
				{Prop: "Value", Type: "string", Default: `""`, Description: "Option value submitted by the native select; unknown values return HTTP 422."},
				{Prop: "Label", Type: "string", Default: `""`, Description: "Visible option text in the picker."},
				{Prop: "Selected", Type: "bool", Default: "false", Description: "Marks the server-side default selection."},
			},
		},
	},
}

// examplesFor returns the pilot Examples + API reference sections for a
// component slug. Non-pilot slugs return nil values so the layout renders
// nothing (no sections, no headings).
func examplesFor(slug string) ([]exampleBlock, *apiRefView) {
	p, ok := pilotPages[slug]
	if !ok {
		return nil, nil
	}
	return p.Examples, p.APIRef
}

// renderExampleDemos executes each example's live demo with the real
// component partial — the same partial its Code snippet names — so the shown
// code is the actual invocation that produced the demo. It returns a copy so
// the package-level registry is never mutated by request handling.
func (s *server) renderExampleDemos(examples []exampleBlock) ([]exampleBlock, error) {
	out := make([]exampleBlock, len(examples))
	copy(out, examples)
	for i := range out {
		var buf bytes.Buffer
		for _, view := range out[i].Views {
			if err := s.templates.ExecuteTemplate(&buf, out[i].Partial, view); err != nil {
				return nil, err
			}
		}
		out[i].Demo = template.HTML(buf.String()) // #nosec G203 -- rendered from trusted embedded partials.
	}
	return out, nil
}
