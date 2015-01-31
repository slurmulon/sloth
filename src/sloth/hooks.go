// TODO - hmac

package sloth

import (
  "fmt"
  "net/http"
  "net/url"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

// Hooks

var _ RestfulHook = (*RestHook)(nil)

type RestfulHook interface {
  Ping()
  Kill()
}

// TODO - make persistable
type RestHook struct {
  subscriberUrl    string
  subscriberMethod string
}

func (hook *RestHook) Url() string {
  return hook.subscriberUrl
}

func (hook *RestHook) Method() string {
  return hook.subscriberMethod
}

func (hook *RestHook) Ping() {
  
}

func (hook *RestHook) Kill() {
  
}

// Hook resource

var _ HookableResource = (*HookResource)(nil)

type HookableResource interface {
  RestfulResource

  Subscribe(subUrl string, subMethod string)
  Broadcast(data interface{})
}

type HookResource struct {
  UrlSlug string

  RestResource

  Hooks []RestHook
}

func (resource *HookResource) Slug() string {
  return resource.UrlSlug
}

func (resource *HookResource) Put(values url.Values) (int, interface{}) {
  fmt.Println("SUBSCRIBE PUT", values, values["subscriber_url"])
  resource.Subscribe(values["subscriber_url"][0], values["subscriber_method"][0])

  return 200, "TODO"
}

func (resource *HookResource) Subscribe(subUrl string, subMethod string) {
  newHook := RestHook {
    subscriberUrl    : subUrl,
    subscriberMethod : subMethod,
  }

  resource.Hooks = append(resource.Hooks, newHook)

  fmt.Println("Successful subscription!", subUrl, subMethod)
}

func (resource *HookResource) Broadcast(data interface{}) {
  for _, hook := range resource.Hooks {
    http.NewRequest(hook.subscriberMethod, hook.subscriberUrl, nil)
    
    // if err != nil {
    //   // handle error
    // }
  }
}

// Hook repository

type HookRepo struct { }
// type HookDB sql.DB

func (repo *HookRepo) Db() *sql.DB {
  db, err := sql.Open("mysql", "user:password@/hooks") // FIXME - integrate with config

  if err != nil {
    // log.Fatal(err)
    fmt.Println("improve this error")
  }

  return db
}

func (repo *HookRepo) All() (*sql.Rows, error) {
  return repo.Db().Query("select id, subscriber_url, subscriber_method from hooks")
}

// TODO
// func (repo *HookRepo) Add(hook *RestfulHook) {
//   // repo.Db().Query("select id, subscriber_url, subscriber_method from hooks")
// }

func (repo *HookRepo) Delete(hook *RestHook) {
  repo.Db().Exec("delete from hooks where subscriber_url = ? and subscriber_method = ?", hook.subscriberUrl, hook.subscriberMethod)
}

func (repo *HookRepo) Close() {
  repo.Db().Close()
}