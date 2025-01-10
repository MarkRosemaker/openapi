package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestSecurityRequirements_Contains(t *testing.T) {
	t.Parallel()

	srs := openapi.SecurityRequirements{
		{"petstore_auth": {"write:pets", "read:pets"}},
	}

	if !srs.Contains(openapi.SecurityRequirement{
		"petstore_auth": []string{"write:pets", "read:pets"},
	}) {
		t.Fatal("should contain")
	}

	if srs.Contains(openapi.SecurityRequirement{
		"petstore_auth": []string{"read:pets"},
	}) {
		t.Fatal("should not contain")
	}
}
