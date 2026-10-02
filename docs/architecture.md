# Stargate Endpoint Architecture

Stargate Endpoint is a standalone node product. It does not depend on the legacy Stargate UI repository.

## Responsibilities

The Endpoint owns:

- local protocol daemons
- local Xray/runtime configuration
- nftables and kernel data-plane state
- local accounting and session state
- execution of Manager commands
- local persistence needed for safe offline operation

The Manager owns:

- endpoint inventory
- global accounts and endpoint memberships
- desired state
- fleet telemetry/history
- jobs and audit logs
- reseller permissions
- orchestration

## Design rule

The Manager must never become the VPN traffic path. It sends commands and receives telemetry; user traffic remains on Endpoint servers.

## Lifecycle

1. Endpoint is installed.
2. Manager creates a one-time enrollment.
3. Endpoint bootstraps and receives an endpoint credential.
4. Endpoint stores its identity locally.
5. Endpoint sends authenticated heartbeats.
6. Later protocol adapters expose capabilities and declarative operations.
