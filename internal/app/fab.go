package app

import (
	"html/template"
	"net/http"
	"strings"
)

// fabView is the server view model for the Floating Action Button. A FAB is
// always icon-anchored: the icon carries no visible text except in the
// extended variant (Label). The accessible name is mandatory and comes from
// AriaLabel for icon-only FABs or from the visible Label for the extended
// variant. Variant and Size use closed vocabularies; Href switches the root to
// an <a> for navigation, otherwise a native <button type="button"> is used.
type fabView struct {
	AriaLabel  string
	Label      string // extended variant visible text (optional; supplies the name when AriaLabel is empty)
	Variant    string // primary | surface | secondary
	Size       string // small | medium | large
	Lowered    bool
	Href       string
	IconSVG    template.HTML
	Command    string
	CommandFor string
	Value      string
	Disabled   bool
}

// Class returns the closed-vocabulary element classes. The extended form is
// derived from the presence of a visible Label; otherwise Size selects the
// container size. Lowered composes an extra elevation class.
func (v fabView) Class() string {
	classes := []string{"ui-fab", "ui-fab-" + v.Variant}
	if v.Label != "" {
		classes = append(classes, "ui-fab-extended")
	} else {
		classes = append(classes, "ui-fab-"+v.Size)
	}
	if v.Lowered {
		classes = append(classes, "ui-fab-lowered")
	}
	return strings.Join(classes, " ")
}

// fabEditIcon is a trusted, internal decorative FAB icon (Material "edit"
// glyph). It is aria-hidden and unfocusable; the FAB's accessible name comes
// from its aria-label or extended label, never from this icon.
const fabEditIcon template.HTML = `<svg aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25Zm17.71-10.21a1 1 0 0 0 0-1.41l-2.34-2.34a1 1 0 0 0-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83Z"></path></svg>` // #nosec G203 -- trusted, internal decorative glyph.

// fabAddIcon is a trusted, internal decorative FAB icon (Material "add"
// glyph) used where the action implies creation.
const fabAddIcon template.HTML = `<svg aria-hidden="true" focusable="false" viewBox="0 0 24 24" fill="currentColor"><path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2Z"></path></svg>` // #nosec G203 -- trusted, internal decorative glyph.

// fabDemo is the view model for the FAB documentation preview sections.
type fabDemo struct {
	FabDemoViews []fabView
}

// defaultFabDemo enumerates the FAB variants and states the docs page teaches.
// Every icon-only FAB must carry a non-empty AriaLabel so its accessible name
// is never empty; the extended FABs supply visible labels.
func defaultFabDemo() *fabDemo {
	return &fabDemo{
		FabDemoViews: []fabView{
			{AriaLabel: "Edit profile", Variant: "primary", Size: "medium", IconSVG: fabEditIcon},
			{AriaLabel: "Edit profile", Variant: "surface", Size: "medium", IconSVG: fabEditIcon},
			{AriaLabel: "Edit profile", Variant: "secondary", Size: "medium", IconSVG: fabEditIcon},
			{AriaLabel: "Edit profile", Variant: "surface", Size: "small", IconSVG: fabEditIcon},
			{AriaLabel: "Edit profile", Variant: "surface", Size: "large", IconSVG: fabEditIcon},
			{AriaLabel: "Unavailable action", Variant: "surface", Size: "medium", IconSVG: fabEditIcon, Disabled: true},
			{AriaLabel: "Edit profile", Variant: "primary", Size: "medium", IconSVG: fabEditIcon, Lowered: true},
			{AriaLabel: "Compose a new email", Variant: "primary", Size: "medium", Label: "Compose", IconSVG: fabAddIcon},
			{AriaLabel: "Add a new item", Variant: "surface", Size: "medium", Label: "Add", IconSVG: fabAddIcon},
			{AriaLabel: "Navigate to settings", Variant: "surface", Size: "medium", Href: "/components/fab", IconSVG: fabEditIcon},
		},
	}
}

func (s *server) fabDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:   "Floating action button (FAB)",
		FabDemo: defaultFabDemo(),
	}, "content/fab.md")
}
