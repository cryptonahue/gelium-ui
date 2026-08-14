# Gelium UI — Prompt para implementar un componente con otra IA

Este documento es un prompt reusable. Copia desde “PROMPT COPIABLE”, completa los parámetros y entrégalo a la IA responsable de exactamente un componente.

## Cómo usarlo

1. Elegí un componente de `COMPONENT-ROADMAP.md`.
2. Completá todos los placeholders de asignación `{{...}}`. El estado final lo elige el worker después de ejecutar.
3. Asigná un modo de concurrencia y ownership explícito.
4. No permitas que dos workers escriban el mismo archivo físico.
5. Designá un único integrador para el checkout canónico.
6. El worker termina con uno de los cuatro estados permitidos: `COMPLETE_AWAITING_USER_ACCEPTANCE`, `READY_FOR_INTEGRATION`, `BLOCKED` o `ABORTED_ON_DRIFT`; nunca declara “terminado” sin evidencia.

## Modos de concurrencia

- `SHARED_HANDOFF`: el worker sólo crea/edita archivos nuevos o exclusivos. Para archivos compartidos entrega patches; no los aplica.
- `ISOLATED_PHYSICAL_WORKSPACE`: trabaja en una copia física autorizada, nunca en el checkout canónico. Entrega manifests y patches.
- `EXCLUSIVE_INTEGRATION`: único modo que permite modificar archivos compartidos canónicos; requiere reserva literal vigente para cada path.

Hoy la opción recomendada para varias IAs es investigación paralela + workspaces aislados + integración serial.

---

# PROMPT COPIABLE

````text
# Agente de componente Gelium UI — contrato operativo seguro

Implementa exactamente UN componente de Gelium UI siguiendo este contrato.

## 0. Parámetros

COMPONENTE: {{COMPONENT_NAME}}
SLUG: {{COMPONENT_SLUG}}
CATEGORÍA: {{CORE | LABS | GELIUM_ONLY}}
VARIANTES REQUERIDAS: {{REQUIRED_VARIANTS}}
ESTADOS REQUERIDOS: {{REQUIRED_STATES}}
COMPORTAMIENTOS REQUERIDOS: {{REQUIRED_BEHAVIORS}}
FLUJO SERVER/HTMX, SI APLICA: {{SERVER_FLOW_OR_NONE}}
CRITERIOS ESPECÍFICOS DE ACEPTACIÓN: {{ACCEPTANCE_CRITERIA}}
MATRIZ DE NAVEGADORES SOPORTADA: {{SUPPORTED_BROWSER_MATRIX}}
FECHA/SNAPSHOT BASELINE WEB: {{WEB_BASELINE_SNAPSHOT_DATE}}

REPOSITORIO CANÓNICO: D:\repos\loom-ui
MATERIAL WEB UPSTREAM: D:\repos\material-web-upstream
FORK DE REFERENCIA: D:\repos\material-web-tailwind
UPSTREAM SNAPSHOT ID PROVISTO POR COORDINADOR: {{UPSTREAM_SNAPSHOT_ID}}
UPSTREAM MANIFEST/HASH PROVISTO: {{UPSTREAM_MANIFEST_OR_HASH}}
ROADMAP: D:\repos\loom-ui\COMPONENT-ROADMAP.md
AUDITORÍA HISTÓRICA: D:\repos\loom-ui\MATERIAL-WEB-PROGRESS.md

MODO ASIGNADO: {{SHARED_HANDOFF | ISOLATED_PHYSICAL_WORKSPACE | EXCLUSIVE_INTEGRATION}}
WORKSPACE AUTORIZADO: {{AUTHORIZED_WORKSPACE_PATH}}
WORKER ID: {{WORKER_ID}}
RESERVA ID: {{RESERVATION_ID}}
PATHS POSEÍDOS: {{OWNED_PATHS}}
PATHS PROHIBIDOS ADICIONALES: {{FORBIDDEN_PATHS_OR_NONE}}
BASELINE SHA-256: {{EXPECTED_HASHES_OR_NONE}}
NUEVA VERSIÓN DE ASSETS PROPUESTA: {{NEW_ASSET_VERSION_OR_INTEGRATOR_OWNED}}

No amplíes el alcance a un segundo componente. Registra necesidades transversales no indispensables como follow-up.

## 1. Resultado obligatorio

Preserva por separado:

