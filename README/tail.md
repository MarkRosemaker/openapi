## The openapi family

This module is the foundation of a family of composable tools. Together they
document an API that is not properly documented and then make it easy to use:
record its traffic, get a specification that says what it actually does, and
generate a Go library to call it with. Every tool operates on the
`*openapi.Document` defined here, so they can also be combined freely.

The table lists them in dependency order: each builds only on those above it.

| Module | Purpose | Builds on |
|---|---|---|
| **openapi** (this module) | Parse, validate, and write OpenAPI 3.x specifications | |
| [openapi-edit](https://github.com/MarkRosemaker/openapi-edit) | Safe structural edits, such as renaming a schema and rewriting every `$ref` to it | openapi |
| [openapi-compare](https://github.com/MarkRosemaker/openapi-compare) | Compare specification objects — exact equality and shape equivalence | openapi |
| [openapi-merge](https://github.com/MarkRosemaker/openapi-merge) | Merge schemas that were inferred independently from different samples | openapi |
| [openapi-enrich](https://github.com/MarkRosemaker/openapi-enrich) | Infer specification content from observed HTTP traffic | openapi, edit, merge |
| [openapi-flatten](https://github.com/MarkRosemaker/openapi-flatten) | Promote inline definitions into named `components` entries | openapi |
| [openapi-compress](https://github.com/MarkRosemaker/openapi-compress) | Deduplicate and merge equivalent component schemas | openapi, compare, edit, enrich, merge |
| [openapi-codegen](https://github.com/MarkRosemaker/openapi-codegen) | Generate Go types, clients, and servers from a specification | openapi, compare, edit, enrich, flatten, compress |

The last four form the pipeline, in that order: `openapi-enrich` records traffic
into a specification, `openapi-flatten` and `openapi-compress` normalize its
structure, and `openapi-codegen` generates the library. Where that library fails
to decode a response, recording the call and enriching again closes the gap.
