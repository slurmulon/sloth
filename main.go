package main

import (
  "fmt"
  "sloth"
  "net/url"
)

// https://cloud.google.com/appengine/docs/go/gettingstarted/helloworld

type Foo struct {
  sloth.RestResource
}

func (Foo) Get(values url.Values) (int, interface{}) {
  return 200, "hello world"
}

func (Foo) Post(values url.Values) (int, interface{}) {
  data := map[string]string{"yum": "thanks"}
  return 200, data
}

func main() {
  fmt.Println("Sloth example")

  slothResource := Foo{ sloth.RestResource{ UrlSlug: "/hello" } }
  slothService  := sloth.RestService{ BaseUrl: "http://foo.bar/api" }

  slothService.AddResource(&slothResource)
  slothService.Start(3000)

  // fmt.Println("sloth resource " + slothResource)
}