1. anatomía visual Material;
2. interacción y precedencia de estados;
3. semántica y accesibilidad;
4. comportamiento de datos/forms;
5. arquitectura Gelium.

Cambiar Lit/Web Components por HTML server-rendered NO autoriza simplificar el diseño Material, omitir estados ni alterar el contrato visual.

Stack objetivo:

- Go `net/http` y `html/template`;
- `embed` y Markdown interno confiable;
- Tailwind CSS 4 y CSS propio;
- tokens públicos `--ui-*`;
- HTMX local sólo como progressive enhancement;
- HTML/CSS moderno antes que JavaScript.

## 2. Restricciones no negociables

NO uses ni introduzcas:

- React;
- Lit;
- Shadow DOM;
- Astro;
- `templ`;
- JSX;
- Custom Elements como runtime obligatorio;
- CDN;
- dependencias externas innecesarias;
- JavaScript del componente sin una brecha de plataforma demostrada.

No agregues JavaScript por conveniencia ni para reemplazar semántica nativa.

No inicialices ni uses Git:

- no `git init`, status, diff, branches, commits, worktrees, stash, reset, checkout, fetch, pull o push.

No leas, imprimas, solicites ni modifiques:

- credenciales, tokens, cookies, claves;
- `.env`;
- credential stores;
- autenticación Git/npm.

Repositorios read-only:

- `D:\repos\material-web-upstream`;
- `D:\repos\material-web-tailwind`;
- `D:\repos\material-web-tailwind\material-tailwind`.

Documentos protegidos, read-only salvo tarea explícita distinta:

- `MATERIAL-WEB-PROGRESS.md`;
- `PROMPT-MATERIAL-WEB-INVENTORY.md`;
- `COMPONENT-ROADMAP.md`;
- `AI-COMPONENT-IMPLEMENTER-PROMPT.md`.

No inventes compatibilidad, rutas, tests, reviews ni resultados. Usa `UNKNOWN` cuando corresponda.

## 3. Seguridad de procesos y puertos

La app aceptada puede estar en `http://localhost:8787` y puede existir `gelium.exe`.

Reglas absolutas del worker:

- nunca inicies nada en `:8787`;
- nunca detengas, reinicies o señales el proceso de `:8787`;
- nunca sobrescribas, reconstruyas, ejecutes o reemplaces `gelium.exe`;
- nunca mates un proceso ajeno;
- no pruebes cambios contra `:8787`;
- usa exclusivamente `:8788` para smoke y sólo si está libre.

La terminal de Windows corre Bash/MSYS: usa sintaxis POSIX, no PowerShell.

Inspección inicial read-only permitida:

```bash
netstat -ano | grep -E '[:.]8787[[:space:]]' || true
netstat -ano | grep -E '[:.]8788[[:space:]]' || true
tasklist.exe /FI "IMAGENAME eq gelium.exe" || true
curl -fsS --max-time 3 http://localhost:8787/healthz || true
```

No derives acciones destructivas de esos resultados.

Para smoke, inicia `PORT=8788 go run ./cmd/gelium` con la capacidad de procesos background de tu herramienta, conserva su session/process ID y detén únicamente ese proceso al finalizar. No uses `&`, `nohup` ni procesos huérfanos.

Si `:8788` está ocupado por un proceso que no iniciaste:

1. no lo mates;
2. no lo reutilices;
3. reporta `BLOCKED_PORT_8788`;
4. pide decisión al coordinador.

## 4. Concurrencia sin Git

### 4.1 Usa sólo el modo asignado

#### `SHARED_HANDOFF`

- Sólo edita archivos nuevos o exclusivos listados en `PATHS POSEÍDOS`.
- No edites archivos compartidos.
- Para cada cambio compartido entrega un integration manifest y patch textual.
- Si falta integración, termina como `READY_FOR_INTEGRATION`, no `COMPLETE`.

#### `ISOLATED_PHYSICAL_WORKSPACE`

- Trabaja únicamente en `WORKSPACE AUTORIZADO`.
- No escribas en `D:\repos\loom-ui`.
- Builds, tests y `:8788` ocurren en la copia aislada.
- No copies automáticamente resultados al canónico.
- Entrega manifest completo y patches reproducibles.

#### `EXCLUSIVE_INTEGRATION`

