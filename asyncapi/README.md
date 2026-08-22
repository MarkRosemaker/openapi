<div align="center" id=badges>

[![Go Reference](https://pkg.go.dev/badge/github.com/MarkRosemaker/asyncapi.svg)](https://pkg.go.dev/github.com/MarkRosemaker/asyncapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/MarkRosemaker/asyncapi)](https://goreportcard.com/report/github.com/MarkRosemaker/asyncapi)
![Code Coverage](https://img.shields.io/badge/coverage-86.2%25-brightgreen)
[![License: Apache](https://img.shields.io/badge/License-Apache-yellow.svg)](./LICENSE)

</div>

<h3 align="center">
  Transform and master your event-driven API specs with ease.
</h3>

Package asyncapi provides a suite of tools for working with [AsyncAPI](https://www.asyncapi.com) specifications, making it easier to parse, format, manipulate, and generate code from these specs.

It is the counterpart of [MarkRosemaker/openapi](https://github.com/MarkRosemaker/openapi) for event-driven APIs and follows the same design, so both packages can be used side by side.

## Introduction

The primary goals of this package are:

- **Parsing** AsyncAPI specifications into a structured format.
- **Validating** the specifications strictly against the rules of the specification.
- **Formatting** the parsed specifications, including sorting maps and merging duplicate content.
- **Adding information programmatically** to the specifications.
- **Marshalling** the modified specifications back into their original format.
- **Utilizing** the parsed specification for code generation.

## Features

- **Comprehensive parsing** of [AsyncAPI 3.1.0](https://www.asyncapi.com/docs/reference/specification/v3.1.0) specifications, in JSON as well as in YAML.
- **Strict validation** against the rules of the specification: required fields, enumerations, key and address patterns, absolute URLs, [runtime expressions](https://www.asyncapi.com/docs/reference/specification/v3.1.0#runtimeExpression), and the rules that span several objects, e.g. that the messages of an [operation](https://www.asyncapi.com/docs/reference/specification/v3.1.0#operationObject) "MUST contain a subset of the messages defined in the channel referenced in this operation". Every error names the exact location of the problem, e.g. `channels["userSignedup"].messages["userSignedUp"].contentType: mime: expected slash after first token`.
- **Marshalling** back to JSON and to YAML, to a file, to a writer or to a byte slice.
- **Reference resolution** of every referencable object, including references that point to other references, e.g. an operation that refers to a message of a channel which in turn refers to a message of the [components object](https://www.asyncapi.com/docs/reference/specification/v3.1.0#componentsObject).
- **Order preservation**: maps keep the order in which their keys were defined, so writing a specification back doesn't reshuffle it.
- **Multi format schemas**: schemas in other formats (Avro, Protobuf, RAML, ...) are kept as they are, [AsyncAPI schemas](https://www.asyncapi.com/docs/reference/specification/v3.1.0#schemaObject) are parsed, including boolean schemas and multiple types.
- **Bindings** of all [20 protocols](https://www.asyncapi.com/docs/reference/specification/v3.1.0#serverBindingsObject) are preserved as they were given, so nothing is lost when a specification is written back.
- **Documented in line with the specification**: every object, field and rule quotes the [official documentation](https://www.asyncapi.com/docs/reference/specification/v3.1.0) and links to the section it comes from.

## Usage

```go
package main

import (
    "fmt"

    "github.com/MarkRosemaker/asyncapi"
)

func main() {
    doc, err := asyncapi.LoadFromFile("path/to/asyncapi.json") // or asyncapi.yaml
    if err != nil {
        fmt.Println("Error parsing spec:", err)
        return
    }

    if err := doc.Validate(); err != nil {
        fmt.Println("Error validating spec:", err)
        return
    }

    // sort the keys of the servers, channels, operations and components in alphabetical order
    doc.SortMaps()

    // write an improved version of your spec, as JSON or as YAML
    if err := doc.WriteToFile("path/to/asyncapi.json"); err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }
}
```

## Additional Information

- [**Go Reference**](https://pkg.go.dev/github.com/MarkRosemaker/asyncapi): The Go reference documentation for the asyncapi package.
- [**Go Report Card**](https://goreportcard.com/report/github.com/MarkRosemaker/asyncapi): Check the code quality report.

## Contributing

If you have any contributions to make, please submit a pull request or open an issue on the [GitHub repository](https://github.com/MarkRosemaker/asyncapi).

## License

This project is licensed under the [Apache 2.0 License](./LICENSE).
