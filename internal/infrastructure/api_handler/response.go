package api_handler

import (
	"encoding/json"
	"net/http"

	"github.com/inna-maikut/avito-pvz/internal/api"
)

func InternalError(w http.ResponseWriter, description string) {
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(api.Error{
		Message: description,
	})
}

func BadRequest(w http.ResponseWriter, description string) {
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(api.Error{
		Message: description,
	})
}

func Unauthorized(w http.ResponseWriter, description string) {
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(api.Error{
		Message: description,
	})
}

func Forbidden(w http.ResponseWriter, description string) {
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(api.Error{
		Message: description,
	})
}

func OK[T any](w http.ResponseWriter, t T) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(t)
}

func Created[T any](w http.ResponseWriter, t T) {
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(t)
}
