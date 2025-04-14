package model

import "errors"

var (
	ErrReceptionNotFound      = errors.New("reception not found")
	ErrReceptionAlreadyExists = errors.New("reception already exists")

	ErrProductNotFound = errors.New("product not found")
)
