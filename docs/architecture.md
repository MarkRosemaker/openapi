# Architecture

- **`encoding/json/v2`** (`encoding/json/jsontext`) — stable in Go 1.27 standard library.
- **`refOrValue[T, O]`** (`ref.go`) — generic type backing all `*Ref` aliases (SchemaRef, HeaderRef, etc.). Implements custom `UnmarshalJSONFrom` / `MarshalJSONTo`. Probes for `$ref` by attempting to unmarshal as `Reference`; falls back to the value type if `$ref` is absent.
- **`loader`** (`loader.go`) — two-pass load: unmarshal → `collectResolveRefs` (collect component schemas, then resolve all `$ref`s).
- **`Schema.Enum`** is `[]jsontext.Value`, **`Schema.Const`** and **`Schema.Default`** are `jsontext.Value` — raw JSON is preserved exactly as written. Kind-based validation (`enumKindMatchesType`, `isJSONInteger`) checks types without decoding to Go values.
- **`Schema` (un)marshalling** (`schema_json.go`) — decodes through `schemaJSON`, whose shallower `type` field overrides the embedded one: a `type` of `[X, "null"]` becomes `Type` X with `Nullable` set. An array of several non-null types fails to decode.
- **`Schema.Type` is optional** (JSON Schema 2020-12): without it a schema constrains no type, and the empty schema accepts any value. Keywords tied to one type (`properties`, `items`, `format`, `pattern`, …) still require it in `Validate`.
- **`Schema.AdditionalProperties`** is `*AdditionalProperties`: a schema for the extra properties' values, or a bare boolean. A nil pointer means the keyword is absent, which is not the same as `true` downstream.
- **Patterns** (`pattern.go`) are compiled with Go's regexp (RE2). `jsonOpts` translates ECMA-262 escapes RE2 spells differently on load, and back on write, so a pattern round-trips as written. Marshalling without `jsonOpts` writes the RE2 form.
