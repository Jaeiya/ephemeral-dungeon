package engine

import (
	"math/rand/v2"
)

// RollDodge checks if a character can evade an attack
// based on the defenders dodge and attackers agility.
func RollDodge(dodge, agility int) bool {
	return relativeRatio(dodge, agility, 0.3)
}

// RollFirstAttack checks to see who strikes first in
// combat based on the characters initiative.
func RollFirstAttack(init1, init2 int) bool {
	return relativeRatio(init1, init2, 0.5)
}

// RollPercent returns true if pRNG returns a value within
// the % chance specified.
//
// 🔵 Chance is the percentage that a roll succeeds. If
// chance is 10 then there is a 10% chance that a roll
// will succeed.
//
// 🔴 Panics when chance < 0 or chance > 100
func RollPercent(chance float64) bool {
	if chance < 0 {
		panic("cannot calculate a negative percent")
	}
	if chance > 100 {
		panic("cannot calculate a chance greater than 100%")
	}
	return rand.Float64() < (float64(chance) / 100.0)
}

func relativeRatio(stat1, stat2 int, weight float64) bool {
	if stat2 == 0 && stat1 > 0 {
		return true
	}

	coefficient := float64(stat1) / float64(stat2)
	// If both stats are equal, then the weight becomes the chance
	chance := coefficient * weight

	if chance >= 1.0 {
		return true
	}
	return rand.Float64() < chance
}
