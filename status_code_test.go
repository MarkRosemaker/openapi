package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestStatusCode_Validate(t *testing.T) {
	t.Parallel()

	for _, tc := range []openapi.StatusCode{
		"200", "2XX", "default",
	} {
		t.Run(string(tc), func(t *testing.T) {
			if err := tc.Validate(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStatusCode_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		s   openapi.StatusCode
		err string
	}{
		{"", `invalid status code ""`},
		{"600", `invalid status code "600"`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.s.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}

func TestStatusCode_IsSuccess(t *testing.T) {
	for sc, success := range map[openapi.StatusCode]bool{
		"200":     true,
		"2XX":     true,
		"300":     false,
		"3XX":     false,
		"400":     false,
		"404":     false,
		"4XX":     false,
		"default": false, // default does not count as success
		"foo":     false, // invalid status code
	} {
		t.Run(string(sc), func(t *testing.T) {
			if got := sc.IsSuccess(); got != success {
				t.Fatalf("want: %t, got: %t", success, got)
			}
		})
	}
}

func TestStatusCode(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		sc  openapi.StatusCode
		exp string
	}{
		{"200", "OK"},
		{"2XX", ""},
		{"404", "Not Found"},
		{"418", "I'm a teapot"},
		{"default", ""},
	} {
		t.Run(string(tc.sc), func(t *testing.T) {
			if got := tc.sc.StatusText(); got != tc.exp {
				t.Fatalf("want: %s, got: %s", tc.exp, got)
			}
		})
	}
}
