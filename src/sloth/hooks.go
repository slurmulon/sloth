// TODO - hmac

package sloth

import (
  "fmt"
  "bytes"
  "net/http"
  "net/url"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

// Hooks

var _ RestfulHook = (*RestHook)(nil)

type RestfulHook interface {
  Mesg(interface{}) (*http.Response, error)
  Kill()
}

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

func (hook *RestHook) Mesg(data interface{}) (*http.Response, error) {
  hookCaller   := &http.Client{}
  dataBytes, _ := AsBytes(data)

  req, err := http.NewRequest(hook.subscriberMethod, hook.subscriberUrl, bytes.NewBuffer(dataBytes))

  if err != nil {
    panic(err)
  }

  req.Header.Add("X-Sloth-Hook-Id",        `W/"something"`)
  req.Header.Add("X-Sloth-Hook-Signature", `W/"abc123"`) // TODO - HMAC

  resp, err := hookCaller.Do(req)

  defer resp.Body.Close()

  return resp, err
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
  resource.Subscribe(values["subscriber_url"][0], values["subscriber_method"][0])
  return 200, "FIXME"
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
  go func() {
    for _, hook := range resource.Hooks {
      hook.Mesg(data)
    }
  }()
}

// Hook repository

type HookRepo struct { }

func (repo *HookRepo) Db() *sql.DB {
  db, err := sql.Open("mysql", "user:password@/hooks") // FIXME - integrate with config

  if err != nil {
    fmt.Println("improve this error")
  }

  return db
}

func (repo *HookRepo) All() (*sql.Rows, error) {
  return repo.Db().Query("select id, subscriber_url, subscriber_method from hooks")
}

func (repo *HookRepo) Add(hook *RestHook) (sql.Result, error) {
  return repo.Db().Exec("insert into hooks (subscriber_url, subscriber_method) values (?, ?)", hook.subscriberUrl, hook.subscriberMethod) 
}

func (repo *HookRepo) Delete(hook *RestHook) (sql.Result, error) {
  return repo.Db().Exec("delete from hooks where subscriber_url = ? and subscriber_method = ?", hook.subscriberUrl, hook.subscriberMethod)
}

func (repo *HookRepo) Close() {
  repo.Db().Close()
}
