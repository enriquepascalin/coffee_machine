package machine

import (
	"fmt"

	"github.com/enriquepascalin/coffee_machine/internal/catalog"
	"github.com/enriquepascalin/coffee_machine/internal/money"
)

type Machine struct {
	catalog          catalog.Catalog
	inventory        Inventory
	cashBox          money.CashBox
	maintenance      Maintenance
	priceOverrides   map[catalog.ProductID]money.Cents
	wearPerDrink     int
	extraFlavorML    int
	extraFlavorPrice money.Cents
}

func New(
	productCatalog catalog.Catalog,
	inventory Inventory,
	cashBox money.CashBox,
	maintenance Maintenance,
	wearPerDrink int,
	extraFlavorML int,
	extraFlavorPrice money.Cents,
) *Machine {
	return &Machine{
		catalog:          productCatalog,
		inventory:        inventory,
		cashBox:          cashBox,
		maintenance:      maintenance,
		priceOverrides:   make(map[catalog.ProductID]money.Cents),
		wearPerDrink:     wearPerDrink,
		extraFlavorML:    extraFlavorML,
		extraFlavorPrice: extraFlavorPrice,
	}
}

func (m *Machine) Products() []catalog.Product {
	return m.catalog.All()
}

func (m *Machine) FindProduct(id catalog.ProductID) (catalog.Product, bool) {
	return m.catalog.FindByID(id)
}

func (m *Machine) InventorySnapshot() InventorySnapshot {
	return m.inventory.Snapshot()
}

func (m *Machine) CashTotal() money.Cents {
	return m.cashBox.Total()
}

func (m *Machine) MaintenanceGauge() int {
	return m.maintenance.GaugePercent()
}

func (m *Machine) CurrentPrice(product catalog.Product) money.Cents {
	if override, ok := m.priceOverrides[product.ID]; ok {
		return override
	}

	return product.BasePrice
}

func (m *Machine) ChangePrice(id catalog.ProductID, newPrice money.Cents) error {
	if newPrice < 0 {
		return fmt.Errorf("negative price: %v", newPrice)
	}

	m.priceOverrides[id] = newPrice
	return nil
}

func (m *Machine) PriceFor(product catalog.Product, extraFlavor catalog.FlavorType) money.Cents {
	price := m.CurrentPrice(product)

	if extraFlavor != catalog.FlavorNone && product.AllowsExtraFlavor {
		price += m.extraFlavorPrice
	}

	return price
}

func (m *Machine) CanPrepare(
	product catalog.Product,
	selectedMilk catalog.MilkType,
	extraFlavor catalog.FlavorType,
) error {
	if m.maintenance.NeedsService() {
		return ErrMaintenanceRequired
	}

	recipe := product.Recipe

	if m.inventory.WaterML() < recipe.WaterML {
		return ErrInsufficientWater
	}

	if m.inventory.BeansG() < recipe.BeansG {
		return ErrInsufficientBeans
	}

	if m.inventory.Cups() < 1 {
		return ErrInsufficientCups
	}

	resolvedMilk, err := m.resolveMilk(product, selectedMilk)
	if err != nil {
		return err
	}

	if recipe.MilkML > 0 && m.inventory.MilkAmount(resolvedMilk) < recipe.MilkML {
		return fmt.Errorf("%w: %s", ErrInsufficientMilk, resolvedMilk.String())
	}

	if recipe.FlavorML > 0 && product.DefaultFlavor != catalog.FlavorNone {
		if m.inventory.FlavorAmount(product.DefaultFlavor) < recipe.FlavorML {
			return fmt.Errorf("%w: %s", ErrInsufficientFlavor, product.DefaultFlavor.String())
		}
	}

	if extraFlavor != catalog.FlavorNone {
		if !product.AllowsExtraFlavor {
			return ErrExtraFlavorNotAllowed
		}

		if m.inventory.FlavorAmount(extraFlavor) < m.extraFlavorML {
			return fmt.Errorf("%w: %s", ErrInsufficientFlavor, extraFlavor.String())
		}
	}

	return nil
}

