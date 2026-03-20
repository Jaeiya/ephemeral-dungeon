package engine

import (
	"math"
)

const (
	baseXP     int     = 50   // It takes this much XP to get to level 1
	xpGainRate float64 = 0.13 // 0.10 would be 10% compounding over baseXP per level
	oneThird   float64 = 1.0 / 3.0
)

type Inventory struct {
	Items []Item
}

type Character struct {
	strength     int // Physical formidability
	dexterity    int // Agility, evasion, ranged weapons, dodging...
	intelligence int // Crafting and magic affinity
	vitality     int // Physical resistances and stamina
	charisma     int // Benefits or Deficits to social interactions
	armor        int
	level        int
	inventory    Inventory
	xp           struct {
		total   int
		current int
	}
	location string
}

func NewCharacter() Character {
	cs := Character{
		strength:     3,
		dexterity:    0,
		intelligence: 0,
		vitality:     3,
		charisma:     3,
	}
	return cs
}

func (cs Character) Strength() int {
	return cs.strength + cs.level*3
}

func (cs Character) Dexterity() int {
	return cs.dexterity + cs.level*3
}

func (cs Character) Intelligence() int {
	return cs.intelligence + cs.level*3
}

func (cs Character) Vitality() int {
	return cs.vitality + cs.level*5
}

func (cs Character) Charisma() int {
	return cs.charisma
}

func (cs Character) Level() int {
	return cs.level
}

func (cs *Character) AddStrength(str int) {
	cs.strength += str
}

func (cs *Character) AddDexterity(dex int) {
	cs.dexterity += dex
}

func (cs *Character) AddIntelligence(in int) {
	cs.intelligence += in
}

func (cs *Character) AddVitality(vit int) {
	cs.vitality += vit
}

func (cs *Character) AddCharisma(cha int) {
	cs.charisma += cha
}

func (cs *Character) AddXP(xp int) {
	xpToLevel := calcXPToLevel(cs.level)
	currentXP := cs.xp.current + xp

	for currentXP >= xpToLevel {
		currentXP -= xpToLevel
		cs.level += 1
		xpToLevel = calcXPToLevel(cs.level)
	}

	cs.xp.current = currentXP
	cs.xp.total += xp
}

func (cs Character) Agility() int {
	return int(math.Round(oneThird * float64(cs.dexterity)))
}

func (cs Character) Dodge() int {
	a := float64(cs.Agility())
	i := float64(cs.Intelligence())

	rating := 0.15 * (a + 0.25*i)

	return int(math.Round(rating))
}

// Initiative
//
// 🔵 Used to calculate who attacks first during combat.
// 🔵 Used to calculate counter-attack chance during combat.
func (cs Character) Initiative() int {
	i := float64(cs.Intelligence())
	d := float64(cs.Dexterity())

	rating := oneThird * (oneThird*i + d)

	return int(math.Round(rating))
}

// Resistance
//
// 🔵 Used to calculate status-effect resistance
func (cs Character) Resistance() int {
	return int(math.Round(oneThird * float64(cs.vitality)))
}

// Stamina
//
// 🔵 Used to calculate exhaustion during combat.
// 🔵 Used to calculate exhaustion during training.
func (cs Character) Stamina() int {
	rating := 0.75 * float64(cs.vitality)
	return int(math.Round(rating))
}

// HitPoints
//
// 🔵 Used to create life pool for character.
func (cs Character) HitPoints() int {
	return int(math.Round(0.5 * float64(cs.vitality)))
}

// calcXPToLevel uses a simple compounding rate algorithm to return
// the XP required for the specified level.
func calcXPToLevel(lvl int) int {
	const gr = 1 + xpGainRate
	return int(math.Round(float64(baseXP) * math.Pow(gr, float64(lvl))))
}
