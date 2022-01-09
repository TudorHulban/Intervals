package scheduler

import (
	"os"
	"testing"
	"time"

	"github.com/TudorHulban/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var l = log.NewLogger(log.DEBUG, os.Stdout, true)

func TestNoRestrictions(t *testing.T) {
	r := Restriction{
		L:         l,
		GMTOffset: 3,
	}

	assert.False(t, r.Check(3, time.Now().Unix(), 0), "no restriction")
}

func TestAllRestrictions(t *testing.T) {
	r := allRestrict()
	r.L = l

	assert.True(t, r.Check(3, time.Now().Unix(), 0), "all restricted")
}

// https://www.unixtimestamp.com/
func TestRestrictions(t *testing.T) {
	// Thu Dec 31 2020 13:00:00 GMT+0200 (Eastern European Standard Time)
	reqTime := int64(1609412400)

	cases := []struct {
		restrict     Restriction
		description  string
		offset       float64
		timestamp    int64
		isRestricted bool
	}{
		{Restriction{L: l, GMTOffset: 3}, "No Restrictions", 0, reqTime, false},
		{Restriction{L: l, Intervals: Intervals{Months: []IntervalMonth{IntervalMonth{From: Month(12), To: Month(12)}}}, GMTOffset: 2}, "Restricted for December", 0, reqTime, true},
		{Restriction{L: l, Intervals: Intervals{Months: []IntervalMonth{IntervalMonth{From: Month(10), To: Month(11)}}}, GMTOffset: 3}, "Available for December", 0, reqTime, false},
		{Restriction{L: l, Intervals: Intervals{DaysOfMonth: []IntervalDay{IntervalDay{From: Day(31), To: Day(31)}, IntervalDay{From: Day(1), To: Day(1)}}}, GMTOffset: 3}, "Restricted for 31st", 0, reqTime, true},
		{Restriction{L: l, Intervals: Intervals{DaysOfMonth: []IntervalDay{IntervalDay{From: Day(1), To: Day(1)}}}, GMTOffset: 3}, "No Restriction for 31st", 0, reqTime, false},
		{Restriction{L: l, Intervals: Intervals{DaysOfWeek: []IntervalDay{IntervalDay{From: Day(4), To: Day(4)}}}, GMTOffset: 3}, "Restricted for Thursday", 0, reqTime, true},
		{Restriction{L: l, Intervals: Intervals{Hours: []IntervalHour{IntervalHour{From: Hour(6), To: Hour(13)}}}, GMTOffset: 3}, "Restricted for morning hours", 0, reqTime, true},
		{Restriction{L: l, Intervals: Intervals{Minutes: []IntervalMinute{IntervalMinute{From: Minute(0), To: Minute(0)}}}, GMTOffset: 3}, "Restricted for sharp", 0, reqTime, true},
	}

	for _, tc := range cases {
		require.Equal(t, tc.restrict.Check(tc.offset, tc.timestamp, 0), tc.isRestricted, tc.description)
	}
}

func TestSpecials1(t *testing.T) {
	// Thu Dec 31 2020 23:00:01 GMT+0200 (Eastern European Standard Time)
	reqTime := int64(1609448401)

	cases := []struct {
		restrict     Restriction
		description  string
		offset       float64
		timestamp    int64
		isRestricted bool
	}{
		{Restriction{L: l, Intervals: Intervals{Months: []IntervalMonth{IntervalMonth{From: Month(12), To: Month(12)}}}, GMTOffset: 2}, "Available for January in Moscow", 3, reqTime, false},
	}

	for _, tc := range cases {
		require.Equal(t, tc.restrict.Check(tc.offset, tc.timestamp, 0), tc.isRestricted, tc.description)
	}
}

func TestSpecials2(t *testing.T) {
	// Mon Dec 28 2020 23:00:01 GMT+0200 (Eastern European Standard Time)
	reqTime := int64(1609189201)

	cases := []struct {
		restrict     Restriction
		description  string
		offset       float64
		timestamp    int64
		isRestricted bool
	}{
		{Restriction{L: l, Intervals: Intervals{DaysOfWeek: []IntervalDay{IntervalDay{From: Day(1), To: Day(1)}}}, GMTOffset: 2}, "Restricted for Monday", 2, reqTime, true},
		{Restriction{L: l, Intervals: Intervals{DaysOfWeek: []IntervalDay{IntervalDay{From: Day(1), To: Day(1)}}}, GMTOffset: 2}, "Available for Tuesday in Moscow", 3, reqTime, false},
		{Restriction{L: l, Intervals: Intervals{Hours: []IntervalHour{IntervalHour{From: Hour(23), To: Hour(23)}}}, GMTOffset: 2}, "Restricted for 23", 2, reqTime, true},
		{Restriction{L: l, Intervals: Intervals{Hours: []IntervalHour{IntervalHour{From: Hour(23), To: Hour(23)}}}, GMTOffset: 2}, "Available for 24 in Moscow", 3, reqTime, false},
		{Restriction{L: l, Intervals: Intervals{Hours: []IntervalHour{IntervalHour{From: Hour(23), To: Hour(23)}}}, GMTOffset: 2}, "Available for 22 in Budapest", 1, reqTime, false},
	}

	for _, tc := range cases {
		require.Equal(t, tc.restrict.Check(tc.offset, tc.timestamp, 0), tc.isRestricted, tc.description)
	}
}

func TestNoOffset(t *testing.T) {
	// Thu Dec 31 2020 23:00:01 GMT+0200 (Eastern European Standard Time)
	reqTime := int64(1609448401)

	cases := []struct {
		restrict     Restriction
		description  string
		timestamp    int64
		isRestricted bool
	}{
		{Restriction{L: l, Intervals: Intervals{Hours: []IntervalHour{IntervalHour{From: Hour(23), To: Hour(23)}}}}, "Restricted for 23", reqTime, true},
	}

	for _, tc := range cases {
		require.Equal(t, tc.restrict.CheckNoOffset(tc.timestamp, 0), tc.isRestricted, tc.description)
	}
}
