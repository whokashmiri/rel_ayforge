# Phase 04 --- Event Ingestion and Transactional Outbox

**Status:** Not Started\
**Milestone:** M4 --- Events can be accepted and durably recorded\
**Prerequisites:** Phase 03 completed.

## 1. Objective

Build the core event ingestion workflow. Authenticated tenants will
submit event payloads; RelayForge will validate them, store them in
PostgreSQL, determine eligible endpoints, and create durable work
records.

This phase introduces the transactional outbox pattern so that event
persistence and the intention to publish work are committed together.

## 2. Why this phase matters

If the API writes an event to PostgreSQL and then tries to publish
directly to a message broker, the process can crash between those
operations. The event might exist in the database but never be
published. The outbox pattern records the event and pending publication
in one database transaction.

## 3. Scope and features

-   Event ingestion endpoint.
-   Tenant-scoped event storage.
-   Event type and payload validation.
-   Idempotency key support.
-   Endpoint subscription matching.
-   Delivery record creation.
-   Transactional outbox record creation.
-   Event retrieval and filtering.

**Out of scope:** Actual NATS publishing and outbound webhook HTTP
requests; those arrive in later phases.

## 4. Concepts to learn

-   JSON payload handling and validation.
-   SQL transactions and atomicity.
-   Unique constraints and conflict handling.
-   Idempotency and duplicate requests.
-   Outbox pattern and eventual consistency.
-   Event envelope design.
-   Tenant-scoped queries.
-   Database indexing and pagination.

## 5. Proposed data model

### events

-   `id`, `tenant_id`, `event_type`, `payload`, `idempotency_key`,
    `created_at`.

### deliveries

-   `id`, `tenant_id`, `event_id`, `endpoint_id`, `status`,
    `attempt_count`, `next_attempt_at`, timestamps.

### outbox

-   `id`, `event_id`, `message_type`, `payload` or reference,
    `published_at`, `created_at`, retry metadata.

Add appropriate unique constraints, foreign keys and indexes. Define
idempotency scope explicitly, such as unique per tenant and key. Decide
how long idempotency records remain valid and document the policy.

## 6. Task sequence

### 4.1 --- Define the event envelope

Specify required fields, event ID, tenant ownership, event type,
payload, creation time and optional metadata. Keep user payload separate
from RelayForge metadata.

**Result:** All event producers use a consistent contract.

### 4.2 --- Create migrations

Create events, deliveries and outbox tables. Add tenant-scoped indexes
and constraints. Include rollback scripts.

**Result:** The database can represent accepted events and pending
delivery work.

### 4.3 --- Implement event validation

Validate content type, body size, JSON syntax, event type format and
required fields. Reject unsupported or malformed input.

**Result:** Invalid requests are rejected before persistence.

### 4.4 --- Implement idempotent ingestion

Require or support an idempotency key. On a repeated request within the
defined scope, return the original result rather than creating a second
event. Define behavior when the same key is reused with a different
payload.

**Result:** Client retries do not create unintended duplicate events.

### 4.5 --- Resolve endpoint subscriptions

Find active endpoints belonging to the authenticated tenant that
subscribe to the event type. Do not allow a tenant to target another
tenant's endpoint.

**Result:** RelayForge knows which destinations should receive an event.

### 4.6 --- Persist atomically

Within one PostgreSQL transaction: 1. Insert the event. 2. Insert one
delivery record per eligible endpoint. 3. Insert the outbox record(s)
needed to publish delivery work. 4. Commit only if the required writes
succeed.

**Result:** The event and its pending work are recorded consistently.

### 4.7 --- Implement event APIs

-   `POST /api/v1/events`
-   `GET /api/v1/events`
-   `GET /api/v1/events/{id}`

Support pagination and tenant-scoped filters. Avoid returning another
tenant's data even when a valid foreign ID is supplied.

**Result:** Tenants can submit and inspect their own events.

### 4.8 --- Add tests

Test invalid payloads, large payloads, duplicate keys, conflicting
payload reuse, no matching endpoints, multiple matching endpoints,
transaction rollback and tenant isolation.

**Result:** Ingestion behavior is predictable under normal and failure
conditions.

## 7. Acceptance criteria

-   [ ] Valid authenticated requests create event records.
-   [ ] Invalid or oversized payloads are rejected.
-   [ ] Idempotent repeated submissions do not create duplicate events.
-   [ ] Reusing a key with a different payload follows documented
    behavior.
-   [ ] Delivery records are created for matching active endpoints.
-   [ ] Event, delivery and outbox writes are transactionally
    consistent.
-   [ ] A forced transaction failure leaves no partial event workflow.
-   [ ] Event list and detail endpoints are tenant-scoped and paginated.
-   [ ] Integration tests run against PostgreSQL.

## 8. Important design note

The outbox publisher may publish a record more than once if it crashes
after broker publication but before marking the record published.
Downstream processing must therefore be idempotent. The outbox improves
reliability; it does not magically provide exactly-once delivery.

## 9. Interview checkpoints

Explain ACID transactions, the dual-write problem, transactional outbox,
idempotency keys, eventual consistency, and why a database unique
constraint is stronger than an in-memory duplicate check.

## 10. Completion record

Record event schema version, idempotency policy, endpoint matching rules
and migration versions.
