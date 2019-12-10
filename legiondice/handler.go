package dice

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type Roll struct {
	Attack  AttackResult
	Defense DefenseResult
}

type Result struct {
	Successes float64
	Rolls     []Roll
}

func Test(attack *Attack, defense *Defense, rolls int, logs int) Result {
	// seed random generator
	rand.Seed(time.Now().Unix())

	// hold the first X results
	results := make([]Roll, logs)

	// collect total number of hits
	sum := 0
	for i := 0; i < rolls; i++ {
		// 4 Roll Attack Dice (includes step 4a, 4b, 4c, 5 and 6)
		// The attacker rolls the dice in the attack pool
		hits, attackResult := AttackRollResult(attack, defense)

		// 7 Roll Defense Dice (includes 7a, 7b, 7c and 8)
		blocks, defenseResult := DefenseRoleResult(hits, attack, defense)

		// 9 Compare Results:
		remainingHits := hits - blocks

		// increase total with number of Successes
		sum += remainingHits

		if i < logs {
			results[i] = Roll{
				attackResult,
				defenseResult,
			}
		}
	}

	return Result{
		float64(sum) / float64(rolls),
		results,
	}
}

func RollDice(response http.ResponseWriter, request *http.Request) {
	allowedOrigins := [8]string{
		"http://legion.localhost",
		"http://legion.localhost.charlesproxy.com",
		"http://legion.localhost:81",
		"http://legion.localhost.charlesproxy.com:81",
		"http://www.swlegion.space",
		"https://www.swlegion.space",
		"http://swlegion.space",
		"https://swlegion.space",
	}

	origin := request.Header.Get("origin")

	response.Header().Set("Content-Type", "application/json")

	for i := 0; i < len(allowedOrigins); i++ {
		if allowedOrigins[i] == origin {
			response.Header().Set("Access-Control-Allow-Origin", origin)

			// Set CORS headers for the preflight request
			if request.Method == http.MethodOptions {
				response.Header().Set("Access-Control-Allow-Methods", "GET,POST")
				response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				response.Header().Set("Access-Control-Max-Age", "3600")
				response.WriteHeader(http.StatusNoContent)
				return
			}
			break
		}
	}

	attack := AttackFromRequest(request)
	defense := DefenseFromRequest(request)
	result := Test(&attack, &defense, 10000, 25)

	json.NewEncoder(response).Encode(result)
}
