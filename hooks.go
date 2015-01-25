// TODO - hmac

package sloth

import (
  "time"
)

type RestHook struct {
  subscriberUrl, subscriberMethod string

  Kill() (string, error)
}

type RestHookResource interface {
  RestResource

  hooks: []RestHook

  Subscribe() RestHook // new rest hook from resource (might want to return response, not RestHook)
  Broadcast(data)
}

func (resource *RestHookResource) Broadcast(data) {
  // for each hook
  //    find their preferred method, then do it
}

func (hookSub *RestHook) Kill() {
  
}
