# Gelium UI — State Patterns Audit (Phase D, handoff)

> **Alcance**: auditoría read-only del estado real de los patrones de estado de pantalla del roadmap (`docs/gelium-ui-system-roadmap.md` Phase D) contra el código, templates, estilos y contratos server-driven existentes. No modifica código, templates, CSS, tests ni docs. Única escritura: este handoff.
>
> **Baseline**: `docs/gelium-ui-system-roadmap.md` (Phase D), `docs/gelium-ui-vocabulary.md` (§3 patrones de estado), `docs/handoffs/{composition-audit,ux-accessibility-audit,mozilla-protocol-audit}.md`, `docs/gelium-ui-composition-rules.md` (§4.8, §5, §8, §9), `docs/gelium-ui-core.md` (§6 HTML-first, §7 server-first), `internal/app/{data_table,toast,text_field,select,demo_whatsapp,server,routes}.go`, `web/templates/{data-table,toast,text-field,select,button,progress,layout,demo-whatsapp,demo-whatsapp-admin}.html`, `web/styles/{data-table,toast,text-field,demo-whatsapp}.css`, `web/static/app.js`, `COMPONENT-ROADMAP.md`, `web/styles_contract_test.go`.

---

## 1. RESUMEN — estado de los 10 patrones

| # | Patrón | Estado |
|---|---|---|
| 1 | Empty state | ◐ parcial — ad-hoc no reusable (`demo-whatsapp.html:51-53`); data table vacía sin guía (`data_table.go:239`) |
| 2 | Loading state | ✅ parcial — botón `aria-busy` (`button.html:4,9`) + `<progress>` determinate/indeterminate (`progress.html:5-23`); sin estado de carga inicial de datos |
| 3 | Skeleton | ✖ no existe (cero referencias en templates/styles) |
| 4 | Inline alert | ◐ parcial — solo error de campo (`text-field.html:5,8`, `select.html:89`); sección/form es ad-hoc |
| 5 | Banner / Notification Bar | ✖ no existe (solo ad-hoc `demo-wa-expired`, `demo-whatsapp.html:132-144`) |
| 6 | Callout | ✖ no existe (+ colisión de naming con Protocol por resolver) |
| 7 | Error state (página/recurso) | ✖ no existe; tampoco feedback de transporte HTMX (`app.js:1-9` solo 422) |
| 8 | Validation summary | ✖ no existe (solo mensajes por campo) |
| 9 | Success feedback | ✖ persistente no existe; solo toast transitorio + `role="status"` helper de campo (`text_field.go:70-71`) |
| 10 | Toast | ✅ completo (`toast.html:1-11`, `app.js:11-77`, contrato `loom:toast`, fallback no-JS) |

De 10: **1 completo**, **3 parciales**, **6 no existen**.

---

## 2. TABLA PATRONES

Leyenda: P = persistente-contextual, T = transitorio-de-acción. "JS" = enhancement requerido (0 = flujo principal sin JS).

