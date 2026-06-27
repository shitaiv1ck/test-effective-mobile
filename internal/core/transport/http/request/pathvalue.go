package httprequest

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	errs "github.com/shitaiv1ck/test-effective-mobile/internal/core/errors"
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	value := r.PathValue(key)
	if value == "" {
		return uuid.UUID{}, fmt.Errorf("empty UUID path value: %w", errs.ErrInvalidArg)
	}

	uuidValue, err := uuid.Parse(value)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to convert `%v` to UUID: %w", value, errs.ErrInvalidArg)
	}

	return uuidValue, nil
}
