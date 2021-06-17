package scheduler

import (
	"errors"
	"time"

	"github.com/TudorHulban/log"
)

type Month int
type IntervalMonth struct {
	From Month
	To   Month
}

type Week int
type IntervalWeek struct {
	From Week
	To   Week
}

type Day int
type IntervalDay struct {
	From Day
	To   Day
}

type Hour int
type IntervalHour struct {
	From Hour
	To   Hour
}

type Minute int
type IntervalMinute struct {
	From Minute
	To   Minute
}

// Intervals Time intervals for the action. Consider UNIX time is for GMT.
type Intervals struct {
	Months      []IntervalMonth
	DaysOfMonth []IntervalDay
	DaysOfWeek  []IntervalDay
	Weeks       []IntervalWeek
	Hours       []IntervalHour
	Minutes     []IntervalMinute
}

// Restriction Provides restriction intervals.
// If not in provided intervals the action can trigger.
// TODO: add a way to trigger only in last day of month.
type Restriction struct {
	Intervals

	L *log.Logger

	GMTOffset      float64 // to accomodate India for ex.
	translatedTime int64   // UNIX time for request given availability offset
}

type entry struct {
	pos              int
	maxActiveSeconds int64
}

var errRestrict = errors.New("is restricted")

func allRestrict() *Restriction {
	return &Restriction{
		Intervals: Intervals{
			Months: []IntervalMonth{
				IntervalMonth{
					From: 1,
					To:   12,
				},
			},
			Weeks: []IntervalWeek{
				IntervalWeek{
					From: 1,
					To:   52,
				},
			},
			DaysOfMonth: []IntervalDay{
				IntervalDay{
					From: 1,
					To:   31,
				},
			},
			DaysOfWeek: []IntervalDay{
				IntervalDay{
					From: 1,
					To:   7,
				},
			},
			Hours: []IntervalHour{
				IntervalHour{
					From: Hour(0),
					To:   Hour(23),
				},
			},
			Minutes: []IntervalMinute{
				IntervalMinute{
					From: 0,
					To:   59,
				},
			},
		},

		GMTOffset: 3, // Eastern European Standard Time
	}
}

// Check Method returns if the action is restricted to take place.
func (r *Restriction) Check(requestOffset float64, requestTime int64, estimatedDurationSecs uint) bool {
	r.translatedTime = requestTime + int64(requestOffset-r.GMTOffset)*3600 + int64(estimatedDurationSecs)

	return r.checks()
}

// Check Method returns if the action is restricted to take place.
func (r *Restriction) CheckNoOffset(requestTime int64, estimatedDurationSecs uint) bool {
	r.translatedTime = requestTime + int64(estimatedDurationSecs)

	return r.checks()
}

func (r *Restriction) checks() bool {
	if r.checkMonth() {
		r.L.Debug("restricted month - yes\n")
		return true
	}

	if r.checkWeek() {
		r.L.Debug("restricted week - yes\n")
		return true
	}

	if r.checkDayOfMonth() {
		r.L.Debug("restricted day of month - yes\n")
		return true
	}

	if r.checkDayOfWeek() {
		r.L.Debug("restricted day of week - yes\n")
		return true
	}

	if r.checkHour() {
		r.L.Debug("restricted hour - yes\n")
		return true
	}

	if r.checkMinute() {
		r.L.Debug("restricted minute - yes\n")
		return true
	}

	r.L.Debug("check - no restrictions\n")

	return false
}

func (r *Restriction) checkMonth() bool {
	if r.Months == nil {
		r.L.Debug("month - no restriction")
		return false
	}

	reqMonth := Month(time.Unix(r.translatedTime, 0).Month())

	for _, interval := range r.Months {
		r.L.Debugf("Interval from: %v to %v. Request: %v. Translated time: %v.", interval.From, interval.To, reqMonth, r.translatedTime)
		if interval.From <= reqMonth && reqMonth <= interval.To {

			return true
		}
	}

	r.L.Debug("Restriction does not apply to request.")
	return false
}

func (r *Restriction) checkWeek() bool {
	if r.Weeks == nil {
		r.L.Debug("week - no restriction")
		return false
	}

	_, reqWeek := time.Unix(r.translatedTime, 0).ISOWeek()

	for _, interval := range r.Weeks {
		if interval.From <= Week(reqWeek) && Week(reqWeek) <= interval.To {
			return true
		}
	}

	return false
}

func (r *Restriction) checkDayOfMonth() bool {
	if r.DaysOfMonth == nil {
		r.L.Debug("days of month - no restriction")
		return false
	}

	reqDay := Day(time.Unix(r.translatedTime, 0).Day())

	for _, interval := range r.DaysOfMonth {
		if interval.From <= reqDay && reqDay <= interval.To {
			return true
		}
	}

	return false
}

func (r *Restriction) checkDayOfWeek() bool {
	if r.DaysOfWeek == nil {
		r.L.Debug("days of week - no restriction")
		return false
	}

	reqDay := Day(time.Unix(r.translatedTime, 0).Weekday())

	for _, interval := range r.DaysOfWeek {
		if interval.From <= reqDay && reqDay <= interval.To {
			return true
		}
	}

	return false
}

func (r *Restriction) checkHour() bool {
	if r.Hours == nil {
		r.L.Debug("hour - no restriction")
		return false
	}

	reqHour := Hour(time.Unix(r.translatedTime, 0).Hour())

	for _, interval := range r.Hours {
		if interval.From <= reqHour && reqHour <= interval.To {
			// r.L.Debugf("from: %v - passed hour: %v - to: %v", interval.From, reqHour, interval.To)
			return true
		}
	}

	return false
}

func (r *Restriction) checkMinute() bool {
	if r.Minutes == nil {
		r.L.Debug("minute - no restriction")
		return false
	}

	reqMinute := Minute(time.Unix(r.translatedTime, 0).Minute())

	for _, interval := range r.Minutes {
		if interval.From <= reqMinute && reqMinute <= interval.To {
			return true
		}
	}

	return false
}

// func (r *Restriction) secondsToRestrictionMonthly() (int64, error) {
// 	if r.Months == nil {
// 		r.L.Debug("month - no restriction")
// 		return 0, nil
// 	}

// 	req := time.Unix(r.translatedTime, 0)

// 	var res []entry // would contain entries with availability in seconds up when the restriction occurs

// 	for i, interval := range r.Months {
// 		r.L.Debugf("Interval from: %v to %v. Request: %v. Translated time: %v.", interval.From, interval.To, req.Month(), r.translatedTime)

// 		if interval.From <= Month(req.Month()) && Month(req.Month()) <= interval.To {
// 			return 0, errRestrict
// 		}

// 		deltaFrom := date(req.Year(), int(interval.From), 1).Sub(req).Seconds() // seconds from request to interval from
// 		deltaTo := date(req.Year(), 12, 1).Sub(req).Seconds()                   // seconds from request to December 1st same year

// 		res = append(res, entry{
// 			pos:              i,
// 			maxActiveSeconds: int64(math.Max(deltaFrom, deltaTo)),
// 		})
// 	}

// 	// sort res and return smallest available interval

// 	r.L.Debug("Restriction does not apply to request.")
// 	return 0, nil
// }

// // https://yourbasic.org/golang/days-between-dates/
// func date(year, month, day int) time.Time {
// 	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
// }
