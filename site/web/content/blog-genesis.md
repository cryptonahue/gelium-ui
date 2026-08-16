# How Gelium UI was born

Gelium UI began as a side project with a simple question. Could server-rendered Go apps get a component library without shipping JavaScript? Most libraries assumed a client app. Gelium chose the opposite path: HTML partials and tokens, rendered by the server, with zero required JS. The name says it: the library is UI built around the Gelium model, not a framework to learn.

## Server-rendered first

Every component is a Go template partial and a CSS file. The server renders the markup, so pages work before any script runs. Progressive enhancement stays optional: plain forms and links carry the interaction. HTMX layers on when a team wants it, never as a requirement.

## The 0-JS contract

Gelium made zero required JS a hard contract, not a marketing claim. Tests scan the bundle and the markup to prove no component depends on a script. Native semantics come first: buttons are buttons, dialogs are real dialogs, and states map to real attributes. The result is UI that survives scripts being blocked or slow.

## Token-driven themes

Themes are token values, not separate component styles. A `--ui-*` vocabulary feeds every rule in every theme. Material ships as the default, and Basecoat proves the model: a second theme added with no markup changes. Switching themes never touches component CSS, because components never hardcode color.

## Built by contract, phase by phase

The roadmap ran as ten phases, from A to J. Every phase shipped with contract tests that pin its behavior. If a phase regresses, the build fails. This is why the docs can say "contract-tested" without hedging: the tests are the receipts.

## The rename

Gelium UI did not start with its name. Early builds used a placeholder brand, Loom UI, and the wire contracts carried that legacy prefix. Renaming a pre-release project is cheap, so the owner migrated the prefix to match the product. The wire contract now says `gelium:*`, and the docs never look back.

## Docs as a product

The docs are not an afterthought. Every page is dogfooded: it renders the real component it documents. The Handbook explains concepts before the reference catalog lists components. A style guide sets the voice, and tests enforce plain English with sentences under 25 words. The blog you are reading runs on the same layout, the same tokens, and the same zero-JS philosophy.
