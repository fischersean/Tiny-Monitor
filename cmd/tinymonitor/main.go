package main

import (
	"fmt"
	"os"
	"toptray/internal/config"
	"toptray/internal/gui"
)

func main() {

    f := fmt.Sprintf("%s/Library/Application Support/TinyMonitor/preferences.json", os.Getenv("HOME"))
    fmt.Println(f)
	config.InitFromFile(f)

	// Set to defaults if we cant load the config file
	gui.Draw()
}
