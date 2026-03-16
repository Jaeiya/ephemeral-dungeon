package lib

import (
	"github.com/jaeiya/monster/lib/utils"
)

const worldInfo = `


DUNGEON LAYOUT
	- Rooms should be dynamically generated

	- Rooms are not navigable backwards

	- Diagonal directions should always be hallways with high-risk, high-reward events
		- Groups of enemies
		- Low-percent high-reward chest spawn

	- Use efficient memory for room generation
		- Copy adjacent room data into current room data when player moves
		- Generate adjacent rooms into existing adjacent room array

	- Rooms should be stylized based on the level of dungeon
		- The deeper you go, the more interesting it should become

	- Every 10 rooms a boss will appear unless a hallway is taken
		- The next room after the hallway will contain the boss


`

type Direction uint8

const (
	None  Direction = 0
	North Direction = 1 << (iota - 1)
	South
	East
	West
	NorthEast
	NorthWest
	SouthEast
	SouthWest
)

func (d Direction) String() string {
	switch d {
	case North:
		return "North"
	case South:
		return "South"
	case East:
		return "East"
	case West:
		return "West"
	case NorthEast:
		return "NorthEast"
	case NorthWest:
		return "NorthWest"
	case SouthEast:
		return "SouthEast"
	case SouthWest:
		return "SouthWest"
	default:
		return "<N/A>"
	}
}

type BaseRoom struct {
	location Direction
	desc     string
	objects  []string
}

type Room interface {
	findObj(obj string)
	interact(obj string)
}

type Dungeon struct{}

var (
	roomDirections = [4]Direction{North, South, East, West}
	hallDirections = [4]Direction{NorthEast, NorthWest, SouthEast, SouthWest}
)

// GenerateExits randomizes the possible exit directions and returns them
// as rooms and hallways.
//
// 🔵 Hallways have a significantly lower chance of appearing
func (d Dungeon) GenerateExits(lastDir Direction) (rooms, halls []Direction) {
	return d.genRoomDirs(lastDir), d.genHallDirs(lastDir)
}

func (d Dungeon) genRoomDirs(lastDir Direction) []Direction {
	rooms := make([]Direction, 0, 4)
	for _, d := range roomDirections {
		if d == lastDir {
			continue
		}
		rooms = append(rooms, d)
	}

	utils.Shuffle(rooms)

	chances := [3]float64{45, 30, 20}

	roomLimit := 1

	for i := 1; i < len(rooms); i++ {
		if RollPercent(chances[i-1]) {
			roomLimit++
		} else {
			break
		}
	}

	return rooms[:roomLimit]
}

func (d Dungeon) genHallDirs(lastDir Direction) []Direction {
	chances := [4]float64{23, 27, 20, 10}

	if !RollPercent(chances[0]) {
		return []Direction{}
	}

	hallLimit := 1

	halls := make([]Direction, 0, 4)
	for _, d := range hallDirections {
		if d == lastDir {
			continue
		}
		halls = append(halls, d)
	}

	utils.Shuffle(halls)

	for i := 1; i < len(halls); i++ {
		if RollPercent(chances[i]) {
			hallLimit++
		} else {
			break
		}
	}

	return halls[:hallLimit]
}
