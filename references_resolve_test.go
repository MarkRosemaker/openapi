package openapi_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
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
		`"components":{"requestBodies": {
		"MyReqBody": {
			"$ref": "#/components/requestBodies/MyActualReqBody"
		},
		"MyActualReqBody": {"content": {"application/json": {}}}
		}}`,
		`"components":{"headers": {
			"MyHeader": {
				"$ref": "#/components/headers/MyActualHeader"
			},
			"MyActualHeader": {"schema": {"type": "string"}}
			}}`,
		`"components":{"securitySchemes": {
			"MyScheme": {
				"$ref": "#/components/securitySchemes/MyActualScheme"
			},
			"MyActualScheme": {
				"type": "apiKey",
				"name": "myApiKey",
				"in": "header"
			}
		}}`,
		`"components":{"links": {
			"MyLink": {
				"$ref": "#/components/links/MyActualLink"
			},
			"MyActualLink": {
				"operationRef": "myOperationRef"
			}
		}}`,
		`"components":{"callbacks": {
			"MyCallback": {
				"$ref": "#/components/callbacks/MyActualCallback"
			},
			"MyActualCallback": {
			}
		}}`,
		`"components":{"pathItems": {
			"MyPathItem": {
				"$ref": "#/components/pathItems/MyActualPathItem"
			},
			"MyActualPathItem": {
			}
		}}`,
		`"components":{
			"callbacks": {
				"MyCallback": {
					"$request.body#/url": {
						"$ref": "#/components/pathItems/MyPathItem"
					}
				}
			},
			"pathItems": {
				"MyPathItem": {}
			}
		}`,
		`"components":{
			"parameters": {
				"MyParameter": {
					"name": "myParamName",
					"in": "query",
					"content": {
						"application/json": {"schema": {"$ref": "#/components/schemas/MySchema"}}
					}
				}
			},
			"schemas": {
				"MySchema": {"type": "string"}
			}
		}`,
		`"components":{
			"parameters": {
				"MyParameter": {
					"name": "myParamName",
					"in": "query",
					"schema": {"type": "string"},
					"examples": {
						"example": {"$ref": "#/components/examples/MyExample"}
					}
				}
			},
			"examples": {
				"MyExample": {}
			}
		}`,
		`"paths":{"/": {
			"parameters": [{"$ref": "#/components/parameters/myparam"}]
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
		`"components": {
			"schemas": {"MySchema": {"type": "object"}},
			"responses": {
				"MyResponse": {
					"description": "A response",
					"content": {
						"application/json": {
							"schema": {
								"$ref": "#/components/schemas/MySchema"
							}
						}
					}
				}
			}
		}`,
	} {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			data := []byte(fmt.Sprintf(
				`{
"openapi": "3.1.0",
"info": {
"title": "test",
"version": "1.0"
},%s}`, data))

			t.Run("from data", func(t *testing.T) {
				doc, err := openapi.LoadFromData(data)
				if err != nil {
					t.Fatalf("load from data: %v", err)
				}

				if err := doc.Validate(); err != nil {
					t.Fatalf("validate: %v", err)
				}
			})

			t.Run("from reader", func(t *testing.T) {
				doc, err := openapi.LoadFromReader(bytes.NewReader(data))
				if err != nil {
					t.Fatalf("load from data: %v", err)
				}

				if err := doc.Validate(); err != nil {
					t.Fatalf("validate: %v", err)
				}
			})

			t.Run("from file", func(t *testing.T) {
				path := filepath.Join(t.TempDir(), "test.json")
				if err := os.WriteFile(path, data, 0o644); err != nil {
					t.Fatalf("write file: %v", err)
				}

				doc, err := openapi.LoadFromFile(path)
				if err != nil {
					t.Fatalf("load from data: %v", err)
				}

				if err := doc.Validate(); err != nil {
					t.Fatalf("validate: %v", err)
				}
			})
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
}}}`, `components.responses["PetResponse"].headers["myheader"]: couldn't resolve "#/components/headers/someheader"`},
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
		{`{"components":{"requestBodies": {
		"MyReqBody": {
			"$ref": "#/components/requestBodies/MyActualReqBody"
		}
		}}}`, `components.requestBodies["MyReqBody"]: couldn't resolve "#/components/requestBodies/MyActualReqBody"`},
		{`{"components":{"headers": {
			"MyHeader": {
				"$ref": "#/components/headers/MyActualHeader"
			}
			}}}`, `components.headers["MyHeader"]: couldn't resolve "#/components/headers/MyActualHeader"`},
		{`{"components":{"securitySchemes": {
				"MyScheme": {
					"$ref": "#/components/securitySchemes/MyActualScheme"
				}
				}}}`, `components.securitySchemes["MyScheme"]: couldn't resolve "#/components/securitySchemes/MyActualScheme"`},
		{`{"components":{"links": {
			"MyLink": {
				"$ref": "#/components/links/MyActualLink"
			}
		}}}`, `components.links.MyLink: couldn't resolve "#/components/links/MyActualLink"`},
		{`{"components":{"callbacks": {
			"MyCallback": {
				"$ref": "#/components/callbacks/MyActualCallback"
			}
		}}}`, `components.callbacks["MyCallback"]: couldn't resolve "#/components/callbacks/MyActualCallback"`},
		{`{"components":{"pathItems": {
			"MyPathItem": {
				"$ref": "#/components/pathItems/MyActualPathItem"
			}
		}}}`, `components.pathItems["MyPathItem"]: couldn't resolve "#/components/pathItems/MyActualPathItem"`},
		{`{"components":{"callbacks": {
			"MyCallback": {
				"$request.body#/url": {
					"$ref": "#/components/pathItems/MyPathItem"
				}
			}
		}}}`, `components.callbacks["MyCallback"]["$request.body#/url"]: couldn't resolve "#/components/pathItems/MyPathItem"`},
		{`{"components":{"parameters": {
			"MyParameter": {
				"content": {
					"application/json": {"schema": {"$ref": "#/components/schemas/MySchema"}}
				}
			}
		}}}`, `components.parameters["MyParameter"].content["application/json"].schema: couldn't resolve "#/components/schemas/MySchema"`},
		{`{"components":{"parameters": {
			"MyParameter": {
				"examples": {
					"example": {"$ref": "#/components/examples/MyExample"}
				}
			}
		}}}`, `components.parameters["MyParameter"].examples["example"]: couldn't resolve "#/components/examples/MyExample"`},
		{`{"paths":{"/": {
			"parameters": [{"$ref": "#/components/parameters/myparam"}]
	}}}`, `paths["/"].parameters[0]: couldn't resolve "#/components/parameters/myparam"`},
	} {
		data := []byte(tc.in)

		t.Run("from data", func(t *testing.T) {
			if _, err := openapi.LoadFromData(data); err == nil || err.Error() != tc.err {
				t.Fatalf("expected error: %s, got: %v", tc.err, err)
			}
		})

		t.Run("from reader", func(t *testing.T) {
			if _, err := openapi.LoadFromReader(bytes.NewReader(data)); err == nil || err.Error() != tc.err {
				t.Fatalf("expected error: %q, got: %v", tc.err, err)
			}
		})

		t.Run("from file", func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "test.json")
			if err := os.WriteFile(path, data, 0o644); err != nil {
				t.Fatalf("write file: %v", err)
			}

			if _, err := openapi.LoadFromFile(path); err == nil || err.Error() != tc.err {
				t.Fatalf("expected error: %q, got: %v", tc.err, err)
			}
		})
	}
}
