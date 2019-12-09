package dice

import "fmt"

func CalculateHits(result AttackResult, attack Attack, defense Defense) int {
	misses := getAttackMisses(result, attack)

	if attack.config.tokens.aim > 0 && misses > 0 {
		fmt.Printf("Rerolling %d misses with %d aimtoken\n", misses, attack.config.tokens.aim)
		red, black, white := getAttackDicesToReroll(&result, attack, misses)
		extraResult := AttackRoll(red, black, white)
		result = combineAttackResults(result, extraResult)
	}

	val := result.Red.H + result.Red.C + result.Black.H + result.Black.C + result.White.H + result.White.C

	if attack.config.surgesToCrits || attack.config.surgesToHits {
		val += result.Red.S + result.Black.S + result.White.S
	}
	return val
}

func getAttackMisses(result AttackResult, attack Attack) int {
	misses := result.Red.N + result.Black.N + result.White.N

	if !attack.config.surgesToCrits && !attack.config.surgesToHits {
		misses += result.Red.S + result.Black.S + result.White.S
	}

	return misses
}

func getAttackDicesToReroll(result *AttackResult, attack Attack, misses int) (red int, black int, white int) {
	// we can only reroll up to the number of aimtokens we have. So either all the misses, or the max allowed rerolls
	count := min(misses, attack.config.tokens.aim * (2 + attack.config.keywords.preciseX))
	convertsSurges := attack.config.surgesToHits || attack.config.surgesToCrits

	whiteToReroll := 0
	blackToReroll := 0
	redToReroll := 0

	// subtract from original result
	for tot := 0; tot < count; {
		if result.White.N > 0 {
			whiteToReroll++
			result.White.N--
		} else if result.White.S > 0 && !convertsSurges {
			whiteToReroll++
			result.White.S--
		} else if result.Black.N > 0 {
			blackToReroll++
			result.Black.N--
		} else if result.Black.S > 0 && !convertsSurges {
			blackToReroll++
			result.Black.S--
		} else if result.Red.N > 0 {
			redToReroll++
			result.Red.N--
		} else if result.Red.S > 0 && !convertsSurges {
			redToReroll++
			result.Red.S--
		}

		tot = whiteToReroll + blackToReroll + redToReroll
	}

	return redToReroll, blackToReroll, whiteToReroll
}

func combineAttackResults(a AttackResult, b AttackResult) AttackResult {
	a.Red.H += b.Red.H
	a.Red.C += b.Red.C
	a.Red.S += b.Red.S
	a.Red.N += b.Red.N

	a.Black.H += b.Black.H
	a.Black.C += b.Black.C
	a.Black.S += b.Black.S
	a.Black.N += b.Black.N

	a.White.H += b.White.H
	a.White.C += b.White.C
	a.White.S += b.White.S
	a.White.N += b.White.N

	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}