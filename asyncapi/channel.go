package asyncapi

import (
	"regexp"
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// reChannelAddressExpression matches a Channel Address Expression, i.e. a name enclosed in curly braces.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#channelAddressExpressions
var reChannelAddressExpression = regexp.MustCompile(`\{([^{}]+)\}`)

// Channel describes a shared communication channel.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#channelObject
type Channel struct {
	// An optional string representation of this channel's address.
	// The address is typically the "topic name", "routing key", "event type", or "path".
	// When null or absent, it MUST be interpreted as unknown.
	// This is useful when the address is generated dynamically at runtime or can't be known upfront.
	// It MAY contain Channel Address Expressions.
	// Query parameters and fragments SHALL NOT be used, instead use bindings to define them.
	Address string `json:"address,omitempty" yaml:"address,omitempty"`
	// A map of the messages that will be sent to this channel by any application at any time.
	// Every message sent to this channel MUST be valid against one, and only one, of the message objects defined in this map.
	Messages Messages `json:"messages,omitempty" yaml:"messages,omitempty"`
	// A human-friendly title for the channel.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// A short summary of the channel.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// An optional description of this channel. CommonMark syntax can be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// An array of $ref pointers to the definition of the servers in which this channel is available.
	// If the channel is located in the root Channels Object, it MUST point to a subset of server definitions located in the root Servers Object.
	// If `servers` is absent or empty, this channel MUST be available on all the servers defined in the Servers Object.
	Servers ServerRefList `json:"servers,omitempty" yaml:"servers,omitempty"`
	// A map of the parameters included in the channel address.
	// It MUST be present only when the address contains Channel Address Expressions.
	Parameters Parameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// A list of tags for logical grouping of channels.
	Tags Tags `json:"tags,omitempty" yaml:"tags,omitempty"`
	// Additional external documentation for this channel.
	ExternalDocs *ExternalDocsRef `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// A map where the keys describe the name of the protocol and the values describe protocol-specific definitions for the channel.
	Bindings *BindingsRef `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the channel for correctness.
func (c *Channel) Validate() error {
	c.Description = strings.TrimSpace(c.Description)

	if err := c.Messages.Validate(); err != nil {
		return &errpath.ErrField{Field: "messages", Err: err}
	}

	if err := c.Servers.Validate(); err != nil {
		return &errpath.ErrField{Field: "servers", Err: err}
	}

	if err := c.Parameters.Validate(); err != nil {
		return &errpath.ErrField{Field: "parameters", Err: err}
	}

	// the parameters map MUST contain all the parameters used in the channel address
	for _, name := range c.AddressExpressions() {
		if _, ok := c.Parameters[name]; !ok {
			return &errpath.ErrField{Field: "parameters", Err: &errpath.ErrKey{
				Key: name, Err: &errpath.ErrRequired{},
			}}
		}
	}

	if err := c.Tags.Validate(); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if c.ExternalDocs != nil {
		if err := c.ExternalDocs.Validate(); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	if c.Bindings != nil {
		if err := c.Bindings.Validate(); err != nil {
			return &errpath.ErrField{Field: "bindings", Err: err}
		}
	}

	return validateExtensions(c.Extensions)
}

// AddressExpressions returns the names of the Channel Address Expressions used in the address,
// i.e. the names of the parameters that are enclosed in curly braces.
func (c *Channel) AddressExpressions() []string {
	matches := reChannelAddressExpression.FindAllStringSubmatch(c.Address, -1)

	names := make([]string, 0, len(matches))
	for _, m := range matches {
		names = append(names, m[1])
	}

	return names
}

func (l *loader) collectChannelRef(c *ChannelRef, ref ref) {
	if !collectRef(l, c, l.channels, ref) {
		return
	}

	l.collectChannel(c.Value, ref)
}

func (l *loader) collectChannel(c *Channel, ref ref) {
	l.collectMessages(c.Messages, append(ref, "messages"))
	l.collectParameters(c.Parameters, append(ref, "parameters"))
	l.collectTags(c.Tags, append(ref, "tags"))

	if c.ExternalDocs != nil {
		l.collectExternalDocsRef(c.ExternalDocs, append(ref, "externalDocs"))
	}

	if c.Bindings != nil {
		l.collectBindingsRef(c.Bindings, append(ref, "bindings"))
	}
}

func (l *loader) resolveChannelRef(c *ChannelRef) error {
	return resolveRef(c, l.channels, l.resolveChannel)
}

func (l *loader) resolveChannel(c *Channel) error {
	if err := l.resolveMessages(c.Messages); err != nil {
		return &errpath.ErrField{Field: "messages", Err: err}
	}

	if err := l.resolveServerRefList(c.Servers); err != nil {
		return &errpath.ErrField{Field: "servers", Err: err}
	}

	if err := l.resolveParameters(c.Parameters); err != nil {
		return &errpath.ErrField{Field: "parameters", Err: err}
	}

	if err := l.resolveTags(c.Tags); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if c.ExternalDocs != nil {
		if err := l.resolveExternalDocsRef(c.ExternalDocs); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	if c.Bindings != nil {
		if err := l.resolveBindingsRef(c.Bindings); err != nil {
			return &errpath.ErrField{Field: "bindings", Err: err}
		}
	}

	return nil
}
