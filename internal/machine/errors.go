package machine

import "errors"

var (
	ErrMaintenanceRequired   = errors.New("machine needs maintenance")
	ErrInsufficientWater     = errors.New("not enough water")
	ErrInsufficientBeans     = errors.New("not enough coffee beans")
	ErrInsufficientMilk      = errors.New("not enough milk")
	ErrInsufficientFlavor    = errors.New("not enough flavor syrup")
	ErrInsufficientCups      = errors.New("not enough disposable cups")
	ErrMilkChoiceNotAllowed  = errors.New("milk choice is not allowed for this product")
	ErrExtraFlavorNotAllowed = errors.New("extra flavor is not allowed for this product")
)