- Lista antes de escribir cada archivo compartido a modificar.
- Verifica que todos aparezcan literalmente en la reserva.
- Captura hashes baseline.
- Una reserva del componente no implica ownership global.
- Si falta ownership de un path, no lo edites.

### 4.2 Compartidos por defecto

Salvo reserva explícita, no edites concurrentemente:

```text
internal/app/server.go
internal/app/server_test.go
web/templates/layout.html
web/styles/app.css
themes/theme-material/theme.css
web/styles_contract_test.go
web/static/app.css
web/static/app.js
web/static/htmx.min.js
web/assets.go
README.md
package.json
package-lock.json
go.mod
go.sum
cmd/gelium/main.go
cmd/gelium/main_test.go
gelium.exe
```

También son compartidos router/mux, page view central, navegación, registry, layout, bundles, índices y archivos generados.

### 4.3 Ownership y drift

Antes de escribir, presenta:

| Ruta | Nueva/existente | Exclusiva/compartida | Owner | Reserva | Acción |
|---|---|---|---|---|---|

Para cada archivo existente autorizado:

1. registra SHA-256 y tamaño;
2. vuelve a calcular el hash inmediatamente antes de escribir;
3. vuelve a comprobar antes del handoff;
4. compara contra tu último estado conocido.

En Bash/MSYS:

```bash
sha256sum path/to/file
wc -c < path/to/file
```

Si cambió sin una escritura tuya:

- detente;
- no sobrescribas ni fusiones automáticamente;
- no reviertas trabajo ajeno;
- conserva tu propuesta como handoff;
- reporta `ABORTED_ON_DRIFT`, path, hash esperado y observado.

### 4.4 Integration manifest

Para cada shared change no aplicado entrega:

```text
FILE: ruta
BASELINE_SHA256: hash
OWNER_REQUIRED: integration-owner
ANCHOR: texto estable y único
PURPOSE: motivo
DEPENDENCIES: contratos relacionados
PATCH:
...diff o bloque exacto...
POSTCONDITION:
...resultado esperado...
TESTS THAT PROVE IT:
...tests/comandos...
```

`STALE_REBASE_REQUIRED` no es un estado final del worker: es una anotación que el coordinador aplica a un handoff previamente emitido cuando su baseline dejó de ser vigente. Si el worker detecta drift durante su ejecución, su estado final es `ABORTED_ON_DRIFT` y no genera un patch contra el baseline cambiado.

## 5. Descubrimiento obligatorio — todavía no escribas

1. Lee `README.md`, `COMPONENT-ROADMAP.md` y `MATERIAL-WEB-PROGRESS.md`.
2. Inspecciona componentes Gelium existentes similares.
3. Inspecciona Go, templates, tests, CSS, theme, docs, build y versión de assets relevantes.
4. Reporta modo, ownership, paths protegidos y estado de puertos/procesos.
5. Distingue hechos, inferencias, `UNKNOWN` y decisiones que requieren aprobación.

No escribas producción durante esta fase.

## 6. Evidencia Material upstream

No implementes desde memoria, capturas o resúmenes genéricos.

Ubica rutas literales del componente concreto para:

1. documentación y demos;
2. render/source;
3. primitives compartidas;
4. Sass/CSS;
5. tokens públicos;
6. geometría hard-coded;
7. tests de comportamiento/accesibilidad.

Registra:

- `UPSTREAM SNAPSHOT ID` provisto por el coordinador; no afirmes haber observado un commit mediante Git porque Git está prohibido;
- verificación contra `UPSTREAM MANIFEST/HASH` cuando fue provisto;
- si falta o no puede verificarse el snapshot requerido, usa `UNKNOWN`/`BLOCKED` y pide al coordinador un baseline;
- rutas exactas;
- variantes y estados;
- anatomía, dimensiones, spacing, type, colors, motion;
- teclado, foco, ARIA y dismissal;
- qué pertenece a Lit/Shadow DOM y no se porta;
- divergencias propuestas.

Material Web es referencia, no dependencia runtime. Reimplementa contratos; no copies automáticamente CSS/código del fork sin revisar procedencia/licencia.

## 7. Auditoría platform-first obligatoria

Antes de proponer JavaScript crea:

| Capacidad | HTML/CSS nativo | Compatibilidad actual | Baseline | No-JS | Gap real | Solución |
|---|---|---|---|---|---|---|

Audita según aplique:

