package main

import (
	"fmt"
	dice "legion/legiondice"
)

func main() {
	attack := dice.CreateAttack(
		3,
		0,
		4,
		"none",
	)

	dice.AddAimToAttack(8, &attack)
	//dice.AddPreciseXToAttack(1, &attack)

	defense := dice.CreateDefense("white", false, 0)

	//dice.AddSurgeToAttack(2,  &attack)
	//dice.AddShieldToDefense(3, &defense)
	//dice.AddDodgeToDefense(1, &defense)
	//dice.AddCoverXToDefense(1, &defense)
	//dice.AddPierceXToAttack(1, &attack)
	//dice.AddCriticalXToAttack(1, &attack)
	//dice.AddRamXToAttack(1, &attack)
	//dice.AddImpactXToAttack(1, &attack)

	result := dice.Test(&attack, &defense, 100000, 10)

	defer fmt.Printf("This attack would result in %f hits", result.Successes)
}
