# Agent Workflow Specification

## Purpose

Define a proportional Gelium UI workflow that preserves direct work for small understood changes while requiring explicit product, architecture, approval, and evidence gates for new screens, new flows, and substantial redesigns.

## Requirements

### Requirement: Route UI work proportionally

The workflow MUST classify a UI task as `direct-exempt`, `design-gated`, `escalate`, or `full-sdd` before applying conditional gates. The classification MUST be recorded for `design-gated`, `escalate`, and `full-sdd` work. A file-count heuristic MAY trigger focused exploration but MUST NOT itself determine risk or force SDD.

#### Scenario: Small correction stays direct
- **GIVEN** a change corrects an existing token, selector, copy string, documented contract, or bounded accessibility issue without changing page architecture
- **WHEN** the agent classifies the task
- **THEN** it MAY use `direct-exempt` and MUST run only the applicable checks without manufacturing a wireframe approval ceremony.

#### Scenario: Substantial redesign is design-gated
- **GIVEN** a change materially changes hierarchy, major regions, primary action, reading order, or audience/owner boundary
- **WHEN** the agent classifies the task
- **THEN** it MUST use `design-gated` and complete Orient, Plan, Architect, and Approve before Build.

#### Scenario: Cross-cutting ambiguity uses durable planning
- **GIVEN** a change spans multiple normative layers and durable proposal/spec/design/task artifacts materially reduce ambiguity
- **WHEN** the maintainer accepts the planning route
- **THEN** the work MUST use `full-sdd` OpenSpec artifacts without making every consumer screen an SDD transaction.

### Requirement: Separate Orient, Plan, and Architect

A design-gated task MUST keep constraints, product intent, and implementation viability as distinct outputs.

#### Scenario: Orient establishes non-negotiable constraints
- **GIVEN** a design-gated UI task begins
- **WHEN** Orient runs
- **THEN** the agent MUST inspect applicable product/design artifacts or record their absence, the Gelium decision pack, relevant vocabulary/registry, and existing hard contracts before presenting a wireframe as viable.

#### Scenario: Plan produces intent without fabricated implementation
- **GIVEN** Orient is complete
- **WHEN** Plan runs
- **THEN** it MUST state job, audience, surface/screen, primary action, states, non-goals, and an intent wireframe without inventing components, data, or media metadata.

#### Scenario: Architect produces the approval packet
- **GIVEN** a Plan wireframe exists
- **WHEN** Architect runs
- **THEN** it MUST inspect the relevant route, handler, templates, data, contracts, components, and no-JS behavior and produce a buildable wireframe plus section/contract mapping.

### Requirement: Require a human architecture decision only when gated

For a design-gated task, Build MUST wait for an `approved`, `exempt`, or bounded `exception` Architecture packet. `draft`, `changes-requested`, and `declined` MUST block Build. A direct-exempt task MUST NOT be blocked by an unnecessary approval ceremony.

#### Scenario: Approved packet permits Build
- **GIVEN** a design-gated packet includes route/contracts, desktop/mobile order, actions, states/recovery, accessibility/no-JS, reuse, and trade-offs
- **WHEN** a human records `approved` with approver, date/channel, and packet version after seeing the buildable wireframe in the conversation
- **THEN** Build MAY begin within the approved scope.

#### Scenario: Unseen packet does not permit Build
- **GIVEN** a design-gated packet exists only in a plan file, or the user said to make the page / `continua` / resume after a model switch without having been shown the wireframe
- **WHEN** the agent would start markup
- **THEN** it MUST stop with `Needs your decision`, show the buildable desktop and mobile wireframes in the conversation, and MUST NOT treat that continue as approval.

#### Scenario: Material deviation reopens the gate
- **GIVEN** Build reveals a change to approved hierarchy, primary action, route contract, or ownership boundary
- **WHEN** the deviation is material
- **THEN** the agent MUST return to Route/Architect and request a revised decision before claiming the implementation conforms.

### Requirement: Keep public progress simple

The workflow SHOULD expose `Working`, `Needs your decision`, `Checking`, and `Ready` to normal users. Detailed gate statuses MAY exist in the ledger but MUST NOT be required for a user to understand the next decision.

### Requirement: Bound subagent and memory use

The workflow MAY use at most two focused read-only explorers for broad discovery, one production writer, and a fresh audit reviewer. Explorers MUST receive exact paths and return bounded findings. Project memory MAY seed Orient but MUST be validated against current repository artifacts and MUST NOT act as approval or audit evidence.

#### Scenario: Memory conflicts with repository state
- **GIVEN** project memory says a component or contract has a prior decision
- **WHEN** current repository artifacts disagree
- **THEN** the current repository artifacts MUST govern and the agent MUST record the discrepancy if it affects the task.
