# Gelium UI — G1 Overlays Audit (Invoker Commands / no-Chromium fallback, handoff)

> **Alcance**: auditoría **read-only** del gap G1 de accesibilidad (overlays Dialog + Select menu) en Gelium UI. No modifica código, templates, CSS, tests ni docs. Única escritura: este handoff.
>
> **Baseline**: `docs/gelium-ui-accessibility-contract.md` (§1.8, §1.9, §3 G1), `docs/gelium-ui-system-roadmap.md` (Phase E, fila `Fallback overlays no-Chromium`), `web/templates/{dialog,select,button}.html`, `internal/app/{dialog,select,button,server}.go`, `web/styles/select-menu.css`, `web/content/{dialog,select,navigation-drawer}.md`, `README.md` (§ Dialog open-code), tests `internal/app/{dialog,select}_test.go`, `web/styles_{dialog,select_menu}_test.go`.
>
> **Fecha**: 2026-08-10. Compatibilidad verificada contra MDN BCD y caniuse (jul-2026).

---

## 1. RESUMEN EJECUTIVO

1. **Dialog y Select menu abren HOY únicamente vía Invoker Commands** (`command`/`commandfor` + `<dialog closedby="any">`). No existe ningún fallback server-rendered en el markup, ni detección de soporte (no hay `@supports` para command ni hook JS; `app.js` solo maneja el swap htmx 422).
2. **`command`/`commandfor` ya NO es Chromium-only**: MDN lo marca **Baseline 2025 — Newly available (dic 2025)**. Los docs lo describen como "Baseline Low" (obsoleto). El botón muerto queda para versiones NO actuales (Chrome ≤133, Safari ≤24/25, Firefox ≤152), no para "Firefox/Safari" en general.
3. **`closedby` sigue sin Baseline** (caniuse 71.54%): **Safari 26.4–27 NO lo soporta** (solo TP); Chrome/Edge 134+ y **Firefox 141+ sí**. En Safari `closedby="any"` se ignora y el dialog se comporta como `"closerequest"` (Escape y Cancel siguen cerrando; el light-dismiss no). Degradación elegante, no botón muerto.
4. **El sub-gap REALMENTE crítico es el Select menu**: `select-menu.css:2-8` y `select.go:22-24` **prometen** "el `<select>` nativo queda como fallback", pero **ese `<select>` NO se renderiza** (`select.html:73-92`). Sin Invoker Commands el form entero queda muerto: trigger inerte + dialog cerrado para siempre + único submit dentro del dialog. Es un incumplimiento docs-vs-código, peor que el Dialog (que al menos está documentado como gap aceptado).
5. **Recomendación**: fix **ahora, en Phase E**, en un slice enfocado. Select menu → **`<select>` nativo real** como control base (es el propio componente Select del sistema). Dialog → **variante página/confirmación** (trigger = `<a>` a una ruta que renderiza la misma acción inline), que es literalmente la "page/detail variant" que el contrato §1.9 ya nombra; el `<dialog>`+command queda como mejora progresiva opt-in, sin dejar controles inertes en ninguna página Gelium.

---

## 2. ESTADO EXACTO POR OVERLAY

### 2.1 DIALOG — `web/templates/dialog.html`, `internal/app/dialog.go`

**Cómo se abre hoy** (markup dogfooded en `/components/dialog`):

- Trigger: `{{template "button" .Trigger}}` (`dialog.html:2`) → `button.html:9` renderiza `<button … command="show-modal" commandfor="confirm-dialog">` (`dialog.go:17`, `Command/CommandFor` en `buttonView` `button.go:13-14`).
- `<dialog id="confirm-dialog" closedby="any" aria-labelledby="…-title" aria-describedby="…-description">` (`dialog.html:3`).
- Cancel: `command="request-close"` + `autofocus` (`dialog.go:18`). Confirm: `command="close"` `value="confirm"` (`dialog.go:19`).
- **Confirm NO ejecuta ninguna acción server**: es puramente presentacional (cierra el dialog, no hay POST, no hay ruta). El "Confirm action" es un demo de confirmación vacío.

