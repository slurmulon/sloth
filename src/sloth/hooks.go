package sloth

import (
  "fmt"
  "bytes"
  "net/http"
  "net/url"
  "database/sql"
  // "crypto/hmac"
  _ "github.com/go-sql-driver/mysql"
)

// Hooks

var _ RestfulHook = (*RestHook)(nil)

type RestfulHook interface {
  Mesg(interface{}) (*http.Response, error)
  Kill() (sql.Result, error) 
}

type RestHook struct {
  id int
  resourceSlug     string
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

func (hook *RestHook) Kill() (sql.Result, error) {
  return new(HookRepo).Db().Exec("delete from hooks where subscriber_url = ? and subscriber_method = ?", hook.subscriberUrl, hook.subscriberMethod)
}

// Hook resource

var _ HookableResource = (*HookResource)(nil)

type HookableResource interface {
  RestfulResource

  Hooks() []RestHook
  Subscribe(subUrl string, subMethod string)
  Broadcast(data interface{})
}

type HookResource struct {
  UrlSlug string
  RestResource
}

func (resource *HookResource) Hooks() []RestHook {
  var parsedHooks []RestHook

  storedHooks, err := new(HookRepo).ForResource(resource)

  defer storedHooks.Close()

  for storedHooks.Next() {
    var id  int
    var subscriberUrl    string
    var subscriberMethod string

    err := storedHooks.Scan(&id, &subscriberUrl, &subscriberMethod)

    if err != nil {
      panic("Failed to parse RestHook from repository")
    }

    newHook := RestHook{resourceSlug: resource.Slug(), subscriberUrl: subscriberUrl, subscriberMethod: subscriberMethod}

    parsedHooks = append(parsedHooks, newHook)
  }

  err = storedHooks.Err()

  if err != nil {
    panic("Failed to acquire RestHooks from repository")
  }

  return parsedHooks
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

  new(HookRepo).Add(&newHook)

  fmt.Println("Successful subscription!", subUrl, subMethod)
}

func (resource *HookResource) Broadcast(data interface{}) {
  for _, hook := range resource.Hooks() {
    go func() { // WARN - this can get a bit crazy if we have thousands of subscribers. need to make this scale properly
      hook.Mesg(data)
    }()
  }
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
  return repo.Db().Query("select id, resource_slug, subscriber_url, subscriber_method from hooks")
}

func (repo *HookRepo) ForResource(resource *HookResource) (*sql.Rows, error) {
  return repo.Db().Query("select id, subscriber_url, subscriber_method from hooks where resource_slug = ?", resource.Slug())
}

func (repo *HookRepo) Add(hook *RestHook) (sql.Result, error) {
  return repo.Db().Exec("insert into hooks (resource_slug, subscriber_url, subscriber_method) values (?, ?)", hook.resourceSlug, hook.subscriberUrl, hook.subscriberMethod) 
}

func (repo *HookRepo) Delete(hook *RestHook) (sql.Result, error) {
  return repo.Db().Exec("delete from hooks where resource_slug = ? and subscriber_url = ? and subscriber_method = ?", hook.resourceSlug, hook.subscriberUrl, hook.subscriberMethod)
}

func (repo *HookRepo) Close() {
  repo.Db().Close()
}
