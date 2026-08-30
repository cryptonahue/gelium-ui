# Gate Ledger and Preflight Specification

## Purpose

Define task-local gate evidence that can be validated in phases without pretending to prove cognitive reading, prevent arbitrary local edits, or authorize delivery.

## Requirements

### Requirement: Select one deterministic ledger

A design-gated task MUST select its ledger through one canonical consumer convention or an explicit command input. The ledger MUST declare `schema_version`, route classification, scope, owned/shared paths, gate statuses, and evidence references.

#### Scenario: Explicit ledger selection
- **GIVEN** a consumer has multiple surfaces under review
- **WHEN** preflight is invoked with `--ledger <path>`
- **THEN** it MUST validate only that ledger and report the selected path in all output.

#### Scenario: Invalid ledger fails closed
- **GIVEN** a ledger is malformed, uses an unknown schema version, has an unknown status, or lacks required ownership/evidence fields
- **WHEN** preflight validates it
- **THEN** it MUST report `invalid-configuration` and MUST NOT claim a passing gate.

### Requirement: Treat reading as attestation

A reading record MAY attest that a required file/path was reviewed. Tooling MUST validate its required shape and referenced file existence where possible, but MUST NOT claim to prove that an agent cognitively read or applied its content.

### Requirement: Separate prebuild and release preflight

The future preflight MUST expose distinct `prebuild` and `release` modes.

#### Scenario: Prebuild validates only prebuild evidence
- **GIVEN** a design-gated surface has a complete Plan and Architecture packet but no candidate implementation
- **WHEN** `prebuild` runs
- **THEN** it MUST require required artifacts, scope declaration, criteria plan, and approval/exemption as applicable, and MUST NOT require rendered viewport/theme/no-JS evidence.

#### Scenario: Release requires rendered evidence
- **GIVEN** a candidate implementation is ready for audit
- **WHEN** `release` runs
- **THEN** it MUST require the rendered audit references, detector result, applicable test/build evidence, and authority-matrix checks.

### Requirement: Validate scope coverage without semantic inference

Preflight MUST compare changed UI-relevant paths with declared owned/shared paths and report uncovered changes. It MUST NOT infer whether a redesign is substantial solely from paths or diff size.

#### Scenario: Uncovered changed template blocks release
- **GIVEN** a candidate changes a UI template outside the declared owned/shared paths
- **WHEN** preflight evaluates the declared diff scope
- **THEN** it MUST fail with the uncovered path and require ledger correction or scope escalation.

### Requirement: Keep approval and delivery authority separate

A preflight outcome MUST state evidence/gate status only. It MUST NOT commit, push, publish, deploy, or claim authority over those operations.

### Requirement: Require authority matrices before coherence checks

Before package/version or wire-contract checks can be enabled, maintainers MUST define equivalent authority surfaces and permitted historical references. A checker MUST report conflicts with each path/value and MUST NOT automatically rename a contract.