**Qué pasa sin Invoker Commands (no-Chromium / versiones no actuales)**: el `<button>` es un botón muerto (`type="button"`, sin handler). El contenido del dialog (Cancel/Confirm) es **inalcanzable**. No hay `<details>`, no hay link, no hay URL, no hay adaptador JS.

**Estado documental**: gap **aceptado y documentado** — `dialog.md:7,9`, `README.md:120`, contrato §1.8 (`:165`), §1.9 (`:182`), §3 G1 (`:347`). Es decir: para Dialog, G1 es una limitación conocida y explícita; el incumplimiento real está en Select menu.

### 2.2 SELECT MENU — `web/templates/select.html:73-92`, `internal/app/select.go`

**Cómo se abre hoy**:

- Form `POST /examples/select/menu` con `hx-post`/`hx-target="this"`/`hx-swap="outerHTML"` (`select.html:74-76`).
- Trigger: `<button type="button" class="m3-select-trigger" command="show-modal" commandfor="select-menu" aria-haspopup="menu">` (`select.html:77-78`).
- `<dialog class="ui-select-menu" id="select-menu" closedby="any" aria-label="Plan options">` con items `type="submit" name="value"` (`select.html:82-88`) + hidden `id`/`value` (`:90-91`).
- `selectMenuDemo.Open` controla **solo** el render de `aria-selected` dentro del dialog siempre-cerrado (`select.html:84`, `select.go:12-14,111`); en el load inicial `Open=true` → marca `aria-selected` en un dialog cerrado (suma a G10).

**Qué pasa sin Invoker Commands**: trigger inerte + dialog cerrado para siempre. **El único camino de submit del form está DENTRO del dialog.** Resultado: **el form completo es un botón muerto**. Y contradice los comentarios:

- `select-menu.css:2-8` — "without command support the dialog stays closed so the native `<select>` field remains the fallback control".
- `select.go:22-24` — idéntica promesa.

**Ninguno de los dos renderiza un `<select>`.** La promesa de fallback es falsa. Este es el hallazgo central del audit.

### 2.3 NAVIGATION DRAWER modal (mismo patrón, fuera del scope pedido)

`navigation-drawer.html:6-12` + `navigation_drawer.go:121-126`: trigger `command="show-modal" commandfor="navigation-drawer-modal"` sobre `<dialog closedby="any">`. Mismo gap, bien documentado (`navigation-drawer.md:92-96,125-127`). Comparte el arquetipo; el fix de Dialog debería decidir si lo incluye o lo deja como contrato documentado.

### 2.4 Resumen

| Overlay | Apertura | Fallback server hoy | Form funcional sin JS | Severidad G1 |
|---|---|---|---|---|
| Dialog | `command`/`commandfor` | **Ninguno** | Parcial (Confirm es presentacional) | Alta (documentada) |
| Select menu | `command`/`commandfor` | **Ninguno (prometido, no renderizado)** | **No** — form muerto | **Crítica** |
| Navigation drawer modal | `command`/`commandfor` | Ninguno | N/A (destinos son links reales) | Alta (documentada) |

---

## 3. COMPATIBILIDAD REAL (verificada 2026-08, no "caniuse-like" sino datos actuales)

| Feature | Estado Baseline | Soporte motor | Consecuencia para Gelium |
|---|---|---|---|
| `<dialog>` element | **Baseline Widely available** (mar-2022) | todos | El elemento en sí no es el problema |
| `command`/`commandfor` (Invoker Commands) | **Baseline 2025 — Newly available (dic-2025)** — MDN. **NO es "Baseline Low 2024"** | Chrome 134+ (mar-2025), Safari 25/26 (2025), Firefox 153 (2026) | Botón muerto SOLO en versiones no actuales (Chrome ≤133, Safari ≤24, Firefox ≤152). Premisa "Firefox/Safari" ya no es categórica |
| `closedby` | **No Baseline** — caniuse global **71.54%** | Chrome/Edge **134+**, **Firefox 141+** (sí, Firefox ya lo tiene), **Safari 26.4/26.5/27: NO** (solo TP) | En Safari se ignora → el dialog abierto con `showModal()` se comporta como `"closerequest"`: Escape + Cancel cierran, light-dismiss no. Degradación elegante |
| `request-close` | Más nuevo que `show-modal`/`close` (`dialog.md:7`) | — | Cancel depende de la versión más nueva |

