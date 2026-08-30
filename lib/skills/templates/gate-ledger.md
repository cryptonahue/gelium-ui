# Gelium gate ledger

> Copy this file into the change record for a `design-gated` task. Until a
> future preflight defines its canonical consumer path, keep it beside the
> approved change artifacts. It is structured attestation and evidence, never
> proof of cognitive reading, approval, commit, publish, or deployment authority.

```json
{
  "schema_version": 1,
  "route": "design-gated",
  "scope": {
    "routes": ["/example"],
    "owned_paths": ["templates/example.html"],
    "shared_paths": []
  },
  "reading": [
    {
      "path": "PRODUCT.md",
      "status": "attested",
      "note": ""
    }
  ],
  "gates": {
    "plan": {
      "status": "pass",
      "evidence": ["packet.md#plan"]
    },
    "architecture": {
      "status": "pass",
      "evidence": ["packet.md#architecture"]
    },
    "criteria_plan": {
      "status": "pass",
      "evidence": ["packet.md#criteria"]
    },
    "approval": {
      "status": "approved",
      "packet": "packet.md",
      "approver": "",
      "channel": "",
      "date": ""
    },
    "rendered_audit": {
      "status": "pending",
      "evidence": []
    }
  },
  "release": {
    "detector": {
      "status": "pending",
      "evidence": []
    },
    "checks": {
      "tests": [],
      "builds": []
    },
    "authority_matrix": {
      "status": "pending",
      "evidence": []
    }
  },
  "exceptions": []
}
```

## Reading attestations

List the Gelium entrypoint, `llms-ux.txt`, `SKILLS.md`, applicable skills,
consumer product/design artifacts, vocabulary/registry, and any route or
contract source actually used. `attested` means the worker attests it considered
the named artifact; it does not assert that a tool can prove mental inspection.

## Prebuild gates

- **Plan:** user job, audience, SURFACE/screen, one primary action, states,
  constraints/non-goals, and the intent wireframe.
- **Architecture:** routes, handler, data, permissions, templates, components,
  tokens, and no-JS/server contracts yield the buildable wireframe.
- **Criteria plan:** intended hierarchy, DOM order, actions, responsive/theme
  intent, states, accessibility/no-JS, preserved contracts, and DESIGN-MEMORY.
- **Approval:** required for `design-gated` work; a chat decision is valid only
  when scope, approver, date/channel, and packet version are recorded.

## Rendered audit

After Build, attach evidence for wide/narrow, selected light/dark class-routed
themes, realistic content, representative states, keyboard/focus/touch,
no-JS/server behavior, detector output, and applicable tests/builds. Preserve
raw detector findings. An exception requires a bounded scope, reason, risk,
owner, and verification evidence; it is never a clean pass.
