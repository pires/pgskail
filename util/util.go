package util

import (
	"io/ioutil"
	"strings"
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

func IsDirEmpty(path string) (bool, error) {
	files, err := ioutil.ReadDir(path)
	return len(files) == 0, err
}

func MakePath(nodes []string) string {
	path := "/"
	if len(nodes) > 0 {
		path = path + strings.Join(nodes, "/")
	}
	return path
}
