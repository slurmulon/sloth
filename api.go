// http://golang.org/pkg/net/http/

package sloth

import (
  "encoding/json"
  "fmt"
  "net/http"
  "net/url"
  "time"
)

const (
  GET    = "GET"
  POST   = "POST"
  PUT    = "PUT"
  DELETE = "DELETE"
)

// Ensures type checking during json marshalling
var _ json.Marshaler = (*RawMessage)(nil)

type RestRequest  (int, interface{})
type RestResponse (int, interface{})

// TODO - how can i just initialize the contentType variable...
// type Jsonable struct {
//   contentType := 'application/json'
// }

type RestError interface {
  Error() string
}

type Service interface {
  uri, string

  // FIXME - move to resource level for greater flexibility, but allow here as well
  MarshalContent(data)
  RequestHandler(resource RestResource) http.HandlerFunc
}

func (service *Service) MarshalContent(data) {
  return data
}

// type JsonService interface {
//   contentType := 'application/json'

//   Service
// }

func (service *JsonService) MarshalContent(data) {
  return json.Marshal(data)
}

func (service *Service) AddResource(resource RestResource, path string) { // TODO - make path deprecated, get it from resource
  http.HandleFunc(path, resource.requestHandler())
}

func (service *Service) Start(port int) {
  portStr := fmt.Sprintf(":%d", port)

  http.ListenAndServe(portStr, nil)
}

func (service *Service) Abort(rw http.ResponseWriter, statusCode int) {
  rw.WriteHeader(statusCode)
}

type Getable interface {
  Get(values url.Values) RestResponse
}

type Postable interface {
  Post(values url.Values) RestResponse
}

type Putable interface {
  Put(values url.Values) RestResponse
}

type Deletable interface {
  Delete(values url.Values) RestResponse
}

type RestResource interface {
  baseUrl string

  all()    RestResource
  byId(id) RestResource

  MarshalContent(data)
  RequestHandler() http.HandlerFunc
}

func (resource *RestResource) MarshalContent(data) {
  return data
}

// type RestRequestFilter func(int, interface{}) RestRequest

func (resource *RestResource) RequestHandler() http.HandlerFunc {
  return func(rw http.ResponseWriter, request *http.Request) {
    var data interface{}
    var stat int

    request.ParseForm()
    method := request.Method
    values := request.Form

    // TODO - move this all 
    switch method {
    case GET:
      stat, data = resource.Get(values)
    case POST:
      stat, data = resource.Post(values)
    case PUT:
      stat, data = resource.Put(values)
    case DELETE:
      stat, data = resource.Delete(values)
    default:
      resource.AbortRequest(rw, 405)
      return
    }

    content, err := service.MarshalContent(data)

    if err != nil {
      // log - failed to marshal content
      service.Abort(rw, 500)
    }

    rw.WriteHeader(stat)
    rw.Write(content)
  }
}

func (resource *RestResource) AbortRequest(rw http.ResponseWriter, statusCode int) {
  rw.WriteHeader(statusCode)
}

// default GET (remove eventually)
func (resource *RestResource) Get(values url.Values) RestResponse {
  data := map[string]string{"hello": "world"}

  return StatusOK, data
}

type RestAPI interface {
  host, base string

  resources []RestResource
}
