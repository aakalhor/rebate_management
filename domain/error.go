package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidInterval     = errors.New("invalid interval")
	ErrNotEligible         = errors.New("criteria is not eligible")
)
