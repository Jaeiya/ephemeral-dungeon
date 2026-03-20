package main

import (
	"fmt"
	"math"

	"github.com/jaeiya/monster/internal/engine"
)

func main() {
	w := engine.Dungeon{}

	// baseRoom := w.GenRoom(engine.West)

	// for i, r := range baseRoom.Rooms {
	// 	fmt.Printf("Room %d: %+v\n", i+1, r.Direction)
	// }

	// return

	worldSize := 500
	rooms := make([]engine.Directions, worldSize*4)
	halls := make([]engine.Directions, worldSize*4)
	totalRooms := 0

	for i := range worldSize {
		r, h, _ := w.GenerateExits(engine.Origin)
		totalRooms += len(r)
		copy(rooms[i*4:], r)
		copy(halls[i*4:], h)
	}

	// roomCounts := make([]int, worldSize)
	// for i, r := range rooms {
	// 	roomCounts[i] = len(r)
	// }

	// roomAvg := 0
	// for _, count := range roomCounts {
	// 	roomAvg += count
	// }

	oneCount := 0
	maxRoomOneCount := 0
	roomSpacingCounts := []int{}

	roomSizes := make([]int, 4)
	for i := 0; i < len(rooms); i += 4 {
		r := rooms[i : i+4]

		sizeIdx := 0
		roomLen := 0
		for k := len(r) - 1; k > -1; k-- {
			if r[k] > 0 {
				sizeIdx = k
				roomLen = k + 1
				break
			}
		}

		roomSizes[sizeIdx] += 1
		if roomLen == 1 {
			oneCount += 1
		} else {
			if oneCount > maxRoomOneCount {
				maxRoomOneCount = oneCount
			}
			if oneCount > 0 {
				roomSpacingCounts = append(roomSpacingCounts, oneCount)
			}
			oneCount = 0
		}
	}

	hallSizes := make([]int, 5)
	hallCounts := make([]int, worldSize)
	spacingCounts := []int{}
	zeroCount := 0
	maxHallZeroCount := 0

	track := 0
	for i := 0; i < len(halls); i += 4 {
		h := halls[i : i+4]

		sizeIdx := 0
		roomLen := 0
		for k := len(h) - 1; k > -1; k-- {
			if h[k] > 0 {
				sizeIdx = k + 1
				roomLen = k + 1
				break
			}
		}

		hallCounts[track] = roomLen
		track++
		hallSizes[sizeIdx] += 1
		if roomLen == 0 {
			zeroCount += 1
		} else {
			if zeroCount > maxHallZeroCount {
				maxHallZeroCount = zeroCount
			}
			if zeroCount > 0 {
				spacingCounts = append(spacingCounts, zeroCount)
			}
			zeroCount = 0
		}
	}

	hallSum := 0
	for _, count := range hallCounts {
		hallSum += count
	}

	spacingSum := 0
	for _, count := range spacingCounts {
		spacingSum += count
	}

	roomSpacingSum := 0
	for _, count := range roomSpacingCounts {
		roomSpacingSum += count
	}

	calcPercent := func(size, maxSize int) float64 {
		return float64(size) / float64(maxSize) * 100.0
	}

	fmt.Printf("Avg Hall Spacing: %d\n", spacingSum/len(spacingCounts))
	fmt.Printf("Avg Room Spacing: %d\n", roomSpacingSum/len(roomSpacingCounts))
	fmt.Printf("Max Room Spacing: %d\n", maxRoomOneCount)
	fmt.Printf("Max Hall Spacing: %d\n\n", maxHallZeroCount)
	for i, size := range roomSizes {
		fmt.Printf("Room%d: %d (%0.f%%)\n", i+1, size, calcPercent(size, worldSize))
	}

	fmt.Println("")

	for i, size := range hallSizes {
		fmt.Printf("Hall%d: %d (%.2f%%)\n", i, size, calcPercent(size, worldSize))
	}

	fmt.Println("")

	hallChance := 0.0
	for i := range 4 {
		hallChance += calcPercent(hallSizes[i+1], worldSize)
	}
	fmt.Printf("Hall Chance: %.2f%%\n", hallChance)

	engine.EncodeQuestion("how do I lift the garage door")
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
