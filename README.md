# Sloth
RESTful Web API library with support for Webhooks in Go

Based on: http://dougblack.io/words/a-restful-micro-framework-in-go.html

## Setup

```bash
cd YOUR_SLOTH_DIR
export GOPATH=$(pwd)
```

## Example

```go
package main

import (
  "fmt"
  "sloth"
  "net/url"
)

// Basic text resource
type FooText struct {
  sloth.RestResource
}

func (FooText) Get(values url.Values) (int, interface{}) {
  return 200, "hello world!"
}

func (FooText) Post(values url.Values) (int, interface{}) {
  return 200, "yum, thanks!"
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
  slothTextResource := FooText{ sloth.RestResource{ UrlSlug: "hello", ContentType: "text/html; charset=utf8" } }
  slothJsonResource := FooJson{ sloth.JsonResource{ UrlSlug: "json" } }
  slothRestService  := sloth.RestService{ Port: 3000 }

  slothRestService.AddResource(&slothTextResource) // http://localhost:3000/hello
  slothRestService.AddResource(&slothJsonResource) // http://localhost:3000/json

  slothRestService.Start()
}

```