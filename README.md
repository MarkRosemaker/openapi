# OpenAPI
[![Go Reference](https://pkg.go.dev/badge/github.com/MarkRosemaker/openapi.svg)](https://pkg.go.dev/github.com/MarkRosemaker/openapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/MarkRosemaker/openapi)](https://goreportcard.com/report/github.com/MarkRosemaker/openapi)
![Code Coverage](https://img.shields.io/badge/coverage-96.4%25-brightgreen)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
<p align="center">
  <img alt="openapi logo" src=https://avatars3.githubusercontent.com/u/16343502 width=200>
</p>

Package openapi provides a suite of tools for working with OpenAPI specifications, making it easier to parse, format, manipulate, and generate code from these specs. 

Whether you're looking to clean up existing API documentation or integrate API design into your development pipeline, this package is built to streamline your workflow.

## Introduction

The primary goals of this package are:

- **Parsing** OpenAPI specifications into a structured format.
- **Formatting** the parsed specifications, including sorting maps and merging duplicate content.
- **Adding information programmatically** to the specifications.
- **Marshalling** the modified specifications back into their original format.
- **Utilizing** the parsed specification for code generation.

# Features

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

    if err := doc.WriteToFile("path/to/openapi.json"); err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }
}
```

## Additional Information

- [**Go Reference**](https://pkg.go.dev/github.com/MarkRosemaker/openapi): The Go reference documentation for the errpath package.
- [**Go Report Card**](https://goreportcard.com/report/github.com/MarkRosemaker/openapi): Check the code quality report.

## Contributing

If you have any contributions to make, please submit a pull request or open an issue on the [GitHub repository](https://github.com/MarkRosemaker/openapi).

## License

This project is licensed under the [Apache 2.0 License](./LICENSE).