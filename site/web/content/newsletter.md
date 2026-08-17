# Newsletter

Newsletter is a zero-JavaScript subscription form for audience conversion. Use it when a page invites people to subscribe to a newsletter and the whole contract must work without client-side scripting — the server owns the submit, the validation, and the success view. The form is an `<aside>` styled with the `ui-newsletter` classes; HTMX is an optional enhancement that swaps only the aside fragment.

## Examples

The initial subscribe state, ready for a no-JS POST.

<div class="specimen-block">
<aside class="ui-newsletter" aria-labelledby="newsletter-title">
  <h2 id="newsletter-title" class="ui-newsletter-title">Stay in the loop</h2>
  <p class="ui-newsletter-description">Product updates and Gelium UI releases. No spam, unsubscribe anytime.</p>
  <form class="ui-newsletter-form" method="post" action="/examples/newsletter" hx-post="/examples/newsletter" hx-target="closest .ui-newsletter" hx-swap="outerHTML">
    <div class="ui-newsletter-field">
      <label class="ui-newsletter-label" for="newsletter-email">Email address</label>
      <div class="ui-newsletter-row">
        <input class="ui-newsletter-input" id="newsletter-email" name="email" type="email" autocomplete="email" required>
        <button class="ui-button ui-button-primary" type="submit"><span>Subscribe</span></button>
      </div>
    </div>
  </form>
</aside>
</div>

An invalid email re-renders the aside with an inline error, the submitted value preserved, and the input flagged.

<div class="specimen-block">
<aside class="ui-newsletter" aria-labelledby="newsletter-title">
  <h2 id="newsletter-title" class="ui-newsletter-title">Stay in the loop</h2>
  <p class="ui-newsletter-description">Product updates and Gelium UI releases. No spam, unsubscribe anytime.</p>
  <form class="ui-newsletter-form" method="post" action="/examples/newsletter" hx-post="/examples/newsletter" hx-target="closest .ui-newsletter" hx-swap="outerHTML">
    <div class="ui-newsletter-alert" id="newsletter-error">
      <div class="ui-inline-alert ui-inline-alert--error" role="alert">
        <p class="ui-inline-alert-title">Invalid email address</p>
        <p class="ui-inline-alert-body">Enter a valid email address to subscribe.</p>
      </div>
    </div>
    <div class="ui-newsletter-field">
      <label class="ui-newsletter-label" for="newsletter-email">Email address</label>
      <div class="ui-newsletter-row">
        <input class="ui-newsletter-input" id="newsletter-email" name="email" type="email" autocomplete="email" required aria-invalid="true" aria-describedby="newsletter-error" value="not-an-email">
        <button class="ui-button ui-button-primary" type="submit"><span>Subscribe</span></button>
      </div>
    </div>
  </form>
</aside>
</div>

A valid address replaces the form entirely with a persistent confirmation.

<div class="specimen-block">
<aside class="ui-newsletter" aria-labelledby="newsletter-title">
  <h2 id="newsletter-title" class="ui-newsletter-title">Stay in the loop</h2>
  <p class="ui-newsletter-description">Product updates and Gelium UI releases. No spam, unsubscribe anytime.</p>
  <p class="ui-newsletter-success" role="status">You're subscribed — watch your inbox for a confirmation email.</p>
</aside>
</div>

The specimens above follow the contract the template `newsletter.html` defines and the data the server example in `internal/app/newsletter.go` renders.

## Guidance

### When to use

Use the newsletter when a page invites subscription and the signup must work without JavaScript — a footer aside, a blog sidebar, a landing block. It earns its place when the server can own the whole round-trip: validate, re-render, and confirm.

### When not to use

Do not use it for anything beyond a single email subscription: it has one field and one action. Do not rely on client-side validation alone — the 422 contract is the real gate. If the surface is a general contact form, use the [Text field](/components/text-field) and form patterns instead.

### Usability

- The form is a real `<form method="post">` with the action URL, so it submits without any script.
- HTMX is optional and additive: the same POST targets the same URL and swaps only the `closest .ui-newsletter` fragment.
- The email input is `type="email"` with `autocomplete="email"` and `required`, so the browser and the server validate together.

### Accessibility

- The aside is labelled by its title via `aria-labelledby`.
- The error reuses the inline-alert primitive with `role="alert"`, and the input carries `aria-invalid` plus `aria-describedby` pointing at it.
- The success view is a `role="status"` paragraph that replaces the form, so the confirmation is announced politely and never vanishes.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.

## Anatomy

- **`ui-newsletter`** — the aside: the surface-container background, the theme radius, the border, and a flex column with `--ui-newsletter-*` tokens.
- **`ui-newsletter-title` / `ui-newsletter-description`** — the heading and the one-line pitch, using the headline-sm and body-sm typescales.
- **`ui-newsletter-form`** — the POST form with the email field row and the submit button.
- **`ui-newsletter-alert`** — the wrapper for the inline error, rendered only when validation fails.
- **`ui-newsletter-field` / `ui-newsletter-row` / `ui-newsletter-label` / `ui-newsletter-input`** — the field structure: label, flexible input (min 14rem), and the primary submit button in one row.
- **`ui-newsletter-success`** — the persistent confirmation paragraph (`role="status"`) that replaces the form.

## Server contract

- **POST + 422** — an invalid email is rejected with status 422 and the `X-Gelium-Validation: true` header, re-rendering the aside with the inline error and the submitted value preserved.
- **POST → 200 success** — a valid email replaces the form with the persistent confirmation.
- In production the success landing is a POST + 303 redirect to a GET page that re-renders the success view; the example route demonstrates the direct 200 form.

## Anti-patterns

- Do not auto-submit on keystroke: the submit button is always visible and the round-trip is explicit.
- Do not hide the email field behind a consent checkbox or extra steps — the pattern is one field, one action.
- Do not dismiss the success view after a timeout: it is persistent until the person navigates.

## Checklist

- The aside is labelled and the form posts to a real action URL.
- The email input is `type="email"`, `autocomplete="email"`, and `required`.
- The 422 path preserves the submitted value and flags the input with `aria-invalid` and `aria-describedby`.
- The success view replaces the form and carries `role="status"`.
- Nothing requires JavaScript to subscribe.

## Accessibility

The whole flow works with no script: the native form submits, the server re-renders the aside with an alert or a confirmation, and every state is announced through the right role (`alert` for the error, `status` for the success). The input keeps its label, its autocomplete hint, and its error association, so assistive technology users get the same path as everyone else.
