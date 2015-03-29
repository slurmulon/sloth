# Sloth
RESTful Web API library with support for Resthooks in Go

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

// Basic hook resource
type FooHook struct {
  sloth.HookResource
}

func (hook *FooHook) NotifyHello() {
  hook.Broadcast("hello world!") // sends message to all hook subscribers
}

func main() {
  fmt.Println("Sloth - Restful APIs in Go")

  slothTextResource := FooText{ sloth.RestResource{ UrlSlug: "hello", ContentType: "text/html; charset=utf8" } }
  slothJsonResource := FooJson{ sloth.JsonResource{ UrlSlug: "json" } }
  slothHookResource := FooHook{ sloth.HookResource{ UrlSlug: "hook" } }
  slothRestService  := sloth.RestService{ Port: 3000 }

  slothRestService.AddResource(&slothTextResource) // http://localhost:3000/hello
  slothRestService.AddResource(&slothJsonResource) // http://localhost:3000/json
  slothRestService.AddResource(&slothHookResource) // http://localhost:3000/hook - automatically supports subscriptions via PUT

  slothRestService.Start()
}
```

## TODO

- [ ] Unit tests
- [ ] Resource-level request handlers
- [ ] Statically-typed Content-Types
- [ ] httprouter
- [ ] Redirects
- [ ] Static assets
- [ ] E-Tags
- [ ] HMAC
- [ ] Hypermedia (links and rels)