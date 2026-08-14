# Gelium UI — Prompts estructurados para agentes

> Guía de trabajo para un agente (IA) en el repo Gelium UI.
> Fase J del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Principio: **referenciar los docs existentes, no duplicarlos**. Cada flujo apunta al documento canónico y extrae solo el checklist operativo.

---

## 1. Reglas de oro (en todo trabajo, sin excepción)

1. **No-JS end-to-end**: el flujo principal DEBE completarse con JS/HTMX deshabilitado (`AI-COMPONENT-IMPLEMENTER-PROMPT.md` §12). JS solo si queda una brecha de plataforma demostrada (auditoría platform-first, §7 del mismo prompt).
2. **HTML-first**: elementos nativos antes que ARIA; `div`/`span` nunca reemplazan controles (`core.md`, prompt §11).
3. **Tokens `--ui-*`**: todo valor visual público es token; cero literales de color/geometría en componentes (guard `TestNoColorLiteralsInComponents`). Los mappings Material viven en el theme, no en componentes.
4. **Server-first**: estado navegable = URL; validación = 422 + `X-Loom-Validation`; feedback persistente ≠ toast (`composition-rules.md` §9).
5. **No tocar**: `app.js`, `tokens.css`/`theme.css` (salvo theme work), `go.mod`/`package.json`, `COMPONENT-ROADMAP.md`, `MATERIAL-WEB-PROGRESS.md`, `docs/handoffs/*` — salvo tarea explícita.
6. **Worktree aislado**: trabaja en una copia física autorizada, nunca en el checkout canónico (`AI-COMPONENT-IMPLEMENTER-PROMPT.md` §4.1 `ISOLATED_PHYSICAL_WORKSPACE`).

---

## 2. Flujo: implementar un componente nuevo

> El prompt operativo completo (parámetros, concurrencia, TDD, reviews, estados finales) está en **`AI-COMPONENT-IMPLEMENTER-PROMPT.md`** (repo root). Este es el checklist condensado que ese prompt ejecuta:

1. **Descubrimiento (read-only)** — leer README, roadmap, docs de contrato (`composition-rules.md`, `gelium-ui-core.md`), componente similar existente; reportar ownership/paths protegidos (prompt §5).
2. **Auditoría platform-first** — tabla capacidad vs HTML/CSS/forms/HTMX vs baseline de browsers; JS solo con brecha demostrada (§7).
3. **Especificación antes de código** — matriz Feature × contrato upstream × estrategia × test × divergencia; divergencia no trivial requiere aprobación (§8).
4. **TDD estricto** — `NO PRODUCTION CODE WITHOUT A FAILING TEST FIRST`; ciclos verticales RED/GREEN, log por ciclo (§9).
5. **Archivos** (convención Phase D, `public-content-patterns.md:81`):
   - `web/templates/<x>.html` (`{{define "x"}}`)
   - `web/styles/<x>.css` (`@layer components`, tokens scoped en root, forced-colors)
   - `@import` en `web/styles/app.css`
   - `web/styles_<x>_test.go` + `sourceAppCSS` en `styles_contract_test.go`
   - `internal/app/<x>.go` (view model) + `<x>_test.go`
   - ruta en `routes.go` + sección en `docs.go` + `web/content/<x>.md` (si docs page)
6. **Contrato server** (si aplica) — 422/`loom:toast`/GET params/POST+303, vocabularios cerrados, view model Go concreto (nunca `map[string]any`, prompt §10).
7. **Registrar en el registry** — actualizar `docs/gelium-ui-component-registry.md` (tabla maestra) y `gelium-ui-dependency-metadata.md` (dependencias) — este es el paso que Phase J agrega.
8. **Verificación** — `npm run build`, `go test ./...`, `go vet ./...`, `go mod verify`, `node --check web/static/app.js`.
9. **Smoke en puerto propio** (nunca `:8787`, no tocar `gelium.exe`) — light/dark, narrow/wide, keyboard, no-JS, HTMX, reduced motion, forced colors (§17).
10. **Estado final** — exactamente uno de `COMPLETE_AWAITING_USER_ACCEPTANCE` / `READY_FOR_INTEGRATION` / `BLOCKED` / `ABORTED_ON_DRIFT`; entrega con checklist observable (§18-20).

**Plantilla mínima de asignación**: ver `AI-COMPONENT-IMPLEMENTER-PROMPT.md` ("Plantilla mínima de asignación", ejemplo Checkbox).

---

## 3. Flujo: implementar una screen recipe (19 campos)

> Los 3 ejemplos implementados están en **`docs/gelium-ui-screen-recipes.md`**; la plantilla de composición (19 campos) en `docs/gelium-ui-screen-composition.md` §3. Checklist:

