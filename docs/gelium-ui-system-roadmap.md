# Gelium UI — System Roadmap

> Roadmap superior del sistema Gelium UI. Define la arquitectura que convierte la colección de componentes Material en un sistema UI reutilizable para distintos proyectos server-rendered.
>
> Este documento es el roadmap de SISTEMA. El roadmap específico de componentes continúa en `COMPONENT-ROADMAP.md`; no se reemplaza ni se modifica aquí.
>
> **Nota de naming**: los documentos del sistema usan el prefijo `gelium-ui-*` (renombrado desde `loom-ui-*` por decisión del mantenedor). `loom-ui` es solo la ruta física del repositorio, no branding.

## Objetivo

Gelium UI =

```text
un core UI universal
+ componentes semánticos
+ patrones UX
+ screen recipes
+ contratos server-driven
+ accesibilidad
+ SEO
+ GEO
+ themes intercambiables
```

No es solamente una colección de componentes Material: es un sistema UI reutilizable. El segundo theme, **Basecoat UI**, está implementado como dirección visual (`themes/theme-basecoat/theme.css`, Phase I). Shadcn no se implementa todavía como runtime ni como port directo; más adelante puede convertirse en un preset o dirección visual "shadcn-like".

## Estado actual del sistema (baseline)

| Área | Estado | Evidencia |
|---|---|---|
| Core Go + Markdown + embed | Completado | `cmd/gelium`, `internal/app`, `web/assets.go` |
| Build Tailwind CSS 4 + HTMX local | Completado | `package.json`, `web/static/` |
| Theme Material light/dark | Completado | `themes/theme-material/theme.css` |
| Theme Basecoat light/dark (Phase I) | Completado | `themes/theme-basecoat/theme.css`, `themes/theme-basecoat/README.md` |
| Componentes Material (contract Gelium) | 20+ entregados | `COMPONENT-ROADMAP.md`, `web/templates/`, tests |
| Tokens públicos `--ui-*` | Parcial; no formalizado como contrato | `themes/theme-material/theme.css` |
| Core agnóstico (desacoplado de Material) | **No existe todavía** | — |
| Theme contract | **No existe todavía** | — |
| Patrones de estado de pantalla | **No existen** (solo ad-hoc) | `demo-whatsapp.html:51-53`, `data_table.go:239` |
| UX/accessibility/SEO/GEO contracts | **No existen** | `layout.html:3-10` (head mínimo) |
| Registry / tooling para agentes | No existe todavía | — |

---

## Decisiones de las auditorías (baseline auditado)

Resultado de la investigación paralela (`docs/handoffs/{core,vocabulary,composition,theme-architecture,basecoat,ux-accessibility,seo-geo,mozilla-protocol}-audit.md`). Estas decisiones quedan fijadas:

### Core / themes (core-audit, theme-architecture-audit)

1. **No existe core agnóstico todavía**: el theme Material define ~157 tokens y 8 componentes (List, Menu, Data table, Navigation bar/tab/drawer, Segmented, Tooltip) declaran tokens scoped que el theme no sobrescribe. La **token ownership** (core vs theme vs component vs pattern) es el primer trabajo real.
2. **Familias de foundation a crear**: spacing, density, sizes, border, colores semánticos complementarios, tipografía composable, motion completo, focus. **Tratar con cuidado**: breakpoints y z-index — no convertirlos automáticamente en tokens públicos; usar top layer para dialog/popover; solo escalas con consumidor real.
3. **Acoplamientos Material a tokenizar**: state-layer colors hardcodeados (`button.css:17-18`, `icon-button.css:22,25`, `chips.css:32,35,107,141-142`) → `color-mix()` theme-aware; geometría px → tokens de anatomía del core; clase `m3-select-trigger` → `ui-select-trigger`; `class="theme-material"` hardcodeado en `layout.html:2` → theme identity data-driven.
4. **Gaps de tokens a cerrar**: `--ui-color-surface-container`, `--ui-type-display-lg`, `--ui-type-title-md`, `--ui-font-mono`; unificar `--ui-color-error`/`--ui-color-danger` → `danger`; eliminar tokens muertos (`--ui-radius-xl`, `--ui-state-dragged-opacity`, `--ui-select-menu-item-icon`) o mantenerlos por compatibilidad si hay consumidores.
5. **Dark mode duplicado con drift** (`theme.css:203-251` vs `253-299`): unificar en una sola rutina; Basecoat valida el patrón de clase única.

