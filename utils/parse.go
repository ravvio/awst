package utils

import (
	"regexp"
	"strconv"
	"time"
)

var (
	second = "s"
	minute = "m"
	hour = "h"
	day = "d"
	week = "w"
)

func ParseDuration(duration string) (int64, error) {
	reg, err := regexp.CompilePOSIX( "[0-9]+[smhdw]")
	if err != nil {
		return 0, err
	}

	var t int64 = 0;
	for _, m := range reg.FindAll([]byte(duration), -1) {
		value, err := strconv.ParseInt(string(m[0:len(m)-2]), 10, 64)
		if err != nil {
			return 0, err
		}

		unit := string(m[len(m)-1])
		switch unit {
		case second:
			t -= value * time.Second.Milliseconds()
		case minute:
			t -= value * time.Minute.Milliseconds()
		case hour:
			t -= value * time.Hour.Milliseconds()
		case day:
			t -= value * 24 * time.Hour.Milliseconds()
		case week:
			t -= value * 7 * 24 * time.Hour.Milliseconds()
		}
	}

	return t, nil
}
