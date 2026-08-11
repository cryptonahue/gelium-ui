# Gelium UI — Phase D+E Coherence Review (handoff)

> **Alcance**: revisión read-only de coherencia de Phase D (state patterns) + Phase E (branding/metadata/a11y). No se modificó código, templates, CSS, tests ni docs. Única escritura: este handoff.
>
> **Baseline**: commits `94340dc` (Success), `157dc7f` (Error state), `42aef8e` (contratos E), `9504216` (código E). Docs `docs/gelium-ui-accessibility-contract.md`, `docs/gelium-ui-seo-contract.md`, `docs/gelium-ui-geo-contract.md`.

---

## 1. RESUMEN EJECUTIVO

Phase D + E están **coherentes**. Los 8 state patterns cumplen el contrato de 4 puntos (partial + CSS + import en `app.css` + entrada en `sourceAppCSS`) y todos tienen tests. Los roles ARIA de los partials coinciden **1:1** con la tabla documentada en `gelium-ui-accessibility-contract.md:325-334` — sin inconsistencias. Branding limpio en las superficies servidas (`web/templates`, `internal/app`, `web/content`): cero residuos de "Gelidium UI"/"Loom UI"; LoomChat queda solo en demos (app ficticia) como corresponde. Metadata server-driven completa y **ningún template se rompe**: `Meta` es un campo de valor (nunca nil-pointer) y el layout protege cada tag con `{{if}}`. Contratos E2 intactos (422 + `X-Loom-Validation`, `loom:toast`, GET params, POST+303). Verificación: `go test -count=1 ./...`, `go vet ./...` y `npm run build` pasan.

Hallazgos menores: residuo de branding en `README.md` (gap abierto documentado por el propio contrato GEO) y en un comentario de `demo-whatsapp.css`; error pages (404/500) no pasan por `resolveMeta` (sin canonical/robots/description). Ninguno es bloqueante.

---

## 2. TABLA DE COHERENCIA DE PATRONES

| Patrón | Partial | CSS | Import app.css | sourceAppCSS | Tests | Wiring | Role en partial |
|---|---|---|---|---|---|---|---|
| Empty state | `empty-state.html` | `empty-state.css` | :33 | :57 | `styles_empty_state_test.go` | `data-table.html:69` (celda `colspan`, `data_table.go:82`) | `status` |
| Skeleton | `skeleton.html` | `skeleton.css` | :38 | :60 | `styles_skeleton_test.go` | primitiva (sr-only "Loading") | `status` + `sr-only` + `aria-hidden` |
| Inline alert | `inline-alert.html` | `inline-alert.css` | :37 | :59 | `styles_inline_alert_test.go` | sección/form-level | `alert` (error) / `status` (resto) |
| Validation summary | `validation-summary.html` | `validation-summary.css` | :40 | :62 | `styles_validation_summary_test.go` | form-level, anchors `#field-error` | `alert` |
| Banner | `banner.html` | `banner.css` | :33 | :55 | `styles_banner_test.go` | slot layout `layout.html:24` | `alert` (error) / `status` (resto) |
| Callout | `callout.html` | `callout.css` | :34 | :56 | `styles_callout_test.go` | tip box inline | **sin role** (`aside`) |
| Error state | `error-state.html` | `error-state.css` | :36 | :58 | `styles_error_state_test.go` | 404 catch-all + 500 (`server.go:355-392`) | `alert` |
| Success (reuso) | — (reusa `inline-alert`/`banner` tone success) | tones `success` en ambos CSS | n/a | n/a | `TestPersistentSuccessPartialsNeverToast` (`styles_contract_test.go:842`) | POST+303 re-render | `status` |

**Sincronía app.css ↔ sourceAppCSS**: el orden es idéntico (el theme se pre-pende aparte vía `themeCSS`). Sin archivos huérfanos: los 7 CSS de state patterns tienen partial y viceversa; `demo-whatsapp.css` tiene sus 2 templates demo.

**Roles ARIA — verificación**: tabla esperada vs real → **100% consistente** (Empty `status`; Skeleton `status`+`sr-only`; Inline `alert`/`status`; Validation `alert`; Banner `alert`/`status`; Callout sin role; Error `alert`; Success `status`). Sin inconsistencias.

---

## 3. HALLAZGOS