### Vocabulario / composición (vocabulary-audit, composition-audit)

6. **Naming resuelto**: Popover = mecanismo web (no patrón); Multi-select = capacidad (no widget); Drawer → canónico "Navigation drawer"; Snackbar → alias de Toast; **Callout** (Gelium = tip box) y **Notification Bar** (≈ Banner) requieren resolución contra Protocol (colisión detectada por mozilla-protocol-audit).
7. **Del vocabulario objetivo (29 términos)**: 8 ya implementados, 8 parciales, 13 sin equivalente. Los **patrones de estado** (Empty, Loading/Skeleton, Inline alert, Banner, Callout, Error state, Validation summary, Success feedback) son requisito bloqueante de toda screen recipe.
8. **Contratos server-driven canónicos** (no inventar otros): (a) HTTP 422 + `X-Loom-Validation`; (b) `HX-Trigger {"loom:toast":…}`; (c) GET con params estables; (d) POST + 303 redirect.
9. **Anti-reglas confirmadas**: validación nunca toast; no toast para feedback persistente; no table para ≤5-8 filas; no board para FIFO; no dialog para flujos largos; estado nunca color-only; URL es el estado.

### UX / accesibilidad (ux-accessibility-audit)

10. **Base accesible sólida**: HTML nativo antes que ARIA, cero roles falsos, `:focus-visible` global, forced colors casi completo, teclado resuelto por contratos nativos, contrato 422 con autofocus de recuperación ejemplar.
11. **Gaps críticos**: overlays (Dialog/Select menu) sin fallback real en no-Chromium (trigger `command`/`commandfor` no-Baseline = botón muerto en Firefox/Safari); `lang="en"` en demos en español; form webhook a ruta inexistente (405); data table sin empty state; sin feedback de errores de red/500 en HTMX.
12. **Confusión de patrones confirmada**: al no existir Banner/Callout, el feedback persistente cae a toast (transitorio) o roles ad-hoc. **Feedback persistente-contextual ≠ feedback transitorio-de-acción** (Phase D debe formalizarlo).

### SEO / GEO (seo-geo-audit)

13. **Base SEO sólida** (server-rendered, `<title>` por ruta, h1 único, HTML semántico, URLs limpias, sin CDN) pero **cero contrato de metadata**: falta meta description, canonical, OG, Twitter, robots, sitemap, JSON-LD; falta `<footer>`; caching débil sin gzip.
14. **GEO**: la fundación es citable (URLs estables, contenido factual) pero invisible para citación enriquecida: identidad de marca dividida (Gelium UI / Gelium UI / LoomChat), cero machine-readable metadata, sin provenance/autoría en docs, headings no interrogativos, sin breadcrumb.

### Public content / Mozilla Protocol (mozilla-protocol-audit)

15. **Protocol es referencia, no runtime**: MPL 2.0, Sass+JS npm, grid propietario y naming `mzp-*` NO se copian. Ninguno de los 14 patrones exige JS para su esencia → valida el principio no-JS.
16. **Clasificación**: public/content patterns son capa separada de componentes de aplicación (criterio: presentación/conversión vs workflow con estado) y NO son themes.
17. **Mapping**: Card ✅, CTA ≈ Button link ✅, Article ≈ `.prose` ✅; los 11 restantes son gaps. Colisiones de naming: "Callout" (Protocol hero ≠ Gelium tip box) y "Notification Bar" (≈ Banner Gelium) a resolver.

### Basecoat (basecoat-audit)