- elementos/atributos HTML;
- forms y navegación;
- selectors/properties CSS;
- Popover API, Invoker Commands, top layer, inertness y dismissal;
- keyboard y focus;
- anchor positioning;
- reduced motion y forced colors;
- RTL y responsive;
- semántica server/network;
- browser probes reales.

Fuentes, en orden:

1. WHATWG/W3C;
2. MDN;
3. Browser Compat Data;
4. Web Features/Baseline;
5. probe real.

Incluye fecha y evalúa cada capacidad contra `MATRIZ DE NAVEGADORES SOPORTADA` y `FECHA/SNAPSHOT BASELINE WEB` provistos. Si falta alguno, no elijas silenciosamente una matriz: marca la decisión `UNKNOWN`/`BLOCKED`. No uses una matriz histórica como si fuera actual.

JavaScript sólo se acepta si queda una brecha funcional imposible de resolver correctamente con:

1. HTML semántico;
2. CSS declarativo;
3. formulario/navegación server-rendered;
4. HTMX como enhancement.

Si sigue siendo necesario:

- documenta la brecha exacta;
- escribe primero un test RED;
- usa módulo vanilla/framework-free mínimo;
- conserva flujo no-JS real;
- no conviertas JS en requisito para leer, navegar, enviar o completar el flujo principal.

Para navegadores sin la primitive moderna, usa una ruta/página server-rendered real. No simules modal, tabs, menú o validación con hacks CSS inaccesibles.

## 8. Especificación antes de código

Produce una matriz:

| Feature | Contrato upstream | Estrategia Gelium | Test | Divergencia |
|---|---|---|---|---|

Debe cubrir:

- raíz semántica y anatomía;
- variants y states;
- combinaciones y precedencia;
- rest/hover/focus-visible/active;
- disabled/loading/error/empty cuando aplican;
- teclado y focus lifecycle;
- labels, names y descriptions;
- forms, values, submission y códigos HTTP;
- no-JS y HTMX;
- light/dark, narrow/wide, RTL;
- reduced motion y forced colors;
- assets y trust boundaries.

Toda divergencia no trivial requiere aprobación ANTES de implementar. No llames paridad a una aproximación no verificada.

## 9. TDD estricto

Ley:

NO PRODUCTION CODE WITHOUT A FAILING TEST FIRST.

Usa ciclos verticales, uno por comportamiento:

1. escribe test mínimo;
2. ejecuta sólo ese test;
3. observa y registra el RED esperado;
4. confirma que no es typo/setup;
5. implementa lo mínimo;
6. observa GREEN;
7. ejecuta suite relacionada;
8. refactoriza sólo en verde;
9. repite.

No hagas todos los tests y después toda la implementación.

Mantén este log:

| Ciclo | Test | Comando RED | Fallo observado | Cambio GREEN | Comando GREEN | Resultado |
|---|---|---|---|---|---|---|

Si un test pasa al primer intento, no cuenta como RED; indica que el comportamiento ya existía o corrige el test. Nunca borres o reviertas trabajo ajeno.

Comandos base:

```bash
go test ./internal/app -run 'TestNombreExacto' -count=1 -v
go test ./... -count=1
go vet ./...
go mod verify
```

Usa primero tests normales. Si Windows Application Control bloquea ejecutables temporales, reporta el bloqueo; sólo usa un `GOTMPDIR` interno si el coordinador lo autoriza. No borres `.tmp` automáticamente.

## 10. Templates, atributos y trust

Usa `html/template` y view models Go concretos.

Obligatorio:

- vocabularios cerrados para variant/type/command/placement;
- defaults explícitos;
- booleans reales para atributos booleanos;
- strings escapables para texto/values;
- atributos emitidos individualmente;
- tests de escaping, omisión e inputs inválidos.

Prohibido:

- `map[string]any` como API;
- attrs/raw HTML strings;
- `template.HTMLAttr` dinámico;
- clases arbitrarias del caller;
- nombres de tags/attrs controlados por usuario;
- `template.HTML` o `template.URL` sobre contenido no confiable.

`template.HTML` sólo para markup interno confiable, con comentario de trust boundary y tests. SVG decorativo trusted debe incluir `aria-hidden="true"` y `focusable="false"`; el texto visible aporta el nombre accesible.

## 11. CSS, tokens y accesibilidad

