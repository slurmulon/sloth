// http://golang.org/pkg/net/http/

package sloth

import (
  "fmt"
  "net/http"
  "net/url"
)

// Methods

const (
  GET    = "GET"
  POST   = "POST"
  PUT    = "PUT"
  DELETE = "DELETE"
)

type Getable   interface { Get(values url.Values)    (int, interface{}) }
type Postable  interface { Post(values url.Values)   (int, interface{}) }
type Putable   interface { Put(values url.Values)    (int, interface{}) }
type Deletable interface { Delete(values url.Values) (int, interface{}) }

func (resource *RestResource) Get(values url.Values)    (int, interface{}) { return 405, "" }
func (resource *RestResource) Put(values url.Values)    (int, interface{}) { return 405, "" }
func (resource *RestResource) Post(values url.Values)   (int, interface{}) { return 405, "" }
func (resource *RestResource) Delete(values url.Values) (int, interface{}) { return 405, "" }

// func (getable *Getable) Get(values url.Values) (int, interface{}) {
//   return 405, ""
// }

// Resources

// var _ RestfulResource = (*RestResource)(nil)

type RestError interface {
  Error() string
}

type RestfulResource interface {
  all()     (int, interface{})
  byId(int) (int, interface{})

  //MarshalContent(data interface{}) (interface{}, interface{}) 
  MarshalContent(data interface{}) ([]byte, error)
  RequestHandler() http.HandlerFunc
}

type RestResource struct {
  baseUrl, contentType string
}

func (resource *RestResource) MarshalContent(data interface{}) ([]byte, error) {//(interface{}, interface{}) {
  return AsBytes(data)
}

// type RestRequestInterceptor func(int, interface{})

func (resource *RestResource) RequestHandler() http.HandlerFunc {
  return func(rw http.ResponseWriter, request *http.Request) {
    var data interface{}
    var stat int

    request.ParseForm()
    method := request.Method
    values := request.Form

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

    // request filter TODO
    // requestInterceptor

    content, err := resource.MarshalContent(data)

    if err != nil {
      resource.AbortRequest(rw, 500)
    }

    rw.WriteHeader(stat)
    rw.Write(content)
  }
}

func (resource *RestResource) AbortRequest(rw http.ResponseWriter, statusCode int) {
  rw.WriteHeader(statusCode)
}

type RestAPI struct {
  host, base string

  resources []RestResource
}

// Services

// var _ RestfulService = (*RestService)(nil)

type RestfulService interface {
  MarshalContent(data interface{})
  RequestHandler(resource RestResource) http.HandlerFunc
}

type RestService struct {
  baseUri string
}

func (service *RestService) MarshalContent(data interface{}) ([]byte, error) {
  return AsBytes(data)
}

func (service *RestService) AddResource(resource RestResource, path string) { // TODO - make path deprecated, get it from resource
  http.HandleFunc(path, resource.RequestHandler())
}

func (service *RestService) Start(port int) {
  portStr := fmt.Sprintf(":%d", port)

  http.ListenAndServe(portStr, nil)
}

func (service *RestService) Abort(rw http.ResponseWriter, statusCode int) {
  rw.WriteHeader(statusCode)
}