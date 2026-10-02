# Configuration

`previewdock.yml` is committed beside the application code and controls the deployment contract. The controller must parse it with a strict YAML decoder, reject unknown fields, and run `internal/config.Config.Validate` before any build begins.

Supported version: `1`.

| Field | Requirement |
| --- | --- |
| `build.context` / `build.dockerfile` | Relative paths below the checked-out repository; parent traversal and absolute paths are rejected. |
| `service.container_port` | Integer from 1 through 65535. |
| `health.path` | Relative HTTP path beginning with `/`; use a dedicated health endpoint. |
| `health.expected_status` | HTTP status from 100 through 599. |
| `resources.cpu` / `resources.memory` | Positive bounded limits enforced by the execution agent. |
| `preview.ttl` | Between one hour and 30 days. |
| `preview.authentication` | `required`, `password`, or explicit `public`. |

Use `previewdock.yml.example` as the initial file. Build and runtime secret values do not belong in this file or in Git; they are stored separately and never returned by list APIs.
