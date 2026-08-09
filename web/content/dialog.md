# Dialog

Dialog is an open-code component built with the broadly available native `<dialog>` element. In supporting browsers, its declarative invoker commands need no component JavaScript or HTMX.

## Native modal behavior

The trigger uses `command="show-modal"` with `commandfor="confirm-dialog"`. These invoker attributes are recent Baseline Low features. Loom UI intentionally includes no component JavaScript fallback, so in older browsers the trigger does nothing; consumers supporting them need a server-rendered fallback or adapter. `request-close`, used by Cancel, is newer than the invoker commands; Confirm uses `command="close"` with the value `confirm`.

The visible headline and description are referenced by `aria-labelledby` and `aria-describedby`. `closedby` is not Baseline: `closedby="any"` adds light dismiss only in supporting browsers. The explicit Cancel action and native Escape behavior remain available in compatible browsers.

## Motion and accessibility

Opening and closing motion is progressive enhancement. `@starting-style` support varies, and the `overlay` close/top-layer transition is Chromium-only rather than interoperable, so motion may be instant or asymmetric elsewhere. Functional behavior in compatible browsers does not depend on animation. Under `prefers-reduced-motion: reduce`, dialog, content, and backdrop transitions are disabled. In forced colors, the dialog keeps a visible 2 px system-color boundary.
