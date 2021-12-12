package calendar_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/ghstats/pkg/calendar"
)

func TestLastNWeeks(t *testing.T) {
	date := calendar.LastNWeeks(3)
	for i, d := range date {
		fmt.Printf("%v: %v %v\n", i, d.Start, d.End)
	}
}