18. **Basecoat es dirección visual válida, NO runtime**: no se instala `basecoat-css`; el theme es `themes/theme-basecoat/theme.css` traduciendo el style pack (default Vega) al vocabulario `--ui-*`.
19. **Mapeo 5/7 limpio** (Button, Card, Badge, Dialog, Data table); **2 requieren decisión**: Text field (floating label = variante de theme/pattern, no requisito universal) y Toast (contrato `loom:toast` no-JS se conserva; solo estética).
20. **Basecoat no tokeniza** elevation/spacing/motion/typescale/state-opacities → derivar del style pack; oklch → hex.

---

## Fases — orden por dependencias arquitectónicas

> El orden anterior (0–8 lineal) se reemplaza por estas fases A–J. La razón: **no construir un theme nuevo sobre una arquitectura que todavía no distingue core / theme / component / pattern / screen / public content**.

### Phase A — Contract and architecture

Cerrar el contrato antes de tocar foundations. Decisiones obligatorias:

- **Token ownership**: qué tokens son core, cuáles theme, cuáles component, cuáles pattern; quién declara y quién sobrescribe (incluye los 8 componentes con tokens scoped).
- **Theme identity**: `class="theme-material"` sale de `layout.html:2`; selección data-driven (server o documento raíz).
- **No hardcoded Material classes**: renombrar `m3-select-trigger` → `ui-select-trigger`; auditar prefijos no-`ui-*`.
- **Semantic color vocabulary**: roles canónicos (canvas/surface/primary/secondary/danger/warning/success/info/outline/scrim/focus-ring + `-fg`).
- **Dark-mode routine única**: eliminar duplicación clase + media query (drift `--ui-switch-track-unselected`).
- **State layers theme-aware**: `color-mix()` sobre token `-fg`; eliminar `rgb()` fijos.
- **Remove/rename broken token references**: `surface-container`, `display-lg`, `title-md`, `font-mono`, `error`/`danger`.
- **Decidir qué tokens son públicos** (contrato), **cuáles internos** (core) y **cuáles muertos** se eliminan o requieren compatibilidad.

**Salida**: `docs/gelium-ui-core.md` actualizado (capítulo Token Ownership) + fixes en `web/styles/`.

**DoD**: token ownership documentado; theme identity data-driven; cero clases Material hardcodeadas; dark unificado; state layers theme-aware; sin referencias rotas; `npm run build` + `go test ./...` + `go vet ./...` verdes; Material sigue pasando smoke.

### Phase B — Theme-neutral foundations

Implementar las foundations con prioridad:

- spacing (`--ui-space-*`);
- density / sizes (`--ui-density-*`, `--ui-size-*`);
- borders (`--ui-border-*`);
- semantic colors complementarios (success/warning/info/outline/scrim);
- typography decomposition (size/weight/line-height/letter-spacing por step);
- motion completo (short/medium/long + easings);
- focus;
- reduced motion;
- forced colors.

**Cuidado**:

- **breakpoints y z-index NO se convierten automáticamente en tokens públicos**; usar **top layer** para dialog/popover; definir solo las escalas con un **consumidor real** (evitar `--ui-z-*`/`--ui-breakpoint-*` sin uso).

**Salida**: tokens core en `web/styles/` + actualización de `docs/gelium-ui-core.md`.

**DoD**: cada familia nueva tiene consumidor real en componentes; sin tokens huérfanos; build+test+smoke verdes; light/dark, reduced motion, forced colors.

### Phase C — Theme-agnostic verification

Antes de cualquier theme nuevo, la suite debe ser theme-agnóstica:

- parametrizar tests (13 hardcodean ruta `theme-material/theme.css`);
- **no asertar hexes Material** (3 tests aseveran valores: `styles_fab_test.go:85-87`, `styles_dialog_test.go:15-29`, `styles_toast_test.go:37-42`);
- asertar **presencia de contratos** (familias de tokens, cobertura light/dark), no valores;
- verificar que **Material sigue pasando**;
- preparar la misma suite para Basecoat;
- revisar `styles_contract_test.go` (lista `sourceAppCSS` desactualizada);
- verificar **generated assets** (`web/static/app.css` embebido).

