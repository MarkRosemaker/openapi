package openapi

import (
	"fmt"
	"regexp"

	"github.com/MarkRosemaker/errpath"
)

var reKey = regexp.MustCompile(`^[a-zA-Z0-9\.\-_]+$`)

// The Components object holds a set of reusable objects for different aspects of the OAS.
// All objects defined within the components object will have no effect on the API unless they are explicitly referenced from properties outside the components object.
// ([Specification])
//
// All the fixed fields are objects that MUST use keys that match the regular expression:
//
//	^[a-zA-Z0-9\.\-_]+$
//
// Field name examples:
//
//	User
//	User_1
//	User_Name
//	user-name
//	my.org.User
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#components-object
type Components struct {
	// An object to hold reusable Schema Objects.
	Schemas Schemas `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	// An object to hold reusable Response Objects.
	Responses ResponsesByName `json:"responses,omitempty" yaml:"responses,omitempty"`
	// An object to hold reusable Parameter Objects.
	Parameters Parameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// An object to hold reusable Example Objects.
	Examples Examples `json:"examples,omitempty" yaml:"examples,omitempty"`
	// An object to hold reusable Request Body Objects.
	RequestBodies RequestBodies `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	// An object to hold reusable Header Objects.
	Headers Headers `json:"headers,omitempty" yaml:"headers,omitempty"`
	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes SecuritySchemes `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	// An object to hold reusable Link Objects.
	Links Links `json:"links,omitempty" yaml:"links,omitempty"`
	// An object to hold reusable Callback Objects.
	Callbacks CallbackRefs `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	// An object to hold reusable Path Item Object.
	PathItems PathItems `json:"pathItems,omitempty" yaml:"pathItems,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

func validateKey(key string) error {
	if reKey.MatchString(key) {
		return nil
	}

	return &errpath.ErrKey{Key: key, Err: &errpath.ErrInvalid[string]{
		Value:   key,
		Message: fmt.Sprintf(`must match the regular expression %q`, reKey),
	}}
}

func (c *Components) Validate() error {
	if err := c.Schemas.Validate(); err != nil {
		return &errpath.ErrField{Field: "schemas", Err: err}
	}

	// validate the key: check if it is a valid key
	for name := range c.Responses.ByIndex() {
		if err := validateKey(name); err != nil {
			return &errpath.ErrField{Field: "responses", Err: err}
		}
	}

	if err := c.Responses.Validate(); err != nil {
		return &errpath.ErrField{Field: "responses", Err: err}
	}

	if err := c.Parameters.Validate(); err != nil {
		return &errpath.ErrField{Field: "parameters", Err: err}
	}

	if err := c.Examples.Validate(); err != nil {
		return &errpath.ErrField{Field: "examples", Err: err}
	}

	if err := c.RequestBodies.Validate(); err != nil {
		return &errpath.ErrField{Field: "requestBodies", Err: err}
	}

	if err := c.Headers.Validate(); err != nil {
		return &errpath.ErrField{Field: "headers", Err: err}
	}

	if err := c.SecuritySchemes.Validate(); err != nil {
		return &errpath.ErrField{Field: "securitySchemes", Err: err}
	}

	if err := c.Links.Validate(); err != nil {
		return &errpath.ErrField{Field: "links", Err: err}
	}

	if err := c.Callbacks.Validate(); err != nil {
		return &errpath.ErrField{Field: "callbacks", Err: err}
	}

	if err := c.PathItems.Validate(); err != nil {
		return &errpath.ErrField{Field: "pathItems", Err: err}
	}

	if err := validateExtensions(c.Extensions); err != nil {
		return err
	}

	return nil
}

// For each field that is a map, sorts the map by key.
func (c *Components) SortAll() {
	c.Schemas.Sort()
	c.Responses.Sort()
	c.Parameters.Sort()
	c.Examples.Sort()
	c.RequestBodies.Sort()
	c.Headers.Sort()
	c.SecuritySchemes.Sort()
	c.Links.Sort()
	c.Callbacks.Sort()
	c.PathItems.Sort()
}

func (c Components) isEmpty() bool {
	return len(c.Schemas) == 0 &&
		len(c.Responses) == 0 &&
		len(c.Parameters) == 0 &&
		len(c.Examples) == 0 &&
		len(c.RequestBodies) == 0 &&
		len(c.Headers) == 0 &&
		len(c.SecuritySchemes) == 0 &&
		len(c.Links) == 0 &&
		len(c.Callbacks) == 0 &&
		len(c.PathItems) == 0 &&
		len(c.Extensions) == 0
}

func (l *loader) collectComponents(cs Components, ref ref) {
	l.collectSchemas(cs.Schemas, append(ref, "schemas"))
	l.collectResponses(cs.Responses, append(ref, "responses"))
	l.collectParameters(cs.Parameters, append(ref, "parameters"))
	l.collectExamples(cs.Examples, append(ref, "examples"))
	l.collectRequestBodies(cs.RequestBodies, append(ref, "requestBodies"))
	l.collectHeaders(cs.Headers, append(ref, "headers"))
	l.collectSecuritySchemes(cs.SecuritySchemes, append(ref, "securitySchemes"))
	l.collectLinks(cs.Links, append(ref, "links"))
	l.collectCallbackRefs(cs.Callbacks, append(ref, "callbacks"))
	l.collectPathItems(cs.PathItems, append(ref, "pathItems"))
}

func (l *loader) resolveComponents(c Components) error {
	if err := l.resolveSchemas(c.Schemas); err != nil {
		return &errpath.ErrField{Field: "schemas", Err: err}
	}

	if err := l.resolveResponses(c.Responses); err != nil {
		return &errpath.ErrField{Field: "responses", Err: err}
	}

	if err := l.resolveParameters(c.Parameters); err != nil {
		return &errpath.ErrField{Field: "parameters", Err: err}
	}

	if err := l.resolveExamples(c.Examples); err != nil {
		return &errpath.ErrField{Field: "examples", Err: err}
	}

	if err := l.resolveRequestBodies(c.RequestBodies); err != nil {
		return &errpath.ErrField{Field: "requestBodies", Err: err}
	}

	if err := l.resolveHeaders(c.Headers); err != nil {
		return &errpath.ErrField{Field: "headers", Err: err}
	}

	if err := l.resolveSecuritySchemes(c.SecuritySchemes); err != nil {
		return &errpath.ErrField{Field: "securitySchemes", Err: err}
	}

	if err := l.resolveLinks(c.Links); err != nil {
		return &errpath.ErrField{Field: "links", Err: err}
	}

	if err := l.resolveCallbackRefs(c.Callbacks); err != nil {
		return &errpath.ErrField{Field: "callbacks", Err: err}
	}

	if err := l.resolvePathItems(c.PathItems); err != nil {
		return &errpath.ErrField{Field: "pathItems", Err: err}
	}

	return nil
}
