package ui

import (
	"fmt"
	"io"
	"sort"

	"github.com/enriquepascalin/coffee_machine/internal/catalog"
	"github.com/enriquepascalin/coffee_machine/internal/machine"
	"github.com/enriquepascalin/coffee_machine/internal/money"
)

type Printer struct {
	out io.Writer
}

func NewPrinter(out io.Writer) *Printer {
	return &Printer{out: out}
}

func (p *Printer) Line(message string) {
	fmt.Fprintln(p.out, message)
}

func (p *Printer) Blank() {
	fmt.Fprintln(p.out)
}

func (p *Printer) Header(message string) {
	fmt.Fprintln(p.out, colorize(ansiBold+ansiBlue, message))
}

func (p *Printer) Info(message string) {
	fmt.Fprintln(p.out, colorize(ansiCyan, message))
}

func (p *Printer) Success(message string) {
	fmt.Fprintln(p.out, colorize(ansiGreen, message))
}

func (p *Printer) Warn(message string) {
	fmt.Fprintln(p.out, colorize(ansiYellow, message))
}

func (p *Printer) Error(message string) {
	fmt.Fprintln(p.out, colorize(ansiRed, message))
}

func (p *Printer) Prompt(message string) {
	fmt.Fprint(p.out, colorize(ansiBold, message))
}

func (p *Printer) PrintCatalog(products []catalog.Product, priceFn func(catalog.Product) money.Cents) {
	p.Header("Products")

	for _, product := range products {
		price := priceFn(product)
		fmt.Fprintf(
			p.out,
			"  %d. %s %s - %s\n",
			product.ID,
			product.Emoji,
			product.Name,
			price.String(),
		)
	}

	p.Blank()
}

func (p *Printer) PrintInventory(
	inventory machine.InventorySnapshot,
	cash money.Cents,
	maintenanceGauge int,
) {
	p.Header("Machine Status")
	fmt.Fprintf(p.out, "Water: %d ml\n", inventory.WaterML)
	fmt.Fprintf(p.out, "Beans: %d g\n", inventory.BeansG)
	fmt.Fprintf(p.out, "Cups: %d\n", inventory.Cups)

	p.Line("Milk:")
	milkTypes := []catalog.MilkType{catalog.MilkWhole, catalog.MilkAlmond, catalog.MilkOat}
	for _, milkType := range milkTypes {
		fmt.Fprintf(p.out, "  - %s: %d ml\n", milkType.String(), inventory.MilkML[milkType])
	}

	p.Line("Flavor:")
	flavorTypes := []catalog.FlavorType{catalog.FlavorVanilla, catalog.FlavorCaramel, catalog.FlavorChocolate}
	for _, flavorType := range flavorTypes {
		fmt.Fprintf(p.out, "  - %s: %d ml\n", flavorType.String(), inventory.FlavorML[flavorType])
	}

	fmt.Fprintf(p.out, "Cash: %s\n", cash.String())
	fmt.Fprintf(p.out, "Maintenance: %d%%\n", maintenanceGauge)
	p.Blank()
}

func (p *Printer) PrintChange(change map[money.Denomination]int) {
	if len(change) == 0 {
		return
	}

	p.Header("Change Returned")

	denominations := make([]int, 0, len(change))
	for denomination := range change {
		denominations = append(denominations, int(denomination))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(denominations)))

	for _, denominationValue := range denominations {
		denomination := money.Denomination(denominationValue)
		count := change[denomination]
		fmt.Fprintf(p.out, "  %s x %d\n", denomination.Cents().String(), count)
	}

	p.Blank()
}
