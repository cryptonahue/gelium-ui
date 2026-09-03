# Hero / product direction

Persuade · hub. A product hero should establish the promise, give brief context, and offer one obvious way to begin.

**[See Gelium remake](#gelium-remake)**

[Open original](https://linear.app) (Linear, verified 30 Aug 2026; HTTP 200)

[Section references](/docs/section-references)

## Original (cited)

Live URL: [linear.app](https://linear.app)

The Linear home page opens with a product-development promise, a short positioning line, and a Coding Sessions action. It then moves from product proof into a sequence of numbered capabilities before a closing conversion region. This block map records the section order, not Linear's brand, assets, or motion.

```text
skip link + product navigation (chrome, out of scope)
H1  The product development system for teams and agents
lead  Purpose-built for planning and building products
primary action  New Coding Sessions
product proof / issue activity demonstration
positioning statement  A new species of product tool
three principles  Purpose-built · Powered by agents · Designed for speed
numbered capability sequence  Intake → Plan → Build → Diffs → Monitor
customer proof and closing conversion links
```

The source is a Persuade landing page: the first scan establishes identity and the next move, while the later sections substantiate the promise.

## Gelium remake

This is how **that landing structure** lands in Gelium: the registered `ui-hero` composition owns the page heading, subtitle, and one primary link. A native `ui-list` carries the supporting principles without importing a logo, animation, gradient, or product screenshot.

<section id="gelium-remake" aria-labelledby="hero-linear-title">
  <section class="ui-hero">
    <div class="ui-hero-content">
      <p class="ui-hero-eyebrow">Product direction</p>
      <h2 class="ui-hero-title" id="hero-linear-title">Make the next release easier to see</h2>
      <p class="ui-hero-subtitle">A calm workspace for turning intent into shipped work, with the important state visible at every step.</p>
      <div class="ui-hero-actions">
        <a class="ui-button ui-button-primary" href="/docs/section-references">Explore the system</a>
      </div>
    </div>
  </section>
  <section aria-labelledby="hero-linear-principles">
    <h3 id="hero-linear-principles">Built for forward motion</h3>
    <ul class="ui-list">
      <li class="ui-list-item"><strong>Purpose-built</strong><span>Keep planning, decisions, and delivery in one legible flow.</span></li>
      <li class="ui-list-item"><strong>Agent-ready</strong><span>Give people and automation the same visible work context.</span></li>
      <li class="ui-list-item"><strong>Focused</strong><span>Reduce noise so the next useful action is easy to find.</span></li>
    </ul>
  </section>
</section>

The hero has one primary action and no background media. If the product later supplies a real asset, it belongs in `ui-hero-media` with honest dimensions and a decorative `alt`, not in a copied landing-page screenshot.

## Keep / adapt

| Kept | Gelium primitive | Rejected | Why |
|---|---|---|---|
| Promise → context → next move | `ui-hero` + Button link | Product logo and navigation clone | The reference is a hero structure, not a shell to copy |
| One dominant conversion path | `ui-button ui-button-primary` | Multiple hero primaries | Persuade screens need one clear first move |
| Three supporting principles | `ui-list` | Icon tiles and decorative badges | Text remains scannable and meaningful without assets |
| Product proof after the first scan | later semantic sections | Auto-playing video or cursor animation | Motion adds skin and a runtime dependency, not proof |
| Responsive message hierarchy | token-driven `ui-hero` | Hard-coded gradients and breakpoints | The Gelium primitive owns theme and narrow behavior |

## Ask before copying

- What is the one action the landing page must earn, and what happens after it succeeds?
- Which proof is verifiable product evidence rather than a decorative demo or unowned claim?
- Does the page need media at all, and who provides accessible, licensed, responsive assets if it does?
