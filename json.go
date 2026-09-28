package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"

	"github.com/MarkRosemaker/jsonutil"
)

var jsonOpts = json.JoinOptions([]json.Options{
	// unevaluatedProperties is set to false in most objects according to the OpenAPI specification
	// also protect against deleting unknown fields when overwriting later
	json.RejectUnknownMembers(true),
	json.WithMarshalers(json.JoinMarshalers(
		json.MarshalToFunc(jsonutil.URLMarshal),
		json.MarshalToFunc(patternMarshal),
	)),
	json.WithUnmarshalers(json.JoinUnmarshalers(
		json.UnmarshalFromFunc(jsonutil.URLUnmarshal),
		json.UnmarshalFromFunc(patternUnmarshal),
	)),
	jsontext.WithIndent("  "), // indent with two spaces
}...)
