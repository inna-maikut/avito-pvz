package model

import "errors"

var (
	ErrReceptionNotFound      = errors.New("reception not found")
	ErrReceptionAlreadyExists = errors.New("reception already exists")

	ErrProductNotFound = errors.New("product not found")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrWrongUserPassword = errors.New(("wrong user password"))
	ErrUserNotFound      = errors.New("user not found")
)
