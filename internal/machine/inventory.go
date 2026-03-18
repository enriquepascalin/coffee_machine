package machine

import "github.com/enriquepascalin/coffee_machine/internal/catalog"

type Inventory struct {
	waterML  int
	beansG   int
	cups     int
	milkML   map[catalog.MilkType]int
	flavorML map[catalog.FlavorType]int
}

type InventorySnapshot struct {
	WaterML  int
	BeansG   int
	Cups     int
	MilkML   map[catalog.MilkType]int
	FlavorML map[catalog.FlavorType]int
}

func NewInventory(
	waterML int,
	beansG int,
	cups int,
	milkML map[catalog.MilkType]int,
	flavorML map[catalog.FlavorType]int,
) Inventory {
	clonedMilk := make(map[catalog.MilkType]int, len(milkML))
	for milkType, amount := range milkML {
		clonedMilk[milkType] = amount
	}

	clonedFlavor := make(map[catalog.FlavorType]int, len(flavorML))
	for flavorType, amount := range flavorML {
		clonedFlavor[flavorType] = amount
	}

	return Inventory{
		waterML:  waterML,
		beansG:   beansG,
		cups:     cups,
		milkML:   clonedMilk,
		flavorML: clonedFlavor,
	}
}

func (i *Inventory) WaterML() int {
	return i.waterML
}

func (i *Inventory) BeansG() int {
	return i.beansG
}

func (i *Inventory) Cups() int {
	return i.cups
}

func (i *Inventory) MilkAmount(milkType catalog.MilkType) int {
	return i.milkML[milkType]
}

func (i *Inventory) FlavorAmount(flavorType catalog.FlavorType) int {
	return i.flavorML[flavorType]
}

func (i *Inventory) Snapshot() InventorySnapshot {
	clonedMilk := make(map[catalog.MilkType]int, len(i.milkML))
	for milkType, amount := range i.milkML {
		clonedMilk[milkType] = amount
	}

	clonedFlavor := make(map[catalog.FlavorType]int, len(i.flavorML))
	for flavorType, amount := range i.flavorML {
		clonedFlavor[flavorType] = amount
	}

	return InventorySnapshot{
		WaterML:  i.waterML,
		BeansG:   i.beansG,
		Cups:     i.cups,
		MilkML:   clonedMilk,
		FlavorML: clonedFlavor,
	}
}

func (i *Inventory) AddWaterML(amount int) {
	i.waterML += amount
}

func (i *Inventory) AddBeansG(amount int) {
	i.beansG += amount
}

func (i *Inventory) AddCups(amount int) {
	i.cups += amount
}

func (i *Inventory) AddMilkML(milkType catalog.MilkType, amount int) {
	i.milkML[milkType] += amount
}

func (i *Inventory) AddFlavorML(flavorType catalog.FlavorType, amount int) {
	i.flavorML[flavorType] += amount
}

func (i *Inventory) ConsumeWaterML(amount int) {
	i.waterML -= amount
}

func (i *Inventory) ConsumeBeansG(amount int) {
	i.beansG -= amount
}

func (i *Inventory) ConsumeCup() {
	i.cups--
}

func (i *Inventory) ConsumeMilkML(milkType catalog.MilkType, amount int) {
	i.milkML[milkType] -= amount
}

func (i *Inventory) ConsumeFlavorML(flavorType catalog.FlavorType, amount int) {
	i.flavorML[flavorType] -= amount
}
