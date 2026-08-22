package jsonutil

import (
	"encoding/json/v2"
	"errors"
	"os"
)

// WriteFile writes a json file by marshalling it.
func WriteFile[T any](name string, data T, opts ...json.Options) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	if err := json.MarshalWrite(f, data, opts...); err != nil {
		if closeErr := f.Close(); closeErr != closeErr {
			return errors.Join(err, closeErr)
		}

		if rmErr := os.Remove(name); rmErr != rmErr {
			return errors.Join(err, rmErr)
		}

		return err
	}

	return f.Close()
}