| Patrón | Estado actual (evidencia) | P/T | Semántica propuesta | JS |
|---|---|---|---|---|
| **Empty state** | ◐ Ad-hoc: `demo-whatsapp.html:51-53` (`.demo-wa-empty` "Sin resultados…" en `<li>`, estilos `demo-whatsapp.css:199-203`). Data table: caption `0 rows · page 1 of 1` (`data_table.go:239`) + `<tbody>` vacío (`data-table.html:53-66`) sin mensaje/CTA; bug select-all `checked` con 0 filas (`data-table.html:42`, `eq .SelectedCount .Total` → `0==0`). G4 de `ux-accessibility-audit.md` | P | Reusable `.ui-empty-state`: título (`<h3>` configurable o `<p>` con `<strong>`) + `<p>` cuerpo + CTA opcional (`.ui-button` real `<a>`/`<button>`, nunca div-spand como control). En Data table: `<tr><td colspan="N">` con el empty dentro de la tabla (la fila es parte del contrato `<table>`, no un div suelto). Sin role inventado: el contenedor del listado lleva `aria-live="polite"` opcional para swaps HTMX; mensaje no-color-only | 0 |
| **Loading state** | ✅ Botón: `button.html:4,9` (`aria-busy`, `.ui-button-spinner`, sr-only "Loading {Label}"). Operación: `<progress>` nativo determinate/indeterminate (`progress.html:5-23`); refresh demo reusa `.ui-progress` determinate + toast (`data-table.html:81`, `data_table.go:354-389`). Gap: sin estado de carga inicial de listados/feed | T | Ya resuelto por contrato nativo: `aria-busy` en controles, `<progress>` en operaciones (`composition-rules.md:180`, anti-rule 8: nunca spinner ad-hoc). Para carga inicial de datos se compone con Skeleton (ver abajo) | 0 |
| **Skeleton** | ✖ No existe. Vocabulario lo marca ✖ (`gelium-ui-vocabulary.md:103-110`); state matrix GAP (`composition-rules.md:162`); cero coincidencias `skeleton` en templates/styles | P (placeholder contextual mientras carga; la UX-a11y audit lo clasifica persistente-contextual) | `.ui-skeleton` 100% CSS: contenedor con `role="status"` + texto `.sr-only` "Loading…" (existe `.sr-only` en `toast.css:70-80`) y bloques placeholder `aria-hidden="true"`; `aria-busy="true"` en la región contenedora (feed/lista). Shimmer/fade bajo `prefers-reduced-motion` desactivado. Sin animación sin JS (CSS puro); el reemplazo por datos es el siguiente request server-rendered | 0 |
| **Inline alert** | ◐ Solo campo: `text-field.html:5,8` (`aria-invalid` + `role="alert"` + `aria-describedby="{id}-error"`), `select.html:89` (`ui-select-menu-error` `role="alert"`), contrato 422 (`text_field.go:64-68`, `select.go:94-96`). Sección/form solo ad-hoc: `demo-whatsapp-admin.html:48,81` (`.demo-wa-notice` sin role), `demo-whatsapp.html:133` (`role="note"`) | P | Reusable `.ui-inline-alert--{info|success|warning|error}` ligado a sección/form: `role="alert"` para error (assertive), `role="status"` para info/success; icono `aria-hidden` + texto; opcional `aria-describedby` desde el `fieldset`/región. Genérico, reemplaza los notices ad-hoc | 0 |
| **Banner (≈ Notification Bar)** | ✖ No existe. Vocabulario lo define (`vocabulary.md:122-129`) sin componente; ad-hoc `demo-wa-expired` (`demo-whatsapp.html:132-144`) y `role="note"` funcionan como banner. Mozilla Protocol: Notification Bar ≈ Banner (`mozilla-protocol-audit.md:59,86,113`) | P | `.ui-banner--{tone}` a nivel página/sitio, sin auto-dismiss: `role="alert"` (error/acción requerida) o `role="status"` (info/success); icono `aria-hidden` + `<p>` mensaje + CTA real + dismiss `<button>` (form POST + 303, patrón WhatsApp `demo_whatsapp.go:559,573`). Ubicado al tope de `<main>`/layout. Vocabulario de tono cerrado `info|success|warning|error` (reusa contrato de Toast `toast.go:45`) | 0 (dismiss = POST + 303) |
| **Callout** | ✖ No existe. `vocabulary.md:131-138` (nota informativa ignorable). Colisión de naming: Callout Protocol (hero full-width) ≠ Callout Gelium (tip box) — `mozilla-protocol-audit.md:16,196`; resolver canónico antes de implementar | P | `<aside class="ui-callout">` (contexto complementario, `vocabulary.md:135`) con heading opcional (`<h3>`) + `<p>` cuerpo + CTA opcional; tono neutral/info preferente, variantes solo con texto (nunca color-only). Sin `role` especial (contenido estático informativo; no es alert ni status) | 0 |
| **Error state (página/recurso)** | ✖ No existe. Sin página de error custom (`server.go` usa `http.Error`); sin banner/callout de error global (`composition-audit.md:226`); sin feedback de transporte HTMX: `app.js:1-9` solo swapea 422 con header, un 500/network no muestra nada (G5 de `ux-accessibility-audit.md:88`) | P | Página/recurso: `.ui-error-state` server-rendered con `<h1>` único + `<p>` descriptivo + retry `.ui-button` real (GET al recurso) + status HTTP real (404/500/503). Fragmento/región en HTMX: el servidor re-renderiza el fragmento con `.ui-error-state`/`.ui-inline-alert` `role="alert"`; errores de transporte (500/network) quedan para Phase E (contrato de transporte, G5) — no se inventa contrato en Phase D | 0 (server-rendered) |
| **Validation summary** | ✖ No existe. Solo errores por campo; Steps lo necesita (`vocabulary.md:181`); state matrix `composition-rules.md:160` (Form: 422 inline sin resumen) | P | `.ui-validation-summary` form-level arriba del form: `role="alert"` + `<h2>`/`<h3>` ("N campos requieren atención") + `<ul>` de `<li>` con `<a href="#{campo}-error">` reales (salto al campo con error, nativo sin JS); los campos mantienen `aria-invalid` + `aria-describedby`. Focus move al summary es enhancement opcional HTMX | 0 (links reales) |
| **Success feedback** | ✖ Persistente no existe. Toast success transitorio (`toast.html`, `toast.go:108-127`); único "persistente" es el helper de campo `role="status"` (`text_field.go:70-71`). Anti-regla: feedback persistente nunca toast (`composition-rules.md:126`) | P | Reusa Banner o Inline alert con tone success: `role="status"` (polite, no alert), persistente, ligado a la sección/página. NUNCA `loom:toast` (transitorio). Post-POST+303 el redirect re-renderiza la página con el success server-rendered; en HTMX el fragmento post-submit incluye el success persistente | 0 |
| **Toast** | ✅ Completo: `toast.html:1-11` (región `#loom-toast-region aria-live="polite" aria-atomic="false" aria-relevant="additions text"`, `role="status"`/`alert` según tipo), `toast.go:13-14,45,56-61,108-127` (contrato `{"loom:toast":{...}}`, vocabulario cerrado), `app.js:11-77` (auto-dismiss 4s/8s pausable), fallback no-JS inline (`toast.go:161-164`), validación nunca toast (`toast.go:129-133`). Tests: `toast_test.go`, `styles_toast_test.go:21` | T | Ya definido; sin cambios. Único patrón transitorio-de-acción completo | J* (ya existe; fallback no-JS real) |

