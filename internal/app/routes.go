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
		{Path: "/components/accordion", Label: "Accordion", Handler: (*server).accordionDocs},
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
		// Composition & content primitives (documented from lib/templates):
		{Path: "/components/hero", Label: "Hero", Handler: (*server).heroDocs},
		{Path: "/components/avatar", Label: "Avatar", Handler: (*server).avatarDocs},
		{Path: "/components/breadcrumb", Label: "Breadcrumb", Handler: (*server).breadcrumbDocs},
		{Path: "/components/footer", Label: "Footer", Handler: (*server).footerDocs},
		{Path: "/components/pagination", Label: "Pagination", Handler: (*server).paginationDocs},
		{Path: "/components/section-heading", Label: "Section heading", Handler: (*server).sectionHeadingDocs},
		{Path: "/components/feature-card", Label: "Feature card", Handler: (*server).featureCardDocs},
		{Path: "/components/split", Label: "Split", Handler: (*server).splitDocs},
		{Path: "/components/image", Label: "Image", Handler: (*server).imageDocs},
		{Path: "/components/media", Label: "Media", Handler: (*server).mediaDocs},
		{Path: "/components/video", Label: "Video", Handler: (*server).videoDocs},
		{Path: "/components/newsletter", Label: "Newsletter", Handler: (*server).newsletterDocs},
		{Path: "/components/language-switcher", Label: "Language switcher", Handler: (*server).languageSwitcherDocs},
		// Feedback & status primitives (the state patterns handbook pages name):
		{Path: "/components/banner", Label: "Banner", Handler: (*server).bannerDocs},
		{Path: "/components/inline-alert", Label: "Inline alert", Handler: (*server).inlineAlertDocs},
		{Path: "/components/callout", Label: "Callout", Handler: (*server).calloutDocs},
		{Path: "/components/skeleton", Label: "Skeleton", Handler: (*server).skeletonDocs},
		{Path: "/components/empty-state", Label: "Empty state", Handler: (*server).emptyStateDocs},
		{Path: "/components/error-state", Label: "Error state", Handler: (*server).errorStateDocs},
		{Path: "/components/validation-summary", Label: "Validation summary", Handler: (*server).validationSummaryDocs},
	}
}

// navLinks is the flat header/legacy nav: Docs hub plus every component link
// derived from docsSections only so registration order in componentRoutes
// cannot drift from the docs IA categories.
func navLinks() []navLink {
	n := 1
	for _, section := range docsSections {
		n += len(section.Links)
	}
	links := make([]navLink, 0, n)
	links = append(links, navLink{Path: "/docs", Label: "Docs"})
	for _, section := range docsSections {
		links = append(links, section.Links...)
	}
	return links
}

type navLink struct {
	Path  string
	Label string
}