**Salida**: suite theme-agnóstica + `docs/gelium-ui-theme-contract.md` actualizado.

**DoD**: suite pasa sin depender de hexes Material; tests de contrato (presencia) para todos los componentes; generated assets verificados.

### Phase D — Universal state patterns

Definir e implementar los contratos de estado:

- Empty state;
- Loading state;
- Skeleton;
- Inline alert;
- Banner;
- Callout;
- Error state;
- Validation summary;
- Success feedback;
- Toast.

**No confundir**:

```text
persistent contextual feedback  ≠  transient action feedback
```

(Persistente-contextual: Empty, Inline alert, Banner, Callout, Error state, Validation summary, Success persistente. Transitorio-de-acción: Toast, Loading de botón/operación.)

**Salida**: componentes en `web/templates/` + `web/styles/` + `docs/gelium-ui-vocabulary.md` actualizado.

**DoD**: cada patrón con contrato (semántica HTML, estados, accesibilidad, server contract), sin ambigüedad persistente/transitorio, build+test+smoke verdes.

### Phase E — UX, accessibility, content, SEO and GEO

Crear los contratos e integrarlos (no documentos aislados):

- `docs/gelium-ui-ux-principles.md`;
- `docs/gelium-ui-ux-patterns.md`;
- `docs/gelium-ui-content-rules.md`;
- `docs/gelium-ui-accessibility-contract.md`;
- `docs/gelium-ui-seo-contract.md`;
- `docs/gelium-ui-seo-patterns.md`;
- `docs/gelium-ui-geo-contract.md`;
- `docs/gelium-ui-geo-patterns.md`.

**Integración obligatoria**: referenciarlos desde `composition-rules.md` y `screen-recipes.md`; implementar metadata server-driven (title/description/canonical/OG/robots/JSON-LD) en el handler Go + layout; unificar identidad de marca; fix de accesibilidad (fallback de overlays, `lang`, webhook 405, empty state, errores de red).

**DoD**: contratos escritos Y referenciados desde composition rules; metadata server-driven funcionando; fixes de accesibilidad aplicados; Material + smoke verdes.

### Phase F — Public content patterns

Fase inspirada en Mozilla Protocol, sin copiar branding. Clasificar como **public/content patterns, NO como theme**:

- Article;
- Billboard/Hero;
- Breadcrumb;
- Callout (resolver colisión de naming con Gelium);
- Card (slots públicos media/tag/CTA);
- CTA Link;
- Feature Card;
- Footer;
- Language Switcher;
- Newsletter;
- Notification Bar (≈ Banner Gelium);
- Section Heading;
- Split;
- Video.

**Reglas**: ninguno exige JS para su esencia (Footer → `<details>`, Language Switcher → GET form, Newsletter → POST+422); no portar grid `mzp-*` ni naming `mzp-*`; Feature Card horizontal descartada (deprecada upstream).

**Salida**: `docs/gelium-ui-public-content-patterns.md` + componentes donde aplique.

**DoD**: patrones clasificados public/content (no theme), mapping a Gelium documentado, colisiones de naming resueltas, build+test+smoke verdes.

### Phase G — First screen recipes

Implementar primero **solo**:

1. Admin Resource;
2. Ops Queue;
3. Public/Social Feed.

Las demás quedan después: Resource Detail, Resource Editor, Dashboard, Settings, Booking, Search Results, Authentication.

Cada recipe incluye:

```text
surface
user
primary task
UX pattern
visual vocabulary
components
states
accessibility
content rules
SEO requirements
GEO requirements
server contract
no-JS flow
HTMX enhancement
responsive behavior
theme requirements
alternatives rejected
rationale
```

**Salida**: `docs/gelium-ui-screen-recipes.md` (solo las 3 primeras).

**DoD**: 3 recipes completas con los 17 campos, usando componentes reales, con server contract y fallback no-JS; smoke en browser.

### Phase H — Theme mechanism

Implementar **solo el mecanismo mínimo**:

