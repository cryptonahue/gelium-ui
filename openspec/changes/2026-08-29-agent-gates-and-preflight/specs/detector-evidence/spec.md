# Detector Evidence Specification

## Purpose

Extend Gelium UX evidence so detector findings are scoped, attributable, and honest about bounded exceptions, media metadata gaps, and shared layouts.

## Requirements

### Requirement: Preserve raw findings

The detector MUST preserve a raw finding even when a declared exception matches it. A matching exception MUST never transform the result into a clean pass.

#### Scenario: Approved exception remains visible
- **GIVEN** a finding matches an approved bounded exception with rule/fingerprint, path/selector scope, reason, risk, owner, ledger evidence, and deterministic expiry
- **WHEN** detector output is produced
- **THEN** it MUST report `pass-with-exceptions`, include the raw finding and exception ID, and NOT report `clean-pass`.

#### Scenario: Expired exception fails
- **GIVEN** a finding matches an exception whose required `expires_at` timestamp has passed
- **WHEN** detector output is produced
- **THEN** it MUST report the exception as expired and fail unless a new approved exception exists.

> Schema v1 uses `expires_at` as its single deterministic exception-expiry mechanism. Version-based expiry is not accepted by this schema.

### Requirement: Attribute shared layouts explicitly

The detector MUST distinguish an owned finding from a declared shared-layout finding. Shared findings MAY be out of the audited surface but MUST remain visible in machine and audit output.

#### Scenario: Shared logout form is not silently discarded
- **GIVEN** a screen ledger declares a logout form path as shared and outside its owned form surface
- **WHEN** the detector finds a form-contract issue there
- **THEN** it MUST emit an attributed shared finding and retain it in the audit record without calling the screen a clean detector pass solely by path exclusion.

### Requirement: Support machine-readable result classes

The detector MUST expose text and machine-readable output with one of `clean-pass`, `pass-with-exceptions`, `failed`, or `invalid-configuration`. Existing positional/default behavior MUST remain compatible until migration guidance changes it.

### Requirement: Represent unknown media dimensions honestly

For informative media with meaningful alt text and safe responsive containment but no trustworthy intrinsic dimensions, guidance and detector output MUST record `media-metadata-unknown` or `pass-with-escalation`. They MUST NOT require or generate fabricated `width`/`height` values.

#### Scenario: Unknown external image metadata
- **GIVEN** an external image URL has no trustworthy dimensions and the consumer cannot provide server-known intrinsic metadata
- **WHEN** the image is rendered
- **THEN** the audit MUST record the metadata limitation, preserve alt/recovery/responsive obligations, and not claim an unconditional clean media pass.

