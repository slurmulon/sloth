package sloth

import (
  "encoding/json"
)

// Ensures type checking during json marshalling
// var _ json.Marshaler = (*RawMessage)(nil)

type JsonService struct {
  RestService
}

type JsonResource struct {
  UrlSlug string

  RestResource
}

func (service *JsonService) MarshalContent(data interface{}) ([]byte, error) {
  return json.Marshal(data)
}

func (resource *JsonService) Type() string {
  return "application/json"
}

func (resource *JsonResource) MarshalContent(data interface{}) ([]byte, error) {
  return json.Marshal(data)
}

func (resource *JsonResource) Type() string {
  return "application/json"
}

func (resource *JsonResource) Slug() string {
  return resource.UrlSlug
}