package main

import (
  "fmt"
  "sloth"
  "net/url"
)

// https://cloud.google.com/appengine/docs/go/gettingstarted/helloworld

// intended syntax
// sloth.Resource('quotes').byId('somsflkdj').Push()
// &sloth.RestResource { 
//   Getable
//   Postable
// }

type FooResource struct {
  UrlSlug string
  
  sloth.RestResource
  // sloth.Getable
  // sloth.Postable
}

func (FooResource) Get(values url.Values) (int, interface{}) {
  fmt.Println("weeeeeee")
  data := map[string]string{"hello": "world"}
  return 200, data
}

func (FooResource) Post(values url.Values) (int, interface{}) {
  data := map[string]string{"yum": "thanks"}
  return 200, data
}

func main() {
  // TODO
  fmt.Println("Sloth example")

  slothResource := &FooResource{UrlSlug: "slkdfj"}

  var api = sloth.RestService{BaseUrl: "http://foo.bar/api"}

  api.AddResource(slothResource, "/hello")
  api.Start(3000)

  // fmt.Println("sloth resource " + slothResource)
}
