package engine

import (
	"math/rand/v2"
)

// RollDodge checks if a character can evade an attack
// based on the defenders dodge and attackers agility.
func RollDodge(dodge, agility int) bool {
	return RollRelativeChance(dodge, agility)
}

// RollFirstAttack checks to see who strikes first in
// combat based on the characters initiative.
func RollFirstAttack(init1, init2 int) bool {
	return RollRelativeChance(init1, init2)
}

// RollPercent returns true if pRNG returns a value within
// the % chance specified.
//
// 🔵 Chance is the percentage that a roll succeeds. If
// chance is 10 then there is a 10% chance that a roll
// will succeed.
//
// 🔴 Panics when chance < 0 or > 100
func RollPercent(chance float64) bool {
	if chance < 0 {
		panic("cannot calculate a negative percent")
	}
	if chance > 100 {
		panic("cannot calculate a chance greater than 100%")
	}
	return rand.Float64() < (float64(chance) / 100.0)
}

// RollRelativeChance checks whether the primary stat has a chance
// to succeed over the secondary stat.
//
// 🔵 The larger the primary stat, the higher the chance to
// return true. The larger the secondary stat, the higher the
// chance to return false.
//
// 🔵 If both stats are equal, the chance to succeed is 50%
func RollRelativeChance(primary, secondary int) bool {
	if primary < 0 || secondary < 0 {
		panic("stats should never be less than zero")
	}

	var chance float64
	f1, f2 := float64(primary), float64(secondary)
	pivot := 0.5

	switch {
	case (f1 == 0 && f2 == 0) || f1 == f2:
		chance = pivot

	case f1 == 0:
		return false

	case f2 == 0:
		return true

	case f2 > f1:
		chance = f1 / f2 * pivot

	default:
		chance = 1.0 - (f2 / f1 * pivot)
	}

	if chance > 0.99 {
		return true
	}

	if chance < 0.01 {
		return false
	}

	return rand.Float64() < chance
}
