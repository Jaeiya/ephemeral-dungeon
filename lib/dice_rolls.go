package lib

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
