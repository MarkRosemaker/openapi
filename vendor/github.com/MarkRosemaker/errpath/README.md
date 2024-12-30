# Error Path
[![Go Reference](https://pkg.go.dev/badge/github.com/MarkRosemaker/errpath.svg)](https://pkg.go.dev/github.com/MarkRosemaker/errpath)
[![Go Report Card](https://goreportcard.com/badge/github.com/MarkRosemaker/errpath)](https://goreportcard.com/report/github.com/MarkRosemaker/errpath)
![Code Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
<p align="center">
  <img alt="errpath logo: golang gopher determinedly walking through a blue and red maze" src=logo.jpg width=300>
</p>

Package errpath provides utilities for creating and managing detailed error paths.
It allows users to construct error messages that include the full path to the error,
which can be particularly useful when traversing complex data structures such as JSON
or YAML files.

Example for an error in an OpenAPI: `components.schemas["Pet"].allOf[0]: invalid schema`

The package defines several error types that can be used to represent different kinds
of errors, such as missing required values, invalid values, and errors occurring at
specific fields, indices, or keys within a data structure. These error types implement
a chaining mechanism that builds a detailed error path.

## Installation

To use this package, import it as follows:

```go
import "github.com/MarkRosemaker/errpath"
```

## Creating Errors

There are several types of errors you can create with this package:

### ErrRequired

Signals that a required value is missing.

```go
	err := &errpath.ErrRequired{}
```

### ErrInvalid

Signals that a value is invalid. You can optionally provide valid values and an explanatory message.

```go
	err := &errpath.ErrInvalid[string]{
	    Value:   "invalid_value",
	    Enum:    []string{"valid1", "valid2"},
	    Message: "must be one of the valid values",
	}
```

### ErrField

Represents an error that occurred in a specific field.

```go
	err := &errpath.ErrField{
	    Field: "fieldName",
	    Err:   &errpath.ErrRequired{},
	}
```

### ErrIndex

Represents an error that occurred at a specific index in a slice.

```go
	err := &errpath.ErrIndex{
	    Index: 3,
	    Err:   &errpath.ErrInvalid[int]{Value: 42},
	}
```

### ErrKey

Represents an error that occurred at a specific key in a map.

```go
	err := &errpath.ErrKey{
	    Key: "keyName",
	    Err: &errpath.ErrRequired{},
	}
```

## Error Chaining

Errors can be nested to form detailed error paths. For example:

```go
	err := &errpath.ErrField{
	    Field: "foo",
	    Err: &errpath.ErrField{
	        Field: "bar",
	        Err: &errpath.ErrKey{
	            Key: "baz",
	            Err: &errpath.ErrField{
	                Field: "qux",
	                Err: &errpath.ErrIndex{
	                    Index: 3,
	                    Err: &errpath.ErrField{
	                        Field: "quux",
	                        Err: &errpath.ErrInvalid[string]{
	                            Value: "corge",
	                        },
	                    },
	                },
	            },
	        },
	    },
	}
```

This will produce an error message like:

```
	foo.bar["baz"].qux[3].quux ("corge") is invalid
```

## Additional Information

- [**Go Reference**](https://pkg.go.dev/github.com/MarkRosemaker/errpath): The Go reference documentation for the errpath package.
- [**Go Report Card**](https://goreportcard.com/report/github.com/MarkRosemaker/errpath): Check the code quality report.

## Contributing

If you have any contributions to make, please submit a pull request or open an issue on the [GitHub repository](https://github.com/MarkRosemaker/errpath).

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
