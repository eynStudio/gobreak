package utils

import (
	"fmt"
	"strconv"
	"time"
)

func Today() time.Time {
	n := time.Now()
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
}

func FmtYyyyMmDd(t time.Time) string {
	return t.Format("2006年01月02日")
}

func FmtYyyyMm(t time.Time) string {
	return t.Format("2006年01月")
}

func FmtNy(ny string) string {
	if len(ny) != 6 {
		return ""
	}
	y, _ := strconv.Atoi(ny[:4])
	m, _ := strconv.Atoi(ny[4:])
	return fmt.Sprintf("%d年%d月", y, m)
}

func FmtRq(ny string) string {
	if len(ny) != 6 {
		return ""
	}
	y, _ := strconv.Atoi(ny[:4])
	m, _ := strconv.Atoi(ny[4:6])
	d, _ := strconv.Atoi(ny[6:])
	return fmt.Sprintf("%d年%d月%d日", y, m, d)
}
