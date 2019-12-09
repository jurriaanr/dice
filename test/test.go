package main

import (
	"fmt"
	dice "legion/legiondice"
)

func main() {
	attack := dice.CreateAttack(
		5,
		0,
		0,
		"crits",
	)

	//dice.AddAimToAttack(1, &attack)

	defense := dice.CreateDefense("red", true)
	result := dice.Test(attack, defense, 5, 10)

	fmt.Printf("This attack would result in %f hits", result.Successes)
}
