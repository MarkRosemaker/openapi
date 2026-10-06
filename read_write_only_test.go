package openapi_test

import (
	"encoding/json/v2"
	"strings"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestSchema_ReadOnlyWriteOnly(t *testing.T) {
	t.Parallel()

	const spec = `{"openapi":"3.1.0","info":{"title":"t","version":"1"},"paths":{},"components":{"schemas":{"User":{` +
		`"type":"object","properties":{"id":{"type":"string","readOnly":true},"password":{"type":"string","writeOnly":true}}}}}}`

	doc, err := openapi.LoadFromDataJSON([]byte(spec))
	if err != nil {
		t.Fatal(err)
	}

	props := doc.Components.Schemas["User"].Properties
	if !props["id"].ReadOnly || props["id"].WriteOnly || !props["password"].WriteOnly || props["password"].ReadOnly {
		t.Errorf("got id %+v and password %+v", props["id"], props["password"])
	}

	out, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{`"readOnly":true`, `"writeOnly":true`} {
		if !strings.Contains(string(out), want) {
			t.Errorf("%s lost in %s", want, out)
		}
	}
}