**Clasificación resumen** (Phase D, `gelium-ui-system-roadmap.md:168`): persistente-contextual = Empty, Skeleton, Inline alert, Banner, Callout, Error state, Validation summary, Success persistente. Transitorio-de-acción = Toast, Loading de botón/operación. **La regla a codificar: nada persistente se anuncia con `loom:toast`; nada transitorio ocupa un slot persistente.**

---

## 3. GAPS BLOQUEANTES para Phase G

| Gap | Por qué bloquea | Evidencia |
|---|---|---|
| **Empty state reusable** | Toda pantalla de admin/search/feed con listados server-side va a mostrar vacíos sin guía; hoy la tabla muestra `0 rows` como dato y "Select all" queda checked (`data-table.html:42`). Bloquea Admin Resource, Search Results, Feed, Dashboard | `data_table.go:239`, `data-table.html:42,53-66`, `composition-audit.md:224` |
| **Skeleton / Loading inicial** | Feed (recipe G) exige estado de carga de datos; hoy solo existe loading de botón y `<progress>` en operaciones puntuales. Bloquea Feed y Dashboard | `composition-audit.md:225`, state matrix `composition-rules.md:156-158` |
| **Inline alert genérico (sección/form)** | Auth (credenciales), Settings y Resource Editor necesitan error persistente a nivel form/sección; hoy solo existe error de campo + notices ad-hoc sin role | `composition-audit.md:226`, `demo-whatsapp-admin.html:48,81` |
| **Validation summary** | Formularios multi-campo (Resource Editor, Auth) sin resumen form-level; el contrato 422 existe pero no agrega resumen. Steps (Resource Editor/Booking) lo requiere | `vocabulary.md:181`, `composition-rules.md:160` |
| **Banner** | Errores globales (sesión expirada, mantenimiento) y avisos de página sin patrón; hoy caen a toast (transitorio) o `role="note"` ad-hoc | `ux-accessibility-audit.md:48`, `composition-audit.md:201` |
| **Error state (página/recurso)** | Auth flow y Resource Detail no pueden mostrar error persistente global; sin página 404/500 custom | `composition-audit.md:226` |
| **Success persistente** | Settings/editor necesitan confirmación no efímera; hoy solo toast transitorio (`composition-audit.md:227`). Menos crítico que los anteriores pero sin él las recipes de edición quedan cojas | `composition-audit.md:227` |

