package spruce

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
)

const DeviceType = "SPRUCE_DEVICE"

//go:embed input_mappings/*.json
var embeddedInputMappings embed.FS

type Device string

const (
	DeviceA30       Device = "a30"
	DeviceMiyoomini Device = "miyoomini"
	DeviceMiyooFlip Device = "miyoomini-flip"
	DeviceTrimui    Device = "trimui"
	DeviceUnknown   Device = "unknown"
)

// DetectDevice detects the device type when running on Spruce by checking environment variables.
func DetectDevice() Device {
	logger := gaba.GetLogger()
	logger.Info("Detecting spruce device type...")

	switch os.Getenv(DeviceType) {
	case "A30":
		return DeviceA30
	case "MIYOOMINI":
		return DeviceMiyoomini
	case "MIYOOFLIP":
		return DeviceMiyooFlip
	case "TRIMUI":
		return DeviceTrimui
	default:
		logger.Warn("Unknown device type, defaulting to miyoomini")
		return DeviceUnknown
	}
}

// GetInputMappingBytes returns the embedded input mapping JSON for the detected muOS device
func GetInputMappingBytes() ([]byte, error) {
	device := DetectDevice()
	return GetInputMappingBytesForDevice(device)
}

// GetInputMappingBytesForDevice returns the embedded input mapping JSON for a specific device
func GetInputMappingBytesForDevice(device Device) ([]byte, error) {
	var filename string
	switch device {
	case DeviceMiyoomini:
		filename = "input_mappings/miyoo.json"
	case DeviceA30:
		filename = "input_mappings/a30.json"
		//case DeviceTrimui:
		//	filename = "input_mappings/trimui.json"
	}

	if filename == "" {
		return nil, fmt.Errorf("no input mapping found for %s", device)
	}

	overridePath := filepath.Join("overrides", "cfw", "spruce", filename)
	data, err := os.ReadFile(overridePath)
	if err != nil {
		data, err = embeddedInputMappings.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded input mapping %s: %w", filename, err)
		}
	}

	return data, nil
}
