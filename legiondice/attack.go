package dice

import (
	"math/rand"
	"strings"
)

var redAttackValues = [...]string{"H", "H", "H", "H", "H", "C", "S", "N"}
var blackAttackValues = [...]string{"H", "H", "H", "C", "S", "N", "N", "N"}
var whiteAttackValues = [...]string{"H", "C", "S", "N", "N", "N", "N", "N"}

type AttackDiceResult struct {
	H int
	C int
	S int
	N int
}

type AttackResult struct {
	Red   AttackDiceResult
	Black AttackDiceResult
	White AttackDiceResult
}

type Attack struct {
	red    int
	black  int
	white  int
	config AttackConfig
}

type AttackConfig struct {
	surgesToHits  bool
	surgesToCrits bool
	keywords      AttackKeywords
}

type AttackKeywords struct {
	pierceX  int
	impactX  int
	ramX     int
	aim      bool
	preciseX int
}

func redAttackDice() string {
	return redAttackValues[rand.Intn(len(redAttackValues))]
}

func blackAttackDice() string {
	return blackAttackValues[rand.Intn(len(blackAttackValues))]
}

func whiteAttackDice() string {
	return whiteAttackValues[rand.Intn(len(whiteAttackValues))]
}

func AttackRoll(redDice, blackDice, whiteDice int) AttackResult {
	redDiceValues := make([]string, redDice)
	blackDiceValues := make([]string, blackDice)
	whiteDiceValues := make([]string, whiteDice)

	for i, rb, rbw := 0, redDice+blackDice, redDice+blackDice+whiteDice; i < rbw; i++ {
		switch {
		case i < redDice:
			redDiceValues[i] = redAttackDice()
		case i < rb:
			blackDiceValues[i-redDice] = blackAttackDice()
		case i < rbw:
			whiteDiceValues[i-rb] = whiteAttackDice()
		}
	}

	return AttackResult{
		Red: AttackDiceResult{
			H: strings.Count(strings.Join(redDiceValues[:], ""), "H"),
			C: strings.Count(strings.Join(redDiceValues[:], ""), "C"),
			S: strings.Count(strings.Join(redDiceValues[:], ""), "S"),
			N: strings.Count(strings.Join(redDiceValues[:], ""), "N"),
		},
		Black: AttackDiceResult{
			H: strings.Count(strings.Join(blackDiceValues[:], ""), "H"),
			C: strings.Count(strings.Join(blackDiceValues[:], ""), "C"),
			S: strings.Count(strings.Join(blackDiceValues[:], ""), "S"),
			N: strings.Count(strings.Join(blackDiceValues[:], ""), "N"),
		},
		White: AttackDiceResult{
			H: strings.Count(strings.Join(whiteDiceValues[:], ""), "H"),
			C: strings.Count(strings.Join(whiteDiceValues[:], ""), "C"),
			S: strings.Count(strings.Join(whiteDiceValues[:], ""), "S"),
			N: strings.Count(strings.Join(whiteDiceValues[:], ""), "N"),
		},
	}
}

func CalculateHits(result AttackResult, attack Attack, defense Defense) int {
	val := result.Red.H + result.Red.C + result.Black.H + result.Black.C + result.White.H + result.White.C
	if attack.config.surgesToCrits || attack.config.surgesToHits {
		val += result.Red.S + result.Black.S + result.White.S
	}
	return val
}

func AttackRollResult(attack Attack, defense Defense) (int, AttackResult) {
	result := AttackRoll(attack.red, attack.black, attack.white)
	return CalculateHits(result, attack, defense), result
}

func AttackTest(attack Attack, defense Defense, rolls int) float64 {
	sum := 0
	for i := 0; i < rolls; i++ {
		hits, _ := AttackRollResult(attack, defense)
		sum += hits
	}
	return float64(sum) / float64(rolls)
}
