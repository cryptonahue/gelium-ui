# Prompt para IA: primera auditoría Material Web → Gelium UI

Copiar el siguiente prompt en la IA/coding agent que vaya a realizar la primera auditoría.

---

## Prompt

Estás trabajando en el proyecto `Gelium UI`, ubicado en:

```text
D:\repos\loom-ui
```

El objetivo de esta tarea es producir la **primera auditoría técnica y funcional de Material Web** para decidir qué convertir a Gelium UI.

### Fuentes locales obligatorias

Analiza estas tres fuentes:

```text
D:\repos\material-web-upstream
D:\repos\material-web-tailwind
D:\repos\material-web-tailwind\material-tailwind
```

`material-web-upstream` es un clone del repositorio oficial:

```text
https://github.com/material-components/material-web
```

No asumas que el fork local está actualizado. Compara siempre contra el upstream actual.

### Contexto de Gelium UI

Gelium UI será:

```text
Tailwind CSS + HTML/templates + HTMX + themes + server-rendered applications
```

La primera implementación y la documentación se harán en Go usando:

```text
net/http
html/template
embed
Markdown
SQLite
SSE
HTMX
Tailwind CSS
```

No usar Astro ni `templ` en esta fase.

Material Web y Basecoat UI son referencias/inspiraciones. Gelium UI no es un fork oficial de ninguno de los dos.

- Material Web: https://github.com/material-components/material-web
- Basecoat UI: https://basecoatui.com/
- shadcn/ui: referencia para open-code y registry.
- shadcn-templ: referencia para distribución y presets, no para atar Gelium a Go/templ.
- HTMX 4: referencia para interacción server-driven, history, multi-target, validation, SSE y WebSockets.

### Regla de seguridad

No modifiques código de `material-web-upstream`, `material-web-tailwind` ni `material-tailwind`.

En esta fase sólo analiza y escribe el informe dentro de `D:\repos\loom-ui`.

No inventes estados, componentes, versiones, tests o resultados de build. Si algo no se puede confirmar, márcalo como `UNKNOWN` y explica cómo verificarlo.

## Estrategia de subagentes

Usa subagentes en paralelo para analizar áreas independientes. No lances un subagente por cada archivo pequeño: eso genera overhead y resultados inconsistentes. Divide inicialmente en estas tres líneas:

### Subagente A — Foundations y tokens

Analiza:

- `color/`
- `tokens/`
- `typography/`
- `elevation/`
- `focus/`
- `ripple/`
- `common.ts`
- `internal/`
- documentación de theming
- diferencias de tokens entre upstream y fork
- `tokens/versions/latest` y versiones anteriores

Debe devolver:

- inventario de tokens;
- tokens nuevos respecto al fork;
- tokens que conviene mapear a `@theme` de Tailwind;
- tokens que son específicos de Web Components y no conviene portar;
- recomendaciones para el theme Material de Gelium.

### Subagente B — Componentes

Analiza por lotes los componentes core y labs.

Core:

- button
- checkbox
- chips
- dialog
- divider
- elevation
- fab
- field
- focus
- icon
- iconbutton
- list
- menu
- progress
- radio
- ripple
- select
- slider
- switch
- tabs
- textfield

Labs:

- aria
- badge
- behaviors
- card
- item
- navigationbar
- navigationdrawer
- navigationtab
- segmentedbutton
- segmentedbuttonset
- gb

Para cada componente, registra:

- existe en upstream;
- existe en el fork;
- existe en `material-tailwind`;
- variantes;
- estados;
- tokens relevantes;
- requisitos de accesibilidad;
- tests y documentación disponibles;
- complejidad de portarlo a HTML + Tailwind;
- necesidad de JavaScript local;
- posible patrón HTMX;
- prioridad Gelium: `P0`, `P1`, `P2`, `defer` o `skip`.

No confundas un token con un componente implementado.

### Subagente C — Build, documentación, tests y port Tailwind

Analiza:

- `package.json` de upstream y del fork;
- scripts de build;
- Sass y generación de CSS/TS;
- manifiesto de Custom Elements;
- catálogo y documentación;
- tests;
- diferencias de archivos tracked;
- estado de `material-tailwind`;
- qué partes pueden reutilizarse en Gelium sin arrastrar Lit/Web Components;
- qué partes de la documentación y playground pueden inspirar las docs Go.

