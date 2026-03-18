package config

import (
	"github.com/enriquepascalin/coffee_machine/internal/catalog"
	"github.com/enriquepascalin/coffee_machine/internal/machine"
	"github.com/enriquepascalin/coffee_machine/internal/money"
)

const (
	AdminPassword   = "admin123"
	WaterJugML      = 20000
	BeanBagG        = 1000
	CupPackCount    = 50
	MilkCartonML    = 1000
	FlavorBottleML  = 750
	ExtraFlavorML   = 20
	ExtraFlavorCost = money.Cents(100)
	WearPerDrink    = 10
)

type Defaults struct {
	AdminPassword string
	Machine       *machine.Machine
}

func NewDefaults() Defaults {
	return Defaults{
		AdminPassword: AdminPassword,
		Machine: machine.New(
			catalog.DefaultCatalog(),
			defaultInventory(),
			defaultCashBox(),
			machine.NewMaintenance(0),
			WearPerDrink,
			ExtraFlavorML,
			ExtraFlavorCost,
		),
	}
}

func defaultInventory() machine.Inventory {
	return machine.NewInventory(
		WaterJugML,
		BeanBagG,
		CupPackCount,
		map[catalog.MilkType]int{
			catalog.MilkWhole:  2 * MilkCartonML,
			catalog.MilkAlmond: 1 * MilkCartonML,
			catalog.MilkOat:    1 * MilkCartonML,
		},
		map[catalog.FlavorType]int{
			catalog.FlavorVanilla:   FlavorBottleML,
			catalog.FlavorCaramel:   FlavorBottleML,
			catalog.FlavorChocolate: FlavorBottleML,
		},
	)
}

func defaultCashBox() money.CashBox {
	return money.NewCashBox(map[money.Denomination]int{
		money.Denom20D: 5,
		money.Denom10D: 5,
		money.Denom5D:  10,
		money.Denom2D:  10,
		money.Denom1D:  20,
		money.Denom50C: 20,
		money.Denom25C: 30,
		money.Denom10C: 30,
		money.Denom5C:  30,
	})
}
