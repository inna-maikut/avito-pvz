package model

import "errors"

var (
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrWrongEmployeePassword = errors.New("wrong employee password")
	ErrEmployeeAlreadyExists = errors.New("employee already exists")

	ErrReceptionNotFound      = errors.New("reception not found")
	ErrReceptionAlreadyExists = errors.New("reception already exists")

	ErrProductNotFound = errors.New("product not found")

	ErrPVZNotFound = errors.New("PVZ not found")

	ErrNotEnoughBalance               = errors.New("not enough balance")
	ErrSendingCoinsToMyselfNotAllowed = errors.New("sending coins to myself not allowed")
)
