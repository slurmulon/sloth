package sloth

// Ensures type checking during json marshalling
var _ json.Marshaler = (*RawMessage)(nil)

type JsonRestService interface {
  // contentType := 'application/json' // can do?
  RestService
}

type JsonRestResource interface {
  RestResource
}

func (service *JsonRestService) MarshalContent(data) {
  return json.Marshal(data)
}

func (resource *JsonRestResource) MarshalContent(data) {
  return json.Marshal(data)
}
