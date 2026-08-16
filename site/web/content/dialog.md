# Dialog

Dialog is an open-code component for confirmations. Use a dialog when a decision needs explicit confirmation before an action completes — and the flow must work without JavaScript. Its base path is a **page variant**: the trigger is a real link to a server-rendered confirmation page. Confirm is a real form POST that redirects back, and Cancel is a link back. The flow works in every browser with **no component JavaScript** and no overlay markup.

## Guidance

### When to use

Use a dialog when a decision needs explicit confirmation before an action completes. The flow must work without JavaScript. The base path is a server-rendered page variant; supporting browsers can opt into the native `<dialog>` modal.

### When not to use

Do not use a dialog for a long or deep flow — that belongs on a page or in steps. Never make the overlay the only path: the no-JS page variant must always exist. For transient feedback after an action completes, use a [Toast](/components/toast); for persistent inline messages, use an inline alert or a [Card](/components/card).

### Usability

- The trigger is a real link styled as a button. Confirm is a real form POST answered with a `303 See Other`. Cancel is a link back.
- The modal is an opt-in enhancement via native invoker commands (`command`/`commandfor`) for supporting browsers.
- Keep the headline and description short — a dialog is for a decision, not a tutorial.

### Accessibility

- The page variant needs no focus trap, overlay or Escape contract — navigation and form submission are platform behavior.
- The modal is a native `<dialog>`: the browser manages the top layer, focus, and Escape; wire `aria-labelledby` and `aria-describedby`.
- Under `prefers-reduced-motion: reduce`, dialog, content, and backdrop transitions are disabled.
- In forced colors, the dialog keeps a visible 2 px system-color boundary.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Alternative names

- Modal, modal dialog, confirmation dialog, alert dialog.

## Agent prompt

Use Dialog for short, focused, reversible decisions that need explicit confirmation. The flow must work without JavaScript. The base path is a page variant: the trigger is a real link to a server-rendered confirmation page. Confirm is a real form POST answered with a `303 See Other`, and Cancel is a link back. Supporting browsers can opt into the native `<dialog>` modal through declarative invoker commands (`command`/`commandfor`). Never use a dialog for a long or deep flow — that belongs on a Page or in Steps. Never make the overlay the only path.

## Page variant (base)

The trigger is an anchor styled as a button:

```html
<a class="ui-button ui-button-primary" href="/components/dialog/confirm">Open confirmation dialog</a>
```

`/components/dialog/confirm` renders the same headline and description inline as normal page content. Confirm submits a real form POST answered with a `303 See Other` back to the docs page. The result shows in a persistent inline alert. Cancel is a link back. This is the "page/detail variant" of the Dialog contract: no overlay, no focus trap, no Escape contract. Navigation and form submission are platform behavior that works in every browser, including browsers without Invoker Commands support. Invoker Commands is the native `command`/`commandfor` API that opens and closes dialogs declaratively.

## How does the modal variant handle dismiss and Escape?

Consumers that target **supporting browsers** can use the native `<dialog>` modal with declarative invoker commands instead.

```html
<button type="button" command="show-modal" commandfor="confirm-dialog">Open confirmation dialog</button>
<dialog id="confirm-dialog" closedby="any" aria-labelledby="confirm-dialog-title" aria-describedby="confirm-dialog-description">
  <h2 id="confirm-dialog-title">Confirm action</h2>
  <p id="confirm-dialog-description">This action will apply the selected changes. Do you want to continue?</p>
  <div class="ui-dialog-actions">
    <button type="button" command="request-close" autofocus>Cancel</button>
    <button type="button" command="close" value="confirm">Confirm</button>
  </div>
</dialog>
```

`command`/`commandfor` are **Baseline 2025 — Newly available** (December 2025), not a Chromium-only feature. The modal is a legitimate enhancement for current browsers. It is deliberately **opt-in** and never a page's only path. `request-close`, used by Cancel, is newer than the invoker commands. `closedby` is **not Baseline**: `closedby="any"` adds light dismiss only in supporting browsers (Safari ignores it and keeps Escape/cancel). The explicit Cancel action and native Escape behavior remain available in compatible browsers.

## Motion and accessibility

The page variant has no overlay motion. The modal's opening and closing motion is progressive enhancement: `@starting-style` support varies, and the `overlay` close/top-layer transition is Chromium-only. Modal motion may be instant or asymmetric elsewhere. Functional behavior in compatible browsers does not depend on animation. Under `prefers-reduced-motion: reduce`, dialog, content, and backdrop transitions are disabled. In forced colors, the dialog keeps a visible 2 px system-color boundary.
