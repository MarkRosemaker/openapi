package openapi_test

import (
	"fmt"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestResolve(t *testing.T) {
	t.Parallel()

	for i, data := range []string{
		`"paths":{"/": {
			"get": {"parameters": [{"$ref": "#/components/parameters/myparam"}]}
		}},
		"components": {
			"parameters": {
				"myparam": {
					"name": "myParamName",
					"in": "query",
					"schema": {"type": "string"}
				}
			}
		}`,
		`"webhooks":{"/": {
			"get": {"parameters": [{"$ref": "#/components/parameters/myparam"}]}
		}},
		"components": {
			"parameters": {
				"myparam": {
					"name": "myParamName",
					"in": "query",
					"schema": {"type": "string"}
				}
			}
		}`,
		`"components":{"schemas": {
			"Pet": {"allOf": [{"$ref": "#/components/schemas/Dog"}]},
			"Dog": {"type": "object"}
		}}`,
		`"components":{
			"responses": {
				"PetResponse": {
					"headers": {"myheader": {"$ref": "#/components/headers/someheader"}},
					"description": "A pet response"
				}
			},
			"headers": {
				"someheader": {"schema": {"type": "string"}}
			}
		}`,
		`"components":{
			"parameters": {
				"MyParameter": {
					"name": "myParamName",
					"in": "query",
					"schema": {
						"type": "array",
						"items": {"$ref": "#/components/schemas/MyItems"}
					}
				}
			},
			"schemas": {
				"MyItems": {"type": "string"}
			}
		}`,
		`"components":{"examples": {
			"MyExample": {
				"$ref": "#/components/examples/MyActualExample"
			},
			"MyActualExample": {
			}
		}}`,
	} {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			doc, err := openapi.LoadFromDataJSON([]byte(fmt.Sprintf(
				`{
	"openapi": "3.1.0",
	"info": {
		"title": "test",
		"version": "1.0"
	},%s}`, data)))
			if err != nil {
				t.Fatalf("load from data: %v", err)
			}

			if err := doc.Validate(); err != nil {
				t.Fatalf("validate: %v", err)
			}
		})
	}
}

func TestResolve_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in  string
		err string
	}{
		{`{"paths":{"/": {
		"get": {"parameters": [{"$ref": "#/components/parameters/myparam"}]}
}}}`, `paths["/"].GET.parameters[0]: couldn't resolve "#/components/parameters/myparam"`},
		{`{"webhooks":{"/": {
"get": {"parameters": [{"$ref": "#/components/parameters/myparam"}]}
}}}`, `webhooks["/"].GET.parameters[0]: couldn't resolve "#/components/parameters/myparam"`},
		{`{"components":{"schemas": {
	"Pet": {"allOf": [{"$ref": "#/components/schemas/Dog"}]}
}}}`, `components.schemas["Pet"].allOf[0]: couldn't resolve "#/components/schemas/Dog"`},
		{`{"components":{"responses": {
"PetResponse": {"headers": {"myheader": {"$ref": "#/components/headers/someheader"}}}
}}}`, `components.responses["PetResponse"]["myheader"]: couldn't resolve "#/components/headers/someheader"`},
		{`{"components":{"parameters": {
"MyParameter": {
	"schema": {"items": {"$ref": "#/components/schemas/MyItems"}}
}
}}}`, `components.parameters["MyParameter"].schema.items: couldn't resolve "#/components/schemas/MyItems"`},
		{`{"components":{"examples": {
	"MyExample": {
		"$ref": "#/components/examples/MyActualExample"
	}
	}}}`, `components.examples["MyExample"]: couldn't resolve "#/components/examples/MyActualExample"`},
	} {
		_, err := openapi.LoadFromDataJSON([]byte(tc.in))
		if err == nil || err.Error() != tc.err {
			t.Fatalf("expected error: %q, got: %v", tc.err, err)
		}
	}
}
