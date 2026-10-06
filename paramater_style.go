package openapi

import (
	"slices"

	"github.com/MarkRosemaker/errpath"
)

// ParameterStyle is how a parameter is serialized. In order to support common ways of serializing simple parameters, a set of `style` values are defined.
// Assume a parameter named `color` has one of the following values:
//
//	string -> "blue"
//	array -> ["blue","black","brown"]
//	object -> { "R": 100, "G": 200, "B": 150 }
//
// The following table shows examples of rendering differences for each value.
//
//	| `style`        | `explode` | `empty` | `string`    | `array`                             | `object`                               |
//	|----------------|-----------|---------|-------------|-------------------------------------|----------------------------------------|
//	| matrix         | false     | ;color  | ;color=blue | ;color=blue,black,brown             | ;color=R,100,G,200,B,150               |
//	| matrix         | true      | ;color  | ;color=blue | ;color=blue;color=black;color=brown | ;R=100;G=200;B=150                     |
//	| label          | false     | .       | .blue       | .blue.black.brown                   | .R.100.G.200.B.150                     |
//	| label          | true      | .       | .blue       | .blue.black.brown                   | .R=100.G=200.B=150                     |
//	| form           | false     | color=  | color=blue  | color=blue,black,brown              | color=R,100,G,200,B,150                |
//	| form           | true      | color=  | color=blue  | color=blue&color=black&color=brown  | R=100&G=200&B=150                      |
//	| simple         | false     | n/a     | blue        | blue,black,brown                    | R,100,G,200,B,150                      |
//	| simple         | true      | n/a     | blue        | blue,black,brown                    | R=100,G=200,B=150                      |
//	| spaceDelimited | false     | n/a     | n/a         | blue%20black%20brown                | R%20100%20G%20200%20B%20150            |
//	| pipeDelimited  | false     | n/a     | n/a         | blue|black|brown                    | R|100|G|200|B|150                      |
//	| deepObject     | true      | n/a     | n/a         | n/a                                 | color[R]=100&color[G]=200&color[B]=150 |
type ParameterStyle string

// The styles a parameter can be serialized in.
const (
	ParameterStyleMatrix ParameterStyle = "matrix"
	ParameterStyleLabel  ParameterStyle = "label"
	ParameterStyleForm   ParameterStyle = "form"
	ParameterStyleSimple ParameterStyle = "simple"
	// ParameterStyleSpaceDelimited separates array or object values by spaces. This option replaces `collectionFormat` equal to `ssv` from OpenAPI 2.0.
	ParameterStyleSpaceDelimited ParameterStyle = "spaceDelimited"
	// ParameterStylePipeDelimited separates array or object values by pipes. This option replaces `collectionFormat` equal to `pipes` from OpenAPI 2.0.
	ParameterStylePipeDelimited ParameterStyle = "pipeDelimited"
	// ParameterStyleDeepObject provides a simple way of rendering nested objects using form parameters.
	ParameterStyleDeepObject ParameterStyle = "deepObject"
)

var allParameterStyles = []ParameterStyle{
	ParameterStyleMatrix, ParameterStyleLabel, ParameterStyleForm, ParameterStyleSimple,
	ParameterStyleSpaceDelimited, ParameterStylePipeDelimited, ParameterStyleDeepObject,
}

// Validate returns an error if the ParameterStyle breaks the specification.
func (s ParameterStyle) Validate() error {
	if slices.Contains(allParameterStyles, s) {
		return nil
	}

	return &errpath.ErrInvalid[ParameterStyle]{
		Value: s,
		Enum:  allParameterStyles,
	}
}
