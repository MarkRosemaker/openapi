# Types
[![Go Reference](https://pkg.go.dev/badge/github.com/go-api-libs/types.svg)](https://pkg.go.dev/github.com/go-api-libs/types)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-api-libs/types)](https://goreportcard.com/report/github.com/go-api-libs/types)
![Code Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

This library offers a collection of commonly used data types in API contexts, ensuring correct marshalling, unmarshalling, and validation. By standardizing these types across all our API libraries, we aim to reduce redundancy, increase consistency, and enhance type safety.

## Features

- **Standardized Types**: Tailored for API data interchange
- **Validation**: Built-in validation rules to ensure data integrity before processing or transmission.
- **JSON Marshalling/Unmarshalling**: Custom JSON handling for each type to handle special cases (e.g., ensuring email addresses are formatted correctly).
- **Used Across go-api-libs**: These types are employed by all API libraries within the [go-api-libs](https://github.com/go-api-libs/) organization for consistency.

## Installation

To install the library, use the following command:

```shell
go get github.com/go-api-libs/types
```

## Usage

When using any of the API libraries from the [go-api-libs](https://github.com/go-api-libs/) organization, you can trust that types in successful API responses are valid.

Beyond that, here's how you can use some types from this library:

```go
package main

import (
	"fmt"

	"github.com/go-api-libs/types"
)

func main() {
	// Using Email type
	email := types.Email("user@example.com")
	if err := email.Validate(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Valid email address")
	}
}
```

## Contributing

If you have any contributions to make, please submit a pull request or open an issue on the [GitHub repository](https://github.com/go-api-libs/types).

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
