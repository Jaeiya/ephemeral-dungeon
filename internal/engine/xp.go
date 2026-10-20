package engine

import "math"

// CalcXP calculates the XP required for a level using a power curve.
//
// 🔵 a rate of 0 is linear and a rate of 1 is quadratic
func CalcXP(level, baseXP int, rate float64) int {
	return baseXP + int(math.Pow(float64(level), rate+1))
}

// CalcXPToLevel calculates the XP required to go from
// currentLevel to goalLevel.
func CalcXPToLevel(currentLevel, goalLevel, baseXP int, rate float64) int {
	if currentLevel >= goalLevel {
		panic("cannot calculate a starting level > level end")
	}
	xpSum := 0
	for l := range goalLevel - currentLevel {
		l += currentLevel
		xpSum += CalcXP(l+1, baseXP, rate)
	}
	return xpSum
}
