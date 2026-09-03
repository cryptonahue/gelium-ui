# 404 / missing route

Read · detail. A not-found page should name the failure plainly and give the reader one known way forward.

**[See Gelium remake](#gelium-remake)**

[Open original](https://vercel.com/docs/errors/does-not-exist) (Vercel, verified 30 Aug 2026; HTTP 404)

[Section references](/docs/section-references)

## Original (cited)

Live URL: [vercel.com/docs/errors/does-not-exist](https://vercel.com/docs/errors/does-not-exist)

The official route responds with HTTP 404. Its visible not-found content is intentionally small; this block map records the observed reading order rather than reproducing the source page's skin or assets.

```text
404
This page doesn’t exist.
Return Home
Docs · Knowledge Base · Blog
```

The cited page offers one recovery link (`Return Home`) plus links to the Vercel Docs, Knowledge Base, and Blog. The missing-route message and lightweight recovery navigation are the structure we are borrowing.

## Gelium remake

This is how **that structure** lands in Gelium: the same status-first message, with our recovery copy and a real GET link back to the section-reference index. The component is the registered server-rendered `ui-error-state`; it is not a screenshot, a Vercel clone, or a link to an unrelated recipe.

<div id="gelium-remake" class="ui-error-state" role="alert">
  <p class="ui-error-state-code" aria-hidden="true">404</p>
  <h2 class="ui-error-state-title">Page not found</h2>
  <p class="ui-error-state-body">This page doesn’t exist. Return to the section references index to continue.</p>
  <a class="ui-button" href="/docs/section-references">Back to section references</a>
</div>

## Keep / adapt

| Kept | Gelium primitive | Rejected | Why |
|---|---|---|---|
| Visible 404 status | `ui-error-state-code` | Vercel logo or product chrome | The reference is a not-found state, not a brand shell |
| Plain missing-page message | `ui-error-state-title` + `ui-error-state-body` | Copied wording as product copy | The job is shared; the voice belongs to Gelium |
| Single missing-route composition | `ui-error-state` with `role="alert"` | Decorative illustration or animation | Adds skin without improving recovery |
| Recovery destination | real `.ui-button` GET link | Dead or JS-only retry control | The user needs a known, keyboard-accessible next step |

## Ask before copying

- Should a missing route return to the documentation index, the product home, or the user's previous known context?
- Does the product need a support or status link for failures beyond a simple 404?
- Should the page retain the product shell, or is a focused error state safer for the failure context?
