package openapi

import "github.com/MarkRosemaker/errpath"

// Callbacks holds a set of reusable callbacks.
type Callbacks map[string]Callback

// Validate the values of the map.
func (cs Callbacks) Validate() error {
	for op, c := range cs {
		if err := c.Validate(); err != nil {
			return &errpath.ErrKey{Key: op, Err: err}
		}
	}

	return nil
}
