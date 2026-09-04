# Gelium UI — Dependency Metadata

> Grafo de dependencias del sistema Gelium UI.
> Fase J del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Responde: qué consume cada pieza (tokens/primitivas/componentes) y qué desbloquea cada pieza (p. ej. Avatar → Queue/Feed; Empty state → todas las recipes).
> Fuentes: registries (`gelium-ui-component-registry.md`, `gelium-ui-pattern-registry.md`, `gelium-ui-theme-registry.md`), `docs/gelium-ui-core.md` (token ownership), `docs/gelium-ui-screen-recipes.md`.

---

## 1. Capas del sistema y dirección de dependencia

```text
core (tokens --ui-* + base.css + tokens.css)
  └─► componentes (partial + CSS + view model Go)
        └─► patterns (composición: state D / UX E / public F)
              └─► screen recipes (composición final sobre contrato server)
```

**Regla de dirección**: core NO depende de componentes; componentes NO dependen de patterns; patterns NO dependen de recipes. Cada capa solo consume de las inferiores. Violar la dirección = un componente que reusa otro componente sin contrato (prohibido; la composición es responsabilidad de patterns/recipes).

**Stack tecnológico fijo** (`gelium-ui-core.md`, `lib/skills/14-component-implementation.md`): Go `net/http` + `html/template`, CSS propio con tokens `--ui-*`, Tailwind solo como bundler, HTMX solo como enhancement, cero JS en el flujo principal.

## 2. Dependencias por componente

> "Tokens" = familias `--ui-*` que consume (ver componente por componente en `gelium-ui-component-registry.md` §2). "Primitivas" = otros partials que renderiza internamente. "Contrato" = server contract que depende (o none).

| Componente | Tokens (familias) | Primitivas que reusa | Contrato server |
|---|---|---|---|
| Button | color, state, size, radius, focus, motion, type, border | Icon (SVG inline), spinner propio | none (form nativo); `aria-busy` loading |
| Text field | field, color, size, radius, border, state, type | Icon (error) | 422 + `X-Gelium-Validation`; valor + focus preservados |
| Data table | data-table (scoped), color, state, size, type | Pagination (footer), Empty state, Skeleton, Progress, Toast (refresh) | `GET ?q=&sort=&dir=&page=&selection=`; `HX-Request` bifurca |
| Toast | toast, color, shadow, radius, type | — | `HX-Trigger gelium:toast` (vocabulario cerrado) |
| Banner | banner, color | Button (CTA), Icon button (dismiss) | dismiss `POST + 303` |
| Empty state | empty-state, color, type | Button (CTA link) | none (server output) |
| Error state | error-state, color, type | Button (retry link) | status HTTP real (404/500/503) |
| Inline alert | inline-alert, color | Icon | — |
| Validation summary | validation-summary, color | — | 422 |
| Dialog | dialog, color, motion | Button (trigger/confirm/cancel) | `GET/POST` page variant (confirm) |
| List | list (scoped), color, size, state | Checkbox/Radio (selection), Icon | — |
| Card | card, color, shadow, radius | Button (CTA, vía CTA Link) | — |
| Avatar | avatar, color, radius, size | — | none (decorativo `aria-hidden`) |
| Pagination | pagination, color, radius, state | — | `GET ?page=` (clamping) |
| Newsletter | newsletter, color | Inline alert, Button | `POST` + 422 `X-Gelium-Validation` |
| Footer | footer, color, type, space | Language Switcher (nav secundaria) | none (slot en layout, `footerView`) |
| Breadcrumb | breadcrumb, color, type | — | none (data de `componentRoutes()`/`navLinks()`) |
| Feature card | (media aspect) + reuso card | Card, Button (CTA Link) | — |
| Hero | hero, color, type | Button (CTA) | — |
| Split | split, color, radius | — | — |
| Video | video, radius | — | none (controles nativos) |
| Resto (checkbox, radio, switch, select, slider, chips, menu, icon-button, fab, segmented-button, tabs, nav-bar, nav-tab, navigation-drawer, tooltip, progress, icon, divider, elevation, focus-ring, section-heading, language-switcher, callout, skeleton, banner) | su familia + color/state/focus/type | Icon (SVG), Button donde aplica | none o form nativo |

## 3. Dependencias por pattern

| Pattern | Depende de (componentes) | Contrato | Desbloquea |
|---|---|---|---|
| Resource list (E3) | Data table/List, Empty state, Skeleton, Pagination | GET params | Admin Resource |
| Search (E4) | GET form, Empty state, Data table/List | `GET ?q=` | Admin Resource (`?q=`) |
| Filters (E5) | Chips, Select, Segmented buttons | GET params | Admin Resource, Ops Queue |
| Pagination (E6) | Pagination standalone / data-table footer | `GET ?page=` | TODAS las recipes |
| Destructive action (E10) | Button danger, Dialog confirm | POST+303 | Admin Resource delete |
| Bulk action (E11) | Data table checkboxes, Dialog confirm | `?selection=` + POST+303 | Admin Resource |
| Confirmation (E18) | Dialog confirm | POST+303 | Admin Resource delete |
| Notifications (E15) | Toast, Banner, toast region | `gelium:toast` | todas (refresh) |
| Queue (vocabulario) | List two-line, Avatar, Badge tone, Button, Toast, Empty state, Banner | POST+303 | **Ops Queue** |
| Feed (vocabulario) | Card, Avatar, Badge, Tabs, Pagination, Skeleton, Empty state | POST+303 react + GET | **Public Feed** |
| Empty state (D1) | Empty state + Button CTA | server output | TODAS las recipes |

