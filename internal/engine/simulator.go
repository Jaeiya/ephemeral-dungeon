package engine

import "fmt"

type Simulator struct{}

type SimulatedWorld struct {
	RoomExits      []Direction
	HallExits      []Direction
	RoomExitsCount int
	HallExitsCount int
}

type SimulatedRoomStats struct {
	ExitCounts       [4]int // How many rooms with 1, 2, 3, or 4 exits
	SingleExitRuns   []int  // How many single-exit rooms between each multi-exit room
	MaxSingleExitRun int    // The max consecutive single-exit rooms before a multi-exit room
}

type SimulatedHallStats struct {
	HallCounts     [5]int // How many rooms with 0, 1, 2, 3, or 4 halls
	ZeroHallRuns   []int  // How many zero-hall rooms between each multi-hall room
	MaxZeroHallRun int    // The max consecutive zero-hall rooms before a multi-hall room
}

func NewSimulator() Simulator {
	return Simulator{}
}

func (s Simulator) SimulateRooms(roomCount int) {
	world := s.GenerateWorld(roomCount)
	h := s.GenerateHallStats(world.HallExits, roomCount)
	r := s.GenerateRoomStats(world.RoomExits)

	spacingSum := 0
	for _, count := range h.ZeroHallRuns {
		spacingSum += count
	}

	roomSpacingSum := 0
	for _, count := range r.SingleExitRuns {
		roomSpacingSum += count
	}

	calcPercent := func(size, maxSize int) float64 {
		return float64(size) / float64(maxSize) * 100.0
	}

	fmt.Println()
	fmt.Printf("Total Rooms: %d\n", roomCount)
	fmt.Printf("Total Exits: %d\n", world.RoomExitsCount)
	fmt.Printf("Total Halls: %d\n\n", world.HallExitsCount)
	fmt.Printf(
		"Avg Room Spacing: %3d (avg # of 1-exit rooms until a multi-exit room)\n",
		roomSpacingSum/len(r.SingleExitRuns),
	)
	fmt.Printf(
		"Max Room Spacing: %3d (max # of 1-exit rooms until a multi-exit room)\n",
		r.MaxSingleExitRun,
	)
	fmt.Printf(
		"Avg Hall Spacing: %3d (avg # of rooms between 2 halls)\n",
		spacingSum/len(h.ZeroHallRuns),
	)
	fmt.Printf("Max Hall Spacing: %3d (max # of rooms between 2 halls)\n\n", h.MaxZeroHallRun)
	for i, size := range r.ExitCounts {
		fmt.Printf("Rooms %d-exit: %d (%0.f%%)\n", i+1, size, calcPercent(size, roomCount))
	}

	fmt.Println()

	for i, size := range h.HallCounts {
		fmt.Printf("Rooms %d-halls: %d (%.2f%%)\n", i, size, calcPercent(size, roomCount))
	}

	fmt.Println()

	hallChance := 0.0
	for i := range 4 {
		hallChance += calcPercent(h.HallCounts[i+1], roomCount)
	}
	fmt.Printf("Hall Chance: %.2f%%\n", hallChance)
}

func (s Simulator) GenerateWorld(size int) SimulatedWorld {
	roomExits := make([]Direction, size*4)
	hallExits := make([]Direction, size*4)
	var totalRoomExits int
	var totalHallExits int

	for i := range size {
		r, h, _ := Dungeon{}.GenerateExits(Origin)
		totalRoomExits += len(r)
		totalHallExits += len(h)
		copy(roomExits[i*4:], r)
		copy(hallExits[i*4:], h)
	}
	return SimulatedWorld{
		RoomExits:      roomExits,
		HallExits:      hallExits,
		RoomExitsCount: totalRoomExits,
		HallExitsCount: totalHallExits,
	}
}

func (s Simulator) GenerateRoomStats(exits []Direction) SimulatedRoomStats {
	singleExitCount := 0
	maxSingleExitRun := 0

	var singleExitRuns []int
	exitCounts := [4]int{}

	for i := 0; i < len(exits); i += 4 {
		room := exits[i : i+4]

		roomLen := 0
		for _, v := range room {
			if v == 0 {
				break
			}
			roomLen++
		}

		exitCounts[roomLen-1] += 1
		if roomLen == 1 {
			singleExitCount += 1
			continue
		}

		if singleExitCount > maxSingleExitRun {
			maxSingleExitRun = singleExitCount
		}
		if singleExitCount > 0 {
			singleExitRuns = append(singleExitRuns, singleExitCount)
		}
		singleExitCount = 0
	}

	return SimulatedRoomStats{
		ExitCounts:       exitCounts,
		SingleExitRuns:   singleExitRuns,
		MaxSingleExitRun: maxSingleExitRun,
	}
}

func (s Simulator) GenerateHallStats(halls []Direction, worldSize int) SimulatedHallStats {
	hallCounts := [5]int{}

	zeroHallRuns := []int{}
	zeroHallCount := 0
	maxZeroHallRun := 0

	hallIdx := 0
	for i := 0; i < len(halls); i += 4 {
		hall := halls[i : i+4]

		hallLen := 0
		for k := range hall {
			if hall[k] == 0 {
				break
			}
			hallLen = k + 1
		}

		hallIdx++
		hallCounts[hallLen] += 1

		if hallLen == 0 {
			zeroHallCount += 1
			continue
		}

		if zeroHallCount > maxZeroHallRun {
			maxZeroHallRun = zeroHallCount
		}
		if zeroHallCount > 0 {
			zeroHallRuns = append(zeroHallRuns, zeroHallCount)
		}
		zeroHallCount = 0
	}

	return SimulatedHallStats{
		HallCounts:     hallCounts,
		ZeroHallRuns:   zeroHallRuns,
		MaxZeroHallRun: maxZeroHallRun,
	}
}