```text
<body class="theme-material">
<body class="theme-basecoat">
```

- La selección debe poder venir del **servidor** o del **documento raíz**.
- **No agregar todavía un runtime de themes basado en JavaScript**.
- Bundle de themes en `web/static/app.css` (único asset embebido) + selección por clase.

**Salida**: mecanismo mínimo funcionando (Material sigue siendo el default).

**DoD**: cambiar la clase del theme cambia la dirección visual sin rebuild ni JS; ambos themes coexisten en el bundle; tests verdes.

### Phase I — Basecoat

Implementar después de las fases anteriores:

```text
Button
Text field
Card
Badge
Dialog
Toast
Data table
```

- Basecoat debe **mapearse al contrato Gelium** (`themes/theme-basecoat/theme.css`).
- **No crear una segunda API** ni reemplazar la anatomía HTML universal sin justificación.

**Decisiones obligatorias**:

- **floating label de Text field**: tratarlo como **variante de theme o pattern**, no como requisito universal;
- **Toast**: conservar `loom:toast`, `aria-live` y fallback no-JS;
- verificar **state layers, borders, density y radii**;
- **auditar cualquier JavaScript requerido**.

**Salida**: `themes/theme-basecoat/theme.css`, `themes/theme-basecoat/README.md`.

**DoD**: matriz theme × component × variant × state ejecutada con Material y Basecoat; tests theme-agnósticos pasando; smoke completo; aceptación visual manual.

### Phase J — Registry

Dejar para el final:

- component registry;
- pattern registry;
- theme registry;
- dependency metadata;
- prompts para agentes;
- screen composition output.

**Regla**: no agregar un registry hasta que los patterns tengan contratos estables.

**Salida**: registries + prompts.

---

## Dependencias entre fases

```text
Phase A (contract/architecture)
  └─► Phase B (theme-neutral foundations)
        └─► Phase C (theme-agnostic verification)
              └─► Phase D (universal state patterns)
                    └─► Phase E (UX/a11y/content/SEO/GEO)
                          ├─► Phase F (public content patterns)
                          └─► Phase G (first screen recipes)
                                └─► Phase H (theme mechanism)
                                      └─► Phase I (Basecoat)
                                            └─► Phase J (registry)
```

**Invariante**: Phase A es la raíz; nada se implementa antes del contrato. Basecoat (I) nunca antes de A–H.

## Trabajo paralelo vs serial

### Puede ser paralelo (investigación/auditoría, read-only)

- Auditorías de core/tokens, UX/accessibility, SEO/GEO, Mozilla Protocol, Basecoat (ya entregadas en `docs/handoffs/`).
- Preparación de contratos (esqueleto de docs) en copias aisladas.
- Reviews read-only.

### Debe ser serial

- Escritura del core y token ownership (Phase A).
- Foundations (Phase B).
- Parametrización de tests (Phase C).
- Implementación de state patterns (Phase D).
- Contratos UX/a11y/SEO/GEO y su integración (Phase E).
- Screen recipes (Phase G).
- Mecanismo de themes (Phase H) y Basecoat (Phase I).
- Build, tests, smoke y release.
- Integración al checkout canónico (un integrador único; workers usan `SHARED_HANDOFF` o copias aisladas).

## Riesgos

| Riesgo | Mitigación |
|---|---|
| Core se contamina con estilos Material | Phase A fija token ownership antes de extraer |
| Tokenizar todo automáticamente (breakpoints/z-index sin consumidor) | Phase B: solo escalas con consumidor real; top layer para overlays |
| Theme nuevo sobre arquitectura sin contrato | Fases A–H son gate duro antes de I |
| Overlays rotos en no-Chromium (fallback `command`/`commandfor`) | Fix de accesibilidad en Phase E; fallback server-rendered |
| Feedback persistente cae a toast | Phase D formaliza persistente vs transitorio |
| Identidad de marca dividida (Gelidium/Gelium/LoomChat) | Phase E unifica naming en metadata y contenido |
| Tests dependientes de hexes Material | Phase C parametriza; aserciones de presencia |
| Colisión de naming public patterns (Callout, Notification Bar) | Phase F resuelve contra vocabulario Gelium |
| Basecoat como segunda API | Phase I: solo tokens sobre contrato Gelium; markup intacto |
| Registry prematuro | Phase J: solo después de contratos estables |

