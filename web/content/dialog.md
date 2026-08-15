# Dialog

Dialog is an open-code component for confirmations. Use a dialog when a decision needs explicit confirmation before an action completes — and the flow must work without JavaScript. Its base path is a **page variant**: the trigger is a real link to a server-rendered confirmation page, Confirm is a real form POST that redirects back, and Cancel is a link back. The flow works in every browser with **no component JavaScript** and no overlay markup.

## Page variant (base)

The trigger is an anchor styled as a button:

```html
<a class="ui-button ui-button-primary" href="/components/dialog/confirm">Open confirmation dialog</a>
```

`/components/dialog/confirm` renders the same headline and description inline as normal page content. Confirm submits a real form POST answered with a `303 See Other` back to the docs page, which shows the result in a persistent inline alert; Cancel is a link back. This is the "page/detail variant" of the Dialog contract: no overlay, no focus trap, no Escape contract — navigation and form submission are platform behavior that works in every browser, including browsers without Invoker Commands support (Invoker Commands is the native `command`/`commandfor` API that opens and closes dialogs declaratively).

## How does the modal variant handle dismiss and Escape?

Consumers that target **supporting browsers** can use the native `<dialog>` modal with declarative invoker commands instead:

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

`command`/`commandfor` are **Baseline 2025 — Newly available** (December 2025), not a Chromium-only feature, so the modal is a legitimate enhancement for current browsers; it is deliberately **opt-in** and never a page's only path. `request-close`, used by Cancel, is newer than the invoker commands. `closedby` is **not Baseline**: `closedby="any"` adds light dismiss only in supporting browsers (Safari ignores it and keeps Escape/cancel). The explicit Cancel action and native Escape behavior remain available in compatible browsers.

## Motion and accessibility

The page variant has no overlay motion. The modal's opening and closing motion is progressive enhancement: `@starting-style` support varies, and the `overlay` close/top-layer transition is Chromium-only rather than interoperable, so modal motion may be instant or asymmetric elsewhere. Functional behavior in compatible browsers does not depend on animation. Under `prefers-reduced-motion: reduce`, dialog, content, and backdrop transitions are disabled. In forced colors, the dialog keeps a visible 2 px system-color boundary.
