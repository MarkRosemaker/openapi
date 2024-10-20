package json

import (
	"github.com/MarkRosemaker/openapi/internal/jsonutil"
	"github.com/go-json-experiment/json"
)

var Options = json.JoinOptions([]json.Options{
	// unevaluatedProperties is set to false in most objects according to the OpenAPI specification
	// also protect against deleting unknown fields when overwriting later
	json.RejectUnknownMembers(true),
	json.WithMarshalers(jsonutil.URLMarshal),
	json.WithUnmarshalers(jsonutil.URLUnmarshal),
}...)
