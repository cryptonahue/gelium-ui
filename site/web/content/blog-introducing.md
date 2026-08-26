# Introducing Gelium UI

Gelium UI is out: themeable, open-code components for server-rendered apps. Copy semantic HTML, pick a theme class, and ship. The hard contract: zero required JavaScript. Tests prove it, 770 of them. This post explains what shipped and why.

## The problem

Most component libraries assume a client app. They want a React or Vue runtime, hydration, and a JavaScript bundle. If your app renders HTML from Go, Rails, Laravel, or Django, your options shrink. Gelium UI started from that question: can a server-rendered app get polished, accessible, themeable components without shipping JavaScript?

## What it is

Gelium UI is open-code. Components are semantic HTML partials you copy into your project, styled by CSS tokens. There is no client runtime, no hydration, and no CDN. It is built on Tailwind CSS 4 and installs from npm: `npm install gelium-ui`. HTMX is wired in for teams that want it, never as a requirement.

## The zero-JS contract

Zero required JavaScript is a hard contract, not a marketing claim. Tests scan the bundle and the markup to prove no component depends on a script. Native semantics come first: buttons are `<button>`, dialogs are real `<dialog>` elements, and states map to real attributes. The JavaScript that exists is progressive enhancement, about 5 KB and optional. The result is UI that survives scripts being blocked or slow.

## Themes are tokens

Themes are token values, not separate component styles. One class on `<html>` swaps Material 3 for Basecoat, Linear, Vercel, or Alden, with no rebuild and no markup changes. All eight Basecoat style packs ship as skins. The Basecoat translation is honest: the official package is audited, and every token was translated against it. Where no official evidence exists, like Base UI visuals, the theme says so.

## One component, many looks

The architecture separates behavior, visual reference, and skin. One semantic component dresses as Material, Basecoat, or Linear without forking markup. Density is a contract too: touch targets never drop below 44 pixels, even when the skin asks for tighter. Accessibility does not negotiate with aesthetics.

## Docs as a product

The docs are not an afterthought. Every page renders the real component it documents, on the same layout and the same tokens. The catalog covers 48 components, from buttons to data tables, with empty, error, and loading states. Roadmap checkmarks mean contract-tested, not aspirational.

## Who it is for

Gelium UI is for server-rendered apps: Go, Rails, Laravel, Django, or plain HTML. For teams that want HTMX or progressive enhancement. For design systems that must work with JavaScript disabled. It is not for React or Vue SPAs, or for Shadow DOM web components. Gelium sits closer to open-code HTML systems like GOV.UK Frontend than to a React design system. It does not replace them; it is another path.

## What is next

The project is pre-1.0 and moves fast: an icon pack gallery, docs localization, and more screen recipes. The code is MIT and lives on [GitHub](https://github.com/cryptonahue/gelium-ui). The package is on [npm](https://www.npmjs.com/package/gelium-ui). If you build server-rendered apps, take a look and tell us what you think.
