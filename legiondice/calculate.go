package dice

func CalculateHits(result AttackResult, attack *Attack, defense *Defense) (int, AttackResult) {
	// 4B Reroll Dice
	// The attacker can resolve any abilities that allow the attacker to reroll attack dice.
	misses := getAttackMisses(&result, attack)

	if attack.config.tokens.aim > 0 && misses > 0 {
		red, black, white := getAttackDicesToReroll(&result, attack, misses)
		extraResult := AttackRoll(red, black, white)
		combineAttackResults(&result, &extraResult)
	}

	// 4C Convert attack surges
	// The attacker changes its attack surge results to the result indicated on its unit card by turning the die.
	// If no result is indicated, the attacker changes the result to a blank.
	applyAttackSurges(&result, attack)

	// 5 Apply Dodge & Cover
	// If the defender has a dodge token for is in cover, the defender may spend dodge tokens and apply cover to
	// cancel hit results. Dodge tokens and cover cannot be used to cancel critical results.
	// A unit can apply cover only against ranged attacks
	applyDodgeAndCover(&result, attack, defense)

	// 6 Modify Attack Dice
	// The attacker can resolve any card abilities that modify the attack dice.
	// Then, the defender can resolve any card abilities that modify the attack dice
	// attack
	applyImpactX(&result, attack, defense)
	applyRamX(&result, attack)
	// defense
	applyArmor(&result, defense)

	// count hits
	val := result.Red.H + result.Red.C + result.Black.H + result.Black.C + result.White.H + result.White.C

	return val, result
}

func CalculateBlocks(result DefenseResult, attack *Attack, defense *Defense) (int, DefenseResult) {
	// add shield blocks
	applyShieldToDefenseResult(&result, defense)

	// 7b Reroll Dice
	// The defender can resolve any abilities that allow the defender to reroll defense dice.
	misses := getDefenseMisses(&result, defense)

	// While defending, you may reroll up to X defense dice
	if defense.config.keywords.uncannyLuckX > 0 && misses > 0 {
		red, white := getDefenseDicesToReroll(&result, defense, misses)
		extraResult := DefenseRoll(red, white)
		combineDefenseResults(&result, &extraResult)
	}

	// 7c Convert Defense Surges:
	// The defender changes its defense surge results to the result indicated on its unit card by turning the die.
	// If no result is indicated, the defender changes the result to a blank.
	applyDefenseSurges(&result, defense)

	// 8 Modify Defense Dice
	// The defender can resolve any card abilities that modify the defense dice.
	// Then, the attacker can resolve any card abilities that modify the defense dice
	applyPierce(&result, attack)

	return result.Red.B + result.White.B, result
}

func getAttackMisses(result *AttackResult, attack *Attack) int {
	misses := result.Red.N + result.Black.N + result.White.N

	if !attack.config.surgesToCrits && !attack.config.surgesToHits {
		misses += result.Red.S + result.Black.S + result.White.S
	}

	return misses
}

func getDefenseMisses(result *DefenseResult, defense *Defense) int {
	misses := result.Red.N + result.White.N

	if !defense.config.surgesToBlock {
		misses += result.Red.S + result.White.S
	}

	return misses
}

func getAttackDicesToReroll(result *AttackResult, attack *Attack, misses int) (red int, black int, white int) {
	// we can only reroll up to the number of aimtokens we have. So either all the misses, or the max allowed rerolls
	count := min(misses, attack.config.tokens.aim*(2+attack.config.keywords.preciseX))

	convertsSurges := attack.config.surgesToHits || attack.config.surgesToCrits

	// first we will temporarily remove surges that would be converted by other keywords like criticalX
	// But only if surges are not converted in general, other wise we do not care at this point
	whiteSurges, blackSurges, redSurges, whiteBlanks, blackBlanks, redBlanks := 0, 0, 0, 0, 0, 0

	if !convertsSurges {
		whiteSurges, blackSurges, redSurges, whiteBlanks, blackBlanks, redBlanks = saveDiceBeforeReroll(result, attack)
	}

	// subtract from original result. Start with red, because it has the most chance of an extra hit
	for tot := 0; tot < count; {
		if result.Red.N > 0 {
			red++
			result.Red.N--
		} else if result.Red.S > 0 && !convertsSurges {
			red++
			result.Red.S--
		} else if result.Black.N > 0 {
			black++
			result.Black.N--
		} else if result.Black.S > 0 && !convertsSurges {
			black++
			result.Black.S--
		} else if result.White.N > 0 {
			white++
			result.White.N--
		} else if result.White.S > 0 && !convertsSurges {
			white++
			result.White.S--
		}

		tot = red + black + white
	}

	// add back the savedSurges
	if !convertsSurges {
		result.White.S += whiteSurges
		result.White.N += whiteBlanks
		result.Black.S += blackSurges
		result.Black.N += blackBlanks
		result.Red.S += redSurges
		result.Red.N += redBlanks
	}

	return red, black, white
}

