package openapi

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// WriteJSON writes the document to w as JSON.
func (d *Document) WriteJSON(w io.Writer) error {
	return json.MarshalWrite(w, d, jsonOpts)
}

// ToJSON returns the document as JSON.
func (d *Document) ToJSON() ([]byte, error) {
	return json.Marshal(d, jsonOpts)
}

// WriteToFile writes the document as JSON to the file at path, which must end in .json, creating its directory.
func (d *Document) WriteToFile(path string) error {
	if ext := filepath.Ext(path); ext != ".json" {
		return fmt.Errorf("unsupported file extension: %s", ext)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	return errors.Join(d.WriteJSON(f), f.Close())
}
