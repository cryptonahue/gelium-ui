# Gelium UI — Roadmap inspirado en Filament

> Documento de handoff para otro agente. Describe el trabajo realizado, los contratos vigentes y la secuencia recomendada para continuar. No autoriza commit, push, publicación ni deploy.

**Fecha:** 2026-09-05  
**Repositorio:** `gelium-ui`  
**Rama:** `main`  
**Stack:** Go + `html/template` + Tailwind CSS 4 + HTMX opcional + HTML nativo  
**Workflow:** Gentle AI, TDD, review independiente por slice

---

## 1. Objetivo del programa

Inspirarse en Filament para que Gelium ofrezca una experiencia coherente para construir aplicaciones server-rendered complejas:

```text
resource list
→ search / filters / pagination
→ selection / actions / confirmation
→ record detail
→ forms / validation
→ dashboard metrics
→ feedback / states
→ recipes documentadas
```

La inspiración es conceptual:

- composición orientada a tareas;
- recipes completas, no solo componentes aislados;
- acciones declarativas y verificables;
- estados completos;
- documentación y extensibilidad.

No se debe copiar:

- Laravel, PHP o Livewire;
- un runtime reactivo obligatorio;
- CRUD mágico basado en modelos;
- branding, markup o assets de Filament;
- una API que oculte HTML, URLs, permisos o contratos HTTP.

Gelium mantiene su identidad:

```text
HTML-first
server-rendered
no-JS por defecto
GET para estado navegable
POST + 303 para mutaciones
422 + X-Gelium-Validation para validación
HTMX solo como mejora progresiva
```

---

## 2. Estado actual — slices completados

Todos los slices siguientes fueron ejecutados con Gentle AI, TDD, review independiente y baseline completo.

### 2.1 Data Table

**Archivos principales:**

- `internal/app/data_table.go`
- `internal/app/data_table_test.go`
- `lib/templates/data-table.html`

**Resultado:**

- Las filas derivan sus celdas de la misma definición de columnas que la cabecera.
- Se eliminó el cuerpo hardcodeado `Name / Status / Date`.
- Test de view model y test de template renderizado cubren regresión.

### 2.2 Search

**Archivos:**

- `internal/app/recipe_admin_resource.go`
- `internal/app/recipe_admin_resource_test.go`

**Resultado:**

- El campo `q` usa `type="search"`.
- Se preserva el GET y el camino no-JS.

### 2.3 Filtros

**Archivos:**

- `internal/app/recipe_admin_resource.go`
- `internal/app/recipe_admin_resource_test.go`
- `site/web/templates/recipe-admin-resource.html`

**Contrato:**

```text
GET ?q=&status=&sort=&dir=&page=&selection=
```

`status` usa vocabulario cerrado:

```text
Active | Pending | Done
```

El filtro:

- renderiza un `<select>` nativo;
- conserva la opción seleccionada;
- filtra server-side;
- se preserva en sort/pagination;
- no agrega un nuevo componente público.

### 2.4 Empty states

**Resultado:**

- búsqueda sola sin resultados → `Clear search`;
- búsqueda + estado sin resultados → `Clear filters`;
- recovery mediante links GET reales;
- la documentación fue alineada con el comportamiento.

### 2.5 Selection y navegación

**Resultado:**

- IDs válidos se deduplican;
- IDs desconocidos se descartan;
- `selection` se preserva en sort y pagination;
- `selection=all` significa todos los registros del conjunto filtrado;
- no se permite que el estado dependa de JavaScript.

**Advertencia:** preservar todos los IDs en la URL puede requerir un token server-side para datasets grandes. La demo actual es in-memory y pequeña.

### 2.6 Bulk delete

**Ruta:**

```text
GET  /recipes/admin-resource/bulk-delete
POST /recipes/admin-resource/bulk-delete
```

**Contrato:**

1. Selección en la lista.
2. GET de confirmación server-rendered.
3. `<dialog open>` con campos hidden `selection`.
4. POST de mutación.
5. Revalidación server-side contra el store.
6. `303 /recipes/admin-resource`.
7. Banner persistente de éxito.

`selection=all` se reevalúa contra el conjunto filtrado en el momento de la mutación.

