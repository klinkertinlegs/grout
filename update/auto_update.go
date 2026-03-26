package update

import (
	"grout/cfw"
	"grout/internal"
	"grout/romm"
	"sync"
	"sync/atomic"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
)

const updateIcon = "\U000F06B0"

type AutoUpdate struct {
	cfwType         cfw.CFW
	host            *romm.Host
	icon            *gaba.DynamicStatusBarIcon
	running         atomic.Bool
	updateAvailable atomic.Bool
	done            chan struct{}
	mu              sync.Mutex
	releaseChannel  internal.ReleaseChannel
	updateInfo      atomic.Pointer[Info]
}

func NewAutoUpdate(c cfw.CFW, r internal.ReleaseChannel, host *romm.Host) *AutoUpdate {
	return &AutoUpdate{
		cfwType:        c,
		releaseChannel: r,
		host:           host,
		icon:           gaba.NewDynamicStatusBarIcon(""), // Start empty, will show icon if update available
		done:           make(chan struct{}),
	}
}

func (a *AutoUpdate) Icon() gaba.StatusBarIcon {
	return gaba.StatusBarIcon{
		Dynamic: a.icon,
	}
}

func (a *AutoUpdate) Start() {
	a.running.Store(true)
	a.done = make(chan struct{})
	go a.run()
}

func (a *AutoUpdate) IsRunning() bool {
	return a.running.Load()
}

func (a *AutoUpdate) UpdateAvailable() bool {
	return a.updateAvailable.Load()
}

func (a *AutoUpdate) UpdateInfo() *Info {
	return a.updateInfo.Load()
}

// Recheck updates the release channel and re-runs the update check.
// This should be called when the user changes the release channel in settings.
func (a *AutoUpdate) Recheck(releaseChannel internal.ReleaseChannel) {
	if a.running.Load() {
		return // Already running, skip
	}

	a.mu.Lock()
	a.releaseChannel = releaseChannel
	a.mu.Unlock()

	a.updateAvailable.Store(false)
	a.updateInfo.Store(nil)
	a.icon.SetText("") // Clear the icon

	a.Start()
}

func (a *AutoUpdate) run() {
	logger := gaba.GetLogger()
	defer func() {
		a.running.Store(false)
		close(a.done)
	}()

	logger.Debug("AutoUpdate: Checking for updates in background")

	a.mu.Lock()
	channel := a.releaseChannel
	a.mu.Unlock()

	info, err := CheckForUpdate(a.cfwType, channel, a.host)
	if err != nil {
		logger.Debug("AutoUpdate: Failed to check for updates", "error", err)
		return
	}

	a.updateInfo.Store(info)

	if info.UpdateAvailable {
		logger.Debug("AutoUpdate: Update available", "current", info.CurrentVersion, "latest", info.LatestVersion)
		a.updateAvailable.Store(true)
		a.icon.SetText(updateIcon)
	} else {
		logger.Debug("AutoUpdate: Already up to date", "version", info.CurrentVersion)
	}
}
