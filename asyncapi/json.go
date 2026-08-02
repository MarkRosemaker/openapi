package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"

	"github.com/MarkRosemaker/jsonutil"
)

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