**Detección en el sistema hoy: NINGUNA.**

- `rg "@supports" web internal themes` → solo `anchor-name` (`menu.css:62`, `select-menu.css:26`). No hay `@supports` para command.
- `app.js:1-9` → únicamente swap htmx de fragmentos 422 con `X-Loom-Validation`. Sin feature-detection.
- **Punto arquitectónico clave**: `@supports` de CSS **no puede gatear atributos HTML** — Invoker Commands los consume el UA, no el CSS. Y la detección por JS está excluida por la regla "no componente JavaScript". **Conclusión: dentro de esta arquitectura un overlay condicionado a detección es imposible; el camino base DEBE ser universalmente funcional (navegación o form real).** Esto valida la premisa del G1 y descarta cualquier "imitar overlay con CSS".

---

## 4. OPCIONES DE FALLBACK Y RECOMENDACIÓN

### 4.1 DIALOG

Opciones del enunciado:

- **(a) El dialog se convierte en página/URL** (confirmación como página de detalle).
- **(b) El trigger es un link a una ruta que renderiza la misma acción inline.**
- **(c) `<details>` nativo como degradación** → **RECHAZADA**: semántica de disclosure, no de modal; sin focus trap ni top-layer; sin contrato Escape/cancel; es una "imitación" de overlay, exactamente lo que el sistema prohíbe. Además el contenido de un `<details>` cerrado es inaccesible hasta expandir.

**Recomendación: (a) + (b) convergen en "confirmación-como-página"** — literalmente la *"page/detail variant"* que el contrato §1.9 ya nombra como fallback esperado de G1.

Diseño coherente con la arquitectura (no-JS, server-rendered, base = navegación real):

1. **Base universal**: el trigger se vuelve un **link real** (`<a class="ui-button" href="/components/dialog/confirm">`) → navega → una ruta server renderiza el mismo título/descripción **inline** como contenido normal de página, con Confirm = **form POST real** (+303 back) o link `?confirmed=1`, y Cancel = link de vuelta. Funciona en todo navegador, 0 JS, 0 CSS imitation.
2. **Mejora progresiva opt-in**: el `<dialog>` + `command` + `closedby="any"` se conserva como **componente open-code** (consumidores que verifiquen soporte Invoker Commands lo usan, como ya documenta `dialog.md:7`). Pero **ninguna página Gelium deja un control inerte**: el preview dogfoodea la variante página; la variante modal queda como ejemplo secundario claramente etiquetado o solo en docs.
3. Se conserva el contrato §1.9 (Escape/cancel/closedby) donde el modal existe; la variante página hereda "restore" vía Cancel-link (hay que ajustar la redacción de §1.8).

### 4.2 SELECT MENU

Opciones del enunciado:

- **(a) `<select>` nativo real como control del form.**
- **(b) Página de selección server-driven.**

**Recomendación: (a)** — y es doblemente obligatoria porque **es la promesa que ya está escrita en el código** (`select-menu.css:2-8`, `select.go:22-24`). Razones:

- Es **el propio componente Select de Gelium** (`select.html:1-71`): semántica nativa, teclado nativo, AT, submit de form, `aria-invalid`, forced-colors — todo ya resuelto y testeado.
- Elimina de raíz el trigger inerte: en todo navegador el campo es un `<select>` real.
- **Advertencia de diseño**: un `<select>` nativo y un `<dialog>`-menú **no pueden ser el mismo campo sin JS** (duplicarían controles/submits). Por eso el menú M3 debe (i) **eliminarse** del demo, o (ii) quedar como **variante secundaria etiquetada** ("requiere Invoker Commands; el `<select>` de abajo es el control base"). Recomendado: (i) eliminar el dialog del demo base y conservar el CSS del menú solo si se mantiene una variante de ejemplo; el `<select>` nativo cubre la función del form sin aportar el menú nada a nivel a11y.