| Tipo | Ubicación | Gravedad | Fix sugerido |
|---|---|---|---|
| Branding residuo (público) | `README.md:1` `# Loom UI` | **Media** | Renombrar a "Gelium UI". Es el gap de brand split que el propio `gelium-ui-geo-contract.md:66,199` documenta como pendiente y quedó abierto (fuera del subconjunto `web/templates|internal/app|web/content` pero es la portada del repo). |
| Branding residuo (comentario) | `web/styles/demo-whatsapp.css:2` ("…composes the Loom primitives") | Baja | Reescribir el comentario: "Gelium UI primitives". Es un asset importado por `app.css`. |
| Metadata gap | `internal/app/server.go` `renderErrorPage` (404/500) no pasa por `resolveMeta` → sin `description`/`canonical`/`robots` | Baja | Pasar el error page por `resolveMeta` (el status 404 ya impide indexación, pero el contrato §15 exige canonical en páginas indexables y coherencia de head). O documentar la excepción en el contrato SEO. |
| A11y matiz (documentable) | `web/templates/skeleton.html` `role="status"` | Baja | Cuando HTMX reemplaza el skeleton, la región viva se elimina sin anunciar "loaded" (solo se anuncia el insert). Cumple el contrato actual; documentar el comportamiento si se quiere el anuncio de fin de carga. |
| Branding residuo (histórico) | `PROMPT-MATERIAL-WEB-INVENTORY.md`, `MATERIAL-WEB-PROGRESS.md`, `AI-COMPONENT-IMPLEMENTER-PROMPT.md`, `COMPONENT-ROADMAP.md` | Baja | Opcional — docs de proceso/roadmap, no superficies servidas. Unificar naming al tocarlas. |
| Branding residuo (temp) | `.tmp/*.html` ("Loom UI") | Despreciable | Limpiar `.tmp/` o ignorar (snapshots de trabajo). |

---

## 4. VERIFICACIÓN

| Comando | Resultado |
|---|---|
| `go test -count=1 ./...` | **ok** `loomui/internal/app` 2.387s; **ok** `loomui/web` 1.409s |
| `go vet ./...` | limpio (sin salida) |
| `npm run build` | **Done in 5s** (tailwindcss v4.3.3 `--minify` + `copy-htmx`) |

**Contratos server verificados tras E2** (todos con tests verdes):
- **422 + `X-Loom-Validation`**: `text_field.go:87-89` + hook `htmx:beforeSwap` en `app.js:1-9`; cubierto por `text_field_test.go`.
- **`loom:toast`**: `toast.go:108-160` (HX-Trigger) + listener `app.js:77-81`; guard `TestPersistentSuccessPartialsNeverToast` asegura que inline-alert/banner NUNCA emiten toast. Cubierto por `toast_test.go`.
- **GET params**: demos `q`/`c` (`demo_whatsapp.go:506-507`), canonical sin query (`TestCanonicalIsCleanWithoutQuery`).
- **POST+303**: `whatsAppWebhookSave`/`send`/`send-template`/`typing`/`read` → `StatusSeeOther`; cubierto por `TestWhatsAppAdminWebhookSaveRedirects`.

**Metadata verificada**:
- Layout (`layout.html:6-14`) emite `title · Gelium UI`, `description`, `canonical`, `robots`, bloque OG (`og:type|title|description|url`) y `JSON-LD` (home) — todos tras `{{if}}`.
- Toda página de docs pasa por `resolveMeta` (`server.go:449-464`) → `Title`, `Description` por ruta, `Canonical = siteBaseURL + clean path`, `Robots = "index, follow"` (docs) / `"noindex, nofollow"` (rutas `/demo/` y `/examples/`).
- Demos: `lang="es"` + `noindex` (`demoMetaES` `demo_whatsapp.go:563-566`, tests `demo_whatsapp_test.go:188-224`).
- Ningún template rompe: fragmentos (`validation-form`, `toast-demo-form`, `toast-region`) no usan el layout; error pages usan `Meta` cero y los guards lo toleran.

---

## 5. RECOMENDACIONES

1. **Cerrar el gap de branding en `README.md`** — es el único residuo de marca que el contrato GEO deja explícitamente pendiente.
2. **Unificar la metadata de error pages** pasando `renderErrorPage` por `resolveMeta` (o declarar la excepción en el contrato SEO §15).
3. **Actualizar el comentario** en `web/styles/demo-whatsapp.css:2`.
4. Mantener sincronizado `sourceAppCSS` con `app.css` y el guard anti-toast de persistent success: hoy está consistente y es lo que protege la coherencia de Phase D.

---

## Ruta del handoff

`docs/handoffs/phase-de-review.md`
