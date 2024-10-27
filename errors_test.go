package openapi_test

import (
	"errors"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("error chain", func(t *testing.T) {
		err := &openapi.ErrField{
			Field: "foo",
			Err: &openapi.ErrField{
				Field: "bar",
				Err: &openapi.ErrKey{
					Key: "baz",
					Err: &openapi.ErrField{
						Field: "qux",
						Err: &openapi.ErrIndex{
							Index: 3,
							Err: &openapi.ErrField{
								Field: "quux",
								Err: &openapi.ErrInvalid[string]{
									Value: "corge",
								},
							},
						},
					},
				},
			},
		}
		if want := `foo.bar["baz"].qux[3].quux ("corge") is invalid`; err.Error() != want {
			t.Fatalf("want: %s, got: %s", want, err.Error())
		}
	})

	t.Run("multiple errors", func(t *testing.T) {
		err := &openapi.ErrField{
			Field: "foo",
			Err: &openapi.ErrKey{
				Key: "bar",
				Err: errors.Join(
					&openapi.ErrIndex{
						Index: 3,
						Err: &openapi.ErrField{
							Field: "name",
							Err: &openapi.ErrInvalid[string]{
								Value:   "corge",
								Message: "duplicate name",
							},
						},
					},
					&openapi.ErrIndex{
						Index: 5,
						Err: &openapi.ErrField{
							Field: "name",
							Err: &openapi.ErrInvalid[string]{
								Value:   "corge",
								Message: "duplicate name",
							},
						},
					},
				),
			},
		}
		if want := `foo["bar"][3].name ("corge") is invalid: duplicate name
foo["bar"][5].name ("corge") is invalid: duplicate name`; err.Error() != want {
			t.Fatalf("want: %s, got: %s", want, err.Error())
		}
	})

	t.Run("required", func(t *testing.T) {
		err, want := &openapi.ErrRequired{}, `a value is required`
		if err.Error() != want {
			t.Fatalf("want: %s, got: %s", want, err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		err, want := &openapi.ErrInvalid[string]{}, `a value is invalid`
		if err.Error() != want {
			t.Fatalf("want: %s, got: %s", want, err)
		}
	})
}
