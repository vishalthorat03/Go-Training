package utils

import (
	"sync"
	"testing"
	"time"
)

func TestMonitorResources(t *testing.T) {
	// Mock function to simulate CPU usage
	mockCPUUsage := func(cpuUsage float64) func() ([]float64, error) {
		return func() ([]float64, error) {
			return []float64{cpuUsage}, nil
		}
	}

	// Override the global cpuUsageFunc with a mock function
	originalCPUUsageFunc := cpuUsageFunc
	defer func() { cpuUsageFunc = originalCPUUsageFunc }() // Restore the original after the test

	// Channels for signaling
	pauseProcessing := make(chan bool, 1)
	stopMonitor := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1) // Wait for the MonitorResources goroutine to finish

	cpuUsageFunc = mockCPUUsage(90.0) // Over CPUThreshold
	go func() {
		defer wg.Done()
		MonitorResources(pauseProcessing, stopMonitor)
	}()

	// Allow time for the goroutine to process
	time.Sleep(300 * time.Millisecond)

	// Verify that processing has been paused due to high CPU usage
	select {
	case paused := <-pauseProcessing:
		if !paused {
			t.Errorf("Expected processing to pause, but it didn't")
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("Test timed out waiting for pause signal")
	}

	// After the pause, change the CPU usage to below 60.0 and check for resume
	cpuUsageFunc = mockCPUUsage(55.0) // Below 60.0 (should trigger resume)

	// Allow more time for the goroutine to process the change and resume
	time.Sleep(800 * time.Millisecond) // Increased sleep time

	// Verify that processing is resumed due to CPU recovery
	select {
	case resumed := <-pauseProcessing:
		if resumed {
			t.Errorf("Expected processing to resume, but it didn't")
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("Test timed out waiting for resume signal")
	}

	// Stop the monitor goroutine
	close(stopMonitor)
	wg.Wait() // Ensure the goroutine finishes cleanly
}
