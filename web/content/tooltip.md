# Tooltip

Tooltip is a short, contextual label that appears on hover or keyboard focus near a control. Use a tooltip when a control needs a clarification that should not permanently occupy the layout — icon buttons, abbreviated labels, unfamiliar terms. Gelium reimplements the Material 3 tooltip over **server-rendered HTML with pure CSS interaction**: the styled `.ui-tooltip` surface lives inside a `.ui-tooltip-host` wrapper and is revealed with the native `:hover` / `:focus-within` pseudo-classes. There is **no component JavaScript** — the show/hide is declarative CSS, works without JavaScript, and is keyboard-focus accessible.

## Anatomy

- **Host** — `ui-tooltip-host`, the `position: relative` inline wrapper that pairs a control with its tooltip and scopes the reveal.
- **Tooltip** — `ui-tooltip`, the floating surface (container color, `--ui-radius-xs` shape, 4px/8px padding, body-small supporting text). It uses `role="tooltip"` and the trigger references it with `aria-describedby`.
- **Subhead** — `ui-tooltip-subhead`, the rich variant's title text (title-small type).
- **Supporting text** — `ui-tooltip-supporting-text`, the body copy that carries the explanation.
- **Action** — `ui-tooltip-action`, an optional real `<a href>` link in the rich variant.

## Variants

- **Plain tooltip** — a single line of supporting text, the default `ui-tooltip`.
- **Rich tooltip** — `ui-tooltip--rich` adds a subhead, supporting text and an optional action link.

The demo below demonstrates both variants on a Button (reusing the Button contract), an Icon button (reusing the Icon button contract) and a link.

## What triggers a tooltip?

- **Rest** — the tooltip is hidden (`visibility: hidden`, zero opacity, no layout impact).
- **Hover** — `:hover` on the host reveals it.
- **Focus** — `:focus-within` on the host reveals it for keyboard users the moment the control is focused.
- **Position** — the surface sits below the anchor by default; the `ui-tooltip--top` modifier places it above.
- **Reduced motion** — `prefers-reduced-motion: reduce` disables the fade transition.

## Accessibility

- The trigger is a real control (button, link, icon button) whose accessible name comes from its own visible label or `aria-label`. The tooltip surface has `role="tooltip"` and is linked with `aria-describedby`, so screen readers announce the description when the control is focused and the surface is shown.
- The tooltip never hides essential information: every demo control is self-explanatory on its own, and the tooltip only adds context. Nothing that a user must know is exclusive to the tooltip.
- Focus does not shift geometry and the reveal causes no layout shift.

## No-JS behavior

Show, hide and focus-reveal are pure CSS (`:hover` / `:focus-within`), so the whole flow works with JavaScript disabled. The rich action link is a real `<a href>` that navigates. The platform-first audit (see Compatibility) rejected component JavaScript because CSS hover/focus reveal fully covers the interaction contract.

## Compatibility

Audited August 2026:

- `:hover` and `:focus-within` — universally supported (Baseline).
- Popover API (`popover`, `popovertarget`) — Baseline 2024, but popover surfaces are toggle- or script-opened and need the Interest Invokers API (the not-yet-standard mechanism that would open a popover on hover/focus via `interesttarget`/`interestaction`) to open on hover/focus, so they are not used here.
- Interest Invokers (`interesttarget` / `interestaction`) — **not Baseline as of August 2026** (no caniuse/MDN/Chrome-Status support). The roadmap's preferred mechanism (a `popover` tooltip opened by `interestaction="show-popover"`) is therefore not viable without JavaScript, so Gelium uses the CSS hover/focus reveal instead — the roadmap's "accessible visible fallback".
- CSS anchor positioning — Baseline 2026, but the tooltip is a small surface attached to its own host, so plain `position: absolute` inside the relative host covers below/above placement with broader support and no extra dependency.
- Native `title` — universal, but not styleable to Material anatomy and not reliably keyboard-revealed, so it is not used as the primary mechanism.

## Trust boundary

The demo markup is static and trusted: `html/template` renders it verbatim, the inline SVG glyphs are trusted internal icons (`aria-hidden`, `focusable="false"`), and there is no user-supplied content.

## Visual checklist

- [ ] Inverse-surface (plain) / surface (rich) container with the right shape and padding
- [ ] Plain and rich variants; rich subhead + supporting text + optional action link
- [ ] Reveals on hover and on keyboard focus, hidden at rest
- [ ] Below by default, above with the `ui-tooltip--top` modifier
- [ ] Focus does not shift geometry; no layout shift
- [ ] Trigger accessible name independent of the tooltip; `aria-describedby` wiring
- [ ] Reduced motion disables the fade; forced-colors keeps the surface discernible
- [ ] Light/dark, narrow/wide, RTL
