# Phase B Review — Gelium UI (carpeta física `loom-ui`)

> Handoff read-only de coherencia. No modifica código, templates, CSS ni tests.
> Alcance: migración a tokens de la Phase B (commit `d86c03a`): `web/styles/tokens.css` (core neutral) vs `themes/theme-material/theme.css` (3 rutas: light, dark clase, dark media).
> Sub-slices: B1 spacing, B2 density/sizes, B3 borders, B4 semantic colors, B5 motion.
> Fecha: 2026-08-10. Fuentes: `web/styles/tokens.css`, `themes/theme-material/theme.css`, `web/styles/toast.css`, `web/styles/app.css`, `web/static/app.css` (build), `web/styles_contract_test.go`, `web/styles_toast_test.go`, diff `d86c03a~1..d86c03a`, `docs/gelium-ui-theme-contract.md`, y barrido completo de `var(--ui-*)` en los 27 archivos de `web/styles/`.

---

## 1. Resumen ejecutivo

La migración es **sólida**: el core es realmente agnóstico (defaults neutros, sin valores Material), los consumidores están tokenizados (cero hex literales en la capa de componentes), los gaps de tokens inexistentes quedaron cerrados, y la arquitectura 3-rutas es coherente salvo un drift puntual.

Se detectan **1 problema grave**, **1 drift medio**, **3 tokens muertos**, **4 migraciones parciales cosméticas** y **2 pulidos de estilo**:

1. **GRAVE — Regresión de contraste en los iconos del toast (las 2 rutas)**. La migración re-apuntó `--ui-toast-icon-*` de valores explícitos correctos por ruta a `var(--ui-color-*)`. El toast usa **superficie inversa** (light → contenedor oscuro `#322f35`; dark → contenedor claro `#ece6f0`), pero `--ui-color-*` resuelve a la paleta de superficie (polaridad opuesta). Resultado: en dark los iconos pasaron de 4.19–5.37:1 a **1.14–1.60:1** (casi invisibles); en light pasaron de 6.73–9.40:1 a **1.98–2.57:1**. El test `TestToastIconTokensDeriveFromCore` (`web/styles_contract_test.go:634-644`) **codifica la regresión** en vez de detectarla.
2. **MEDIO — Drift dark clase vs dark media**: `--ui-switch-track-unselected: #4a4458` está en el bloque de clase dark pero **falta** en el bloque de media query (`theme.css:264` vs `292-350`). Con OS dark sin clase explícita, el track no seleccionado del switch queda `#e6e0e9` (lavanda claro). El propio `docs/gelium-ui-theme-contract.md` ya documenta esta duplicación como drift conocido a eliminar.
3. **BAJO — Tokens muertos en core**: `--ui-space-8`, `--ui-state-selected-opacity` (sin consumidor) y `--ui-color-success-fg` (override en 3 rutas sin consumidor). El commit removió `state-dragged-opacity` como muerto pero dejó `state-selected-opacity`.
4. **BAJO — Migraciones parciales** (literales donde existe token): `toast.css:20` (`inset-inline: 1rem`), `select.css:24` (`padding: 1rem … 1rem`), `text-field.css:22,58,63` (`inset-inline-*: 1rem`), `chips.css:134` (`.5rem .625rem` en `.ui-chip-remove`). Solo cosmético; los bloques de demo conservan literales con criterio.
5. **COSMÉTICO**: indentación rota en `theme.css:264-265`; `--ui-color-warning-fg: #fdd663` en dark es idéntico a `--ui-color-warning` (redundante, no rompe el demo WhatsApp que sobre contenedor oscuro da 9.25:1).

**Veredicto**: la Phase B cumple el objetivo de foundation theme-neutral; el único problema que bloquea es el toast (a11y real en el artefacto build). El resto son deuda menor.

---

## 2. TOAST_DARK_RISK — análisis detallado

### 2.1 Por qué pasa

- El toast de Material usa **inverse surface**: en light el contenedor es oscuro (`#322f35`), en dark es claro (`#ece6f0`). El texto (`--ui-toast-fg`) y la acción (`--ui-toast-action`) usan valores por ruta y **no** se vieron afectados.
- La migración re-apuntó los iconos a la paleta semántica de **superficie**:
  `--ui-toast-icon-*: var(--ui-color-info|success|warning|danger)`.
- En dark, esos tokens resuelven a la paleta dark (colores claros tipo `#d0bcff`/`#81c995`/`#fdd663`/`#f2b8b5`), pensados para superficies oscuras — pero el toast dark es **claro**. Polaridad invertida.
- Pre-migración (`d86c03a~1`) los valores eran explícitos por ruta y correctos; la migración los rompió en las dos rutas.

### 2.2 Contrastes medidos (WCAG — iconos son "gráficos" no-texto, umbral 3:1; texto normal 4.5:1)

