package openapi_test

import (
	"errors"
	"testing"
)

func errAs[T any, E interface {
	*T
	error
}](t *testing.T, err error,
) E {
	t.Helper()

	var zero T
	target := E(&zero)
	if !errors.As(err, &target) {
		t.Fatalf("want: %T, got: %T", target, err)
	}

	return target
}
