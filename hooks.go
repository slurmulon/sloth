// TODO - hmac

package sloth

import (
  "time"
)

// TODO - make persistable
type RestHook struct {
  subscriberUrl, subscriberMethod string

  Kill() (string, error)
}

type RestHookResource interface {
  RestResource

  Hooks []RestHook

  Subscribe(subUrl string, subMethod string)
  Broadcast(data)
}

func (resource *RestHookResource) Subscribe(subUrl string, subMethod string) {
  hooks.push(&RestHook {
    subscriberUrl    : subUrl
    subscriberMethod : subMethod
  })
}

func (resource *RestHookResource) Broadcast(data) {
  // for each hook
  //    find their preferred method, then do it
  for _, hook := range resource.Hooks {
    // create resource from url
    subHookResource := new(Resource) {

    }

    // http.Get("http://example.com/")
  }
}

func (hookSub *RestHook) Kill() {
  
}
