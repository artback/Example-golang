package logging

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func Elapsed(what string) func() {
	start := time.Now()
	return func() {
		log.Info("%s took %v\n", what, time.Since(start))
	}
}