## 4. Qué desbloquea cada pieza (mapa inverso)

> Fila = pieza; columna de "desbloquea" = recipes/patterns que dependen de ella. Una pieza sin consumidor real está marcada "sin consumidor directo hoy".

| Pieza | Desbloquea |
|---|---|
| **Avatar** | Queue (Ops Queue), Feed (Public Feed) |
| **Empty state** | TODAS las recipes (Admin Resource 2 variantes, Ops Queue 2, Public Feed por vista) + Data table |
| **Pagination standalone** | Ops Queue, Public Feed (partial fuera del Data table) |
| **Badge tone** | Queue (tone SLA/estado), Feed (marcador "New", counts) |
| **Skeleton** | Feed (loading), loading de cualquier recipe |
| **Banner / Inline alert / Validation summary / Error state / Success feedback** | feedback de todas las recipes (422, POST+303, 404) |
| **Data table** | Admin Resource (vehículo de lista) |
| **Text field / Select / Button / Dialog** | Admin Resource (form, filtros, confirm delete) |
| **Card** | Feed (unidad), Feature card, Split, Dashboard futuro |
| **Tabs** | Feed (vistas), navegación de contexto |
| **List** | Queue (fila), Menu, Navigation, Settings futuro |
| **Breadcrumb** | Settings y páginas `/components/*` (GEO §9/§14) |
| **Footer** | TODAS las recipes (chrome) + contrato SEO §3 |
| **Hero / Split / Video / Feature card** | Landing, páginas editoriales (Phase G futura) |
| **Section heading** | encabezado de secciones de cualquier recipe |
| **Toast** | refresh remoto de todas las recipes (transitorio) |
| **Newsletter / Language switcher** | sin consumidor de recipe (bloqueadas por i18n/ejemplo) |
| **Callout** | sin consumidor directo en recipes Phase G (informativo) |
| **Tokens core (`tokens.css`)** | TODO: cada componente y theme lee `--ui-*` |
| **theme-material** | TODO: es el theme default activo |
| **theme-basecoat** | **Implementado (Phase I)** — theme completo, light + dark clase única, en bundle |

## 5. Grafo visual (recipes)

```text
                         core tokens (tokens.css + theme.css)
                                      │
   ┌──────────────┬───────────────────┼───────────────────┬───────────────────┐
 Admin Resource  │      Ops Queue     │    Public Feed    │    (recetas futuras)│
   │             │                    │                   │                    │
 Data table ── Pagination ─────────── Pagination          │                    │
 Empty state    List two-line + Avatar──► Card + Avatar + Badge                │
 Skeleton       Badge tone + Button   Tabs + Skeleton + Badge                  │
 Text field     Banner + Toast        Empty state + Toast + Banner             │
 Select/Button  Error state (404)     Pagination + Error state (404)           │
 Dialog ── confirm                    Button (like)                            │
 Banner/Inline/Validation summary ────► 422 + POST+303 + gelium:toast ────────────┘
```

## 6. Tooling: endpoint `/registry.json` (pendiente, documentado)

Se evaluó exponer el component registry como JSON servido por el handler (`GET /registry.json`) en el worktree aislado. **Decisión: NO implementado** — documentado como pendiente.

**Por qué no ahora**:

1. El registry es **documentación**, y su tabla maestra ya se genera desde el código real. Un JSON servido crearía una **segunda fuente de verdad** que debe mantenerse en Go en paralelo al markdown → riesgo de drift (el mismo modo de falla que los audits señalan repetidamente: `theme-architecture-audit.md` §8, `styles_contract_test.go` desactualizado).
2. No hay **consumidor real** del JSON hoy (ningún tooling de agentes ni dashboard lo lee).
3. Agregaría superficie de tests y de mantenimiento a un slice que es documentation-first.

**Spec del endpoint cuando haya consumidor**:

- `GET /registry.json` → `application/json; charset=utf-8`, estructura:
  ```json
  {
    "version": "1.0",
    "generated_from": "docs/gelium-ui-component-registry.md",
    "components": [
      {
        "name": "button",
        "root_class": ".ui-button",
        "template": "web/templates/button.html",
        "css": "web/styles/button.css",
        "category": "action",
        "variants": ["primary", "secondary", "outline", "text"],
        "states": ["rest", "hover", "focus", "pressed", "disabled", "loading"],
        "tokens": ["--ui-color-primary", "--ui-size-control"],
        "server_contract": "none"
      }
    ]
  }
  ```
- Fuente única: un slice Go en `internal/app/registry.go` que **sea la misma data** que renderiza el markdown (o un generador markdown→JSON en build, como `scripts/` ya hace con copy-htmx). Regla: nunca dos listas manuales de componentes en el repo.
- Tests: `TestRegistryJSON` (shape, ids únicos, categorías cerradas) + `GET /registry.json` en `postOnlyPaths()`? NO — GET es idempotente y no muta; no entra a `postOnlyPaths()`. Ruta en `sitemap`? NO — `robots: noindex` (es tooling, no página).
- **Candidato de implementación**: post-A-J (expansión opcional; no hay consumidor); requiere aprobar la fuente única de verdad antes de escribir una línea.

---

**Definición de done (Phase J)**: grafo core → components → patterns → recipes documentado, dependencias por componente con evidencia, mapa inverso de desbloqueo, y la decisión de tooling JSON tomada (pendiente con spec) y registrada.
