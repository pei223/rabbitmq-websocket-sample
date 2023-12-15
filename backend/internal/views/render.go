package views

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pei223/rabbitmq-websocket-sample/internal/logger"
)

func RenderJson(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(v)
	if err != nil {
		logger.Logger.Error().Msg("Interface encode failed")
	}
	return err
}

func Render(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(int(status))
}

func RenderError(w http.ResponseWriter, err error) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var apiErr *ApiError
	if !errors.As(err, &apiErr) {
		logger.Logger.Debug().Msg("not ApiError")
		apiErr = NewUnexpectedErr(err.Error())
	}
	w.WriteHeader(int(apiErr.StatusCode))
	encoder := json.NewEncoder(w)
	err = encoder.Encode(apiErr)
	if err != nil {
		logger.Logger.Error().Msg("ApiError encode failed")
	}
	return err
}
