package minui

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
)

//go:embed input_mappings/*.json
var embeddedInputMappings embed.FS

// GetInputMappingBytes returns the embedded input mapping JSON for the current device.
// Only arm32 Miyoo devices need custom keyboard mappings. All arm64 devices
// (TrimUI, Miyoo Flip, MagicX, GKD Pixel, etc.) use standard SDL controller input.
func GetInputMappingBytes() ([]byte, error) {
	logger := gaba.GetLogger()
	logger.Debug("Detecting MinUI device type", "arch", runtime.GOARCH)

	if runtime.GOARCH != "arm" {
		// arm64 devices use standard SDL controller input
		return nil, nil
	}

	// arm32 = Miyoo Mini / Mini Plus / A30 — needs custom keyboard mapping
	filename := "input_mappings/miyoo.json"

	overridePath := filepath.Join("overrides", "cfw", "minui", filename)
	data, err := os.ReadFile(overridePath)
	if err != nil {
		data, err = embeddedInputMappings.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded input mapping %s: %w", filename, err)
		}
	}

	return data, nil
}
