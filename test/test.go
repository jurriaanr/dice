package main

import (
	"fmt"
	dice "legion/legiondice"
)

func main() {
	attack := dice.CreateAttack(
		3,
		0,
		0,
		"crits",
	)
	defense := dice.CreateDefense("red", true)
	result := dice.Test(attack, defense, 100000, 10)

	fmt.Print(result)
}
