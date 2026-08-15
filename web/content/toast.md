# Toast

Toast is a Gelium-only component for transient, server-driven feedback. Use a toast when an action completes server-side and the result should be announced in place. That means a save confirmation, an error, or a status change — without pushing users to another view.

## Guidance

### When to use

Use a toast when an action completes server-side and the result should be announced in place. That means a save confirmation, an error, or a status change — without pushing users to another view.

### When not to use

Never report field validation failures as toasts — show an inline error next to the field instead. Do not use toasts for persistent or critical feedback: that belongs in an inline alert or a banner that stays until the user acts.

### Usability

- Four variants share one anatomy (icon, message, optional Dismiss): `info`, `success`, `warning`, and `error`.
- The server triggers the toast through an `HX-Trigger` header that raises the `gelium:toast` event; without JavaScript the same POST re-renders a persistent inline toast.
- The timer pauses on hover and focus so the message remains readable.

### Accessibility

- `error` uses `role="alert"` (assertive); `info`, `success`, and `warning` use `role="status"` inside an `aria-live="polite"` region.
- The tone is conveyed by a decorative, aria-hidden icon and a semantic role that is never color alone.
- Dismiss is keyboard-accessible and focusable.
- `prefers-reduced-motion` disables the enter/exit transition; `forced-colors` replaces color-only signals with borders and system highlight colors.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Alternative names

- Snackbar, notification, toast notification.

## Agent prompt

Use Toast for transient, server-driven feedback: a save confirmation, an error, or a status change. Announce it in place without pushing users to another view. Four variants share one anatomy (icon, message, optional Dismiss action). `error` uses `role="alert"`; the rest use `role="status"` inside an `aria-live="polite"` region. The server triggers it through an `HX-Trigger` header, and the timer pauses on hover and focus. Never report validation failures as toasts, and never use it for persistent or critical feedback — that belongs in an Inline alert or a Banner.

## What variants can a toast show and how is it dismissed?

Four variants share one anatomy — an icon, a message, and an optional Dismiss action: `info`, `success`, `warning`, and `error`. The tone is conveyed by a decorative, aria-hidden icon and a semantic role that is never color alone. `error` uses `role="alert"` (assertive); `info`, `success`, and `warning` use `role="status"` (polite). Messages are announced by the surrounding `aria-live="polite"` live region.

## Server-driven wire contract

The demo below posts to `/examples/toast/demo`. With HTMX, the response carries an `HX-Trigger` header that raises the `gelium:toast` event.

```json
{"gelium:toast":{"type":"success","message":"Record updated"}}
```

The local `app.js` listens for that event and displays an auto-dismissing toast in a fixed bottom region. The timer pauses on hover and focus so the message stays readable.

## No-JavaScript fallback

Without JavaScript the same POST re-renders the full page with a persistent inline toast. The server-driven feedback is complete before any enhancement loads. Validation failures are never reported as toast notifications.

## Styling

The surface uses dedicated `--ui-toast-*` tokens and `--ui-shadow-3` elevation in both light and dark schemes. Dismiss is keyboard-accessible and focusable, `prefers-reduced-motion` disables the enter/exit transition, and `forced-colors` replaces color-only signals with borders and system highlight colors.
