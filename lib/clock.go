package lib

import (
	"fmt"
	"strings"
)

const (
	//🟠 Should always be a multiple of 60
	secondsPerTick int = 15
	ticksPerMinute int = 60 / secondsPerTick
	ticksPerHour   int = ticksPerMinute * 60
	hourSeconds    int = 60 * 60
	daySeconds     int = hourSeconds * 24
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

func (wc *WorldClock) TickMinutes(value int) {
	wc.Tick(ticksPerMinute * value)
}

func (wc *WorldClock) TickHours(value int) {
	wc.Tick(ticksPerHour * value)
}

func (wc *WorldClock) TickDays(value int) {
	wc.Tick(ticksPerHour * 24 * value)
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

// Returns the elapsed time in a shorthand string format
//
//	`1yr, 5mo, 3d, 4h, 10m`
//	`3d, 0h, 10m`
//	`10m`
//
// 🔵 Will only display time units that have value.
func (wc *WorldClock) ElapsedTime() string {
	ts := wc.totalSeconds()

	min := (ts / 60) % 60
	hour := wc.hour24()
	day := (ts / daySeconds) % 30
	month := (ts / daySeconds / 30) % 12
	year := (ts / daySeconds / 30 / 12)

	sb := strings.Builder{}

	showYear := year >= 1
	showMonth := month >= 1 || showYear
	showDay := day >= 1 || showMonth
	showHour := hour >= 1 || showDay

	if showYear {
		fmt.Fprintf(&sb, "%dyr, ", year)
	}

	if showMonth {
		fmt.Fprintf(&sb, "%dmo, ", month)
	}

	if showDay {
		fmt.Fprintf(&sb, "%dd, ", day)
	}

	if showHour {
		fmt.Fprintf(&sb, "%dh, ", hour)
	}

	fmt.Fprintf(&sb, "%dm", min)
	return sb.String()
}

func (wc *WorldClock) totalSeconds() int {
	return wc.ticks * secondsPerTick
}

func (wc *WorldClock) hour24() int {
	return (wc.totalSeconds() / hourSeconds) % 24
}
