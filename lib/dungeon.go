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

type Directions uint8

const (
	None  Directions = 0
	North Directions = 1 << (iota - 1)
	South
	East
	West
	NorthEast
	NorthWest
	SouthEast
	SouthWest
)

func (d Directions) String() string {
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
	location Directions
	desc     string
	objects  []string
}

type Room interface {
	findObj(obj string)
	interact(obj string)
}

type Exits struct {
	rooms []Directions
	halls []Directions
	mask  Directions
}

var (
	roomDirections = [4]Directions{North, South, East, West}
	hallDirections = [4]Directions{NorthEast, NorthWest, SouthEast, SouthWest}
)

// GenerateExits randomizes the possible exit directions and returns them
// as rooms and hallways.
//
// 🔵 Hallways have a significantly lower chance of appearing
func (d Dungeon) GenerateExits(
	lastDir Directions,
) (rooms, halls []Directions, exits Directions) {
	rooms = d.genRoomDirs(lastDir)
	halls = d.genHallDirs(lastDir)

	for _, r := range rooms {
		exits |= r
	}

	for _, h := range halls {
		exits |= h
	}

	return rooms, halls, exits
}

func (d Dungeon) genRoomDirs(lastDir Directions) []Directions {
	rooms := make([]Directions, 0, 4)
	for _, d := range roomDirections {
		if d == lastDir {
			continue
		}
		rooms = append(rooms, d)
	}

	chances := [3]float64{45, 30, 20}

	roomLimit := 1

	for i := 1; i < len(rooms); i++ {
		if RollPercent(chances[i-1]) {
			roomLimit++
		} else {
			break
		}
	}

	utils.Shuffle(rooms)

	return rooms[:roomLimit]
}

func (d Dungeon) genHallDirs(lastDir Directions) []Directions {
	chances := [4]float64{23, 27, 20, 10}

	if !RollPercent(chances[0]) {
		return []Directions{}
	}

	hallLimit := 1

	halls := make([]Directions, 0, 4)
	for _, d := range hallDirections {
		if d == lastDir {
			continue
		}
		halls = append(halls, d)
	}

	for i := 1; i < len(halls); i++ {
		if RollPercent(chances[i]) {
			hallLimit++
		} else {
			break
		}
	}

	utils.Shuffle(halls)

	return halls[:hallLimit]
}
