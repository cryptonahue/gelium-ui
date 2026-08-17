# Error state

Error state is a server-rendered, full-page status for a resource that could not be delivered — a 404 page, a 500 failure, a section that failed to load. Use an error state when the whole surface cannot render, so the user knows what happened and how to recover. It renders as a real `div` styled with the `ui-error-state` class with `role="alert"` and needs no component JavaScript.

## Examples

<div class="component-preview">
  <div class="ui-error-state" role="alert">
    <p class="ui-error-state-code" aria-hidden="true">404</p>
    <h1 class="ui-error-state-title">Page not found</h1>
    <p class="ui-error-state-body">The page you are looking for does not exist or has moved.</p>
    <a class="ui-button" href="/">Back to home</a>
  </div>
</div>

<div class="component-preview">
  <div class="ui-error-state" role="alert">
    <p class="ui-error-state-code" aria-hidden="true">500</p>
    <h1 class="ui-error-state-title">Something went wrong</h1>
    <p class="ui-error-state-body">We couldn't load this workspace. Try again, or come back in a few minutes.</p>
    <a class="ui-button" href="/workspaces/acme">Try again</a>
  </div>
</div>

## Guidance

### When to use

Use an error state when a resource cannot be delivered: a page that does not exist (404), an internal failure (500), a section whose fetch failed (see [Feedback](/docs/feedback), FEED-LOAD-FAIL and FEED-FAIL). The server picks the real HTTP status and the copy per status; the retry is an optional real GET link back to a known URL.

### When not to use

Do not use an error state for an empty collection — that is an [empty state](/components/empty-state). Do not use an error state for a field or section-level failure that the rest of the page can survive: an [inline alert](/components/inline-alert) or a row-level message fits better. Never show an endless [skeleton](/components/skeleton) instead of the error state: when a fetch fails, fail fast with a real status and a retry.

### Usability

- The oversized status code is decorative — the heading carries the meaning.
- Copy says what happened and what to do next — never a bare "Error 500" with no recovery.
- The optional retry is a real link, not a button masquerading as one.

### Accessibility

- The root is `role="alert"`, so assistive technology announces the failure as an actionable condition.
- The status code is `aria-hidden` decoration; the heading, body and retry carry the meaning.
- Under forced colors the code, title and body repaint `CanvasText` and the retry stays a link (`LinkText`).

## Anatomy

- **Error state** — `ui-error-state`, the centered flex column with the scoped `--ui-error-state-*` tokens (`padding`, `gap`, `code-color`, `title-color`, `body-color`).
- **Code** — `ui-error-state-code`, the oversized decorative status (`--ui-type-display-lg`, `aria-hidden`).
- **Title** — `ui-error-state-title`, the failure headline (`--ui-type-title-lg`).
- **Body** — `ui-error-state-body`, the explanation and next step (`--ui-type-body-sm`).
- **Retry** — an optional `ui-button` link to a known recovery URL.

## States

The error state has no hover or focus states of its own; the only interactive element is the optional retry link, which follows the Button contract. The state is defined by the real HTTP status the server renders (404, 500, 503) and the copy paired with it.

## Checklist

- [ ] Real HTTP status rendered (404/500/503) — not a made-up code.
- [ ] `role="alert"` on the root; status code `aria-hidden`.
- [ ] Copy says what happened and what to do next.
- [ ] Retry, when present, is a real link to a known URL.
- [ ] The failure surfaced fast — no eternal skeleton, no silent blank.
- [ ] Forced-colors output keeps heading and retry readable.

## Accessibility

The error state announces the failure with `role="alert"` and keeps every piece of meaning in text: the heading names the failure and the body names the recovery. The decorative code stays hidden from assistive technology, and the retry is a real link so keyboard and assistive users can act on it. Color is never the only signal.

## See also

- [Feedback](/docs/feedback) — the decision matrix entries FEED-FAIL and FEED-LOAD-FAIL.
- [Empty state](/components/empty-state) — the sibling surface for a genuinely empty (not failed) result.
- [Inline alert](/components/inline-alert) — the section-level failure signal when the page can continue.
- [Skeleton](/components/skeleton) — the loading surface whose failure path is the error state.
- [Choose the right control](/docs/choose-the-right-control) — the cross-component decision.