package web

import (
	"strings"
	"testing"
)

// Priority components migrated onto the Accordion-style reference/skin cascade.
var cascadeComponents = []string{
	"button", "field", "select", "checkbox", "radio", "switch",
	"dialog", "tabs", "card", "badge", "chip",
}

func TestComponentReferenceAndSkinAdaptersSurviveCompilation(t *testing.T) {
	compiled := compactCSS(t, compiledAppCSS(t))
	ref := strings.Index(compiled, `html[data-gelium-reference=material]{`)
	skin := strings.Index(compiled, `html[data-gelium-skin=basecoat]{`)
	if ref < 0 || skin < 0 || skin <= ref {
		t.Fatalf("component adapters must compile with skin after reference: ref=%d skin=%d", ref, skin)
	}

	// Material reference supplies button anatomy identity before Basecoat skin.
	// Accordion and component adapters may emit separate same-selector blocks;
	// search the full reference→skin window and the full skin suffix.
	if !strings.Contains(compiled[ref:skin], "--ui-button-cascade-id:material-reference") {
		t.Error("material reference must supply button cascade identity")
	}
	if !strings.Contains(compiled[skin:], "--ui-button-cascade-id:basecoat-skin") {
		t.Error("basecoat skin must override button cascade identity")
	}
	if !strings.Contains(compiled[ref:skin], "--ui-button-radius:var(--ui-radius-full)") {
		t.Error("material reference must keep pill button radius")
	}
	if !strings.Contains(compiled[skin:], "--ui-button-radius:.5rem") {
		t.Error("basecoat skin must own denser button radius")
	}
}

func TestEachProductSkinOwnsPriorityComponentAnatomy(t *testing.T) {
	source := compactCSS(t, repositoryFile(t, "site", "web", "styles", "component-skin.css")) +
		compactCSS(t, repositoryFile(t, "site", "web", "styles", "basecoat-pack-skins.css"))
	for _, skin := range []string{
		"material", "basecoat", "basecoat-nova", "basecoat-maia", "basecoat-lyra",
		"basecoat-mira", "basecoat-luma", "basecoat-sera", "basecoat-rhea",
		"baseui", "alden", "linear", "vercel",
	} {
		selector := `html[data-gelium-skin="` + skin + `"]{`
		start := strings.Index(source, selector)
		if start < 0 {
			t.Errorf("skin %q missing component-skin selector", skin)
			continue
		}
		// Pack skins may split shared tokens across multiple selectors; search
		// a window large enough to cover the dedicated pack block.
		window := source[start:]
		if len(window) > 4500 {
			window = window[:4500]
		}
		for _, comp := range cascadeComponents {
			want := "--ui-" + comp + "-cascade-id:" + skin + "-skin"
			if skin == "basecoat" {
				// Vega default keeps the historical cascade id.
				want = "--ui-" + comp + "-cascade-id:basecoat-skin"
			}
			if !strings.Contains(window, want) && !strings.Contains(source, want) {
				t.Errorf("skin %q must own %s anatomy identity %q", skin, comp, want)
			}
		}
		for _, token := range []string{
			"--ui-button-radius:",
			"--ui-button-padding-y:",
			"--ui-field-radius:",
			"--ui-field-padding-x:",
			"--ui-select-radius:",
			"--ui-checkbox-size:",
			"--ui-radio-size:",
			"--ui-switch-width:",
			"--ui-dialog-radius:",
			"--ui-tabs-height:",
			"--ui-card-radius:",
			"--ui-card-padding:",
			"--ui-badge-size:",
			"--ui-chip-height:",
			"--ui-icon-button-radius:",
		} {
			// Shared Basecoat-family block supplies many tokens once for packs.
			if !strings.Contains(window, token) && !strings.Contains(source, token) {
				t.Errorf("skin %q missing anatomy token %s", skin, token)
			}
		}
	}
}

