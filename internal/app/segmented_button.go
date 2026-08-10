package app

import (
	"net/http"
)

// segmentedButtonDemo is the server-rendered Segmented buttons documentation
// demo. The markup is fully static (closed vocabulary): the three sets below
// demonstrate the roadmap contract for Wave 4B —
//
//   - single-select set: native <input type="radio"> group sharing one name;
//   - multi-select set: native <input type="checkbox"> group sharing one name;
//   - action set: <button type="button"> group inside role="group".
//
// Selection derives from the native :checked pseudo-class (no JavaScript), so
// the same markup works inside a <form method="get"> that submits the checked
// values. There is deliberately no Go-driven per-segment data: the template is
// the reusable, copyable markup, mirroring Checkbox/Switch/List.
type segmentedButtonDemo struct{}

func (s *server) segmentedButtonDocs(w http.ResponseWriter, _ *http.Request) {
	s.renderMarkdownPage(w, pageView{
		Title:               "Segmented buttons",
		SegmentedButtonDemo: &segmentedButtonDemo{},
	}, "content/segmented-button.md")
}
