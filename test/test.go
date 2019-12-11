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
		"none",
	)

	//dice.AddAimToAttack(2, &attack)
	//dice.AddPreciseXToAttack(1, &attack)

	defense := dice.CreateDefense("white", false, 0)

	//dice.AddDodgeToDefense(1, &defense)
	//dice.AddCoverXToDefense(1, &defense)
	//dice.AddPierceXToAttack(1, &attack)
	//dice.AddCriticalXToAttack(1, &attack)
	//dice.AddRamXToAttack(1, &attack)

	result := dice.Test(&attack, &defense, 100000, 10)

	fmt.Printf("This attack would result in %f hits", result.Successes)
}
