package app

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/enriquepascalin/coffee_machine/internal/auth"
	"github.com/enriquepascalin/coffee_machine/internal/catalog"
	"github.com/enriquepascalin/coffee_machine/internal/config"
	"github.com/enriquepascalin/coffee_machine/internal/machine"
	"github.com/enriquepascalin/coffee_machine/internal/money"
	"github.com/enriquepascalin/coffee_machine/internal/ui"
)

type App struct {
	machine       *machine.Machine
	authenticator auth.Authenticator
	input         *ui.Input
	printer       *ui.Printer
	mode          Mode
}

func New(
	machine *machine.Machine,
	authenticator auth.Authenticator,
	input *ui.Input,
	printer *ui.Printer,
) *App {
	return &App{
		machine:       machine,
		authenticator: authenticator,
		input:         input,
		printer:       printer,
		mode:          ModeUser,
	}
}

func (a *App) Run() error {
	a.printer.Header("Enterprise Coffee Machine")

	for {
		var (
			exit bool
			err  error
		)

		switch a.mode {
		case ModeUser:
			exit, err = a.runUserMode()
		case ModeAdmin:
			exit, err = a.runAdminMode()
		default:
			return fmt.Errorf("unknown app mode: %d", a.mode)
		}

		if err != nil {
			a.printer.Error(err.Error())
		}

		if exit {
			a.printer.Info("Goodbye.")
			return nil
		}
	}
}

func (a *App) runUserMode() (bool, error) {
	a.printer.Prompt("user> [menu|buy|remaining|admin|exit]: ")

	command, err := a.input.ReadLine()
	if err != nil {
		return false, err
	}

	switch strings.ToLower(command) {
	case "menu":
		a.printer.PrintCatalog(a.machine.Products(), a.machine.CurrentPrice)
		return false, nil
	case "buy":
		return false, a.handleBuy()
	case "remaining":
		a.printer.PrintInventory(a.machine.InventorySnapshot(), a.machine.CashTotal(), a.machine.MaintenanceGauge())
		return false, nil
	case "admin":
		return false, a.handleAdminLogin()
	case "exit":
		return true, nil
	default:
		return false, fmt.Errorf("unknown command: %s", command)
	}
}

func (a *App) runAdminMode() (bool, error) {
	a.printer.Prompt("admin> [remaining|refill|prices|take|maintenance|logout|exit]: ")

	command, err := a.input.ReadLine()
	if err != nil {
		return false, err
	}

	switch strings.ToLower(command) {
	case "remaining":
		a.printer.PrintInventory(a.machine.InventorySnapshot(), a.machine.CashTotal(), a.machine.MaintenanceGauge())
		return false, nil
	case "refill":
		return false, a.handleRefill()
	case "prices":
		return false, a.handlePriceChange()
	case "take":
		amount, err := a.machine.CollectMoney()
		if err != nil {
			return false, err
		}
		a.printer.Success("Collected " + amount.String())
		return false, nil
	case "maintenance":
		a.machine.PerformMaintenance()
		a.printer.Success("Maintenance completed.")
		return false, nil
	case "logout":
		a.mode = ModeUser
		a.printer.Info("Switched to user mode.")
		return false, nil
	case "exit":
		return true, nil
	default:
		return false, fmt.Errorf("unknown command: %s", command)
	}
}

func (a *App) handleAdminLogin() error {
	a.printer.Prompt("password: ")

	password, err := a.input.ReadLine()
	if err != nil {
		return err
	}

	if !a.authenticator.Authenticate(password) {
		return errors.New("invalid admin password")
	}

	a.mode = ModeAdmin
	a.printer.Success("Admin mode enabled.")
	return nil
}

