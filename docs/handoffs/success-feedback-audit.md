# Gelium UI — Success Feedback Persistente Audit (Phase D, handoff)

> **Alcance**: inventario read-only del patrón SUCCESS FEEDBACK PERSISTENTE (patrón 8/último de Phase D, después de Empty State `eba1c4c`, Skeleton `0688020`, Inline Alert `43c0dac`, Validation Summary `c9afd49`, Banner `5787799`, Callout `664a40c`, Error State `157dc7f`). No modifica código, templates, CSS, tests ni docs. Única escritura: este handoff.
>
> **Baseline**: `docs/gelium-ui-system-roadmap.md` (Phase D :147-172, matriz :415), `docs/gelium-ui-vocabulary.md` (§Toast, §Inline alert, §Banner), `docs/gelium-ui-composition-rules.md` (§4.8 :110-119, anti-rule 4 :126, §9 :164-173), `docs/handoffs/{state-patterns-audit,inline-alert-audit,banner-audit}.md`, `internal/app/{text_field,toast,data_table,demo_whatsapp,server,routes}.go`, `web/templates/{inline-alert,banner,text-field,data-table,chips,demo-whatsapp,toast}.html`, `web/styles/{inline-alert,banner,tokens}.css`, `web/styles_contract_test.go`, `web/styles_{inline_alert,banner}_test.go`, `internal/app/{server,toast,text_field}_test.go`.

---

## 1. ESTADO ACTUAL — cómo se confirma éxito hoy

**No existe success persistente (nivel sección/página) en ningún template ni handler.** El único `role="status"` persistente ligado a éxito es el helper de campo, no un patrón de feedback:

| Vía | Evidencia | Tipo |
|---|---|---|
| Helper de campo success | `internal/app/text_field.go:70-71` (`field.Helper = "Name accepted"`, `field.MessageRole = "status"`) → render en `web/templates/text-field.html:8` (`role="{{.MessageRole}}"` → `role="status"`). Campo-level, NO sección/página | Persistente, campo |
| Toast success transitorio | `internal/app/toast.go:82` (demo estático "Your changes were saved."), `internal/app/toast.go:99`; `internal/app/data_table.go:399,417` ("Data refreshed." — `HX-Trigger loom:toast` + fallback inline no-JS `data-table.html:81`) | Transitorio |
| WhatsApp demo (POST+303) | `internal/app/demo_whatsapp.go:559,573` → `http.Redirect(..., http.StatusSeeOther)`; la página destino re-renderiza **sin ningún mensaje de éxito** — el éxito se infiere de la mutación visible (mensaje aparece en la lista) | Sin confirmación |
| Notices ad-hoc `role="status"` | `web/templates/data-table.html:21` (`.data-table-demo-notice`, selección) y `web/templates/chips.html:63` (`.chips-demo-notice`, remoción) — status de demo, no reusables, no success | Demo ad-hoc |
| Error de validación | `text_field.go:64-68` (422 + `X-Loom-Validation`), `text_field.html:5,8` (`role="alert"`); NUNCA toast (`toast.go:129-133`) | Persistente, error |

**Conclusión**: hoy una operación exitosa se confirma con **toast success transitorio** (data-table refresh) o **no se confirma** (WhatsApp send → 303 sin mensaje). No hay banner/inline success persistente consumiendo los tones que ya existen. Coincide con `state-patterns-audit.md:21` ("✖ persistente no existe; solo toast transitorio + role=status helper de campo") y con `composition-audit.md:227`.

---

## 2. DECISIÓN — REUSO (no componente nuevo)

**Recommendación: SUCCESS FEEDBACK es un patrón DOCUMENTAL (guía de uso + tests de confirmación). NO se crea archivo de componente nuevo.**

Evidencia verificada:

- **Inline alert ya sirve como success persistente de sección**: `web/templates/inline-alert.html:1` deriva el rol del tone — `role="alert"` solo si `eq .Tone "error"`, **`role="status"` para el resto incluido success**. Tone `--success` implementado en `web/styles/inline-alert.css:31` (`color-mix` sobre `--ui-color-success`, precedente `.demo-wa-quality--GREEN`). Contrato CSS pinneado en `web/styles_inline_alert_test.go:22,61` y `TestInlineAlertTonesUseCoreTokens`.
- **Banner ya sirve como success persistente de página**: `web/templates/banner.html:1` misma derivación `error→alert`, resto→`status` (success incluido). Tone `--success` en `web/styles/banner.css:31`; pinneado en `web/styles_banner_test.go:23,66,102-105,124`. **El render Go ya lo verifica**: `internal/app/server_test.go:296-307` `TestBannerRoleIsDerivedFromTone` (success → `role="status"`).
- **Ambos partials ya derivan `role="status"` para tone success**: confirmado en las dos templates (líneas 1 de cada una). El enunciado del roadmap "reusa Banner o Inline alert con tone success" (`roadmap.md:42` del state-patterns-audit) **ya está entregado** por los commits `43c0dac` y `5787799`.
- **Tono success existe desde B4**: `--ui-color-success` (:45) y `--ui-color-success-fg` (:46) en `web/styles/tokens.css`. `success-container` no existe deliberadamente (`tokens.css:40-43`) → success usa `color-mix`, cero tokens nuevos.

