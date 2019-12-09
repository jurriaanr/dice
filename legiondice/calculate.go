package dice

func CalculateHits(result AttackResult, attack *Attack, defense *Defense) int {
	// 4B Reroll attack dice
	misses := getAttackMisses(&result, attack)

	if attack.config.tokens.aim > 0 && misses > 0 {
		red, black, white := getAttackDicesToReroll(&result, attack, misses)
		extraResult := AttackRoll(red, black, white)
		combineAttackResults(&result, &extraResult)
	}

	// 4C Convert attack surges
	applySurges(&result, attack)

	// 5 Apply Dodge & Cover
	applyDodgeAndCover(&result, defense)

	// count hits
	val := result.Red.H + result.Red.C + result.Black.H + result.Black.C + result.White.H + result.White.C

	return val
}

func getAttackMisses(result *AttackResult, attack *Attack) int {
	misses := result.Red.N + result.Black.N + result.White.N

	if !attack.config.surgesToCrits && !attack.config.surgesToHits {
		misses += result.Red.S + result.Black.S + result.White.S
	}

	return misses
}

func getAttackDicesToReroll(result *AttackResult, attack *Attack, misses int) (red int, black int, white int) {
	// we can only reroll up to the number of aimtokens we have. So either all the misses, or the max allowed rerolls
	count := min(misses, attack.config.tokens.aim*(2+attack.config.keywords.preciseX))
	convertsSurges := attack.config.surgesToHits || attack.config.surgesToCrits

	redToReroll := 0
	blackToReroll := 0
	whiteToReroll := 0

	// subtract from original result. Start with red, because it has the most chance of an extra hit
	for tot := 0; tot < count; {
		if result.Red.N > 0 {
			redToReroll++
			result.Red.N--
		} else if result.Red.S > 0 && !convertsSurges {
			redToReroll++
			result.Red.S--
		} else if result.Black.N > 0 {
			blackToReroll++
			result.Black.N--
		} else if result.Black.S > 0 && !convertsSurges {
			blackToReroll++
			result.Black.S--
		} else if result.White.N > 0 {
			whiteToReroll++
			result.White.N--
		} else if result.White.S > 0 && !convertsSurges {
			whiteToReroll++
			result.White.S--
		}

		tot = redToReroll + blackToReroll + whiteToReroll
	}

	return redToReroll, blackToReroll, whiteToReroll
}

func applySurges(result *AttackResult, attack *Attack) {
	if attack.config.surgesToHits {
		result.Red.H += result.Red.S
		result.Black.H += result.Black.S
		result.White.H += result.White.S
		result.Red.S = 0
		result.Black.S = 0
		result.White.S = 0
	} else if attack.config.surgesToCrits {
		result.Red.C += result.Red.S
		result.Black.C += result.Black.S
		result.White.C += result.White.S
		result.Red.S = 0
		result.Black.S = 0
		result.White.S = 0
	}
}

func applyDodgeAndCover(result *AttackResult, defense *Defense) {
	hitsToRemove := min(defense.config.cover+defense.config.keywords.coverX, 2)
	hitsToRemove += defense.config.tokens.dodge

	for hits := result.White.H + result.Black.H + result.Red.H; hitsToRemove > 0 && hits > 0; {
		if result.White.H > 0 {
			result.White.H--
			hitsToRemove--
		} else if result.Black.H > 0 {
			result.Black.H--
			hitsToRemove--
		} else if result.Red.H > 0 {
			result.Red.H--
			hitsToRemove--
		}

		hits = result.White.H + result.Black.H + result.Red.H
	}
}

func combineAttackResults(a *AttackResult, b *AttackResult) {
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
}
