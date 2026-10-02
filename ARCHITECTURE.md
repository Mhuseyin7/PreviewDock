# Architecture

PreviewDock maintains control-plane records in PostgreSQL and processes durable jobs through Redis. The controller validates every state transition and dispatches a signed, minimal execution request to an agent. Agents own only PreviewDock-labelled resources. Routing changes occur only after the candidate passes its configured health check. The database migration is in `migrations/0001_initial.sql`; its idempotency keys prevent duplicate GitHub delivery processing.

Production implementation boundaries: repository authorization precedes planning; planning pins the head SHA; builders may only use validated context paths; agents enforce resource limits and never mount the management Docker socket into a workload. A reconciliation loop compares owned runtime labels with recorded deployments and cleans only owned resources.