- Elementos nativos antes que ARIA.
- No uses `div`/`span` para reemplazar controles.
- No añadas ARIA redundante.
- Todos los tokens públicos usan `--ui-*`.
- Los mappings Material viven en el theme.
- No agregues familias de tokens que el slice no necesita.
- Define precedencia explícita de estados.
- El foco no debe cambiar geometría ni producir layout shift.
- No elimines outlines sin indicador equivalente.

Prueba:

- light/dark;
- narrow/wide;
- RTL cuando aplica;
- reduced motion;
- forced colors;
- disabled y estados combinados;
- texto/error no comunicado sólo por color;
- contraste antes de afirmar WCAG.

## 12. Flujo completo sin JavaScript

El flujo principal acordado debe completarse con JS/HTMX deshabilitado dentro de la matriz soportada.

Según el caso prueba:

- `href` real;
- form `method`/`action`;
- respuesta HTML completa;
- values/errors preservados;
- status HTTP correcto;
- anchor de retorno;
- acción server-rendered equivalente.

HTMX puede recibir fragments, pero debe existir rama no-HX. Testea requests normales y `HX-Request: true` por separado.

No afirmes “sin JS” si sólo se ve, no completa la acción, depende exclusivamente de `hx-*` o imita semántica con CSS.

## 13. Docs dogfoodeadas

Crea `/components/{{COMPONENT_SLUG}}` usando la implementación real, nunca markup duplicado.

Documenta:

- propósito y anatomy;
- variants/states;
- accesibilidad;
- comportamiento no-JS;
- HTMX si aplica;
- compatibilidad/Baseline;
- trust boundary;
- divergencias;
- checklist visual.

En `SHARED_HANDOFF`, route/nav/layout se entregan como integration manifest.

## 14. Build y assets

Source CSS: `web/styles/app.css` o el módulo componente asignado.
Outputs actuales de `npm run build`:

- `web/static/app.css` — compilado por Tailwind;
- `web/static/htmx.min.js` — sobrescrito por `scripts/copy-htmx.mjs`.

Sólo ejecuta `npm run build` si trabajas en copia aislada o tienes ownership exclusivo de TODOS los outputs reales del script y de cualquier archivo de versionado que deba cambiar. Poseer sólo `web/static/app.css` no autoriza la build completa. Nunca ejecutes build en `SHARED_HANDOFF`. Si sólo posees source CSS, entrega el build a la integration/build lane.

Todo cambio de asset embebido requiere URL versionada nueva o content hash. `Cache-Control` solo no alcanza.

Prueba:

1. HTML referencia versión nueva;
2. CSS servido contiene marker del componente;
3. package/docs/layout no divergen;
4. build reproducible;
5. no CDN;
6. output corresponde al source.

No construyas `gelium.exe`.

## 15. Reviews separadas

### Gate A — Spec review

Después de tests verdes revisa sólo contra parámetros, matriz, upstream, no-JS y aceptación. Marca cada requisito:

`PASS | FAIL | BLOCKED | UNKNOWN`

Incluye evidencia path/test. Corrige FAIL con TDD y repite hasta PASS.

### Gate B — Quality review

Sólo tras Spec PASS revisa seguridad, escaping, trust, semántica, accesibilidad, CSS, mantenibilidad, ownership/drift, scope, deps, assets, docs y fragilidad de tests.

Corrige Critical/Important con TDD y repite hasta APPROVED.

Usa reviewers independientes cuando tu plataforma lo permita. Si no, declara la limitación; no inventes aprobación independiente.

## 16. Verificación final

Si estás autorizado:

```bash
npm run build
go test ./... -count=1
go vet ./...
go mod verify
node --check web/static/app.js
```

No ocultes warnings. Si el build toca un path sin ownership, detente y reporta scope violation.

## 17. Smoke en :8788

Sólo tras build/tests y desde workspace autorizado:

1. confirma `:8788` libre;
2. inicia servidor propio en background;
3. espera `/healthz`;
4. prueba ruta, assets y versión;
5. prueba flujo no-HX;
6. prueba HX si aplica;
7. prueba status/headers;
8. browser real y consola;
9. detén exclusivamente tu proceso.

Valida light/dark, narrow/wide, teclado, focus, disabled/error, reduced motion, forced colors cuando sea viable y comparación visual upstream.

Nunca uses `:8787` para probar cambios.

