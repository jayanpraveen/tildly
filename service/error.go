package service

import (
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