- Efectos colaterales positivos: resuelve parcialmente **G8** (`aria-expanded` en triggers de select) y **G10** (`aria-selected` inválido en botones del menú) al eliminar el trigger/menú del camino base.

---

## 5. PROPUESTA DE IMPLEMENTACIÓN (sin tocar código — alcance estimado)

### 5.1 DIALOG — variante página/confirmación

| Archivo | Cambio |
|---|---|
| `internal/app/dialog.go` | Trigger deja `Command`/`CommandFor` y pasa a `Href: "/components/dialog/confirm"` (renders `<a>` en `button.html:6`). Nuevo handler `GET /components/dialog/confirm` que renderiza la confirmación inline (mismo headline/descripción) vía partial o slot. Confirm = POST real a nueva ruta `POST /examples/dialog/confirm` (→303 a `/components/dialog`) o, mínimo, Confirm como link. Cancel = link back. |
| `web/templates/dialog.html` | Nuevo partial `dialog-confirm` (página inline). El partial `dialog` (modal) se conserva solo como componente de mejora. |
| `internal/app/server.go` | Registrar `GET /components/dialog/confirm` y `POST /examples/dialog/confirm` + entrada en `postOnlyPaths()`. |
| `internal/app/dialog_test.go` | Contratos cambian: trigger pasa a `<a href="/components/dialog/confirm">`; se quitan/reemplazan asserts de `command=` (`dialog_test.go:20-22`); nuevos tests de la ruta confirm (GET 200, POST 303/422). Conservar asserts `closedby`/Escape para la variante modal si se mantiene. |
| `web/content/dialog.md`, `README.md:120` | Corregir "Baseline Low" → **Baseline 2025 (Newly available)**; describir la variante página como base y el modal como mejora opt-in. |
| `docs/gelium-ui-accessibility-contract.md` | Actualizar §1.8/§1.9 y fila G1: fallback página entregado; nota de `closedby` (Safari) y Firefox. |

### 5.2 SELECT MENU — `<select>` nativo real

| Archivo | Cambio |
|---|---|
| `web/templates/select.html` (`select-menu-demo`, :73-92) | Reemplazar `.m3-select-trigger` button + hidden `value` por `<select id="select-menu" name="value" aria-label="…">` con las 3 opciones, con estilos `.ui-select` del propio componente. Decidir: eliminar el `<dialog>` o dejarlo como variante etiquetada. Mantener form `POST` + `hx-post` fragment flow. |
| `internal/app/select.go` | `selectMenuDemo` pasa a portar `Options` a un `<select>` nativo; eliminar render `aria-selected`/`Open` (`:12-14,111`); flujo de validación 422 con vocabulario cerrado intacto. |
| `web/styles/select-menu.css` | Si se elimina el menú: quitar CSS y corregir el comentario falso (`:2-8`). Si se mantiene variante: solo variante etiquetada. |
| Tests | `select_test.go` (quitar asserts `command=`/`aria-selected` `:52-67,86-88` → asserts de `<select>`), `web/styles_select_menu_test.go`, `web/styles_select_test.go`. `web/styles_button_test.go` conserva asserts de `command=` (siguen usados por drawer/fab/icon-button) — sin cambio. |
| Docs/contrato | Fila G1: marcar sub-gap select resuelto; G8/G10 parcialmente resueltos. |

### 5.3 Alcance estimado

- **~9-11 archivos** (5-6 por overlay, sin overlaps) + docs de contrato.
- El grueso del esfuerzo es **tests + prosa de docs**, no markup: los asserts de contrato actuales (`dialog_test.go`, `select_test.go`, `styles_select_menu_test.go`) dependen de `command=`/`closedby`/`aria-selected`.
- Slice de 1 commit enfocado dentro de Phase E (el roadmap ya tiene la fila `Fallback overlays no-Chromium (Dialog/Select) | component | E`).

