package openapi

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (d Document) WriteJSON(w io.Writer) error {
	return json.MarshalWrite(w, d, jsonOpts)
}

func (d *Document) ToJSON() ([]byte, error) {
	return json.Marshal(d, jsonOpts)
}

func (d *Document) WriteToFile(path string) error {
	switch filepath.Ext(path) {
	case ".json": // ok
	default:
		return fmt.Errorf("unsupported file extension: %s", filepath.Ext(path))
	}

	// create the underlying directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	return errorsJoin(d.WriteJSON(f), f.Close())
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
