# Video

Video is a responsive container for native video playback. Use it when a page embeds a local video file — a demo, a walkthrough, a product clip — that needs controls, captions, and a laid-out frame without JavaScript. It is a container, not a content component: the video asset and its captions are provided by the page, and the container supplies the ratio, the frame, and the native player.

## Examples

A 16:9 video with native controls, lazy loading, captions, and a fallback.

<div class="ui-video">
  <video controls poster="/static/video/walkthrough-poster.jpg" loading="lazy" crossorigin="anonymous">
    <source src="/static/video/walkthrough.mp4" type="video/mp4">
    <track kind="captions" src="/static/video/walkthrough.en.vtt" srclang="en" label="English">
    <p class="ui-video-fallback">Your browser does not support HTML video.</p>
  </video>
</div>

The 4:3 variant keeps the same contract with a squarer ratio.

<div class="ui-video ui-video--aspect-4-3">
  <video controls loading="lazy">
    <source src="/static/video/short.mp4" type="video/mp4">
    <p class="ui-video-fallback">Your browser does not support HTML video.</p>
  </video>
</div>

The specimens above follow the contract the template `video.html` defines.

## Guidance

### When to use

Use the video container when the page owns a video file and the player belongs in the layout at a fixed ratio — a product demo, a tutorial, a short clip. It earns its place "best used inside another component": [Split](/components/split), [Card](/components/card), or the [Feature card](/components/feature-card) media slot.

### When not to use

Do not use the container for third-party embeds — that is the [Media](/components/media) embed contract with its allowlist and consent boundary. Do not auto-play video: playback starts only when the person chooses it. If the video is decorative background movement, a CSS background is the right tool, not a player.

### Usability

- Default to the 16:9 ratio; use the `ui-video--aspect-4-3` modifier only when the footage is genuinely 4:3.
- Provide typed `<source>` elements and a text fallback so unsupported browsers still show a message inside the frame.
- Add a poster image when the first frame is not a good resting state; the template emits `poster` only when one is supplied.

### Accessibility

- The player uses native `<video controls>`, so keyboard and assistive-technology users get the platform player.
- Add a captions `<track>` for spoken content; the template sets `crossorigin="anonymous"` on the video so the caption file can be fetched cross-origin.
- The frame keeps a `CanvasText` boundary in forced-colors mode so the video region stays visible when color is removed.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.

## Anatomy

- **`.ui-video`** — the frame: a 16:9 `aspect-ratio` (literal, deliberately not tokenized — structural geometry like breakpoints and z-index), the theme small radius, a border from `--ui-video-border`, and the surface-container background behind the player.
- **`.ui-video--aspect-4-3`** — the 4:3 modifier, same frame, squarer ratio.
- **`.ui-video-fallback`** — the caption inside the player for browsers that cannot play the source.

## Variants

- `ui-video` — the 16:9 default.
- `ui-video ui-video--aspect-4-3` — the 4:3 modifier.

## Anti-patterns

- Do not add `autoplay`; the user starts playback.
- Do not omit captions for video with spoken or meaningful audio.
- Do not tokenize the ratio: aspect-ratio stays literal, same rule as Card media and Feature card media.

## Checklist

- The video is inside `ui-video` with its ratio class.
- The player has `controls` and `loading="lazy"`.
- Spoken content has a captions `track` with `srclang` and a label.
- The sources are typed and a fallback paragraph is present.
- Nothing auto-plays.

## Accessibility

Native controls, a captions track, and a text fallback keep the player operable and understandable without JavaScript. The frame's forced-colors boundary keeps the video region visible when color is removed, and the container's own ratio reserves layout space so the page does not shift while the file loads.