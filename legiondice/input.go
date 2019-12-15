package dice

import (
	"net/http"
	"strconv"
	"strings"
)

func CreateAttack(
	redAttackDice,
	blackAttackDice,
	whiteAttackDice int,
	surge string,
) Attack {
	// setup Attack pool
	attack := Attack{
		red:   redAttackDice,
		black: blackAttackDice,
		white: whiteAttackDice,
	}

	// the way surges are converted
	switch strings.ToLower(surge) {
	case "hits":
		attack.config.surgesToHits = true
	case "crits":
		attack.config.surgesToCrits = true
	}

	return attack
}

func AddAimToAttack(aim int, attack *Attack) {
	attack.config.tokens.aim = aim
}

func AddPreciseXToAttack(preciseX int, attack *Attack) {
	attack.config.keywords.preciseX = preciseX
}

func AddPierceXToAttack(pierceX int, attack *Attack) {
	attack.config.keywords.pierceX = pierceX
}

func AddImpactXToAttack(impactX int, attack *Attack) {
	attack.config.keywords.impactX = impactX
}

func AddCriticalXToAttack(criticalX int, attack *Attack) {
	attack.config.keywords.criticalX = criticalX
}

func AddRamXToAttack(ramX int, attack *Attack) {
	attack.config.keywords.ramX = ramX
}

func AddSharpshooterXToAttack(sharpshooterX int, attack *Attack) {
	attack.config.keywords.sharpshooterX = sharpshooterX
}

func AddBlastToAttack(blast bool, attack *Attack) {
	attack.config.keywords.blast = blast
}

func AddHighVelocityToAttack(highVelocity bool, attack *Attack) {
	attack.config.keywords.highVelocity = highVelocity
}

func CreateDefense(defenseDice string, surges bool, cover int) Defense {
	// setup Defense pool
	defense := Defense{}

	if strings.ToLower(defenseDice) == "red" {
		defense.config.rollsRedDefense = true
	} else {
		defense.config.rollsWhiteDefense = true
	}

	defense.config.surgesToBlock = surges
	defense.config.cover = cover

	return defense
}

func AddDodgeToDefense(dodge int, defense *Defense) {
	defense.config.tokens.dodge = dodge
}

func AddShieldToDefense(shield int, defense *Defense) {
	defense.config.tokens.shield = shield
}

func AddCoverXToDefense(coverX int, defense *Defense) {
	defense.config.keywords.coverX = coverX
}

func AddArmorToDefense(armor bool, defense *Defense) {
	defense.config.keywords.armor = armor
}

func AddArmorXToDefense(armorX int, defense *Defense) {
	defense.config.keywords.armorX = armorX
}

func AddUncannyLuckXToDefense(uncannyLuckX int, defense *Defense) {
	defense.config.keywords.uncannyLuckX = uncannyLuckX
}

func AddLowProfileToDefense(lowProfile bool, defense *Defense) {
	defense.config.keywords.lowProfile = lowProfile
}

func AddImperviousToDefense(impervious bool, defense *Defense) {
	defense.config.keywords.impervious = impervious
}

func AttackFromRequest(request *http.Request) Attack {
	// attack dice
	r := paramToInt("r", request, 25)
	b := paramToInt("b", request, 25)
	w := paramToInt("w", request, 25)
	// attack surges conversion (crits, hits, none)
	as := request.URL.Query().Get("as")

	attack := CreateAttack(
		int(r),
		int(b),
		int(w),
		as,
	)

	aim := paramToInt("aim", request, 10)
	preciseX := paramToInt("preciseX", request, 10)
	pierceX := paramToInt("pierceX", request, 10)
	impactX := paramToInt("impactX", request, 10)
	criticalX := paramToInt("criticalX", request, 10)
	ramX := paramToInt("ramX", request, 10)
	sharpshooterX := paramToInt("sharpshooterX", request, 10)
	blast := paramToBoolean("blast", request)
	highVelocity := paramToBoolean("highVelocity", request)

	AddAimToAttack(aim, &attack)
	AddPreciseXToAttack(preciseX, &attack)
	AddPierceXToAttack(pierceX, &attack)
	AddImpactXToAttack(impactX, &attack)
	AddCriticalXToAttack(criticalX, &attack)
	AddRamXToAttack(ramX, &attack)
	AddSharpshooterXToAttack(sharpshooterX, &attack)
	AddBlastToAttack(blast, &attack)
	AddHighVelocityToAttack(highVelocity, &attack)

	return attack
}

func DefenseFromRequest(request *http.Request) Defense {
	// defense dice type (red/white)
	diceColor := request.URL.Query().Get("d")

	// defense surges conversion (true/false)
	surge := paramToBoolean("ds", request)
	cover := paramToInt("cover", request, 10)

	defense := CreateDefense(diceColor, surge, cover)

	armor := paramToBoolean("armor", request)
	dodge := paramToInt("dodge", request, 10)
	shield := paramToInt("shield", request, 10)
	coverX := paramToInt("coverX", request, 10)
	armorX := paramToInt("armorX", request, 10)
	uncannyLuckX := paramToInt("uncannyLuckX", request, 10)
	lowProfile := paramToBoolean("lowProfile", request)
	impervious := paramToBoolean("impervious", request)

	AddDodgeToDefense(dodge, &defense)
	AddShieldToDefense(shield, &defense)
	AddCoverXToDefense(coverX, &defense)
	AddArmorToDefense(armor, &defense)
	AddArmorXToDefense(armorX, &defense)
	AddUncannyLuckXToDefense(uncannyLuckX, &defense)
	AddLowProfileToDefense(lowProfile, &defense)
	AddImperviousToDefense(impervious, &defense)

	return defense
}

func stringToBoolean(text string) bool {
	switch text {
	case "true":
		return true
	case "t":
		return true
	case "1":
		return true
	case "y":
		return true
	case "yes":
		return true
	}

	return false
}

func stringToInt(text string) int {
	val, _ := strconv.ParseInt(text, 10, 64)
	return int(val)
}

func paramToBoolean(key string, request *http.Request) bool {
	return stringToBoolean(request.URL.Query().Get(key))
}

func paramToInt(key string, request *http.Request, maxVal int) int {
	return min(stringToInt(request.URL.Query().Get(key)), maxVal)
}
