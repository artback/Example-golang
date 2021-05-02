package logging

import (
	"log"
	"time"
)

func Elapsed(what string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", what, time.Since(start))
	}
}
