package stats

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"math"
	"time"
)

type RAMMonitor struct {
	UpdateInterval time.Duration
	Usage          float64
	Prefix         string
	Relative       bool

	Enabled bool
}

func (rm *RAMMonitor) Toggle() {
	rm.Enabled = !rm.Enabled
}

func (rm *RAMMonitor) Enable() {
	rm.Enabled = true
}

func (rm *RAMMonitor) Disable() {
	rm.Enabled = false
}

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

func (rm *RAMMonitor) Start(c chan string) error {
	var stats *mem.VirtualMemoryStat
	var err error
	for {
		stats, err = mem.VirtualMemory()
		time.Sleep(rm.UpdateInterval)

		if err != nil {
			return err
		}
		go func() {
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

func NewRAMMonitor(interval int, enabled bool, relative bool) (rm *RAMMonitor, err error) {
	rm = new(RAMMonitor)
	rm.UpdateInterval = time.Duration(interval) * time.Second
	rm.Prefix = "RAM:"
	rm.Relative = relative
	rm.Enabled = enabled
	return rm, err
}
