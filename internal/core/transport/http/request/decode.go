package httprequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	errs "github.com/shitaiv1ck/test-effective-mobile/internal/core/errors"
)

var validate = validator.New()

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("failed to decode: %v: %w", err, errs.ErrInvalidArg)
		}

		return fmt.Errorf("failed to decode: %w", err)
	}

	if err := validate.Struct(dest); err != nil {
		return fmt.Errorf("failed to validate: %w", errs.ErrInvalidArg)
	}

	return nil
}
