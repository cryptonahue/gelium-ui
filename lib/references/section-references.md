# Section reference fichas

Each ficha below is a portable structural summary. The docs site may contain a richer on-page Gelium remake; this package copy remains usable without that site.

## REF-ARTICLE — `section-references/article.md`

Source: `https://vercel.com/blog/the-end-of-credential-sprawl-for-agents`

Structure: article chrome → one H1 → author/date/read-time → lead → hero media → H2/H3 prose → code/evidence → related content.

Gelium adaptation: `SCREEN=detail`, readable prose measure, native media with alt/captions/transcript when applicable, visible source/provenance, loading/empty/error recovery. Do not copy Vercel skin, logo, assets, or claims.

## REF-404 — `section-references/404.md`

Source: `https://vercel.com/docs/errors/does-not-exist`

Structure: clear not-found identity → short explanation → one path back → optional secondary help.

Gelium adaptation: `ui-error-state`, one H1, canonical recovery link, no screenshot-only error, no unrelated recipe.

## REF-AUTH — `section-references/auth.md`

Source: `https://vercel.com/signup`

Structure: auth identity → labeled fields → consent/context → one primary submit → validation → recovery or alternate auth path.

Gelium adaptation: native labeled controls, server-side 422, validation summary plus inline errors, safe `next`, no nested card for a simple form, no claim that an unregistered route exists.

## REF-FAQ — `section-references/faq.md`

Source: `https://vercel.com/pricing`

Structure: decision context → grouped questions → native disclosure → concise answers → one next action.

Gelium adaptation: native `details`/`summary`, answers owned by the product, current commercial claims only, one decision action, keyboard and no-JS operation.

## REF-HERO — `section-references/hero.md`

Source: `https://linear.app/`

Structure: promise → short context → one primary CTA → optional supporting proof.

Gelium adaptation: `ui-hero`, page-owned H1, tokenized skin, no copied logo/gradient/motion, no extra CTA competition.

## REF-PRICING — `section-references/pricing.md`

Source: `https://linear.app/pricing`

Structure: pricing context → plan ladder → feature comparison → billing basis → FAQ/qualification → one primary CTA.

Gelium adaptation: mobile-safe `ui-data-table` or semantic cards, explicit billing basis, real product claims only, no invented prices, no layout table for non-comparison content.
