package utils

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const (
	CPUThreshold    = 80
	MemoryThreshold = 90
)

func MonitorResources(pauseProcessing chan bool, stopMonitor <-chan struct{}) {
	for {
		select {
		case <-stopMonitor:
			return
		default:
			cpuUsage, _ := cpu.Percent(0, false)
			memStats, _ := mem.VirtualMemory()

			if cpuUsage[0] > CPUThreshold || memStats.UsedPercent > MemoryThreshold {
				pauseProcessing <- true
			} else {
				select {
				case <-pauseProcessing:
				default:
				}
			}

			time.Sleep(10 * time.Second)
		}
	}
}
