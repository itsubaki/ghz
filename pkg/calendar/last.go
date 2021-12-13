package calendar

import (
	"sort"
	"time"
)

type TimeDate struct {
	Start time.Time `json:"start,omitempty"`
	End   time.Time `json:"end,omitempty"`
}

func Last12Weeks() []TimeDate {
	return LastNWeeks(12)
}

func LastNWeeks(n int) []TimeDate {
	return LastNWeeksWith(time.Now(), n)
}

func LastNWeeksWith(now time.Time, n int) []TimeDate {
	zzz := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	laststartday := zzz.AddDate(0, 0, -int(zzz.Weekday()))

	days := make([]time.Time, 0)
	for i := 1; i < n+1; i++ {
		days = append(days, laststartday.AddDate(0, 0, -i*7))
	}

	tmp := make([]TimeDate, 0)
	for _, d := range days {
		tmp = append(tmp, TimeDate{
			Start: d,
			End:   d.AddDate(0, 0, 7),
		})
	}

	out := make([]TimeDate, 0)
	for i := len(tmp) - 1; i > -1; i-- {
		out = append(out, tmp[i])
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Start.After(out[j].Start) })
	return out
}