| Ruta | Contenedor | Icono | Antes | Después | Veredicto |
|---|---|---|---|---|---|
| dark | `#ece6f0` | info `#d0bcff` | `#6750a4` **5.26** | **1.39** | FAIL (invisible) |
| dark | `#ece6f0` | success `#81c995` | `#2e7d32` **4.19** | **1.60** | FAIL |
| dark | `#ece6f0` | warning `#fdd663` | `#7a5700` **5.37** | **1.14** | FAIL (peor caso) |
| dark | `#ece6f0` | error `#f2b8b5` | `#b3261e` **5.34** | **1.39** | FAIL |
| light | `#322f35` | info `#4f5d75` | `#d0bcff` **7.73** | **1.98** | FAIL |
| light | `#322f35` | success `#2e7d32` | `#81c995` **6.73** | **2.57** | FAIL |
| light | `#322f35` | warning `#b9930a` | `#fdd663` **9.40** | **4.54** | pasa 3:1, no 4.5:1 |
| light | `#322f35` | error `#b3261e` | `#f2b8b5` **7.72** | **2.02** | FAIL |

Texto (`--ui-toast-fg`) y acción (`--ui-toast-action`) intactos: dark 13.94:1 y 5.26:1; light 11.46:1 y 7.73:1.

### 2.3 Fix mínimo propuesto

Los iconos deben seguir la paleta **inversa** (la que ya tenía el toast pre-migración), no la de superficie. Opciones:

- **A (recomendada, mínima y reversible)**: restaurar valores explícitos por ruta en `theme.css`:
  - light: `--ui-toast-icon-info: #d0bcff; success: #81c995; warning: #fdd663; error: #f2b8b5;`
  - dark clase **y** dark media (las dos rutas): `info: #6750a4; success: #2e7d32; warning: #7a5700; error: #b3261e;`
  - Resultados: dark 4.19–5.37:1 (todos ≥ 4.5 salvo success 4.19 que pasa no-texto 3:1), light 6.73–9.40:1.
- **B (más limpia, más trabajo)**: introducir tokens de superficie inversa (p. ej. `--ui-inverse-*`) que mapeen la polaridad opuesta, y que el toast los consuma. Correcto para multi-theme (Basecoat) pero excede la Phase B.
- **C (no viable tal cual)**: usar `*-fg`. En dark `success-fg` `#00391f` da 10.68:1 (bien) pero **no existe `info-fg`** y `warning-fg` dark es idéntico a `warning` (`#fdd663`, 1.14:1) — no resuelve el peor caso.

**Obligatorio además**: revisar `TestToastIconTokensDeriveFromCore` (`web/styles_contract_test.go:634-644`), que hoy exige `var(--ui-color-*)` y bloquea la corrección. El refactor de tests a "contrato en vez de hex" de la Phase B ocultó esta regresión; el contrato debería exigir contraste mínimo por ruta o no derivar iconos de inversa desde la paleta de superficie.

**Verificado en el build**: `web/static/app.css` contiene `--ui-toast-container:#ece6f0` con `--ui-toast-icon-info:var(--ui-color-info)` en el bloque dark y la regla `.ui-toast-icon-info{color:var(--ui-toast-icon-info)}`. El riesgo está en producción.

---

## 3. GAPS ENCONTRADOS

| Tipo | Ubicación | Gravedad | Fix sugerido |
|---|---|---|---|
| Regresión a11y toast | `themes/theme-material/theme.css:53-56` (light) y `:278-281` / `:338-341` (dark clase + media) | **ALTA** | Valores explícitos por ruta (sección 2.3, opción A) o tokens de superficie inversa; corregir el test de contrato que la codifica |
| Drift dark clase vs media | `theme.css:264` tiene `--ui-switch-track-unselected:#4a4458`; el bloque media `292-350` NO lo tiene | MEDIA | Agregar la línea al bloque media (o ejecutar el refactor "una sola rutina dark" ya acordado en `docs/gelium-ui-theme-contract.md`) |
| Token muerto | `web/styles/tokens.css:106` `--ui-space-8` | BAJA | Eliminar del core (sin consumidor) |
| Token muerto | `tokens.css:59` `--ui-state-selected-opacity` | BAJA | Eliminar (el commit ya removió `state-dragged-opacity`; este quedó) |
| Token muerto | `tokens.css:46` y `theme.css:27,251,311` `--ui-color-success-fg` (3 overrides, 0 consumidores) | BAJA | Eliminar, o darle consumidor (candidato para icono success del toast dark) |
| Migración parcial | `toast.css:20` `inset-inline: 1rem` | BAJA | `var(--ui-space-4)` |
| Migración parcial | `select.css:24` `padding: 1rem var(--ui-select-caret-reserve) 0 1rem` | BAJA | `padding: var(--ui-space-4) var(--ui-select-caret-reserve) 0 var(--ui-space-4)` |
| Migración parcial | `text-field.css:22,63` `inset-inline-*: 1rem`; `:58` `padding-inline-end: 3.5rem` | BAJA | `--ui-space-4` para los `1rem` (3.5rem es layout de reserva, puede quedar) |
| Migración parcial | `chips.css:134` `margin-inline: .5rem .625rem` | BAJA | `.5rem` → `--ui-space-2` (.625rem sin token, dejar literal) |
| Cosmético | `theme.css:264-265` indentación rota | BAJA | Reindentar |
| Redundancia | `theme.css:253` `--ui-color-warning-fg: #fdd663` = `--ui-color-warning` (dark) | BAJA | Documentar o apuntar a un tono distinto si algún día se usa sobre warning puro |
| Remanente identidad | `.m3-select-trigger`/`.m3-select` en `select-menu.css:23,75`, `app.css:64,175,181`, `select.html:74,77` | BAJA (pre-existente, no de Phase B) | Renombrar a `ui-select-menu-*` (ya señalado en `core-audit.md`) |

