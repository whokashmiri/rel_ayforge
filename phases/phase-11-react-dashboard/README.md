# Phase 11 --- React Developer Dashboard

**Status:** Not Started\
**Milestone:** M11 --- User-facing developer console\
**Prerequisites:** Backend API and security behavior from earlier phases
are stable.

## 1. Objective

Build a React dashboard that allows a tenant to manage its RelayForge
integration without manually calling every API endpoint. The dashboard
is a client of the Go API; it must not bypass backend authorization or
business rules.

## 2. Why this phase matters

A developer platform needs a usable interface for configuring endpoints,
inspecting events and diagnosing delivery problems. The dashboard also
provides a demonstrable end-to-end product experience for the portfolio.

## 3. Scope and features

-   Tenant sign-in or controlled API-key onboarding flow.
-   Overview dashboard.
-   API key creation/revocation.
-   Webhook endpoint management.
-   Event explorer and filters.
-   Delivery detail and attempt timeline.
-   DLQ inspection and replay.
-   Usage and reliability metrics.
-   Loading, empty, error and permission states.

**Out of scope:** Replacing the Go API, storing raw API secrets after
display, or exposing administrative cross-tenant operations.

## 4. Concepts to learn

-   React components and hooks.
-   TypeScript types.
-   API client design.
-   Authentication state and secure browser storage trade-offs.
-   Form validation.
-   Server state and pagination.
-   Error handling and accessibility.
-   Dashboard information architecture.

## 5. Proposed frontend structure

``` text
web/
  src/
    app/
    components/
    features/
      overview/
      api-keys/
      endpoints/
      events/
      deliveries/
      dead-letter-queue/
    lib/
      api-client.ts
    types/
    routes/
```

Choose the frontend tooling and styling system when this phase begins.
Keep API contracts documented independently of the UI.

## 6. Task sequence

### 11.1 --- Define user journeys and wireframes

Map the key journeys: onboard tenant, register endpoint, submit test
event, inspect delivery and replay a failed delivery.

**Result:** Dashboard screens solve defined user tasks.

### 11.2 --- Bootstrap React and TypeScript

Create the frontend application, formatting/linting and environment
configuration. Keep API base URLs configurable.

**Result:** A reproducible frontend development environment exists.

### 11.3 --- Build shared application shell

Implement navigation, layout, route handling, loading states and error
boundaries.

**Result:** Screens share consistent navigation and behavior.

### 11.4 --- Integrate authentication

Implement the chosen session or API credential flow. Do not casually
persist long-lived secret API keys in local storage. If browser-based
user authentication is added, prefer a secure server-managed session
design with HttpOnly, Secure and appropriate SameSite cookies.

**Result:** Dashboard access follows a documented security model.

### 11.5 --- Build endpoint management

Create, edit, disable and delete endpoints; configure event
subscriptions; display signing-secret creation/rotation once where
appropriate.

**Result:** Tenants can manage destinations.

### 11.6 --- Build event explorer

Show event type, ID, creation time and status. Add pagination and
filters. Avoid rendering untrusted payloads as HTML.

**Result:** Users can inspect submitted events.

### 11.7 --- Build delivery details

Display destination, status, attempt timeline, response code, timing and
safe error information.

**Result:** Users can diagnose delivery outcomes.

### 11.8 --- Build DLQ and replay screens

Show dead-lettered deliveries and require a clear confirmation before
replay. Display replay outcome and audit information.

**Result:** Users can recover failed deliveries through authorized
actions.

### 11.9 --- Build overview metrics

Show event volume, success/failure rates, delivery latency and DLQ
counts using backend metrics or API summaries.

**Result:** Users get a useful platform overview.

### 11.10 --- Accessibility and UI tests

Test keyboard navigation, labels, focus behavior, responsive layouts and
loading/error states. Add component and critical journey tests.

**Result:** Dashboard is usable and robust.

## 7. Acceptance criteria

-   [ ] Dashboard can connect to a configurable API environment.
-   [ ] Tenant users cannot access another tenant's data.
-   [ ] Endpoint create/update/disable/delete workflows work.
-   [ ] Event explorer supports pagination and filtering.
-   [ ] Delivery detail shows attempt history.
-   [ ] Replay requires authorization and confirmation.
-   [ ] API errors and loading states are understandable.
-   [ ] No long-lived raw API key is silently stored in browser storage.
-   [ ] UI handles untrusted event payloads safely.
-   [ ] Core journeys have automated tests.
-   [ ] Dashboard works at desktop and mobile viewport sizes.

## 8. Interview checkpoints

Explain component state, server state, API client boundaries,
authentication storage trade-offs, pagination, accessibility and why
frontend visibility is not a substitute for backend authorization.

## 9. Completion record

Record user journeys, screenshots, API contract versions and test
results.
