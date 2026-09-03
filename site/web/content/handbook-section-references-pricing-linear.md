# Pricing / plan comparison

Persuade · detail. A pricing page should make plan differences scannable, state the billing basis, and offer one next step without turning every cell into a sales pitch.

**[See Gelium remake](#gelium-remake)**

[Open original](https://linear.app/pricing) (Linear, verified 30 Aug 2026; HTTP 200)

[Section references](/docs/section-references)

## Original (cited)

Live URL: [linear.app/pricing](https://linear.app/pricing)

The Linear pricing page presents four plans, a feature comparison, customer proof, and a closing availability section. The plan controls and sales path repeat in the source for responsive placements; this block map records the information architecture without importing Linear's logo, artwork, or CSS.

```text
skip link + product navigation (chrome, out of scope)
H1  Pricing
plans in order  Free · Basic · Business · Enterprise
plan price / billing basis + included feature bullets
plan-level get-started or contact-sales links
customer proof and customer-stories link
feature comparison table across the four plans
closing availability statement + get-started / sales / app links
```

The core pricing job is comparison first: people can see the plan ladder and feature boundaries before choosing whether they need self-serve or a sales conversation.

## Gelium remake

This is how **that pricing structure** lands in Gelium: a native, narrow-safe data table handles cross-plan comparison, while one primary link provides the declared next move. Prices and feature claims below are illustrative Gelium copy, not Linear's commercial terms.

<section id="gelium-remake" aria-labelledby="pricing-linear-title">
  <header>
    <p>Plans</p>
    <h2 id="pricing-linear-title">Choose the amount of coordination you need</h2>
    <p>Compare the work, controls, and support that each plan makes available.</p>
  </header>
  <div class="pricing-linear-plans">
    <article class="ui-card ui-card-outlined">
      <h3 class="ui-card-title">Free</h3>
      <p class="ui-card-body">$0 / month · A small team trying the workflow.</p>
      <ul class="ui-list"><li class="ui-list-item">Core workspace</li><li class="ui-list-item">Limited collaboration</li></ul>
    </article>
    <article class="ui-card ui-card-outlined">
      <h3 class="ui-card-title">Basic</h3>
      <p class="ui-card-body">$12 / member / month · Shared delivery.</p>
      <ul class="ui-list"><li class="ui-list-item">Unlimited projects</li><li class="ui-list-item">Team reporting</li></ul>
    </article>
    <article class="ui-card ui-card-outlined">
      <h3 class="ui-card-title">Business</h3>
      <p class="ui-card-body">$24 / member / month · Cross-team coordination.</p>
      <ul class="ui-list"><li class="ui-list-item">Roles and approvals</li><li class="ui-list-item">Usage planning</li></ul>
    </article>
    <article class="ui-card ui-card-outlined">
      <h3 class="ui-card-title">Enterprise</h3>
      <p class="ui-card-body">Custom · Governed organizations.</p>
      <ul class="ui-list"><li class="ui-list-item">Advanced controls</li><li class="ui-list-item">Priority support</li></ul>
    </article>
  </div>
  <div class="ui-data-table-scroll">
    <table class="ui-data-table-table">
      <caption>Plan comparison</caption>
      <thead>
        <tr><th scope="col">Plan</th><th scope="col">Price</th><th scope="col">Best for</th><th scope="col">Includes</th></tr>
      </thead>
      <tbody>
        <tr><th scope="row">Free</th><td>$0 / month</td><td>Trying the workflow</td><td>Core workspace and limited collaboration</td></tr>
        <tr><th scope="row">Team</th><td>$12 / member / month</td><td>Shared delivery</td><td>Unlimited projects, roles, and team reporting</td></tr>
        <tr><th scope="row">Scale</th><td>Custom</td><td>Governed organizations</td><td>Advanced controls, support, and usage planning</td></tr>
      </tbody>
    </table>
  </div>
  <p><a class="ui-button ui-button-primary" href="/docs/section-references">Compare your starting point</a></p>
</section>

The table remains a real `<table>` with a caption and scoped headers. On narrow screens its wrapper provides internal scrolling rather than masking content or shrinking the comparison into unreadable cards.

## Keep / adapt

| Kept | Gelium primitive | Rejected | Why |
|---|---|---|---|
| Plan ladder before detail | native `ui-data-table` | Marketing tiles with hidden differences | Comparison needs aligned columns and headers |
| Billing basis beside plan name | table caption, headers, and cells | Invented “from” prices or feature claims | Commercial facts need an owner and source |
| One clear next step | `ui-button ui-button-primary` link | CTA in every feature cell | Repetition competes with the comparison job |
| Responsive comparison | `ui-data-table-scroll` with `min-width: 0` | `overflow-x: hidden` or clipped columns | Users must be able to reach every value |
| Sales boundary for complex needs | plain explanatory copy | Fake checkout or account workflow | A reference page must not imply unavailable commerce behavior |

## Ask before copying

- Who owns plan names, prices, limits, entitlements, and the date on which each claim was last verified?
- Which features need a contract, security review, or sales handoff instead of self-serve activation?
- Is comparison best served by a table for this audience, or would a smaller plan set be clearer as semantic cards?