**Único hueco real**: **Inline alert no tiene test de render que pruebe `role="status"` para success** (no existe view model Go ni helper `renderInlineAlert`; solo tests CSS que verifican que el selector de tone existe, no el rol derivado). Banner sí lo tiene.

Por lo tanto, el entregable es: **documentación de contrato de uso + tests de confirmación de los tones success existentes**. Sin `success-feedback.html`, sin `success-feedback.css`, sin handler Go nuevo.

---

## 3. CONTRATO DE USO (guía documental)

| Señal | Cuándo | Rol ARIA | Evidencia de base |
|---|---|---|---|
| **`inline-alert--success`** | Éxito persistente de **sección/form**: "Settings saved" inline dentro de la región del form que se guardó; "Contact added" bajo el form de alta. Sobrevive a la navegación si la sección persiste | `role="status"` (polite) | `inline-alert.html:1`, `inline-alert.css:31` |
| **`banner--success`** | Éxito persistente de **página/operación global**: POST+303 aterriza en la página destino con "Settings saved" al tope (slot `layout.html` entre `</header>` y `<main>`); sin auto-dismiss | `role="status"` (polite) | `banner.html:1`, `banner.css:31`, `server_test.go:296-307` |
| **`toast` success** | SOLO transitorio post-acción ("Data refreshed.", `data_table.go:399,417`). **NUNCA** para success persistente | `role="status"`/`alert` en `#loom-toast-region` | `toast.html`, anti-regla `composition-rules.md:126` (anti-rule 4) |

**Reglas**:
- La regla a codificar en docs: *nada persistente se anuncia con `loom:toast`; nada transitorio ocupa un slot persistente* (`roadmap.md:162-168`, `composition-rules.md:119`).
- El partial de success persistente es **output del servidor** (nunca estado cliente); si el estado muere en el render siguiente, el servidor decide (POST+303 lo limpia o lo deja según el flag).
- Dismiss de un banner--success NO es un caso normal (el éxito no pide acción); si se necesita, contrato existente POST+303 (`banner.html:8-10`, `demo_whatsapp.go:559`).
- Nada color-only: icono opcional `aria-hidden` + texto + role + bloque forced-colors ya cubiertos por los partials existentes.
- **Tokens nuevos: NO.** success/info usan `color-mix` sobre `--ui-color-success`/`--ui-color-info` (decisión ya tomada en `inline-alert-audit.md:112` y `banner-audit.md:104`).

---

## 4. CONTRATO SERVER — POST + 303 → success persistente (patrón WhatsApp)

Reusa el contrato existente (d) de `gelium-ui-core.md` (`state-patterns-audit.md:67,74`):

```text
POST {accion} → handler muta el store → http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
              → la página destino re-renderiza el success persistente (role="status")
```

- **No-JS (flujo principal)**: el redirect 303 de `demo_whatsapp.go:559,573` es el patrón exacto; la página destino incluye `{{template "banner" ...}}` con `Tone: "success"` (slot `layout.html` + campo `pageView.Banner`, propuesta `banner-audit.md:139-148`) o el `inline-alert--success` dentro de la sección. 0 JS, `app.js` intacto.
- **HTMX (enhancement)**: el fragmento post-submit re-renderiza el success persistente en la región (mismo mecanismo que `data-table.html:81` incluye el toast inline, pero con `role="status"` persistente en vez de auto-dismiss). **No** se emite `HX-Trigger loom:toast` para este caso (regla transversal `state-patterns-audit.md:78`).
- **Fallo**: validación → 422 + `X-Loom-Validation` + `inline-alert--error`/`validation-summary`; error global → `banner--error`. El success nunca comparte vía con el error (el handler renderiza uno u otro).
- **Ejemplo Phase G — Settings**: `POST /settings` → valida → 303 a `/settings?tab=general` → página con `banner--success` "Settings saved" (`role="status"`, persistente) + los valores guardados re-renderizados. El toast transitorio queda descartado para esta confirmación.

