package sloth

import (
  "fmt"
  "log"
  "net/http"
  "net/url"
)

// Methods

const (
  GET    = "GET"
  POST   = "POST"
  PUT    = "PUT"
  PATCH  = "PATCH"
  DELETE = "DELETE"
  HEAD   = "HEAD"
)

func (resource *RestResource) Get(values url.Values)    (int, interface{}) { log.Println("[WARN] Unimplemented GET",    resource); return 405, "" }
func (resource *RestResource) Put(values url.Values)    (int, interface{}) { log.Println("[WARN] Unimplemented PUT",    resource); return 405, "" }
func (resource *RestResource) Post(values url.Values)   (int, interface{}) { log.Println("[WARN] Unimplemented POST",   resource); return 405, "" }
func (resource *RestResource) Patch(values url.Values)  (int, interface{}) { log.Println("[WARN] Unimplemented PATCH",  resource); return 405, "" }
func (resource *RestResource) Delete(values url.Values) (int, interface{}) { log.Println("[WARN] Unimplemented DELETE", resource); return 405, "" }
func (resource *RestResource) Head(values url.Values)   (int, interface{}) { log.Println("[WARN] Unimplemented HEAD",   resource); return 405, "" }


// Resources

var _ RestfulResource = (*RestResource)(nil)

type RestfulResource interface {
  Slug() string
  Type() string

  All()        RestfulResource
  ById(string) RestfulResource

  MarshalContent(data interface{}) ([]byte, error)
  HeaderHandler(headers http.Header) (http.Header, error)
  RequestHandler() http.HandlerFunc
  AbortRequest(rw http.ResponseWriter, statusCode int)

  Get(values url.Values)    (int, interface{})
  Post(values url.Values)   (int, interface{})
  Put(values url.Values)    (int, interface{})
  Patch(values url.Values)  (int, interface{})
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

func (resource *RestResource) All() RestfulResource {
  return &RestResource { UrlSlug: resource.Slug() + "/", ContentType: resource.ContentType }
}

func (resource *RestResource) ById(id string) RestfulResource {
  return &RestResource { UrlSlug: resource.Slug() + "/" + id, ContentType: resource.ContentType }
}

func (resource *RestResource) MarshalContent(data interface{}) ([]byte, error) {
  return AsBytes(data)
}

func (resource *RestResource) HeaderHandler(header http.Header) (http.Header, error) {
  return header, nil
}

func (resource *RestResource) RequestHandler() http.HandlerFunc {
  return func(rw http.ResponseWriter, request *http.Request) {
    var data interface{}
    var stat int

    request.ParseForm()

    method  := request.Method
    values  := request.Form
    headers := rw.Header()

    switch method {
    case GET:
      stat, data = resource.Get(values)
    case POST:
      stat, data = resource.Post(values)
    case PUT:
      stat, data = resource.Put(values)
    case PATCH:
      stat, data = resource.Patch(values)
    case DELETE:
      stat, data = resource.Delete(values)
    default:
      resource.AbortRequest(rw, 405)
      return
    }

    if resource.Type() != "" {
      headers.Set("Content-Type", resource.Type())
    }

    content, contentErr := resource.MarshalContent(data)
    headers, headersErr := resource.HeaderHandler(headers)
   
    if contentErr != nil || headersErr != nil { // TODO - useful logs
      resource.AbortRequest(rw, 500)
    }

    rw.WriteHeader(stat)
    rw.Write(content)
  }
}

func (resource *RestResource) AbortRequest(rw http.ResponseWriter, statusCode int) {
  rw.WriteHeader(statusCode)
}

// Services

var _ RestfulService = (*RestService)(nil)

type RestfulService interface {
  // MarshalContent(data interface{}) ([]byte, error)
  HeaderHandler(header http.Header) (http.Header, error)
  RequestHandler(resource RestfulResource) http.HandlerFunc
}

type RestService struct {
  Port int
}

func (service *RestService) MarshalContent(data interface{}) ([]byte, error) {
  return AsBytes(data)
}

func (service *RestService) HeaderHandler(header http.Header) (http.Header, error) {
  return header, nil
}

func (service *RestService) RequestHandler(resource RestfulResource) http.HandlerFunc {
  return func(rw http.ResponseWriter, request *http.Request) {
    headers      := rw.Header()
    headers, err := service.HeaderHandler(headers)

    if err != nil { // TODO - useful logs
      resource.AbortRequest(rw, 500)
    }

    resource.RequestHandler()(rw, request)
  }
}

// TODO - consider using http.Server instead to provide greater flexibility (handler in addition to read/write timeouts, header constraints, etc)
func (service *RestService) AddResource(resource RestfulResource) {
  // http.HandleFunc("/" + resource.Slug(), service.RequestHandler(resource))
  http.HandleFunc("/" + resource.Slug(), resource.RequestHandler())
}

func (service *RestService) Start() {
  portStr := fmt.Sprintf(":%d", service.Port)

  log.Println("Binding to port ", portStr)

  http.ListenAndServe(portStr, nil)
}
