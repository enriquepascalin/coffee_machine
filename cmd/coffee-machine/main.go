package main

import (
	"fmt"
	"os"

	"github.com/enriquepascalin/coffee_machine/internal/app"
	"github.com/enriquepascalin/coffee_machine/internal/auth"
	"github.com/enriquepascalin/coffee_machine/internal/config"
	"github.com/enriquepascalin/coffee_machine/internal/ui"
)

func main() {
	defaults := config.NewDefaults()

	application := app.New(
		defaults.Machine,
		auth.New(defaults.AdminPassword),
		ui.NewInput(os.Stdin),
		ui.NewPrinter(os.Stdout),
	)

	if err := application.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
