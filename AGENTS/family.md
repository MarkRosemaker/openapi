# The openapi family

These repositories exist to document an API that is not properly documented,
and then to use it: record its traffic, get an OpenAPI specification that says
what it actually does, and generate a Go library to call it with. Where that
library fails to decode a response, recording the call and enriching again
closes the gap.

| Repository | Builds on | Test data from |
|---|---|---|
| openapi | | |
| openapi-edit | openapi | |
| openapi-compare | openapi | |
| openapi-merge | openapi | |
| openapi-enrich | openapi, edit, merge | |
| openapi-flatten | openapi | enrich's golden files |
| openapi-compress | openapi, compare, edit, enrich, merge | flatten's golden files, enrich's interactions |
| openapi-codegen | openapi, compare, compress, edit, enrich, flatten | compress's golden files, enrich's interactions |

Work goes down the table: while a repository higher up has a pull request
open, finish that first. A change that must reach the repositories below it
waits until it is merged and they are bumped.

This repository is the top of the chain: every other one builds on its `*openapi.Document`, so a change here reaches all of them.