## Artefactos por fase

| Fase | Artefacto |
|---|---|
| A | `docs/gelium-ui-core.md` (Token Ownership), fixes `web/styles/` |
| B | tokens core en `web/styles/`, `docs/gelium-ui-core.md` |
| C | suite theme-agnóstica, `docs/gelium-ui-theme-contract.md` |
| D | componentes de estado + `docs/gelium-ui-vocabulary.md` |
| E | 8 docs de contratos + metadata server-driven + fixes a11y |
| F | `docs/gelium-ui-public-content-patterns.md` + componentes |
| G | `docs/gelium-ui-screen-recipes.md` (3 recipes) |
| H | mecanismo `<body class="theme-*">` |
| I | `themes/theme-basecoat/theme.css` + README |
| J | registries + prompts |

Handoffs: `docs/handoffs/{core,vocabulary,composition,theme-architecture,basecoat,ux-accessibility,seo-geo,mozilla-protocol}-audit.md`.

## Matriz de gaps (inventario de planificación A–J)

> **Histórica + residual.** La mayoría de filas de A–J ya se implementaron en código. No reescribas esta tabla fila a fila como board vivo: usá la sección **Estado actual** y los registries. Abajo solo quedan los **residuals abiertos** que todavía valen trabajo.

| Residual abierto | Tipo | Evidencia / notas |
|---|---|---|
| Dark-mode routine única (posible drift clase vs media) | polish | Ambos themes aún declaran bloque `.theme-*.theme-dark` **y** `@media (prefers-color-scheme: dark)` |
| Familias scoped fuera del theme (List/Menu/Data table/Nav*/Segmented/Tooltip) | architecture | Un theme nuevo debe declarar esos tokens globalmente si quiere pintarlas |
| SEO origin + `og.png` placeholder | product | `siteBaseURL = https://gelium-ui.example`; no hay `web/static/og.png` real |
| Nav discoverability | DX | **Cerrado (docs-shell):** `docsNavFor` + two-pane shell sidebar/topbar on `/docs` + `/components/*` expose Getting started, `docsSections` component groups, Patterns (`/docs/patterns`), Recipes (outbound `/recipes/*`), Themes (`/docs/themes`); theme switcher stays 0-JS `?theme=`. Demos remain secondary (`/demo/*`, noindex) — not a chrome gap. |
| Docs pages para Avatar, Pagination, state/public patterns | dogfood | Partials+CSS+tests sí; `/components/*` no (policy o páginas) |
| Branding operativo residual | docs | Wire `loom:*` / `X-Loom-*` **congelado** (estrategia en `gelium-ui-wire-compatibility.md`); branding humano (LICENSE, prompts, COMPONENT-ROADMAP, README) ya migrado a Gelium |
| Recetas restantes (Detail/Editor/Dashboard/Settings/Booking/Search/Auth) | recipe | Diferidas a propósito post-3 |
| Registry JSON servido | tooling | Registries son Markdown; runtime JSON documentado como pendiente |
| Version bump past 0.4.0 narrative | release | `package.json` + `?v=0.4.0` en layout |
| a11y demo leftovers (G10 / `href="#"` scaffold) | polish | `docs/gelium-ui-accessibility-contract.md`; demo admin |

**Criterio original (archivado)**: "Bloquea" = sin este ítem, Basecoat o una screen recipe no se podían completar. Ese gate **ya se cruzó** para las 3 recipes y Basecoat.

### Inventario histórico A–J (referencia)

