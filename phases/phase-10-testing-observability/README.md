# Phase 10 --- Testing, Observability and Performance

**Status:** Not Started\
**Milestone:** M10 --- A measurable and testable platform\
**Prerequisites:** Phases 01--09 completed.

## 1. Objective

Establish confidence in RelayForge through automated tests, structured
logs, metrics, distributed tracing and performance experiments. This
phase strengthens the whole system rather than adding a new business
capability.

## 2. Why this phase matters

A service is not production-ready merely because its happy path works.
Engineers need to detect failures, understand latency, find bottlenecks
and reproduce bugs. Tests protect behavior as the architecture evolves.

## 3. Scope and features

-   Unit, integration and end-to-end tests.
-   Race detector and static analysis.
-   Structured logs with correlation IDs.
-   Prometheus metrics.
-   Grafana dashboards.
-   OpenTelemetry traces.
-   Health and readiness checks.
-   Load tests and profiling.
-   Operational runbooks.

## 4. Concepts to learn

-   Test doubles, fakes and dependency boundaries.
-   Integration test infrastructure.
-   Race detection and fuzz testing.
-   RED metrics: rate, errors and duration.
-   USE metrics: utilization, saturation and errors.
-   Metrics cardinality.
-   Trace spans and context propagation.
-   pprof CPU and memory profiling.
-   Load-test design and reproducibility.

## 5. Task sequence

### 10.1 --- Define quality gates

Document required commands for formatting, vetting, linting, unit tests,
integration tests and race detection. Establish minimum test coverage
expectations without treating coverage percentage as proof of quality.

**Result:** Every change has repeatable quality checks.

### 10.2 --- Expand unit tests

Test domain rules, validation, retry calculations, state transitions,
signing and rate-limit policy with deterministic inputs.

**Result:** Business rules are tested quickly.

### 10.3 --- Build integration tests

Use disposable or dedicated PostgreSQL, Redis and NATS test services.
Verify repository behavior, transactions, outbox publishing, consumer
handling and delivery persistence.

**Result:** Infrastructure integration is verified.

### 10.4 --- Add end-to-end tests

Run the API, publisher, broker, worker and a local webhook receiver.
Submit an event and verify delivery, failure, retry and replay
workflows.

**Result:** Complete user journeys are tested.

### 10.5 --- Run race and robustness tests

Run `go test -race ./...`. Add fuzz tests for parsers/validators where
useful. Test cancellation, restarts, broker interruption and slow
endpoints.

**Result:** Concurrency and edge cases receive focused testing.

### 10.6 --- Standardize structured logging

Use `log/slog` or the chosen logger. Include request ID, tenant ID where
safe, event/delivery ID and component. Redact secrets and avoid
unbounded payload logging.

**Result:** Logs support diagnosis without leaking sensitive data.

### 10.7 --- Add Prometheus metrics

Track request rate, error rate, latency, event ingestion, delivery
outcomes, worker activity, retry counts, DLQ size and broker lag where
available. Avoid tenant IDs or event IDs as metric labels to prevent
high cardinality.

**Result:** System health can be measured.

### 10.8 --- Add Grafana dashboards

Create dashboards for API health, delivery outcomes, worker utilization,
retry activity and infrastructure health. Add useful panel descriptions
and units.

**Result:** Operators can understand system behavior visually.

### 10.9 --- Add OpenTelemetry tracing

Create spans across API ingestion, outbox publishing, message
consumption and webhook delivery. Propagate trace context where
supported and avoid putting secrets or full payloads into span
attributes.

**Result:** A single event's path can be investigated across components.

### 10.10 --- Load test and profile

Use k6 to create documented workloads. Measure throughput, latency
percentiles, error rate and resource usage. Use pprof to investigate
bottlenecks before optimizing.

**Result:** Performance claims are based on reproducible measurements.

## 6. Acceptance criteria

-   [ ] Unit and integration test suites run from documented commands.
-   [ ] End-to-end test covers successful delivery.
-   [ ] End-to-end tests cover failure and retry.
-   [ ] Race detector reports no known data races in tested paths.
-   [ ] Logs include correlation fields and redact secrets.
-   [ ] Prometheus exposes useful service metrics.
-   [ ] Grafana dashboards load and display data.
-   [ ] Traces connect key stages of event processing.
-   [ ] Load-test workload and environment are documented.
-   [ ] Performance results include hardware/configuration and
    percentile latency.
-   [ ] pprof can be used in a controlled development environment.

## 7. Interview checkpoints

Explain unit versus integration tests, race detection, metrics versus
logs versus traces, cardinality, p95/p99 latency, throughput, profiling
and why optimization should be evidence-driven.

## 8. Completion record

Attach test commands, dashboard screenshots, sample trace, load-test
report and profiling findings. Never claim a throughput figure without
its test conditions.
