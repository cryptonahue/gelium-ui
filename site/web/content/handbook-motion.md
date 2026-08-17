# Motion

Motion should **explain change** or **focus attention**, not decorate. Gelium is 0-JS-first: most UI works without animation. When motion exists, it must respect user preference and stay tokenized.

## Sources

| Topic | Source |
|---|---|
| Reduced motion | [WCAG 2.2 – prefers-reduced-motion](https://www.w3.org/WAI/WCAG22/Understanding/animation-from-interactions.html) (understanding animation from interactions) |
| M3 motion roles | [Material 3 motion](https://m3.material.io/styles/motion/overview/how-it-works) (transitions that clarify hierarchy/navigation) |
| Severity vs interruption | [NNG: error-message guidelines](https://www.nngroup.com/articles/error-message-guidelines/) (don’t over-interrupt) |
| Gelium tokens | `--ui-motion-*`, `--ui-easing-*` in [Tokens](/docs/tokens) |
| Docs VT | Same-document view transitions with reduced-motion guard (implementation in consumer JS) |

## When to move

| ID | Use motion | Example |
|---|---|---|
| **MOTION-FEEDBACK** | Brief confirmation of an action | Toast enter/exit |
| **MOTION-STATE** | Show expand/collapse or open/close | Disclosure, dialog presence |
| **MOTION-NAV** | Optional continuity between peer views | Same-document VT when supported |
| **MOTION-NONE** | Default | Static forms, tables, most admin CRUD |

## When not to move

- Decorative loops on content the user is trying to read.
- Large parallax or continuous animation without an essential purpose.
- Motion that is the **only** way to understand state (always pair with text/structure).
- Ignoring `prefers-reduced-motion: reduce` — under reduce, snap or skip non-essential motion.

## Implementation rules (Gelium)

1. Prefer **CSS** + tokens (`--ui-motion-short`, etc.) over ad-hoc durations.
2. Honor **`prefers-reduced-motion`** in CSS and any JS enhancement.
3. Dialogs/toasts may animate lightly; **duration stays short**.
4. Never block task completion on an animation finishing.
5. View transitions are progressive enhancement only.

## Anti-patterns

- Page-load animation on every navigation in a dense admin.
- Error shake as the only error indicator (still need text + `FEED-VAL`).
- Custom cubic-beziers per screen instead of tokens.

## Checklist

1. Is this MOTION-FEEDBACK, STATE, NAV, or NONE?
2. What is the reduced-motion behavior?
3. Which token duration/easing?
4. Does the UI still work with JS off and motion off?

## See also

- [Feedback](/docs/feedback) · [Accessibility](/docs/accessibility) · [Performance](/docs/performance) · [`/llms-ux.txt`](/llms-ux.txt)
