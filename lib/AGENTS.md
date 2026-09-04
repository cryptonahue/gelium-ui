# AGENTS.md — gelium-ui

Server-rendered UI for **Tailwind CSS 4** + **HTMX**. HTML-first, 0-JS contract by
default; token + class themes; optional small JS enhancement. **Not** a JS framework.

## Before you generate any UI

First read `skills/00-agent-routing.md` and classify the requested outcome.
It is the canonical route selector; the rules below are downstream Gelium
contracts. Do not apply design-gated ceremony to `direct-exempt` work.

For a new library component, load `skills/14-component-implementation.md`
after routing.

Read this entrypoint, then `llms-ux.txt`, then `SKILLS.md`; each layer is part
of the protocol. Start with `skills/01-foundations.md`. After picking a SURFACE
mode (`skills/02`), run `skills/08-product-reasoning.md` to discover missing
product-level UX before drawing. Run `skills/13-reference-selection.md` after
SURFACE/SCREEN to resolve matching section and component references before the
wireframe. Before selecting components or writing markup, run
`skills/10-page-section-architecture.md` to contract each major page region
by purpose, audience, hierarchy, action, revelation, and recovery. Then run
`skills/11-design-criteria.md` for hierarchy, type, and density decisions. For
new screens, new flows, or substantial page redesigns, get the structure
approved through `skills/12-wireframe-approval.md` before writing markup.

## Route UI work proportionally

Classify the request before planning:

- **direct-exempt** — an understood, bounded copy, token, selector,
  accessibility, bug, or existing-contract correction with no page/flow
  architecture change. Inspect the relevant files, build, and run focused
  checks; do not manufacture a wireframe or ledger ceremony.
- **design-gated** — a new screen, new flow, or substantial redesign. This
  route is **required for design-gated** work. Follow
  `ROUTE → ORIENT → PLAN → ARCHITECT → APPROVE → BUILD → AUDIT → RELEASE`.
  Orient reads product/design artifacts, vocabulary, registry, and hard
  contracts; Plan records job, audience, states, and an intent wireframe;
  Architect validates it against real routes, data, permissions, templates,
  components, and no-JS/server behavior. Stop after Architect: show the
  buildable wireframe in the conversation and wait for an explicit approval of
  that packet before any markup or CSS. A request to make the page, `continua`,
  or a model-switch resume is not approval unless the human has seen that
  wireframe.
- **delegated-direct** — broad or multi-file work that benefits from fresh
  context but does not change screen/flow architecture. Delegate one bounded
  explorer or writer; the parent retains scope, integration, and verification.
- **escalate** — product intent, risk, scope, permissions, data, or architecture
  cannot be resolved from the current artifacts. Ask for the smallest concrete
  decision instead of inventing a solution.
- **full-sdd** — use OpenSpec only for cross-cutting work where durable proposal,
  design, specification, and task artifacts materially reduce ambiguity.

For design-gated work, copy the shipped
`skills/templates/gate-ledger.md` and
`skills/templates/wireframe-approval-packet.md` into the change record. The
ledger is structured evidence, not proof that a person read or approved a file.
Prebuild records decisions and approval; rendered evidence is collected only
after Build during Audit.

Golden rules that apply to every task:

1. **HTML-first** — native elements before ARIA; `div`/`span` never replace
   controls (`<button>`, `<input>`, `<select>`, `<dialog>`, `<a>`).
2. **Theme by class, never hex** — put `theme-material | theme-basecoat` (and
   `theme-dark` for dark) on `<html>`. No one-off color literals; use `--ui-*`
   tokens.
3. **0-JS first** — the main flow must complete with JS disabled. Progressive
   enhancement only.
4. **Server-first state** — navigable state = URL; validation = 422 +
   `X-Gelium-Validation`; **POST+303**; persistent feedback ≠ toast.
5. **States always** — empty, error, loading, success on every surface.
6. **Mobile** — touch targets ≥ 44px (`--ui-touch-target`), `min-width: 0` on
   scroll children, **never** `overflow-x: hidden` on `body` (no masking).
7. **DoD before done** — every surface passes `skills/07-dod-and-antislop.md`, starting
   from its step-0 artifacts gate (no `PRODUCT.md`/`DESIGN.md` → ask the user first),
   plus the per-screen usability checklist in `skills/09-usability-checklist.md`.
8. **Registry-first shells** — page-level layouts compose registered components
   (`ui-container`, nav primitives); custom shell CSS is spacing/width only.
9. **Verify mechanically** — run `scripts/ux-detect.sh` before claiming a surface done.
10. **Section purpose before components** — follow `ARCH-PRODUCT → ARCH-PAGE →
    ARCH-SECTION → ARCH-COMPONENTS → ARCH-TOKENS`; every major region needs a
    `SECTION-CONTRACT` before it can be rendered or styled.
11. **Plan before markup** — new screens, new flows, and substantial page
    redesigns require the `skills/12` approval gate and `skills/11` pre-emit
    critique. The agent MUST show the buildable wireframe in chat and wait.
    Accessibility fixes, bug fixes, contract corrections,
    component/mechanical changes, and already-approved small adjustments are
    exempt unless the scope expands into a new or substantially redesigned
    surface.

## Gallery

The searchable icon gallery (item 0 of the roadmap) dogfoods trusted
inline SVG (`.ui-icon`, `currentColor`). Library catalogs:

- Material Symbols rounded — `@material-symbols/svg-400` (fill; default)
- Tabler Icons — `@tabler/icons` (outline stroke + filled; MIT)

Consumers run `scripts/extract-used-icons.mjs` so the app binary embeds
only names referenced as string literals. Unprefixed names use `--set`
(`material` default, or `tabler`). Prefixed names pick a catalog:

- `data-gelium-icon="chevron_right"` / `icons.SVG("settings")`
- `data-gelium-icon="tabler:chevron-right"` / `icons.SVG("ms:home")`
- `icons.SVG("tabler-filled:home")`

Prefer one set per product; prefixes are for an explicit mixed screen.
Packs are SVG server-rendered, never icon-fonts.

## Layout of this package

- `dist/gelium.css` — drop-in prebuilt bundle (themes + tokens + components).
- `styles/` — source CSS (`index.css` manifest, component sheets, `tokens.css`).
- `themes/` — `theme-material.css`, `theme-basecoat.css`.
- `templates/` — server-rendered partials (Go `html/template` blocks).
- `js/gelium.js` — optional enhancement (422 swap, toast, VT guard, slider fill).
- `llms-ux.txt` — the agent decision pack (SURFACE / SCREEN / WF / DATA / FEED /
  JOURNEY / MEDIA / SKEL ids).
- `skills/` — actionable decision skills for agents.
