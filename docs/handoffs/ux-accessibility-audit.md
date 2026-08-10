# Gelium UI — UX & Accessibility Audit (handoff)

> **Alcance**: investigación read-only sobre UX y accesibilidad del estado real del catálogo. No modifica código, templates, CSS ni tests. Única escritura: este handoff.
>
> **Baseline**: `README.md`, `COMPONENT-ROADMAP.md`, `docs/gelium-ui-core.md`, `docs/gelium-ui-vocabulary.md`, `docs/gelium-ui-composition-rules.md`, `docs/handoffs/composition-audit.md`, `internal/app/*`, `web/templates/*.html` (todos), `web/styles/*`, `web/static/app.js`, `themes/theme-material/theme.css`.
>
> **Lectura recomendada en paralelo**: `docs/handoffs/composition-audit.md` (§5 estados, §7 anti-rules) y `docs/gelium-ui-core.md` (§7 contratos server-driven, §9 accesibilidad).

---

## 0. Veredicto general

El sistema es **estructuralmente sano**: elementos nativos antes que ARIA, estado nunca color-only, `:focus-visible` con contrato único, forced colors y reduced motion casi completos, sin roles falsos (`tablist`/`dialog`/`menu` falsos no existen). Los gaps fuertes están en **patrones de estado de pantalla** (Empty/Skeleton/Banner/Callout/Validation summary no existen), en **fallback real de overlays** (Dialog/Select menu dependen de Invoker Commands sin fallback entregado) y en **feedback de errores de transporte** (inexistente). Los problemas más graves file:line están en los **demos** (`demo-whatsapp*.html`), no en las primitives.

---

## 1. State patterns

### 1.1 Clasificación del feedback (persistente-contextual vs transitorio-de-acción)

| Capa | Persistencia | Define | Dónde vive hoy |
|---|---|---|---|
| **Persistente-contextual** | Sobrevive a la interacción, ligado al contexto | Inline alert, Banner, Callout, Error state, Validation summary, Empty state, Skeleton | Casi todo FALTA (solo inline field error y empty ad-hoc) |
| **Transitorio-de-acción** | Resultado efímero de una acción, auto-dismiss | Toast, Loading (button/operation) | Toast ✅ completo, Loading ✅ parcial (sin Skeleton) |

### 1.2 Tabla por patrón

| Patrón | Tipo | Estado | Evidencia |
|---|---|---|---|
| Empty state | persistente-contextual | **◐ parcial (ad-hoc, no reusable)** | Único: `demo-whatsapp.html:51-53` (`demo-wa-empty`, `Sin resultados…`). Data table: tbody vacío + caption `0 rows` sin mensaje/CTA (`data_table.go:239`, `data-table.html:53-66`). Vocabulario lo marca GAP (`gelium-ui-vocabulary.md:94-101`, `composition-audit.md:87`). |
| Loading (botón) | transitorio-acción | ✅ existe | `button.html:4,9` (`aria-busy`, nombre dinámico `Loading {Label}`, spinner). |
| Loading (operación) | transitorio-acción | ✅ existe | `progress.html:5-23` `<progress>` nativo determinate/indeterminate; refresh `data-table.html:81`. |
| Skeleton | persistente-contextual | ✖ no existe | Vocabulario: "Skeleton no existe" (`gelium-ui-vocabulary.md:103-110`). Ningún skeleton en templates. |
| Inline alert (campo) | persistente-contextual | ✅ existe (solo campo) | `text-field.html:5,8` (`aria-invalid` + `role="alert"` + `aria-describedby`), `select.html:89` (`ui-select-menu-error` `role="alert"`), contrato 422 (`text_field.go:64-68`). |
| Inline alert (sección/form) | persistente-contextual | ✖ no existe como componente | Solo ad-hoc: `demo-whatsapp-admin.html:48,81` (`demo-wa-notice--warn`, `demo-wa-notice`), `demo-whatsapp.html:133` (`role="note"`). |
| Banner | persistente-contextual | ✖ no existe | Vocabulario lo define (`gelium-ui-vocabulary.md:122-129`) sin componente. |
| Callout | persistente-contextual | ✖ no existe | `gelium-ui-vocabulary.md:131-138`. |
| Error state (página/recurso) | persistente-contextual | ✖ no existe | `composition-audit.md:226` (Auth flow / resource detail no pueden mostrar error persistente global). |
| Validation summary (form-level) | persistente-contextual | ✖ no existe | No hay resumen de errores a nivel form; solo mensajes por campo. El contrato Steps lo necesita (`gelium-ui-vocabulary.md:181`). |
| Success feedback | persistente-contextual | ✖ no existe | Toast es transitorio; no hay confirmación no efímera (`composition-audit.md:227`). |
| Toast | transitorio-acción | ✅ completo | `toast.html:1-11` (región `aria-live="polite"`, `role="status"/"alert"`), `app.js:17-18,41-76` (auto-dismiss 4s/8s pausable), contrato `loom:toast` (`toast.go:108-127`). |

