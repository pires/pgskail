package util

import (
	"time"
)

/**
 * Schedule a function to run every _interval_ seconds
 */
func Schedule(what func(), interval uint64) chan struct{} {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				what()
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	return stop
}
