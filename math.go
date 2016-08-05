package utils

import (
	"math"
	"time"
)

func MidnightBeforeOrEqual(timestamp time.Time) time.Time {
	return timestamp.UTC().Truncate(time.Hour * 24)
}

func DaysBefore(timestamp time.Time, days int) time.Time {
	return timestamp.AddDate(0, 0, -1*days)
}

func Round(val float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}
