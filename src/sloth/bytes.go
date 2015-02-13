package sloth

import (
  "encoding/gob"
  "bytes"
)

func AsBytes(key interface{}) ([]byte, error) {
  var buf bytes.Buffer

  enc := gob.NewEncoder(&buf)
  err := enc.Encode(key)

  if err != nil {
    return nil, err
  }
  
  return buf.Bytes(), nil
}