---

## 5. TESTS PROPUESTOS

**Confirmación de reuso (el entregable de tests)**:

1. **Inline alert: derivación de rol por tone** (GAP detectado — hoy solo existe a nivel CSS, no a nivel render): agregar helper `renderInlineAlert` + test análogo a `TestBannerRoleIsDerivedFromTone` (`internal/app/server_test.go:296-307`) que pruebe `tone=success → role="status"` y `tone=error → role="alert"` contra `web/templates/inline-alert.html`. (Banner YA cubierto por `server_test.go:296-307`.)
2. **Test de guía "persistente nunca toast"**: test de contrato que verifique que los partials de success persistente (`inline-alert.html`, `banner.html`) NO emiten `HX-Trigger loom:toast` y que la clase `--success` en ambos resuelve a `role="status"` en el template (guard del anti-rule `composition-rules.md:126`). Puede vivir en `web/styles_contract_test.go` (vocabulario cerrado por patrón) o como render test Go.
3. **Los tests CSS existentes YA pasan como confirmación** (no requieren cambios): `TestInlineAlertTonesUseCoreTokens` (`web/styles_inline_alert_test.go`), `TestBannerTonesUseCoreTokens` (`web/styles_banner_test.go:116-136`) y `TestBannerContractCSSWired` — prueban que `--success` está presente y usa solo tokens core sin hex literals.
4. **(Opcional, integración Phase G)** `internal/app/settings_test.go` — POST /settings → 303 → body destino contiene `ui-banner--success` + `role="status"` "Settings saved"; POST inválido → 422 + `inline-alert--error`, sin toast.

Si la decisión hubiera sido componente nuevo, el set sería `web/styles_success_feedback_test.go` + `web/styles/success-feedback.css` — **no aplica** (reuso).

---

## 6. FILES IMPACTADOS (solo read-only)

**Documentación (el entregable real)**:
- `docs/gelium-ui-vocabulary.md` — **no existe sección "Success feedback"** (solo el tone en Toast :144-148); agregar entrada `### Success feedback` con guía de uso (inline-alert--success vs banner--success vs NUNCA toast) y mapeo a los partials existentes. (Mismo patrón de las entradas Inline alert :112-120 / Banner :122-129.)
- `docs/gelium-ui-composition-rules.md` — §4.8 (:110-119) y anti-rule 4 (:126) ya cubren la regla; opcional reforzar con ejemplo success persistente.
- `docs/gelium-ui-system-roadmap.md` — matriz :415 ("Success feedback persistente | … | Crear contrato (no confundir con toast)") marcarlo como resuelto por reuso de tones success.

**Tests**:
- `internal/app/server_test.go` — agregar render test de rol para inline-alert (gap §5.1).
- `web/styles_contract_test.go` — opcional: guard de guía "persistente nunca toast" (§5.2).

**Integración Phase G (sin implementar aquí)**:
- `web/templates/layout.html` — slot `{{if .Banner}}{{template "banner" .Banner}}{{end}}` (propuesta `banner-audit.md:139-147`).
- `internal/app/server.go` — campo `Banner *bannerView` en `pageView` (:32-69).
- `internal/app/settings.go` (futuro) — handler POST+303 con `bannerView{Tone:"success"}`.

**NO tocados**: `web/templates/inline-alert.html`, `web/templates/banner.html`, `web/styles/{inline-alert,banner,tokens}.css`, `web/static/app.js`, `web/static/app.css`, `internal/app/{toast,text_field,data_table,demo_whatsapp}.go` — el reuso de los tones success existentes no requiere cambios. **Ningún archivo nuevo de componente.**

---

## 7. Fuentes de autoridad

`docs/gelium-ui-system-roadmap.md` (Phase D :147-172, matriz :415), `docs/gelium-ui-vocabulary.md` (:110-149), `docs/gelium-ui-composition-rules.md` (:110-119, :126, :164-173), `docs/handoffs/{state-patterns-audit,inline-alert-audit,banner-audit}.md`, commits `43c0dac` (Inline Alert) y `5787799` (Banner), `web/templates/{inline-alert,banner,text-field,data-table,chips,toast,demo-whatsapp}.html`, `web/styles/{inline-alert,banner,tokens}.css`, `web/styles_{inline_alert,banner}_test.go`, `internal/app/{server,toast,text_field,data_table,demo_whatsapp}.go`, `internal/app/{server,toast,text_field}_test.go`.
