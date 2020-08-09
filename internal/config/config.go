package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	// CPUUpdateInterval is the user defined time in seconds to delay update
	CPUUpdateInterval int

	// RAMUpdateInterval is the user defined time in seconds to delay update
	RAMUpdateInterval int

	// RAMDisplayMode specifiies whether to show absolute or relative stats. 0 = %, 1 = Absolute
	RAMDisplayMode int

	// SaveLocation is the HDD path to the app data save dir. This is hardcoded in the app's main() func
	SaveLocation string

	// MinInterval is the smallest amount of time allowed between updates
	MinInterval int
)

type saveConfig struct {
	CPUUpdateInterval int `json:"cpu_interval"`
	RAMUpdateInterval int `json:"ram_interval"`
	RAMDisplayMode    int `json:"ram_display"`
}

func setToDefault() {
	CPUUpdateInterval = 1
	RAMUpdateInterval = 5
	RAMDisplayMode = 1
	MinInterval = 1
}

// InitFromFile initializes the app's configuration from the provided JSON file
func InitFromFile(f string) {

	MinInterval = 1
	SaveLocation = f

	if _, err := os.Stat(f); err != nil {
		if os.IsNotExist(err) {
			// Since the file doesnt exist, it is likely the entire dir is missing
			d := filepath.Dir(f)
			os.Mkdir(d, 0777)

		}
		// Set to defaults and return
		setToDefault()
		return
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		setToDefault()
		return
	}

	var cfig saveConfig
	err = json.Unmarshal(b, &cfig)

	if err != nil {
		setToDefault()
		return
	}
	CPUUpdateInterval = cfig.CPUUpdateInterval
	RAMUpdateInterval = cfig.RAMUpdateInterval
	RAMDisplayMode = cfig.RAMDisplayMode
}

// Save writes the current user preferences to the app's save data location
func Save() (err error) {

	f := SaveLocation

	b, err := json.Marshal(saveConfig{
		CPUUpdateInterval: CPUUpdateInterval,
		RAMUpdateInterval: RAMUpdateInterval,
		RAMDisplayMode:    RAMDisplayMode,
	})

	if err != nil {
		return err
	}

	file, err := os.Create(f)

	if err != nil {
		return err
	}

	_, err = file.Write(b)

	return err
}
