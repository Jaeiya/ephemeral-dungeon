package main

import (
	"fmt"
	"math"

	"github.com/jaeiya/monster/internal/engine"
	_ "modernc.org/sqlite"
)

func main() {
	// w := engine.Dungeon{}
	// baseRoom := w.GenRoom(engine.West)

	// for i, r := range baseRoom.Rooms {
	// 	fmt.Printf("Room %d: %+v\n", i+1, r.Direction)
	// }

	// count := 0
	// for i := range 50 {
	// 	val := engine.CalcXPToLevel(i + 1)
	// 	fmt.Println(i+1, val)
	// 	count += val
	// }

	s := engine.NewSimulator()
	s.SimulateRooms(1000)
}

func calcXP(maxLevel, baseXP int, rate float64) {
	xpMap := make(map[int]int, maxLevel)
	rate = 1 + rate

	for lvl := range maxLevel {
		xpMap[lvl+1] = int(math.Round(float64(baseXP) * math.Pow(rate, float64(lvl))))
	}

	fmt.Println(xpMap)
}

func calcTotalXP(maxLevel, baseXP int, rate float64) {
	xpMap := make(map[int]int, maxLevel)
	xpMap[1] = baseXP

	rate = 1 + rate
	currentXP := baseXP
	totalXP := baseXP

	for lvl := range maxLevel - 1 {
		val := int(math.Round(float64(currentXP) * rate))
		totalXP += val
		xpMap[lvl+2] = totalXP
		currentXP = val
	}

	fmt.Println(totalXP)

	fmt.Println(xpMap)
}

func calcLevel(xp, baseXP int, rate float64) {
	rate = 1 + rate
	currentXP := baseXP
	lvl := 0

	for currentXP <= xp {
		currentXP = int(math.Round(float64(currentXP) * rate))
		lvl += 1
	}

	fmt.Println(xp, currentXP, lvl)
}

func calcLevelV2(xp, baseXP int, rate float64) {
	rate = 1 + rate
	currentXP := baseXP
	totalXP := currentXP
	lvl := 0

	for totalXP <= xp {
		val := int(math.Round(float64(currentXP) * rate))
		totalXP += val
		currentXP = val
		lvl += 1
	}

	fmt.Printf("XP: %d\nCurrentXP: %d\nLevel: %d\n", xp, currentXP, lvl)
}

func calcXPV3(lvl, baseXP int) int {
	const length = 10
	const ratio = 1.10

	deltaXP := float64(baseXP) * math.Pow(ratio, float64(lvl-1))
	return int(math.Round(deltaXP))
}
