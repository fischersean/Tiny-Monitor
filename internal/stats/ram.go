package stats

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"math"
	"time"
)

// RAMMonitor defines the data necessary to read stats from the systems memory
type RAMMonitor struct {
	// UpdateInterval is the delay between updates on the channel
	UpdateInterval time.Duration

	// Usage is the percent usage of the systems total RAM resources
	Usage float64

	// Prefix is the formatting prefix to be send on the channel via String()
	Prefix string

	// Channel is the channel that data will be sent over
	Channel chan string

	// Enabled is the flag to track is the monitor is still needed
	Enabled bool

	// Relative is the flag to determine if the usage will be a percentage or absolute value
	Relative bool
}

// Stop disables the monitor and allows for the channel to be safely closed
func (rm *RAMMonitor) Stop() {
	rm.Enabled = false
}

// String provides formatting for the monitor
func (rm *RAMMonitor) String() string {

	if !rm.Enabled {
		return ""
	}

	var suffix string
	var usStr string
	if rm.Relative {
		suffix = "%"
		usStr = fmt.Sprintf("%2v", int(rm.Usage))
	} else {
		suffix = "GB"
		usStr = fmt.Sprintf("%2.1f", float32(rm.Usage))
	}

	return fmt.Sprintf("%s %v%-2s", rm.Prefix, usStr, suffix)
}

// Start begins monitoring and sending data on the provided chan
func (rm *RAMMonitor) Start(c chan string) error {
	var stats *mem.VirtualMemoryStat
	var err error

	rm.Channel = c
	for {
		if !rm.Enabled {
			return nil
		}
		stats, err = mem.VirtualMemory()
		time.Sleep(rm.UpdateInterval)

		if err != nil {
			return err
		}
		go func() {
			// Test to make sure the reciever still wants data. If not, close and return.
			// The outer function will return on the next loop
			if !rm.Enabled {
				close(rm.Channel)
				return
			}
			if rm.Relative {
				rm.Usage = stats.UsedPercent
			} else {
				rm.Usage = float64(stats.Used) / math.Pow10(9)
			}
			c <- fmt.Sprintf("%s", rm)
		}()

	}
}

// ChangeInterval changes the update interval for the monitor
func (rm *RAMMonitor) ChangeInterval(interval int) {
	rm.UpdateInterval = time.Duration(interval) * time.Second
}

// NewRAMMonitor creates and returns a new RAMMonitor instance
func NewRAMMonitor(interval int, enabled bool, relative bool) (rm *RAMMonitor, err error) {
	rm = new(RAMMonitor)
	rm.UpdateInterval = time.Duration(interval) * time.Second
	rm.Prefix = "RAM:"
	rm.Relative = relative
	rm.Enabled = enabled
	return rm, err
}
