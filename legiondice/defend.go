package dice

import (
	"math/rand"
	"strings"
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
	redDiceValues := make([]string, redDice)
	whiteDiceValues := make([]string, whiteDice)

	for i, rw := 0, redDice+whiteDice; i < rw; i++ {
		switch {
		case i < redDice:
			redDiceValues[i] = redDefenseDice()
		case i < rw:
			whiteDiceValues[i-redDice] = whiteDefenseDice()
		}
	}

	return DefenseResult{
		Red: DefenseDiceResult{
			B: strings.Count(strings.Join(redDiceValues[:], ""), "B"),
			S: strings.Count(strings.Join(redDiceValues[:], ""), "S"),
			N: strings.Count(strings.Join(redDiceValues[:], ""), "N"),
		},
		White: DefenseDiceResult{
			B: strings.Count(strings.Join(whiteDiceValues[:], ""), "B"),
			S: strings.Count(strings.Join(whiteDiceValues[:], ""), "S"),
			N: strings.Count(strings.Join(whiteDiceValues[:], ""), "N"),
		},
	}
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
