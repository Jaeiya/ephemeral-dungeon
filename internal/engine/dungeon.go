package engine

import (
	"fmt"

	"github.com/jaeiya/monster/internal/utils"
)

const worldInfo = `


DUNGEON LAYOUT
	- Rooms should be dynamically generated
		- Exits should be randomly generated
			- Exits do NOT need a room until necessary
				- If the player looks into an exit, generate the room
				- If there is only 1 exit, generate the room

	- Rooms are not navigable backwards (Ephemeral)

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
	Origin Direction = 0
	North  Direction = 1 << (iota - 1)
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
		return fmt.Sprintf("nil%d", d)
	}
}

type (
	DirectionChances [4]float64
)

var (
	// The chances for exits to be generated in the first room
	// of a dungeon.
	originChances = DirectionChances{100, 100, 100, 100}

	// The chances for an exit to generate. One cardinal direction
	// always generates, but the subsequent 3 directions require
	// the previous to have succeeded. If the 2nd exit fails
	// to generate, then only 1 exit will generate.
	exitChances = DirectionChances{100, 40, 25, 20}

	// The chances for a hallway to generate. Subsequent hallways
	// can only generate if the previous one succeeds. If the
	// first hallway fails, no hallways are generated.
	hallChances = DirectionChances{15, 20, 20, 20}
)

type BaseRoom struct {
	Rooms       []Room
	Exits       []Direction
	Description string
	Direction   Direction
}

type Room struct {
	Description string
	Exits       []Direction
	Halls       []Direction
}

type Exits struct {
	cardinals []Direction
	hallways  []Direction
	mask      Direction
}

type Dungeon struct{}

var (
	cardinals = [4]Direction{North, South, East, West}
	diagonals = [4]Direction{NorthEast, NorthWest, SouthEast, SouthWest}
)

// GenerateExits randomizes the possible exit directions and returns them
// as rooms and hallways.
//
// 🔵 Hallways have a significantly lower chance of appearing
func (d Dungeon) GenerateExits(lastDir Direction) (cardinals, hallways []Direction, exits Direction) {
	cardinals = d.GenExits(lastDir, exitChances)
	hallways = d.GenHalls(lastDir, hallChances)

	for _, r := range cardinals {
		exits |= r
	}

	for _, h := range hallways {
		exits |= h
	}

	return cardinals, hallways, exits
}

func (d Dungeon) GenRoom(lastDir Direction) Room {
	room := Room{}

	if lastDir == Origin {
		room.Exits = d.GenExits(lastDir, originChances)
	} else {
		room.Exits = d.GenExits(lastDir, exitChances)
	}

	if lastDir != Origin {
		room.Halls = d.GenHalls(lastDir, hallChances)
	}

	return room
}

func (d Dungeon) GenDirections(lastDir Direction, dirs [4]Direction, chances DirectionChances) []Direction {
	utils.Shuffle(dirs[:])

	enteredFrom := d.getOppositeDir(lastDir)

	dirCount := 0
	for i := range dirs {
		dir := dirs[i]

		if enteredFrom == dir {
			continue
		}
		if !Chance(chances[dirCount]) {
			break
		}
		dirs[dirCount] = dir
		dirCount++
	}

	if dirCount == 0 {
		return nil
	}

	return dirs[:dirCount]
}

func (d Dungeon) GenExits(lastDir Direction, chances DirectionChances) []Direction {
	return d.GenDirections(lastDir, cardinals, chances)
}

func (d Dungeon) GenHalls(lastDir Direction, chances DirectionChances) []Direction {
	return d.GenDirections(lastDir, diagonals, chances)
}

func (d Dungeon) getOppositeDir(dir Direction) Direction {
	switch dir {
	case North:
		return South
	case South:
		return North
	case East:
		return West
	case West:
		return East
	case NorthEast:
		return SouthWest
	case SouthWest:
		return NorthEast
	case NorthWest:
		return SouthEast
	case SouthEast:
		return NorthWest
	default:
		return Origin
	}
}
