package catalog

import "github.com/enriquepascalin/coffee_machine/internal/money"

type Product struct {
	ID                ProductID
	Name              string
	Emoji             string
	BasePrice         money.Cents
	Recipe            Recipe
	DefaultMilk       MilkType
	DefaultFlavor     FlavorType
	AllowsMilkChoice  bool
	AllowsExtraFlavor bool
	BrewSeconds       int
}