**Importante:** esto es una recipe demo. No contiene todavía autorización real ni auditoría persistente.

### 2.7 Pagination

**Archivos:**

- `lib/styles/pagination.css`
- `lib/styles/data-table.css`
- `lib/styles_pagination_test.go`
- `lib/styles_data_table_test.go`

**Resultado:**

Los controles standalone e inline usan:

```css
min-width: var(--ui-touch-target);
min-height: var(--ui-touch-target);
```

Se eliminó la dimensión fija de `2rem`.

### 2.8 Record Detail / Infolist recipe-level

**Ruta:**

```text
GET /recipes/admin-resource/{id}
```

**Archivos:**

- `internal/app/recipe_admin_resource.go`
- `internal/app/recipe_admin_resource_test.go`
- `internal/app/server.go`
- `site/web/templates/recipe-admin-resource.html`
- `site/web/styles/recipe-admin-resource.css`

**Composición:**

- `<article>`;
- `<section aria-labelledby>`;
- `<dl>` + `<dt>` + `<dd>`;
- badge textual de estado;
- links hermanos para Edit y Delete;
- 404 real con `error-state`;
- GET-only.

**Decisión:** no crear todavía `ui-infolist`. Hace falta un segundo consumidor real para probar que la anatomía es estable y compartida.

### 2.9 KPI Dashboard

**Ruta:**

```text
GET /recipes/admin-dashboard
```

**Archivos nuevos:**

- `internal/app/recipe_admin_dashboard.go`
- `internal/app/recipe_admin_dashboard_test.go`
- `site/web/templates/recipe-admin-dashboard.html`
- `site/web/styles/recipe-admin-dashboard.css`

**Integración:**

- `internal/app/server.go` registra la ruta;
- `site/web/styles/app.css` importa el CSS recipe-level;
- `applyRequestChrome` cubre theme/scheme.

**Composición:**

- `<section>` con heading;
- grid fluid `auto-fit/minmax`;
- `<article class="ui-card ui-card-outlined">`;
- métricas reales desde `resourceDemoStore`:
  - total projects;
  - active projects;
  - pending projects;
  - done projects;
- recent projects con `ui-list`, máximo cinco, ordenados por fecha.

**Estados de fixture:**

```text
?state=ready
?state=empty
?state=loading
?state=error
```

- value: cards KPI;
- empty: mensaje + link a resources;
- loading: `role="status"` + `aria-busy="true"`;
- error: `role="alert"` + retry link.

No hay chart, polling, polling interval, WebSocket ni runtime reactivo.

---

## 3. Contratos que otro agente debe respetar

### 3.1 No-JS

El flujo principal debe funcionar sin scripts:

- lectura = GET;
- filtros = GET query params;
- pagination = links GET;
- detalle = GET;
- confirmación = página server-rendered con `<dialog open>`;
- mutación = form POST + 303;
- validación = 422 + `X-Gelium-Validation`;
- feedback persistente = Banner/Inline alert;
- Toast solo para resultado transitorio.

### 3.2 Semántica

Preferir:

- `<article>` para una entidad autocontenida;
- `<section aria-labelledby>` para grupos;
- `<dl>` para label/value;
- `<ul>`/`<ol>` para listas;
- `<a>` para navegación;
- `<button>` o `<form>` para mutaciones;
- `<dialog open>` únicamente como variante de confirmación server-rendered.

No usar:

- `div` como botón;
- links GET para mutaciones;
- color como único estado;
- nested interactive controls;
- glyphs Unicode inventados como iconos;
- métricas ficticias;
- dimensiones o identidades inventadas.

### 3.3 Theme y tokens

- Theme por clase en `<html>`.
- `applyRequestChrome` debe contemplar cada nuevo recipe view.
- CSS usa `--ui-*`.
- No introducir colores, spacing, radii, shadows, widths o letter-spacing literales sin contrato.
- Touch targets mínimo `var(--ui-touch-target)`.
- Responsive fluido antes que breakpoints.

### 3.4 Estados

Cada recipe debe declarar:

```text
value/rest
empty
loading
error
success cuando haya mutación
```

Pero no se debe simular loading inicial en una página server-rendered que no tiene fase asíncrona real. Loading solo corresponde a una región realmente en espera o a un fixture explícito/documentado.

