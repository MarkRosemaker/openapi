# Agent Notes for MarkRosemaker/openapi

## Workflow Preferences

- **Start from latest master**: Always `git fetch origin master && git rebase origin/master` (or create the branch from `origin/master`) before starting work. Avoid unnecessary merge commits.
- **One focused PR per feature**: Keep changes small and scoped. Don't bundle unrelated cleanup into a feature PR.
- **No `gh` CLI**: GitHub interactions go through the MCP GitHub tools (`mcp__github__*`). Use `ToolSearch` to load their schemas.
- **Skip YAML handling**: When adding JSON-level features, focus on JSON only. Do not add corresponding YAML plumbing unless explicitly requested.

## Building and Testing

```bash
# This repo requires Go 1.26.3 with the jsonv2 experiment flag
GOEXPERIMENT=jsonv2 go build ./...
GOEXPERIMENT=jsonv2 go test ./...
```

All CI and local testing uses `GOEXPERIMENT=jsonv2`. Never run `go test` without it — the build will fail.

## Key Architecture

- **`encoding/json/v2`** (`encoding/json/jsontext`) — not stable stdlib yet; gated behind `GOEXPERIMENT=jsonv2`. Vendor dir at `vendor/`.
- **`refOrValue[T, O]`** (`ref.go`) — generic type backing all `*Ref` aliases (SchemaRef, HeaderRef, etc.). Implements custom `UnmarshalJSONFrom` / `MarshalJSONTo`. Uses `json.RejectUnknownMembers(false)` when probing for a `$ref` so that extra properties alongside `$ref` are silently discarded (as required by the Reference Object spec) rather than causing a parse error.
- **`loader`** (`loader.go`) — two-pass load: unmarshal → `collectResolveRefs` (collect component schemas, then resolve all `$ref`s).
- **`Schema.Enum`** is `[]any` — JSON numbers decode to `float64`, so integer enum values arrive as `float64` and are validated by `enumValueMatchesType`.

## OAS 3.1 / JSON Schema 2020-12 Notes

- **Reference Object**: OAS 3.1.0 defines a Reference Object with exactly three fields: `$ref`, `summary`, `description`. The spec states "This object cannot be extended with additional properties and any properties added SHALL be ignored." Our `UnmarshalJSONFrom` passes `json.RejectUnknownMembers(false)` when probing for a reference so that extra properties (e.g. `default`) do not prevent `$ref` from being detected — they are silently discarded. Note: this means `default` alongside `$ref` is lost on parse, which is correct per spec.
- **Enum type validation**: Each enum value must match the schema's declared type (`enumValueMatchesType` in `schema.go`).
