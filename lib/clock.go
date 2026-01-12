package lib

import "fmt"

const (
	//🟠 Should always be a multiple of 60
	secondsPerTick int = 15
	ticksPerMinute int = 60 / secondsPerTick
	ticksPerHour   int = ticksPerMinute * 60
	hourSeconds    int = 60 * 60
)

type TimePeriod int

const (
	Dawn      TimePeriod = iota //  4:00 AM
	Morning                     //  6:00 AM
	Day                         // 10:00 PM
	Afternoon                   //  2:00 PM
	Evening                     //  5:00 PM
	Night                       //  8:00 PM
)

func (tp TimePeriod) String() string {
	switch tp {
	case Dawn:
		return "Dawn"
	case Morning:
		return "Morning"
	case Day:
		return "Day"
	case Afternoon:
		return "Afternoon"
	case Evening:
		return "Evening"
	case Night:
		return "Night"
	default:
		panic(fmt.Errorf("invalid time period: %d", tp))
	}
}

type WorldClock struct {
	ticks int
}

func NewWorldClock() *WorldClock {
	wc := WorldClock{}
	// Always start game at 6:00 AM
	wc.TickHours(6)
	return &wc
}

func (wc *WorldClock) Tick(value int) {
	wc.ticks += value
}

func (wc *WorldClock) String() string {
	ts := wc.totalSeconds()

	sec := ts % 60
	min := (ts / 60) % 60
	hour := (ts / hourSeconds) % 24

	amPM := "AM"
	if hour >= 12 {
		hour -= 12
		amPM = "PM"
	}

	if hour == 0 {
		hour = 12
	}

	return fmt.Sprintf(
		"%d:%02d:%02d %s (%s)",
		hour,
		min,
		sec,
		amPM,
		wc.TimePeriod(),
	)
}

func (wc *WorldClock) TickMinutes(value int) {
	wc.Tick(ticksPerMinute * value)
}

func (wc *WorldClock) TickHours(value int) {
	wc.Tick(ticksPerHour * value)
}

func (wc *WorldClock) TickDays(value int) {
	wc.Tick(ticksPerHour * 24 * value)
}

func (wc *WorldClock) TimePeriod() TimePeriod {
	h := wc.hour24()

	switch {

	case h >= 4 && h < 6:
		return Dawn

	case h >= 6 && h < 10:
		return Morning

	// Up to 2:59 PM
	case h >= 10 && h < 15:
		return Day

	// Up to 4:59 PM
	case h >= 15 && h < 17:
		return Afternoon

	// Up to 7:59 PM
	case h >= 17 && h < 20:
		return Evening

	default:
		// Captures 8:00 PM to 3:59 AM
		return Night
	}
}

func (wc *WorldClock) totalSeconds() int {
	return wc.ticks * secondsPerTick
}

func (wc *WorldClock) hour24() int {
	return (wc.totalSeconds() / hourSeconds) % 24
}