func (a *App) handleBuy() error {
	product, ok, err := a.promptProduct()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	selectedMilk, err := a.promptMilk(product)
	if err != nil {
		return err
	}

	selectedFlavor, err := a.promptFlavor(product)
	if err != nil {
		return err
	}

	if err := a.machine.CanPrepare(product, selectedMilk, selectedFlavor); err != nil {
		return err
	}

	price := a.machine.PriceFor(product, selectedFlavor)
	session, ok, err := a.collectPayment(price)
	if err != nil {
		return err
	}
	if !ok {
		a.printer.Warn("Purchase cancelled.")
		return nil
	}

	change, err := a.machine.Purchase(product, selectedMilk, selectedFlavor, session)
	if err != nil {
		return err
	}

	a.printer.Info("Preparing " + product.Name + " " + product.Emoji)
	a.printer.Progress("Brewing", product.BrewSeconds)

	if len(change) > 0 {
		a.printer.PrintChange(change)
	}

	a.printer.Success("Enjoy your " + product.Name + " " + product.Emoji)
	return nil
}

func (a *App) handleRefill() error {
	waterJugs, err := a.promptNonNegativeInt("How many 20L water jugs to add? ")
	if err != nil {
		return err
	}

	beanBags, err := a.promptNonNegativeInt("How many 1kg bean bags to add? ")
	if err != nil {
		return err
	}

	cupPacks, err := a.promptNonNegativeInt("How many 50-cup packs to add? ")
	if err != nil {
		return err
	}

	wholeMilkCartons, err := a.promptNonNegativeInt("How many whole milk cartons (1L) to add? ")
	if err != nil {
		return err
	}

	almondMilkCartons, err := a.promptNonNegativeInt("How many almond milk cartons (1L) to add? ")
	if err != nil {
		return err
	}

	oatMilkCartons, err := a.promptNonNegativeInt("How many oat milk cartons (1L) to add? ")
	if err != nil {
		return err
	}

	vanillaBottles, err := a.promptNonNegativeInt("How many vanilla bottles (750ml) to add? ")
	if err != nil {
		return err
	}

	caramelBottles, err := a.promptNonNegativeInt("How many caramel bottles (750ml) to add? ")
	if err != nil {
		return err
	}

	chocolateBottles, err := a.promptNonNegativeInt("How many chocolate bottles (750ml) to add? ")
	if err != nil {
		return err
	}

	a.machine.AddWaterJugs(waterJugs, config.WaterJugML)
	a.machine.AddBeanBags(beanBags, config.BeanBagG)
	a.machine.AddCupPacks(cupPacks, config.CupPackCount)

	a.machine.AddMilkCartons(catalog.MilkWhole, wholeMilkCartons, config.MilkCartonML)
	a.machine.AddMilkCartons(catalog.MilkAlmond, almondMilkCartons, config.MilkCartonML)
	a.machine.AddMilkCartons(catalog.MilkOat, oatMilkCartons, config.MilkCartonML)

	a.machine.AddFlavorBottles(catalog.FlavorVanilla, vanillaBottles, config.FlavorBottleML)
	a.machine.AddFlavorBottles(catalog.FlavorCaramel, caramelBottles, config.FlavorBottleML)
	a.machine.AddFlavorBottles(catalog.FlavorChocolate, chocolateBottles, config.FlavorBottleML)

	a.printer.Success("Machine refilled.")
	return nil
}

func (a *App) handlePriceChange() error {
	a.printer.PrintCatalog(a.machine.Products(), a.machine.CurrentPrice)

	productIDValue, err := a.promptNonNegativeInt("Product ID to update: ")
	if err != nil {
		return err
	}

	product, ok := a.machine.FindProduct(catalog.ProductID(productIDValue))
	if !ok {
		return fmt.Errorf("product %d not found", productIDValue)
	}

	newPriceCentsValue, err := a.promptNonNegativeInt("New price in cents: ")
	if err != nil {
		return err
	}

	if err := a.machine.ChangePrice(product.ID, money.Cents(newPriceCentsValue)); err != nil {
		return err
	}

	a.printer.Success("Updated price for " + product.Name)
	return nil
}

