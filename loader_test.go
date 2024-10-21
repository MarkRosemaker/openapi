package openapi_test

import (
	"bytes"
	_ "embed"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/go-json-experiment/json"
)

var (
	//go:embed examples/openapi.json
	exampleJSON []byte
	//go:embed examples/openapi.yaml
	exampleYAML []byte

	openapiDocumentType = reflect.TypeFor[openapi.Document]()
)

func TestLoadFromFile(t *testing.T) {
	t.Parallel()

	t.Run("example json file", func(t *testing.T) {
		t.Parallel()

		doc, err := openapi.LoadFromFile("examples/openapi.json")
		if err != nil {
			t.Fatal(err)
		}

		if err := doc.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("example yaml file", func(t *testing.T) {
		t.Parallel()

		doc, err := openapi.LoadFromFile("examples/openapi.yaml")
		if err != nil {
			t.Fatal(err)
		}

		if err := doc.Validate(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestLoadFromFile_Error(t *testing.T) {
	t.Parallel()

	t.Run("missing file", func(t *testing.T) {
		if _, err := openapi.LoadFromFile("missing.yaml"); err == nil {
			t.Fatal("expected error")
		} else if want := "open missing.yaml: no such file or directory"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("invalid extension", func(t *testing.T) {
		if _, err := openapi.LoadFromFile("examples/invalid.txt"); err == nil {
			t.Fatal("expected error")
		} else if want := "unknown file extension: .txt"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}

func TestLoadFromReader(t *testing.T) {
	t.Parallel()

	t.Run("json", func(t *testing.T) {
		t.Parallel()

		doc, err := openapi.LoadFromReader(bytes.NewReader(exampleJSON))
		if err != nil {
			t.Fatal(err)
		}

		if err := doc.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("yaml", func(t *testing.T) {
		t.Parallel()

		doc, err := openapi.LoadFromReader(bytes.NewReader(exampleYAML))
		if err != nil {
			t.Fatal(err)
		}

		if err := doc.Validate(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestLoadFromReader_Error(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		if _, err := openapi.LoadFromReader(strings.NewReader(``)); err == nil {
			t.Fatal("expected error")
		} else if !errors.Is(err, io.EOF) {
			t.Fatalf("got: %v, want: %v", err, io.EOF)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		if _, err := openapi.LoadFromReader(strings.NewReader(`{"openapi":"3.0.`)); err == nil {
			t.Fatal("expected error")
		} else if want := "unexpected EOF"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("extra field", func(t *testing.T) {
		t.Parallel()

		_, err := openapi.LoadFromReader(strings.NewReader(`   {"extra":"foo"}`))
		assertJSONError(t, err, openapiDocumentType, `unknown name "extra"`)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		t.Parallel()

		if _, err := openapi.LoadFromReader(strings.NewReader(
			`openapi: 3.0.0
			`)); err == nil {
			t.Fatal("expected error")
		} else if want := "yaml: line 2: found a tab character that violates indentation"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}

func TestLoadFromData(t *testing.T) {
	t.Parallel()

	t.Run("json", func(t *testing.T) {
		t.Parallel()

		doc, err := openapi.LoadFromData(exampleJSON)
		if err != nil {
			t.Fatal(err)
		}

		if err := doc.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("yaml", func(t *testing.T) {
		t.Parallel()

		doc, err := openapi.LoadFromData(exampleYAML)
		if err != nil {
			t.Fatal(err)
		}

		if err := doc.Validate(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestLoadFromData_Error(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		if _, err := openapi.LoadFromData([]byte(``)); err == nil {
			t.Fatal("expected error")
		} else if !errors.Is(err, io.EOF) {
			t.Fatalf("got: %v, want: %v", err, io.EOF)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		if _, err := openapi.LoadFromData([]byte(`{"openapi":"3.0.`)); err == nil {
			t.Fatal("expected error")
		} else if want := "yaml: found unexpected end of stream"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("extra field", func(t *testing.T) {
		t.Parallel()

		_, err := openapi.LoadFromData([]byte(`   {"extra":"foo"}`))
		assertJSONError(t, err, openapiDocumentType, `unknown name "extra"`)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		t.Parallel()

		if _, err := openapi.LoadFromData([]byte(
			`openapi: 3.0.0
			`)); err == nil {
			t.Fatal("expected error")
		} else if want := "yaml: line 2: found a tab character that violates indentation"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}

func assertJSONError(t *testing.T, err error, wantGoType reflect.Type, wantErr string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected JSON error")
		return
	}

	jsonErr := &json.SemanticError{}
	if !errors.As(err, &jsonErr) {
		t.Fatalf("expected json.SemanticError, got %T", err)
	} else if jsonErr.GoType != wantGoType {
		t.Fatalf("mismatched go type, got: %s, want: %s", jsonErr.GoType, wantGoType)
	} else if jsonErr.Err.Error() != wantErr {
		t.Fatalf("got: %v, want: %v", jsonErr.Err, wantErr)
	}
}
