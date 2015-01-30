// TODO - hmac

package sloth

import "net/http"

// Hooks

var _ RestfulHook = (*RestHook)(nil)

type RestfulHook interface {
  Ping()
  Kill() // rename Unsubscribe?
}

// TODO - make persistable
type RestHook struct {
  subscriberUrl    string
  subscriberMethod string
}

func (hook *RestHook) Ping() {
  
}

func (hook *RestHook) Kill() {
  
}

// Hook resource

var _ RestfulHookResource = (*RestHookResource)(nil)

type RestfulHookResource interface {
  RestfulResource

  Subscribe(subUrl string, subMethod string)
  Broadcast(data interface{})
}

type RestHookResource struct {
  *RestResource

  Hooks []RestHook
}

func (resource *RestHookResource) Subscribe(subUrl string, subMethod string) {
  newHook := RestHook {
    subscriberUrl    : subUrl,
    subscriberMethod : subMethod,
  }

  resource.Hooks = append(resource.Hooks, newHook)
}

func (resource *RestHookResource) Broadcast(data interface{}) {
  for _, hook := range resource.Hooks {
    http.NewRequest(hook.subscriberMethod, hook.subscriberUrl, nil)
    
    // if err != nil {
    //   // handle error
    // }
  }
}
