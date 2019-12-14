package dice

import (
	"math/rand"
)

var redDefenseValues = [...]string{"B", "B", "B", "S", "N", "N"}
var whiteDefenseValues = [...]string{"B", "S", "N", "N", "N", "N"}

type DefenseDiceResult struct {
	B int
	S int
	N int
}

type DefenseResult struct {
	Red   DefenseDiceResult
	White DefenseDiceResult
}

type Defense struct {
	config DefenseConfig
}

type DefenseConfig struct {
	rollsRedDefense   bool
	rollsWhiteDefense bool
	surgesToBlock     bool
	cover             int
	keywords          DefenseKeywords
	tokens            DefenseTokens
}

type DefenseKeywords struct {
	armor        bool
	impervious   bool
	lowProfile   bool
	armorX       int
	coverX       int
	uncannyLuckX int
}

type DefenseTokens struct {
	dodge int
}

func redDefenseDice() string {
	return redDefenseValues[rand.Intn(len(redDefenseValues))]
}

func whiteDefenseDice() string {
	return whiteDefenseValues[rand.Intn(len(whiteDefenseValues))]
}

func DefenseRoll(redDice, whiteDice int) DefenseResult {
	result := DefenseResult{}

	for i, rw := 0, redDice+whiteDice; i < rw; i++ {
		switch {
		case i < redDice:
			d := redDefenseDice()
			switch d {
			case "B":
				result.Red.B++
			case "S":
				result.Red.S++
			case "N":
				result.Red.N++
			}
		case i < rw:
			d := whiteDefenseDice()
			switch d {
			case "B":
				result.White.B++
			case "S":
				result.White.S++
			case "N":
				result.White.N++
			}
		}
	}

	return result
}

func DefenseRoleResult(hits int, attack *Attack, defense *Defense) (blocks int, result DefenseResult, resultAfter DefenseResult) {
	redDice := 0
	whiteDice := 0

	hits = addImperviousToDefense(hits, attack, defense)

	if defense.config.rollsRedDefense {
		redDice = hits
	} else {
		whiteDice = hits
	}

	// 7a Roll Dice
	//  For each hit and critical result on the attacker’s dice, the defender rolls one defense die whose
	// color matches the defender’s defense, which is presented on the defender’s unit card.
	result = DefenseRoll(redDice, whiteDice)
	blocks, resultAfter = CalculateBlocks(result, attack, defense)

	return blocks, result, resultAfter
}