**No bloqueante para las 3 recipes iniciales**: Callout (nota informativa; ninguna de las 3 recipes de Phase G lo requiere; colisión de naming a resolver en Phase F). **En la frontera**: el transporte HTMX (500/network) es G5 de Phase E, pero Admin Resource con refresh remoto lo necesita; se recomienda que Error state cierre el gap de fragmento-error en Phase D y deje el transporte puro a Phase E.

---

## 4. CONTRATOS SERVER — integración SIN inventar contratos nuevos

Contratos canónicos vigentes (`gelium-ui-core.md:251-267`): (a) HTTP 422 + `X-Loom-Validation`; (b) `HX-Trigger {"loom:toast":…}`; (c) GET params estables; (d) POST + 303 redirect.

| Patrón nuevo | Contrato existente | Integración |
|---|---|---|
| **Empty state** | (c) GET params + handler existente | El handler detecta set vacío (como hoy `data_table.go:239` cuenta `total`) y el fragmento/página incluye el empty (fila `<td colspan>` en tabla, `.ui-empty-state` en lista). El empty es **output del servidor**, nunca estado cliente. HTMX: el mismo swap `outerHTML` de `data-table-panel` trae el empty (`data-table.html:8,48`). CTA opcional = link real |
| **Error state (página/recurso)** | (a) status HTTP + (c) GET | Página: status HTTP real (404/500/503) + re-render server-rendered del `.ui-error-state` (igual que `renderMarkdownPageStatus`, `server.go:155`). Recurso/fragmento: el handler detecta el error y devuelve el fragmento con el error inline `role="alert"` (mismo patrón de bifurcación `HX-Request` de `data_table.go:120-130`). Transporte puro (network/500 fuera de handler): **queda para Phase E** (G5) — en Phase D no se toca `app.js` |
| **Validation summary** | (a) 422 + `X-Loom-Validation` (sin header nuevo) | El handler de validación agrega al re-render (fragmento HX o página completa no-JS) el summary + los inline alerts por campo + links `#campo-error`. Mismo flujo que `text_field.go:55-92` y `select.go:72-124`; el hook `app.js:1-9` sigue swapeando solo 422 con header. Validación sigue sin disparar toast |
| **Success persistente** | (d) POST + 303 redirect, o fragmento post-submit | No usa `loom:toast`. No-JS: el POST hace redirect (`demo_whatsapp.go:559`) y la página destino re-renderiza el banner/inline success `role="status"`. HTMX: el fragmento post-submit incluye el success persistente (como `data-table.html:84` incluye toast inline, pero con role status persistente). La operación que falla usa 422 (validación) o error inline |
| **Banner dismiss** | (d) POST + 303 | Dismiss = form POST a ruta de dismiss + 303 (patrón WhatsApp `demo_whatsapp.go:559,573`), sin JS. Sin contrato nuevo |
| **Skeleton** | — (sin server) | Server-rendered inicial (HTML estático con `.sr-only` + bloques `aria-hidden`); el siguiente GET/POST reemplaza por datos. Sin contrato |