---

## 6. RIESGOS Y DECISIÓN RECOMENDADA

### Riesgos

1. **Pérdida de dogfood del modal**: los docs dejarían de mostrar el `<dialog>` modal en vivo. Mitigación: conservar el partial modal como componente + ejemplo etiquetado; el modal ya es Baseline 2025, así que es un ejemplo legítimo — pero **nunca el único demo**.
2. **Trampa del doble control en select**: mantener `<select>` Y `<dialog>` como el mismo campo duplica controles/submits. Mitigación: uno solo como base; el otro es variante etiquetada o se elimina.
3. **Churn de tests**: los contratos actuales asertan `command=`/`closedby`/`aria-selected`; hay que reescribirlos (riesgo bajo, mecánico, pero es donde vive el 60% del trabajo).
4. **Drift factual en docs**: "Baseline Low 2024" está obsoleto (Baseline 2025). Corregirlo en el mismo slice evita re-drift.
5. **Scope creep con Navigation drawer**: comparte el arquetipo. Recomendado: documentarlo con el mismo contrato y resolverlo en el mismo PR o en el siguiente, no olvidarlo.
6. **§1.8 focus restoration**: la variante página no tiene overlay; "restauración" = Cancel-link. Ajustar redacción del contrato para no prometer trap donde no hay overlay.

### Decisión

**Fix AHORA en Phase E**, slice enfocado. Matices:

- La **urgencia real del botón muerto bajó**: `command`/`commandfor` es Baseline desde dic-2025; el dead-button solo afecta versiones no actuales. El sub-gap que NO puede esperar es el **Select menu** (form entero muerto + promesa falsa docs-vs-código).
- Orden recomendado dentro del slice: **Select menu primero** (`<select>` nativo — deja de haber control muerto en el sistema), **Dialog después** (variante página + modal como mejora), **docs/Baseline corregidos en el mismo commit**.
- Diferir (fuera de Phase E) es defendible solo si la política de soporte de Gelium declara "solo navegadores actuales con Invoker Commands", pero eso choca con el pilar no-JS/no-imitación y con la promesa ya escrita en `select-menu.css`/`select.go`. Fijar ahora.

---

## 7. EVIDENCIA CLAVE (file:line)

- `web/templates/dialog.html:2-9` — trigger `button` + `<dialog closedby="any">`.
- `internal/app/dialog.go:17-19` — `Command`/`CommandFor` (`show-modal`, `request-close`, `close`); Confirm sin acción server.
- `web/templates/select.html:73-92` — form + trigger `command` + dialog con items submit + hidden inputs; **sin `<select>`**.
- `web/styles/select-menu.css:2-8` y `internal/app/select.go:22-24` — promesa falsa de fallback `<select>` nativo.
- `internal/app/select.go:12-14,111` — `Open` solo marca `aria-selected` en dialog siempre-cerrado.
- `web/templates/button.html:9` / `internal/app/button.go:13-14` — render de `command`/`commandfor` en el trigger.
- `web/templates/navigation-drawer.html:6-12` / `navigation_drawer.go:121-126` — tercer overlay con el mismo patrón.
- `web/content/dialog.md:7,9` / `README.md:120` / contrato §1.8:165, §1.9:182, §3:347 — estado documentado del gap.
- `internal/app/dialog_test.go:17-24,47-70` / `select_test.go:52-67,86-88` — asserts actuales que cambiarían.
- Detección ausente: `rg "@supports"` solo `anchor-name` (`menu.css:62`, `select-menu.css:26`); `app.js:1-9` solo swap 422.
- Compatibilidad: MDN `HTMLButtonElement.command` (Baseline 2025 Newly available, dic-2025); caniuse `html element: dialog: closedby` (71.54%; Safari 26.4–27 ❌, Firefox 141+ ✅, Chrome/Edge 134+ ✅); MDN `<dialog>` (Baseline Widely available, mar-2022).