Debe devolver:

- inventario del pipeline actual;
- diferencias upstream/fork;
- piezas reutilizables;
- piezas que deben reescribirse;
- riesgos técnicos;
- propuesta de pipeline Go + Tailwind + HTMX.

## Coordinación

Cada subagente debe devolver su informe en formato estructurado, sin modificar archivos fuente.

El agente principal debe:

1. revisar los tres informes;
2. resolver contradicciones leyendo los archivos originales;
3. no repetir trabajo que ya esté confirmado;
4. crear una única tabla consolidada;
5. escribir el documento final sólo después de validar los datos.

Si la plataforma no permite subagentes, simula estas tres fases en orden, conservando exactamente la misma estructura de salida.

## Documento de salida

Actualiza o crea:

```text
D:\repos\loom-ui\MATERIAL-WEB-PROGRESS.md
```

Con estas secciones:

1. **Snapshot de fuentes**
   - rutas;
   - remotes;
   - versiones;
   - commits;
   - fechas;
   - conteo de archivos tracked;
   - método de comparación.

2. **Inspiraciones y atribución**
   - Material Web;
   - Basecoat UI;
   - shadcn/ui;
   - shadcn-templ;
   - HTMX.

3. **Inventario de foundations y tokens**

4. **Inventario completo de componentes**
   - una fila por componente;
   - core/labs;
   - estado upstream;
   - estado fork;
   - estado material-tailwind;
   - prioridad Gelium.

5. **Diferencias upstream vs fork**

6. **Qué rescatar de `material-tailwind`**

7. **Qué no portar**
   - Lit;
   - Shadow DOM;
   - Custom Elements como arquitectura obligatoria;
   - controllers que sólo tengan sentido en el runtime de Material Web.

8. **Traducción conceptual a Gelium**
   - Custom Element → HTML + clases/atributos;
   - Lit property → atributo o estado server-side;
   - Lit event → evento HTML/HTMX;
   - Shadow DOM CSS → Tailwind + tokens;
   - imperative JS → HTMX/server response cuando sea posible.

9. **Matriz de progreso Gelium**
   - `NOT_STARTED`;
   - `IN_REVIEW`;
   - `DESIGN_READY`;
   - `IMPLEMENTED`;
   - `DOCS_READY`;
   - `TESTED`;
   - `BLOCKED`.

10. **Primer vertical slice recomendado**

```text
button → form/input → dialog → toast
```

11. **HTMX/realtime backlog**
   - dialog remoto;
   - tabs + history;
   - validación server-side;
   - tables/pagination;
   - `HX-Trigger`;
   - SSE notifications;
   - SSE progress;
   - WebSockets opcionales.

12. **Go MVP**
   - `net/http`;
   - `html/template`;
   - `embed`;
   - Markdown docs;
   - SQLite/WAL;
   - event broker;
   - SSE;
   - playground dogfooded.

13. **Próximos pasos priorizados**

## Formato por componente

Usa una tabla como ésta:

| Componente | Categoría | Upstream | Fork | Material Tailwind | Variantes | Estados | HTMX pattern | Prioridad | Gelium status | Evidencia |
|---|---|---|---|---|---|---|---|---|---|---|
| Button | Core | Sí | Sí | Sí | filled, outlined... | disabled, focus... | form/action | P0 | NOT_STARTED | rutas exactas |

Cada afirmación importante debe incluir una ruta de archivo concreta, por ejemplo:

```text
material-web-upstream/button/filled-button.ts
material-web-upstream/docs/components/button.md
material-web-tailwind/material-tailwind/src/...
```

## Criterios de calidad

- No usar marketing del README como prueba suficiente.
- Preferir código, tests y documentación específica.
- Diferenciar `documented`, `implemented`, `tested` y `generated`.
- Identificar componentes que existen sólo en `labs`.
- Identificar tokens sin componente asociado.
- Indicar cuando una API sólo existe por Lit/Web Components.
- No declarar un port como terminado sólo porque existe CSS parecido.
- No modificar código fuente durante la auditoría.
- Al final, reportar los archivos creados/modificados y los comandos realmente ejecutados.

## Resultado final esperado

Un documento accionable para que el siguiente trabajo pueda comenzar directamente por el primer vertical slice, sin volver a descubrir qué existe en Material Web, qué quedó en el fork viejo y qué piezas sirven para Gelium UI.

---
