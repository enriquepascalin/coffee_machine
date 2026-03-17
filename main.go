package main

import "fmt"

type Recipe struct {
	water int
	milk  int
	beans int
	price int
	name  string
}

type CoffeeMachine struct {
	water int
	milk  int
	beans int
	cups  int
	money int
}

func (cm *CoffeeMachine) printState() {
	fmt.Println("The coffee machine has:")
	fmt.Printf("%d ml of water\n", cm.water)
	fmt.Printf("%d ml of milk\n", cm.milk)
	fmt.Printf("%d g of coffee beans\n", cm.beans)
	fmt.Printf("%d disposable cups\n", cm.cups)
	fmt.Printf("$%d of money\n", cm.money)
	fmt.Println()
}

func getRecipe(choice string) (Recipe, bool) {
	switch choice {
	case "1":
		return Recipe{water: 250, milk: 0, beans: 16, price: 4, name: "Espresso"}, true
	case "2":
		return Recipe{water: 350, milk: 75, beans: 20, price: 7, name: "Latte"}, true
	case "3":
		return Recipe{water: 200, milk: 100, beans: 12, price: 6, name: "Cappuccino"}, true
	default:
		return Recipe{}, false
	}
}

func (cm *CoffeeMachine) canMake(recipe Recipe) bool {
	switch {
	case cm.water < recipe.water:
		fmt.Println("Sorry, not enough water!")
		return false
	case cm.milk < recipe.milk:
		fmt.Println("Sorry, not enough milk!")
		return false
	case cm.beans < recipe.beans:
		fmt.Println("Sorry, not enough coffee beans!")
		return false
	case cm.cups < 1:
		fmt.Println("Sorry, not enough disposable cups!")
		return false
	default:
		fmt.Println("I have enough resources, making you a coffee!")
		return true
	}
}

func (cm *CoffeeMachine) buy(choice string) {
	recipe, ok := getRecipe(choice)
	if !ok {
		return
	}

	if !cm.canMake(recipe) {
		return
	}

	cm.water -= recipe.water
	cm.milk -= recipe.milk
	cm.beans -= recipe.beans
	cm.cups--
	cm.money += recipe.price
}

func (cm *CoffeeMachine) fill() {
	var waterToAdd, milkToAdd, beansToAdd, cupsToAdd int

	fmt.Println("Write how many ml of water you want to add:")
	fmt.Scan(&waterToAdd)

	fmt.Println("Write how many ml of milk you want to add:")
	fmt.Scan(&milkToAdd)

	fmt.Println("Write how many grams of coffee beans you want to add:")
	fmt.Scan(&beansToAdd)

	fmt.Println("Write how many disposable cups you want to add:")
	fmt.Scan(&cupsToAdd)

	cm.water += waterToAdd
	cm.milk += milkToAdd
	cm.beans += beansToAdd
	cm.cups += cupsToAdd
}

func (cm *CoffeeMachine) take() {
	fmt.Printf("I gave you $%d\n", cm.money)
	fmt.Println()
	cm.money = 0
}

func main() {
	coffeeMachine := CoffeeMachine{
		water: 400,
		milk:  540,
		beans: 120,
		cups:  9,
		money: 550,
	}

	for {
		var action string

		fmt.Println("Write action (buy, fill, take, remaining, exit):")
		fmt.Scan(&action)

		switch action {
		case "buy":
			var choice string
			fmt.Println("What do you want to buy? 1 - espresso, 2 - latte, 3 - cappuccino, back - to main menu:")
			fmt.Scan(&choice)

			if choice == "back" {
				fmt.Println()
				continue
			}

			coffeeMachine.buy(choice)
			fmt.Println()

		case "fill":
			coffeeMachine.fill()
			fmt.Println()

		case "take":
			coffeeMachine.take()

		case "remaining":
			coffeeMachine.printState()

		case "exit":
			return
		}
	}
}
