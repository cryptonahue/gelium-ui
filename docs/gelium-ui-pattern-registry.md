# Gelium UI — Pattern Registry

> Inventario de **patterns** (composición) del sistema Gelium UI.
> Fase J del system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Un pattern es una **composición** de componentes reales sobre un contrato server; NO es un componente. Los registries de componentes y patterns no se mezclan (ver `gelium-ui-component-registry.md`).
> Fuentes de autoridad: `docs/gelium-ui-composition-rules.md` (reglas de selección), `docs/gelium-ui-ux-patterns.md` (Phase E, 19), `docs/gelium-ui-vocabulary.md`, `docs/gelium-ui-public-content-patterns.md` (Phase F, 14), `docs/gelium-ui-screen-recipes.md` (Phase G), `docs/handoffs/{state-patterns,public-patterns,composition}-audit.md`.

---

## 1. Capas de patterns

| Capa | Fase | Nº | Ejemplos |
|---|---|---|---|
| State patterns | D | 8 | Empty state, Skeleton, Inline alert, Banner, Callout, Error state, Validation summary, Success feedback |
| UX patterns | E | 19 | Resource list, Search, Destructive action, Confirmation, Notifications… |
| Public/content patterns | F | 14 | Article, Hero, Breadcrumb, Footer, Newsletter, Split… |

Los state patterns tienen partial real (ver component registry, categoría SP). Los public patterns tienen partial real (categoría P). Los UX patterns son composiciones documentadas; su readiness está marcado en `gelium-ui-ux-patterns.md` (**Ready** vs **Phase G**).

---

## 2. State patterns (Phase D — 8)

> Regla rectora (`state-patterns-audit.md:45`): **persistente-contextual** nunca viaja por `loom:toast`; **transitorio-de-acción** nunca ocupa un slot persistente.

| # | Pattern | Componente(s) | Cuándo usarlo (referencia) | Consumido por | Ejemplo real |
|---|---|---|---|---|---|
| D1 | Empty state | `empty-state.html` | Todo listado server-side sin datos; mensaje + CTA real; NUNCA flash de loading (`composition-rules.md` §8 GAP cerrado) | Admin Resource (2 variantes), Ops Queue (2), Public Feed (por vista), Data table (`<td colspan>`) | `data-table.html:68-70`, recipes |
| D2 | Skeleton | `skeleton.html` | Carga inicial de regiones de datos; `role="status"` + sr-only "Loading" | Public Feed (placeholder documentado), Admin Resource (refresh) | `skeleton.css` |
| D3 | Inline alert | `inline-alert.html` | Error de validación 422 junto al campo; advertencia de sección persistente | Admin Resource (form), Newsletter, Data table refresh | `inline-alert--error` en recipes |
| D4 | Banner | `banner.html` | Aviso persistente nivel página/sitio que exige acción; success persistente post-303 | Admin Resource (post-create/delete), Ops Queue (post-advance/dequeue) | `banner--success`/`banner--error` |
| D5 | Callout | `callout.html` | Contenido informativo/ignorable sin urgencia | (ninguna recipe Phase G; documental) | tip box `<aside>` |
| D6 | Error state | `error-state.html` | Recurso inexistente o error global; status HTTP real + retry GET | Admin Resource (404), Ops Queue (404), Public Feed (404) | `recipe_*.go` 404 |
| D7 | Validation summary | `validation-summary.html` | Errores 422 multi-campo con links a `#campo-error` | Admin Resource (form), Newsletter (email) | `validation-summary.css` |
| D8 | Success feedback | reusa `inline-alert--success` / `banner--success` | Confirmación NO efímera de operación exitosa (POST+303 → página re-renderiza) | TODAS las recipes | `banner--success` post-303 |

Nota: **Toast** es el mecanismo transitorio-de-acción (contrario directo de D8); vive en el component registry (categoría B) y en UX pattern #15 (Notifications).

---

## 3. UX patterns (Phase E — 19)

> Documentación completa de cada pattern (problem/user/context/happy/empty/loading/error/recovery/mobile/a11y/server contract/when-not) en `docs/gelium-ui-ux-patterns.md`. Aquí solo el índice de composición.

