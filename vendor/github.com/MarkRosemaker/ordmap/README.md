# Ordered Map
[![Go Reference](https://pkg.go.dev/badge/github.com/MarkRosemaker/ordmap.svg)](https://pkg.go.dev/github.com/MarkRosemaker/ordmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/MarkRosemaker/ordmap)](https://goreportcard.com/report/github.com/MarkRosemaker/ordmap)
![Code Coverage](https://img.shields.io/badge/coverage-98.5%25-brightgreen)
[![License: **MIT**](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
<p align="center">
  <img alt="ordmap logo: a gopher holding a map, surrounded by keys" src=logo.jpg width=300>
</p>

`ordmap` is a Go package that provides a generic ordered map implementation, primarily designed for [JSON v2](https://pkg.go.dev/encoding/json/v2) marshalling and unmarshalling.

An ordered map maintains the order of keys based on insertion, allowing you to iterate over the map in the order in which entries were added. This can be particularly useful for applications where the order of elements is important, such as in JSON serialization or when maintaining the sequence of operations.

## Features

- **Seamless JSON v2 Integration:** Directly integrates with the [JSON v2](https://pkg.go.dev/encoding/json/v2) library for efficient and order-preserving marshalling and unmarshalling.
- **Custom Ordered Maps:** Provides robust helper functions to easily define your own custom ordered maps with minimal boilerplate code.
- **Pre-Defined Ordered Map Alias:** Simplifies usage by offering a pre-defined ordered map type that can be conveniently aliased for specific key and value types.
- **Efficient Ordered Operations:** Ensures efficient insertion, retrieval, and iteration while maintaining the order of elements, making it ideal for use cases where order matters.
- **Ordered Iteration:** Leverages the `ByIndex` method to iterate over the map in an ordered manner based on the insertion sequence.

## Installation

To install the library, use the following command:

```shell
go get github.com/MarkRosemaker/ordmap
```

## Usage

### Custom Ordered Map

To create your own custom ordered map, you can utilize helper functions to define its methods:

```go
package main

import (
	"encoding/json/v2"
	"encoding/json/jsontext"
	"iter"

	"github.com/MarkRosemaker/ordmap"
)

type MyOrderedMap map[string]*ValueWithIndex

type ValueWithIndex struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`

	idx int // to order a map of this type
}

func getIndex(v *ValueWithIndex) int                    { return v.idx }
func setIndex(v *ValueWithIndex, i int) *ValueWithIndex { v.idx = i; return v }

// ByIndex returns a sequence of key-value pairs ordered by index.
func (om MyOrderedMap) ByIndex() iter.Seq2[string, *ValueWithIndex] {
	return ordmap.ByIndex(om, getIndex)
}

// Sort sorts the map by key and sets the indices accordingly.
func (om MyOrderedMap) Sort() {
	ordmap.Sort(om, setIndex)
}

// Set sets a value in the map, adding it at the end of the order.
func (om *MyOrderedMap) Set(key string, v *ValueWithIndex) {
	ordmap.Set(om, key, v, getIndex, setIndex)
}

// MarshalJSONTo marshals the key-value pairs in order.
func (om *MyOrderedMap) MarshalJSONTo(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONTo(om, enc, opts)
}

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (om *MyOrderedMap) UnmarshalJSONFrom(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONFrom(om, dec, opts, setIndex)
}
```

If you prefer the map values to be non-pointer types, you can adjust the implementation as follows:

```go
type MyOrderedMap map[string]ValueWithIndex

func getIndex(v ValueWithIndex) int                   { return v.idx }
func setIndex(v ValueWithIndex, i int) ValueWithIndex { v.idx = i; return v }

func (om MyOrderedMap) ByIndex() iter.Seq2[string, ValueWithIndex] {
	return ordmap.ByIndex(om, getIndex)
}

func (om *MyOrderedMap) Set(key string, v ValueWithIndex) {
	ordmap.Set(om, key, v, getIndex, setIndex)
}
```

### Using The Pre-Defined Ordered Map

For simplicity, an ordered map type is already defined for you. You only need to specify the key and value types:

```go
package main

import (
	"github.com/MarkRosemaker/ordmap"
)

type MyOrderedMap = ordmap.OrderedMap[string, *MyValue]

type MyValue struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}
```

## Contributing

If you have any contributions to make, please submit a pull request or open an issue on the [GitHub repository](https://github.com/MarkRosemaker/ordmap).

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
