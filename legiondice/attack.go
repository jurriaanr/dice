package dice

import (
	"math/rand"
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
	tokens        AttackTokens
	keywords      AttackKeywords
}

type AttackKeywords struct {
	blast         bool
	highVelocity  bool
	pierceX       int
	impactX       int
	ramX          int
	preciseX      int
	criticalX     int
	sharpshooterX int
}

type AttackTokens struct {
	aim int
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
	result := AttackResult{}

	for i, rb, rbw := 0, redDice+blackDice, redDice+blackDice+whiteDice; i < rbw; i++ {
		switch {
		case i < redDice:
			d := redAttackDice()
			switch d {
				case "H":
					result.Red.H++
				case "C":
					result.Red.C++
				case "S":
					result.Red.S++
				case "N":
					result.Red.N++
			}
		case i < rb:
			d := blackAttackDice()
			switch d {
			case "H":
				result.Black.H++
			case "C":
				result.Black.C++
			case "S":
				result.Black.S++
			case "N":
				result.Black.N++
			}
		case i < rbw:
			d := whiteAttackDice()
			switch d {
			case "H":
				result.White.H++
			case "C":
				result.White.C++
			case "S":
				result.White.S++
			case "N":
				result.White.N++
			}
		}
	}

	return result
}

func AttackRollResult(attack *Attack, defense *Defense) (hits int, result AttackResult, resultAfter AttackResult) {
	result = AttackRoll(attack.red, attack.black, attack.white)
	hits, resultAfter = CalculateHits(result, attack, defense)
	return hits, result, resultAfter
}
