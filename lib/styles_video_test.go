package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestVideoContainerCSSMapsTokens(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "video.css"), " ")

	for _, contract := range []string{
		`.ui-video {`,
		`width: 100%;`,
		`aspect-ratio: 16 / 9;`,
		`overflow: hidden;`,
		`border: 1px solid var(--ui-video-border);`,
		`border-radius: var(--ui-video-radius);`,
		`.ui-video video {`,
		`.ui-video--aspect-4-3 {`,
		`aspect-ratio: 4 / 3;`,
		`.ui-video-fallback {`,
		`--ui-video-radius: var(--ui-radius-sm);`,
		`--ui-video-border: var(--ui-color-border);`,
		`--ui-video-bg: var(--ui-color-surface-container);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source video CSS is missing contract %q", contract)
		}
	}

	// aspect-ratio is structural geometry and must stay literal — it is never
	// tokenized (roadmap rule, like breakpoints and z-index).
	if strings.Contains(css, `aspect-ratio: var(--ui-`) {
		t.Error("video.css must not tokenize aspect-ratio (structural geometry, roadmap rule)")
	}
	if strings.Contains(css, `transition:`) || strings.Contains(css, `animation:`) {
		t.Error("video.css must not declare transitions or animations")
	}
}

func TestVideoContractCSSWired(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-video`,
		`.ui-video--aspect-4-3`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled video CSS is missing %q", contract)
		}
	}
}

func TestVideoClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "lib", "templates", "video.html")

	for _, contract := range []string{
		`class="ui-video`,
		`ui-video--aspect-4-3`,
		`<video controls`,
		`loading="lazy"`,
		`<track kind="captions"`,
		`ui-video-fallback`,
	} {
		if !strings.Contains(tmpl, contract) {
			t.Errorf("video.html is missing contract %q", contract)
		}
	}
	// Zero-JS: native controls only, never autoplay (a11y + no-JS).
	if strings.Contains(tmpl, "autoplay") {
		t.Error("video.html must not render autoplay (no-JS and a11y)")
	}
	if strings.Contains(tmpl, "ui-video-demo") {
		t.Error("video.html must not ui-prefix demo scaffolding")
	}
}
