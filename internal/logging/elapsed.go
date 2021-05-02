package logging

import (
	"time"
)

func Elapsed(what string) func() {
	start := time.Now()
	return func() {
		Info.Printf("%s took %v\n", what, time.Since(start))
	}
}
