package openapi

import (
	"errors"
	"strings"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// Extensions represents additional fields that can be added to OpenAPI objects.
//
// While the OpenAPI Specification tries to accommodate most use cases, additional data can be added to extend the specification at certain points.
//
// The field name MUST begin with x-, for example, x-internal-id. Field names beginning x-oai- and x-oas- are reserved for uses defined by the OpenAPI Initiative. The value can be null, a primitive, an array or an object.
// ([Specification])
//
// It is here an alias of jsontext.Value to allow inlining within structs, enabling
// seamless marshalling and unmarshalling. Using jsontext.Value preserves the order
// of fields, preventing unnecessary changes when parsing and writing OpenAPI
// specifications. Although a map could be used, it doesn't maintain the order,
// leading to potential inconsistencies in the output. Custom marshalling for an
// inlined object is not possible, which prevents the use of an ordered map.
//
// Note: For convenience, certain common extensions are implemented as fields
// directly within the respective structs.
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#specification-extensions
type Extensions = jsontext.Value

// ErrUnknownField is returned when a field is not recognized and also doesn't have a "x-" prefix signifying it is an extension.
var ErrUnknownField = errors.New(`unknown field or extension without "x-" prefix`)

func validateExtensions(ext Extensions) error {
	if len(ext) == 0 {
		return nil
	}

	m := map[string]any{}
	if err := json.Unmarshal(ext, &m); err != nil {
		return err
	}

	for k := range m {
		if !strings.HasPrefix(k, "x-") {
			return &ErrField{Field: k, Err: ErrUnknownField}
		}
	}

	return nil
}