Nota: `--ui-slider-fill` (única referencia "sin definición") NO es gap: es custom property set por runtime en `app.js` con fallback `0%` (`slider.css:32`).

---

## 4. CONSISTENCIA — confirmaciones core ↔ theme

- **Color**: las 26 declaraciones core de color (canvas, surface, surface-container, fg, fg-muted, primary, primary-fg, secondary, secondary-fg, border, border-strong, danger, danger-fg, focus-ring + success, success-fg, warning, warning-fg, warning-container, info, danger-container, scrim) están sobrescritas en las **3 rutas**. `--ui-color-error` correctamente NO se sobrescribe (alias `var(--ui-color-danger)` que sigue a cada ruta).
- **State / focus / motion / size / radius / elevation / fonts / typescale**: override solo en el bloque light base; las rutas dark **heredan** (no hay valores que dependan del modo). Consistente.
- **Spacing / border**: el theme NO sobrescribe `--ui-space-*` ni `--ui-border-*` (usa los defaults neutros del core, que ya son Material: grid 0.25rem, 1px/2px solid). Intencional y consistente entre rutas.
- **Z-cero hex literales** en la capa de componentes: todos los colores pasan por tokens.
- **Sin selectores dark en componentes**: `prefers-color-scheme`/`.theme-dark`/`data-theme` solo viven en el theme. Los state layers migraron a `color-mix(in srgb, currentColor <opacity>%, transparent)` — correctos para variantes filled/tonal y multi-theme.
- **B5 motion / reduced motion**: los 150ms/500ms quedaron en `--ui-motion-short/long` y `--ui-easing-standard`; no quedó ningún literal de duración. Cobertura `prefers-reduced-motion: reduce` completa: 13 componentes con bloque propio + bloque central en `app.css:53-71` para button/text-field/dialog/toast/elevation/switch/select/select-menu/slider/progress/fab/list. Cada componente con `transition` está cubierto por uno de los dos lados.
- **Forced colors**: uso consistente de `CanvasText`/`Canvas`/`Highlight`/`ButtonText`/`GrayText`/`LinkText`/`WindowText`; toast (`app.css:118`) y dialog (`app.css:117`) ganaron borde forzado en la migración.
- **Build**: `web/static/app.css` regenerado contiene todas las familias core, el bloque dark media, los tokens del toast y las reglas de consumo (verificado por grep).
- **Orden de import**: `app.css` importa `tokens.css` antes del theme — el contrato de override se respeta.

---

## 5. RECOMENDACIONES — prioridad

1. **P0** — Corregir contraste de iconos del toast en las 2 rutas (sección 2.3 opción A) + ajustar `TestToastIconTokensDeriveFromCore` para que exija contraste mínimo en vez de la derivación actual. Bloqueante para a11y; hoy el defecto está en producción.
2. **P1** — Cerrar el drift `--ui-switch-track-unselected` en el bloque dark media (`theme.css:292-350`). Es el síntoma concreto de la duplicación clase+media que el theme-contract ya manda eliminar.
3. **P2** — Podar tokens muertos del core (`--ui-space-8`, `--ui-state-selected-opacity`, `--ui-color-success-fg`) o darle consumidor a `success-fg`.
4. **P3** — Migrar los literales señalados en `toast.css:20`, `select.css:24`, `text-field.css:22,63`, `chips.css:134` a tokens de escala existentes.
5. **P3** — Reindentar `theme.css:264-265` y documentar `warning-fg` dark redundante.

---

## Ruta del handoff

`docs/handoffs/phase-b-review.md`
