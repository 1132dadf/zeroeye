# WebSocket origin security

## Current review item

The WebSocket upgrader in `server.go` currently accepts every origin. That may be intentional for local development, but it is unsafe as a blanket production default when browser clients can reach the endpoint.

Origin checks are not a replacement for authentication, but they help prevent an untrusted website from using a victim's browser to open a connection to a service that trusts ambient credentials.

## Deployment guidance

- For production, configure an explicit allowlist of trusted web origins.
- Keep local-development origins opt-in rather than accepting every origin globally.
- Do not treat the Origin header as user authentication; non-browser clients can forge it.
- Authenticate connections and authorize each subscription independently.
- Apply connection, message-size, and per-client rate limits.
- Log rejected connection metadata without recording credentials or full sensitive payloads.

## Follow-up implementation checklist

- Add a server configuration field for allowed origins.
- Make the default behavior explicit and safe for production deployments.
- Add tests for an allowed origin, a rejected origin, a missing Origin header, and local-development configuration.
- Document migration steps for deployments that currently depend on the permissive policy.

This note identifies a security concern; it does not itself change runtime behavior. Do not mark the concern resolved until the origin policy is implemented and tested.
