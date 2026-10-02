# Threat model

Assets include GitHub App keys, session credentials, deployment secrets, Docker host control, preview data, and routing authority. Threats include forged/replayed webhooks, malicious PR Dockerfiles, SSRF through health URLs, cross-tenant object access, and accidental deletion of unrelated Docker resources.

Controls: HMAC SHA-256 comparison in constant time, bounded payloads, delivery idempotency in PostgreSQL, authorization at every tenant query, SHA-pinned checkouts, strict paths and config validation, a separate authenticated agent, labelled-resource cleanup, and segregated untrusted workers. Remaining risk: a Docker-only worker is not adequate for hostile code; use a separate host or microVM isolation.
