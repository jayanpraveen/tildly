package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

type errorResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func SetError(status int, msg string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	e := errorResponse{
		Status: status,
		Msg:    msg,
	}

	err := EncodeJson(e, w)

	if err != nil {
		panic(err)
	}
}

func EncodeJson(v interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(v)
}

func DecodeJson(v interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(v)
}
