# DESIGN.md template (consumer apps)

Copy this file to your application repo root as `DESIGN.md`. It records the visual and interaction direction without forking Gelium's token system.

You can answer the starter questions in plain language. If you do not know an answer, write `Unknown`, `To decide`, or `N/A` with a short reason. The agent may propose a safe default, but must record it as an assumption and ask for confirmation when it affects page architecture.

## Start in plain language

1. Which screens or URLs do you expect to have first?
2. Should each surface mainly help people **do a task**, **read**, **decide**, or **experience** something?
3. Should navigation be mainly a top bar, a side navigation, tabs, or something else?
4. What products or sites have a useful quality to borrow, and what should we avoid copying?

The agent should turn these answers into the structured sections below and ask only the follow-up questions needed for a safe implementation.

## Screen inventory

One row per important URL or surface. The agent can propose the Gelium mode and screen type; do not invent them when the product decision is unclear.

| URL | Surface: Operate / Read / Persuade / Experience | Screen type | User job | Primary action |
| --- | --- | --- | --- | --- |
| `/` | … | hub / list / detail / form / confirm / settings / queue / result | … | … |

## Chrome and navigation

- Navigation model: top bar / side navigation / tabs / other / undecided
- Footer: none / legal / contact / links / undecided
- What stays stable across screens: …
- Role-specific navigation or boundaries: …

## Theme

- HTML class: `theme-material` | `theme-basecoat` | `theme-neubrutalism` | undecided
- Dark mode: explicit `theme-dark` class route; never media-only as the sole authority
- Brand direction in one sentence: …

## Density

- Operate surfaces: cozy | compact
- Read/Persuade surfaces: comfortable | cozy
- Exception and reason: …

## Data display defaults

- Collections: table / list / cards / decide per screen
- Detail records: description list / sections / other
- Filters, sorting, or pagination: …
- Never use a display choice that conflicts with the data's purpose.

## Components and composition

Prefer registered Gelium `ui-*` partials from `gelium-ui/templates`. A new primitive needs a token-first rationale and a vocabulary/registry check.

- Reusable components or recipes: …
- Components explicitly not wanted: …

## States and recovery defaults

For every applicable screen, plan rest, loading, empty, error, success, and partial states.

- Empty state: explain what is empty, why it matters, and the next action.
- Error state: explain the problem and a recovery action.
- Success state: use a result page when the outcome must persist; do not rely on a toast alone.
- No-JS path: …

## Motion and accessibility

- Motion: `MOTION-NONE` by default on Operate; honor `prefers-reduced-motion`.
- Keyboard and focus: …
- Contrast and non-color cues: …
- Narrow layout and touch target notes: …

## Anti-slop checklist

Keep Gelium's anti-slop rules and check project-specific constraints before markup.

- [ ] No purple-blue gradient used as a default identity.
- [ ] No centered hero plus three identical feature cards by default.
- [ ] No invented metrics, avatars, testimonials, or “trusted by” logos.
- [ ] No rounded surface and shadow applied to every element.
- [ ] One highlighted primary action per screen.
- [ ] Operate surfaces do not use decorative motion.
- [ ] Empty, loading, error, and success states have a server-rendered path.
- [ ] No one-off feature colors, fonts, or spacing scales.
- [ ] …

## References to borrow from

| Reference | Borrow | Avoid | Reason |
| --- | --- | --- | --- |
| … | … | … | … |

## Brand notes

Logo, typeface, and any extra brand color must map into Gelium tokens. Do not hardcode one-off hex values in features.

- Logo: …
- Typeface: use theme `--ui-font-*` unless a brand choice is recorded here.
- Brand tokens: …

## Visual definition of done

- [ ] A stranger understands the page purpose in five seconds.
- [ ] The screen inventory has one primary action per URL.
- [ ] Desktop and narrow layouts preserve the same regions and reading order.
- [ ] Theme, states, focus, contrast, and reduced motion are checked.
- [ ] Recovery works without JavaScript where the flow requires it.

## Open decisions and assumptions

- Decision: …
- Status: open / assumed / confirmed
- Impact if wrong: …
- Needed before implementation: yes / no

Gelium docs: `/docs/themes`, `/docs/tokens`, `/docs/density`, `/docs/agent-workflow`.
