# YAML Marshalling and Unmarshalling

This repository provides functionality for YAML marshalling and unmarshalling similar to [JSON v2](https://github.com/go-json-experiment/json). It allows users to use the same JSON options they would use for JSON v2 to encode and decode YAML.

## Features

- Marshal Go structs to YAML
- Unmarshal YAML to Go structs
- Support for JSON v2 options

## Installation

To install the package, use the following command:

```sh
go get github.com/MarkRosemaker/yaml
```

## Usage

Here is a basic example of how to use the package:

```go
package main

import (
	"fmt"
	"github.com/MarkRosemaker/yaml"
)

type Example struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	ex := Example{Name: "John Doe", Age: 30}

	// Marshal to YAML
	yamlData, err := yaml.Marshal(ex)
	if (err != nil) {
		fmt.Println("Error marshalling to YAML:", err)
		return
	}
	fmt.Println("YAML Data:", string(yamlData))

	// Unmarshal from YAML
	var ex2 Example
	err = yaml.Unmarshal(yamlData, &ex2)
	if (err != nil) {
		fmt.Println("Error unmarshalling from YAML:", err)
		return
	}
	fmt.Println("Unmarshalled Struct:", ex2)
}
```

## JSON v2 Options

This package supports the same options available in JSON v2 for encoding and decoding. You can use json struct tags to customize the YAML output just as you would with JSON.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.
