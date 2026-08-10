package app

import (
	"net/http"
)

// menuDemo is the server-rendered Menu documentation demo. The markup is fully
// static (closed vocabulary) and demonstrates the Wave 5 roadmap contract:
//
//   - a menu surface opened declaratively with the native popover attribute;
//   - the trigger is a real <button popovertarget> — no component JavaScript;
//   - items are real navigation links (<a href>), real action buttons
//     (<button type="button">) or native checkbox/radio selection rows inside
//     a real <form>, so the no-JS flow is genuine navigation/submission;
//   - disabled items, optional leading icons, a divider and the Material menu
//     anatomy (48px items, container radius, state layers, focus ring).
//
// The menu surface itself is a <ul> so the item rows keep native list
// semantics; selection rows reuse the native checkbox/radio pattern from the
// List and Segmented button components. There is deliberately no Go-driven
// per-item data: the template is the reusable, copyable markup.
type menuDemo struct{}

func (s *server) menuDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:    "Menu",
		MenuDemo: &menuDemo{},
	}, "content/menu.md")
}
