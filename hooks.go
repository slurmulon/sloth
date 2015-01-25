// TODO - rest hooks!
// TODO - hmac

package sloth

import (
  "time"
)

type RestHookResource interface {
  RestResource

  subscribe() RestHook // new rest hook from resource (might want to return response, not RestHook)
}

type RestHook struct {
  subscriberUrl string

  // Confirm(key) (string, error)
  Broadcast(data) (string, error)
  Kill() (string, error)
}

func (hookSub *RestHook) Confirm(key) (string, error) {

}

func (hookSub *RestHook) Broadcast(data) (string, error) {
  
}

func ( hookSub *RestHook) Kill() (string, error) {
  
}