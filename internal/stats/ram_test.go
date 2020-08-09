package stats

import (
	"testing"
	"time"
)

// TestRAMInterfaceCompliant ensures that ram monitor complies with the stats monitor interface
func TestRAMInterfaceCompliant(t *testing.T) {

	tFunc := func(a Stat) {
		return
	}

	m, _ := NewRAMMonitor(5, true, true)

	tFunc(m)
}

/// TestNewRAMMonitor tests that a new stats monitor can be created
func TestNewRAMMonitor(t *testing.T) {
	monitor, err := NewRAMMonitor(5, true, true)

	if err != nil {
		t.Fatal(err)
	}

	if monitor.Enabled != true {
		t.Fatalf("Expected Enabled=true, found %v", monitor.Enabled)
	}

	if monitor.Prefix != "RAM:" {
		t.Fatalf("Expected Prefix=RAM:, found %v", monitor.Prefix)
	}

	if monitor.Relative != true {
		t.Fatalf("Expected Relative=true, found %v", monitor.Relative)
	}

	if monitor.UpdateInterval != time.Duration(5)*time.Second {
		t.Fatalf("Expected UpdateInterval=5s, found %v", monitor.UpdateInterval)
	}

}

// TestRAMStart tests whether the monitor can be started and stopped
func TestRAMStart(t *testing.T) {

	monitor, err := NewRAMMonitor(1, true, true)

	if err != nil {
		t.Fatal(err)
	}

	c := make(chan string)

	go monitor.Start(c)

	msg := <-c

	if len(msg) == 0 {
		t.Fatal("Could not read from channel")
	}

    monitor.Stop()
    
}

