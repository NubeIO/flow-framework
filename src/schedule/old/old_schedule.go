package old

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"
)

/*
Implements functionality to define a weekly prototypical schedule in JSON format. You can then
https://github.com/rickar/cal TODO add in exceptions


TODO add in a master/global schedule feature

Existing JSON
{"events":{"2ee2a7b9-cf34-4b6d-b9d3-819c33d86040":{"name":"Back-Zone-1","dates":[{"start":"2021-08-21T13:45:00.000Z","end":"2021-08-21T17:45:00.000Z"}],"value":20,"color":"#9013fe"}},"weekly":{"33874fa3-dc9b-42de-a358-5bf83b2b1f1e":{"name":"Front-Zone-1","days":["sunday","monday","tuesday","wednesday","thursday","friday","saturday"],"start":"17:00","end":"03:00","value":20,"color":"#d0021b"}},"exception":{}}

*/

// AnsiFormat The 'time' standard library has this wacky construct in which date parsing relies on specific
const AnsiFormat = "1504"

// IntervalJSONEncoding Struct used in the decoding of the JSON-encoded schedule
type IntervalJSONEncoding struct {
	Start    string `json:"start"`
	Duration string `json:"duration"`
	Name     string `json:"name"`
}

// IntervalValue Identifies unique intervals in the ScheduleDefinition
type IntervalValue struct {
	start time.Time
	end   time.Time
	name  string
}

// ByStart Declarations and functions required for IntervalValue sorting
type ByStart []IntervalValue

func (a ByStart) Len() int           { return len(a) }
func (a ByStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStart) Less(i, j int) bool { return a[i].start.Before(a[j].start) }

func (iv IntervalValue) String() string {
	return fmt.Sprintf("%s: %s ---> %s", iv.name, iv.start, iv.end)
}

// A Definition represents intervals which comprise a weekly schedule, calculated from a prototypical definition.
// Definition Multiple entries covering the same intervals can be declared, and intervals can overflow into other intervals.
type Definition struct {
	Description string `json:"description"`
	Schedule    struct {
		Mon []IntervalJSONEncoding `json:"mon"`
		Tue []IntervalJSONEncoding `json:"tue"`
		Wed []IntervalJSONEncoding `json:"wed"`
		Thu []IntervalJSONEncoding `json:"thu"`
		Fri []IntervalJSONEncoding `json:"fri"`
		Sat []IntervalJSONEncoding `json:"sat"`
		Sun []IntervalJSONEncoding `json:"sun"`
	}
	currentJSON string
	first       time.Time
	last        time.Time
	Intervals   []IntervalValue
}

// New Construct a new ScheduleDefinition from the supplied json string
func New(jsonString string) (sd *Definition, err error) {
	return NewFromTime(jsonString, time.Now())
}

// NewFromTime Construct a new ScheduleDefinition, using the supplied json string and a specific time
func NewFromTime(jsonString string, t time.Time) (sd *Definition, err error) {
	sd = new(Definition)
	if sd.currentJSON != jsonString {
		sd.Intervals = nil
		sd.currentJSON = jsonString
		err = json.Unmarshal([]byte(sd.currentJSON), sd)
		if err != nil {
			return
		}
	}
	err = sd.loadNew(t)
	return
}

// Loads new Intervals from the defined prototype schedule, based on the time supplied
func (sd *Definition) loadNew(t time.Time) (err error) {
	epoch := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	for epoch.Weekday() != time.Monday {
		epoch = epoch.AddDate(0, 0, -1)
	}
	for _, value := range [...][]IntervalJSONEncoding{sd.Schedule.Mon, sd.Schedule.Tue, sd.Schedule.Wed, sd.Schedule.Thu, sd.Schedule.Fri, sd.Schedule.Sat, sd.Schedule.Sun} {
		for _, interval := range value {
			var start, base time.Time
			base, _ = time.Parse(AnsiFormat, "0000")
			start, err = time.Parse(AnsiFormat, interval.Start)
			if err != nil {
				sd = nil
				return
			}
			seconds := start.Sub(base)
			eventStart := epoch.Add(seconds)
			var duration time.Duration
			duration, err = time.ParseDuration(interval.Duration)
			if err != nil {
				sd = nil
				return
			}
			duration -= 1 * time.Second
			eventEnd := eventStart.Add(duration)
			sd.Intervals = append(sd.Intervals, IntervalValue{eventStart, eventEnd, interval.Name})
		}
		epoch = epoch.Add(24 * time.Hour)
	}
	sort.Sort(ByStart(sd.Intervals))
	sd.first = sd.Intervals[0].start
	sd.last = sd.Intervals[len(sd.Intervals)-1].end
	return
}

// String implementation for printing
func (sd *Definition) String() string {
	var s string
	s = fmt.Sprintf("Schedule Defintion: %s, %d intervals, \nspan %s -> %s\n", sd.Description, len(sd.Intervals), sd.first, sd.last)
	for _, v := range sd.Intervals {
		s += fmt.Sprintf("%s\n", v)
	}
	return s
}

// GetDescription implementation for printing
func (sd *Definition) GetDescription() string {
	return sd.Description
}

// Within Determines if the supplied time is within the defined prototypical schedule, using Second precision. This function
// allocates new intervals if the supplied time is outside the intervals currently calculated by previous lookups.
func (sd *Definition) Within(t time.Time) bool {
	if t.After(sd.last) || t.Before(sd.first) {
		err := sd.loadNew(t)
		if err != nil {
			log.Println("error on load")
		}
	}
	for _, v := range sd.Intervals {
		if (t.After(v.start) || t.Unix() == v.start.Unix()) && (t.Before(v.end) || t.Unix() == v.end.Unix()) {
			return true
		}
	}
	return false
}

// MatchingIntervals Returns all matching intervals in the weekly schedule based on the supplied time. schedule allows overlapping
// intervals (events)–you can differentiate intervals using the "name" attribute.
func (sd *Definition) MatchingIntervals(t time.Time) []IntervalValue {
	var a []IntervalValue
	for _, v := range sd.Intervals {
		if (t.After(v.start) || t.Unix() == v.start.Unix()) && (t.Before(v.end) || t.Unix() == v.end.Unix()) {
			a = append(a, v)
		}
	}
	return a
}
