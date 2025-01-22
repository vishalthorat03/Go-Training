package utils

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

const (
	CPUThreshold  = 80
	RetryInterval = 500 * time.Millisecond // Reduced interval to 500ms for faster detection
)

// cpuUsageFunc is a variable that holds the function to get CPU usage.
// This allows us to override it in tests.
var cpuUsageFunc = func() ([]float64, error) {
	return cpu.Percent(0, false)
}

// MonitorResources monitors CPU resources and signals workers to pause or resume processing.
func MonitorResources(pauseProcessing chan bool, stopMonitor <-chan struct{}) {
	var paused bool

	for {
		select {
		case <-stopMonitor:
			return
		default:
			cpuUsage, _ := cpuUsageFunc()

			// Log current CPU usage for debugging
			fmt.Printf("Current CPU Usage: %.2f\n", cpuUsage[0])

			if cpuUsage[0] > CPUThreshold {
				if !paused {
					pauseProcessing <- true // Signal to pause processing
					paused = true
					fmt.Println("Paused processing due to high CPU usage")
				}
			} else if paused && cpuUsage[0] < 60.0 { // Change recovery condition to 60.0
				pauseProcessing <- false // Signal to resume processing
				paused = false
				fmt.Println("Resumed processing due to low CPU usage")
			}

			time.Sleep(RetryInterval) // Check resources every 500 milliseconds
		}
	}
}
