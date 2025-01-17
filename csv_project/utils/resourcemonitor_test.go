package utils

import (
	"testing"
	"time"
)

func TestMonitorResources(t *testing.T) {
	pause := make(chan bool, 1)
	stop := make(chan struct{})

	go func() {
		time.Sleep(2 * time.Second)
		close(stop)
	}()

	go MonitorResources(pause, stop)
	time.Sleep(1 * time.Second)
	pause <- true
}
