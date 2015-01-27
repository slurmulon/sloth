package sloth

import (
  "encoding/json"
)

// Ensures type checking during json marshalling
// var _ json.Marshaler = (*RawMessage)(nil)

type JsonRestService struct {
  *RestService
}

type JsonRestResource struct {
  *RestResource
}

func (service *JsonRestService) MarshalContent(data interface{}) {
  return json.Marshal(data)
}

func (resource *JsonRestResource) MarshalContent(data interface{}) {
  return json.Marshal(data)
}
