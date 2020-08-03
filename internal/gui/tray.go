package gui

import (
	"fmt"
	"github.com/getlantern/systray"
	"toptray/internal/config"
	"toptray/internal/stats"
)

func toggleStat(item *systray.MenuItem, stat stats.Stat) (err error) {
	stat.Toggle()
	if item.Checked() {
		item.Uncheck()
	} else {
		item.Check()
	}
	return nil
}

// toggleGroup checks the first item and unchecks the others. Saves to global config
func toggleGroup(index int, group map[int]*systray.MenuItem, config *int) {
	for key, value := range group {
		if key == index {
			value.Check()
			*config = key
			continue
		}
		value.Uncheck()
	}
}

// showStats displays the String() of each stat in the order they exist in the slice.
// We are jut going to hardcode the number of stats. Theres probably a more elegant solution
// but I ant to get past this.
func showStats(cm *stats.CPUMonitor, rm *stats.RAMMonitor) {

	tFormat := func(s1 string, s2 string) string {
		return fmt.Sprintf("%s | %s", s1, s2)
	}

	ccm := make(chan string)
	crm := make(chan string)

	go cm.Start(ccm)
	go rm.Start(crm)

	var title string

	cmStat := <-ccm
	rmStat := <-crm
	title = tFormat(cmStat, rmStat)
	systray.SetTitle(title)

	cm.ChangeInterval(config.CPUUpdateInterval)
	rm.ChangeInterval(config.RAMUpdateInterval)

	for {
		select {
		case cmStat = <-ccm:
		case rmStat = <-crm:
		}
		title = tFormat(cmStat, rmStat)
		systray.SetTitle(title)
	}
}

func onReady() {

	// Set both to 1 for initial draw. Will update at first change
	cm, err := stats.NewCPUMonitor(config.MinInterval, true)
	rm, err := stats.NewRAMMonitor(config.MinInterval, true, config.RAMDisplayMode == 0) // 0 is relative

	if err != nil {
		panic(err)
	}

	go showStats(cm, rm)

	// I have no idea why this fixes the drawing issues
	t := systray.AddMenuItem("Update Intervals:", "")
	t.Disable()

	cpuOptList := make(map[int]*systray.MenuItem)
	ramOptList := make(map[int]*systray.MenuItem)
	ramDisplayList := make(map[int]*systray.MenuItem)

	cpuSelect := systray.AddMenuItem("  CPU", "")
	ramSelect := systray.AddMenuItem("  RAM", "")

	ramOptList[1] = ramSelect.AddSubMenuItem("1 s", "")
	ramOptList[2] = ramSelect.AddSubMenuItem("2 s", "")
	ramOptList[5] = ramSelect.AddSubMenuItem("5 s", "")

	ramOptList[config.RAMUpdateInterval].Check()

	ramDisplayMode := ramSelect.AddSubMenuItem("Format", "")

	ramDisplayList[0] = ramDisplayMode.AddSubMenuItem("Percentage", "")
	ramDisplayList[1] = ramDisplayMode.AddSubMenuItem("Absolute", "")

	ramDisplayList[config.RAMDisplayMode].Check()

	cpuOptList[1] = cpuSelect.AddSubMenuItem("1 s", "")
	cpuOptList[2] = cpuSelect.AddSubMenuItem("2 s", "")
	cpuOptList[5] = cpuSelect.AddSubMenuItem("5 s", "")

	cpuOptList[config.CPUUpdateInterval].Check()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	for {
		select {
		case <-cpuOptList[1].ClickedCh:
			cm.ChangeInterval(1)
			toggleGroup(1, cpuOptList, &config.CPUUpdateInterval)

		case <-cpuOptList[2].ClickedCh:
			cm.ChangeInterval(2)
			toggleGroup(2, cpuOptList, &config.CPUUpdateInterval)

		case <-cpuOptList[5].ClickedCh:
			cm.ChangeInterval(5)
			toggleGroup(5, cpuOptList, &config.CPUUpdateInterval)

			//-----------------------------------------------------------
		case <-ramOptList[1].ClickedCh:
			rm.ChangeInterval(1)
			toggleGroup(1, ramOptList, &config.RAMUpdateInterval)

		case <-ramOptList[2].ClickedCh:
			rm.ChangeInterval(2)
			toggleGroup(2, ramOptList, &config.RAMUpdateInterval)

		case <-ramOptList[5].ClickedCh:
			rm.ChangeInterval(5)
			toggleGroup(5, ramOptList, &config.RAMUpdateInterval)

		//-----------------------------------------------------------
		case <-ramDisplayList[0].ClickedCh:
			rm.Relative = true
			toggleGroup(0, ramDisplayList, &config.RAMDisplayMode)

		case <-ramDisplayList[1].ClickedCh:
			rm.Relative = false
			toggleGroup(1, ramDisplayList, &config.RAMDisplayMode)
		}
	}
}

func onExit() {
	// clean up here
	err := config.Save()
	if err != nil {
		// Fail silently
	}
	fmt.Println("Quiting")
}

// Draw renders the systray object
func Draw() {
	systray.Run(onReady, onExit)
}
