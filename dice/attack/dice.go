package attack

import (
	"math/rand"
	"strings"
)

var redDice = [...]string{"h", "h", "h", "h", "h", "c", "s", "b"}
var blackDice = [...]string{"h", "h", "h", "c", "s", "n", "n", "n"}
var whiteDice = [...]string{"h", "c", "s", "n", "n", "n", "n", "n"}

type Result struct {
	h int
	c int
	s int
	n int
}

func red() string {
	return redDice[rand.Intn(len(redDice))]
}

func black() string {
	return blackDice[rand.Intn(len(blackDice))]
}

func white() string {
	return whiteDice[rand.Intn(len(whiteDice))]
}

func roll(redDice, blackDice, whiteDice int) Result {
	dice := make([]string, redDice+blackDice+whiteDice)

	for i := 0; i < redDice+blackDice+whiteDice; i++ {
		switch {
		case i < redDice:
			dice[i] = red()
		case i < redDice+blackDice:
			dice[i] = black()
		case i < redDice+blackDice+whiteDice:
			dice[i] = white()
		}
	}

	result := strings.Join(dice[:], "")

	return Result{
		h: strings.Count(result, "h"),
		c: strings.Count(result, "c"),
		s: strings.Count(result, "s"),
		n: strings.Count(result, "n"),
	}
}

func RoleResult(redDice, blackDice, whiteDice int, hasSurge bool) int {
	result := roll(redDice, blackDice, whiteDice)
	val := result.h + result.c
	if hasSurge {
		val += result.s
	}
	return val
}

func Test(redDice, blackDice, whiteDice, rolls int, hasSurge bool) float64 {
	l := rolls
	sum := 0
	for rolls > 0 {
		sum += RoleResult(redDice, blackDice, whiteDice, hasSurge)
		rolls -= 1
	}
	return float64(sum) / float64(l)
}
