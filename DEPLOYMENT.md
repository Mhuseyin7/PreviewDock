# Deployment

Use Compose for the management services only. Place the API behind TLS, restrict PostgreSQL and Redis to the internal network, provision persistent database backups, and inject secrets through a secret manager rather than committing `.env`. Configure wildcard DNS and a proxy only after ownership and certificates are verified. A production runner must be a separate restricted host; do not expose its Docker socket through the API or previews.
