package json_test

import (
	"errors"
	"iter"
	"reflect"
	"slices"
	"testing"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

var orderdMapType = reflect.TypeFor[orderedMap]()

type orderedMap map[string]*value

func (om orderedMap) ByIndex() iter.Seq2[string, *value] {
	return _json.OrderedMapByIndex(om, func(v *value) int { return v.idx })
}

type value struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`

	idx int
}

func (om *orderedMap) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(om, dec, opts, func(v *value, i int) { v.idx = i })
}

func (om orderedMap) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(&om, enc, opts)
}

func TestOrderedMap_ByIndex(t *testing.T) {
	t.Parallel()

	// test that the keys are sorted by index
	want := []string{"foo", "bar", "baz"}
	om := orderedMap{
		"foo": &value{idx: 1},
		"bar": &value{idx: 2},
		"baz": &value{idx: 3},
	}

	i := 0
	for k := range om.ByIndex() {
		if want[i] != k {
			t.Fatalf("got: %v, want: %v", k, want[i])
		}
		i++
	}

	for range om.ByIndex() {
		break // i.e. yield returns true
	}

	// test that additional keys are included but come after the sorted keys
	want = append(want, "qux", "moo", "one", "two", "three")
	om["qux"] = &value{}
	om["moo"] = &value{}
	om["one"] = &value{}
	om["two"] = &value{}
	om["three"] = &value{}

	i = 0
	for k := range om.ByIndex() {
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

func TestOrderedMap_JSON(t *testing.T) {
	t.Parallel()

	want := `{"foo":{"foo":"foo","bar":1},"bar":{"foo":"foo","bar":1},"baz":{"foo":"foo","bar":1},"qux":{"foo":"","bar":0},"moo":{"foo":"","bar":0},"one":{"foo":"","bar":0},"two":{"foo":"","bar":0},"three":{"foo":"","bar":0}}`
	om := &orderedMap{}
	if err := json.Unmarshal([]byte(want), om); err != nil {
		t.Fatal(err)
	}

	got, err := json.Marshal(om)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != want {
		t.Fatalf("got: %v, want: %v", string(got), want)
	}
}

func TestOrderedMap_Unmarshal_Errors(t *testing.T) {
	t.Parallel()

	valueType := reflect.TypeFor[value]()

	for _, tc := range []struct {
		name  string
		data  string
		types []reflect.Type
		err   string
	}{
		{"empty", ``, nil, `EOF`},
		{"string instead of object", `""`, nil, `expected {, got string`},
		{"missing string for object name", `{1}`, nil, `jsontext: missing string for object name`},
		{"missing string for object name", `{"foo":1}`, []reflect.Type{valueType}, ``},
		{
			"missing string for object name", `{"foo":{"foo":"foo","bar":1`,
			nil, `unmarshal foo: unexpected EOF`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(tc.data), &orderedMap{})
			for _, tp := range append([]reflect.Type{orderdMapType}, tc.types...) {
				err = unwrapJSONError(t, err, tp)
			}

			if tc.err != "" {
				if err == nil {
					t.Fatal("expected error")
				} else if err.Error() != tc.err {
					t.Fatalf("got: %v, want: %v", err, tc.err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}

type impossibleToMarshal struct{ err error }

func (imp impossibleToMarshal) MarshalJSONV2(*jsontext.Encoder, json.Options) error {
	return imp.err
}

type cannotMarshalKey map[impossibleToMarshal]*value

func (om cannotMarshalKey) ByIndex() iter.Seq2[impossibleToMarshal, *value] {
	return _json.OrderedMapByIndex(om, func(v *value) int { return v.idx })
}

func (om cannotMarshalKey) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(&om, enc, opts)
}

func (om *cannotMarshalKey) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(om, dec, opts, func(v *value, i int) { v.idx = i })
}

type cannotMarshalValue map[string]*impossibleToMarshal

func (om cannotMarshalValue) ByIndex() iter.Seq2[string, *impossibleToMarshal] {
	return _json.OrderedMapByIndex(om, func(v *impossibleToMarshal) int { return 0 })
}

func (om cannotMarshalValue) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(&om, enc, opts)
}

func (om *cannotMarshalValue) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(om, dec, opts, func(v *impossibleToMarshal, i int) {})
}

func TestOrderedMap_Marshal_Errors(t *testing.T) {
	t.Parallel()

	someErr := errors.New("some error")

	t.Run("marshalling key", func(t *testing.T) {
		om := cannotMarshalKey{impossibleToMarshal{err: someErr}: &value{}}
		_, err := json.Marshal(om)
		err = unwrapJSONError(t, err, reflect.TypeFor[cannotMarshalKey]())
		err = unwrapJSONError(t, err, reflect.TypeFor[impossibleToMarshal]())
		if err == nil {
			t.Fatal("expected error")
		} else if err.Error() != someErr.Error() {
			t.Fatalf("got: %v, want: %v", err, someErr)
		}
	})

	t.Run("marshalling value", func(t *testing.T) {
		om := cannotMarshalValue{"foo": &impossibleToMarshal{err: someErr}}
		_, err := json.Marshal(om)
		err = unwrapJSONError(t, err, reflect.TypeFor[cannotMarshalValue]())
		err = unwrapJSONError(t, err, reflect.TypeFor[impossibleToMarshal]())
		if err == nil {
			t.Fatal("expected error")
		} else if err.Error() != someErr.Error() {
			t.Fatalf("got: %v, want: %v", err, someErr)
		}
	})
}

func unwrapJSONError(t *testing.T, err error, wantGoType reflect.Type) error {
	t.Helper()

	if err == nil {
		t.Fatalf("expected JSON error for %s", wantGoType)
		return nil
	}

	jsonErr := &json.SemanticError{}
	if !errors.As(err, &jsonErr) {
		t.Fatalf("expected json.SemanticError for %s, got %T", wantGoType, err)
		return nil
	}

	if jsonErr.GoType != wantGoType {
		t.Fatalf("mismatched go type, got: %s, want: %s", jsonErr.GoType, wantGoType)
		return nil
	}

	return jsonErr.Err
}
