package sloth

import (
  "encoding/json"
)

// Ensures type checking during json marshalling
var _ json.Marshaler = (*json.RawMessage)(nil)

type JsonService struct {
  RestService
}

type JsonResource struct {
  UrlSlug string

  RestResource
}

func (*JsonService) MarshalContent(data interface{}) ([]byte, error) {
  return json.Marshal(data)
}
 
func (*JsonService) Type() string {
  return "application/json"
}

func (*JsonResource) MarshalContent(data interface{}) ([]byte, error) {
  return json.Marshal(data)
}

func (*JsonResource) Type() string {
  return "application/json"
}

func (resource *JsonResource) Slug() string {
  return resource.UrlSlug
}
