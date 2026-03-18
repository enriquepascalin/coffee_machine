package catalog

import "github.com/enriquepascalin/coffee_machine/internal/money"

type Catalog struct {
	products []Product
}

func DefaultCatalog() Catalog {
	return Catalog{products: []Product{
		{
			ID:            1,
			Name:          "Espresso",
			Emoji:         "☕",
			BasePrice:     400,
			Recipe:        Recipe{WaterML: 250, BeansG: 16},
			DefaultMilk:   MilkNone,
			DefaultFlavor: FlavorNone,
			BrewSeconds:   2,
		},
		{
			ID:                2,
			Name:              "Latte",
			Emoji:             "🥛",
			BasePrice:         700,
			Recipe:            Recipe{WaterML: 350, BeansG: 20, MilkML: 75},
			DefaultMilk:       MilkWhole,
			DefaultFlavor:     FlavorNone,
			AllowsMilkChoice:  true,
			AllowsExtraFlavor: true,
			BrewSeconds:       3,
		},
		{
			ID:                3,
			Name:              "Cappuccino",
			Emoji:             "☁️",
			BasePrice:         600,
			Recipe:            Recipe{WaterML: 200, BeansG: 12, MilkML: 100},
			DefaultMilk:       MilkWhole,
			DefaultFlavor:     FlavorNone,
			AllowsMilkChoice:  true,
			AllowsExtraFlavor: true,
			BrewSeconds:       3,
		},
		{
			ID:                4,
			Name:              "Mocha",
			Emoji:             "🍫",
			BasePrice:         850,
			Recipe:            Recipe{WaterML: 250, BeansG: 18, MilkML: 80, FlavorML: 20},
			DefaultMilk:       MilkWhole,
			DefaultFlavor:     FlavorChocolate,
			AllowsMilkChoice:  true,
			AllowsExtraFlavor: false,
			BrewSeconds:       4,
		},
		{
			ID:                5,
			Name:              "Vanilla Latte",
			Emoji:             "🌼",
			BasePrice:         money.Cents(900),
			Recipe:            Recipe{WaterML: 300, BeansG: 18, MilkML: 100, FlavorML: 20},
			DefaultMilk:       MilkWhole,
			DefaultFlavor:     FlavorVanilla,
			AllowsMilkChoice:  true,
			AllowsExtraFlavor: false,
			BrewSeconds:       4,
		},
	}}
}

func (c Catalog) All() []Product {
	out := make([]Product, len(c.products))
	copy(out, c.products)
	return out
}

func (c Catalog) FindByID(id ProductID) (Product, bool) {
	for _, p := range c.products {
		if p.ID == id {
			return p, true
		}
	}
	return Product{}, false
}
