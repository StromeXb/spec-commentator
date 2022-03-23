package openapi

import (
	"encoding/json"
	"net/http"
)

func Resp(w http.ResponseWriter, status int, response interface{}) {
	bytes, _ := json.Marshal(response)
	w.WriteHeader(status)
	_, _ = w.Write(bytes)
}

func ErrorResp(w http.ResponseWriter, status int, code string, description string) {
	e := struct {
		Code        *string `json:"code,omitempty"`
		Description *string `json:"description,omitempty"`
	}{
		Code:        &code,
		Description: &description,
	}
	bytes, _ := json.Marshal(e)
	w.WriteHeader(status)
	_, _ = w.Write(bytes)
}
