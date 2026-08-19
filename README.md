<div align="center" id=badges>

[![Go Reference](https://pkg.go.dev/badge/github.com/MarkRosemaker/openapi.svg)](https://pkg.go.dev/github.com/MarkRosemaker/openapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/MarkRosemaker/openapi)](https://goreportcard.com/report/github.com/MarkRosemaker/openapi)
![Code Coverage](https://img.shields.io/badge/coverage-94.5%25-brightgreen)
[![License: Apache](https://img.shields.io/badge/License-Apache-yellow.svg)](./LICENSE)

</div>

<p align="center">
  <img alt="OpenAPI Logo" src=openapi-logo.svg width=500>
</p>

<h3 align="center">
  Transform and master your API specs with ease.
</h3>

Package openapi provides a suite of tools for working with OpenAPI specifications, making it easier to parse, format, manipulate, and generate code from these specs.

Whether you're looking to clean up existing API documentation or integrate API design into your development pipeline, this package is built to streamline your workflow.

**This package is currently being utilized to format OpenAPI specifications in the [go-api-libs](https://github.com/go-api-libs) project.**

## Introduction

The primary goals of this package are:

- **Parsing** OpenAPI specifications into a structured format.
- **Formatting** the parsed specifications, including sorting maps and merging duplicate content.
- **Adding information programmatically** to the specifications.
- **Marshalling** the modified specifications back into their original format.
- **Utilizing** the parsed specification for code generation.

This module is deliberately kept focused on representing and validating a
specification. Transformations that not everyone needs — flattening, deduplicating,
enriching, generating code — live in [separate modules](#the-openapi-family) so that
users who only want to parse, validate, and prettify a spec don't pay for them.

## Features

- **Comprehensive parsing** of OpenAPI specifications.
- **Flexible formatting** options to improve readability and consistency.
- **Ability to merge and deduplicate** content within specifications.
- **Programmatic modification** of specifications before marshalling.
- **Code generation capabilities** based on parsed specifications.

## Usage

```go
package main

import (
    "fmt"

    "github.com/MarkRosemaker/openapi"
)

func main() {
    doc, err := openapi.LoadFromFile("path/to/openapi.json") // or openapi.yaml
    if err != nil {
        fmt.Println("Error parsing spec:", err)
        return
    }

    if err := doc.Validate(); err != nil {
        fmt.Println("Error validating spec:", err)
        return
    }

    // sort keys of each component in alphabetical order
    doc.Components.SortMaps()

	// write an improved version of your spec
    if err := doc.WriteToFile("path/to/openapi.json"); err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }
}
```

## The openapi family

This module is the foundation of a family of composable tools. Every one of them
operates on the `*openapi.Document` defined here, so they can be combined freely.

| Module | Purpose |
|---|---|
| **openapi** (this module) | Parse, validate, and write OpenAPI 3.x specifications |
| [openapi-compare](https://github.com/MarkRosemaker/openapi-compare) | Compare specification objects — exact equality and shape equivalence |
| [openapi-edit](https://github.com/MarkRosemaker/openapi-edit) | Safe structural edits, such as renaming a schema and rewriting every `$ref` to it |
| [openapi-flatten](https://github.com/MarkRosemaker/openapi-flatten) | Promote inline definitions into named `components` entries |
| [openapi-compress](https://github.com/MarkRosemaker/openapi-compress) | Deduplicate and merge equivalent component schemas |
| [openapi-merge](https://github.com/MarkRosemaker/openapi-merge) | Merge schemas that were inferred independently from different samples |
| [openapi-enrich](https://github.com/MarkRosemaker/openapi-enrich) | Infer specification content from observed HTTP traffic |
| [openapi-codegen](https://github.com/MarkRosemaker/openapi-codegen) | Generate Go types, clients, and servers from a specification |

A typical pipeline records traffic with `openapi-enrich`, normalizes structure with
`openapi-flatten` and `openapi-compress`, and finally generates code with
`openapi-codegen`.

## Additional Information

- [**Go Reference**](https://pkg.go.dev/github.com/MarkRosemaker/openapi): The Go reference documentation for the openapi package.
- [**Go Report Card**](https://goreportcard.com/report/github.com/MarkRosemaker/openapi): Check the code quality report.

## Contributing

If you have any contributions to make, please submit a pull request or open an issue on the [GitHub repository](https://github.com/MarkRosemaker/openapi).

## License

This project is licensed under the [Apache 2.0 License](./LICENSE).
