package web

import (
	"strings"
	"testing"
)

// TestOnThisPageStylesInBundle proves the On-this-page rail and previous/next
// pagination styles compile into the single served CSS bundle. The rail is
// hidden on narrow screens and becomes a sticky third pane at >= 64rem.
func TestOnThisPageStylesInBundle(t *testing.T) {
	compiled := compactCSS(t, compiledAppCSS(t))
	for _, contract := range []string{
		".docs-on-this-page{display:none}",
		".docs-on-this-page-link",
		".docs-prev-next{",
		".docs-prev-next-direction{",
		// compactCSS strips whitespace: margin:0 auto compacts to margin:0auto.
		".docs-shell-content{width:min(65ch,100%);margin:0auto}",
	} {
		if !strings.Contains(compiled, contract) {
			t.Errorf("compiled bundle is missing On-this-page contract %q", contract)
		}
	}
}
