package engine

import "github.com/jaeiya/monster/internal/utils"

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
		return "nil"
	}
}

type CardinalChances [3]float64

var (
	// Hub room
	originChances = CardinalChances{100, 100, 100}

	// Default Room: chances for 2, 3, & 4 rooms to generate
	// with diminishing returns: 3rd room can't generate
	// unless the 2nd room successfully generated, etc...
	defaultChances = CardinalChances{40, 25, 20}
)

var (
	// The chances for a hallways to generate where each hallway
	// after the first, is dependant on the previous one being
	// successful in order to roll the chance.
	//
	// 0: chance for 1 hallway to generate
	// 1: chance for 2nd hallway if 1st hallway succeeds, etc...
	hallChances = [4]float64{15, 20, 20, 20}
	roomBuf     = [8]Room{}
)

type BaseRoom struct {
	Rooms       []Room
	Exits       []Direction
	Description string
	Direction   Direction
}

type Room struct {
	Description string
	Direction   Direction
}

type Exits struct {
	cardinals []Direction
	hallways  []Direction
	mask      Direction
}

type Dungeon struct{}

var (
	cardinals = [4]Direction{North, South, East, West}
	hallways  = [4]Direction{NorthEast, NorthWest, SouthEast, SouthWest}
)

// GenerateExits randomizes the possible exit directions and returns them
// as rooms and hallways.
//
// 🔵 Hallways have a significantly lower chance of appearing
func (d Dungeon) GenerateExits(
	lastDir Direction,
) (cardinals, hallways []Direction, exits Direction) {
	cardinals = d.genCardinalDirs(lastDir, defaultChances)
	hallways = make([]Direction, 0, 4)
	d.genHallDirs(&hallways, lastDir)

	for _, r := range cardinals {
		exits |= r
	}

	for _, h := range hallways {
		exits |= h
	}

	return cardinals, hallways, exits
}

func (d Dungeon) GenRoom(lastDir Direction) BaseRoom {
	var dirArray [8]Direction
	dirs := dirArray[:0]

	if lastDir == Origin {
		dirs = d.genCardinalDirs(lastDir, originChances)
	} else {
		dirs = d.genCardinalDirs(lastDir, defaultChances)
	}

	if lastDir != Origin {
		d.genHallDirs(&dirs, lastDir)
	}

	for i := range dirs {
		roomBuf[i] = Room{
			// TODO: procedurally generate room description and objects
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
func (d Dungeon) genCardinalDirs(lastDir Direction, chances CardinalChances) []Direction {
	dirs := [4]Direction{}
	dirCount := 0
	for _, dir := range cardinals {
		if dir == d.getOppositeDir(lastDir) {
			continue
		}
		dirs[dirCount] = dir
		dirCount++
	}

	dirLimit := 1
	for ; dirLimit < len(dirs); dirLimit++ {
		if !Chance(chances[dirLimit-1]) {
			break
		}
	}

	utils.Shuffle(dirs[:])

	return dirs[:dirLimit]
}

func (d Dungeon) genHallDirs(destDirs *[]Direction, lastDir Direction) {
	if !Chance(hallChances[0]) {
		return
	}

	dirLimit := 1

	var dirsArray [4]Direction
	dirs := dirsArray[:0]

	for _, dir := range hallways {
		if dir == d.getOppositeDir(lastDir) {
			continue
		}
		dirs = append(dirs, dir)
	}

	for i := 1; i < len(dirs); i++ {
		if Chance(hallChances[i]) {
			dirLimit++
		} else {
			break
		}
	}

	utils.Shuffle(dirs)

	*destDirs = append(*destDirs, dirs[:dirLimit]...)
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