func (m *Machine) Purchase(
	product catalog.Product,
	selectedMilk catalog.MilkType,
	extraFlavor catalog.FlavorType,
	session *money.PaymentSession,
) (map[money.Denomination]int, error) {
	if err := m.CanPrepare(product, selectedMilk, extraFlavor); err != nil {
		return nil, err
	}

	price := m.PriceFor(product, extraFlavor)
	total := session.Total()

	if total < price {
		return nil, money.ErrInsufficientPayment
	}

	changeAmount := total - price
	available := mergeCounts(m.cashBox.Snapshot(), session.Snapshot())

	change := make(map[money.Denomination]int)
	if changeAmount > 0 {
		var err error
		change, err = money.MakeChange(changeAmount, available)
		if err != nil {
			return nil, err
		}
	}

	if err := session.Commit(&m.cashBox); err != nil {
		return nil, err
	}

	for denomination, count := range change {
		if err := m.cashBox.Remove(denomination, count); err != nil {
			return nil, err
		}
	}

	if err := m.brew(product, selectedMilk, extraFlavor); err != nil {
		return nil, err
	}

	return change, nil
}

func (m *Machine) brew(
	product catalog.Product,
	selectedMilk catalog.MilkType,
	extraFlavor catalog.FlavorType,
) error {
	recipe := product.Recipe

	resolvedMilk, err := m.resolveMilk(product, selectedMilk)
	if err != nil {
		return err
	}

	m.inventory.ConsumeWaterML(recipe.WaterML)
	m.inventory.ConsumeBeansG(recipe.BeansG)
	m.inventory.ConsumeCup()

	if recipe.MilkML > 0 {
		m.inventory.ConsumeMilkML(resolvedMilk, recipe.MilkML)
	}

	if recipe.FlavorML > 0 && product.DefaultFlavor != catalog.FlavorNone {
		m.inventory.ConsumeFlavorML(product.DefaultFlavor, recipe.FlavorML)
	}

	if extraFlavor != catalog.FlavorNone && product.AllowsExtraFlavor {
		m.inventory.ConsumeFlavorML(extraFlavor, m.extraFlavorML)
	}

	m.maintenance.Increase(m.wearPerDrink)

	return nil
}

func (m *Machine) resolveMilk(
	product catalog.Product,
	selectedMilk catalog.MilkType,
) (catalog.MilkType, error) {
	if product.Recipe.MilkML == 0 {
		return catalog.MilkNone, nil
	}

	if selectedMilk == catalog.MilkNone {
		return product.DefaultMilk, nil
	}

	if !product.AllowsMilkChoice && selectedMilk != product.DefaultMilk {
		return catalog.MilkNone, ErrMilkChoiceNotAllowed
	}

	return selectedMilk, nil
}

func (m *Machine) PerformMaintenance() {
	m.maintenance.Reset()
}

func (m *Machine) AddWaterJugs(units int, waterJugML int) {
	if units > 0 {
		m.inventory.AddWaterML(units * waterJugML)
	}
}

func (m *Machine) AddBeanBags(units int, beanBagG int) {
	if units > 0 {
		m.inventory.AddBeansG(units * beanBagG)
	}
}

func (m *Machine) AddCupPacks(units int, cupPackCount int) {
	if units > 0 {
		m.inventory.AddCups(units * cupPackCount)
	}
}

func (m *Machine) AddMilkCartons(milkType catalog.MilkType, units int, milkCartonML int) {
	if units > 0 {
		m.inventory.AddMilkML(milkType, units*milkCartonML)
	}
}

func (m *Machine) AddFlavorBottles(flavorType catalog.FlavorType, units int, flavorBottleML int) {
	if units > 0 {
		m.inventory.AddFlavorML(flavorType, units*flavorBottleML)
	}
}

func (m *Machine) CollectMoney() (money.Cents, error) {
	total := m.cashBox.Total()
	snapshot := m.cashBox.Snapshot()

	for denomination, count := range snapshot {
		if count > 0 {
			if err := m.cashBox.Remove(denomination, count); err != nil {
				return 0, err
			}
		}
	}

	return total, nil
}

func mergeCounts(
	left map[money.Denomination]int,
	right map[money.Denomination]int,
) map[money.Denomination]int {
	merged := make(map[money.Denomination]int, len(left)+len(right))

	for denomination, count := range left {
		merged[denomination] = count
	}

	for denomination, count := range right {
		merged[denomination] += count
	}

	return merged
}
