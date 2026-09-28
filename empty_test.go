package openapi_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestJSON_EmptyMeansAbsent(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		v    any
		want string
	}{
		{"schema required", &openapi.Schema{Type: openapi.TypeObject, Required: []string{}}, `{"type":"object"}`},
		{"schema allOf", &openapi.Schema{AllOf: openapi.SchemaRefList{}}, `{}`},
		{"operation parameters and tags", &openapi.Operation{Parameters: openapi.ParameterList{}, Tags: []string{}}, `{}`},
		{"path item parameters", &openapi.PathItem{Parameters: openapi.ParameterList{}}, `{}`},
		{"response content", &openapi.Response{Description: "OK", Content: openapi.Content{}}, `{"description":"OK"}`},
		{"server variables", &openapi.Server{URL: "/", Variables: openapi.ServerVariables{}}, `{"url":"/"}`},
		// an empty value that means something is still written.
		{"schema enum", &openapi.Schema{Enum: []jsontext.Value{}}, `{"enum":[]}`},
		{"operation security", &openapi.Operation{Security: openapi.SecurityRequirements{}}, `{"security":[]}`},
		{"operation servers", &openapi.Operation{Servers: openapi.Servers{}}, `{"servers":[]}`},
		{"schema properties", &openapi.Schema{Type: openapi.TypeObject, Properties: openapi.SchemaRefs{}}, `{"type":"object","properties":{}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			b, err := json.Marshal(tc.v, jsonOpts)
			if err != nil {
				t.Fatal(err)
			}

			got := jsontext.Value(b)
			if err := got.Compact(); err != nil {
				t.Fatal(err)
			}

			if string(got) != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}
