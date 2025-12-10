# JSON v2 Utilities
[![Go Reference](https://pkg.go.dev/badge/github.com/MarkRosemaker/jsonutil.svg)](https://pkg.go.dev/github.com/MarkRosemaker/jsonutil)
[![Go Report Card](https://goreportcard.com/badge/github.com/MarkRosemaker/jsonutil)](https://goreportcard.com/report/github.com/MarkRosemaker/jsonutil)
![Code Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
<p align="center">
  <img alt="jsonutil logo: cute golang gopher in JSON brackets with gear in the background and text saying v2" src=logo.jpg width=300>
</p>

`jsonutil` is a Go package that provides custom [JSON v2](https://pkg.go.dev/encoding/json/v2) marshaling and unmarshaling functions for specific data types.

This package is particularly useful when you need to handle JSON encoding and decoding for types like `time.Duration` and `url.URL` in a customized manner.

## Features

* Custom marshaler and unmarshaler for `url.URL`:
  * `URLMarshal` marshals `url.URL` as a string.
  * `URLUnmarshal` unmarshals `url.URL` from a string.
* Custom marshaler and unmarshaler for `time.Duration`:
  * `DurationMarshalIntSeconds` marshals `time.Duration` as an integer representing seconds.
  * `DurationUnmarshalIntSeconds` unmarshals `time.Duration` from an integer assuming it represents seconds.

## Installation

To install the library, use the following command:

```shell
go get github.com/MarkRosemaker/jsonutil
```

## Usage


### Custom Marshaling and Unmarshaling for `url.URL`

To use the custom marshaler and unmarshaler for `url.URL`, you can import the package and use the provided functions:

```go
package main

import (
	"encoding/json/v2"
	"fmt"
	"net/url"

	"github.com/MarkRosemaker/jsonutil"
)

var jsonOpts = json.JoinOptions(
	json.WithMarshalers(json.NewMarshalers(json.MarshalFuncV2(jsonutil.URLMarshal))),
	json.WithUnmarshalers(json.NewUnmarshalers(json.UnmarshalFuncV2(jsonutil.URLUnmarshal))),
)

type MyStruct struct {
	Foo  string  `json:"foo"`
	Link url.URL `json:"link"`
}

func main() {
	out := &MyStruct{}
	if err := json.Unmarshal([]byte(`{"foo":"bar","link":"https://example.com/"}`), out, jsonOpts); err != nil {
		panic(err)
	}

	if out.Link.String() != "https://example.com/" {
		panic("something went wrong")
	}

	res, err := json.Marshal(MyStruct{
		Foo:  "baz",
		Link: url.URL{Scheme: "http", Host: "example.com", Path: "/some-path"},
	}, jsonOpts)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
	// Output: {"foo":"baz","link":"http://example.com/some-path"}
}
```

### Custom Marshaling and Unmarshaling for `time.Duration`

To use the custom marshaler and unmarshaler for `time.Duration`, you can import the package and use the provided functions:

```go
package main

import (
	"encoding/json/v2"
	"fmt"
	"time"

	"github.com/MarkRosemaker/jsonutil"
)

var jsonOpts = json.JoinOptions(
	json.WithMarshalers(json.MarshalFuncV2(jsonutil.DurationMarshalIntSeconds)),
	json.WithUnmarshalers(json.UnmarshalFuncV2(jsonutil.DurationUnmarshalIntSeconds)),
)

type MyStruct struct {
	Foo      string        `json:"foo"`
	Duration time.Duration `json:"duration"`
}

func main() {
	out := &MyStruct{}
	if err := json.Unmarshal([]byte(`{"foo":"bar","duration":3}`), out, jsonOpts); err != nil {
		panic(err)
	}

	if out.Duration != 3*time.Second {
		panic("something went wrong")
	}

	res, err := json.Marshal(MyStruct{
		Foo:      "baz",
		Duration: time.Minute,
	}, jsonOpts)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
	// Output: {"foo":"baz","duration":60}
}
```

## Contributing

If you have any contributions to make, please submit a pull request or open an issue on the [GitHub repository](https://github.com/MarkRosemaker/jsonutil).

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
