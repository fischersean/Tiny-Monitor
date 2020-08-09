package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestSetDefault(t *testing.T) {

	setToDefault()

	if CPUUpdateInterval != 1 {
		t.Fatalf("Expected default cpu interval to be 1, got %d", CPUUpdateInterval)
	}

	if RAMUpdateInterval != 5 {
		t.Fatalf("Expected default cpu interval to be 5, got %d", RAMUpdateInterval)
	}

	if RAMDisplayMode != 1 {
		t.Fatalf("Expected default ram display mode to be 1, got %d", RAMDisplayMode)
	}

	if MinInterval != 1 {
		t.Fatalf("Expected default min interval to be 1, got %d", MinInterval)
	}

}

func TestInitFromFile(t *testing.T) {

	createTestData := func(f string) {
		// Create a test file
		sConfig := saveConfig{
			CPUUpdateInterval: 6,
			RAMUpdateInterval: 6,
			RAMDisplayMode:    6,
		}

		b, err := json.Marshal(sConfig)

		if err != nil {
			t.Fatal("Could not marshal test data")
		}

		file, err := os.Create(f)

		if err != nil {
			t.Fatal("Could not create test data")
		}

		_, err = file.Write(b)

	}

	tPath := "testconfig.json"
	createTestData(tPath)

	InitFromFile(tPath)

	if CPUUpdateInterval != 6 {
		t.Fatalf("Expected cpu interval to be 6, got %d", CPUUpdateInterval)
	}

	if RAMUpdateInterval != 6 {
		t.Fatalf("Expected ram interval to be 6, got %d", RAMUpdateInterval)
	}

	if RAMDisplayMode != 6 {
		t.Fatalf("Expected ram display mode to be 6, got %d", RAMDisplayMode)
	}

	if MinInterval != 1 {
		t.Fatalf("Expected min interval to be 1, got %d", MinInterval)
	}

	// Cleanup
	os.Remove(tPath)

}

func TestSave(t *testing.T) {

	SaveLocation = "testconfig.json"

	if err := Save(); err != nil {
		t.Fatal(err)
	}

	_, err := os.Stat(SaveLocation)

	if os.IsNotExist(err) {
		t.Fatal("Could not find saved file")
	}

	os.Remove(SaveLocation)

}
