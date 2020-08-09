package stats

import (
	"testing"
	"time"
)

// TestCPUInterfaceCompliant tests whether the cpu monitor implements the stats monitor interface
func TestCPUInterfaceCompliant(t *testing.T) {

	tFunc := func(a Stat) {
		return
	}

	m, _ := NewCPUMonitor(1, true)

	tFunc(m)
}

// TestNewCPUMonitor tests that a new cpu monitor can be created
func TestNewCPUMonitor(t *testing.T) {
	monitor, err := NewCPUMonitor(1, true)

	if err != nil {
		t.Fatal(err)
	}

	if monitor.Enabled != true {
		t.Fatalf("Expected Enabled=true, found %v", monitor.Enabled)
	}

	if monitor.Prefix != "CPU:" {
		t.Fatalf("Expected Prefix=CPU:, found %v", monitor.Prefix)
	}

	if monitor.UpdateInterval != time.Duration(1)*time.Second {
		t.Fatalf("Expected UpdateInterval=5s, found %v", monitor.UpdateInterval)
	}

}

// TestCPUStart tests whether the stats monitor can be started and stopped
func TestCPUStart(t *testing.T) {

	monitor, err := NewCPUMonitor(1, true)

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

