package main

import (
  "fmt"
  "sloth"
  "net/url"
  // "database/sql"
)

// Basic text resource
type FooText struct {
  sloth.RestResource
}

func (FooText) Get(values url.Values) (int, interface{}) {
  return 200, "hello world!"
}

func (FooText) Post(values url.Values) (int, interface{}) {
  data := map[string]string{"yum": "thanks"}
  return 200, data
}

// Basic json resource
type FooJson struct {
  sloth.JsonResource
}

func (FooJson) Get(values url.Values) (int, interface{}) {
  data := map[string]string{"hello": "json!"}
  return 200, data
}

func main() {
  fmt.Println("Sloth - Restful APIs in Go")

  slothTextResource := FooText{ sloth.RestResource{ UrlSlug: "/hello", ContentType: "text/html; charset=utf8" } }
  slothJsonResource := FooJson{ sloth.JsonResource{ UrlSlug: "/json" }}
  slothService  := sloth.RestService{ Port: 3000 }

  slothService.AddResource(&slothTextResource)
  slothService.AddResource(&slothJsonResource)
  slothService.Start()
}
