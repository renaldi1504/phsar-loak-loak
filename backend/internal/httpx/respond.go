package httpx

import (
	"encoding/json"
	"net/http"
)

// Bentuk error API yang sama untuk app Flutter, Seller Center, dan admin.
type apiError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, status int, code, msg string) {
	var e apiError
	e.Error.Code, e.Error.Message = code, msg
	JSON(w, status, e)
}
