# Stargate Endpoint

Stargate Endpoint is the standalone server-side Stargate node.

It is intentionally independent from the legacy `Stargate-ui` repository. A single Endpoint instance runs on one VPN/server location and connects to `Stargate-endpoint-manager` for fleet registration, health reporting and future declarative remote control.

## Architecture

- **Stargate Endpoint**: local data plane and runtime authority.
- **Stargate Endpoint Manager**: central control plane.
- **Endpoint ↔ Manager**: server-to-server protocol only.
- Browser panel sessions/cookies are not used for Manager authentication.
- Endpoint enrollment uses a short-lived one-time token.
- After enrollment, the Endpoint receives a dedicated credential for future communication.

## Current phase

Phase 1 establishes:

1. standalone Endpoint binary
2. Manager enrollment bootstrap
3. persistent local Endpoint identity/configuration
4. authenticated heartbeat
5. capability reporting
6. clean separation from the old Stargate UI codebase

The local runtime adapters and protocol implementations will be added behind explicit interfaces instead of coupling the Endpoint to the old panel repository.

## Development

```bash
go run ./cmd/stargate-endpoint
```

Default listen address:

```
127.0.0.1:8090
```

Configuration can be supplied with environment variables:

- `STARGATE_MANAGER_URL`
- `STARGATE_ENROLLMENT_ID`
- `STARGATE_ENROLLMENT_TOKEN`
- `STARGATE_ENDPOINT_CONFIG`
- `STARGATE_ENDPOINT_NAME`
- `STARGATE_REGION`
- `STARGATE_COUNTRY`
- `STARGATE_CITY`
