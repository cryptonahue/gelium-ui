# Presentamos Gelium UI (borrador en español)

Borrador de trabajo del post de lanzamiento. La versión publicada en el sitio
(`site/web/content/blog-introducing.md`) es la traducción inglesa de este texto.
Slug: `introducing`. Fecha: 2026-08-21. Autor: Gelium UI team.

---

# Presentamos Gelium UI

Gelium UI ya está disponible: una librería de componentes themeable y de código
abierto para aplicaciones server-rendered. HTML semántico que copias, temas por
tokens que se cambian con una clase, y un contrato estricto: cero JavaScript
obligatorio. No es un eslogan; hay 770 tests que lo verifican.

## El problema

La mayoría de las librerías de componentes asumen que tienes una app cliente.
Piden un runtime de React o Vue, hidratación y un bundle de JavaScript. Si tu
app renderiza HTML desde Go, Rails, Laravel o Django, las opciones se reducen.
Gelium nació de esa pregunta: ¿puede una app server-rendered tener componentes
bonitos, accesibles y themeable sin shipping JavaScript?

## Qué es

Gelium UI es open-code. Los componentes son parciales de HTML semántico que
copias a tu proyecto, estilados con tokens CSS (`--ui-*`). No hay runtime
cliente, ni hidratación, ni CDN. Está construida sobre Tailwind CSS 4 y se
instala desde npm: `npm install gelium-ui`. HTMX está integrado para quien lo
quiera, pero nunca es un requisito.

## El contrato 0-JS

Cero JavaScript requerido es un contrato duro, no un claim de marketing. Los
tests escanean el bundle y el markup para demostrar que ningún componente
depende de un script. La semántica nativa va primero: los botones son
`<button>`, los diálogos son `<dialog>` reales, y los estados mapean a
atributos reales. El JavaScript que existe es mejora progresiva: alrededor de
5KB opcionales. El resultado es UI que sobrevive aunque los scripts fallen o
tarden en cargar.

## Temas por tokens

Los temas son valores de tokens, no estilos de componentes separados. Una
clase en `<html>` cambia Material 3 por Basecoat, Linear, Vercel o Alden, sin
rebuild y sin tocar el markup. La traducción de Basecoat es honesta: el
tarball oficial está auditado y cada style pack (Vega, Nova, Mira, Lyra y el
resto) se tradujo token a token, con las divergencias documentadas. Donde no
hay evidencia oficial, como los visuales de Base UI, el tema se etiqueta como
"inspirado en la docs, autoría de Gelium".

## Un componente, muchos looks

La arquitectura separa comportamiento, referencia visual y skin. En la
práctica: un mismo componente semántico se viste como Material, Basecoat o
Linear sin bifurcar el markup. La densidad también es un contrato: por
defecto, los touch targets nunca bajan de 44px, aunque el skin de turno pida
más compacto. La accesibilidad no se negocia por estética.

## Docs como producto

La documentación no es un accesorio. Cada página renderiza el componente real
que documenta, con el mismo layout y los mismos tokens. El catálogo cubre 48
componentes, de botones a data tables, con estados vacíos, de error y de
carga. Los checkmarks del roadmap significan "verificado por tests", no
aspiracional.

## Para quién es (y para quién no)

Gelium es para apps server-rendered: Go, Rails, Laravel, Django o HTML plano.
Para equipos que quieren HTMX o mejora progresiva. Para sistemas de diseño que
deben funcionar con JavaScript deshabilitado. No es para SPAs de React o Vue,
ni para Web Components con Shadow DOM. No pretende reemplazarlos; es otro
camino.

## Qué sigue

El proyecto está en pre-1.0 y avanza rápido: galería de packs de iconos,
localización de docs y más screen recipes. El código es MIT y está en GitHub.
Si construyes apps server-rendered, dale una mirada y cuéntanos qué piensas.
