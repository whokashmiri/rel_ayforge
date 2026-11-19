# Phase 12 --- Containerization, CI/CD and Deployment

**Status:** Not Started\
**Milestone:** M12 --- Reproducible deployment and portfolio release\
**Prerequisites:** Core backend functionality and tests are stable.

## 1. Objective

Package RelayForge for repeatable local and cloud deployment. Establish
automated quality checks, environment configuration, operational
documentation, backups and a public demonstration environment where
feasible.

## 2. Why this phase matters

A backend project is stronger when another developer can build, test and
run it reliably. Deployment also demonstrates operational awareness:
configuration, secrets, health checks, logs, monitoring, upgrades and
recovery.

## 3. Scope and features

-   Multi-stage Docker builds.
-   Docker Compose for local development and demo.
-   Separate API, worker and publisher processes from one codebase.
-   GitHub Actions CI pipeline.
-   Automated tests, linting and image builds.
-   Environment-specific configuration and secret management.
-   Cloud deployment.
-   TLS, reverse proxy and network restrictions.
-   Database backup and migration procedure.
-   Monitoring, alerting and operational runbooks.
-   Optional Kubernetes deployment after the simpler deployment is
    stable.

## 4. Concepts to learn

-   Container images and layers.
-   Multi-stage builds and minimal runtime images.
-   Process separation and configuration.
-   CI versus CD.
-   Secrets management.
-   Health checks and deployment readiness.
-   Rolling updates and backward-compatible migrations.
-   Backups, restore tests and disaster recovery.
-   Kubernetes deployments, services, probes and secrets (optional
    advanced extension).

## 5. Proposed files

``` text
Dockerfile
docker-compose.yml
.env.example
.github/workflows/ci.yml
deployments/docker/
deployments/kubernetes/
docs/operations/
```

## 6. Task sequence

### 12.1 --- Prepare production configuration

Separate development and production configuration. Validate required
environment variables at startup. Never commit real secrets. Document
rotation and access controls.

**Result:** Environments are configured without source-code edits.

### 12.2 --- Build multi-stage Docker images

Compile static or appropriately linked Go binaries in a build stage. Use
a small runtime image, non-root user, explicit entry point and no
development tools in production.

**Result:** API, worker and publisher can run in containers.

### 12.3 --- Complete Docker Compose

Define API, worker, publisher, PostgreSQL, Redis, NATS and monitoring
services. Add health checks, dependency readiness and named volumes.
Keep local credentials clearly development-only.

**Result:** A developer can run the complete stack locally.

### 12.4 --- Implement CI

On pull requests and pushes: - Format and verify Go code. - Run
`go vet`. - Run unit tests. - Run selected integration tests. - Run race
tests where practical. - Run frontend checks if dashboard exists. -
Build container images.

**Result:** Changes receive automated quality checks.

### 12.5 --- Add image and release workflow

Tag images with immutable commit identifiers. Publish only through a
controlled workflow. Avoid deploying unreviewed pull-request code to
production.

**Result:** Releases are traceable to source commits.

### 12.6 --- Prepare cloud environment

Choose a suitable hosting environment. Configure firewall rules, DNS,
TLS, reverse proxy, persistent database storage and private service
networking. Do not expose PostgreSQL, Redis or NATS publicly.

**Result:** Infrastructure is reachable only through intended paths.

### 12.7 --- Deploy and verify

Deploy the API, publisher and workers. Apply migrations using a
controlled release step. Verify health checks, event ingestion,
successful delivery, retry and dashboard connectivity.

**Result:** A working environment demonstrates the complete product.

### 12.8 --- Add backups and recovery

Document database backup frequency, retention, encryption and restore
procedure. Perform a restore test. Define what happens to pending outbox
and delivery records during recovery.

**Result:** Recovery is practiced, not merely assumed.

### 12.9 --- Add operational runbooks

Document restart procedures, broker outage response, growing DLQ,
database connection exhaustion, secret rotation, rollback and incident
investigation.

**Result:** Operators know how to respond to common failures.

### 12.10 --- Optional Kubernetes extension

Only after Docker deployment is stable, define Deployments, Services,
ConfigMaps/Secrets, readiness/liveness probes, resource requests/limits
and scaling policies. Keep stateful infrastructure managed or
deliberately designed; do not casually deploy databases in an ephemeral
setup.

**Result:** Kubernetes knowledge is demonstrated through a working
deployment, not just YAML files.

## 7. Acceptance criteria

-   [ ] A clean checkout can build the required images.
-   [ ] Docker Compose starts the documented local stack.
-   [ ] Containers run with least privilege where practical.
-   [ ] Production secrets are not committed or baked into images.
-   [ ] CI runs formatting, vetting, tests and builds.
-   [ ] Images are traceable to a commit.
-   [ ] Production data services are not publicly exposed.
-   [ ] TLS is enabled for public traffic.
-   [ ] Migrations have a documented release procedure.
-   [ ] Health checks and monitoring work after deployment.
-   [ ] A test event can be submitted and delivered in the deployed
    environment.
-   [ ] Backup restore has been tested.
-   [ ] Operational runbooks are available.
-   [ ] Public demo contains no real customer data or secrets.

## 8. Interview checkpoints

Explain multi-stage builds, container layers, CI/CD, secrets management,
health probes, migration compatibility, backup/restore, horizontal
scaling and when Kubernetes is justified.

## 9. Completion record

Record deployment architecture, cloud resources, release commit, CI run,
verification steps, backup restore evidence and known limitations.