| # | Pattern | Compuesto por | Cuándo usarlo (ref. composition-rules) | Consumido por | Estado |
|---|---|---|---|---|---|
| 1 | Authentication | Text field, Inline alert, Validation summary, Banner, Button, 422 | pantalla completa, nunca dialog | — | Phase G |
| 2 | Onboarding | Steps (gap), Callout, Text field, Banner, 422 | multi-paso | — | Phase G |
| 3 | Resource list | Data table o List, Empty state, Skeleton, Pagination, GET params | set de registros (`composition-rules.md` §4.1) | **Admin Resource** | Ready |
| 4 | Search | GET form, Empty state, Data table/List, Skeleton | find by free text | Admin Resource (`?q=`) | Ready |
| 5 | Filters | Chips, Select, Segmented buttons, GET params | narrow by category | Admin Resource, Ops Queue (`?status=&kind=`) | Ready |
| 6 | Pagination | Pagination standalone / data-table footer, GET params | navegar páginas | **TODAS las recipes** | Ready |
| 7 | Empty state | Empty state | (== D1) | todas | Ready |
| 8 | Loading | Button `aria-busy`, Progress, Skeleton | operación en curso | todas | Ready |
| 9 | Error recovery | Error state, Inline alert, Validation summary, 422 | fallo entendible + retry | todas | Ready* |
| 10 | Destructive action | Button (danger), Dialog confirm, POST+303 | irreversible | Admin Resource (delete) | Ready |
| 11 | Bulk action | Data table checkboxes, Dialog confirm, POST+303 | actuar sobre múltiples | Admin Resource (`?selection=`) | Ready* |
| 12 | Multi-step form | Steps (gap), Progress, Validation summary, 422 | input complejo | — | Phase G |
| 13 | Checkout | Steps (gap), Text field, 422, Validation summary | purchase | — | Phase G |
| 14 | Booking | Steps (gap), date input, 422, Success feedback | reserva | — | Phase G |
| 15 | Notifications | Toast, Banner, toast region `aria-live` | resultado transitorio vs persistente | Admin Resource (refresh), Ops Queue, Public Feed | Ready |
| 16 | Settings | Panel, List, Text field, Banner/Inline success, 422 | configuración | — | Phase G |
| 17 | Permissions | List + checkboxes, Segmented, Dialog, 422 | quién puede qué | — | Phase G |
| 18 | Confirmation | Dialog confirm/cancel, `closedby="any"` | confirmar acción consecuente | Admin Resource (delete) | Ready |
| 19 | Undo / recovery | Toast action, 422 value preservation, POST+303 | revertir error | — | Phase G |

`*` Ready a nivel de primitiva; la recipe vive en Phase G (`ux-patterns.md:38`).

---

## 4. Public/content patterns (Phase F — 14)

> Regla: ninguno exige JS para su esencia; secciones NO son theme; no portar naming `mzp-*`. Estado completo en `docs/gelium-ui-public-content-patterns.md` (13/14 ✅, Card slots ◐).

| # | Pattern | Compuesto por | Tier | Desbloquea |
|---|---|---|---|---|
| F1 | Article | `.prose` (base.css) + intro opcional | 1 | contrato tipográfico de cualquier página |
| F2 | Billboard/Hero | Hero (`h1` display-lg + subtitle + CTA Button link + media con scrim) | 2 | Landing / Public Feed |
| F3 | Breadcrumb | `breadcrumb.html` | 1 | GEO §9/§14 (BreadcrumbList JSON-LD) |
| F4 | Callout | `callout.html` (Phase D) | — | (== D5) |
| F5 | Card (slots públicos) | `card.html` + media/tag/meta/CTA | 2 | Admin Resource, Public Feed (◐ pendiente) |
| F6 | CTA Link | **Button link** (`button.html` con `Href`) — no componente propio | 1 | Empty state / Banner / Callout / Hero / Feature Card / Split |
| F7 | Feature Card | Card + media + CTA Link | 2 | páginas editoriales |
| F8 | Footer | `footer.html` + slot en `layout.html` + `footerView` | 3 | TODAS las recipes (chrome) |
| F9 | Language Switcher | `language-switcher.html` (GET form, submit visible) | 3 | i18n futura |
| F10 | Newsletter | `newsletter.html` + handler (POST + 422) | 3 | conversión |
| F11 | Notification Bar | alias documental de Banner | 4 | (== D4) |
| F12 | Section Heading | `section-heading.html` (siempre `h2`) | 1 | encabezado de secciones |
| F13 | Split | `split.html` (grid 2 col, stack narrow) | 2 | Landing / Public Feed |
| F14 | Video | `video.html` (`<video controls>` + captions) | 1 | media promocional |

---

## 5. Cuándo usar qué pattern (resumen ejecutivo)

La tabla de decisión completa vive en `composition-rules.md` §4 (table vs list, list vs queue, queue vs board, card vs panel, feed vs collection, timeline vs activity list, dialog vs page, toast vs inline alert vs banner vs callout). Anti-reglas duras en `composition-rules.md` §5:

- ≤5-8 filas estáticas → List, nunca Data table.
- Cola FIFO estricta → Queue, nunca Board.
- Validación → 422 + Inline alert, NUNCA toast.
- Feedback persistente/crítico → Inline alert/Banner, nunca toast.
- Estados navegables → URL (GET params), nunca estado cliente.
- Flujo largo → Page/Steps, nunca Dialog.

---

**Definición de done (Phase J)**: los 41 patterns (8 state + 19 UX + 14 public) indexados con componentes reales, consumidores reales (recipes) y referencia a la regla de selección; sin colisión con el component registry.
