package openapi

import (
	"github.com/MarkRosemaker/jsonutil"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

var jsonOpts = json.JoinOptions([]json.Options{
	// unevaluatedProperties is set to false in most objects according to the OpenAPI specification
	// also protect against deleting unknown fields when overwriting later
	json.RejectUnknownMembers(true),
	json.WithMarshalers(json.NewMarshalers(
		json.MarshalToFunc(jsonutil.URLMarshal))),
	json.WithUnmarshalers(json.NewUnmarshalers(
		json.UnmarshalFromFunc(jsonutil.URLUnmarshal))),
	jsontext.WithIndent("  "), // indent with two spaces
}...)
