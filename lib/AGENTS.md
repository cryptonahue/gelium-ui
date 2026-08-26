# AGENTS.md — gelium-ui

Server-rendered UI for **Tailwind CSS 4** + **HTMX**. HTML-first, 0-JS contract by
default; token + class themes; optional small JS enhancement. **Not** a JS framework.

## Before you generate any UI

Read `skills/` — each skill is one actionable decision or flow. Start with
`skills/01-foundations.md`. After picking a SURFACE mode (`skills/02`), run
`skills/08-product-reasoning.md` to discover missing product-level UX before
drawing. The full decision pack is in `llms-ux.txt`.

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

## Gallery

The searchable icon gallery (item 0 of the roadmap) exposes top icon SVG sets as
opt-in `lib/icons/<pack>` packs — Tabler, Lucide, Heroicons, Phosphor — each with a
searchable gallery on `/docs/icon`, themeable via `currentColor`. Packs are SVG
server-rendered, never icon-fonts.

## Layout of this package

- `dist/gelium.css` — drop-in prebuilt bundle (themes + tokens + components).
- `styles/` — source CSS (`index.css` manifest, component sheets, `tokens.css`).
- `themes/` — `theme-material.css`, `theme-basecoat.css`.
- `templates/` — server-rendered partials (Go `html/template` blocks).
- `js/gelium.js` — optional enhancement (422 swap, toast, VT guard, slider fill).
- `llms-ux.txt` — the agent decision pack (SURFACE / SCREEN / WF / DATA / FEED /
  JOURNEY / MEDIA / SKEL ids).
- `skills/` — actionable decision skills for agents.
