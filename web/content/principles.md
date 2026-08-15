# Design principles

Four principles decide every Gelium component. They come from the foundation (see [Documentation](/docs)) and are enforced by tests, not by taste.

## Native semantics first

**What it means.** Every component starts from the platform element that already provides the behavior: `<button>` for actions, `<a>` for navigation, `<dialog>` for dialogs, native form controls for input. ARIA is used only where the platform has no equivalent.

**Why.** Native elements ship behavior, keyboard support, and accessibility for free — and keep working when JavaScript never loads.

**Example.** [Button](/components/button) is a real `<button>` (or a real `<a>` for navigation); a loading button is a disabled native button with `aria-busy="true"`, not a styled div.

## No component JavaScript

**What it means.** Components ship no JavaScript of their own. Selection, menus, tooltips, and dialogs run on declarative platform features (`:checked`, the Popover API, Invoker Commands, Interest Invokers) or plain server round-trips. HTMX is an optional enhancement, never a requirement.

**Why.** No bundles to build, no runtime to version, nothing to break — the same markup works in any server-rendered application.

**Example.** [Switch](/components/switch) and [Checkbox](/components/checkbox) run on `:checked`; [Menu](/components/menu) uses the Popover API; [Dialog](/components/dialog) uses native `<dialog>` with `closedby`.

## Server-rendered by design

**What it means.** Every component is a server-rendered HTML partial with a documented [server contract](/docs/server-contracts). The URL is the state, mutations are POST + 303, validation is 422, and transient feedback is server-driven.

**Why.** The code is open: copy the partial and its contract into your own templates — nothing is generated, nothing is hidden behind a runtime.

**Example.** This documentation site is the first Gelium application, and [Data table](/components/data-table) keeps list state in stable GET parameters that work with and without HTMX.

## Accessible by default

**What it means.** Light/dark, forced colors, reduced motion, keyboard focus, and RTL are part of every component contract — implemented in the core, not patched per component.

**Why.** Accessibility decided at design time costs less than accessibility retrofitted after a release, and every page ships accessible from the first render.

**Example.** [Toast](/components/toast) renders inside an `aria-live="polite"` region and uses `role="alert"` for errors; [Icon button](/components/icon-button) pairs decorative SVG with a real accessible name.

## How we verify

The build enforces the principles: `go test ./...` asserts server contracts, `aria-*` attributes, live regions, and 422 recovery; the style contract tests assert token-driven states plus reduced-motion and forced-colors coverage; the theme matrix tests keep themes swappable without touching markup. A change that weakens a principle needs a stronger reason.
