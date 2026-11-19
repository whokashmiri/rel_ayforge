# Phase 01 --- Go Foundation and Project Bootstrap

**Status:** Not Started\
**Milestone:** M1 --- A runnable Go application foundation\
**Prerequisites:** Go installed; Git available.

## 1. Objective

Create the initial RelayForge Go codebase and a small HTTP API process.
This phase establishes the conventions we will use throughout the
project. It is deliberately limited: no database, message broker,
authentication system, or worker pool yet.

By the end, a developer can clone the repository, configure it, start
the API, call a health endpoint, run tests, and understand how the code
is organized.

## 2. Why this phase matters

A reliable backend begins with a predictable application lifecycle and
understandable boundaries. We will learn Go by building real project
pieces rather than completing unrelated syntax exercises.

## 3. Scope and features

-   Initialize the Go module and pin dependencies.
-   Load configuration from environment variables with safe defaults.
-   Start an HTTP server using Chi.
-   Add health and readiness endpoints (readiness initially checks only
    that the process is initialized).
-   Add request ID and structured request logging middleware.
-   Handle graceful shutdown on operating-system signals.
-   Add a minimal Makefile and development instructions.

**Out of scope:** PostgreSQL, Redis, NATS, tenant authentication,
business endpoints, React, Kubernetes.

## 4. Concepts to learn

-   Go modules, packages, imports and `go.mod`.
-   `package main`, `func main`, exported identifiers and naming
    conventions.
-   Structs, methods, pointers and composition.
-   Interfaces and why small interfaces are useful.
-   Explicit error handling and error wrapping.
-   Environment variables and configuration validation.
-   `net/http`, handlers, middleware and routing.
-   `context.Context`, cancellation and server shutdown.
-   Unit tests with the standard `testing` package.

## 5. Proposed files

``` text
cmd/api/main.go
internal/config/config.go
internal/httpserver/server.go
internal/handler/health_handler.go
internal/middleware/request_id.go
internal/middleware/logging.go
.env.example
.gitignore
Makefile
go.mod
```

Create files only when their responsibility is needed. Avoid empty
packages and speculative abstractions.

## 6. Task sequence

### 1.1 --- Initialize the repository and module

-   Create the `relayforge` repository.
-   Run `go mod init` using the chosen repository module path.
-   Add `.gitignore` entries for local environment files, binaries,
    coverage output and editor files.
-   Add a concise `.env.example` containing non-secret development
    settings.

**Result:** Go recognizes the repository as a module and generated/local
files are not committed.

### 1.2 --- Learn the application entry point

-   Create `cmd/api/main.go`.
-   Keep `main` focused on loading configuration, constructing
    dependencies, starting the server and handling shutdown.
-   Avoid placing business logic in `main`.

**Result:** The API has one clear startup path.

### 1.3 --- Implement configuration

-   Define a `Config` struct.
-   Read `HTTP_ADDR`, `HTTP_READ_TIMEOUT`, `HTTP_WRITE_TIMEOUT` and
    `HTTP_SHUTDOWN_TIMEOUT`.
-   Parse durations and return useful errors for invalid values.
-   Keep secrets out of source control and logs.

**Result:** Configuration is centralized and validated at startup.

### 1.4 --- Create the HTTP server

-   Use Chi for routing.
-   Configure server address and timeouts.
-   Register `/health` and `/ready`.
-   Return JSON responses with appropriate status codes and content
    types.

**Result:** The API can accept and respond to HTTP requests.

### 1.5 --- Add middleware

-   Generate or accept a safe request ID.
-   Attach it to the request context and response header.
-   Log method, path, status, duration and request ID.
-   Do not log request bodies, credentials or secret headers.

**Result:** Requests can be correlated in logs without exposing
sensitive data.

### 1.6 --- Implement graceful shutdown

-   Listen for interrupt/termination signals.
-   Cancel the application context.
-   Call `http.Server.Shutdown` with a bounded timeout.
-   Log shutdown errors and ensure the process exits predictably.

**Result:** The server stops accepting new requests and allows in-flight
requests a chance to finish.

### 1.7 --- Add tests and developer commands

-   Test configuration parsing and invalid values.
-   Test health/readiness responses.
-   Test request ID behavior.
-   Add `make run`, `make test`, `make fmt` and `make lint` placeholders
    or working targets as appropriate.

**Result:** Common development tasks are repeatable.

## 7. Acceptance criteria

-   [ ] `go test ./...` succeeds.
-   [ ] `go vet ./...` succeeds.
-   [ ] `gofmt` has been applied.
-   [ ] API starts using documented environment variables.
-   [ ] `GET /health` returns HTTP 200 and JSON.
-   [ ] `GET /ready` returns HTTP 200 after startup.
-   [ ] Each response includes a request ID.
-   [ ] Logs contain request metadata but no credentials.
-   [ ] Ctrl+C triggers graceful shutdown.
-   [ ] A new developer can follow the README to run the API.

## 8. Verification commands

``` bash
go version
go mod tidy
gofmt -w .
go test ./...
go vet ./...
go run ./cmd/api
```

Then call `GET http://localhost:8080/health` (or the configured
address).

## 9. Interview checkpoints

Be able to explain: 1. What is a Go module, and what does `go.mod`
contain? 2. Why does `main` belong in `cmd/api`? 3. What is an interface
in Go? 4. Why does Go return errors instead of using exceptions for
ordinary error flow? 5. What is middleware? 6. Why should an HTTP server
have timeouts? 7. What does graceful shutdown protect?

## 10. Completion record

When this phase is finished, update the status at the top, record the
commit hash, note any deviations from this plan, and link any
architecture decisions. Do not mark a task complete until its acceptance
criteria have been verified.
