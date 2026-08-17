# Media

Media is the figure system for non-image content: native audio players, transcripts, and third-party embeds. Use it when a page carries audio, a transcript that pairs with it, or a video embed — each with a safe, accessible, zero-JavaScript contract. The markup styles through `media.css` and follows the templates `media.html` defines (`audio`, `transcript`, and `embed`).

## Examples

An audio figure with native controls and a transcript link.

<div class="specimen-block">
<figure class="ui-media ui-media-audio">
  <audio controls preload="metadata">
    <source src="/static/audio/launch-intro.mp3" type="audio/mpeg">
    <p class="ui-media-fallback">Your browser does not support audio.</p>
  </audio>
  <figcaption>The launch episode — eight minutes.</figcaption>
  <p><a href="#transcript-launch">Read transcript</a></p>
</figure>

<section class="ui-transcript" id="transcript-launch" aria-labelledby="transcript-launch-heading">
  <h2 id="transcript-launch-heading">Transcript</h2>
  <div class="ui-transcript-content">Gelium UI ships with native HTML semantics and zero component JavaScript. This episode walks through the server-driven state contract, the token system, and the themes.</div>
</section>
</div>

An embed that is not on the allowlist renders a consent boundary instead of the iframe.

<div class="specimen-block">
<figure class="ui-media ui-media-embed">
  <div class="ui-embed-boundary">
    <p>This optional embed is unavailable until the source is approved.</p>
    <a href="/docs/media">Review our embed policy</a>
  </div>
  <figcaption>Third-party embeds require an allowlisted source and explicit consent.</figcaption>
</figure>
</div>

The specimens above follow the contracts the template `media.html` defines.

## Guidance

### When to use

Use the media pattern for audio with a transcript, a transcript on its own, or an embeddable video that is on the allowlist. It earns its place when the media is content people consume or read on the page — not a decorative flourish.

### When not to use

Do not use media for images — that is the [Image](/components/image) pattern. Do not embed a third-party iframe that is not allowlisted; the consent boundary renders instead. Do not auto-play audio or video: playback starts only when the person chooses it.

### Usability

- Audio renders with native controls and `preload="metadata"` so the page does not download the whole file up front.
- Pair audio with a transcript link: the transcript lives in its own `ui-transcript` section targetable by that link.
- Embeds render a lazy `iframe` only for allowlisted sources; everything else gets the consent boundary with a link to the policy.

### Accessibility

- The audio element exposes native controls, so keyboard and assistive-technology users get the platform player.
- The transcript is a labelled section (`aria-labelledby` pointing at its heading), so it is navigable and its content is plain text.
- The embed iframe always carries a `title` describing its content, and the consent fallback is real text with a real link.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.

## Anatomy

- **`ui-media ui-media-audio`** — the audio figure: a native `<audio controls>` with typed `<source>` elements and a text fallback paragraph.
- **`ui-media-fallback`** — the frame caption shown when the format is unsupported.
- **`ui-transcript`** — the transcript section: a labelled heading plus `ui-transcript-content`, rendered with preserved white space.
- **`ui-media ui-media-embed`** — the embed figure: an allowlisted iframe with `loading="lazy"` and a `title`, or a `ui-embed-boundary` consent block with a policy link.

## Template contract

The `audio` template renders the figure shell, the player with `preload` (default `metadata`), typed sources, the fallback text, an optional caption, and an optional transcript link. The `transcript` template renders the section with its heading and content. The `embed` template renders the iframe only when both the source is allowed and a `Src` is present; otherwise it renders the boundary with the fallback copy and optional consent link.

## Anti-patterns

- Do not auto-play audio or video; the user starts playback.
- Do not embed unknown third-party origins directly — route them through the allowlist and consent boundary.
- Do not drop the transcript for audio that carries information; the transcript is part of the accessible contract.

## Checklist

- The audio figure has native controls and typed sources with a fallback.
- The transcript section is labelled and linkable from the figure.
- The embed is allowlisted, lazy-loaded, and titled; unapproved sources render the consent boundary.
- Nothing auto-plays.

## Accessibility

Native audio controls and lazy, titled iframes keep the media operable by keyboard and assistive technology. The transcript is plain, labelled text — the durable form of the audio content — and the consent boundary for third-party embeds keeps the user informed before any network request to a foreign origin is made.
