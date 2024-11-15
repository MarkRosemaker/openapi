package openapi

type (
	// SchemaRef is a reference to a Schema or an actual Schema.
	SchemaRef = refOrValue[Schema, *Schema]
	// HeaderRef is a reference to a Header or an actual Header.
	HeaderRef = refOrValue[Header, *Header]
	// ResponseRef is a reference to a Response or an actual Response.
	ResponseRef = refOrValue[Response, *Response]
	// ParameterRef is a reference to a Parameter or an actual Parameter.
	ParameterRef = refOrValue[Parameter, *Parameter]
	// RequestBodyRef is a reference to a RequestBody or an actual RequestBody.
	RequestBodyRef = refOrValue[RequestBody, *RequestBody]
	// LinkRef is a reference to a Link or an actual Link.
	LinkRef = refOrValue[Link, *Link]
	// ExampleRef is a reference to an Example or an actual Example.
	ExampleRef = refOrValue[Example, *Example]
	// PathItemRef is a reference to a PathItem or an actual PathItem.
	PathItemRef = refOrValue[PathItem, *PathItem]
	// SecuritySchemeRef is a reference to a SecurityScheme or an actual SecurityScheme.
	SecuritySchemeRef = refOrValue[SecurityScheme, *SecurityScheme]
	// CallbackRef is a reference to a Callback or an actual Callback.
	CallbackRef = refOrValue[Callback, *Callback]

	// SchemaRefList is a slice of SchemaRef.
	SchemaRefList []*SchemaRef
)

func getIndexRef[T any, O referencable[T]](ref *refOrValue[T, O]) int    { return ref.idx }
func setIndexRef[T any, O referencable[T]](ref *refOrValue[T, O], i int) { ref.idx = i }
