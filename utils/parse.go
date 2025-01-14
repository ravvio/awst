package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	second = "s"
	minute = "m"
	hour   = "h"
	day    = "d"
	week   = "w"
)

// Parse a date or datetime in ISO 8601 formats
// Available formats are
// * 2006-01-02
// * 2006-01-02T15:04:05
// * 2006-01-02T15:04:05-0700
func ParseDatetime(value string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05-0700",
	}

	var t = time.Time{}
	var err error
	for _, format := range formats {
		if t, err = time.Parse(format, value); err == nil {
			return t, err
		}
	}
	return t, err
}

func ParseDuration(duration string) (int64, error) {
	reg, err := regexp.CompilePOSIX("[0-9]+[smhdw]")
	if err != nil {
		return 0, err
	}

	var l = 0
	var t int64 = 0
	for _, m := range reg.FindAll([]byte(duration), -1) {
		value, err := strconv.ParseInt(string(m[0:len(m)-1]), 10, 64)
		if err != nil {
			return 0, err
		}

		unit := string(m[len(m)-1])
		switch unit {
		case second:
			t += value * time.Second.Milliseconds()
		case minute:
			t += value * time.Minute.Milliseconds()
		case hour:
			t += value * time.Hour.Milliseconds()
		case day:
			t += value * 24 * time.Hour.Milliseconds()
		case week:
			t += value * 7 * 24 * time.Hour.Milliseconds()
		}

		l += len(m)
	}

	if l != len(duration) {
		return 0, fmt.Errorf("could not parse full duration expression")
	}

	return t, nil
}
