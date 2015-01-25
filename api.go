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

// type RestRequest struct {
//   *http.Request
// }

// type RestResponse struct {
//   *http.Response
// }

type RestResponse (int, interface{})

type RestError interface {
  Error() string
}

type Service interface {
  uri, contentType, string

  // FIXME - move to resource level for greater flexibility
  MarshalContent(data)
  RequestHandler(resource RestResource) http.HandlerFunc
}

func (service *Service) MarshalContent(data) {
  return data
}

type JsonService interface {
  contentType := 'application/json'

  Service
}

func (service *JsonService) MarshalContent(data) {
  return json.Marshal(data)
}

func (service *Service) RequestHandler(resource RestResource) http.HandlerFunc {
  return func(rw http.ResponseWriter, request *http.Request) {
    var data interface{}
    var code int

    request.ParseForm()
    method := request.Method
    values := request.Form

    switch method {
    case GET:
      code, data = resource.Get(values)
    case POST:
      code, data = resource.Post(values)
    case PUT:
      code, data = resource.Put(values)
    case DELETE:
      code, data = resource.Delete(values)
    default:
      service.Abort(rw, 405)
      return
    }

    content, err := service.MarshalContent(data)

    if err != nil {
      // log - failed to marshal content
      service.Abort(rw, 500)
    }

    rw.WriteHeader(code)
    rw.Write(content)
  }
}

func (service *Service) AddResource(resource RestResource, path string) {
  http.HandleFunc(path, service.requestHandler(resource))
}

func (service *Service) Start(port int) {
  portStr := fmt.Sprintf(":%d", port)

  http.ListenAndServe(portStr, nil)
}

func (service *Service) Abort(rw http.ResponseWriter, statusCode int) {
  rw.WriteHeader(statusCode)
}

type RestResource interface {
  baseUrl string

  Get(values url.Values)    RestResponse
  Post(values url.Values)   RestResponse
  Put(values url.Values)    RestResponse
  Delete(values url.Values) RestResponse

  all()    RestResource
  byId(id) RestResource
}

// default GET (remove eventually)
func (resource *RestResource) Get(values url.Values) RestResponse {
  data := map[string]string{"hello": "world"}
  return 200, data
}

type RestAPI interface {
  host, base string

  resources []RestResource
}
