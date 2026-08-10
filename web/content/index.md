# Gelidium UI

Themeable, open-code UI components for server-rendered applications. Native HTML semantics, zero component JavaScript, and a Material 3 design system built with Tailwind CSS v4 and HTMX.

## The foundation

Every Gelidium component starts from the platform, not from a framework:

- **Native semantics first** — real `<button>`, `<dialog>`, `<input type="radio">`, `<table>`, `<nav>`; ARIA only where the platform has no equivalent.
- **No component JavaScript** — selection, menus, tooltips, and dialogs run on declarative platform features (`:checked`, Popover API, Invoker Commands, Interest Invokers) or plain server round-trips.
- **Server-rendered by design** — HTMX as an enhancement, never a requirement. The docs you are reading are themselves the first Gelidium application.
- **Accessible by default** — light/dark, forced-colors, reduced-motion, keyboard focus, and RTL are part of every component contract, tested in the build.

## Get started

Read the [documentation](/docs) or jump straight into a component: [Button](/components/button), [Text field](/components/text-field), [Dialog](/components/dialog), or [Toast](/components/toast).

See the full library in action with the [WhatsApp manager demo](/demo/whatsapp).
