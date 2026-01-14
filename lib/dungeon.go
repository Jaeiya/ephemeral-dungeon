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

type Direction int

const (
	None Direction = iota
	North
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
		return "<missing_direction>"
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
	roomDirections = []Direction{North, South, East, West}
	hallDirections = []Direction{NorthEast, NorthWest, SouthEast, SouthWest}
)

// GenerateExits randomizes the possible exit directions and returns them
// as rooms and hallways.
//
// 🔵 Hallways have a significantly lower chance of appearing
func (d Dungeon) GenerateExits(lastDir Direction) (rooms, halls []Direction) {
	return d.genRoomDirs(lastDir), d.genHallDirs(lastDir)
}

func (d Dungeon) genRoomDirs(lastDir Direction) []Direction {
	rooms := utils.FilterSlice(
		roomDirections,
		func(d Direction) bool { return lastDir == d },
	)

	utils.Shuffle(rooms)

	roomLimit := 1
	if RollPercent(45) {
		roomLimit += 1
	}
	if RollPercent(30) && roomLimit > 1 {
		roomLimit += 1
	}
	if RollPercent(20) && roomLimit > 2 {
		roomLimit += 1
	}

	return rooms[:roomLimit]
}

func (d Dungeon) genHallDirs(lastDir Direction) []Direction {
	if RollPercent(23) {
		halls := utils.FilterSlice(
			hallDirections,
			func(d Direction) bool { return lastDir == d },
		)
		utils.Shuffle(halls)
		hallLimit := 1
		if RollPercent(27) {
			hallLimit += 1
		}
		if RollPercent(20) && hallLimit > 1 {
			hallLimit += 1
		}
		if RollPercent(10) && hallLimit > 2 {
			hallLimit += 1
		}
		return halls[:hallLimit]
	}

	return []Direction{}
}