---

## 4. Qué falta — roadmap priorizado

## Fase A — Seguridad y contratos de aplicación

### A1. Autorización para acciones destructivas — boundary resuelto

La autorización pertenece al consumer/application, no al paquete UI. Gelium no debe conocer sesiones, roles, tenants, ownership ni políticas de negocio.

La recipe demo conserva su funcionalidad mediante una política explícita `allow-all`, pero expone un boundary testeable `(request, action, record) -> bool`:

- `401` lo resuelve el middleware de autenticación del consumer;
- `403` lo resuelve la policy/gate/hook del consumer;
- las acciones no autorizadas no se renderizan;
- las rutas vuelven a verificar autorización;
- bulk delete revalida cada registro justo antes de mutar;
- selecciones mixtas producen un resultado parcial explícito.

La auditoría (`actor`, timestamp, motivo, IDs y resultado) sigue siendo responsabilidad de la aplicación y queda fuera de Gelium.

**Archivos implementados en esta recipe demo:**

- `internal/app/server.go`
- `internal/app/recipe_admin_resource.go`
- `internal/app/recipe_admin_resource_test.go`
- `site/web/templates/recipe-admin-resource.html`
- `docs/gelium-ui-screen-recipes.md`

No inventar un middleware genérico desde Gelium: autorización pertenece al consumer/application. InnovaCMS debe conectar Laravel Policies/Gates; InnovaGO ya cuenta con el patrón equivalente `admin.Allows`.

### A2. Auditoría de acciones — contrato documentado

El contrato genérico de integración está documentado en `docs/gelium-ui-application-integration.md`.

La aplicación consumer debe documentar o emitir:

```text
actor
action
resource IDs
result
partial failures
timestamp
```

Gelium comunica el resultado mediante sus primitives de feedback, pero no agrega autenticación, autorización ni persistencia de auditoría dentro del package UI.

## Fase B — Profundizar Record Detail

### B1. Segundo consumidor de detalle — implementado

Se implementó el detalle de un Ops Queue item en `GET /recipes/ops-queue/{id}`. El detalle reutiliza la anatomía semántica, estados, tokens y contratos server-first de la recipe Admin Resource sin extraer todavía `ui-infolist`.

Objetivo: comprobar si `<dl>` + grouped sections + status badge se repite con el mismo contrato.

Solo después considerar extraer:

```text
ui-infolist
```

La extracción debe esperar hasta que existan:

- dos consumidores;
- misma anatomía;
- mismos estados;
- mismos tokens;
- mismos requisitos responsive/a11y;
- tests comunes.

### B2. Relaciones y campos opcionales — contrato resuelto

El detalle de Ops Queue documenta y prueba:

- valores ausentes (`Unassigned`);
- timestamps con timezone explícita;
- textos largos con wrapping seguro;
- ausencia de relaciones inventadas;
- ausencia de campos sensibles por defecto;
- permisos por registro/campo delegados al consumer.

Nunca mostrar campos por defecto solo porque existen en el modelo.

## Fase C — Dashboard

### C1. Data contract de métricas — contrato documentado

El contrato genérico está documentado en `docs/gelium-ui-dashboard-metrics.md`.

El KPI dashboard actual sigue usando datos del demo store y no se presenta como producción. Antes de agregar charts o métricas reales deben definirse fuente autorizada, período, timezone, unidad, freshness, delta, permisos y estados empty/error.

### C2. Dashboard con filtros

Solo si existe una tarea real:

- período (`from`, `to`);
- status;
- owner;
- organización/tenant.

Toda vista navegable debe ser GET y deep-linkable.

### C3. Charts

Diferidos. Antes requieren decidir:

- formato accesible equivalente;
- fallback textual/tabular;
- dimensiones y responsive;
- fuente de datos;
- no-JS path.

## Fase D — Feedback integrado — recipe documentada

Crear documentación/recipe de decisión para:

| Situación | Feedback |
|---|---|
| Resultado transitorio | Toast |
| Éxito persistente después de POST | Banner |
| Error contextual | Inline alert |
| Validación de formulario | 422 + validation summary |
| Error de recurso | Error state |
| Sin datos | Empty state |

