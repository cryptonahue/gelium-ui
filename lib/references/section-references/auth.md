# REF-AUTH

Source: https://vercel.com/signup
Kind: section · type: auth · tags: form, register, validation, recovery

## Structure

```text
auth identity / title
labeled fields
consent / context
one primary submit
validation summary + field errors
recovery or alternate auth path
```

## Gelium adaptation

- Native labeled controls; server-side validation (422)
- Validation summary plus inline errors
- Safe `next` / return path when the product has one
- No nested card theater for a simple form
- Do not claim routes that are not registered

## Reject

Fake social buttons without providers, invented password rules UI, client-only validation as the only gate.
