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
- **`refOrValue[T, O]`** (`ref.go`) — generic type backing all `*Ref` aliases (SchemaRef, HeaderRef, etc.). Implements custom `UnmarshalJSONFrom` / `MarshalJSONTo`. Probes for `$ref` by attempting to unmarshal as `Reference`; falls back to the value type if `$ref` is absent.
- **`loader`** (`loader.go`) — two-pass load: unmarshal → `collectResolveRefs` (collect component schemas, then resolve all `$ref`s).
- **`Schema.Enum`** is `[]jsontext.Value` and **`Schema.Default`** is `jsontext.Value` — raw JSON is preserved exactly as written. Kind-based validation (`enumKindMatchesType`, `isJSONInteger`) checks types without decoding to Go values.

## OAS 3.1 / JSON Schema 2020-12 Notes

- Any implementation needs to be justified by the [official OpenAPI Specification](https://spec.openapis.org/oas/v3.1.0), reference links to it in relevant comments