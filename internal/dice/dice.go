package dice

import (
	"math/rand/v2"
	"sort"
)

// Log describes one roll for the UI animation.
type Log struct {
	Label   string `json:"label"`
	Sides   int    `json:"sides"`
	Dice    []int  `json:"dice"`
	Dropped []int  `json:"dropped,omitempty"`
	Bonus   int    `json:"bonus,omitempty"`
	Total   int    `json:"total"`
}

func D(sides int) int { return rand.IntN(sides) + 1 }

func Roll(label string, n, sides int) Log {
	l := Log{Label: label, Sides: sides}
	for range n {
		v := D(sides)
		l.Dice = append(l.Dice, v)
		l.Total += v
	}
	return l
}

// DropLowest rolls n dice and drops the k lowest.
func DropLowest(label string, n, sides, k int) Log {
	l := Roll(label, n, sides)
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(a, b int) bool { return l.Dice[idx[a]] < l.Dice[idx[b]] })
	l.Total = 0
	drop := map[int]bool{}
	for _, i := range idx[:k] {
		drop[i] = true
	}
	for i, v := range l.Dice {
		if drop[i] {
			l.Dropped = append(l.Dropped, i)
		} else {
			l.Total += v
		}
	}
	return l
}

func Pick[T any](xs []T) T { return xs[rand.IntN(len(xs))] }

func PickN[T any](xs []T, n int) []T {
	p := rand.Perm(len(xs))
	if n > len(xs) {
		n = len(xs)
	}
	out := make([]T, 0, n)
	for _, i := range p[:n] {
		out = append(out, xs[i])
	}
	return out
}

// RollPlus rolls n dice and adds a flat bonus.
func RollPlus(label string, n, sides, bonus int) Log {
	l := Roll(label, n, sides)
	l.Bonus = bonus
	l.Total += bonus
	return l
}

func Chance(p float64) bool { return rand.Float64() < p }

func Mod5e(score int) int {
	if score < 10 {
		return (score - 11) / 2
	}
	return (score - 10) / 2
}
