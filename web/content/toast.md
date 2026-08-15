# Toast

Toast is a Gelium-only component for transient, server-driven feedback. Use a toast when an action completes server-side and the result should be announced in place — a save confirmation, an error, a status change — without pushing users to another view.

## What variants can a toast show and how is it dismissed?

Four variants share one anatomy — an icon, a message, and an optional Dismiss action: `info`, `success`, `warning`, and `error`. The tone is conveyed by a decorative, aria-hidden icon and a semantic role that is never color alone: `error` uses `role="alert"` (assertive), while `info`, `success`, and `warning` use `role="status"` (polite). Messages are announced by the surrounding `aria-live="polite"` live region.

## Server-driven wire contract

The demo below posts to `/examples/toast/demo`. With HTMX, the response carries an `HX-Trigger` header that raises the `loom:toast` event:

```json
{"loom:toast":{"type":"success","message":"Record updated"}}
```

The local `app.js` listens for that event and displays an auto-dismissing toast in a fixed bottom region, pausing its timer on hover and focus so the message remains readable.

## No-JavaScript fallback

Without JavaScript the same POST re-renders the full page with a persistent inline toast, so the server-driven feedback is complete before any enhancement loads. Validation failures are never reported as toast notifications.

## Styling

The surface uses dedicated `--ui-toast-*` tokens and `--ui-shadow-3` elevation in both light and dark schemes. Dismiss is keyboard-accessible and focusable, `prefers-reduced-motion` disables the enter/exit transition, and `forced-colors` replaces color-only signals with borders and system highlight colors.
