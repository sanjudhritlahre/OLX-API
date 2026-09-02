package httpx

import (
	"encoding/json"
	"net/http"
)

type Code string

const (
	InvalidID        Code = "invalid_id"
	InternalError    Code = "internal_error"
	NotFound         Code = "not_found"
	MalformedJSON    Code = "malformed_json"
	ValidationFailed Code = "validation_failed"
	Unauthenticated  Code = "unauthenticated"
	Forbidden        Code = "forbidden"
	Conflict         Code = "conflict"
	RateLimited      Code = "rate_limited"
)

type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Code    Code `json:"code"`
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, status int, message string, code Code) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(errorEnvelope{Error: errorPayload{Code: code, Message: message}})
}