package calendar

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	monthly = "monthly"
	weekly  = "weekly"
	daily   = "daily"
)

type Date struct {
	Period string `json:"period,omitempty"`
	Start  string `json:"start,omitempty"`
	End    string `json:"end,omitempty"`
}

func (d Date) String() string {
	if d.Period == monthly {
		return d.YYYYMM()
	}

	return d.YYYYMMDD()
}

func (d Date) YYYYMM() string {
	return d.Start[:7]
}

func (d Date) YYYYMMDD() string {
	return d.Start
}

func Last(period string) ([]Date, error) {
	n, err := strconv.Atoi(period[:len(period)-1])
	if err != nil {
		return []Date{}, fmt.Errorf("invalid period=%v: %v", period, err)
	}

	var date []Date
	if strings.HasSuffix(period, "m") {
		date = LastNMonths(n)
	}
	if strings.HasSuffix(period, "d") {
		date = LastNDays(n)
	}
	if strings.HasSuffix(period, "w") {
		date = LastNWeeks(n)
	}

	return date, nil
}

func Last12Months() []Date {
	return LastNMonths(12)
}

func LastNMonths(n int) []Date {
	return LastNMonthsWith(time.Now(), n)
}

func LastNMonthsWith(now time.Time, n int) []Date {
	if n < 1 || n > 12 {
		panic(fmt.Sprintf("parameter=%v is not in 0 < n < 13", n))
	}

	months := make([]time.Time, 0)
	for i := 1; i < n+1; i++ {
		months = append(months, now.AddDate(0, -i, -now.Day()+1))
	}

	tmp := make([]Date, 0)
	for _, m := range months {
		tmp = append(tmp, Date{
			Period: monthly,
			Start:  m.Format("2006-01") + "-01",
			End:    m.AddDate(0, 1, 0).Format("2006-01") + "-01",
		})
	}

	out := make([]Date, 0)
	for i := len(tmp) - 1; i > -1; i-- {
		out = append(out, tmp[i])
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Start < out[j].Start })

	return out
}

func Last31Days() []Date {
	return LastNDays(31)
}

func LastNDays(n int) []Date {
	return LastNDaysWith(time.Now(), n)
}

func LastNDaysWith(now time.Time, n int) []Date {
	days := make([]time.Time, 0)
	for i := 1; i < n+1; i++ {
		days = append(days, now.AddDate(0, 0, -i))
	}

	tmp := make([]Date, 0)
	for _, d := range days {
		tmp = append(tmp, Date{
			Period: daily,
			Start:  d.Format("2006-01-02"),
			End:    d.AddDate(0, 0, 1).Format("2006-01-02"),
		})
	}

	out := make([]Date, 0)
	for i := len(tmp) - 1; i > -1; i-- {
		out = append(out, tmp[i])
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Start < out[j].Start })

	return out
}

func Last12Weeks() []Date {
	return LastNWeeks(12)
}

func LastNWeeks(n int) []Date {
	return LastNWeeksWith(time.Now(), n)
}

func LastNWeeksWith(now time.Time, n int) []Date {
	laststartday := now.AddDate(0, 0, -int(now.Weekday()))

	days := make([]time.Time, 0)
	for i := 1; i < n+1; i++ {
		days = append(days, laststartday.AddDate(0, 0, -i*7))
	}

	tmp := make([]Date, 0)
	for _, d := range days {
		tmp = append(tmp, Date{
			Period: weekly,
			Start:  d.Format("2006-01-02"),
			End:    d.AddDate(0, 0, 7).Format("2006-01-02"),
		})
	}

	out := make([]Date, 0)
	for i := len(tmp) - 1; i > -1; i-- {
		out = append(out, tmp[i])
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Start < out[j].Start })

	return out
}

func Parse(value string) (time.Time, error) {
	out, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Now(), fmt.Errorf("parse time=%s: %v", value, err)
	}

	return out, nil
}
