# Phase 06 --- Concurrent Delivery Worker Pool

**Status:** Not Started\
**Milestone:** M6 --- Bounded concurrent job processing\
**Prerequisites:** Phase 05 completed.

## 1. Objective

Build the worker runtime that processes delivery jobs concurrently.
Workers will consume durable jobs, coordinate work safely, respect
concurrency limits and shut down gracefully.

## 2. Why this phase matters

A sequential processor can become slow when many external endpoints
respond slowly. Go's goroutines make concurrent I/O practical, but
concurrency must be bounded and managed. Unlimited goroutines can
exhaust memory, sockets, database connections or downstream services.

## 3. Scope and features

-   Worker process entry point.
-   Configurable worker count.
-   Job dispatch and worker lifecycle.
-   Controlled concurrency and backpressure.
-   Context cancellation.
-   Graceful shutdown.
-   Panic containment and worker error reporting.
-   Metrics for active, completed and failed work.

**Out of scope:** Complete retry policy and final webhook request
implementation; the worker will initially call a delivery interface that
later receives the real implementation.

## 4. Concepts to learn

-   Goroutines and scheduler basics.
-   Channels and `select`.
-   Worker pool pattern.
-   `sync.WaitGroup`, mutexes and atomics.
-   Context cancellation and deadlines.
-   Data races and `go test -race`.
-   Bounded queues and backpressure.
-   Avoiding goroutine leaks.

## 5. Proposed files

``` text
cmd/worker/main.go
internal/worker/pool.go
internal/worker/delivery_worker.go
internal/worker/job.go
internal/worker/processor.go
```

## 6. Design principles

-   Set a configurable maximum number of concurrent jobs.
-   Do not create one goroutine per unbounded incoming message.
-   Use the broker's pending limits and worker capacity to apply
    backpressure.
-   Define ownership of message acknowledgment clearly.
-   Do not share mutable state without synchronization.
-   A worker failure must not silently kill the entire process.
-   Do not recover panics indiscriminately; if recovery is used at a
    worker boundary, log context and mark the job outcome safely.

## 7. Task sequence

### 6.1 --- Define the job contract

Create a delivery job type containing the identifiers needed to load
canonical records. Avoid passing unnecessary payload copies through
channels.

**Result:** Workers receive a consistent job representation.

### 6.2 --- Define the processor interface

Define the smallest interface the worker needs, such as processing a
delivery ID with a context. Use a fake implementation in early tests.

**Result:** Pool mechanics can be tested independently of HTTP delivery.

### 6.3 --- Implement one worker

Create a worker loop that receives jobs, handles cancellation, calls the
processor and reports outcomes.

**Result:** One worker processes jobs predictably.

### 6.4 --- Implement the pool

Start a configured number of workers and feed them through a bounded
jobs channel or broker consumer concurrency mechanism. Avoid introducing
an unnecessary in-memory queue if JetStream already provides the durable
queue.

**Result:** Work is processed concurrently with a clear limit.

### 6.5 --- Connect the JetStream consumer

Pull or receive messages according to worker capacity. Ensure
acknowledgment, negative acknowledgment and in-progress extension
behavior match the processing duration and failure policy.

**Result:** Broker messages are handed to available workers safely.

### 6.6 --- Implement graceful shutdown

Stop receiving new jobs, cancel work only according to the documented
shutdown policy, wait for in-flight jobs up to a deadline and leave
unfinished messages unacknowledged or negatively acknowledged for
redelivery.

**Result:** Shutdown does not falsely report unfinished jobs as
successful.

### 6.7 --- Add observability hooks

Track active workers, jobs started, completed, failed, duration and
queue/consumer lag where available.

**Result:** Worker behavior can be inspected.

### 6.8 --- Test concurrency and lifecycle

Use controlled fake processors to test maximum concurrency,
cancellation, worker errors, shutdown and race safety.

**Result:** The pool is verified under concurrent load.

## 8. Acceptance criteria

-   [ ] Worker count is configurable.
-   [ ] Concurrent processing never exceeds the configured limit.
-   [ ] A slow job does not block every other worker.
-   [ ] Cancellation is propagated.
-   [ ] Shutdown stops intake and waits within a defined deadline.
-   [ ] Unfinished jobs are not acknowledged as successful.
-   [ ] Worker errors are visible in logs and metrics.
-   [ ] No data races are reported by `go test -race ./...`.
-   [ ] Tests demonstrate bounded concurrency and no leaked workers.

## 9. Interview checkpoints

Explain goroutines versus OS threads at a high level, channels,
WaitGroups, mutexes, race conditions, backpressure, context cancellation
and why bounded concurrency matters.

## 10. Completion record

Record default worker count, tuning configuration, shutdown policy,
race-test output and benchmark observations.
