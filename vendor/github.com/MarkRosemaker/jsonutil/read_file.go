package jsonutil

import (
	"encoding/json/v2"
	"errors"
	"os"
)

// ReadFile reads a json file and unmarshals it.
func ReadFile[T any](name string, opts ...json.Options) (T, error) {
	var v T

	f, err := os.Open(name)
	if err != nil {
		return v, err
	}

	if err := json.UnmarshalRead(f, &v, opts...); err != nil {
		if closeErr := f.Close(); closeErr != nil {
			return v, errors.Join(err, closeErr)
		}

		return v, err
	}

	return v, f.Close()
}
