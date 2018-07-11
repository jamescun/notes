package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/jamescun/notes/models"

	"github.com/rs/zerolog"
)

// App contains methods that are common to all parts of the application.
type App struct {
	Logger zerolog.Logger
}

// GetLogger returns a logger configured with contextual information
// about the request and the caller.
func (a *App) GetLogger(ctx context.Context) *zerolog.Logger {
	return a.getLogger(ctx, 1)
}

func (a *App) getLogger(ctx context.Context, skip int) *zerolog.Logger {
	_, file, line, _ := runtime.Caller(skip)

	log := a.Logger.With().
		Str("caller", fmt.Sprintf("%s:%d", filepath.Base(file), line)).
		Logger()

	return &log
}

// API contains methods that are common to all API applications returning
// a machine readable response.
type API struct {
	App
}

// StatusCoder is implemented by response bodies which require a non HTTP-200
// response to be sent to the client.
type StatusCoder interface {
	StatusCode() int
}

// OK returns a HTTP 200 (unless body implements StatusCoder) and marshals body
// as JSON to the client.
func (a *API) OK(w http.ResponseWriter, r *http.Request, body interface{}) {
	if body == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	status := http.StatusOK
	if sc, ok := body.(StatusCoder); ok {
		status = sc.StatusCode()
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		a.getLogger(r.Context(), 2).Error().Err(err).Msg("json marshal failed")
	}
}

// Fail returns a HTTP 500 (or other if err implements StatusCode() int) and
// marshals an Error object to the client.
func (a *API) Fail(w http.ResponseWriter, r *http.Request, err error) {
	var body models.Error

	// catch common IO errors where the client has unexpectedly disconnected
	// and return the appropriate error message and status code.
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		err = models.Error{
			Code:    "eof",
			Message: "unexpected end of file",

			Status: http.StatusBadRequest,
		}
	}

	switch err := err.(type) {
	case models.Error:
		body = err

	default:
		body = models.Error{
			Code:    "internal_server_error",
			Message: "an unknown error occurred",

			Status: http.StatusInternalServerError,
		}

		a.getLogger(r.Context(), 2).Error().Err(err).Msg("request failed")
	}

	a.OK(w, r, models.ErrorWrapper{Error: body})
}
