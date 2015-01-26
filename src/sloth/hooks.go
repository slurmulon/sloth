// TODO - hmac

package sloth

// var _ RestfulHook = (*RestfulHookResource)(nil)
// var _ RestfulHookResource = (*RestHookResource)(nil)

type RestfulHook interface {
  Kill() (string, error) // rename Unsubscribe?
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

func (hook *RestHook) Kill() {
  
}