**Regla transversal**: ningún patrón persistente nuevo emite `HX-Trigger loom:toast`; el toast queda reservado a resultados de acción transitorios. Esto formaliza la anti-regla "feedback persistente nunca toast" (`composition-rules.md:126`).

---

## 5. ORDEN SUGERIDO de implementación (Phase D)

Basado en: bloqueo a Phase G, dependencia entre patrones, y esfuerzo (todo 100% estático salvo Toast ya entregado).

```text
1. Empty state            — desbloquea Data table (fix G4 select-all + tbody vacío) y Search/Admin.
                             Reescribe `data_table.go:239` caption → fila empty con mensaje + CTA.
2. Skeleton + Loading     — CSS puro; desbloquea Feed (recipe G) y Dashboard.
3. Inline alert genérico  — reemplaza notices ad-hoc; desbloquea Auth/Settings; reutiliza 422.
4. Validation summary     — form-level sobre el contrato 422 existente; desbloquea Resource Editor/Auth.
5. Banner                 — errores globales (sesión/mantenimiento); dismiss POST+303; desbloquea Admin Resource.
6. Error state            — página 404/500 + fragmento-error HTMX; desbloquea Resource Detail/Auth.
7. Success persistente    — reusa Banner/Inline alert con role="status"; Settings/Editor.
8. Callout                — DIFERIBLE a post-recipes; resolver colisión de naming con Protocol primero (Phase F).
```

Prioridad: **1→6 son bloqueantes de Phase G** (matriz `gelium-ui-system-roadmap.md:408-415` marca Empty/Loading/Inline/Banner/Error/Validation como bloqueantes; Callout y Success no en negrita). Success persistente puede postergarse si la primera recipe no edita.

**JS por patrón**: 0 JS en los 8 patrones nuevos (servidor + HTML/CSS nativo; dismiss y workflow por POST+303; saltos por anchors reales). Enhancement HTMX opcional: swap de fragmentos (ya existe), focus al summary al llegar un 422 con summary, `aria-live` en contenedores de listado para anunciar empty/error en swaps. Toast conserva su JS ya entregado.

**DoD Phase D adicional** (sobre `roadmap.md:172`): cada patrón nuevo requiere `styles_<pattern>_test.go` + `internal/app/<pattern>.go` + `web/templates/<pattern>.html` + entrada en `sourceAppCSS` de `styles_contract_test.go:24-58` (lista a sincronizar) + ruta en `routes.go:16-47` + `pageView` en `server.go:32-69` + actualización de `gelium-ui-vocabulary.md` (estado ✖→✅) y de la state matrix `composition-rules.md:149-162`.

---

## 6. Fuentes de autoridad usadas

`docs/gelium-ui-system-roadmap.md`, `docs/gelium-ui-vocabulary.md`, `docs/gelium-ui-composition-rules.md`, `docs/gelium-ui-core.md`, `docs/handoffs/{composition-audit,ux-accessibility-audit,mozilla-protocol-audit}.md`, `COMPONENT-ROADMAP.md`, `internal/app/{data_table,toast,text_field,select,demo_whatsapp,progress,server,routes}.go`, `web/templates/{data-table,toast,text-field,select,button,progress,layout,demo-whatsapp,demo-whatsapp-admin}.html`, `web/styles/{data-table,toast,text-field,demo-whatsapp}.css`, `themes/theme-material/theme.css`, `web/static/app.js`, `web/styles_contract_test.go`.
