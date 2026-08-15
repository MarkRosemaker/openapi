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
- **`refOrValue[T, O]`** (`ref.go`) — generic type backing all `*Ref` aliases (SchemaRef, HeaderRef, etc.). Implements custom `UnmarshalJSONFrom` / `MarshalJSONTo`.
- **`loader`** (`loader.go`) — two-pass load: unmarshal → `collectResolveRefs` (collect component schemas, then resolve all `$ref`s).
- **`emptier` interface** (`ref.go`) — implemented by `*Schema` to detect sibling keywords alongside `$ref` during unmarshal.
- **`Schema.Enum`** is `[]any` — JSON numbers decode to `float64`, so integer enum values arrive as `float64` and are validated by `enumValueMatchesType`.

## OAS 3.1 / JSON Schema 2020-12 Notes

- **`$ref` + sibling keywords**: In OAS 3.1+ (which uses JSON Schema 2020-12), a `$ref` in a Schema Object may carry sibling keywords (`description`, `default`, etc.). These apply alongside the referenced schema. Our implementation merges the referenced schema's fields into the sibling schema (sibling values win), then clears `s.Ref` so the result validates as a plain schema.
- **Reference Object vs Schema `$ref`**: The Reference Object (used for parameters, responses, headers) allows only `$ref`, `summary`, `description` as siblings. Schema `$ref` (JSON Schema context) allows any keyword. Our `emptier` interface ensures sibling capture only happens for `*Schema`, not other ref types.
- **Enum type validation**: Each enum value must match the schema's declared type (`enumValueMatchesType` in `schema.go`).
