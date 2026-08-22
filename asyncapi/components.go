package asyncapi

import (
	"fmt"
	"regexp"

	"github.com/MarkRosemaker/errpath"
)

// reKey is the regular expression all keys of the fixed fields of the components object must match.
var reKey = regexp.MustCompile(`^[a-zA-Z0-9\.\-_]+$`)

// Components holds a set of reusable objects for different aspects of the AsyncAPI specification.
// All objects defined within the components object will have no effect on the API unless they are
// explicitly referenced from properties outside the components object.
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
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#componentsObject
type Components struct {
	// An object to hold reusable Schema Objects.
	Schemas Schemas `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	// An object to hold reusable Server Objects.
	Servers Servers `json:"servers,omitempty" yaml:"servers,omitempty"`
	// An object to hold reusable Channel Objects.
	Channels Channels `json:"channels,omitempty" yaml:"channels,omitempty"`
	// An object to hold reusable Operation Objects.
	Operations Operations `json:"operations,omitempty" yaml:"operations,omitempty"`
	// An object to hold reusable Message Objects.
	Messages Messages `json:"messages,omitempty" yaml:"messages,omitempty"`
	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes SecuritySchemes `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	// An object to hold reusable Server Variable Objects.
	ServerVariables ServerVariables `json:"serverVariables,omitempty" yaml:"serverVariables,omitempty"`
	// An object to hold reusable Parameter Objects.
	Parameters Parameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// An object to hold reusable Correlation ID Objects.
	CorrelationIDs CorrelationIDs `json:"correlationIds,omitempty" yaml:"correlationIds,omitempty"`
	// An object to hold reusable Operation Reply Objects.
	Replies Replies `json:"replies,omitempty" yaml:"replies,omitempty"`
	// An object to hold reusable Operation Reply Address Objects.
	ReplyAddresses ReplyAddresses `json:"replyAddresses,omitempty" yaml:"replyAddresses,omitempty"`
	// An object to hold reusable External Documentation Objects.
	ExternalDocs ExternalDocsByName `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// An object to hold reusable Tag Objects.
	Tags TagsByName `json:"tags,omitempty" yaml:"tags,omitempty"`
	// An object to hold reusable Operation Trait Objects.
	OperationTraits OperationTraits `json:"operationTraits,omitempty" yaml:"operationTraits,omitempty"`
	// An object to hold reusable Message Trait Objects.
	MessageTraits MessageTraits `json:"messageTraits,omitempty" yaml:"messageTraits,omitempty"`
	// An object to hold reusable Server Bindings Objects.
	ServerBindings BindingsByName `json:"serverBindings,omitempty" yaml:"serverBindings,omitempty"`
	// An object to hold reusable Channel Bindings Objects.
	ChannelBindings BindingsByName `json:"channelBindings,omitempty" yaml:"channelBindings,omitempty"`
	// An object to hold reusable Operation Bindings Objects.
	OperationBindings BindingsByName `json:"operationBindings,omitempty" yaml:"operationBindings,omitempty"`
	// An object to hold reusable Message Bindings Objects.
	MessageBindings BindingsByName `json:"messageBindings,omitempty" yaml:"messageBindings,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// validateKey checks that a key of a fixed field of the components object is well-formed.
func validateKey(key string) error {
	if reKey.MatchString(key) {
		return nil
	}

	return &errpath.ErrKey{Key: key, Err: &errpath.ErrInvalid[string]{
		Value:   key,
		Message: fmt.Sprintf(`must match the regular expression %q`, reKey),
	}}
}

// component is one of the fixed fields of the components object.
type component struct {
	field string
	keys  func(func(string) bool)
	valid func() error
	sort  func()
}

// fields returns the fixed fields of the components object in the order they are defined.
func (c *Components) fields() []component {
	return []component{
		{"schemas", keys(c.Schemas), c.Schemas.Validate, c.Schemas.Sort},
		{"servers", keys(c.Servers), c.Servers.Validate, c.Servers.Sort},
		{"channels", keys(c.Channels), c.Channels.Validate, c.Channels.Sort},
		{"operations", keys(c.Operations), c.Operations.Validate, c.Operations.Sort},
		{"messages", keys(c.Messages), c.Messages.Validate, c.Messages.Sort},
		{"securitySchemes", keys(c.SecuritySchemes), c.SecuritySchemes.Validate, c.SecuritySchemes.Sort},
		{"serverVariables", keys(c.ServerVariables), c.ServerVariables.Validate, c.ServerVariables.Sort},
		{"parameters", keys(c.Parameters), c.Parameters.Validate, c.Parameters.Sort},
		{"correlationIds", keys(c.CorrelationIDs), c.CorrelationIDs.Validate, c.CorrelationIDs.Sort},
		{"replies", keys(c.Replies), c.Replies.Validate, c.Replies.Sort},
		{"replyAddresses", keys(c.ReplyAddresses), c.ReplyAddresses.Validate, c.ReplyAddresses.Sort},
		{"externalDocs", keys(c.ExternalDocs), c.ExternalDocs.Validate, c.ExternalDocs.Sort},
		{"tags", keys(c.Tags), c.Tags.Validate, c.Tags.Sort},
		{"operationTraits", keys(c.OperationTraits), c.OperationTraits.Validate, c.OperationTraits.Sort},
		{"messageTraits", keys(c.MessageTraits), c.MessageTraits.Validate, c.MessageTraits.Sort},
		{"serverBindings", keys(c.ServerBindings), c.ServerBindings.Validate, c.ServerBindings.Sort},
		{"channelBindings", keys(c.ChannelBindings), c.ChannelBindings.Validate, c.ChannelBindings.Sort},
		{"operationBindings", keys(c.OperationBindings), c.OperationBindings.Validate, c.OperationBindings.Sort},
		{"messageBindings", keys(c.MessageBindings), c.MessageBindings.Validate, c.MessageBindings.Sort},
	}
}

// resolvers returns the field name and the reference resolver of each fixed field,
// in the order the fields are defined.
func (l *loader) resolvers(c Components) []struct {
	field   string
	resolve func() error
} {
	return []struct {
		field   string
		resolve func() error
	}{
		{"schemas", func() error { return l.resolveSchemas(c.Schemas) }},
		{"servers", func() error { return l.resolveServers(c.Servers) }},
		{"channels", func() error { return l.resolveChannels(c.Channels) }},
		{"operations", func() error { return l.resolveOperations(c.Operations) }},
		{"messages", func() error { return l.resolveMessages(c.Messages) }},
		{"securitySchemes", func() error { return l.resolveSecuritySchemes(c.SecuritySchemes) }},
		{"serverVariables", func() error { return l.resolveServerVariables(c.ServerVariables) }},
		{"parameters", func() error { return l.resolveParameters(c.Parameters) }},
		{"correlationIds", func() error { return l.resolveCorrelationIDs(c.CorrelationIDs) }},
		{"replies", func() error { return l.resolveReplies(c.Replies) }},
		{"replyAddresses", func() error { return l.resolveReplyAddresses(c.ReplyAddresses) }},
		{"externalDocs", func() error { return l.resolveExternalDocsByName(c.ExternalDocs) }},
		{"tags", func() error { return l.resolveTagsByName(c.Tags) }},
		{"operationTraits", func() error { return l.resolveOperationTraits(c.OperationTraits) }},
		{"messageTraits", func() error { return l.resolveMessageTraits(c.MessageTraits) }},
		{"serverBindings", func() error { return l.resolveBindingsByName(c.ServerBindings) }},
		{"channelBindings", func() error { return l.resolveBindingsByName(c.ChannelBindings) }},
		{"operationBindings", func() error { return l.resolveBindingsByName(c.OperationBindings) }},
		{"messageBindings", func() error { return l.resolveBindingsByName(c.MessageBindings) }},
	}
}

// keys returns a sequence of the keys of a map, in the order they were defined.
func keys[M ~map[string]V, V any](m M) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

// Validate checks the components object for correctness.
func (c *Components) Validate() error {
	for _, f := range c.fields() {
		for key := range f.keys {
			if err := validateKey(key); err != nil {
				return &errpath.ErrField{Field: f.field, Err: err}
			}
		}

		if err := f.valid(); err != nil {
			return &errpath.ErrField{Field: f.field, Err: err}
		}
	}

	return validateExtensions(c.Extensions)
}

// SortMaps sorts each field that is a map by key.
func (c *Components) SortMaps() {
	for _, f := range c.fields() {
		f.sort()
	}

	for _, s := range c.Schemas {
		s.Value.SortMaps()
	}
}

// isEmpty reports whether the components object holds no objects at all.
func (c Components) isEmpty() bool {
	for _, f := range c.fields() {
		for range f.keys {
			return false
		}
	}

	return len(c.Extensions) == 0
}

func (l *loader) collectComponents(c Components, ref ref) {
	l.collectSchemas(c.Schemas, append(ref, "schemas"))
	l.collectServers(c.Servers, append(ref, "servers"))
	l.collectChannels(c.Channels, append(ref, "channels"))
	l.collectOperations(c.Operations, append(ref, "operations"))
	l.collectMessages(c.Messages, append(ref, "messages"))
	l.collectSecuritySchemes(c.SecuritySchemes, append(ref, "securitySchemes"))
	l.collectServerVariables(c.ServerVariables, append(ref, "serverVariables"))
	l.collectParameters(c.Parameters, append(ref, "parameters"))
	l.collectCorrelationIDs(c.CorrelationIDs, append(ref, "correlationIds"))
	l.collectReplies(c.Replies, append(ref, "replies"))
	l.collectReplyAddresses(c.ReplyAddresses, append(ref, "replyAddresses"))
	l.collectExternalDocsByName(c.ExternalDocs, append(ref, "externalDocs"))
	l.collectTagsByName(c.Tags, append(ref, "tags"))
	l.collectOperationTraits(c.OperationTraits, append(ref, "operationTraits"))
	l.collectMessageTraits(c.MessageTraits, append(ref, "messageTraits"))
	l.collectBindingsByName(c.ServerBindings, append(ref, "serverBindings"))
	l.collectBindingsByName(c.ChannelBindings, append(ref, "channelBindings"))
	l.collectBindingsByName(c.OperationBindings, append(ref, "operationBindings"))
	l.collectBindingsByName(c.MessageBindings, append(ref, "messageBindings"))
}

func (l *loader) resolveComponents(c Components) error {
	for _, r := range l.resolvers(c) {
		if err := r.resolve(); err != nil {
			return &errpath.ErrField{Field: r.field, Err: err}
		}
	}

	return nil
}
