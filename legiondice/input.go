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

func AddCriticalXToAttack(criticalX int, attack *Attack) {
	attack.config.keywords.criticalX = criticalX
}

func AddRamXToAttack(ramX int, attack *Attack) {
	attack.config.keywords.ramX = ramX
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

func AddCoverXToDefense(coverX int, defense *Defense) {
	defense.config.keywords.coverX = coverX
}

func AttackFromRequest(request *http.Request) Attack {
	// attack dice
	r := paramToInt("r", request)
	b := paramToInt("b", request)
	w := paramToInt("w", request)
	// attack surges conversion (crits, hits, none)
	as := request.URL.Query().Get("as")

	attack := CreateAttack(
		int(r),
		int(b),
		int(w),
		as,
	)

	aim := paramToInt("aim", request)
	preciseX := paramToInt("preciseX", request)
	pierceX := paramToInt("pierceX", request)
	criticalX := paramToInt("criticalX", request)
	ramX := paramToInt("ramX", request)

	AddAimToAttack(aim, &attack)
	AddPreciseXToAttack(preciseX, &attack)
	AddPierceXToAttack(pierceX, &attack)
	AddCriticalXToAttack(criticalX, &attack)
	AddRamXToAttack(ramX, &attack)

	return attack
}

func DefenseFromRequest(request *http.Request) Defense {
	// defense dice type (red/white)
	d := request.URL.Query().Get("d")

	// defense surges conversion (true/false)
	ds := paramToBoolean("ds", request)
	cover := paramToInt("cover", request)

	defense := CreateDefense(d, ds, cover)

	dodge := paramToInt("dodge", request)
	coverX := paramToInt("coverX", request)

	AddDodgeToDefense(dodge, &defense)
	AddCoverXToDefense(coverX, &defense)

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

func paramToInt(key string, request *http.Request) int {
	return stringToInt(request.URL.Query().Get(key))
}