func getDefenseDicesToReroll(result *DefenseResult, defense *Defense, misses int) (red int, white int) {
	count := min(misses, defense.config.keywords.uncannyLuckX)

	for tot := 0; tot < count; {
		if result.Red.N > 0 {
			red++
			result.Red.N--
		} else if result.Red.S > 0 && !defense.config.surgesToBlock {
			red++
			result.Red.S--
		} else if result.White.N > 0 {
			white++
			result.White.N--
		} else if result.White.S > 0 && !defense.config.surgesToBlock {
			white++
			result.White.S--
		}

		tot = red + white
	}

	return red, white
}

func saveDiceBeforeReroll(result *AttackResult, attack *Attack) (whiteSurges, blackSurges, redSurges, whiteBlanks, blackBlanks, redBlanks int) {
	numberOfSurgesToSave := attack.config.keywords.criticalX
	numberOfSurgesOrBlanksToSave := attack.config.keywords.ramX

	whiteSurges = 0
	blackSurges = 0
	redSurges = 0
	whiteBlanks = 0
	blackBlanks = 0
	redBlanks = 0

	for numberOfSurgesToSave > 0 {
		// makes sense to keep the white surges first do better dice can be used for rerolls
		if result.White.S > 0 {
			numberOfSurgesToSave--
			whiteSurges++
		} else if result.Black.S > 0 {
			numberOfSurgesToSave--
			blackSurges++
		} else if result.Red.S > 0 {
			numberOfSurgesToSave--
			redSurges++
		} else {
			break
		}
	}

	for numberOfSurgesOrBlanksToSave > 0 {
		// makes sense to keep the white surges first do better dice can be used for rerolls
		if result.White.N > 0 {
			numberOfSurgesOrBlanksToSave--
			whiteBlanks++
		} else if result.Black.N > 0 {
			numberOfSurgesOrBlanksToSave--
			blackBlanks++
		} else if result.Red.N > 0 {
			numberOfSurgesOrBlanksToSave--
			redBlanks++
		} else if result.White.S > 0 && !attack.config.surgesToCrits {
			numberOfSurgesOrBlanksToSave--
			whiteSurges++
		} else if result.Black.S > 0 && !attack.config.surgesToCrits {
			numberOfSurgesOrBlanksToSave--
			blackSurges++
		} else if result.Red.S > 0 && !attack.config.surgesToCrits {
			numberOfSurgesOrBlanksToSave--
			redSurges++
		} else {
			break
		}
	}

	return whiteSurges, blackSurges, redSurges, whiteBlanks, blackBlanks, redBlanks
}

// 4C Convert attack surges
// The attacker changes its attack surge results to the result indicated on its unit card by turning the die.
// If no result is indicated, the attacker changes the result to a blank.
func applyAttackSurges(result *AttackResult, attack *Attack) {
	// first we apply critical X converting surges to crits
	applyCriticalX(result, attack)

	// now handle remaining surges
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
	} else {
		result.Red.N += result.Red.S
		result.Black.N += result.Black.S
		result.White.N += result.White.S
		result.Red.S = 0
		result.Black.S = 0
		result.White.S = 0
	}
}

// while converting offensive surges, change up to X Dice surge results to crit results
func applyCriticalX(result *AttackResult, attack *Attack) {
	for tot := attack.config.keywords.criticalX; tot > 0; {
		if result.Red.S > 0 {
			result.Red.C++
			result.Red.S--
			tot--
		} else if result.Black.S > 0 {
			result.Black.S--
			result.Black.C++
			tot--
		} else if result.White.S > 0 {
			result.White.S--
			result.White.C++
			tot--
		} else {
			break
		}
	}
}

// While attacking a unit that has Armor, change up to X hit results to crits
func applyImpactX(result *AttackResult, attack *Attack, defense *Defense) {
	if defense.config.keywords.armor || defense.config.keywords.armorX > 0 {
		for tot := attack.config.keywords.impactX; tot > 0; {
			if result.White.H > 0 {
				result.White.H--
				result.White.C++
				tot--
			} else if result.Black.H > 0 {
				result.Black.H--
				result.Black.C++
				tot--
			} else if result.Red.H > 0 {
				result.Red.H--
				result.Red.C++
				tot--
			} else {
				break
			}
		}
	}
}