### 1.3 Qué se confunde hoy

- **Validación ≠ toast**: regla vigente y testeada (`toast.go:129-133`, `COMPONENT-ROADMAP.md:49`) — se respeta. El riesgo real es el inverso: al no existir Banner/Callout/Validation summary, **no hay patrón persistente para errores de página o éxito guardado**, lo que empuja a usar toast para feedback que debería sobrevivir.
- **Transitorio disfrazado de persistente**: en el flujo no-JS el toast inline "persistente" de la demo (`toast.html:21`, `data-table.html:84`) se mezcla con el `role="status"` de los notices (`data-table.html:21`, `chips.html:63`) — tres mecanismos con apariencia similar pero semánticas distintas (alert/status/region), sin guía de composición que los distinga en pantalla.
- **Empty state ad-hoc**: `demo-wa-empty` (`demo-whatsapp.html:51-53`) es el único empty del sistema y no es reusable; la tabla muestra `0 rows` como si fuera un dato (ver §6).
- **`role="note"` de ventana expirada** (`demo-whatsapp.html:133`) funciona como Banner/Callout ad-hoc sin componente detrás.

---

## 2. Keyboard

No hay focus management JS; todo el teclado es el contrato nativo. Verificación por control:

| Control | Teclado | Estado |
|---|---|---|
| Dialog (`dialog.html:3`) | Focus trap, Escape y light dismiss nativos (`closedby="any"`). | ✅ nativo, PERO la **apertura depende de `command="show-modal"`** (`dialog.go:17-19`, `button.html:9`), no-Baseline. En Firefox/Safari el trigger es un botón muerto (ver §8-G1). |
| Select menu (`select.html:77-88`) | `<dialog>` nativo; items submit buttons; Tab navega, Escape/light dismiss cierran. | ✅ nativo, PERO mismo problema de `command` (G1). |
| Menu popover (`menu.html:44-127`) | Popover nativo: Escape + light dismiss; items son buttons/links/checkbox/radio reales; **Tab** navega (sin roving focus, aceptable porque no reclama `role="menu"`). | ✅ nativo. Faltan `aria-expanded` en triggers (G8). |
| Tabs (`tabs.html:6-18`, `tabs.go:8-14`) | Links reales + `aria-current`; sin arrows, sin `tablist` (decisión documentada). | ✅ cumplido; el teclado de links (Tab/Enter) es suficiente. |
| Segmented buttons (`segmented-button.html:10,50`) | Radios/checkboxes nativos en `fieldset` → **flechas funcionan nativamente**. | ✅ |
| Radio/Checkbox/Switch/Slider/Select nativos | Contrato nativo completo (flechas en radio/select/range). | ✅ |
| Tooltip (`tooltip.css:52-56`) | Reveal por `:hover`/`:focus-within`. | ✅ foco teclado revela. |
| Chips input remove (`chips.html:57`) | Botón real. | ✅ |
| Data table (`data-table.html:40-72`) | Sort/filter/pagination = links reales; checkboxes nativos. | ✅ |
| FAB/Icon button/Button disabled (`fab.html:5,15`, `icon-button.html:4`, `button.html:4`) | `tabindex="-1"` + `aria-disabled`. | ✅ |

