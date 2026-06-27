package httprequest

import (
	"fmt"
	"net/http"
	"strconv"

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
