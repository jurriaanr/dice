package defend

import (
	"math/rand"
	"strings"
)

var redDice = [...]string{"b", "b", "b", "s", "n", "n"}
var whiteDice = [...]string{"b", "s", "n", "n", "n", "n"}

type Result struct {
	b int
	s int
	n int
}

func red() string {
	return redDice[rand.Intn(len(redDice))]
}

func white() string {
	return whiteDice[rand.Intn(len(whiteDice))]
}

func roll(redDice, whiteDice int) Result {
	dice := make([]string, redDice+whiteDice)

	for i := 0; i < redDice+whiteDice; i++ {
		switch {
		case i < redDice:
			dice[i] = red()
		case i < redDice+whiteDice:
			dice[i] = white()
		}
	}

	result := strings.Join(dice[:], "")

	return Result{
		b: strings.Count(result, "b"),
		s: strings.Count(result, "s"),
		n: strings.Count(result, "n"),
	}
}

func RoleResult(redDice, whiteDice int, hasSurge bool) int {
	result := roll(redDice, whiteDice)
	val := result.b
	if hasSurge {
		val += result.s
	}
	return val
}

func Test(redDice, whiteDice, rolls int, hasSurge bool) float64 {
	l := rolls
	sum := 0
	for rolls > 0 {
		sum += RoleResult(redDice, whiteDice, hasSurge)
		rolls -= 1
	}
	return float64(sum) / float64(l)
}
