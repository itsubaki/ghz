package calendar

import (
	"sort"
	"time"
)

type Date struct {
	Start time.Time `json:"start,omitempty"`
	End   time.Time `json:"end,omitempty"`
}

func Last12Weeks() []Date {
	return LastNWeeks(12)
}

func LastNWeeks(n int) []Date {
	return LastNWeeksWith(time.Now(), n)
}

func LastNWeeksWith(now time.Time, n int) []Date {
	zzz := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	laststartday := zzz.AddDate(0, 0, -int(zzz.Weekday()))

	days := make([]time.Time, 0)
	for i := 1; i < n+1; i++ {
		days = append(days, laststartday.AddDate(0, 0, -i*7))
	}

	out := make([]Date, 0)
	for _, d := range days {
		out = append(out, Date{
			Start: d,
			End:   d.AddDate(0, 0, 7),
		})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Start.After(out[j].Start) })
	return out
}
