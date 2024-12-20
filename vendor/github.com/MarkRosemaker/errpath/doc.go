/*
Package errpath provides utilities for creating and managing detailed error paths.
It allows users to construct error messages that include the full path to the error,
which can be particularly useful when traversing complex data structures such as JSON
or YAML files.

The package defines several error types that can be used to represent different kinds
of errors, such as missing required values, invalid values, and errors occurring at
specific fields, indices, or keys within a data structure. These error types implement
a chaining mechanism that builds a detailed error path.

Usage:

To use this package, import it as follows:

	import "github.com/MarkRosemaker/errpath"

Creating Errors:

There are several types of errors you can create with this package:

1. ErrRequired: Signals that a required value is missing.

	err := &errpath.ErrRequired{}

2. ErrInvalid: Signals that a value is invalid. You can optionally provide valid values and an explanatory message.

	err := &errpath.ErrInvalid[string]{
	    Value:   "invalid_value",
	    Enum:    []string{"valid1", "valid2"},
	    Message: "must be one of the valid values",
	}

3. ErrField: Represents an error that occurred in a specific field.

	err := &errpath.ErrField{
	    Field: "fieldName",
	    Err:   &errpath.ErrRequired{},
	}

4. ErrIndex: Represents an error that occurred at a specific index in a slice.

	err := &errpath.ErrIndex{
	    Index: 3,
	    Err:   &errpath.ErrInvalid[int]{Value: 42},
	}

5. ErrKey: Represents an error that occurred at a specific key in a map.

	err := &errpath.ErrKey{
	    Key: "keyName",
	    Err: &errpath.ErrRequired{},
	}

Error Chaining:

Errors can be nested to form detailed error paths. For example:

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

This will produce an error message like:

	foo.bar["baz"].qux[3].quux ("corge") is invalid
*/
package errpath
