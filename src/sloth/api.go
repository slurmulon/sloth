package sloth

import (
  "fmt"
  "net/http"
  "net/url"
  "github.com/golang/glog"
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

func (resource *RestResource) Get(values url.Values)    (int, interface{}) { glog.Warning("[WARN] Unimplemented GET",    resource); return 405, "" }
func (resource *RestResource) Put(values url.Values)    (int, interface{}) { glog.Warning("[WARN] Unimplemented PUT",    resource); return 405, "" }
func (resource *RestResource) Post(values url.Values)   (int, interface{}) { glog.Warning("[WARN] Unimplemented POST",   resource); return 405, "" }
func (resource *RestResource) Patch(values url.Values)  (int, interface{}) { glog.Warning("[WARN] Unimplemented PATCH",  resource); return 405, "" }
func (resource *RestResource) Delete(values url.Values) (int, interface{}) { glog.Warning("[WARN] Unimplemented DELETE", resource); return 405, "" }
func (resource *RestResource) Head(values url.Values)   (int, interface{}) { glog.Warning("[WARN] Unimplemented HEAD",   resource); return 405, "" }


// Resources

var _ RestfulResource = (*RestResource)(nil)

type RestfulResource interface {
  Slug() string
  Type() string

  All()        RestfulResource
  ById(string) RestfulResource

  MarshalContent(data interface{}) ([]byte, error)
  HeaderHandler(headers http.Header) (http.Header, error)
  // RequestHandler() http.HandlerFunc // TODO

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

// type RestRequestInterceptor func(int, interface{})

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

// TODO - custom headers, ETags
func (service *RestService) RequestHandler(resource RestfulResource) http.HandlerFunc {
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
      service.AbortRequest(rw, 405)
      return
    }

    if resource.Type() != "" {
      headers.Set("Content-Type", resource.Type())
    }

    content, contentErr := resource.MarshalContent(data)
    headers, headersErr := resource.HeaderHandler(headers)

    if contentErr != nil || headersErr != nil {
      service.AbortRequest(rw, 500)
    }

    rw.WriteHeader(stat)
    rw.Write(content)
  }
}

// TODO - consider using http.Server instead to provide greater flexibility (handler in addition to read/write timeouts, header constraints, etc)
func (service *RestService) AddResource(resource RestfulResource) {
  http.HandleFunc("/" + resource.Slug(), service.RequestHandler(resource))
}

func (service *RestService) Start() {
  portStr := fmt.Sprintf(":%d", service.Port)

  glog.Info("Binding to port ", portStr)

  http.ListenAndServe(portStr, nil)
}

func (service *RestService) AbortRequest(rw http.ResponseWriter, statusCode int) {
  rw.WriteHeader(statusCode)
}
