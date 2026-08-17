package app

import (
	"strings"
	"testing"
)

// TestRecipeCriteriaBridgesMapToHandbook proves each live recipe list page
// exposes the SKEL/FEED/DATA bridge and links Patterns/Feedback.
func TestRecipeCriteriaBridgesMapToHandbook(t *testing.T) {
	cases := []struct {
		path string
		want []string
	}{
		{
			path: "/recipes/admin-resource",
			want: []string{
				"Maps to Gelium criteria",
				"SKEL-ADMIN-RESOURCE",
				"DATA-TABLE",
				"FEED-VAL",
				`href="/docs/patterns"`,
				`href="/docs/feedback"`,
			},
		},
		{
			path: "/recipes/ops-queue",
			want: []string{
				"Maps to Gelium criteria",
				"SKEL-OPS-QUEUE",
				"DATA-QUEUE",
				"JOURNEY-QUEUE",
				`href="/docs/patterns"`,
			},
		},
		{
			path: "/recipes/public-feed",
			want: []string{
				"Maps to Gelium criteria",
				"SKEL-FORUM",
				"DATA-FEED",
				`href="/docs/data-display"`,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			body := getOKBody(t, tc.path)
			for _, w := range tc.want {
				if !strings.Contains(body, w) {
					t.Errorf("%s missing bridge contract %q", tc.path, w)
				}
			}
		})
	}
}

// TestContentStyleDefinesHeadingGrammar proves content-style documents H1–H3
// and list/table/quote rules for agents and editors.
func TestContentStyleDefinesHeadingGrammar(t *testing.T) {
	body := getOKBody(t, "/docs/content-style")
	for _, want := range []string{
		"Content structure",
		"Canonical handbook outline",
		"H1",
		"H2",
		"H3",
		"Blockquote",
		"DOC-H1",
		"nngroup.com",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("content-style missing %q", want)
		}
	}
}
