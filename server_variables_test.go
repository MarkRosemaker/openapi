package openapi

import (
	"slices"
	"testing"
)

func TestServerVariables_ByIndex(t *testing.T) {
	t.Parallel()

	// test that the keys are sorted by index
	want := []string{"foo", "bar", "baz"}
	sv := ServerVariables{
		"foo": &ServerVariable{idx: 1},
		"bar": &ServerVariable{idx: 2},
		"baz": &ServerVariable{idx: 3},
	}

	i := 0
	for k := range sv.ByIndex() {
		if want[i] != k {
			t.Fatalf("got: %v, want: %v", k, want[i])
		}
		i++
	}

	// test that additional keys are included but come after the sorted keys
	want = append(want, "qux", "moo", "one", "two", "three")
	sv["qux"] = &ServerVariable{}
	sv["moo"] = &ServerVariable{}
	sv["one"] = &ServerVariable{}
	sv["two"] = &ServerVariable{}
	sv["three"] = &ServerVariable{}

	i = 0
	for k := range sv.ByIndex() {
		if i < 3 {
			if want[i] != k {
				t.Fatalf("got: %v, want: %v", k, want[i])
			}
		} else {
			if slices.Contains(want[3:], k) {
				// remove k from want
				want = slices.DeleteFunc(want, func(w string) bool { return w == k })
			} else {
				t.Fatalf("unexpected key: %v", k)
			}
		}
		i++
	}
}

func TestServerVariables_Validate(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		var sv ServerVariables
		if err := sv.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("empty", func(t *testing.T) {
		sv := ServerVariables{}
		if err := sv.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("empty enum array", func(t *testing.T) {
		sv := ServerVariables{"default": &ServerVariable{Enum: []string{}, idx: 1}}
		if err := sv.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := `"default": enum array must not be empty`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}
