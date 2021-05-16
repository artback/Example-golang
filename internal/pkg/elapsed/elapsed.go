package elapsed

import (
	"github.com/go-log/log/info"
	"time"
)

type Log struct {
	info.Info
}

func (l Log) Elapsed(what string) func() {
	start := time.Now()
	return func() {
		l.Infof("%s took %v\n", what, time.Since(start))
	}
}
