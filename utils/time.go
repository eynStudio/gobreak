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

func IsSameDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return d1 == d2 && m1 == m2 && y1 == y2
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

func FmtRq(rq string) string {
	if len(rq) != 8 {
		return ""
	}
	y, _ := strconv.Atoi(rq[:4])
	m, _ := strconv.Atoi(rq[4:6])
	d, _ := strconv.Atoi(rq[6:])
	return fmt.Sprintf("%d年%d月%d日", y, m, d)
}
