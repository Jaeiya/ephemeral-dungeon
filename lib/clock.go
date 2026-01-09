package lib

import "fmt"

const (
	//🟠 Should always be a multiple of 60
	secondsPerTick = 15
	ticksPerMinute = 60 / secondsPerTick
	ticksPerHour   = ticksPerMinute * 60
)

type TimePeriod int

const (
	Dawn       TimePeriod = iota //  4:00 AM
	Morning                      //  6:00 AM
	MidMorning                   //  9:00 AM
	Day                          // 12:00 PM
	Afternoon                    //  2:30 PM
	Evening                      //  5:00 PM
	Twilight                     //  7:00 PM
	Night                        //  8:00 PM
)

type WorldClock struct {
	ticks       int
	elapsedTime struct {
		days    int
		hours   int
		minutes int
		seconds int
	}
	relativeTime struct {
		day     int
		hours   int
		minutes int
		seconds int
		isPM    bool
	}
}

func NewWorldClock() WorldClock {
	wc := WorldClock{}
	// Always start game at 6:00 AM
	wc.TickHours(6)
	return wc
}

func (wc *WorldClock) Tick(value int) {
	wc.ticks += value

	et := &wc.elapsedTime
	et.seconds = wc.ticks * secondsPerTick
	et.minutes = et.seconds / 60
	et.hours = et.minutes / 60
	et.days = et.hours / 24

	rt := &wc.relativeTime
	rt.seconds = et.seconds % 60
	rt.minutes = et.minutes % 60
	rt.hours = et.hours % 24
	rt.day = et.days

	if rt.hours < 12 {
		rt.isPM = false
	}

	if rt.hours >= 12 {
		rt.hours -= 12
		if rt.hours == 0 {
			rt.hours = 12
		}
		rt.isPM = true
	}
}

func (wc WorldClock) String() string {
	rt := wc.relativeTime
	amPM := "AM"
	if rt.isPM {
		amPM = "PM"
	}
	return fmt.Sprintf("%d:%02d:%02d %s", rt.hours, rt.minutes, rt.seconds, amPM)
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

func (wc WorldClock) ElapsedTime() (int, int, int, int) {
	et := wc.elapsedTime
	return et.days, et.hours, et.minutes, et.seconds
}

func (wc WorldClock) Time() (int, int, int, int, bool) {
	rt := wc.relativeTime
	return rt.day, rt.hours, rt.minutes, rt.seconds, rt.isPM
}