| Gap | Tipo | Fase | Bloquea Basecoat | Bloquea recipe | Acción |
|---|---|---|---|---|---|
| Token ownership (core vs theme vs component vs pattern) | contract | A | **Sí** | Sí | Definir propiedad de cada familia (incluye 8 componentes scoped) |
| Theme identity data-driven (`layout.html:2`) | contract | A | **Sí** | No | Clase del theme desde server/raíz |
| Clases Material hardcodeadas (`m3-select-trigger`) | cleanup | A | Sí | No | Renombrar a `ui-select-trigger` |
| Semantic color vocabulary (roles canónicos) | contract | A | **Sí** | Sí | Fijar roles + `-fg` |
| Dark-mode routine única (drift) | contract | A | **Sí** | No | Unificar clase + media query |
| State layers theme-aware (`color-mix`) | foundation | A | **Sí** | No | Reemplazar `rgb()` fijos |
| Referencias rotas (`surface-container`, `display-lg`, `title-md`, `font-mono`, `error/danger`) | cleanup | A | **Sí** | No | Definir/renombrar/unificar |
| Tokens muertos (`radius-xl`, `state-dragged-opacity`, `select-menu-item-icon`) | cleanup | A | No | No | Eliminar o compat |
| Spacing (`--ui-space-*`) | foundation | B | **Sí** | Sí | Crear escala y migrar |
| Density / sizes (`--ui-density-*`, `--ui-size-*`) | foundation | B | **Sí** | Sí | Crear tokens y migrar alturas |
| Borders (`--ui-border-*`) | foundation | B | Sí | No | Crear tokens |
| Colores complementarios (success/warning/info/outline/scrim) | foundation | B | Sí | Sí | Crear tokens con consumidor |
| Typography decomposition | foundation | B | **Sí** | No | Descomponer `--ui-type-*` |
| Motion completo (medium/long + easings) | foundation | B | Sí | No | Crear escala; reemplazar literales |
| Focus / reduced motion / forced colors | foundation | B | No | No | Centralizar (existe parcial) |
| Breakpoints / z-index | foundation | B | No | No | **NO tokenizar**; top layer; solo con consumidor |
| Tests parametrizados (ruta theme) | test | C | **Sí** | No | Parametrizar 13 tests |
| Tests sin hexes Material | test | C | **Sí** | No | Aserciones de presencia |
| `styles_contract_test.go` desactualizado | test | C | No | No | Sincronizar lista `sourceAppCSS` |
| Generated assets verificados | test | C | No | No | Verificar `web/static/app.css` |
| Empty state | component | D | No | **Sí** | Crear componente reusable |
| Loading state / Skeleton | component | D | No | **Sí** | Crear componente |
| Inline alert (sección) | component | D | No | **Sí** | Crear componente genérico |
| Banner | component | D | No | **Sí** | Crear componente |
| Callout | component | D | No | Sí | Crear componente (resolver naming Protocol) |
| Error state (página/recurso) | component | D | No | **Sí** | Crear contrato |
| Validation summary | component | D | No | **Sí** | Crear contrato |
| Success feedback persistente | component | D | No | Sí | Crear contrato (no confundir con toast) |
| Fallback overlays no-Chromium (Dialog/Select) | component | E | No | Sí | Fallback server-rendered real |
| `lang="en"` en demos español | cleanup | E | No | No | Corregir demos |
| Form webhook 405 | bugfix | E | No | No | Corregir ruta o remover |
| Errores de red/500 en HTMX | component | E | No | Sí | Feedback de transporte |
| Metadata server-driven (description/canonical/OG/robots) | contract | E | No | Sí | Contrato SEO + handler |
| JSON-LD / structured data | contract | E | No | Sí | Contrato GEO |
| Sitemap.xml / robots.txt | contract | E | No | No | Generar server-side |
| `<footer>` faltante | component | E | No | No | Crear en layout |
| Identidad de marca unificada | documentation | E | No | No | Unificar Gelium UI |
| Provenance/autoría en docs | documentation | E | No | No | Fechas/autor en content |
| UX principles / patterns | documentation | E | No | Sí | Contratos referenciados |
| Content rules | documentation | E | No | Sí | Contrato de contenido |
| Accessibility contract | documentation | E | No | **Sí** | Contrato referenciado |
| SEO/GEO patterns | documentation | E | No | Sí | Contratos referenciados |
| Article / Billboard-Hero / CTA Link / Section Heading / Split / Video | pattern | F | No | Sí | Public/content patterns (no theme) |
| Breadcrumb | pattern | F | No | Sí | Crear (`<nav>`+`<ol>`) |
| Footer / Language Switcher / Newsletter | pattern | F | No | Sí | Alternativas no-JS (details/GET/POST+422) |
| Feature Card / Notification Bar | pattern | F | No | Sí | Composición Card+CTA; ≈ Banner |
| Admin Resource recipe | recipe | G | No | — | Implementar con estado completo |
| Ops Queue recipe | recipe | G | No | — | Composición List+Badge+POST |
| Public/Social Feed recipe | recipe | G | No | — | Composición Card/List+Skeleton |
| Recetas restantes (Detail/Editor/Dashboard/Settings/Booking/Search/Auth) | recipe | G | No | — | Diferidas post-3 |
| Theme mechanism (`<body class="theme-*">`) | theme | H | **Sí** | No | Mecanismo mínimo sin JS |
| Basecoat theme | theme | I | — | No | Solo tokens sobre contrato Gelium |
| Floating label Text field | theme | I | — | No | Variante de theme/pattern, no universal |
| Registry (component/pattern/theme) | registry | J | No | No | Solo tras contratos estables |

