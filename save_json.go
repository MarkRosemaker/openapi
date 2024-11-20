package openapi

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
)

func (d Document) WriteJSON(w io.Writer) error {
	return json.MarshalWrite(w, d, _json.Options)
}

func (d *Document) ToJSON() ([]byte, error) {
	return json.Marshal(d, _json.Options)
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
	defer f.Close()

	return d.WriteJSON(f)
}
