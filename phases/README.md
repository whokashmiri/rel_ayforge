# RelayForge --- Implementation Guide Index

**Status:** Planning\
**Purpose:** Detailed, task-by-task implementation specifications for
RelayForge.

This directory is the working implementation plan. The root `README.md`
describes the product; these phase guides describe how we intend to
build it.

## Development sequence

  ----------------------------------------------------------------------------------------------------------
  Phase                   Guide                                                      Main outcome
  ----------------------- ---------------------------------------------------------- -----------------------
  01                      [Go Foundation](phase-01-go-foundation/README.md)          Runnable Go HTTP
                                                                                     service

  02                      [PostgreSQL and REST                                       Persistent tenant API
                          API](phase-02-postgresql-rest-api/README.md)               

  03                      [Authentication](phase-03-authentication/README.md)        API keys and tenant
                                                                                     isolation

  04                      [Event Ingestion](phase-04-event-ingestion/README.md)      Durable event and
                                                                                     outbox records

  05                      [NATS Integration](phase-05-nats-integration/README.md)    Durable asynchronous
                                                                                     messaging

  06                      [Worker Pool](phase-06-worker-pool/README.md)              Bounded concurrent
                                                                                     processing

  07                      [Webhook Delivery](phase-07-webhook-delivery/README.md)    Signed HTTP webhook
                                                                                     delivery

  08                      [Retry and DLQ](phase-08-retry-dlq/README.md)              Failure recovery and
                                                                                     replay

  09                      [Redis and Rate                                            Tenant limits and
                          Limiting](phase-09-redis-rate-limiting/README.md)          delivery throttles

  10                      [Testing and                                               Automated verification
                          Observability](phase-10-testing-observability/README.md)   and monitoring

  11                      [React Dashboard](phase-11-react-dashboard/README.md)      Developer-facing
                                                                                     console

  12                      [Deployment](phase-12-deployment/README.md)                CI/CD and deployed
                                                                                     system
  ----------------------------------------------------------------------------------------------------------

## How to use these guides

1.  Read the current phase from beginning to end before coding.
2.  Work through tasks in order unless an explicit dependency decision
    changes the sequence.
3.  Ask questions and understand the Go concepts before implementing a
    task.
4.  Run the stated tests and verify the acceptance criteria.
5.  Update task checkboxes and phase status only after the feature
    works.
6.  Record implementation deviations and architectural decisions.
7.  Keep this documentation aligned with the actual code.

## Status convention

-   **Not Started** --- No implementation work has been verified.
-   **In Progress** --- Work has begun; some tasks may be complete.
-   **Blocked** --- A dependency or unresolved decision prevents
    progress.
-   **Completed** --- All acceptance criteria have been verified.

A planned feature is not a completed feature. Do not claim unimplemented
functionality in the root README, CV or portfolio.

## Project-wide engineering rules

-   PostgreSQL is the durable source of truth for business records.
-   NATS JetStream transports asynchronous work; it does not replace
    PostgreSQL.
-   Delivery is at-least-once. Design consumers and receivers for
    idempotency; do not claim exactly-once external effects.
-   Keep concurrency bounded and propagate cancellation.
-   Enforce tenant isolation in backend services and database queries.
-   Treat outbound webhook URLs as untrusted input and defend against
    SSRF.
-   Never commit real credentials, API keys or signing secrets.
-   Introduce infrastructure only when the current phase needs it.
-   Prefer a modular Go codebase with separately runnable processes
    before considering microservices.

## Architecture decisions

Record significant decisions in `docs/architecture/` or as ADRs,
including: - Module path and package layout. - API key format and
hashing strategy. - Event envelope and schema versioning. - Idempotency
scope and retention. - NATS stream/consumer configuration. - Delivery
state machine and retry classification. - Outbound URL validation
policy. - Dashboard authentication model. - Deployment topology and
backup policy.