func TestReferencePresetsDefinePriorityComponentAnatomy(t *testing.T) {
	source := compactCSS(t, repositoryFile(t, "site", "web", "styles", "component-reference.css"))
	for _, ref := range []string{"material", "basecoat", "baseui"} {
		selector := `html[data-gelium-reference="` + ref + `"]{`
		start := strings.Index(source, selector)
		if start < 0 {
			t.Errorf("reference %q missing component-reference selector", ref)
			continue
		}
		end := strings.Index(source[start:], "}")
		if end < 0 {
			t.Errorf("reference %q selector not closed", ref)
			continue
		}
		block := source[start : start+end]
		for _, comp := range cascadeComponents {
			want := "--ui-" + comp + "-cascade-id:" + ref + "-reference"
			if !strings.Contains(block, want) {
				t.Errorf("reference %q must define %s identity %q", ref, comp, want)
			}
		}
	}
}

func TestAldenSkinButtonTopologyDoesNotInheritBasecoatReference(t *testing.T) {
	source := compactCSS(t, repositoryFile(t, "site", "web", "styles", "component-skin.css"))
	selector := `html[data-gelium-skin="alden"]{`
	start := strings.Index(source, selector)
	if start < 0 {
		t.Fatal("Alden component skin selector missing")
	}
	end := strings.Index(source[start:], "}")
	if end < 0 {
		t.Fatal("Alden component skin selector not closed")
	}
	block := source[start : start+end]
	for _, want := range []string{
		"--ui-button-cascade-id:alden-skin",
		"--ui-button-radius:12px",
		"--ui-field-radius:8px",
		"--ui-card-padding:1.5rem",
		"--ui-dialog-radius:16px",
	} {
		if !strings.Contains(block, want) {
			t.Errorf("Alden skin must own distinctive anatomy %q", want)
		}
	}
	if strings.Contains(block, "--ui-button-radius:.375rem") {
		t.Error("Alden skin must not inherit Basecoat button radius")
	}

	compiled := compactCSS(t, compiledAppCSS(t))
	if !strings.Contains(compiled, "--ui-button-cascade-id:alden-skin") {
		t.Error("compiled bundle must retain Alden button cascade identity")
	}
}

func TestPriorityComponentCoreAvoidsBehaviorVisualSelectors(t *testing.T) {
	files := []string{
		"button.css", "icon-button.css", "text-field.css", "select.css",
		"checkbox.css", "radio.css", "switch.css", "dialog.css", "tabs.css",
		"card.css", "badge.css", "chips.css",
	}
	for _, file := range files {
		css := sourceComponentCSS(t, file)
		if strings.Contains(css, "--behavior-") || strings.Contains(css, "[data-behavior") {
			t.Errorf("%s must not select behavior for visuals", file)
		}
		// No component-root custom property declarations that would shadow
		// html[data-gelium-*] anatomy adapters for the migrated families.
		for _, forbidden := range []string{
			"--ui-dialog-min-width:",
			"--ui-tabs-height:",
			"--ui-chip-height:",
			"--ui-button-padding-y:",
			"--ui-icon-button-radius:",
		} {
			if strings.Contains(css, "  "+forbidden) || strings.Contains(css, "\t"+forbidden) {
				t.Errorf("%s must not declare shadowing anatomy token %s on the component root", file, forbidden)
			}
		}
	}
}

func TestVercelSkinWinsOverBaseUIReferenceForButtonAnatomy(t *testing.T) {
	compiled := compactCSS(t, compiledAppCSS(t))
	ref := strings.Index(compiled, `html[data-gelium-reference=baseui]{`)
	skin := strings.Index(compiled, `html[data-gelium-skin=vercel]{`)
	if ref < 0 || skin <= ref {
		t.Fatalf("vercel skin must follow baseui reference: ref=%d skin=%d", ref, skin)
	}
	if !strings.Contains(compiled[ref:skin], "--ui-button-cascade-id:baseui-reference") {
		t.Error("baseui reference must identify button anatomy")
	}
	if !strings.Contains(compiled[skin:], "--ui-button-cascade-id:vercel-skin") {
		t.Error("vercel skin must override button cascade identity")
	}
	if !strings.Contains(compiled[skin:], "--ui-dialog-radius:12px") {
		t.Error("vercel skin must own dialog radius")
	}
}
