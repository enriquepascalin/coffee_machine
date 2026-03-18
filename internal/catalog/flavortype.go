package catalog

import "fmt"

type FlavorType int

const (
	FlavorNone FlavorType = iota
	FlavorVanilla
	FlavorCaramel
	FlavorChocolate
)

func (f FlavorType) String() string {
	switch f {
	case FlavorNone:
		return "none"
	case FlavorVanilla:
		return "vanilla"
	case FlavorCaramel:
		return "caramel"
	case FlavorChocolate:
		return "chocolate"
	default:
		return fmt.Sprintf("FlavorType(%d)", int(f))
	}
}
