package service

import (
	"encoding/json"
	"io"
)

func EncodeJson(v interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(v)
}

func DecodeJson(v interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(v)
}
