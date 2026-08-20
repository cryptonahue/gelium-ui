package app

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type accordionItem struct {
	Value     string
	Heading   string
	Body      string
	Open      bool
	TriggerID string
	PanelID   string
}

// accordionBehavior describes interaction rules independently of the visual
// theme (skin). All profiles in this slice use native details/summary HTML.
type accordionBehavior string

const (
	accordionBehaviorNative   accordionBehavior = "native"
	accordionBehaviorBasecoat accordionBehavior = "basecoat"
	accordionBehaviorMaterial accordionBehavior = "material"
	accordionBehaviorBaseUI   accordionBehavior = "baseui"
)

type accordionExecution string

const (
	accordionExecutionNative accordionExecution = "native"
	accordionExecutionHTMX   accordionExecution = "htmx"
)

func normalizeAccordionBehavior(raw string) accordionBehavior {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case string(accordionBehaviorBasecoat):
		return accordionBehaviorBasecoat
	case string(accordionBehaviorMaterial):
		return accordionBehaviorMaterial
	case string(accordionBehaviorBaseUI):
		return accordionBehaviorBaseUI
	default:
		return accordionBehaviorNative
	}
}

func normalizeAccordionExecution(raw string) accordionExecution {
	if strings.EqualFold(strings.TrimSpace(raw), string(accordionExecutionHTMX)) {
		return accordionExecutionHTMX
	}
	return accordionExecutionNative
}

type accordionView struct {
	ID string
	// Heading is an optional visible name for the accordion region. When it is
	// empty, Label is retained as the backwards-compatible aria-label fallback.
	Heading  string
	Label    string
	Items    []accordionItem
	Multiple bool
	// MultipleSet distinguishes an explicit false from the zero value. This
	// lets behavior profiles choose their own default without taking control
	// away from callers that explicitly selected multiple-open mode.
	MultipleSet bool
	// Name is emitted on details only for the exclusive (Multiple=false) mode.
	// The browser's named-details support is progressive, not a JS polyfill.
	Name      string
	Behavior  accordionBehavior
	Execution accordionExecution
}

var accordionIDPartRE = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

func accordionIDPart(value string) string {
	value = strings.TrimSpace(value)
	value = accordionIDPartRE.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-_")
	if value == "" {
		return "item"
	}
	return strings.ToLower(value)
}

// prepareAccordion assigns deterministic IDs from the caller's stable Value.
// Values remain data attributes exactly as supplied; IDs are normalized only
// to satisfy HTML identifier conventions and avoid unsafe attribute content.
// It also normalizes the typed profiles so templates never emit arbitrary
// behavior or execution values.
func prepareAccordion(view *accordionView) *accordionView {
	if view == nil {
		return view
	}
	view.Behavior = normalizeAccordionBehavior(string(view.Behavior))
	view.Execution = normalizeAccordionExecution(string(view.Execution))
	root := view.ID
	if root == "" {
		root = "accordion"
		view.ID = root
	}
	// Basecoat and Material are exclusive by default, using the native named
	// details group. A true Multiple value remains backwards-compatible as an
	// explicit opt-in; MultipleSet is needed when explicitly selecting false.
	if view.Behavior == accordionBehaviorBasecoat || view.Behavior == accordionBehaviorMaterial {
		if !view.Multiple {
			if view.Name == "" {
				view.Name = accordionIDPart(root) + "-group"
			}
			keepOpen := true
			for i := range view.Items {
				if view.Items[i].Open {
					if keepOpen {
						keepOpen = false
					} else {
						view.Items[i].Open = false
					}
				}
			}
		}
	}
	used := map[string]int{}
	for i := range view.Items {
		item := &view.Items[i]
		part := accordionIDPart(item.Value)
		used[part]++
		if used[part] > 1 {
			part += "-" + strconv.Itoa(used[part])
		}
		// The index suffix is only a collision disambiguator; normal items are
		// keyed by Value so IDs remain stable when unrelated items are added.
		item.TriggerID = root + "-trigger-" + part
		item.PanelID = root + "-panel-" + part
	}
	return view
}

type accordionBehaviorContextKey struct{}
type accordionExecutionContextKey struct{}

func accordionBehaviorFromRequest(r *http.Request) accordionBehavior {
	if v, ok := r.Context().Value(accordionBehaviorContextKey{}).(accordionBehavior); ok {
		return v
	}
	return accordionBehaviorNative
}

func accordionExecutionFromRequest(r *http.Request) accordionExecution {
	if v, ok := r.Context().Value(accordionExecutionContextKey{}).(accordionExecution); ok {
		return v
	}
	return accordionExecutionNative
}

func accordionForRequest(r *http.Request) *accordionView {
	view := defaultAccordion()
	view.Behavior = accordionBehaviorFromRequest(r)
	view.Execution = accordionExecutionFromRequest(r)
	// The built-in demo historically defaults to multiple-open. Profiles on
	// the request path may resolve that demo default independently; explicitly
	// constructed accordionView values remain authoritative.
	if view.Behavior == accordionBehaviorBasecoat || view.Behavior == accordionBehaviorMaterial {
		view.Multiple = false
		view.MultipleSet = false
	}
	return prepareAccordion(view)
}

func defaultAccordion() *accordionView {
	return prepareAccordion(&accordionView{
		ID: "accordion-demo", Heading: "Frequently asked questions", Label: "Frequently asked questions", Multiple: true,
		Items: []accordionItem{
			{Value: "native", Heading: "Is this native HTML?", Body: "Yes. Each item is a details/summary disclosure and remains usable when JavaScript is disabled.", Open: true},
			{Value: "themes", Heading: "Does it follow the active theme?", Body: "The same semantic markup is styled with the active theme's accordion and focus tokens."},
			{Value: "htmx", Heading: "Where does HTMX fit?", Body: "HTMX may progressively enhance a server-rendered panel; it is not required for opening or closing items."},
		},
	})
}

func (s *server) accordionDocs(w http.ResponseWriter, r *http.Request) {
	accordion := accordionForRequest(r)
	// Basecoat's official example is a flat disclosure list without an inner
	// card heading; the page H1 already names the component. Keep the optional
	// heading available for gallery/embedded compositions.
	if accordion.Behavior == accordionBehaviorBasecoat {
		accordion.Heading = ""
	}
	s.renderMarkdownPage(w, r, pageView{Title: "Accordion", AccordionDemo: accordion}, "content/accordion.md")
}