La matriz de decisión y la composición server-rendered están documentadas en `docs/gelium-ui-feedback-recipe.md`.

La recipe consolida Toast, Banner, Inline alert, 422 + validation summary, Error state y Empty state; además fija fallback no-JS, accesibilidad y el boundary consumer-owned. Los primitives existentes no cambian.

## Fase E — Extensibilidad tipo Filament — contrato documentado

Antes de pensar en plugins runtime, documentar:

- cómo crear una recipe;
- cómo conectar una fuente de datos;
- cómo declarar acciones;
- cómo respetar estados;
- cómo registrar un componente nuevo;
- cómo extender un theme;
- qué contratos son propiedad de Gelium y cuáles del consumer.

**Referencias existentes:**

- `lib/skills/14-component-implementation.md`
- `docs/gelium-ui-component-registry.md`
- `docs/gelium-ui-pattern-registry.md`
- `docs/gelium-ui-screen-recipes.md`
- `lib/README.md`

La guía de extensibilidad está documentada en `docs/gelium-ui-extensibility.md`.

No crear un sistema de plugins de ejecución hasta que exista un consumer real.

## Fase F — Package/release hygiene — gate verificado

La gate actual fue verificada sin publicar: `npm run release:check` pasó para `gelium-ui@0.6.5`. Antes de publicar una versión:

```bash
go test ./...
go vet ./internal/... ./site/... ./lib/...
npm run build
bash scripts/ux-detect.sh
git diff --check
npm run release:check
```

Además:

- separar commits por slice lógico;
- revisar `.atl/` y `.pi/` antes de incluirlos;
- no limpiar archivos dirty ajenos;
- no incluir `.env`, credenciales ni stores de auth;
- no publicar sin solicitud explícita.

---

## 5. Próximo trabajo recomendado

El gap auditado está en `docs/plans/2026-09-05-filament-gap-audit.md`.

El boundary multi-tenant y el template de testing transversal ya están documentados en `docs/gelium-ui-application-integration.md` y `docs/gelium-ui-recipe-testing.md`.

El próximo slice de implementación debe ser **uno** de estos problemas respaldados por un consumer real:

1. Global search server-rendered para Ops Queue — implementado como primer consumer.
2. Export seguro — contrato documentado y demo CSV acotada en Ops Queue; export production/async pendiente de consumer real.
3. Import seguro — contrato documentado y demo CSV síncrona en Admin Resource; import production/async pendiente de consumer real.
4. Relaciones/nested resources — contrato documentado y demo explícita Project → Tasks en Admin Resource; persistencia/policy productiva sigue siendo consumer-owned.

No implementar todos a la vez. Elegir un caso real, escribir primero el contrato framework-neutral y recién después construir la recipe.

---

## 6. Verificación conocida al momento de este handoff

Las verificaciones completas ejecutadas durante los slices fueron:

- `go test ./...` — PASS;
- `go vet ./internal/... ./site/... ./lib/...` — PASS;
- `npm run build` — PASS;
- `bash scripts/ux-detect.sh` — PASS;
- `git diff --check` — PASS;
- `npm run release:check` — PASS para `gelium-ui@0.6.5`;
- audit de capacidades contrastado con la documentación pública de Filament;
- review Gentle AI del candidato actual — BLOQUEADA por `candidate-view-invalid`; RDD queda desactivado solo a nivel de este clon hasta resolver el runtime.

El worktree contiene cambios acumulados de los slices anteriores, cambios previos de landing/footer y archivos de tooling (`.atl/`, `.pi/`). Otro agente debe revisar `git status --short` antes de editar y no asumir que todos los cambios pertenecen a una sola tarea.

---

## 7. Regla de handoff

Otro agente debe:

1. Leer este roadmap.
2. Leer `AGENTS.md`, `lib/AGENTS.md` y las skills aplicables.
3. Ejecutar `git status --short`.
4. Elegir un solo slice de la sección “Próximo trabajo recomendado”.
5. Hacer RED → GREEN → refactor.
6. Pedir review independiente Gentle AI.
7. Ejecutar el baseline completo.
8. No commit/push/publish/deploy sin autorización explícita.
