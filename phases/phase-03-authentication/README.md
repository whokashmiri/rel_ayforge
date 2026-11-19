# Phase 03 --- API Key Authentication and Multi-Tenancy

**Status:** Not Started\
**Milestone:** M3 --- Authenticated tenant-scoped API\
**Prerequisites:** Phase 02 completed.

## 1. Objective

Allow each tenant to authenticate using API keys and ensure every tenant
can access only its own resources. This phase introduces the security
boundary that all later event and endpoint features will rely on.

## 2. Why this phase matters

RelayForge is a multi-tenant platform. Many organizations use the same
service, but their data must remain isolated. Authentication answers who
is making a request; tenant authorization determines which data that
caller may access.

## 3. Scope and features

-   Generate high-entropy API keys.
-   Display the raw key only once at creation.
-   Store only a cryptographic hash of the key.
-   Authenticate requests through middleware.
-   Attach the authenticated tenant identity to request context.
-   Support revocation and rotation.
-   Enforce tenant scoping in repository queries.
-   Add authentication and authorization tests.

**Out of scope:** User passwords, OAuth, SSO, role-based dashboard
access and billing.

## 4. Concepts to learn

-   Authentication versus authorization.
-   Cryptographically secure random generation.
-   Hashing versus encryption.
-   Constant-time comparison where applicable.
-   HTTP middleware and request context.
-   Least privilege and defense in depth.
-   Tenant isolation and insecure direct object reference (IDOR).
-   Secret handling and log redaction.

## 5. Proposed files

``` text
internal/domain/api_key.go
internal/service/auth_service.go
internal/repository/api_key_repository.go
internal/handler/api_key_handler.go
internal/middleware/authentication.go
migrations/002_create_api_keys.sql
```

## 6. Security design

-   Generate keys using `crypto/rand`, not `math/rand`.
-   Use a recognizable public prefix to help identify key types, but
    never rely on the prefix as authentication.
-   Hash the full secret using a suitable one-way cryptographic hash.
    Since keys are high entropy, a fast cryptographic hash such as
    SHA-256 can be appropriate; use a keyed HMAC hash if the design
    includes a separately protected server-side pepper.
-   Store metadata such as key ID, tenant ID, prefix, hash, created
    time, last-used time and revoked time.
-   Return the raw key once and never store or log it.
-   Compare derived hashes safely.
-   Require TLS in deployed environments.

## 7. Task sequence

### 3.1 --- Define API key lifecycle

Document creation, one-time display, authentication, revocation and
rotation behavior before coding.

**Result:** The credential lifecycle is unambiguous.

### 3.2 --- Create database migration

Create an API key table with tenant foreign key, key metadata, hash,
timestamps and revocation state. Add indexes for lookup and tenant
listing.

**Result:** API key metadata is stored safely.

### 3.3 --- Implement key generation

Use `crypto/rand` to create a high-entropy secret. Separate the public
identifier/prefix from the secret material.

**Result:** New keys are unpredictable.

### 3.4 --- Implement key storage and lookup

Hash the secret before storage. Authenticate by looking up the key
identifier and verifying the supplied secret against the stored hash.

**Result:** Raw API secrets are never persisted.

### 3.5 --- Implement API key management endpoints

-   `POST /api/v1/auth/api-keys`
-   `GET /api/v1/auth/api-keys`
-   `DELETE /api/v1/auth/api-keys/{id}`

Return the raw secret only from the creation response. Ensure management
operations are themselves protected by an appropriate bootstrap/admin
mechanism; do not leave key creation publicly available without a
controlled registration flow.

**Result:** Keys can be created, listed and revoked safely.

### 3.6 --- Implement authentication middleware

-   Read credentials from the documented authorization header.
-   Reject missing, malformed, invalid and revoked credentials.
-   Load the associated tenant.
-   Add a typed tenant identity to request context.
-   Avoid returning details that help attackers distinguish unknown keys
    from revoked keys.

**Result:** Protected routes know the authenticated tenant.

### 3.7 --- Enforce tenant isolation

-   Pass tenant identity explicitly into service/repository calls.
-   Include tenant ID in every tenant-owned query.
-   Do not trust tenant IDs supplied by request bodies for
    authorization.
-   Test attempts to access another tenant's resource.

**Result:** One tenant cannot read or modify another tenant's data.

### 3.8 --- Add security tests

Test missing key, invalid key, revoked key, valid key, tenant mismatch,
key rotation and secret redaction.

**Result:** The authentication boundary is verified.

## 8. Acceptance criteria

-   [ ] API keys are generated using a cryptographically secure random
    source.
-   [ ] Raw keys are shown only once and are not stored.
-   [ ] Revoked keys cannot authenticate.
-   [ ] Protected routes reject unauthenticated requests.
-   [ ] Authenticated tenant identity is available through request
    context.
-   [ ] Tenant-owned queries are always tenant-scoped.
-   [ ] Cross-tenant access tests fail safely.
-   [ ] Secrets do not appear in logs or error responses.
-   [ ] Key creation has a controlled authorization path.

## 9. Interview checkpoints

Explain authentication versus authorization, why passwords and API keys
are handled differently, why `crypto/rand` matters, how middleware
passes identity, and how you prevent IDOR and cross-tenant data leakage.

## 10. Completion record

Document the exact key format, hash strategy, revocation behavior and
tests. Never paste real API keys into documentation or Git history.
