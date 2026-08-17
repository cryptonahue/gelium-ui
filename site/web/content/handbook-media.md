# Media

Media is content, not decoration. Keep the same HTML-first contract across editorial image, picture, video, audio, transcript, and optional embeds. Use local fixtures or a neutral placeholder in examples; never add copyrighted media to the repository.

## Image and picture

`MEDIA-IMAGE` is the baseline contract for informative editorial images.

Use a meaningful `alt` for informative images and `alt=""` only when the image is genuinely decorative. Reserve intrinsic `width` and `height` for every editorial image to prevent layout shift. Add `loading="lazy"` for below-the-fold content; keep the first meaningful image eager. Use `srcset` and `sizes` when you have real variants, and keep an aspect-ratio fallback when media is delayed or missing.

```html
<img src="/static/media/editorial-placeholder.svg" alt="Abstract green and blue geometric placeholder" width="1200" height="675" loading="lazy" srcset="..." sizes="(max-width: 48rem) 100vw, 65ch">
```

`<picture>` is for art direction or format selection, not a reason to duplicate the same source. Preserve the `img` fallback and its accessible name.

## Video and captions

Use the existing native `<video controls>` primitive. Provide a poster, one or more typed `<source>` elements, and a captions `<track kind="captions">` when speech or meaningful audio exists. Include a textual fallback so a failed or unsupported video still explains what happened.

Captions are not a transcript substitute: captions synchronize sound and speech with the timeline; a transcript gives the complete text as a readable alternative.

## Audio and transcript

Use native `<audio controls>` with typed sources and `preload="metadata"` by default. Link an optional transcript immediately after the player and keep the transcript in the document, not behind JavaScript. A transcript must identify speakers or meaningful sound when that context matters.

<div class="docs-media-demo">
  <audio controls preload="metadata"><source src="/static/media/empty-audio.ogg" type="audio/ogg"><p>Audio is unavailable. Read the <a href="#transcript-example">transcript</a> instead.</p></audio>
</div>

## Safe embed boundary

Do not invent arbitrary third-party iframe URLs. The embed primitive has an explicit server-side approval boundary: render an iframe only when the source is allowlisted and consent has been established. Otherwise render the fallback with a clear link to the canonical source. Keep the fallback useful without requiring a script, cookie, or external request.

This is intentionally a policy boundary, not an embed registry. Product code owns the allowlist, title, dimensions, privacy review, and consent copy.

## Responsive behavior and states

Media must shrink to its container (`max-width: 100%`) without `overflow-x: hidden` on the page. Use an aspect-ratio wrapper or intrinsic dimensions to reserve space, and let narrow embeds use an internal responsive frame. Test a 390px viewport as well as desktop.

- **Loading:** reserve dimensions; lazy-load below the fold; avoid skeletons that hide the real control.
- **Error:** keep native fallback text and offer a transcript, download, or canonical link.
- **Empty:** say why no media is present and offer the next useful action; never leave a blank frame.

## When / when not

| Use this slice when | Choose something else when |
|---|---|
| The user needs to read, hear, watch, or compare editorial media. | A decorative shape can be CSS or an icon with no content value. |
| A responsive image has real crops or encoded variants. | You only have one source; do not add a pretend `srcset`. |
| An external player is approved, consented, and has a canonical fallback. | The provider, URL, privacy, or fallback is unknown. |
| Audio/video has a transcript or captions appropriate to its content. | The media is a silent decorative loop; use a decorative image instead. |

## Sources

- [W3C WCAG 2.2 — 1.1.1 Non-text Content](https://www.w3.org/WAI/WCAG22/Understanding/non-text-content.html): informative images need an equivalent text alternative.
- [W3C WCAG 2.2 — 1.2.2 Captions (Prerecorded)](https://www.w3.org/WAI/WCAG22/Understanding/captions-prerecorded.html): prerecorded synchronized media with audio needs captions.
- [W3C WCAG 2.2 — 1.2.3 Audio Description or Media Alternative](https://www.w3.org/WAI/WCAG22/Understanding/audio-description-or-media-alternative-prerecorded.html): provide a media alternative or audio description where needed.
- [GOV.UK Design System — Images](https://design-system.service.gov.uk/styles/images/): official guidance for responsive, purposeful images.

See also: [Responsive](/docs/responsive), [Accessibility](/docs/accessibility), [Performance](/docs/performance).
