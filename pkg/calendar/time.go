package calender

import (
	"fmt"
	"time"
)

func parse(value string) (time.Time, error) {
	out, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Now(), fmt.Errorf("parse time=%s: %v", value, err)
	}

	return out, nil
}
