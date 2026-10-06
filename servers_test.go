package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestServers_Validate(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		s := openapi.Servers{}
		if err := s.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("invalid server", func(t *testing.T) {
		s := openapi.Servers{{}}
		if err := s.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := "[0].url is required"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}
