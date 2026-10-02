# Legacy Stargate UI source snapshot

This directory contains a standalone source snapshot copied from the stable `Stargate-ui` working tree.

## Purpose

The snapshot exists so the new `Stargate-endpoint` project can reuse and progressively extract mature protocol, backend, Xray, subscription, and deployment code without modifying the original `Stargate-ui` repository.

## Isolation rules

- Do not edit files in `Stargate-ui` through this project.
- Do not add a Git remote from this directory back to the original repository.
- The nested `go.mod` intentionally keeps the legacy source buildable as its own Go module.
- Root `Stargate-endpoint` builds/tests do not include this nested module.
- New Endpoint runtime code should be added outside this snapshot first, then mature legacy components can be migrated deliberately.

## Source origin

Source origin: `github.com/Arshak888/Stargate-ui`

The copy was taken from the local stable working tree at the time of migration. It is a snapshot, not a live submodule or Git history link.

## Migration strategy

1. Keep the snapshot unchanged while the Endpoint architecture is established.
2. Identify reusable protocol/backend components.
3. Extract or adapt components into Endpoint-owned packages.
4. Add tests around each migrated component.
5. Keep the original `Stargate-ui` repository independently stable.
