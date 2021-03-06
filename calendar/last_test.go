package calendar_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/ghz/calendar"
)

func TestLastNWeeks(t *testing.T) {
	date := calendar.Last12Weeks()
	for i, d := range date {
		fmt.Printf("%v: %v %v\n", i, d.Start, d.End)
	}
}