## 18. Aceptación del usuario

Entrega checklist observable y específica. Como mínimo:

- anatomía y variantes;
- rest/hover/focus/active;
- disabled/loading/error;
- keyboard/focus;
- no-JS;
- HTMX como enhancement;
- light/dark;
- narrow/wide/RTL;
- reduced motion/forced colors;
- consola;
- docs reales;
- asset versionado;
- confirmación de que `:8787`/`gelium.exe` no fueron tocados.

No integres otro componente ni reemplaces la app estable. Espera aceptación explícita.

## 19. Estado final permitido

Usa exactamente uno:

- `COMPLETE_AWAITING_USER_ACCEPTANCE`: implementación e integración completas, reviews y smoke aprobados;
- `READY_FOR_INTEGRATION`: handoff listo, shared patches aún no aplicados;
- `BLOCKED`: falta evidencia, decisión, ownership o puerto;
- `ABORTED_ON_DRIFT`: cambió un baseline concurrentemente.

`SHARED_HANDOFF` no puede declarar COMPLETE si quedan patches.

## 20. Entrega obligatoria

### Estado
`[ONE_ALLOWED_STATUS_CHOSEN_BY_WORKER]`

### Resumen
Componente, alcance, modo y resultado.

### Seguridad de app
Estado de `:8787`, `gelium.exe`, `:8788` y proceso propio.

### Evidencia upstream
Snapshot ID/manifest provistos, paths y divergencias.

### Auditoría platform-first
Fecha, tabla y gap JS o “no requerido”.

### Ownership y drift
Paths, reservas, hashes baseline/final y drift.

### Log TDD
Cada RED/GREEN con comandos y resultados reales.

### Archivos
Separar creados, modificados, generados, propuestos no aplicados y protegidos verificados.

### Integration manifest
Patches para shared files no editados.

### Reviews
Spec y Quality con findings/resolución.

### Verificación
Build, tests, vet, mod verify, smoke, browser, no-JS y HTMX.

### Checklist del usuario
Casos observables.

### Issues/blockers
UNKNOWN, riesgos y trabajo pendiente.

### Confirmación
Declara explícitamente: no Git, no credenciales, no CDN, no React/Lit/Shadow DOM/Astro/templ, no JS sin gap, upstream read-only, `:8787`/`gelium.exe` intactos, sin overwrite concurrente y sin resultados inventados.
````

## Plantilla mínima de asignación

```text
COMPONENTE: Checkbox
SLUG: checkbox
CATEGORÍA: CORE
VARIANTES REQUERIDAS: checked, unchecked, indeterminate si se aprueba
ESTADOS REQUERIDOS: rest, hover, focus, pressed, disabled, error
COMPORTAMIENTOS REQUERIDOS: label clickeable, form submit normal, keyboard nativo
FLUJO SERVER/HTMX, SI APLICA: submit normal obligatorio; HTMX opcional
CRITERIOS ESPECÍFICOS DE ACEPTACIÓN: completar tras auditoría upstream/platform-first
MATRIZ DE NAVEGADORES SOPORTADA: {{COORDINATOR_MUST_DEFINE}}
FECHA/SNAPSHOT BASELINE WEB: {{COORDINATOR_MUST_DEFINE}}
UPSTREAM SNAPSHOT ID PROVISTO POR COORDINADOR: {{COORDINATOR_MUST_DEFINE}}
UPSTREAM MANIFEST/HASH PROVISTO: {{COORDINATOR_MUST_DEFINE}}
MODO ASIGNADO: SHARED_HANDOFF
WORKSPACE AUTORIZADO: D:\repos\loom-ui
WORKER ID: checkbox-worker-01
RESERVA ID: checkbox-exclusive-files-only
PATHS POSEÍDOS: web/templates/checkbox.html, web/content/checkbox.md
PATHS PROHIBIDOS ADICIONALES: todos los shared files; entregar patches
BASELINE SHA-256: hashes provistos por integrador
NUEVA VERSIÓN DE ASSETS PROPUESTA: integrator-owned
```

## Nota de coordinación

Mientras Gelium no complete la Wave P del roadmap, asignar dos componentes a dos IAs NO significa permitir dos escritores sobre `D:\repos\loom-ui`. Los workers preparan contratos/artifacts en paralelo; un integrador único incorpora cada lane de forma serial.
