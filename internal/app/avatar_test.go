package app

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	webassets "geliumui/web"
)

// renderPartial executes one recipe/primitive partial from the same embedded
// template set the server uses, for focused primitive tests.
func renderPartial(t *testing.T, name string, data any) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/*.html"))
	var b bytes.Buffer
	if err := tmpl.ExecuteTemplate(&b, name, data); err != nil {
		t.Fatalf("render template %q: %v", name, err)
	}
	return b.String()
}

// TestAvatarDecorativeInitialsAreAriaHidden proves a decorative initials avatar
// (paired with a visible name) is aria-hidden and carries the initials text.
func TestAvatarDecorativeInitialsAreAriaHidden(t *testing.T) {
	body := renderPartial(t, "avatar", avatarView{Initials: "AR", Decorative: true, Size: "sm"})
	for _, contract := range []string{
		`class="ui-avatar ui-avatar--sm"`,
		`aria-hidden="true"`,
		`class="ui-avatar-initials">AR</span>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("decorative initials avatar is missing %q", contract)
		}
	}
}

// TestAvatarDefaultsToMdSize proves an empty size modifier renders the md size.
func TestAvatarDefaultsToMdSize(t *testing.T) {
	body := renderPartial(t, "avatar", avatarView{Initials: "BT", Decorative: true})
	if !strings.Contains(body, `class="ui-avatar ui-avatar--md"`) {
		t.Error("avatar must default to the md size modifier")
	}
}

// TestAvatarImageCarriesAltWhenMeaningful proves an image avatar that carries
// meaning on its own is announced via its alt text and is not aria-hidden.
func TestAvatarImageCarriesAltWhenMeaningful(t *testing.T) {
	body := renderPartial(t, "avatar", avatarView{ImageSrc: "/img/alicia.png", ImageAlt: "Alicia R.", Size: "lg"})
	for _, contract := range []string{
		`class="ui-avatar ui-avatar--lg"`,
		`src="/img/alicia.png"`,
		`alt="Alicia R."`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("meaningful image avatar is missing %q", contract)
		}
	}
	if strings.Contains(body, "aria-hidden") {
		t.Error("a meaningful image avatar must not be aria-hidden")
	}
}

// TestAvatarImageDecorativeHasEmptyAlt proves a decorative image avatar renders
// aria-hidden with an empty alt (never a missing alt or the filename).
func TestAvatarImageDecorativeHasEmptyAlt(t *testing.T) {
	body := renderPartial(t, "avatar", avatarView{ImageSrc: "/img/alicia.png", Decorative: true})
	for _, contract := range []string{
		`aria-hidden="true"`,
		`alt=""`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("decorative image avatar is missing %q", contract)
		}
	}
}

// TestAvatarSizeSanitizesToClosedSet proves avatarSize clamps unknown sizes to
// the md default.
func TestAvatarSizeSanitizesToClosedSet(t *testing.T) {
	cases := map[string]string{"sm": "sm", "md": "md", "lg": "lg", "xl": "md", "": "md"}
	for in, want := range cases {
		if got := avatarSize(in); got != want {
			t.Errorf("avatarSize(%q) = %q, want %q", in, got, want)
		}
	}
}
