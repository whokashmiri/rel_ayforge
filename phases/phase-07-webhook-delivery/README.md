# Phase 07 --- Webhook Delivery Engine

**Status:** Not Started\
**Milestone:** M7 --- Deliver signed HTTP webhooks to customer
endpoints\
**Prerequisites:** Phase 06 completed.

## 1. Objective

Implement the real delivery processor. For each delivery job, RelayForge
loads the event and endpoint, constructs a signed HTTP request, sends it
with strict safety controls, and records the result.

## 2. Why this phase matters

This is the point where RelayForge communicates with external software.
External endpoints may be slow, unavailable, misconfigured or malicious.
The delivery engine must be secure, observable and careful with
resources.

## 3. Scope and features

-   Endpoint registration and event subscriptions.
-   HTTP POST delivery.
-   HMAC-SHA256 signatures and timestamp headers.
-   Configurable request timeout.
-   Delivery status and attempt persistence.
-   Response code and bounded response-body capture.
-   Safe outbound network validation and SSRF defenses.
-   Per-endpoint and global concurrency controls.

## 4. Concepts to learn

-   `net/http` clients and transports.
-   Request contexts and timeouts.
-   HMAC signing and verification.
-   HTTP response handling.
-   Connection reuse and transport tuning.
-   SSRF and outbound request security.
-   Safe body limits and secret management.
-   Testing HTTP clients with `httptest`.

## 5. Proposed files

``` text
internal/webhook/client.go
internal/webhook/signer.go
internal/webhook/validator.go
internal/service/delivery_service.go
internal/repository/endpoint_repository.go
internal/repository/delivery_repository.go
```

## 6. Delivery contract

A webhook request should include: - A stable event ID. - Event type. -
Event creation timestamp. - JSON payload. - Signature header. -
Signature timestamp. - A delivery/attempt identifier where useful.

Define the exact signed bytes and header names in the API documentation.
Receivers must be able to verify signatures consistently.

## 7. Task sequence

### 7.1 --- Define endpoint model and subscription rules

Store tenant ownership, URL, active status, subscribed event types and
encrypted signing secret (or a carefully designed secret-storage
approach). Define secret rotation behavior.

**Result:** Each tenant can configure destinations safely.

### 7.2 --- Implement endpoint APIs

Create, list, update, disable and delete endpoints. Validate event
subscriptions and enforce tenant ownership.

**Result:** Tenants can manage destinations through the API.

### 7.3 --- Implement outbound URL validation

Reject unsupported schemes, credentials in URLs, loopback, private,
link-local, multicast and reserved addresses unless explicitly allowed
for local development. Revalidate DNS resolutions and defend against DNS
rebinding. Prevent redirects to prohibited destinations or validate
every redirect hop. Document local testing exceptions.

**Result:** RelayForge cannot be casually abused to reach internal
infrastructure.

### 7.4 --- Build the HTTP client

Use a reusable `http.Client` and configured transport. Set request
deadlines, connection limits, TLS verification and response-body limits.
Never disable TLS verification in production.

**Result:** Outbound requests use controlled reusable resources.

### 7.5 --- Implement signature generation

Use HMAC-SHA256 over a documented canonical byte sequence, including
timestamp and payload. Compare signatures with constant-time comparison
in verification tests. Support secret rotation where appropriate.

**Result:** Receivers can verify that webhook content came from
RelayForge and was not modified.

### 7.6 --- Implement delivery execution

Load delivery, event and endpoint records. Build the request, attach
headers, sign it, send it and classify the outcome. Persist status,
response code, duration and a safely truncated response excerpt.

**Result:** RelayForge can send and record webhook attempts.

### 7.7 --- Integrate with workers

Replace the fake processor with the real delivery service. Ensure only
one component owns acknowledgment decisions.

**Result:** Worker jobs perform real external HTTP deliveries.

### 7.8 --- Test delivery behavior

Use `httptest` servers to simulate success, 4xx, 5xx, timeout, malformed
response, slow response and connection failure. Test signatures and SSRF
validation.

**Result:** Delivery behavior is reproducible without relying on
third-party services.

## 8. Acceptance criteria

-   [ ] Tenants can create and manage their own endpoints.
-   [ ] Only subscribed endpoints receive matching events.
-   [ ] Requests include documented event and signature headers.
-   [ ] A test receiver can verify HMAC signatures.
-   [ ] Request timeouts and response-body limits are enforced.
-   [ ] Delivery attempts record outcome, duration and response status.
-   [ ] Credentials and signing secrets never appear in logs.
-   [ ] SSRF protections block prohibited destinations and redirects.
-   [ ] HTTP tests cover success and failure scenarios.
-   [ ] Worker acknowledgment behavior matches persisted delivery state.

## 9. Interview checkpoints

Explain HTTP transport reuse, request context, HMAC, why webhook
delivery is at-least-once, how to defend against SSRF, and how you
classify retryable versus permanent errors.

## 10. Completion record

Document signature format, headers, endpoint validation policy, timeout
defaults and test matrix.
