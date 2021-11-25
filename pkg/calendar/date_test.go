package calendar_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/prstats/pkg/calendar"
)

func TestLastNWeeks(t *testing.T) {
	date := calendar.LastNWeeks(3)
	for _, d := range date {
		fmt.Printf("%v %v %v\n", d.Period, d.Start, d.End)
	}
}
