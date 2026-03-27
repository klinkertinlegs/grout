package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	versionsURL    = "https://grout.romm.app/versions.json"
	defaultTimeout = 30 * time.Second
)

type VersionsFile struct {
	Stable *ChannelRelease            `json:"stable"`
	Beta   *ChannelRelease            `json:"beta"`
	RomM   map[string]*ChannelRelease `json:"romm"`
}

type ChannelRelease struct {
	Version string                   `json:"version"`
	Notes   string                   `json:"notes"`
	Assets  map[string]*ChannelAsset `json:"assets"`
}

type ChannelAsset struct {
	URL    string `json:"url"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

func FetchVersionsFile() (*VersionsFile, error) {
	client := &http.Client{Timeout: defaultTimeout}

	req, err := http.NewRequest(http.MethodGet, versionsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Grout-Updater")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch versions: status %d", resp.StatusCode)
	}

	var versions VersionsFile
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("failed to decode versions file: %w", err)
	}

	return &versions, nil
}
