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
	dice   int
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

func CalculateBlocks(result *DefenseResult, defense *Defense) int {
	val := result.Red.B + result.White.B
	if defense.config.surgesToBlock {
		val += result.Red.S + result.White.S
	}
	return val
}

func DefenseRoleResult(defense *Defense) (int, DefenseResult) {
	redDice := 0
	whiteDice := 0
	if defense.config.rollsRedDefense {
		redDice = defense.dice
	} else {
		whiteDice = defense.dice
	}

	result := DefenseRoll(redDice, whiteDice)
	return CalculateBlocks(&result, defense), result
}

func DefenseTest(defense *Defense, rolls int) float64 {
	sum := 0
	for i := 0; i < rolls; i++ {
		blocks, _ := DefenseRoleResult(defense)
		sum += blocks
	}
	return float64(sum) / float64(rolls)
}
