package asyncapi

import (
	"io"

	"github.com/MarkRosemaker/yaml"
)

// WriteYAML writes the document in YAML format to the given writer.
//
// "An AsyncAPI document can be JSON or YAML format. [...] In order to preserve the ability to
// round-trip between YAML and JSON formats, YAML version 1.2 is RECOMMENDED along with some
// additional constraints." ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#format
func (d Document) WriteYAML(w io.Writer) error {
	data, err := d.ToYAML()
	if err != nil {
		return err
	}

	_, err = w.Write(data)

	return err
}

// ToYAML marshals the document into YAML.
func (d *Document) ToYAML() ([]byte, error) {
	return yaml.Marshal(d, jsonOpts)
}
