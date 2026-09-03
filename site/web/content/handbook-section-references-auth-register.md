# Auth / register

Read · form. Registration should make identity, requirements, and recovery visible before the user commits.

**[See Gelium remake](#gelium-remake)**

[Open original](https://vercel.com/signup) (Vercel, verified 30 Aug 2026; HTTP 200)

[Section references](/docs/section-references)

## Original (cited)

Live URL: [vercel.com/signup](https://vercel.com/signup)

The official page is a registration surface with social alternatives, account fields, an agreement notice, and proof that the product is used in production. This block map records the observed structure, not Vercel's logos, assets, or copy.

```text
Vercel wordmark
H1  Your first deploy is just a sign-up away.
social sign-in options (Google · GitHub · OpenAI · Apple)
email + password registration form
Terms of Service + Privacy Policy agreement
customer proof / release-cycle and uptime metrics
```

The page presents the same registration message in responsive variants. The reusable jobs are identity entry, explicit consent, and a credible path into the product.

## Gelium remake

This is how **that registration structure** lands in Gelium: a native form, explicit labels, a terms checkbox, one primary submit, and a recovery link. The action is a product integration seam, not a claim that Gelium ships an account service here.

<section id="gelium-remake" aria-labelledby="auth-register-title">
  <header>
    <p>Account start</p>
    <h2 id="auth-register-title">Create a workspace account</h2>
    <p>Start with the identity you will use to deploy, review, and recover access.</p>
  </header>
  <div class="ui-card ui-card-outlined">
    <form class="auth-register-form" method="post" action="/account/register">
    <div class="ui-inline-alert ui-inline-alert--error" role="alert">
      <p class="ui-inline-alert-title">Example validation state</p>
      <p>Resolve the highlighted fields before creating the account.</p>
    </div>
    <div class="ui-validation-summary" role="alert">
      <h3 class="ui-validation-summary-title">Check the following</h3>
      <ul class="ui-validation-summary-list">
        <li class="ui-validation-summary-item"><a href="#auth-email">Enter a valid email address.</a></li>
        <li class="ui-validation-summary-item"><a href="#auth-password">Use at least 8 characters for your password.</a></li>
      </ul>
    </div>
    <div class="ui-text-field ui-text-field-filled">
      <span class="ui-text-field-control">
        <label for="auth-email">Email address</label>
        <input id="auth-email" name="email" type="email" autocomplete="email" required aria-describedby="auth-email-help">
      </span>
      <p class="ui-text-field-message" id="auth-email-help">Use an address you can access for recovery.</p>
    </div>
    <div class="ui-text-field ui-text-field-filled">
      <span class="ui-text-field-control">
        <label for="auth-password">Password</label>
        <input id="auth-password" name="password" type="password" autocomplete="new-password" minlength="8" required>
      </span>
    </div>
    <p><label class="ui-checkbox"><input id="auth-terms" name="terms" value="accepted" type="checkbox" required><span class="ui-checkbox-mark"></span><span class="ui-checkbox-label">I agree to the Terms and Privacy Policy.</span></label></p>
    <p><button class="ui-button ui-button-primary" type="submit">Create account</button></p>
    <p><a href="/account/recover">Forgot your password?</a></p>
  </form>
  </div>
</section>

A real registration handler should return the same form with HTTP 422, preserved values where safe, `aria-invalid="true"`, and `X-Gelium-Validation: true` for an optional HTMX fragment swap. It should never repopulate a password value.

## Keep / adapt

| Kept | Gelium primitive | Rejected | Why |
|---|---|---|---|
| Clear registration entry point | native form + labeled `ui-text-field` inputs | Copied wordmark or social logos | Identity is the job; another product's brand is not |
| Email and password sequence | native `input` with `autocomplete` and constraints | Unlabeled floating controls | Labels and browser semantics must survive no-JS |
| Explicit agreement | `ui-checkbox` inside the form | Consent hidden in small print | The user must understand what the submit includes |
| One commit action | `ui-button ui-button-primary` | Multiple competing primary buttons | Registration needs one next move |
| Field-level recovery | `ui-validation-summary` links + recovery anchor | Silent client-only errors | 422 feedback must be persistent and navigable |

## Ask before copying

- Which identity providers are actually approved, and what account-linking policy protects users from duplicate identities?
- What password, passkey, rate-limit, bot-defense, and email-verification requirements apply before this form is connected?
- Where should failed registration return the user, and which fields are safe to preserve in a 422 response?
