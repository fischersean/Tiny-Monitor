package stats

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

// CPUMonitor defines the data necessary to read stats from the systems CPU
type CPUMonitor struct {
	// UpdateInterval is the delay between updates on the channel
	UpdateInterval time.Duration

	// Usage is the percent usage of the systems total CPU resources
	Usage float64

	// Prefix is the formatting prefix to be send on the channel via String()
	Prefix string

	// Channel is the channel that data will be sent over
	Channel chan string

	// Enabled is the flag to track is the monitor is still needed
	Enabled bool
}

// Stop disables the monitor and allows for the channel to be safely closed
func (cm *CPUMonitor) Stop() {
	cm.Enabled = false
}

// String provides formatting for the monitor
func (cm *CPUMonitor) String() string {
	return fmt.Sprintf("%s %3v%%", cm.Prefix, int(cm.Usage))
}

// Start begins monitoring and sending data on the provided chan
func (cm *CPUMonitor) Start(c chan string) error {
	var tmpUsage []float64
	var err error

	cm.Channel = c
	for {
		// if the monitor has been stopped, return with a nil error
		if !cm.Enabled {
			return nil
		}
		tmpUsage, err = cpu.Percent(cm.UpdateInterval, false)
		if err != nil {
			return err
		}
		go func() {
			// Test to make sure the reciever still wants data. If not, close and return.
			// The outer function will return on the next loop
			if !cm.Enabled {
				close(cm.Channel)
				return
			}
			cm.Usage = tmpUsage[0]
			cm.Channel <- fmt.Sprintf("%s", cm)
		}()
	}
}

// ChangeInterval changes the update interval for the monitor
func (cm *CPUMonitor) ChangeInterval(interval int) {
	cm.UpdateInterval = time.Duration(interval) * time.Second
}

// NewCPUMonitor creates and returns a new CPUMonitor instance
func NewCPUMonitor(interval int, enabled bool) (cm *CPUMonitor, err error) {
	cm = new(CPUMonitor)
	cm.UpdateInterval = time.Duration(interval) * time.Second
	cm.Prefix = "CPU:"
	cm.Enabled = enabled
	return cm, err
}
