package utils

import (
	"time"
)

func Today() time.Time {
	n := time.Now()
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
}

func FmtYyyyMmDd(t time.Time) string {
	return t.Format("2006年01月02日")
}