**Conclusión teclado**: el flujo de teclado está resuelto por nativos, sin roving focus falso. El único hueco real de teclado es G1 (overlays sin fallback en navegadores no-Chromium).

---

## 3. Focus

- **`:focus-visible` global** en `focus-ring.css:7-14` + reglas por componente + override forced-colors (`app.css:122`). Ningún elemento pierde indicador. ✅
- **Overlays**:
  - Dialog: al abrir, el navegador enfoca el primer foco/`autofocus` (Cancel lleva `autofocus`, `dialog.go:18`); **focus trap y retorno de foco son nativos de `showModal`**. ✅ (gated por G1).
  - Menu popover: entrada/salida de foco nativa del popover; **sin `aria-expanded` en el trigger** (G8).
- **Recuperación en 422**: `autofocus` de vuelta al campo con error en no-JS (`text_field.go:67`, `toast.go:149`); se omite en rama HX para no robar foco. ✅ patrón ejemplar.
- **Gap**: no hay `aria-expanded` ni gestión de foco explícita para popover; `focus-ring.html:5` enseña un `<span tabindex="0" role="link">` que es enfocable pero **no activable** (sin handler Enter) — demo de mal patrón.

---

## 4. Errors y recovery

- **Contrato 422 + `X-Gelium-Validation`**: `text_field.go:64-91`, `select.go:94-114`, `toast.go:147-179`; hook `htmx:beforeSwap` solo swapea 422 con header (`app.js:1-9`). El valor se preserva (`text_field.go:62`) y el foco vuelve al campo (`text_field.go:67`). ✅
- **Reporte al usuario**: `aria-invalid="true"` + `aria-describedby` + mensaje persistente `role="alert"` (`text-field.html:5,8`); success/helper con `role="status"` (`text_field.go:70-71`). ✅
- **Errores de red/remotos**: **NO EXISTEN**. El hook `app.js:1-9` solo trata 422 con header; un 500/network en HTMX **no muestra nada** (sin `hx-on::response-error`, sin región de error de transporte, sin retry). El refresh de data table ante fallo no tiene feedback (G5). La regla "422 ≠ error de transporte" documentada (`composition-audit.md:115`) protege el swap, pero no define qué pasa cuando el transporte falla.
- **Recovery de campo**: buena (re-submit re-renderiza; helper reemplaza error). No hay clear-on-input (intencional, server-driven). ✅
- **Bugs concretos**:
  - `POST /demo/whatsapp/admin` **no está registrado** (`server.go:83` solo `GET`); el form de webhook (`demo-whatsapp-admin.html:72-83`) recibe **405** sin feedback (G3).
  - Estados `aria-invalid` en checkbox/radio (`checkbox.html:33`, `radio.html:35`) **sin mensaje ni descripción** (solo demo visual).

---

## 5. ARIA

### 5.1 De más (redundante/incorrecto)

- `role="list"` en `<ul>` (`demo-whatsapp.html:30`) — redundante.
- `aria-label` **pisando** el `<label>` visible: filter `aria-label="Filter by name or status"` vs visible "Filter" (`data-table.html:9-10`); slider `aria-label` vs caption visible (`slider.html:5,11,17`); progress `aria-label="Progress 30 percent"` vs caption "30%" (`progress.html:5-23`). SR anuncia un nombre ≠ texto visible.
- `aria-selected` sobre `<button>` (`select.html:84`) — **inválido para `role=button`** (aria-selected es de listbox/option/tab). Menor.
- `title` duplicando info visible (`demo-whatsapp.html:44,148`).
- `<span tabindex="0" role="link">` no activable (`focus-ring.html:5`).
- `role="link"` redundante en `<a>` (`button.html:4`, `fab.html:5,15`, `icon-button.html:4`) — intencional en el patrón disabled-link, inofensivo.

### 5.2 De menos

