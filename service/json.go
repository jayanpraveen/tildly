package service

import (
	"encoding/json"
	"io"
)

func DecodeJson(v interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(v)
}
