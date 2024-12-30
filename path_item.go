package openapi

import (
	"net/http"
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// PathItem describes the operations available on a single path.
// A Path Item MAY be empty, due to ACL constraints.
// The path itself is still exposed to the documentation viewer but they will not know which operations and parameters are available.
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#path-item-object
type PathItem struct {
	// An optional, string summary, intended to apply to all operations in this path.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// An optional, string description, intended to apply to all operations in this path. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// An alternative `server` array to service all operations in this path.
	Servers Servers `json:"servers,omitempty" yaml:"servers,omitempty"`
	// A list of parameters that are applicable for all the operations described under this path. These parameters can be overridden at the operation level, but cannot be removed there. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
	Parameters ParameterList `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// A definition of a GET operation on this path.
	Get *Operation `json:"get,omitempty" yaml:"get,omitempty"`
	// A definition of a PUT operation on this path.
	Put *Operation `json:"put,omitempty" yaml:"put,omitempty"`
	// A definition of a POST operation on this path.
	Post *Operation `json:"post,omitempty" yaml:"post,omitempty"`
	// A definition of a DELETE operation on this path.
	Delete *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
	// A definition of a OPTIONS operation on this path.
	Options *Operation `json:"options,omitempty" yaml:"options,omitempty"`
	// A definition of a HEAD operation on this path.
	Head *Operation `json:"head,omitempty" yaml:"head,omitempty"`
	// A definition of a PATCH operation on this path.
	Patch *Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
	// A definition of a TRACE operation on this path.
	Trace *Operation `json:"trace,omitempty" yaml:"trace,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:"-"`

	// an index to the original location of this object
	idx int
}

func getIndexPathItem(p *PathItem) int              { return p.idx }
func setIndexPathItem(p *PathItem, i int) *PathItem { p.idx = i; return p }

// Operations iterates over all operations in the path item.
func (p *PathItem) Operations(yield func(string, *Operation) bool) {
	if op := p.Get; op != nil {
		if !yield(http.MethodGet, op) {
			return
		}
	}
	if op := p.Put; op != nil {
		if !yield(http.MethodPut, op) {
			return
		}
	}
	if op := p.Post; op != nil {
		if !yield(http.MethodPost, op) {
			return
		}
	}
	if op := p.Delete; op != nil {
		if !yield(http.MethodDelete, op) {
			return
		}
	}
	if op := p.Options; op != nil {
		if !yield(http.MethodOptions, op) {
			return
		}
	}
	if op := p.Head; op != nil {
		if !yield(http.MethodHead, op) {
			return
		}
	}
	if op := p.Patch; op != nil {
		if !yield(http.MethodPatch, op) {
			return
		}
	}
	if op := p.Trace; op != nil {
		if !yield(http.MethodTrace, op) {
			return
		}
	}
}

// SetOperation sets the operation for the given method.
func (p *PathItem) SetOperation(method string, op *Operation) {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		p.Get = op
	case http.MethodPut:
		p.Put = op
	case http.MethodPost:
		p.Post = op
	case http.MethodDelete:
		p.Delete = op
	case http.MethodOptions:
		p.Options = op
	case http.MethodHead:
		p.Head = op
	case http.MethodPatch:
		p.Patch = op
	case http.MethodTrace:
		p.Trace = op
	}
}

// Validate validates the path item.
func (p *PathItem) Validate() error {
	if err := p.Parameters.Validate(); err != nil {
		return &errpath.ErrField{Field: "parameters", Err: err}
	}

	for method, op := range p.Operations {
		if err := op.Validate(); err != nil {
			return &errpath.ErrField{Field: method, Err: err}
		}
	}

	return validateExtensions(p.Extensions)
}

func (l *loader) collectPathItemRef(p *PathItemRef, ref ref) {
	if p.Value != nil {
		l.collectPathItem(p.Value, ref)
	}
}

func (l *loader) collectPathItem(p *PathItem, ref ref) {
	l.pathItems[ref.String()] = p
}

func (l *loader) resolvePathItemRef(ref *PathItemRef) error {
	return resolveRef(ref, l.pathItems, l.resolvePathItem)
}

func (l *loader) resolvePathItem(p *PathItem) error {
	if err := l.resolveParameterList(p.Parameters); err != nil {
		return &errpath.ErrField{Field: "parameters", Err: err}
	}

	for method, op := range p.Operations {
		if err := l.resolveOperation(op); err != nil {
			return &errpath.ErrField{Field: method, Err: err}
		}
	}

	return nil
}
