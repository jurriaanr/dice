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

	switch strings.ToLower(surge) {
	case "hits":
		attack.config.surgesToHits = true
	case "crits":
		attack.config.surgesToCrits = true
	}

	return attack
}

func CreateDefense(defenseDice string, surges bool) Defense {
	// setup Defense pool
	defense := Defense{}

	if strings.ToLower(defenseDice) == "red" {
		defense.config.rollsRedDefense = true
	} else {
		defense.config.rollsWhiteDefense = true
	}

	defense.config.surgesToBlock = surges

	return defense
}


func AttackFromRequest(request *http.Request) Attack {
	// attack dice
	r, _ := strconv.ParseInt(request.URL.Query().Get("r"), 10, 64)
	b, _ := strconv.ParseInt(request.URL.Query().Get("b"), 10, 64)
	w, _ := strconv.ParseInt(request.URL.Query().Get("w"), 10, 64)
	// attack surges conversion (crits, hits, none)
	as := request.URL.Query().Get("as")

	return CreateAttack(
		int(r),
		int(b),
		int(w),
		as,
	)
}

func DefenseFromRequest(request *http.Request) Defense {
	// defense dice type (red/white)
	d := request.URL.Query().Get("d")

	// defense surges conversion (true/false)
	ds := stringToBoolean(request.URL.Query().Get("ds"))

	return CreateDefense(d, ds)
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
