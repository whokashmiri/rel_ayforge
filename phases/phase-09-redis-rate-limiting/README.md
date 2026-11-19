# Phase 09 --- Redis and Rate Limiting

**Status:** Not Started\
**Milestone:** M9 --- Protect shared resources and enforce usage limits\
**Prerequisites:** Phase 08 completed.

## 1. Objective

Introduce Redis for fast shared counters and short-lived coordination,
then implement tenant-level API rate limiting and delivery throttles.
PostgreSQL remains the source of truth for durable business data.

## 2. Why this phase matters

A single tenant or burst of traffic should not consume all API capacity
or overwhelm external webhook destinations. Rate limits protect
RelayForge, promote fair usage and reduce accidental or abusive traffic.

## 3. Scope and features

-   Redis connection lifecycle.
-   Tenant-level request limits.
-   Endpoint or destination-level delivery limits where appropriate.
-   Standard rate-limit response headers.
-   Configurable policies.
-   Redis outage behavior.
-   Optional caching only where a measured use case exists.

## 4. Concepts to learn

-   Redis data structures and atomic commands.
-   Fixed-window, sliding-window and token-bucket algorithms.
-   Race conditions across multiple API instances.
-   TTL and key design.
-   Cache invalidation.
-   Fail-open versus fail-closed trade-offs.
-   Distributed rate limiting.

## 5. Proposed files

``` text
internal/cache/redis.go
internal/middleware/rate_limit.go
internal/service/rate_limit_service.go
internal/config/config.go
```

## 6. Task sequence

### 9.1 --- Add Redis to local infrastructure

Add Redis to Docker Compose with health checks and persistent
configuration only if required. Use environment variables for connection
settings.

**Result:** Redis is available for development.

### 9.2 --- Implement Redis client lifecycle

Create a client, validate connectivity, configure pool/timeouts and
close it during shutdown.

**Result:** Redis resources are managed safely.

### 9.3 --- Select and document an algorithm

Choose an algorithm based on desired fairness and burst behavior. A
token bucket is a reasonable starting point. Document key format, time
units and atomicity requirements.

**Result:** Rate-limit behavior is intentional.

### 9.4 --- Implement atomic tenant limits

Use an atomic Redis script or suitable command pattern so concurrent API
instances cannot overspend the same limit.

**Result:** Limits work across multiple API processes.

### 9.5 --- Add HTTP middleware

Apply limits after tenant authentication. Return HTTP 429 and useful
`Retry-After` or rate-limit headers. Do not trust a tenant identifier
supplied by the client.

**Result:** API requests are throttled per authenticated tenant.

### 9.6 --- Add delivery throttles

Add bounded concurrency or rate controls per destination so a failing or
rate-limited customer endpoint is not flooded.

**Result:** Outbound delivery respects configured capacity.

### 9.7 --- Define Redis outage policy

Decide whether each limiter fails open or closed. Document the
availability/security trade-off and emit metrics when Redis is
unavailable.

**Result:** Redis outages do not create surprising behavior.

### 9.8 --- Test distributed behavior

Test boundary conditions, simultaneous requests, TTL expiry, multiple
API instances and Redis outage behavior.

**Result:** Limits are verified under concurrency.

## 7. Acceptance criteria

-   [ ] Redis connects and shuts down cleanly.
-   [ ] Limits are based on authenticated tenant identity.
-   [ ] Concurrent requests cannot exceed the chosen algorithm's defined
    allowance.
-   [ ] HTTP 429 and relevant headers are returned when limited.
-   [ ] Delivery throttles protect external destinations.
-   [ ] Redis outage behavior matches documented policy.
-   [ ] Metrics expose allowed, rejected and backend-error counts.
-   [ ] Tests cover boundary, expiry and concurrent request cases.

## 8. Design caution

Do not put durable event, delivery or audit records only in Redis. Redis
is used for fast coordination and temporary state; PostgreSQL remains
the durable source of truth.

## 9. Interview checkpoints

Explain token bucket versus fixed window, atomic Redis operations, TTL,
distributed rate limiting, cache-aside, and fail-open versus fail-closed
decisions.

## 10. Completion record

Record algorithm choice, limit defaults, key schema, outage policy and
concurrency test results.
