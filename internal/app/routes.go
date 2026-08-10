package app

import (
	"net/http"
)

// componentRoute is the single integration point for a component docs page.
// Label is the navigation label; Path is the GET route.
type componentRoute struct {
	Path    string
	Label   string
	Handler func(*server, http.ResponseWriter, *http.Request)
}

// componentRoutes is ordered: navigation order and route registration order.
func componentRoutes() []componentRoute {
	return []componentRoute{
		{Path: "/components/button", Label: "Button", Handler: (*server).buttonDocs},
		{Path: "/components/text-field", Label: "Text field", Handler: (*server).textFieldDocs},
		{Path: "/components/dialog", Label: "Dialog", Handler: (*server).dialogDocs},
		{Path: "/components/toast", Label: "Toast", Handler: (*server).toastDocs},
		{Path: "/components/elevation", Label: "Elevation", Handler: (*server).elevationDocs},
		{Path: "/components/focus-ring", Label: "Focus ring", Handler: (*server).focusRingDocs},
		{Path: "/components/icon", Label: "Icon", Handler: (*server).iconDocs},
		{Path: "/components/divider", Label: "Divider", Handler: (*server).dividerDocs},
		{Path: "/components/card", Label: "Card", Handler: (*server).cardDocs},
		{Path: "/components/badge", Label: "Badge", Handler: (*server).badgeDocs},
		{Path: "/components/checkbox", Label: "Checkbox", Handler: (*server).checkboxDocs},
		{Path: "/components/radio", Label: "Radio", Handler: (*server).radioDocs},
		{Path: "/components/switch", Label: "Switch", Handler: (*server).switchDocs},
		{Path: "/components/select", Label: "Select", Handler: (*server).selectDocs},
		{Path: "/components/slider", Label: "Slider", Handler: (*server).sliderDocs},
		{Path: "/components/progress", Label: "Progress", Handler: (*server).progressDocs},
		{Path: "/components/icon-button", Label: "Icon button", Handler: (*server).iconButtonDocs},
		{Path: "/components/fab", Label: "FAB", Handler: (*server).fabDocs},
		{Path: "/components/list", Label: "List", Handler: (*server).listDocs},
		{Path: "/components/chips", Label: "Chips", Handler: (*server).chipsDocs},
		{Path: "/components/tabs", Label: "Tabs", Handler: (*server).tabsDocs},
		{Path: "/components/navigation-bar", Label: "Navigation bar", Handler: (*server).navigationBarDocs},
		{Path: "/components/navigation-tab", Label: "Navigation tab", Handler: (*server).navigationTabDocs},
		{Path: "/components/segmented-button", Label: "Segmented buttons", Handler: (*server).segmentedButtonDocs},
		{Path: "/components/menu", Label: "Menu", Handler: (*server).menuDocs},
		{Path: "/components/navigation-drawer", Label: "Navigation drawer", Handler: (*server).navigationDrawerDocs},
		{Path: "/components/data-table", Label: "Data table", Handler: (*server).dataTableDocs},
		{Path: "/components/tooltip", Label: "Tooltip", Handler: (*server).tooltipDocs},
	}
}

func navLinks() []navLink {
	links := make([]navLink, 0, len(componentRoutes()))
	for _, r := range componentRoutes() {
		links = append(links, navLink{Path: r.Path, Label: r.Label})
	}
	return links
}

type navLink struct {
	Path  string
	Label string
}
