package domain

import "errors"

var (
	ErrNotFound = errors.New("object not found")
	InvalidPercentage = errors.New("invalid percentage")
)
