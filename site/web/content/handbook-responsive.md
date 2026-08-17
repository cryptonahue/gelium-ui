# Diseñar para tamaños de pantalla, no para dispositivos

El diseño responsive en Gelium parte de **viewports y contenido**, no de “es un iPhone” o “es una tablet”. El mismo HTML y CSS deben refluir con dignidad desde un panel estrecho hasta un monitor ancho. Esta página fija la postura de producto; el detalle de controles está en [Elegir el control correcto](/docs/choose-the-right-control) y en [Forms](/docs/forms).

## Viewports, no nombres de dispositivo

- Diseña y prueba por **ancho de viewport** (y, cuando importe, altura/`dvh`), no por listas de modelos.
- Los breakpoints son **pasos de layout** (`40rem`, `48rem`, …), no etiquetas `phone` / `tablet` / `desktop` en el CSS de producto.
- Un teléfono plegado, un split-screen o una ventana de escritorio estrecha pueden compartir el mismo ancho: el layout no debe asumir “móvil = touch y poco ancho” como un solo paquete inseparable.

## Contenido que refluye

- Prioriza **reflow**: columnas que se apilan, filas con `flex-wrap`, tablas que hacen scroll **dentro** de su contenedor.
- El texto de lectura larga usa la medida tipográfica del docs chrome: **`.prose` con `max-width: 65ch`**. No estires párrafos a todo el ancho del monitor.
- Los layouts de consumo (p. ej. recipes) respetan **`--ui-container-max`** como techo de columna de aplicación, no como disculpa para forzar un ancho de “desktop forever”.

## De móvil hacia arriba (mobile-first)

- Escribe el **stack por defecto** para el viewport estrecho: una columna, acciones apiladas, cabeceras que no pelean por el mismo renglón.
- **Mejora desde desktop** con media queries `min-width` (o capas equivalentes): más columnas, toolbars en fila, paneles lado a lado.
- Evita un CSS “desktop-only” que luego se parchea con excepciones móviles opacas. Las recipes (admin-resource, ops-queue, public-feed) **apilan cabeceras y contienen tablas** en anchos estrechos; ese es el patrón a copiar.

## Contención: no enmascarar el desborde

- **overflow-x: hidden no enmascara** un layout roto: no lo pongas en el documento ni en el body de la pantalla para “ocultar” el desborde. Corta focus, sombras y contenido real.
- En flex/grid, pon **`min-width: 0`** (o `min-inline-size: 0`) en columnas e hijos que deban encogerse. Sin eso, el min-content del hijo empuja el viewport.
- Tablas densas: envuelve el `<table>` en un contenedor con scroll horizontal local (p. ej. `.ui-data-table-scroll`), no en un clip del `html`.
- Mide el **min-content width** real (DevTools, capturas a ~360–400px y en el umbral del breakpoint). Si algo fuerza ~780px, el fallo está en el layout, no en el usuario.

## Objetivos táctiles y formularios

- Los controles interactivos respetan **`--ui-touch-target`** (suelo de área útil; Material lo eleva donde aplica el tema). No inventes hit-areas ad hoc por pantalla.
- Formularios: etiquetas visibles, `type`/`inputmode` nativos, sin bloquear pegar — ver [Forms](/docs/forms).
- Elige el control por la tarea, no por el “aspecto móvil”: [Choose the right control](/docs/choose-the-right-control).

## Breakpoints por pasos

| Enfoque | Hacer | Evitar |
|---|---|---|
| Nombre | Anchos en `rem` u otras unidades de layout | `@media (device-width: 375px)` o “solo iPad” |
| Cantidad | Pocos escalones donde el **layout cambia de forma** | Un breakpoint por cada gadget del lab |
| Orden | Base estrecha → realce ancho | Base ancha → parches `max-width` interminables |
| Prueba | Viewports continuos y ventanas redimensionables | Solo emuladores con chrome de dispositivo fijo |

## Tokens y piezas Gelium

| Pieza | Rol |
|---|---|
| `--ui-touch-target` | Suelo de tamaño táctil / hit area en botones e icon-buttons |
| `--ui-container-max` | Ancho máximo razonable de shell de aplicación / recipe |
| `.prose` + **65ch** | Medida de lectura en documentación y texto largo |
| Recipes + contención | Stack de headers, `min-width: 0`, scroll de tabla local |

Temas y vocabulario de tokens: [Themes](/docs/themes), [Tokens](/docs/tokens). Accesibilidad (contraste, foco, motion): [Accessibility](/docs/accessibility).

## Rendimiento y honestidad de payload

Un layout que no desborda no sustituye una postura de payload. JS es mejora progresiva; el CSS de tokens y temas es grande **a propósito**. Mide y compara con la misma regla en [Performance](/docs/performance) y el posicionamiento en [Why Gelium](/docs/compare).

## Qué no es esta página

- No es un catálogo de device frames ni un simulador de notch.
- No autoriza `overflow-x: hidden` global como “fix responsive”.
- No sustituye las demos de componente ni las recipes: allí se ve el apilado real.

## Lista rápida de verificación

1. ¿El layout se entiende a ~360px de ancho sin scroll horizontal de página?
2. ¿Los breakpoints nombran **cambios de layout**, no marcas de hardware?
3. ¿Hay `min-width: 0` donde flex/grid se encoge?
4. ¿Las tablas y pre scrollean **dentro** de su caja?
5. ¿Los controles cumplen `--ui-touch-target` y el contrato de [Forms](/docs/forms)?
6. ¿El CSS parte del stack estrecho y realza en anchos mayores?
