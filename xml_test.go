package openapi_test

import (
	"fmt"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func xmlDoc(xml string) []byte {
	return fmt.Appendf(nil, `{
  "openapi": "3.1.0",
  "info": {
    "title": "t",
    "version": "1"
  },
  "components": {
    "schemas": {
      "FlexStatementResponse": {
        "type": "object",
        "properties": {
          "timestamp": {
            "type": "string",
            "xml": {
              "attribute": true
            }
          },
          "Status": {
            "type": "string"
          },
          "Books": {
            "type": "array",
            "items": {
              "type": "string",
              "xml": {
                "name": "Book"
              }
            },
            "xml": {
              "wrapped": true
            }
          }
        },
        "xml": %s
      }
    }
  }
}`, xml)
}

// An element's name, namespace and prefix, an attribute, and a wrapped list are read and written as they are.
func TestXML(t *testing.T) {
	t.Parallel()

	testJSON(t, xmlDoc(`{
          "name": "FlexStatementResponse",
          "namespace": "https://example.com/flex",
          "prefix": "flex"
        }`))
}

func TestXML_RelativeNamespace(t *testing.T) {
	t.Parallel()

	doc, err := openapi.LoadFromDataJSON(xmlDoc(`{"namespace": "flex"}`))
	if err != nil {
		t.Fatal(err)
	}

	want := `components.schemas["FlexStatementResponse"].xml.namespace ("flex") is invalid: must be an absolute URI`
	if err := doc.Validate(); err == nil || err.Error() != want {
		t.Fatalf("got %v, want %s", err, want)
	}
}
