package stats

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

type CPUMonitor struct {
	UpdateInterval time.Duration
	Usage          float64
	Prefix         string

	Enabled bool
}

func (cm *CPUMonitor) Toggle() {
	cm.Enabled = !cm.Enabled
}

func (cm *CPUMonitor) Enable() {
	cm.Enabled = true
}

func (cm *CPUMonitor) Disable() {
	cm.Enabled = false
}

func (cm *CPUMonitor) String() string {
	if !cm.Enabled {
		return ""
	}
	return fmt.Sprintf("%s %3v%%", cm.Prefix, int(cm.Usage))
}

func (cm *CPUMonitor) Start(c chan string) error {
	var tmpUsage []float64
	var err error
	for {
		tmpUsage, err = cpu.Percent(cm.UpdateInterval, false)
		if err != nil {
			return err
		}
		go func() {
			cm.Usage = tmpUsage[0]
			c <- fmt.Sprintf("%s", cm)
		}()
	}
}

// ChangeInterval changes the update interval for the monitor
func (cm *CPUMonitor) ChangeInterval(interval int) {
	cm.UpdateInterval = time.Duration(interval) * time.Second
}

func NewCPUMonitor(interval int, enabled bool) (cm *CPUMonitor, err error) {
	cm = new(CPUMonitor)
	cm.UpdateInterval = time.Duration(interval) * time.Second
	cm.Prefix = "CPU:"
	cm.Enabled = enabled
	return cm, err
}
