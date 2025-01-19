package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestTag_JSON(t *testing.T) {
	t.Parallel()

	testJSON(t, []byte(`{
	"name": "pet",
	"description": "Pets operations"
}`), &openapi.Tag{})
}
