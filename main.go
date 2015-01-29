package main

import (
  "fmt"
  "sloth"
  "net/url"
)

// https://cloud.google.com/appengine/docs/go/gettingstarted/helloworld

type FooResource struct {
  UrlSlug string
  
  sloth.RestResource
}

func (FooResource) Get(values url.Values) (int, interface{}) {
  fmt.Println("Successful custom GET!")

  data := map[string]string{"hello": "world"}
  return 200, data
}

func (FooResource) Post(values url.Values) (int, interface{}) {
  data := map[string]string{"yum": "thanks"}
  return 200, data
}

func main() {
  fmt.Println("Sloth example")

  slothResource := &FooResource{UrlSlug: "/hello"}

  var api = sloth.RestService{BaseUrl: "http://foo.bar/api"}

  api.AddResource(slothResource)
  api.Start(3000)

  // fmt.Println("sloth resource " + slothResource)
}
