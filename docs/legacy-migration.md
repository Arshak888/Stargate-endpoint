# Legacy Stargate-ui migration map

The repository now contains an isolated source snapshot at legacy/stargate-ui/.
The original Stargate-ui repository remains a separate stable project and must not be
modified as part of Endpoint development.

## What the snapshot gives Endpoint

The snapshot contains mature implementation material for:

- Xray configuration, process supervision, API traffic accounting and inbound handling.
- Native Xray protocols including VMess, VLESS, Trojan, Shadowsocks, Mixed, AnyTLS, TUIC and NaiveProxy.
- Server-side VPN protocols including L2TP, PPTP, OpenVPN, OpenConnect, SSTP, IKEv2, WireGuard C, AmneziaWG, GRE, MTProto and SSH.
- Protocol-specific defaults and validation.
- nftables accounting, routing and policy enforcement.
- VPN address allocation and per-client accounting.
- Database models and the existing Account / AccountInbound membership projection.
- Backend daemon bundling for protocols that need external daemons.
- Extensive unit and integration tests around the mature implementation.

## What should NOT be copied wholesale

Endpoint is not a renamed web panel. The following remain panel concerns unless a
specific piece is extracted behind an Endpoint-owned interface:

- Gin web UI and templates.
- Browser sessions and panel authentication.
- Reseller UI and panel-only permissions.
- Telegram/WhatsApp presentation and notification jobs.
- Panel update/install/uninstall flows.
- Browser form defaults and UI-specific JSON projection.
- Panel startup orchestration.

## Extraction rule

When a legacy component is needed:

1. Identify the smallest reusable implementation boundary.
2. Define an Endpoint-owned interface around the behavior.
3. Add tests at the new boundary.
4. Copy or refactor only the required implementation into Endpoint-owned packages.
5. Keep the legacy snapshot unchanged until the extracted implementation has parity tests.
6. Do not make Endpoint runtime behavior depend on panel HTTP handlers.

## Current protocol status

The legacy snapshot demonstrates mature protocol support, but Endpoint MUST NOT advertise
a protocol capability merely because its source exists in the snapshot.

The Endpoint currently advertises only capabilities that are operational in the new
Endpoint runtime:

- endpoint.telemetry.heartbeat
- endpoint.state.local
- legacy.stargate-ui.source

Protocol capabilities will be added individually as Endpoint adapters become operational.

## First extraction targets

The first useful extraction boundaries are:

1. Runtime/system telemetry: CPU, memory, disk and uptime, plus an active session provider interface.
2. Xray runtime adapter: process lifecycle, config generation, traffic and session inspection.
3. Protocol adapter registry: one adapter per protocol family with explicit capability discovery.
4. Account projection: reuse the mature Account / AccountInbound ideas from web/service/accountproject.go while keeping Global Account identity in the Manager.
5. Accounting: extract nftables counters and session reconciliation behind an Endpoint-owned interface.
6. External daemons: reuse mature backend bundle/build mechanisms where appropriate, without carrying panel installation logic into the Endpoint command path.

The goal is not to rewrite 860 files. The goal is to turn the mature parts of those files
into clean Endpoint-owned runtime components while leaving the stable panel intact.