1. **Justificar contra composition-rules** — ¿qué pattern del vocabulario cubre, con qué regla de selección (§4)? ¿qué estados de la state matrix (§8) cubre y cuáles son GAP? ¿qué anti-regla podría violar (§5)? Rationale obligatorio (`composition-rules.md` §11).
2. **Elegir componentes del registry** — SOLO componentes existentes; cero primitivas nuevas salvo aprobación (las recipes Phase G son 100% wiring; Avatar/pagination fueron la excepción aprobada con su propio contrato).
3. **Completar los 19 campos** (SURFACE, USER, PRIMARY_TASK, SECONDARY_TASKS, UX_PATTERN, VISUAL_VOCABULARY, COMPONENTS, STATES, ACCESSIBILITY, CONTENT_RULES, SEO_REQUIREMENTS, GEO_REQUIREMENTS, SERVER_CONTRACT, NO_JS_FLOW, HTMX_ENHANCEMENT, RESPONSIVE_BEHAVIOR, THEME_REQUIREMENTS, ALTERNATIVES_REJECTED, RATIONALE) — ver `gelium-ui-screen-composition.md` §3.
4. **Server contract** — GET params con vocabularios cerrados, POST+303 para mutaciones, 422 + `X-Loom-Validation` para validación, `loom:toast` solo transitorio; rutas `POST` de refresh en `postOnlyPaths()`.
5. **No-JS + HTMX** — rama no-HX completa + enhancement `hx-get`/fragmento; testea ambos por separado (`HX-Request: true`).
6. **SEO/GEO** — `/recipes/*` = `noindex, nofollow`; title/description por ruta; canonical limpio; JSON-LD cuando indexable.
7. **Archivos** — `internal/app/recipe_<x>.go` + `<x>_test.go`, `web/templates/recipe-<x>.html`, `web/styles/recipe-<x>.css` (solo layout, sin literales de color), rutas en `server.go`.
8. **Verificación** — build + test + vet + smoke en puerto propio.

---

## 4. Flujo: implementar un theme nuevo

> Procedimiento completo (12 pasos) en **`docs/gelium-ui-theme-implementation-guide.md`**; contrato en `docs/gelium-ui-theme-contract.md`; registry en `docs/gelium-ui-theme-registry.md`. Checklist condensado:

1. **Precondiciones de gate** — Phase 1 core cerrada, contrato aprobado, tests theme-agnósticos, clase data-driven, audit read-only de la referencia. Si falla: STOP, reportar `BLOCKED` (guía §0).
2. **Inspección** — audit de la referencia → `docs/handoffs/<referencia>-audit.md` + lista de tokens por componente (guía §1).
3. **Token mapping** — tabla referencia → `--ui-*`; hex por defecto; derivar elevation/motion/type/state del CSS de la referencia; marcar tokens fuera de scope (guía §2).
4. **Implementación** — `themes/<theme>/theme.css` (luz + dark, UNA rutina, sin duplicación) + import en `app.css` + selección por clase en runtime (mecanismo Phase H, 2 pasos — `gelium-ui-theme-registry.md` §4).
5. **Cobertura** — todo `var(--ui-*)` del scope definido; matriz componente × variante × estado en browser (guía §4-5).
6. **Dark/light, responsive, reduced motion/forced colors** (guía §6-8).
7. **Documentación** — `themes/<theme>/README.md` (dirección visual, mapping, divergencias, matriz) (guía §9).
8. **Tests** — suite theme-agnóstica (presencia, no valor), `styles_contract_test.go` en sync.
9. **Verificación + smoke + aceptación visual** — build/test/vet + checklist de browser (guía §10-12).
10. **Cierre** — reporte con el formato del roadmap (`STATUS: READY_FOR_INTEGRATION`, FILES, DECISIONS, TESTS, BUILD, SMOKE, EVIDENCE, RISKS).

---

## 5. Flujo: trabajo documental / registries (este slice)

1. **Fuente de verdad = código real**: la tabla de componentes se genera leyendo `web/templates/*.html`, `web/styles/*.css` y `internal/app/*.go`; no se inventa nada.
2. **Categorías cerradas** — foundation/action/input/feedback/navigation/data/public/state-pattern/recipe-primitive; los patterns (D/E/F) se indexan en `gelium-ui-pattern-registry.md`, no en el component registry.
3. **No duplicar** — cada doc referencia los contratos existentes (`composition-rules.md`, `theme-contract.md`, `theme-implementation-guide.md`, `ux-patterns.md`, `public-content-patterns.md`, `screen-recipes.md`, `AI-COMPONENT-IMPLEMENTER-PROMPT.md`) con enlaces relativos; solo lo operativo se extrae.
4. **Todo cambio documental verifica** — `npm run build`, `go test ./...`, `go vet ./...`, `git diff --check`.

---

## 6. Concurrencia y ownership (resumen)

| Modo | Quién escribe | Cuándo |
|---|---|---|
| `SHARED_HANDOFF` | solo archivos nuevos/exclusivos; patches para shared | investigación paralela |
| `ISOLATED_PHYSICAL_WORKSPACE` | copia física autorizada; nunca el canónico | varias IAs |
| `EXCLUSIVE_INTEGRATION` | shared files canónicos con reserva literal | integrador único |

Shared files por defecto (no editar sin reserva): `server.go`, `server_test.go`, `layout.html`, `app.css`, `theme.css`, `styles_contract_test.go`, `static/app.css`, `static/app.js`, `static/htmx.min.js`, `assets.go`, `README.md`, `package.json`, `package-lock.json`, `go.mod`, `go.sum`, `cmd/gelium/*` (`AI-COMPONENT-IMPLEMENTER-PROMPT.md` §4.2). Sin Git en workers (`git init/status/diff/commit` prohibidos); el integrador commitea.

---

**Definición de done (Phase J)**: los 4 flujos (componente, recipe, theme, documental) referencian los docs canónicos, extraen checklist operativos accionables y respetan las reglas no-JS/HTML-first/tokens de §1.