- `aria-expanded` en triggers de popover/menu/select (`menu.html:44,75,90`, `select.html:78`).
- `aria-current` en tabs del admin (`demo-whatsapp-admin.html:23-27`, solo clase `--active`).
- `aria-label` en links-acción placeholder con nombre emoji (`demo-whatsapp-admin.html:42,62,118`).
- `lang="es"` en demos en español (ver G2).
- Skip link / bypass blocks (G7).
- `role="status"`/`aria-live` en notices ad-hoc (`demo-whatsapp-admin.html:48,81`).
- `aria-atomic`/nombre en `role="note"` de ventana expirada (`demo-whatsapp.html:133`).

### 5.3 Falso roles: **ninguno**. No hay `tablist`, `dialog` falso ni `menu` falso — Dialog es `<dialog>`, Menu es `ul` de controles nativos, Tabs son links. Esta es la mayor fortaleza ARIA del sistema y debe preservarse.

---

## 6. Flows (cuellos de botella UX)

1. **Data table sin empty state** (G4): con `?q=zzz` la caption dice `0 rows · page 1 of 1` (`data_table.go:239`) y tbody vacío (`data-table.html:53-66`); además el checkbox "Select all" queda `checked` cuando 0 filas (`eq .SelectedCount .Total` → `0==0`, `data-table.html:42`). Sin guía para el usuario.
2. **Demo WhatsApp como única screen** (`demo-whatsapp.html`, `demo-whatsapp-admin.html`): todo lo demás es galería de componentes. El admin es una segunda pantalla, pero con **interacciones muertas**: `href="#"` en "conectar", ⚙, 🗑, ✎, 👁, copiar, rotar (`demo-whatsapp-admin.html:31,42,52,62,79,104,118`) y el form de webhook **inoperable** (G3). El botón "regenerar ⚠" (`demo-whatsapp-admin.html:79`) es `type="button"` sin acción.
3. **Sin confirmación de éxito persistente**: "Guardar webhook" no llega a ejecutarse; el patrón general de acciones (`POST + 303`) no deja feedback de éxito cuando el resultado no pasa por toast (`demo_whatsapp.go:544-585`).
4. **Sin estado de carga para datos**: feeds/dashboards futuros no tienen Skeleton ni error de página (`composition-audit.md:57-58,226`).

---

## 7. Reduced motion y forced colors

- **Reduced motion**: bloque central en `app.css:52-69` (button, text-field, dialog, toast, elevation, switch, select, select-menu, slider, progress, fab, list) + bloques por componente (tabs, navigation-bar/tab/drawer, segmented-button, icon-button, tooltip, menu, chips, data-table, demo-whatsapp). **Gap**: `checkbox.css:26` y `radio.css:26` (transform `scale(.92)` al activar) **no se desactivan** bajo reduced motion — menor.
- **Forced colors**: bloque central `app.css:71-213` + por componente (menu, chips, data-table, navigation-*, segmented-button, icon-button, tooltip, tabs, demo-whatsapp). Cobertura prácticamente completa de todos los componentes (checkbox/radio/card/badge/divider/elevation/focus-ring cubiertos en el bloque central).
- **Contraste**: tokens light/dark cumplen AA por diseño (p. ej. `--ui-color-fg-muted #49454f` sobre canvas `#fff7ff`; danger `#b3261e` sobre canvas ≈ 8:1). Los tones del demo (warning `#6a4b00` sobre `#fff3d6`, quality GREEN/YELLOW/RED `demo-whatsapp.css:501-503`) también cumplen. No se detectaron fallos de contraste AA salvo verificación visual pendiente.
- **Drift conocido**: `demo-whatsapp.css:403` referencia `var(--ui-color-error, #b3261e)` — token inexistente (`--ui-color-error` vs `--ui-color-danger`), el fallback hardcodeado siempre aplica (ya documentado en `gelium-ui-core.md:104`).

---

## 8. Gaps de accesibilidad concretos (priorizados)

