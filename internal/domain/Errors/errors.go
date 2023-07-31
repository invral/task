package Errors

import (
	"errors"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrZeroBalance     = errors.New("zero balance")
	ErrNegativeBalance = errors.New("negative balance")
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrAccountExists   = errors.New("account already exists")
)
