---
name: gelium-ui
description: Use when composing or reviewing Gelium UI work. Select the smallest outcome-first route, then load the applicable Gelium contracts.
version: 0.6.3
license: MIT
---

# Gelium UI agent skill

Use `skills/00-agent-routing.md` as the canonical first decision layer. It
selects `direct-exempt`, `delegated-direct`, `design-gated`, `escalate`, or
`full-sdd` and defines the boundary between verification and repository
delivery.

After routing, read `AGENTS.md`, `llms-ux.txt`, `SKILLS.md`, and only the
applicable downstream skills. This wrapper exists so agent registries such as
Gentle AI can discover Gelium; the routing and UX rules remain in their linked
source files and are not duplicated here.
