# Loom UI

> Themeable, open-code UI components for Tailwind CSS and HTMX.

Current release: **v0.4.0**, adding a Loom-only Toast component for server-driven, transient feedback — an accessible `aria-live` region, an HTMX `loom:toast` trigger contract, a complete no-JS inline fallback, and a minimal framework-free auto-dismiss enhancement.

Este primer vertical slice es una aplicación de documentación server-rendered en Go. Usa HTML semántico, componentes copiables, un theme Material basado en tokens propios `--ui-*`, Tailwind CSS 4 y HTMX servido localmente.

## Requisitos

- Go 1.24 o posterior.
- Node.js 20 o posterior con npm.

## Setup y build reproducible

Desde la raíz del proyecto:

```bash
npm ci
go mod download
npm run build
```

`npm run build` compila `web/styles/app.css` con Tailwind CSS 4 hacia `web/static/app.css` y copia `htmx.min.js` desde la dependencia npm `htmx.org` hacia `web/static/`. La integración local vive en `web/static/app.js`. No se usa CDN.

Los archivos finales dentro de `web/static/` se mantienen como artefactos de build y son embebidos en el binario Go mediante `embed`.

## Tests y checks

```bash
go test ./...
go vet ./...
```

Los tests usan `httptest` y verifican comportamiento de handlers, forms URL-encoded, HTTP 422/200, render Markdown, contratos HTML/ARIA de Button y Text field, integración HTMX, CSS/tokens y assets embebidos sin snapshots completos.

## Ejecutar

```bash
go run ./cmd/loom
```

La aplicación escucha en `http://localhost:8787`. Para elegir otro puerto:

```bash
PORT=3000 go run ./cmd/loom
```

Rutas disponibles:

- `/` — home renderizada desde Markdown.
- `/components/button` — documentación y preview real del Button.
- `/components/text-field` — documentación, estados y demo HTMX real del Text field.
- `/components/dialog` — documentación y preview del Dialog nativo; comandos declarativos y light dismiss se aplican en navegadores compatibles.
- `/components/toast` — documentation y demo HTMX del Toast: variantes, contrato `loom:toast` y fallback inline sin JS.
- `POST /examples/text-field/validate` — validación server-side; devuelve 422 para vacío/whitespace y 200 para valores válidos. Sin JavaScript renderiza la página de documentación completa; con `HX-Request: true`, HTMX recibe únicamente el form actualizado como progressive enhancement.
- `POST /examples/toast/demo` — feedback server-driven; devuelve `HX-Trigger: {"loom:toast":{...}}` para HTMX y, sin JavaScript, re-renderiza la página con un toast inline persistente. 422 para mensaje vacío.
- `/healthz` — health check de texto plano.
- `/static/app.css` — CSS Tailwind compilado y embebido.
- `/static/htmx.min.js` — HTMX local y embebido.
- `/static/app.js` — hook local que permite a HTMX intercambiar únicamente fragments HTTP 422 marcados con `X-Loom-Validation: true`.

## Estructura real

```text
loom-ui/
├── cmd/loom/
│   ├── main.go
│   └── main_test.go
├── internal/app/
│   ├── server.go
│   └── server_test.go
├── scripts/
│   └── copy-htmx.mjs
├── themes/theme-material/
│   └── theme.css
├── web/
│   ├── assets.go
│   ├── content/
│   │   ├── button.md
│   │   ├── dialog.md
│   │   ├── index.md
│   │   ├── text-field.md
│   │   └── toast.md
│   ├── static/
│   │   ├── app.css
│   │   ├── app.js
│   │   └── htmx.min.js
│   ├── styles/
│   │   └── app.css
│   └── templates/
│       ├── button.html
│       ├── dialog.html
│       ├── layout.html
│       ├── text-field.html
│       └── toast.html
├── go.mod
├── go.sum
├── package.json
└── package-lock.json
```

### Nota sobre `embed`

Go no permite que una directiva `//go:embed` lea rutas padre con `..`. Por eso `web/assets.go` vive junto al árbol `web/` y embebe `templates`, `content` y los assets compilados de `static`. El source independiente del theme permanece en `themes/theme-material/theme.css`; Tailwind lo integra en `web/static/app.css`, que es el asset final embebido.

## Button open-code

El markup reusable está en `web/templates/button.html` y los estilos en `web/styles/app.css`. Incluye:

- `primary`, `secondary` y `outline`;
- estados disabled y loading;
- enlaces con `Href` se renderizan sin `href`, fuera del tab order y con `aria-disabled` cuando están disabled o loading; loading además expone `aria-busy` y el nombre accesible dinámico `Loading {Label}` tanto en enlaces como en botones;
- focus exterior de 3 px con offset de 2 px;
- slot `IconSVG` por instancia para SVG inline arbitrario. Es `template.HTML`: acepta únicamente markup interno/de confianza, nunca input de usuario. Por contrato, el SVG debe incluir `aria-hidden="true"` y `focusable="false"`; el texto `Label` aporta el texto de acción usado por el accessible name;
- semántica separada para acciones (`button`) y navegación (`a`).

## Dialog open-code

El markup reusable está en `web/templates/dialog.html` y los estilos en `web/styles/app.css`. El elemento `<dialog>` tiene soporte amplio, pero `command`/`commandfor` son Baseline Low, `request-close` es aún más reciente y `closedby` no es Baseline. En navegadores compatibles, los comandos declarativos y `closedby="any"` aportan apertura/cierre y light dismiss sin JavaScript del componente; para navegadores anteriores, el consumidor debe añadir un fallback server-rendered o adaptador. `@starting-style` y las transiciones de cierre/top layer con `overlay` son mejoras progresivas (estas últimas, Chromium-only y no interoperables), por lo que el motion puede ser instantáneo o asimétrico en otros navegadores. Incluye reduced motion, forced colors, nombres accesibles y acciones explícitas de Cancel/Confirm.

## Alcance

El alcance actual incluye Button, Text field, Dialog y Toast open-code, con un form HTMX de validación server-side que reemplaza el form completo tanto en HTTP 200 como en HTTP 422, y un Toast server-driven que funciona sin JS y mejora con HTMX. No incorpora Lit, Shadow DOM, Astro, `templ`, SQLite, SSE, registry, CLI, toasts adicionales ni themes adicionales.

## Licencia

Loom UI se distribuye bajo la [licencia MIT](LICENSE).
