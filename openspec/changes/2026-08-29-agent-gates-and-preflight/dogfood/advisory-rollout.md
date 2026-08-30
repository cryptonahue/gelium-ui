# Design-gated required rollout

## Status

For future `design-gated` work, the Orient → Plan → Architect → **visible
Approve** sequence is **required** before markup. A ledger, packet, and
`gelium-preflight` / `gelium-ux-detect` record remain evidence; they still do
not grant permission to commit, push, publish, or deploy.

## What is required now

For a new screen, new flow, or substantial redesign:

```text
ROUTE → ORIENT → PLAN → ARCHITECT → APPROVE → BUILD → AUDIT → RELEASE
```

1. Write the buildable packet.
2. **Show** the desktop and mobile wireframes in the conversation.
3. Wait for an explicit yes to **that** packet.
4. Only then write markup/CSS.
5. Keep a JSON ledger and run preflight with declared changed paths. Preserve
   detector output and any exceptions with rule, path, owner, evidence, and
   `expires_at`.

“Make the page”, `continua`, “oks dale”, or a resume after a model/context
switch is **not** approval unless the human already saw that wireframe and
approved it.

## What remains direct

`direct-exempt` work remains direct: bounded copy, token, selector,
accessibility, behavior-bug, or contract corrections with no page/flow
architecture shift use the narrowest relevant checks. Do not create a ledger
or wireframe packet merely because a detector exists.

## Legacy path

`lib/scripts/ux-detect.sh` remains the consumer-compatible positional/default
detector. `gelium-ux-detect` is the scoped JSON/text command. Do not remove the
legacy command.

## Why required now

The recipe-public-feed dogfood showed detector exceptions and one false
positive, which was fixed. A later DeepFilter discovery implementation skipped
showing the wireframe because a continue after a model-limit interruption was
treated as approval. Required mode exists to prevent that skip, not to invent
CI deploy locks.

## Still not granted

Passing preflight or the detector does not authorize commit, publish, deploy,
or a claim that a visual review is complete.
