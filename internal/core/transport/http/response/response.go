package httpresponse

import (
	"encoding/json"
	"errors"
	"net/http"

	errs "github.com/shitaiv1ck/test-effective-mobile/internal/core/errors"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/logger"
	"go.uber.org/zap"
)

type ResponseWritter struct {
	http.ResponseWriter
	statusCode int
}

func NewRW(w http.ResponseWriter) *ResponseWritter {
	return &ResponseWritter{
		ResponseWriter: w,
	}
}

func (rw *ResponseWritter) ErrorResponse(msg string, err error, logger *logger.Logger) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(setStatusCode(err))

	errDTO := ErrorDTO{
		Message: msg,
		Err:     err.Error(),
	}

	if rw.statusCode == http.StatusInternalServerError {
		logger.Error(msg, zap.Error(err))
	} else {
		logger.Warn(msg, zap.Error(err))
	}

	if err := json.NewEncoder(rw).Encode(errDTO); err != nil {
		panic(err)
	}
}

func (rw *ResponseWritter) JSONResponse(body any, statusCode int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)

	if err := json.NewEncoder(rw).Encode(body); err != nil {
		panic(err)
	}
}

func (rw *ResponseWritter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *ResponseWritter) GetStatusCode() int {
	return rw.statusCode
}

func setStatusCode(err error) int {
	if errors.Is(err, errs.ErrInvalidArg) {
		return http.StatusBadRequest
	}

	if errors.Is(err, errs.ErrNotFound) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
