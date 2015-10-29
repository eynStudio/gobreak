package log

import (
	"testing"
)

func Test_Log(t *testing.T) {
	l := &log{}
	l.Trace("[TRACE] trace")
	l.Tracef("[TRACE] %s", "tracef")
}
