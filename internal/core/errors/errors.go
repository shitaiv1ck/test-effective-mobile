package errs

import "errors"

var (
	ErrInvalidArg = errors.New("Invalid argument")
	ErrNotFound   = errors.New("Not found")
)
