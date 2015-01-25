package gotser

import (
  "time"
)

type AuthToken struct {
  token string
  time Time
}

type OAuth2 interface {
  Auth
}