package models

import (
	"net/http"
)

// Error is a generic model, not stored in any database, containing
// error information to be returned to the caller.
type Error struct {
	Code    string `json:"code"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`

	Reasons []Error `json:"reasons,omitempty"`

	Status int `json:"status,omitempty"`
}

// StatusCode implements app.StatusCoder to override the HTTP status code
// when an error is marshalled. If none is configured, HTTP 500 is used.
func (e Error) StatusCode() int {
	if e.Status < 1 {
		return http.StatusInternalServerError
	}

	return e.Status
}

func (e Error) Error() string {
	if e.Field != "" {
		return e.Code + ": " + e.Field + ": " + e.Message
	}

	return e.Code + ": " + e.Message
}

// ErrorWrapper is a helper for marshalling to embed an Error within a key
// called 'error'.
type ErrorWrapper struct {
	Error `json:"error"`
}
