# Phase 08 --- Retry Scheduling and Dead Letter Queue

**Status:** Not Started\
**Milestone:** M8 --- Recover from failed deliveries and support replay\
**Prerequisites:** Phase 07 completed.

## 1. Objective

Build a deliberate failure-management lifecycle for webhook deliveries.
Temporary failures should be retried with controlled backoff; exhausted
or permanently invalid deliveries should be visible in a Dead Letter
Queue (DLQ) and eligible for authorized manual replay.

## 2. Why this phase matters

External systems fail. A timeout or HTTP 503 should not necessarily mean
a delivery is permanently lost. At the same time, retrying every failure
forever can overload systems and waste resources.

## 3. Scope and features

-   Delivery status state machine.
-   Retry classification.
-   Exponential backoff with jitter.
-   Configurable maximum attempts.
-   Scheduled retry execution.
-   DLQ state and query APIs.
-   Authorized manual replay.
-   Attempt history and audit trail.

## 4. Concepts to learn

-   Transient versus permanent errors.
-   Exponential backoff and jitter.
-   State machines and valid transitions.
-   Retry storms and circuit-breaker concepts.
-   At-least-once delivery.
-   Idempotent replay.
-   Scheduling and durable delayed work.
-   Operational audit trails.

## 5. Suggested delivery states

Define and document a finite state machine, for example:

`pending → processing → succeeded`

`processing → retry_scheduled → processing`

`processing → dead_lettered`

`dead_lettered → retry_scheduled` (manual replay)

The exact states and transition rules must be consistent in code and
database constraints. Handle worker crashes and stale `processing`
records.

## 6. Task sequence

### 8.1 --- Define failure classification

Classify timeouts, connection errors, HTTP 408, 429 and selected 5xx
responses as potentially retryable. Treat most other 4xx responses as
permanent by default, with explicit handling for 429 and any documented
exceptions.

**Result:** Retry decisions are predictable.

### 8.2 --- Define retry policy

Set configurable base delay, maximum delay, maximum attempts and jitter.
Use exponential backoff with a cap. Respect `Retry-After` where
appropriate and safe.

**Result:** Retries are bounded and avoid synchronized retry bursts.

### 8.3 --- Implement delivery state transitions

Centralize allowed state transitions and make updates concurrency-safe.
Use database conditions or transactions to prevent two workers from
claiming the same delivery.

**Result:** Delivery lifecycle cannot be changed arbitrarily.

### 8.4 --- Persist retry schedule

Store next attempt time and attempt count. Use a durable
scheduler/publisher that makes due jobs available without keeping all
scheduled work only in process memory.

**Result:** Scheduled retries survive restarts.

### 8.5 --- Implement retry worker flow

Select due retries, claim them safely, publish or process them and
update state. Prevent duplicate scheduling where possible and rely on
idempotent handling for redelivery.

**Result:** Due deliveries are retried automatically.

### 8.6 --- Implement DLQ management

Mark exhausted or permanent failures as dead-lettered. Expose
tenant-scoped listing and detail APIs with attempt history and failure
reason.

**Result:** Failed deliveries are inspectable.

### 8.7 --- Implement manual replay

Add an authorized replay endpoint. Record who/what initiated replay,
original delivery reference and replay time. Enforce tenant scope and
replay limits.

**Result:** Operators or customers can safely request another attempt.

### 8.8 --- Test failure scenarios

Test retry timing, jitter bounds, max attempts, permanent errors, 429
behavior, process restart, duplicate messages, concurrent claims and
replay authorization.

**Result:** Recovery behavior is proven.

## 7. Acceptance criteria

-   [ ] Retryable and permanent failures follow documented rules.
-   [ ] Backoff is exponential, capped and jittered.
-   [ ] Retry attempts never exceed configured limits automatically.
-   [ ] Retry schedules survive application restarts.
-   [ ] Concurrent workers cannot claim the same delivery unsafely.
-   [ ] Exhausted deliveries appear in the DLQ.
-   [ ] DLQ records include useful failure and attempt history.
-   [ ] Replay is authenticated, tenant-scoped and audited.
-   [ ] Duplicate broker messages do not create duplicate successful
    effects.
-   [ ] Failure and recovery integration tests pass.

## 8. Important delivery guarantee

RelayForge provides at-least-once delivery attempts, not exactly-once
effects at a receiver. A receiver may process a request and then lose
the response before RelayForge records success. Include a stable
event/delivery identifier so receivers can implement deduplication.

## 9. Interview checkpoints

Explain exponential backoff, jitter, retry classification, poison
messages, DLQs, state machines, retry storms and the limits of
exactly-once delivery.

## 10. Completion record

Record retry defaults, classification rules, state diagram, DLQ
retention and replay audit behavior.
