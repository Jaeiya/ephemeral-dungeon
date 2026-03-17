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
	Origin Directions = 0
	North  Directions = 1 << (iota - 1)
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

type CardinalChances [3]float64

var (
	OriginCardinal = CardinalChances{100, 100, 100}
	StdCardinal    = CardinalChances{45, 30, 20}
)

var roomBuf = make([]Room, 8)

type BaseRoom struct {
	Rooms       []Room
	Exits       []Directions
	Description string
	Direction   Directions
}

type Room struct {
	Description string
	Direction   Directions
}

type Exits struct {
	cardinals []Directions
	hallways  []Directions
	mask      Directions
}

type Dungeon struct{}

var (
	cardinals = [4]Directions{North, South, East, West}
	hallways  = [4]Directions{NorthEast, NorthWest, SouthEast, SouthWest}
)

// GenerateExits randomizes the possible exit directions and returns them
// as rooms and hallways.
//
// 🔵 Hallways have a significantly lower chance of appearing
func (d Dungeon) GenerateExits(
	lastDir Directions,
) (cardinals, hallways []Directions, exits Directions) {
	cardinals = d.genCardinalDirs(lastDir, StdCardinal)
	hallways = make([]Directions, 0, 4)
	d.genHallDirs(&hallways, lastDir)

	for _, r := range cardinals {
		exits |= r
	}

	for _, h := range hallways {
		exits |= h
	}

	return cardinals, hallways, exits
}

func (d Dungeon) GenRoom(lastDir Directions) BaseRoom {
	var dirArray [8]Directions
	dirs := dirArray[:0]

	if lastDir == Origin {
		dirs = d.genCardinalDirs(lastDir, OriginCardinal)
	} else {
		dirs = d.genCardinalDirs(lastDir, StdCardinal)
	}

	if lastDir != Origin {
		d.genHallDirs(&dirs, lastDir)
	}

	for i := range dirs {
		roomBuf[i] = Room{
			Direction: dirs[i],
		}
	}

	return BaseRoom{
		Direction: lastDir,
		Rooms:     roomBuf[:len(dirs)],
	}
}

// genCardinalDirs generates a slice of possible cardinal exits
// based on the specified chances slice.
func (d Dungeon) genCardinalDirs(lastDir Directions, chances CardinalChances) []Directions {
	dirs := make([]Directions, 0, 4)
	for _, dir := range cardinals {
		if dir == d.getOppositeDir(lastDir) {
			continue
		}
		dirs = append(dirs, dir)
	}

	dirLimit := 1

	for i := 1; i < len(dirs); i++ {
		if RollPercent(chances[i-1]) {
			dirLimit++
		} else {
			break
		}
	}

	utils.Shuffle(dirs)

	return dirs[:dirLimit]
}

func (d Dungeon) genHallDirs(destDirs *[]Directions, lastDir Directions) {
	chances := [4]float64{23, 27, 20, 10}

	if !RollPercent(chances[0]) {
		return
	}

	dirLimit := 1

	var dirsArray [4]Directions
	dirs := dirsArray[:0]

	for _, dir := range hallways {
		if dir == d.getOppositeDir(lastDir) {
			continue
		}
		dirs = append(dirs, dir)
	}

	for i := 1; i < len(dirs); i++ {
		if RollPercent(chances[i]) {
			dirLimit++
		} else {
			break
		}
	}

	utils.Shuffle(dirs)

	*destDirs = append(*destDirs, dirs[:dirLimit]...)
}

func (d Dungeon) getOppositeDir(dir Directions) Directions {
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
