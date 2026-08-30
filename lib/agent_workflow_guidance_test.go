package lib

import (
	"strings"
	"testing"
)

// TestAgentWorkflowGuidanceContract keeps the distributed agent pack aligned:
// every installed guidance entrypoint must describe the same proportional route
// and point at the shipped records needed for a design-gated change.
func TestAgentWorkflowGuidanceContract(t *testing.T) {
	for _, tc := range []struct {
		name string
		path []string
		want []string
	}{
		{
			name: "entrypoint",
			path: []string{"lib", "AGENTS.md"},
			want: []string{
				"ROUTE → ORIENT → PLAN → ARCHITECT → APPROVE → BUILD → AUDIT → RELEASE",
				"direct-exempt",
				"design-gated",
				"required for design-gated",
				"skills/templates/gate-ledger.md",
			},
		},
		{
			name: "decision pack",
			path: []string{"lib", "llms-ux.txt"},
			want: []string{
				"ROUTE → ORIENT → PLAN → ARCHITECT → APPROVE → BUILD → AUDIT → RELEASE",
				"prebuild criteria plan",
				"rendered audit",
				"media-metadata-unknown",
				"skills/templates/wireframe-approval-packet.md",
			},
		},
		{
			name: "skill index",
			path: []string{"lib", "SKILLS.md"},
			want: []string{
				"ROUTE → ORIENT → PLAN → ARCHITECT → APPROVE → BUILD → AUDIT → RELEASE",
				"skills/templates/gate-ledger.md",
				"skills/templates/wireframe-approval-packet.md",
			},
		},
		{
			name: "criteria skill",
			path: []string{"lib", "skills", "11-design-criteria.md"},
			want: []string{
				"Criteria plan (prebuild)",
				"Rendered audit (postbuild)",
				"pass-with-escalation",
				"Gelium icon allowlist",
				"extract-used-icons",
				"--set material|tabler",
			},
		},
		{
			name: "approval skill",
			path: []string{"lib", "skills", "12-wireframe-approval.md"},
			want: []string{
				"Intent wireframe (Plan)",
				"Buildable wireframe (Architect)",
				"approved | changes-requested | declined | exception",
				"Visible packet",
				"continua",
			},
		},
		{
			name: "gate ledger template",
			path: []string{"lib", "skills", "templates", "gate-ledger.md"},
			want: []string{
				`"schema_version": 1`,
				"Reading attestations",
				"Prebuild gates",
				"Rendered audit",
			},
		},
		{
			name: "approval packet template",
			path: []string{"lib", "skills", "templates", "wireframe-approval-packet.md"},
			want: []string{
				"Packet version:",
				"Plan — intent wireframe",
				"Architect — buildable wireframe",
				"Decision: approved | changes-requested | declined | exception",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			content := repositoryFile(t, tc.path...)
			for _, required := range tc.want {
				if !strings.Contains(content, required) {
					t.Errorf("%s is missing %q", strings.Join(tc.path, "/"), required)
				}
			}
		})
	}
}
