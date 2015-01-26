// TODO - hmac

package sloth

import (
  "time"
)

// TODO - make persistable
type RestHook struct {
  subscriberUrl string

  Kill() (string, error) // rename Unsubscribe?
}

type RestHookResource interface {
  RestResource

  Hooks []RestHook

  Subscribe(subUrl string, subMethod string)
  Broadcast(data)
}

func (resource *RestHookResource) Subscribe(subUrl string, subMethod string) {
  hooks.push(&RestHook {
    subscriberUrl    : subUrl,
    subscriberMethod : subMethod
  })
}

func (resource *RestHookResource) Broadcast(data) {
  for _, hook := range resource.Hooks {
    http.NewRequest(hook.subscriberMethod, hook.subscriberUrls, nil)
    
    // if err != nil {
    //   // handle error
    // }
  }
}

func (hookSub *RestHook) Kill() {
  
}
