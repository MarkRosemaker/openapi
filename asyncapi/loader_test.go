package asyncapi_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

// examplePaths returns the paths of all example documents with the given extension.
func examplePaths(t *testing.T, ext string) []string {
	t.Helper()

	paths, err := filepath.Glob(filepath.Join("examples", "v3.1", "*"+ext))
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) == 0 {
		t.Fatalf("no example documents with extension %q", ext)
	}

	return paths
}

func TestLoadFromFile(t *testing.T) {
	t.Parallel()

	for _, path := range examplePaths(t, ".yaml") {
		t.Run(filepath.Base(path), func(t *testing.T) {
			t.Parallel()

			doc, err := asyncapi.LoadFromFile(path)
			if err != nil {
				t.Fatal(err)
			}

			if err := doc.Validate(); err != nil {
				t.Fatal(err)
			}

			gotJSON, err := doc.ToJSON()
			if err != nil {
				t.Fatal(err)
			}

			// the JSON version of the same document must be identical
			wantJSON, err := os.ReadFile(strings.TrimSuffix(path, ".yaml") + ".json")
			if err != nil {
				t.Fatal(err)
			}

			if got, want := string(gotJSON), strings.TrimSpace(string(wantJSON)); got != want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

func TestLoadFromFile_JSON(t *testing.T) {
	t.Parallel()

	for _, path := range examplePaths(t, ".json") {
		t.Run(filepath.Base(path), func(t *testing.T) {
			t.Parallel()

			doc, err := asyncapi.LoadFromFile(path)
			if err != nil {
				t.Fatal(err)
			}

			if err := doc.Validate(); err != nil {
				t.Fatal(err)
			}

			// writing the document must reproduce the file
			buf := &bytes.Buffer{}
			if err := doc.WriteJSON(buf); err != nil {
				t.Fatal(err)
			}

			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := buf.String(), strings.TrimSpace(string(want)); got != want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

func TestLoadFromFile_Errors(t *testing.T) {
	t.Parallel()

	t.Run("file doesn't exist", func(t *testing.T) {
		t.Parallel()

		if _, err := asyncapi.LoadFromFile("examples/does-not-exist.yaml"); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("unsupported file extension", func(t *testing.T) {
		t.Parallel()

		_, err := asyncapi.LoadFromFile("examples/invalid.txt")
		if err == nil {
			t.Fatal("expected error")
		}

		if want := "unsupported file extension: .txt"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}

func TestLoadFromData(t *testing.T) {
	t.Parallel()

	yamlData, err := os.ReadFile("examples/v3.1/simple.yaml")
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := os.ReadFile("examples/v3.1/simple.json")
	if err != nil {
		t.Fatal(err)
	}

	for name, data := range map[string][]byte{"yaml": yamlData, "json": jsonData} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// the format is detected automatically
			doc, err := asyncapi.LoadFromData(data)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := doc.Info.Title, "Account Service"; got != want {
				t.Fatalf("got: %v, want: %v", got, want)
			}

			// the same document is loaded from a reader
			docFromReader, err := asyncapi.LoadFromReader(bytes.NewReader(data))
			if err != nil {
				t.Fatal(err)
			}

			gotJSON, err := doc.ToJSON()
			if err != nil {
				t.Fatal(err)
			}

			wantJSON, err := docFromReader.ToJSON()
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(gotJSON, wantJSON) {
				t.Fatalf("got:\n%s\nwant:\n%s", gotJSON, wantJSON)
			}
		})
	}
}

func TestLoadFromData_Errors(t *testing.T) {
	t.Parallel()

	t.Run("unresolved reference", func(t *testing.T) {
		t.Parallel()

		_, err := asyncapi.LoadFromDataYAML([]byte(`asyncapi: 3.1.0
info:
  title: Account Service
  version: 1.0.0
channels:
  userSignedup:
    address: user/signedup
    messages:
      UserSignedUp:
        $ref: '#/components/messages/DoesNotExist'
`))
		if err == nil {
			t.Fatal("expected error")
		}

		want := `channels["userSignedup"].messages["UserSignedUp"]: couldn't resolve "#/components/messages/DoesNotExist"`
		if err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		t.Parallel()

		// an unknown field is kept as an extension and only reported when validating
		doc, err := asyncapi.LoadFromDataJSON([]byte(
			`{"asyncapi":"3.1.0","info":{"title":"foo","version":"1.0.0"},` +
				`"channels":{"foo":{}},"doesNotExist":true}`,
		))
		if err != nil {
			t.Fatal(err)
		}

		err = doc.Validate()
		if err == nil {
			t.Fatal("expected error")
		}

		if want := "doesNotExist: " + asyncapi.ErrUnknownField.Error(); err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("unknown field of a channel", func(t *testing.T) {
		t.Parallel()

		doc, err := asyncapi.LoadFromDataJSON([]byte(
			`{"asyncapi":"3.1.0","info":{"title":"foo","version":"1.0.0"},` +
				`"channels":{"foo":{"doesNotExist":true}}}`,
		))
		if err != nil {
			t.Fatal(err)
		}

		err = doc.Validate()
		if err == nil {
			t.Fatal("expected error")
		}

		want := `channels["foo"].doesNotExist: ` + asyncapi.ErrUnknownField.Error()
		if err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()

		if _, err := asyncapi.LoadFromDataJSON([]byte(`{`)); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestWriteToFile(t *testing.T) {
	t.Parallel()

	doc, err := asyncapi.LoadFromFile("examples/v3.1/simple.yaml")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("json", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "sub", "asyncapi.json")
		if err := doc.WriteToFile(path); err != nil {
			t.Fatal(err)
		}

		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		want, err := doc.ToJSON()
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, want) {
			t.Fatalf("got:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("yaml", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "asyncapi.yaml")
		if err := doc.WriteToFile(path); err != nil {
			t.Fatal(err)
		}

		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		want, err := doc.ToYAML()
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, want) {
			t.Fatalf("got:\n%s\nwant:\n%s", got, want)
		}

		// the document that was written can be read again
		again, err := asyncapi.LoadFromFile(path)
		if err != nil {
			t.Fatal(err)
		}

		gotJSON, err := again.ToJSON()
		if err != nil {
			t.Fatal(err)
		}

		wantJSON, err := doc.ToJSON()
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(gotJSON, wantJSON) {
			t.Fatalf("got:\n%s\nwant:\n%s", gotJSON, wantJSON)
		}
	})

	t.Run("unsupported file extension", func(t *testing.T) {
		t.Parallel()

		err := doc.WriteToFile(filepath.Join(t.TempDir(), "asyncapi.txt"))
		if err == nil {
			t.Fatal("expected error")
		}

		if want := "unsupported file extension: .txt"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}
