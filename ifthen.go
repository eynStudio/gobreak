package gobreak

func IfThenStr(t bool, y, n string) string {
	if t {
		return y
	} else {
		return n
	}
}

func IfThenInt(t bool, y, n int) int {
	if t {
		return y
	} else {
		return n
	}
}

func IfThenF64(t bool, y, n float64) float64 {
	if t {
		return y
	} else {
		return n
	}
}

func IfThenBool(t bool, y, n bool) bool {
	if t {
		return y
	} else {
		return n
	}
}
