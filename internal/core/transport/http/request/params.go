package httprequest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	errs "github.com/shitaiv1ck/test-effective-mobile/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, nil
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert `%v` to int: %w", value, errs.ErrInvalidArg)
	}

	return &num, nil
}

func GetUUIDQueryParam(r *http.Request, key string) (*uuid.UUID, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, nil
	}

	uuidValue, err := uuid.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert `%v` to UUID: %w", value, errs.ErrInvalidArg)
	}

	return &uuidValue, nil
}

func GetStringQueryParam(r *http.Request, key string) *string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}

	return &value
}

func GetDateQueryParam(r *http.Request, key string) (*time.Time, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, nil
	}

	date, err := time.Parse("01-2006", value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert `%v` to date format `01-2006`: %w", value, errs.ErrInvalidArg)
	}

	return &date, nil
}