// you may turn 1 attack die to a crit result
func applyRamX(result *AttackResult, attack *Attack) {
	for tot := attack.config.keywords.ramX; tot > 0; {
		if result.White.N > 0 {
			result.White.N--
			result.White.C++
			tot--
		} else if result.Black.N > 0 {
			result.Black.N--
			result.Black.C++
			tot--
		} else if result.Red.N > 0 {
			result.Red.N--
			result.Red.C++
			tot--
		} else if result.White.S > 0 {
			result.White.S--
			result.White.C++
			tot--
		} else if result.Black.S > 0 {
			result.Black.S--
			result.Black.C++
			tot--
		} else if result.Red.S > 0 {
			result.Red.S--
			result.Red.C++
			tot--
		} else if result.White.H > 0 {
			result.White.H--
			result.White.C++
			tot--
		} else if result.Black.H > 0 {
			result.Black.H--
			result.Black.C++
			tot--
		} else if result.Red.H > 0 {
			result.Red.H--
			result.Red.C++
			tot--
		} else {
			break
		}
	}
}

// While defending, cancel (X/all) hit results
func applyArmor(result *AttackResult, defense *Defense) {
	if defense.config.keywords.armor {
		result.Red.H = 0
		result.Black.H = 0
		result.White.H = 0
	} else if defense.config.keywords.armorX > 0 {
		for tot := defense.config.keywords.armorX; tot > 0; {
			if result.White.H > 0 {
				result.White.H--
				tot--
			} else if result.Black.H > 0 {
				result.Black.H--
				tot--
			} else if result.Red.H > 0 {
				result.Red.H--
				tot--
			} else {
				break
			}
		}
	}
}

// 7c Convert Defense Surges:
// The defender changes its defense surge results to the result indicated on its unit card by turning the die.
// If no result is indicated, the defender changes the result to a blank.
func applyDefenseSurges(result *DefenseResult, defense *Defense) {
	if defense.config.surgesToBlock {
		result.Red.B += result.Red.S
		result.White.B += result.White.S
		result.Red.S = 0
		result.White.S = 0
	} else {
		result.Red.N += result.Red.S
		result.White.N += result.White.S
		result.Red.S = 0
		result.White.S = 0
	}
}

func applyDodgeAndCover(result *AttackResult, attack *Attack, defense *Defense) {
	cover := 0

	// Blast: Ignore cover
	if !attack.config.keywords.blast {
		cover += defense.config.cover

		// Low Profile: While defending, if you have light cover, improve your cover by 1
		if cover == 1 && defense.config.keywords.lowProfile {
			cover = 2
		}

		// Cover X: While defending against a ranged attack, improve your cover by X
		cover += cover + defense.config.keywords.coverX

		// Sharpshooter X: While performing a ranged attack, reduce the defender’s cover by X
		if attack.config.keywords.sharpshooterX > 0 {
			cover -= max(attack.config.keywords.sharpshooterX, 0)
		}
	}

	hitsToRemove := min(cover, 2)

	if !attack.config.keywords.highVelocity {
		hitsToRemove += defense.config.tokens.dodge
	}

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

// when attacking, ignore up to X block results
func applyPierce(result *DefenseResult, attack *Attack) {
	for tot := attack.config.keywords.pierceX; tot > 0; {
		if result.White.B > 0 {
			result.White.B--
			tot--
		} else if result.Red.B > 0 {
			result.Red.B--
			tot--
		} else {
			break
		}
	}
}

// While defending, if the attack pool has Pierce X, roll x additional dice
func addImperviousToDefense(dice int, attack *Attack, defense *Defense) int {
	if defense.config.keywords.impervious && attack.config.keywords.pierceX > 0 {
		dice += attack.config.keywords.pierceX
	}
	return dice
}

// While defending against a ranged attack,during the 'Roll Defense Dice' step the defender may flip active shield tokens
// to add 1 block and roll 1 defense dice less for each token
func applyShieldToDefenseDice(dice int, defense *Defense) int {
	return max(dice - defense.config.tokens.shield, 0)
}

func applyShieldToDefenseResult(result *DefenseResult, defense *Defense) {
	if defense.config.rollsRedDefense {
		result.Red.B += defense.config.tokens.shield
	} else {
		result.White.B += defense.config.tokens.shield
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

func combineDefenseResults(a *DefenseResult, b *DefenseResult) {
	a.Red.B += b.Red.B
	a.Red.S += b.Red.S
	a.Red.N += b.Red.N

	a.White.B += b.White.B
	a.White.S += b.White.S
	a.White.N += b.White.N
}
