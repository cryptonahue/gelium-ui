# Banner

Banner is a page-level, full-width status bar that carries a persistent page or site signal — expired session, scheduled maintenance, consent, or a success that must survive the redirect. Use a banner when the whole page (or the whole site) needs one durable message that is not a transient toast and not a field error. It renders as a real `div` styled with the `ui-banner` class and needs no component JavaScript.

## Examples

<div class="component-preview banner-preview">
  <div class="ui-banner ui-banner--error" role="alert">
    <span class="ui-banner-icon" aria-hidden="true">!</span>
    <div class="ui-banner-content">
      <p class="ui-banner-title">Payment failed</p>
      <p class="ui-banner-body">We couldn't process your payment. Check the card details or try again.</p>
    </div>
    <a class="ui-button" href="/account/billing">Review payment method</a>
    <form class="ui-banner-dismiss" method="post" action="/dismiss/payment-failed">
      <button class="ui-icon-button" type="submit" aria-label="Dismiss">×</button>
    </form>
  </div>
  <div class="ui-banner ui-banner--info" role="status">
    <span class="ui-banner-icon" aria-hidden="true">i</span>
    <div class="ui-banner-content">
      <p class="ui-banner-title">Scheduled maintenance</p>
      <p class="ui-banner-body">Tonight 22:00–23:00 UTC. Some reports may be slow to load.</p>
    </div>
    <form class="ui-banner-dismiss" method="post" action="/dismiss/maintenance">
      <button class="ui-icon-button" type="submit" aria-label="Dismiss">×</button>
    </form>
  </div>
</div>

## Guidance

### When to use

Use a banner for a signal that belongs to the whole page or site rather than one element: maintenance windows, session expiry, consent, or a task that succeeded and landed somewhere that must show it. Match the tone to the situation with `ui-banner--{error,warning,success,info}` and let the `role` follow the tone (`alert` for error, `status` otherwise). The dismiss action is a real `POST` form that returns a `303`, so the banner closes without JavaScript.

### When not to use

Do not use a banner for field validation or for a single failed row: those belong to the [validation summary](/components/validation-summary), the [inline alert](/components/inline-alert) or a row-level message. Do not use a banner for an action result that disappears on its own — that is a [toast](/components/toast). The banner is persistent by design: if the signal should vanish after four seconds, it is the wrong surface.

### Usability

- A banner carries an optional icon, a title, a body and an optional call to action: the copy does the explaining, the tone only reinforces it.
- The body and title rest on the page-level tokens, so light and dark schemes stay consistent without extra work.
- When a task succeeds and the next page is the result, the banner in the main column confirms it instead of echoing a toast (see [Feedback](/docs/feedback), FEED-OK-PAGE).

### Accessibility

- The `role` derives from the tone: `error` becomes `alert` and everything else becomes `status`, so assistive technology announces the right kind of signal.
- The icon is decorative: keep it `aria-hidden` and put the meaning in the title or body text — the banner is never color-only.
- In forced-colors mode the banner keeps a `CanvasText` boundary and the error tone paints with `Mark`, so the signal survives without color.

## Anatomy

- **Banner** — `ui-banner`, the full-width flex row with the scoped `--ui-banner-*` tokens (`padding`, `gap`, `radius`, `icon-size`, `bg`, `fg`).
- **Icon** — `ui-banner-icon`, a decorative `aria-hidden` glyph for the tone's visual echo.
- **Content** — `ui-banner-content`, the column that stacks the title and body; it flexes to fill the row.
- **Title** — `ui-banner-title`, the short label line (`--ui-type-label-lg`).
- **Body** — `ui-banner-body`, the supporting sentence (`--ui-type-body-sm`).
- **Dismiss** — `ui-banner-dismiss`, a real `POST` form with a `ui-icon-button` submit so dismissal works with zero JavaScript.

## States

The banner has no hover or focus states of its own: it is a static status surface. Its only state change is dismissal, which is a server round trip (`POST` + `303`) that removes the banner on the next render.

## Checklist

- [ ] One durable, page-level message only — nothing that should expire on its own.
- [ ] Tone class applied (`--error`, `--warning`, `--success` or `--info`) and the role follows the tone.
- [ ] Meaning in text, not color: icon is `aria-hidden`, title/body carry the message.
- [ ] Dismiss (when present) is a real `POST` form with a labeled submit button.
- [ ] Copy says what happened and what to do next — no bare "Something went wrong".

## Accessibility

The banner is announced with the role that matches its tone and anything interactive inside it (CTA link, dismiss button) is a real control with its own accessible name. The dismiss button is labeled (default `aria-label="Dismiss"`). Color is never the sole channel for the message: text, icon and role all reinforce the tone.

## See also

- [Feedback](/docs/feedback) — the decision matrix this component lives in (FEED-PAGE, FEED-FAIL, FEED-OK-PAGE, FEED-SYS, FEED-PARTIAL).
- [Inline alert](/components/inline-alert) — the section-level sibling for smaller, contextual signals.
- [Toast](/components/toast) — the transient, action-result surface banners never replace.
- [Choose the right control](/docs/choose-the-right-control) — the cross-component decision.