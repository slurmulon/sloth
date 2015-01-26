// TODO - hmac

package sloth

var _ RestfulHook = (*RestHook)(nil)
var _ RestfulHookResource = (*RestHookResource)(nil)

type RestfulHook interface {
  Ping()
  Kill() // rename Unsubscribe?
}

type RestfulHookResource interface {
  RestfulResource

  Subscribe(subUrl string, subMethod string)
  Broadcast(data interface{})
}

// TODO - make persistable
type RestHook struct {
  subscriberUrl    string
  subscriberMethod string
}

type RestHookResource struct {
  *RestResource

  Hooks []RestHook
}

func (resource *RestHookResource) Subscribe(subUrl string, subMethod string) {
  hooks.push(&RestHook {
    subscriberUrl    : subUrl,
    subscriberMethod : subMethod,
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

func (hook *RestHook) Ping() {
  
}

func (hook *RestHook) Kill() {
  
}
