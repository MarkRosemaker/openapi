package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"

	"github.com/MarkRosemaker/jsonutil"
)

// jsonOpts are the options used to read and write a document.
//
// "An AsyncAPI document can be JSON or YAML format. All field names in the specification are
// case sensitive. [...] In order to preserve the ability to round-trip between YAML and JSON
// formats, YAML version 1.2 is RECOMMENDED along with some additional constraints."
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#format
var jsonOpts = json.JoinOptions([]json.Options{
	// the AsyncAPI specification doesn't allow unknown fields in most objects
	// also protect against deleting unknown fields when overwriting later
	json.RejectUnknownMembers(true),
	json.WithMarshalers(json.JoinMarshalers(
		json.MarshalToFunc(jsonutil.URLMarshal),
	)),
	json.WithUnmarshalers(json.JoinUnmarshalers(
		json.UnmarshalFromFunc(jsonutil.URLUnmarshal),
	)),
	jsontext.WithIndent("  "), // indent with two spaces
}...)
