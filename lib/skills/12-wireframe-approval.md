# Skill: Wireframe approval gate

Use this skill to decide when a page or flow needs an architecture and wireframe
approval before markup or CSS changes. It complements skills 01–11 and keeps
small, low-risk fixes moving. Approval is a human decision about a proposed
structure, not approval of pixel polish or a runtime slash command.

## When approval is required

Require approval **before implementation** when the change introduces:

- a new screen or route;
- a new user flow, lifecycle, or multi-step interaction; or
- a substantial redesign of an existing page, meaning its information
  hierarchy, major regions, primary action, reading order, or owner/reader
  boundary materially changes.

For these changes, first inspect the real product, route, template, handler,
CSS, components, states, and no-JS/server contract. Do not edit implementation
files during the design pass.

## Explicit exemptions

Do not block implementation approval for:

- accessibility fixes or inclusive responsive improvements;
- bug fixes that restore the documented behavior;
- corrections to an existing Gelium, URL, server, component, or data contract;
- component-level or mechanical changes with no page/flow architecture change
  (for example token replacement, selector repair, copy correction, or a
  mechanical migration); or
- small adjustments that are already covered by an approved architecture or
  wireframe and do not materially change its scope.

An exempt change still follows skills 01–11, the existing contracts, and the
DoD. If its implementation reveals a new screen, flow, or substantial
architectural change, stop at that boundary and apply this gate to the newly
expanded scope.

## Approval packet

For a gated change, prepare a compact packet in English technical artifact
convention:

```text
Change:
Product job / audience:
Existing route and contracts:
Scope: new screen | new flow | substantial redesign
Section inventory: ordered regions and purpose
Primary action and action hierarchy:
Desktop wireframe:
Mobile wireframe:
States and recovery:
Accessibility and no-JS behavior:
Reuse / DESIGN-MEMORY decision:
Open questions and explicit trade-offs:
```

Wireframes are structural, not decorative. Use ASCII or another reviewable
text form to show reading order, major regions, primary/supporting actions,
disclosures, recovery, and owner-only areas. Do not spend the approval packet on
exact colors, shadows, or pixel dimensions; apply skill 11 after structure is
approved.

## Approval workflow

1. **Classify scope.** Name the trigger and confirm whether an exemption
   applies. Completion: the scope classification is recorded.
2. **Read before drawing.** Inspect the existing implementation and apply
   skills 02, 08, and 10. Completion: route, server/no-JS behavior, sections,
   and existing Gelium contracts are listed.
3. **Prepare the packet.** Include desktop and mobile wireframes, states,
   accessibility, reuse, and open decisions. Completion: another person can
   understand the intended structure without seeing markup.
4. **Request a decision.** Record `approved`, `changes requested`, or
   `declined`, with date/author and the packet version. Completion: gated work
   has an explicit approval outcome before implementation.
5. **Implement within scope.** Follow the approved section order and reuse
   existing components/tokens. Completion: implementation matches the packet,
   or any material deviation is recorded and re-approved.
6. **Verify.** Run rendered-HTML, responsive, theme, state, accessibility,
   no-JS, and mechanical checks from skills 07, 09, and 11. Completion: all
   failures are fixed, escalated, or explicitly documented.

## Intentional gate bypass / exception record

An exempt change does not need a ceremonial wireframe. When a change that would
normally be gated intentionally proceeds without prior approval, record an
explicit exception rather than silently bypassing the gate. This is a short
note in the change/PR/task record, not a new blocking workflow:

```text
WIREFRAME gate exception: yes
Scope that would normally require approval:
Reason for bypass:
Risk boundary / files and routes:
Maintainer or user decision:
Follow-up review required: yes | no
Verification evidence:
```

Use this only when timing, incident recovery, a tiny contained expansion, or
another stated constraint makes prior approval impractical. Keep the scope
bounded and complete the record as part of implementation. Accessibility fixes,
bug fixes, contract corrections, component/mechanical changes, and already
approved small adjustments remain ordinary exemptions and do not require an
exception record unless their scope expands into a gated change.

## Handoff

After approval, hand off to implementation with the packet and its decision
record. The implementation handoff is:

```text
skills/01–02 → skills/08 → skills/10 → skill 11 → skill 12 (when gated)
→ registered components → tokens/skin → skill 11 critique → skills/09 → skills/07
```

For exempt small work, use the narrowest relevant path and do not manufacture a
wireframe. In all cases, plan before markup, preserve semantic HTML and
existing Gelium contracts, and keep the primary flow functional with JS
 disabled.

## Verification checklist

- [ ] Scope is new screen, new flow, substantial redesign, or documented exemption.
- [ ] Gated work has a read-only inspection and an approved packet before markup.
- [ ] Exempt work is not delayed by an unnecessary approval ceremony.
- [ ] Any intentional bypass is recorded as an explicit exception with bounded scope.
- [ ] Packet covers desktop/mobile order, actions, states, recovery, accessibility,
      no-JS behavior, and DESIGN-MEMORY reuse.
- [ ] Rendered implementation matches the approved structure or records a
      material deviation and decision.
- [ ] Skills 07, 09, and 11 verification is complete.
