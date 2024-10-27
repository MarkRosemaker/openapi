package openapi

type (
	// NOTE: This does not work because we run into issue:
	// https://github.com/golang/go/issues/50729
	// SchemaRef      = refOrValue[Schema, *Schema]

	HeaderRef      = refOrValue[Header, *Header]
	ResponseRef    = refOrValue[Response, *Response]
	ParameterRef   = refOrValue[Parameter, *Parameter]
	RequestBodyRef = refOrValue[RequestBody, *RequestBody]
	LinkRef        = refOrValue[Link, *Link]
	ExampleRef     = refOrValue[Example, *Example]
	PathItemRef    = refOrValue[PathItem, *PathItem]

	SchemaRefs []*SchemaRef
)

func getIndexRef[T any, O referencable[T]](ref *refOrValue[T, O]) int    { return ref.idx }
func setIndexRef[T any, O referencable[T]](ref *refOrValue[T, O], i int) { ref.idx = i }
