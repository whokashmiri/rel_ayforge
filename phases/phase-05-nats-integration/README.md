# Phase 05 --- NATS JetStream Integration

**Status:** Not Started\
**Milestone:** M5 --- Durable asynchronous messaging\
**Prerequisites:** Phase 04 completed.

## 1. Objective

Introduce NATS JetStream to move delivery work from the API/database
workflow into an asynchronous processing pipeline. Implement an outbox
publisher and durable consumers.

## 2. Why this phase matters

The API should acknowledge an accepted event without waiting for every
destination to respond. A durable broker buffers work, separates
producers from consumers and allows workers to be unavailable
temporarily without losing acknowledged messages.

## 3. Scope and features

-   Run NATS with JetStream enabled.
-   Define subjects and stream configuration.
-   Publish pending outbox messages.
-   Use durable consumers.
-   Acknowledge, negatively acknowledge or leave messages unacknowledged
    according to documented policy.
-   Handle reconnects and broker outages.
-   Make message processing idempotent.

**Out of scope:** Final webhook delivery behavior, retry scheduling
policy and DLQ user interface.

## 4. Concepts to learn

-   Message brokers and pub/sub.
-   Subjects, streams and consumers.
-   Durable versus ephemeral consumers.
-   Acknowledgment and redelivery.
-   At-least-once processing.
-   Message schema/versioning.
-   Outbox publisher loops.
-   Backpressure and broker retention.

## 5. Proposed files

``` text
internal/messaging/connection.go
internal/messaging/publisher.go
internal/messaging/consumer.go
internal/messaging/subjects.go
internal/outbox/publisher.go
cmd/publisher/main.go
```

## 6. Message contract

Define a versioned message containing a delivery ID, event ID, tenant
ID, endpoint ID, event type and schema version. Prefer references to
canonical database records over copying sensitive or unnecessarily large
payloads into messages. Validate message contents before processing.

## 7. Task sequence

### 5.1 --- Run NATS JetStream locally

Add NATS to Docker Compose with JetStream enabled, persistent storage
and a health check. Document the local connection URL and
development-only credentials.

**Result:** A repeatable local broker is available.

### 5.2 --- Define stream and subject conventions

Choose stable names, for example a delivery subject namespace. Document
stream retention, storage, replicas for development versus production,
maximum message size and retention limits.

**Result:** Messaging names and retention behavior are explicit.

### 5.3 --- Implement NATS connection lifecycle

Create a reusable connection component with startup validation,
reconnect behavior, logging and clean shutdown.

**Result:** Processes connect to NATS predictably.

### 5.4 --- Implement outbox polling

Claim pending outbox rows safely using transactions and an appropriate
locking strategy (for example `FOR UPDATE SKIP LOCKED`). Avoid two
publisher instances claiming the same row simultaneously.

**Result:** Multiple publisher instances can work without unnecessary
duplicate claims.

### 5.5 --- Publish messages

Publish versioned delivery messages with a stable message ID. Use broker
deduplication features where appropriate, but retain application-level
idempotency.

**Result:** Outbox work is forwarded to JetStream.

### 5.6 --- Mark publication state

After a successful publish acknowledgment, update the outbox row. If
publication fails, record the failure and retry with bounded backoff.

**Result:** The publisher can resume after interruptions.

### 5.7 --- Create durable consumer

Create a durable consumer with explicit acknowledgment behavior,
sensible pending limits and controlled redelivery settings.

**Result:** Delivery workers can receive work reliably.

### 5.8 --- Implement consumer handler skeleton

Decode and validate the message, load the delivery record and
acknowledge only after the chosen processing handoff is safe. Do not
acknowledge before durable state is secured.

**Result:** Messages enter the processing pipeline without being
silently discarded.

### 5.9 --- Test broker failure scenarios

Stop NATS, restart it, interrupt the publisher and interrupt a consumer.
Verify outbox records remain available and messages are redelivered as
designed.

**Result:** Recovery behavior is demonstrated.

## 8. Acceptance criteria

-   [ ] JetStream starts with persistent local storage.
-   [ ] Stream and durable consumer configuration is reproducible.
-   [ ] Outbox records publish and are marked only after broker
    acknowledgment.
-   [ ] Publisher handles temporary broker unavailability.
-   [ ] Multiple publishers do not corrupt outbox state.
-   [ ] Consumers can resume after restart.
-   [ ] A message is not acknowledged before safe handoff.
-   [ ] Duplicate messages do not create duplicate business effects.
-   [ ] Broker outage and recovery tests pass.

## 9. Interview checkpoints

Explain pub/sub, durable streams, acknowledgments, redelivery,
at-least-once semantics, outbox publishing, consumer idempotency and why
a message broker does not replace the relational database.

## 10. Completion record

Document subject names, stream settings, consumer settings, retention
policy and failure-test results.