**Criterio**: "Bloquea" = sin este ítem, Basecoat o una screen recipe no se pueden completar correctamente. Cleanup/test/polish no se marcan como bloqueantes de recipe salvo que la evidencia lo justifique.

## Verificación (cada fase que modifique código)

```bash
npm run build
go test ./...
go vet ./...
```

Además verificar siempre:

- Material funciona;
- tests no dependen de hexes Material;
- light/dark;
- narrow/wide;
- reduced motion;
- forced colors;
- teclado;
- no-JS;
- HTMX;
- empty/loading/error;
- smoke browser en el puerto permitido;
- consola sin errores.

**Un componente o theme no está terminado sin aceptación visual manual.**

## Definition of done por fase

| Fase | DoD |
|---|---|
| A | Token ownership documentado; theme identity data-driven; cero clases Material hardcodeadas; dark unificado; state layers theme-aware; sin referencias rotas; build+test+smoke verdes |
| B | Cada familia nueva con consumidor real; sin tokens huérfanos; build+test+smoke verdes |
| C | Suite theme-agnóstica (sin hexes Material); tests de contrato; generated assets verificados |
| D | Patrones de estado con contrato; persistente vs transitorio sin ambigüedad; build+test+smoke verdes |
| E | 8 contratos escritos Y referenciados; metadata server-driven; fixes a11y aplicados |
| F | Public/content patterns clasificados (no theme); colisiones resueltas; build+test+smoke verdes |
| G | 3 recipes completas con los 17 campos; server contract y fallback no-JS; smoke |
| H | Mecanismo mínimo: cambiar clase cambia theme sin rebuild ni JS |
| I | Matriz theme × component × variant × state; tests theme-agnósticos; aceptación visual manual |
| J | Registries funcionales; prompts estructurados |

---

## Próximo slice de implementación (post A–J)

Phases **A–J están entregadas en el árbol** (plus docs-shell chrome). El trabajo nuevo ya no es “desbloquear Basecoat”; es pulido, productización y expansión opcional:

```text
1. Truth sync          README · este roadmap · theme registry · cmd/gelium
2. DX / discoverability  ✅ docs-shell: nav Recipes+Themes+Patterns+grouped components (demos still secondary)
3. SEO productization    BASE_URL configurable + og.png real (cuando haya dominio)
4. Theme polish          dark routine única; ownership de familias scoped
5. Release               bump de versión past 0.4.0 + push main (hoy ahead de origin)
6. Optional expansion    más screen recipes · tercer theme · registry JSON runtime
```

**No reabrir A–H** salvo regresión. Wire protocol `loom:toast` / `X-Loom-Validation` permanece congelado.
