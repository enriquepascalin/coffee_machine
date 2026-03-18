package money

import "errors"

var (
	ErrInvalidDenomination = errors.New("invalid denomination")
	ErrInsufficientPayment = errors.New("insufficient payment")
	ErrCannotMakeChange    = errors.New("cannot make change with available cash")
)