func (a *App) promptProduct() (catalog.Product, bool, error) {
	a.printer.PrintCatalog(a.machine.Products(), a.machine.CurrentPrice)
	a.printer.Prompt("Select product ID or type 'back': ")

	line, err := a.input.ReadLine()
	if err != nil {
		return catalog.Product{}, false, err
	}

	if strings.EqualFold(line, "back") {
		return catalog.Product{}, false, nil
	}

	idValue, err := strconv.Atoi(line)
	if err != nil {
		return catalog.Product{}, false, errors.New("invalid product selection")
	}

	product, ok := a.machine.FindProduct(catalog.ProductID(idValue))
	if !ok {
		return catalog.Product{}, false, fmt.Errorf("product %d not found", idValue)
	}

	return product, true, nil
}

func (a *App) promptMilk(product catalog.Product) (catalog.MilkType, error) {
	if !product.AllowsMilkChoice {
		return product.DefaultMilk, nil
	}

	a.printer.Line("Milk options:")
	a.printer.Line("0. default")
	a.printer.Line("1. whole")
	a.printer.Line("2. almond")
	a.printer.Line("3. oat")
	a.printer.Prompt("Select milk option: ")

	line, err := a.input.ReadLine()
	if err != nil {
		return catalog.MilkNone, err
	}

	switch strings.TrimSpace(line) {
	case "", "0":
		return product.DefaultMilk, nil
	case "1":
		return catalog.MilkWhole, nil
	case "2":
		return catalog.MilkAlmond, nil
	case "3":
		return catalog.MilkOat, nil
	default:
		return catalog.MilkNone, errors.New("invalid milk option")
	}
}

func (a *App) promptFlavor(product catalog.Product) (catalog.FlavorType, error) {
	if !product.AllowsExtraFlavor {
		return catalog.FlavorNone, nil
	}

	a.printer.Line("Flavor options:")
	a.printer.Line("0. none")
	a.printer.Line("1. vanilla")
	a.printer.Line("2. caramel")
	a.printer.Line("3. chocolate")
	a.printer.Prompt("Select extra flavor: ")

	line, err := a.input.ReadLine()
	if err != nil {
		return catalog.FlavorNone, err
	}

	switch strings.TrimSpace(line) {
	case "", "0":
		return catalog.FlavorNone, nil
	case "1":
		return catalog.FlavorVanilla, nil
	case "2":
		return catalog.FlavorCaramel, nil
	case "3":
		return catalog.FlavorChocolate, nil
	default:
		return catalog.FlavorNone, errors.New("invalid flavor option")
	}
}

func (a *App) collectPayment(price money.Cents) (*money.PaymentSession, bool, error) {
	session := money.NewPaymentSession()

	for {
		total := session.Total()
		if total >= price {
			return session, true, nil
		}

		remaining := price - total

		a.printer.Info("Inserted: " + total.String())
		a.printer.Info("Remaining: " + remaining.String())
		a.printer.Line("Accepted denominations (cents): 2000 1000 500 200 100 50 25 10 5")
		a.printer.Prompt("Insert denomination or type 'cancel': ")

		line, err := a.input.ReadLine()
		if err != nil {
			return nil, false, err
		}

		if strings.EqualFold(line, "cancel") {
			return nil, false, nil
		}

		denominationValue, err := strconv.Atoi(line)
		if err != nil {
			a.printer.Error("invalid denomination")
			continue
		}

		if err := session.Insert(money.Denomination(denominationValue), 1); err != nil {
			a.printer.Error(err.Error())
			continue
		}
	}
}

func (a *App) promptInt(prompt string) (int, error) {
	a.printer.Prompt(prompt)
	return a.input.ReadInt()
}

func (a *App) promptNonNegativeInt(prompt string) (int, error) {
	value, err := a.promptInt(prompt)
	if err != nil {
		return 0, err
	}

	if value < 0 {
		return 0, errors.New("value cannot be negative")
	}

	return value, nil
}
