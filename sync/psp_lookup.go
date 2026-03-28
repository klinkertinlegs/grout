package sync

import (
	_ "embed"
	"encoding/json"
	"strings"
	"sync"
)

//go:embed data/psp_gamedb.json
var pspGameDBData []byte

var (
	pspGameDB     map[string]string
	pspGameDBOnce sync.Once
)

func loadPSPGameDB() map[string]string {
	pspGameDBOnce.Do(func() {
		pspGameDB = make(map[string]string)
		json.Unmarshal(pspGameDBData, &pspGameDB)
	})
	return pspGameDB
}

// LookupPSPTitle resolves a PPSSPP save folder name to a game title.
// Folder names are typically in the format "UCUS98662_GameData0" or just "UCUS98662".
// Returns the title and true if found, or empty string and false if not.
func LookupPSPTitle(folderName string) (string, bool) {
	db := loadPSPGameDB()

	// Strip common suffixes like _GameData0, _DATA, etc.
	gameID := folderName
	if idx := strings.Index(gameID, "_"); idx > 0 {
		gameID = gameID[:idx]
	}

	// Normalize: remove dashes
	gameID = strings.ReplaceAll(gameID, "-", "")

	title, ok := db[gameID]
	return title, ok
}
