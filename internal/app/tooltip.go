package app

import (
	"net/http"
)

// tooltipDemo is the server-rendered Tooltip documentation demo. The markup is
// fully static (closed vocabulary) and demonstrates the roadmap contract:
//
//   - tooltips are a styled .ui-tooltip surface revealed with pure CSS on
//     :hover / :focus-within of a .ui-tooltip-host wrapper — no component
//     JavaScript, no reliance on the not-yet-Baseline Interest Invokers API;
//   - the trigger carries aria-describedby pointing at the role="tooltip"
//     surface, so screen readers get the description when the trigger is
//     focused (the roadmap's accessible fallback);
//   - demo content never hides essential information exclusively inside the
//     tooltip: every trigger is self-explanatory on its own (the visible label
//     or aria-label names the action);
//   - plain and rich variants, the rich one with a title, supporting text and
//     an optional action link, positioned below the anchor (above via the
//     ui-tooltip--top modifier).
//
// The reusable trigger classes (.ui-button, .ui-icon-button) come from the
// existing Button and Icon button contracts.
type tooltipDemo struct{}

func (s *server) tooltipDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:       "Tooltip",
		TooltipDemo: &tooltipDemo{},
	}, "content/tooltip.md")
}
