// http://golang.org/pkg/net/http/

package sloth

import (
  "fmt"
  "net/http"
  "net/url"
)

// type url.Values map[string][]string

// Methods

const (
  GET    = "GET"
  POST   = "POST"
  PUT    = "PUT"
  DELETE = "DELETE"
)

func (resource *RestResource) Get(values url.Values)    (int, interface{}) { fmt.Println("[WARN] Unimplemented GET",    resource); return 405, "" }
func (resource *RestResource) Put(values url.Values)    (int, interface{}) { fmt.Println("[WARN] Unimplemented PUT",    resource); return 405, "" }
func (resource *RestResource) Post(values url.Values)   (int, interface{}) { fmt.Println("[WARN] Unimplemented POST",   resource); return 405, "" }
func (resource *RestResource) Delete(values url.Values) (int, interface{}) { fmt.Println("[WARN] Unimplemented DELETE", resource); return 405, "" }

// Resources

var _ RestfulResource = (*RestResource)(nil)

type RestError interface {
  Error() string
}

type RestfulResource interface {
  Slug() string
  Type() string

  All()     (int, interface{})
  ById(int) (int, interface{})

  MarshalContent(data interface{}) ([]byte, error)
  // RequestHandler() http.HandlerFunc

  Get(values url.Values) (int, interface{})
  Post(values url.Values) (int, interface{})
  Put(values url.Values) (int, interface{})
  Delete(values url.Values) (int, interface{})
}

type RestResource struct {
  UrlSlug, ContentType string
}

func (resource *RestResource) Slug() string {
  return resource.UrlSlug
}

func (resource *RestResource) Type() string {
  return resource.ContentType
}

func (resource *RestResource) All() (int, interface{}) {
  return 200, "TODO"
}

func (resource *RestResource) ById(id int) (int, interface{}) {
  return 200, "TODO"
}

func (resource *RestResource) MarshalContent(data interface{}) ([]byte, error) {//(interface{}, interface{}) {
  return AsBytes(data)
}

// type RestRequestInterceptor func(int, interface{})

// type RestAPI struct {
//   BaseUrl string

//   resources []RestResource
// }

// Services

var _ RestfulService = (*RestService)(nil)

type RestfulService interface {
  MarshalContent(data interface{}) ([]byte, error)
  RequestHandler(resource RestfulResource) http.HandlerFunc
}

type RestService struct {
  Port int
}

func (service *RestService) MarshalContent(data interface{}) ([]byte, error) {
  return AsBytes(data)
}

func (service *RestService) RequestHandler(resource RestfulResource) http.HandlerFunc {
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
      service.AbortRequest(rw, 405)
      return
    }

    // request filter TODO
    // requestInterceptor

    content, err := resource.MarshalContent(data)

    if err != nil {
      service.AbortRequest(rw, 500)
    }

    if resource.Type() != "" {
      rw.Header().Set("Content-Type", resource.Type()) // "text/html; charset=utf-8")
    }

    rw.WriteHeader(stat)
    rw.Write(content)
  }
}

func (service *RestService) AddResource(resource RestfulResource) { // TODO - make path deprecated, get it from resource
  http.HandleFunc(resource.Slug(), service.RequestHandler(resource))
}

func (service *RestService) Start() {
  portStr := fmt.Sprintf(":%d", service.Port)

  fmt.Println("Binding to port ", portStr)

  http.ListenAndServe(portStr, nil)
}

func (service *RestService) AbortRequest(rw http.ResponseWriter, statusCode int) {
  rw.WriteHeader(statusCode)
}
