# Phase 02 --- PostgreSQL and REST API Foundation

**Status:** Not Started\
**Milestone:** M2 --- A functioning API backed by PostgreSQL\
**Prerequisites:** Phase 01 completed.

## 1. Objective

Connect RelayForge to PostgreSQL and build the first persistent REST
resource: tenant registration and retrieval. Establish database access,
migrations, request validation, and a clear handler/service/repository
separation.

## 2. Why this phase matters

An API that stores data only in memory loses it when the process
restarts. PostgreSQL gives RelayForge durable relational storage. This
phase also introduces database boundaries without prematurely building
the entire event platform.

## 3. Scope and features

-   Run PostgreSQL locally with Docker Compose.
-   Connect using `pgx` and a connection pool.
-   Manage schema changes using Goose migrations.
-   Create the initial `tenants` table.
-   Implement tenant creation and retrieval.
-   Validate inputs and return consistent JSON errors.
-   Add database health checks and integration tests.

**Out of scope:** API keys, tenant authorization, events, NATS, Redis
and dashboard.

## 4. Concepts to learn

-   Relational data, primary keys, unique constraints and indexes.
-   SQL parameterization and SQL injection prevention.
-   Connection pools and resource cleanup.
-   Database migrations and schema versioning.
-   Transactions (introductory use only).
-   Handler, service and repository responsibilities.
-   HTTP methods, status codes and JSON contracts.
-   Unit tests versus integration tests.

## 5. Proposed files

``` text
internal/database/postgres.go
internal/domain/tenant.go
internal/handler/tenant_handler.go
internal/service/tenant_service.go
internal/repository/tenant_repository.go
migrations/001_create_tenants.sql
docker-compose.yml
```

The exact package names can be adjusted to match the codebase, but each
layer must have a clear reason to exist.

## 6. Initial data model

`tenants` should include: - `id` --- UUID primary key. - `name` ---
required organization name. - `created_at` and `updated_at` ---
timestamps. - Optional status field if lifecycle behavior is defined.

Add a unique constraint only for a field whose uniqueness is a real
product requirement. Do not invent unnecessary columns.

## 7. Task sequence

### 2.1 --- Add local PostgreSQL

-   Add a PostgreSQL service to Docker Compose.
-   Configure database name, user, password and port through environment
    variables.
-   Add a persistent development volume and health check.
-   Document how to start and stop the service.

**Result:** PostgreSQL runs locally in a repeatable way.

### 2.2 --- Implement database configuration

-   Parse the connection string from environment configuration.
-   Create a `pgxpool.Pool`.
-   Ping the database during startup.
-   Close the pool during shutdown.
-   Set sensible pool limits and connection lifetimes.

**Result:** The API manages database connections safely.

### 2.3 --- Add migrations

-   Configure Goose.
-   Write the tenant table migration and rollback.
-   Add commands for applying and checking migrations.
-   Test migrations against a clean database.

**Result:** Database schema can be recreated and evolved consistently.

### 2.4 --- Define tenant domain and repository

-   Define the tenant model.
-   Define a small repository interface at the consuming layer where
    practical.
-   Implement SQL queries with parameters.
-   Return clear not-found and conflict errors.

**Result:** Persistence details are separated from HTTP handling.

### 2.5 --- Implement service logic

-   Validate tenant names and business rules.
-   Call the repository.
-   Translate persistence outcomes into application-level errors.
-   Keep HTTP-specific status codes out of domain logic.

**Result:** Tenant behavior can be tested without an HTTP server.

### 2.6 --- Implement REST handlers

-   `POST /api/v1/tenants` creates a tenant.
-   `GET /api/v1/tenants/{id}` retrieves a tenant.
-   Validate JSON and reject malformed or oversized bodies.
-   Return consistent JSON error shapes.
-   Use 201, 200, 400, 404 and 409 appropriately.

**Result:** Tenant data can be created and retrieved through HTTP.

### 2.7 --- Add tests

-   Unit-test service validation with a fake repository.
-   Integration-test repository queries against PostgreSQL.
-   Test handler status codes and response bodies.
-   Verify migrations work from an empty database.

**Result:** API behavior and persistence are verified.

## 8. Acceptance criteria

-   [ ] PostgreSQL starts using documented Docker Compose commands.
-   [ ] Migrations apply and roll back successfully.
-   [ ] API fails clearly when required database configuration is
    invalid.
-   [ ] Database connections are closed on shutdown.
-   [ ] Tenant creation persists across API restarts.
-   [ ] Tenant retrieval returns the correct record.
-   [ ] Invalid JSON and invalid input return 400.
-   [ ] Unknown tenant IDs return 404.
-   [ ] SQL queries use parameters, not string concatenation.
-   [ ] Unit and integration tests pass.

## 9. Example API contract

`POST /api/v1/tenants`

``` json
{
  "name": "Acme Commerce"
}
```

Expected success: `201 Created`, returning the tenant ID and timestamps.
Do not return secrets in API responses.

## 10. Interview checkpoints

Explain connection pooling, migrations, SQL injection, parameterized
queries, repository boundaries, HTTP status codes, and when a database
transaction is needed.

## 11. Completion record

Record migration version, test commands, API examples and any schema
decisions. Update task checkboxes only after verification.