| # | Severidad | Gap | Evidencia |
|---|---|---|---|
| G1 | **Crítico** | Overlays sin fallback: Dialog y Select menu se abren solo con `command`/`commandfor` (Invoker Commands, no-Baseline). En Firefox/Safari el trigger es un botón muerto; el comentario en `select-menu.css:2-8` promete "native `<select>` fallback" que **no se renderiza** (el form `select.html:74-92` solo tiene trigger + dialog + hidden). | `dialog.html:3`, `dialog.go:17-19`, `select.html:77-88`, `button.html:9`, `select-menu.css:2-8` |
| G2 | **Alto** | `lang="en"` en contenido 100% en español: SR lee español con voz inglesa. | `demo-whatsapp.html:3`, `demo-whatsapp-admin.html:3` |
| G3 | **Alto** | Form de webhook muerto: `POST /demo/whatsapp/admin` no registrado → **405** sin feedback. Además 4+ links `href="#"` placeholder y botón "regenerar" inerte. | `server.go:83` vs `demo-whatsapp-admin.html:72-83`; `demo-whatsapp-admin.html:31,42,52,62,79,118` |
| G4 | **Alto** | Data table sin empty state (0 filas = caption + tbody vacío, sin mensaje/CTA; "Select all" checked con 0 filas). | `data_table.go:239`, `data-table.html:42,53-66` |
| G5 | **Alto** | Sin feedback de errores de red/500 en HTMX (sin `hx-on::response-error`/región de transporte/retry). | `app.js:1-9`, `data_table.go:354-389` |
| G6 | **Medio** | Nombres accesibles desincronizados: `aria-label` pisa el `<label>` visible (filter, slider, progress). | `data-table.html:9-10`, `slider.html:5,11,17`, `progress.html:5-23` |
| G7 | **Medio** | Sin skip link ("saltar al contenido"); landmark `main` ausente en el admin; heading/landmark del shell dependen de la demo. | `layout.html:11-16`, `demo-whatsapp-admin.html:21-28` |
| G8 | **Medio** | `aria-expanded` ausente en triggers popover/menu/select. | `menu.html:44,75,90`, `select.html:78` |
| G9 | **Medio** | Admin tabs activos sin `aria-current` (clase-only) y links-acción emoji sin `aria-label`. | `demo-whatsapp-admin.html:23-27,42,62,118` |
| G10 | **Bajo** | `aria-selected` inválido sobre `<button>` en Select menu; `role="list"` redundante; botón Dismiss de toast inline no-JS sin handler (`data-loom-toast-dismiss` nunca se escucha en `app.js`); tooltip rich con action link no clickeable por mouse (`pointer-events:none`). | `select.html:84`, `demo-whatsapp.html:30`, `toast.html:5` vs `app.js:41-76`, `tooltip.css:43` |
| G11 | **Bajo** | Checkbox/radio sin reduced-motion; drift `--ui-color-error`. | `checkbox.css:26`, `radio.css:26`, `demo-whatsapp.css:403` |

---

## 9. Recomendaciones de prioridad (para fases siguientes)

1. **G1 antes que todo**: entregar fallback server-rendered real para Dialog/Select menu (p. ej. variante "page"/"details" o `<select>` nativo como fallback real) — respeta el contrato no-JS end-to-end del core (`gelium-ui-core.md:11`).
2. **State patterns Phase 4**: Empty state reusable (desbloquea data table y search), Banner/Callout/Validation summary (desbloquean Auth/Settings), Skeleton.
3. **Contrato de error de transporte** como tercer pilar server-driven (hoy solo 422 y toast).
4. Corregir demos: `lang`, webhook 405, placeholder links, `aria-expanded`, empty de tabla.

---

**Fuentes de autoridad usadas**: `README.md`, `COMPONENT-ROADMAP.md`, `docs/gelium-ui-core.md`, `docs/gelium-ui-vocabulary.md`, `docs/gelium-ui-composition-rules.md`, `docs/handoffs/composition-audit.md`, `internal/app/{server,text_field,toast,data_table,select,dialog,tabs,routes,demo_whatsapp}.go`, `web/templates/*.html` (todos), `web/styles/*.css`, `web/static/app.js`, `themes/theme-material/theme.css`.
