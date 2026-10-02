# Security

Use a unique webhook secret and reject unsigned or replayed deliveries. Treat fork pull requests as hostile: they cannot use secrets or the privileged worker, and require explicit approval plus a dedicated hardened runner or VM/microVM. Do not assume Docker containers isolate hostile tenants.

Preview access is intended to be authenticated by default. Store token hashes, encrypt secret values at rest using a managed key, redact values from logs, and scope runtime credentials to one deployment. Block cloud metadata and control-plane networks from preview networks. Apply CPU, memory, PID, filesystem, capability, network, and execution-time limits in the runtime adapter.
