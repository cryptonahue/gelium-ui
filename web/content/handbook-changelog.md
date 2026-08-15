# Changelog

This page mirrors the [CHANGELOG.md](https://github.com/cryptonahue/gelium-ui/blob/main/CHANGELOG.md) in the repository root.

Todas las cambios notables del proyecto Gelium UI se documentan en este archivo.

El formato sigue [Keep a Changelog](https://keepachangelog.com/es/1.1.0/) y el proyecto usa [Semantic Versioning](https://semver.org/).

## [0.4.0] — 2026-08-15

### Added
- **Content style guide** (`/docs/content-style`): reglas de copy para errores (recognize + recover, nunca culpar al usuario — NNG), toasts (verbo + resultado), empty states (qué hay / por qué está vacío / qué hacer), banners y validation summary. Patrones de escritura editorial: plain English, voz activa, sin "please", AP Style.
- **Screen-reading rules**: sección "Reading on screen" en el content style guide — fundamento NNG del eye-tracking (la mayoría escanea, no lee), patrón F (lo importante arriba-izquierda), web copy más corto que impreso. Sección "Paragraphs and sentences": párrafos 2-4 oraciones, oraciones ≤ 25 palabras, pirámide invertida, listas sobre prosa.
- **Copy length contract**: `TestComponentPagesKeepSentencesUnder25Words` — ninguna página de componente puede tener oraciones > 25 palabras (179 reescritas).
- **Copy contract test**: `TestRecipeErrorCopyUsesActionPattern` + `TestRecipeEmptyStatesCarryActionLanguage` — errores con patrón de acción, empty states accionables.
- **Acknowledgments page** (`/docs/acknowledgments`): reconocimiento explícito de todas las fuentes de inspiración — Material Design 3, USWDS, GOV.UK, Mozilla Protocol, Base UI, Basecoat UI, Naive UI, Name That UI, Material Web, shadcn/templ — con qué se tomó, cómo se adaptó y licencia.
- **Information architecture page** (`/docs/information-architecture`): regla concepto-antes-que-referencia, criterios para agregar páginas (quién navega / qué tarea / concepto-vs-referencia) y un agent prompt para que LLMs auditen la IA de los docs.
- **Choose the right control page** (`/docs/choose-the-right-control`): tabla de decisión para elegir el componente de input correcto (Radio vs Select vs Checkbox vs Switch vs Slider vs Text field vs Menu) con reglas de oro.
- **Guidance sections en las 28 páginas de componentes**: When to use / When not to use / Usability / Accessibility, con cross-links a la página de decisión.
- **Handbook de 6 páginas**: Themes, Tokens, Server contracts, Accessibility, Design principles, Information architecture — concepto antes que referencia.
- **Búsqueda client-side de docs**: `#docs-search-index` JSON + `search.js`, fallback 0-JS a GET `/docs?q=`.
- **Theme switcher nativo** (`<select name="theme">`) y **scheme switcher real** (`<input type="checkbox" role="switch">`) — ambos 0-JS con progressive enhancement.
- **GitHub link en el topbar** de los docs.

### Changed
- **Wire contract migrado**: `loom:*` / `X-Loom-*` → `gelium:*` / `X-Gelium-*` (decisión del owner: el proyecto es nuevo, migrar ahora). El prefijo wire ahora coincide con el nombre del producto.
- **Legibilidad del prose**: medida 48rem (~90 chars) → **65ch**, `text-wrap: pretty` (sin orphans), `text-wrap: balance` en headings, `hyphens: auto`, `text-box-trim` progresivo.
- **Ritmo vertical**: breadcrumb → título, provenance → título, h2 → h3 con márgenes bidireccionales (antes solo en una dirección).
- **Contraste AA**: el `fg-muted` de Basecoat light corregido (4.35:1 → 4.75:1) para cumplir WCAG AA; test de contrato en ambos themes.
- **Copy de errores unificado** al patrón de acción: "Name is required." → "Enter the project name."; "Choose a valid status." → "Choose a status from the list." (antes mezclaba dos voces).
- **Empty state del public feed**: ahora con CTA ("Follow more people to fill this feed.").
- **Jerarquía del sidebar**: Handbook movido a posición 2 (después de Getting started, antes de componentes).
- **Renombrado de residuos**: "Loom UI" → "Gelium UI" en 4 docs raíz (MATERIAL-WEB-PROGRESS, prompts, roadmap).

### Fixed
- `dependency-metadata.md` claim stale sobre Phase I.
- Roadmap sin marcadores de fase completada (Phase I y J ahora DONE).

### Enforced (nuevos tests de contrato)
- `styles_readability_test.go`: 65ch, text-wrap, hyphens, text-box-trim, line-height ≥ 1.6, ritmo vertical, breadcrumb.
- `styles_prose_contrast_test.go`: WCAG AA para prose en ambos themes.
- `copy_contract_test.go`: patrón de acción en errores, empty states accionables, oraciones ≤ 25 palabras.
- `content_name_that_ui_test.go`: secciones Name That UI (Alternative names + Agent prompt) en páginas servidas.
- `handbook_test.go`: páginas del Handbook renderizan, en nav, en sitemap.

## [0.3.0] — 2026-08-14

### Added
- **SDD completo del roadmap A-J**: verificación formal de las 10 fases (A-J) con pares RED→GREEN.
- **Registry sync guards**: `registry_sync_test.go` — los registries (component, pattern, theme, dependency, agent-prompts, screen-composition) deben referenciar archivos reales.
- **Landing mejorado**: FAQ (Base UI), claims con checkmarks (Naive UI), demo card visual (Basecoat), link a GitHub, BASE_URL documentado.
- **Docs root explicador**: qué es Gelium UI (y qué NO es), Quick start en 4 pasos, temas explicados (Material, Basecoat, Base UI como vocabulario nunca runtime).
- **Basecoat theme** (Phase I, PR #19) — theme completo, light + dark clase única.

### Changed
- `index.md` cableado a `/docs` (estaba embebido pero sin ruta que lo sirviera).
- Wire contract documentado como canónico (`gelium-ui-wire-compatibility.md` reescrito).

## [0.2.0] — 2026-08-13

### Added
- **Theme mechanism** (Phase H, PR #18): selección por clase en el documento raíz (`<html>`), dark por ruta de clase única (sin `@media prefers-color-scheme`), themes swappables sin tocar markup.
- **Screen recipes** (Phase G): Admin Resource, Ops Queue, Public Feed — patrones de composición.
- **Public content patterns** (Phase F): 14 patrones de contenido con card slots.
- **Server contracts**: GET+query estables, POST+303, 422 + header `X-Loom-Validation` (luego `X-Gelium-Validation`), `loom:toast` (luego `gelium:toast`).

## [0.1.0] — 2026-08-12

### Added
- **Núcleo de Gelium UI**: componentes server-rendered en Go + HTMX, cero JS de componentes.
- **Arquitectura en 6 capas**: core tokens → themes → componentes → patrones → recipes → screens.
- **Tokens `--ui-*`**: vocabulario de typography, color roles, spacing, elevation.
- **Dos themes en bundle único**: Material (default, M3) + Basecoat.
- **Contratos de acceso y verificación** (Phases A-D): native semantics, focus rings, aria-* en toda la superficie.
- **Docs shell**: navegación, sidebar, breadcrumbs, búsqueda (deshabilitada), theme/scheme switchers.

[0.4.0]: https://github.com/cryptonahue/gelium-ui/releases
[0.3.0]: https://github.com/cryptonahue/gelium-ui/releases
[0.2.0]: https://github.com/cryptonahue/gelium-ui/releases
[0.1.0]: https://github.com/cryptonahue/gelium-ui/releases
