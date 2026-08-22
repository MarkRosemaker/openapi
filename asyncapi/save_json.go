package asyncapi

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// WriteJSON writes the document in JSON format to the given writer.
//
// "An AsyncAPI document can be JSON or YAML format." ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#format
func (d Document) WriteJSON(w io.Writer) error {
	return json.MarshalWrite(w, d, jsonOpts)
}

// ToJSON marshals the document into JSON.
func (d *Document) ToJSON() ([]byte, error) {
	return json.Marshal(d, jsonOpts)
}

// WriteToFile writes the document to a file, in JSON or in YAML format,
// depending on the extension of the given path.
func (d *Document) WriteToFile(path string) error {
	// determine the file type and write accordingly
	var write func(io.Writer) error

	switch ext := filepath.Ext(path); ext {
	case ".json":
		write = d.WriteJSON
	case ".yaml", ".yml":
		write = d.WriteYAML
	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}

	// create the underlying directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	return errorsJoin(write(f), f.Close())
}

func errorsJoin(err1, err2 error) error {
	if err1 == nil {
		return err2
	}

	if err2 == nil {
		return err1
	}

	return errors.Join(err1, err2)
}
