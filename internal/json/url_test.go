package json

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
)

var urlType = reflect.TypeFor[*url.URL]()

type hasURL struct {
	A *url.URL `json:"a,omitempty"`
	B url.URL  `json:"b,omitempty"`
	C *url.URL `json:"c"`
	D url.URL  `json:"d"`
}

func TestURL(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "empty",
			in:   `{}`,
			out:  `{"c":null,"d":""}`,
		},
		{
			name: "pointer",
			in:   `{"a":"http://example.com","c":"http://example.com"}`,
			out:  `{"a":"http://example.com","c":"http://example.com","d":""}`,
		},
		{
			name: "non-pointer",
			in:   `{"b":"http://example.com","d":"http://example.com"}`,
			out:  `{"b":"http://example.com","c":null,"d":"http://example.com"}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got hasURL
			if err := json.Unmarshal([]byte(tc.in), &got, Options); err != nil {
				t.Fatal(err)
			}

			b, err := json.Marshal(got, Options)
			if err != nil {
				t.Fatal(err)
			}

			if string(b) != tc.out {
				t.Fatalf("got: %v, want: %v", string(b), tc.out)
			}
		})
	}
}

func TestURL_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		data string
		err  string
	}{
		{
			name: "invalid json",
			data: `{"a":"`,
			err:  `unexpected EOF`,
		},
		{
			name: "int instead of string",
			data: `{"a":42}`,
			err:  `expected string, got 42`,
		},
		{
			name: "invalid URL",
			data: `{"a":"\t"}`,
			err:  `parse "\t": net/url: invalid control character in URL`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got hasURL
			jsonErr := &json.SemanticError{}
			if err := json.Unmarshal([]byte(tc.data), &got, Options); err == nil {
				t.Fatal("expected error")
			} else if !errors.As(err, &jsonErr) {
				t.Fatalf("expected json.SemanticError, got %T", err)
			} else if jsonErr.GoType != urlType {
				t.Fatalf("got: %s, want: %s", jsonErr.GoType, urlType)
			} else if jsonErr.Err.Error() != tc.err {
				t.Fatalf("got: %v, want: %v", jsonErr.Err, tc.err)
			}
		})
	}
}
