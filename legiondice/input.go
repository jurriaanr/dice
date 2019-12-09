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
	r, _ := strconv.ParseInt(request.URL.Query().Get("r"), 10, 64)
	b, _ := strconv.ParseInt(request.URL.Query().Get("b"), 10, 64)
	w, _ := strconv.ParseInt(request.URL.Query().Get("w"), 10, 64)
	// attack surges conversion (crits, hits, none)
	as := request.URL.Query().Get("as")

	aim, _ := strconv.ParseInt(request.URL.Query().Get("aim"), 10, 64)
	preciseX, _ := strconv.ParseInt(request.URL.Query().Get("preciseX"), 10, 64)

	attack := CreateAttack(
		int(r),
		int(b),
		int(w),
		as,
	)

	if aim > 0 {
		AddAimToAttack(int(aim), &attack)
	}

	if preciseX > 0 {
		AddPreciseXToAttack(int(preciseX), &attack)
	}

	return attack
}

func DefenseFromRequest(request *http.Request) Defense {
	// defense dice type (red/white)
	d := request.URL.Query().Get("d")

	// defense surges conversion (true/false)
	ds := stringToBoolean(request.URL.Query().Get("ds"))
	cover, _ := strconv.ParseInt(request.URL.Query().Get("cover"), 10, 64)

	defense := CreateDefense(d, ds, int(cover))

	dodge, _ := strconv.ParseInt(request.URL.Query().Get("dodge"), 10, 64)
	coverX, _ := strconv.ParseInt(request.URL.Query().Get("coverX"), 10, 64)

	if dodge > 0 {
		AddDodgeToDefense(int(dodge), &defense)
	}

	if coverX > 0 {
		AddCoverXToDefense(int(coverX), &defense)
	}

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
