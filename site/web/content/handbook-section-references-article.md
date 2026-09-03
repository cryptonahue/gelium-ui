# Article / rich post

Read · detail. One H1, a readable measure, then evidence. The original is a citation, not a skin.

**[See Gelium remake](#gelium-remake)**

[Open original](https://vercel.com/blog/the-end-of-credential-sprawl-for-agents) (Vercel, 25 Aug 2026)

[Section references](/docs/section-references)

## Original (cited)

Live URL: [vercel.com/blog/the-end-of-credential-sprawl-for-agents](https://vercel.com/blog/the-end-of-credential-sprawl-for-agents)

Our crop is a **block map** of reading order (2026-08-30). No Vercel assets in CSS.

```text
product chrome (out of scope)
H1  The end of credential sprawl for agents
    authors · date · 6 min read
lead + in-body product link (Vercel Connect)
hero product image (light/dark pair)
H2  Vaults don't fix long-lived tokens
    prose · CLI snippet · code sample
H2  What changes when access becomes a request
    comparison table
H2  Connectors for the services you already use
H2  Pricing and availability
related / explore cards
Ready to deploy?  (marketing CTA)
site footer
```

## Gelium remake

This is how **that structure** lands in Gelium: same jobs (title, byline, lead, H2 evidence, table, related), our copy, tokens, no their brand. It is not `/recipes/rich-article` (a different media fixture).

<article id="gelium-remake" class="docs-foundation-showcase" aria-label="Gelium reconstruction of the cited article structure">
  <header>
    <p>Read · article fixture</p>
    <h2>Access as a request</h2>
    <p>By Gelium UI · 30 Aug 2026 · 6 min read</p>
    <p>Long-lived secrets turn every agent into a standing grant. Issue access at the moment of the task, scoped to that task, and let it expire.</p>
  </header>
  <p role="note"><strong>No hero pair.</strong> A Read page does not need a light/dark marketing image to start. Evidence belongs in the body.</p>
  <h3>Standing tokens do not expire with the job</h3>
  <p>A vault hides a secret. It does not limit what the secret can do after it leaks. Rotation scripts and copied environment values are a second workload.</p>
  <pre><code>gelium access mint slack --name ops-bot
# identity comes from the runtime, not a second env secret</code></pre>
  <h3>What changes when access is requested</h3>
  <div class="ui-data-table-scroll">
    <table class="ui-data-table-table">
      <caption>Stored token versus request-time access</caption>
      <thead>
        <tr><th scope="col">Property</th><th scope="col">Stored token</th><th scope="col">Request-time</th></tr>
      </thead>
      <tbody>
        <tr><th scope="row">Lifetime</th><td>Until someone rotates it</td><td>Short-lived, refreshed</td></tr>
        <tr><th scope="row">Reach</th><td>Everything the agent might need</td><td>Scoped to this step</td></tr>
        <tr><th scope="row">Revocation</th><td>Redeploy copies</td><td>One command</td></tr>
      </tbody>
    </table>
  </div>
  <h3>Connectors you already run</h3>
  <p>Register a provider once. Attach it to the projects that need it. The app never stores the provider secret.</p>
  <h3>Availability</h3>
  <p>Say who can request access and which environments are allowed. Pricing copy belongs on a pricing screen, not as a second primary on the article.</p>
  <nav aria-label="Related content">
    <h3>Related content</h3>
    <ul class="ui-list">
      <li class="ui-list-item"><a href="/docs/media">Media foundations</a> <span>Captions, transcript, no arbitrary embeds.</span></li>
      <li class="ui-list-item"><a href="/docs/section-references">Section references</a> <span>Cite originals; reconstruct with tokens.</span></li>
    </ul>
  </nav>
</article>

## Keep / adapt

| Kept | Gelium primitive | Rejected | Why |
|---|---|---|---|
| One title + lead + byline | article header | Product mega-nav | Not this page’s chrome |
| Descriptive H2/H3 evidence | headings + prose + `pre` | Copy-link-to-heading chrome | Extra JS |
| Comparison table | `ui-data-table` | Light/dark decorative hero pair | Skin, not structure |
| Related content | `ui-list` + nav | “Ready to deploy?” band | Second primary on a Read page |
| Evidence in the body | note + table | External embed / tweet row | `MEDIA-EMBED` |

## Ask before copying

- Related stories as a list, or omit?
- Native media in the body, or prose + table only?
- Article JSON-LD on the **consumer** host, not copied from Vercel?
