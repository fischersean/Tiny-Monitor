package stats

import (
	"fmt"
	"github.com/shirou/gopsutil/net"
	"math"
	"strconv"
	"strings"
	"time"
)

// NetworkMonitor gives the current systems total network (upload and download) measurements
type NetworkMonitor struct {
	UpdateInterval time.Duration
	BytesUp        uint64
	BytesDown      uint64
	Prefix         string
	MagnitudeMap   map[int]string

	Enabled bool
}

func (nu *NetworkMonitor) Toggle() {
	nu.Enabled = !nu.Enabled
}

func (nu *NetworkMonitor) Enable() {
	nu.Enabled = true
}

func (nu *NetworkMonitor) Disable() {
	nu.Enabled = false
}

// formatter ensures that the read-out is always 3 digits for both upload and download
func (nu *NetworkMonitor) String() string {

	if !nu.Enabled {
		return ""
	}

	adjustMagnitude := func(quant uint64, mmap map[int]string) (s string) {
		fquant := float64(quant)
		var adQuant float64
		var mag int
		switch {
		case fquant < math.Pow10(3):
			mag = 1
			adQuant = fquant
		case fquant < math.Pow10(6): // 10^6
			mag = 3
			adQuant = fquant / math.Pow10(mag)
		case fquant < math.Pow10(9):
			mag = 6
			adQuant = fquant / math.Pow10(mag)
		default:
			mag = 9
			adQuant = fquant / math.Pow10(mag)
		}

		s = strconv.FormatFloat(adQuant, 'f', 1, 64)
		s = strings.TrimRight(strings.TrimRight(s, "0"), ".") // Remove trailing 0's
		return fmt.Sprintf("%5s %1vB/s", s, mmap[mag])
	}

	upStr := adjustMagnitude(nu.BytesUp, nu.MagnitudeMap)
	downStr := adjustMagnitude(nu.BytesDown, nu.MagnitudeMap)

	return fmt.Sprintf("↓%5v  ↑%-5v", downStr, upStr)
}

// Start reads the stuff
func (nu *NetworkMonitor) Start(c chan string) error {
	var stats0 []net.IOCountersStat
	var stats1 []net.IOCountersStat
	var err error
	for {
		stats0, err = net.IOCounters(false)
		// Doesnt really match perfectly to interval but i think it is good enough
		time.Sleep(nu.UpdateInterval)
		stats1, err = net.IOCounters(false)

		if err != nil {
			return err
		}
		go func() {
			nu.BytesUp = stats1[0].BytesSent - stats0[0].BytesSent
			nu.BytesDown = stats1[0].BytesRecv - stats0[0].BytesRecv
			c <- fmt.Sprintf("%s", nu)
		}()
	}
}

// NewNetworkMonitor initializes a new network monitor listening for all activity
func NewNetworkMonitor(interval int, enabled bool) (nu *NetworkMonitor, err error) {
	nu = new(NetworkMonitor)
	nu.UpdateInterval = time.Duration(interval) * time.Second
	nu.MagnitudeMap = map[int]string{1: "", 3: "k", 6: "M", 9: "G"}
	nu.Prefix = ""
	nu.Enabled = enabled
	return nu, err
}
