package sloth

import "time"

type AuthToken struct {
  token string
  time time.Time
}

type OAuth2 interface {
  // Auth
}