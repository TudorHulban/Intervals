package scheduler_test

import (
	"banners/scheduler"
	"os"
	"sort"
	"testing"

	"github.com/TudorHulban/log"
	"github.com/stretchr/testify/require"
)

type task struct {
	scheduler.Restriction

	description       string
	latestTriggerUNIX int64
	estDurationSecs   uint
}

func TestHowToUse(t *testing.T) {
	// Thu Dec 31 2020 13:00:00 GMT+0200 (Eastern European Standard Time)
	reqTime := int64(1609412400)

	restrict := scheduler.Restriction{
		L:         log.NewLogger(log.DEBUG, os.Stdout, true),
		GMTOffset: 3,
	}

	t1 := task{
		Restriction:       restrict,
		description:       "Task 1",
		estDurationSecs:   0,
		latestTriggerUNIX: 1,
	}
	h1 := scheduler.IntervalHour{From: scheduler.Hour(6), To: scheduler.Hour(14)}
	t1.Intervals.Hours = []scheduler.IntervalHour{h1}

	t2 := task{
		Restriction:       restrict,
		description:       "Task 2",
		estDurationSecs:   7200,
		latestTriggerUNIX: 2,
	}
	h2 := scheduler.IntervalHour{From: scheduler.Hour(6), To: scheduler.Hour(14)}
	t2.Intervals.Hours = []scheduler.IntervalHour{h2}

	tasks := []task{t1, t2}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].latestTriggerUNIX < tasks[j].latestTriggerUNIX
	})

	var ta task

	for i := len(tasks) - 1; i >= 0; i-- {
		tasks[i].L.Debug(tasks[i].description)

		if pass := tasks[i].CheckNoOffset(reqTime, tasks[i].estDurationSecs); !pass {
			ta = tasks[i]
			break
		}
	}

	require.Equal(t, t2.description, ta.description)
}
