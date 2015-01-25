package gotser

import (
  "time"
)


type Auth interface {
  token : AuthToken // some other type
}

type AuthToken struct {
  value string // name of the object
  time Time // its value
}

type OAuth2 interface {
  Auth
}