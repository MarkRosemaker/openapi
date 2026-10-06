# The openapi family

These repositories exist to document an API that is not properly documented,
and then to use it: record its traffic, get an OpenAPI specification that says
what it actually does, and generate a Go library to call it with. Where that
library fails to decode a response, recording the call and enriching again
closes the gap.

| Repository | Builds on |
|---|---|
| openapi | |
| openapi-edit | openapi |
| openapi-compare | openapi |
| openapi-merge | openapi |
| openapi-enrich | openapi, edit, merge |
| openapi-flatten | openapi |
| openapi-compress | openapi, compare, edit, enrich, merge |
| openapi-codegen | openapi, compare, compress, edit, enrich, flatten |

Work goes down the table: while a repository higher up has a pull request
open, finish that first. A change that must reach the repositories below it
waits until it is merged and they are bumped.

This repository is the top of the chain: every other one builds on its `*openapi.Document`, so a change here reaches all of them.

## Tests

Each repository tests itself against one golden file in `testdata/`, or a few,
crafted to hold every case it handles once, as small as that allows. Nothing
regenerates a golden file and no repository copies another's: a change in what
comes out is edited in by hand, so it is read before it is accepted. A case is
a test of its own only where it needs a comment to explain it. Every feature
has one regression guard, not several.
