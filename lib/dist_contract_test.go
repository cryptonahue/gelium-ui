package lib

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDistBundleIsCommittedAndComplete pins the S5 dist contract: the
// prebuilt bundle (lib/dist/gelium.css) must exist, be non-trivial, and carry
// the component primitives an unthemed consumer relies on.
func TestDistBundleIsCommittedAndComplete(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(repositoryRoot(t), "lib", "dist", "gelium.css"))
	if err != nil {
		t.Fatalf("read lib/dist/gelium.css: %v", err)
	}
	if len(b) < 50_000 {
		t.Errorf("dist bundle suspiciously small (%d bytes); expected the full component set", len(b))
	}
	for _, contract := range []string{`.ui-button`, `.ui-toast`, `.ui-text-field`, `--ui-touch-target`} {
		if !strings.Contains(string(b), contract) {
			t.Errorf("dist bundle missing %q", contract)
		}
	}
	// The dist bundle MUST carry the theme roots: themes ship WITH the
	// package (lib/themes/*.css) so a consumer installing gelium-ui gets a
	// working theme pair, not bare components.
	if !strings.Contains(string(b), ".theme-material") || !strings.Contains(string(b), ".theme-basecoat") {
		t.Error("dist bundle must embed both theme roots (.theme-material, .theme-basecoat)")
	}
}

// TestLibJsShipsConsumerEnhancements pins the S5 JS contract: lib/js/gelium.js
// carries the consumer enhancements (toast region, 422 contract, VT flag,
// slider fill) and the docs chrome (theme/scheme optimistic toggle) does NOT.
func TestLibJsShipsConsumerEnhancements(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(repositoryRoot(t), "lib", "js", "gelium.js"))
	if err != nil {
		t.Fatalf("read lib/js/gelium.js: %v", err)
	}
	js := string(b)
	for _, contract := range []string{`gelium:toast`, `X-Gelium-Validation`, `startViewTransition`, `--ui-slider-fill`} {
		if !strings.Contains(js, contract) {
			t.Errorf("gelium.js missing consumer enhancement %q", contract)
		}
	}
	if strings.Contains(js, "refreshChromeHrefs") {
		t.Error("gelium.js must not contain docs chrome (refreshChromeHrefs is site-only)")
	}
}
