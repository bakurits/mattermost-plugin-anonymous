package anonymous

import (
	"fmt"
	"sync"
	"time"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

// PluginIdentifier unique plugin identifier
type PluginIdentifier struct {
	ID      string
	Version string
}

type pluginVerificationTracker struct {
	verifiedPlugins map[PluginIdentifier]bool

	unverifiedPluginsList []PluginIdentifier

	// unverifiedPluginsLock guards unverifiedPluginsList
	unverifiedPluginsLock *sync.RWMutex
}

// UnverifiedPlugins returns list of unverified plugins
func (a *anonymous) UnverifiedPlugins() []PluginIdentifier {
	a.pluginVerificationTracker.unverifiedPluginsLock.RLock()
	defer a.pluginVerificationTracker.unverifiedPluginsLock.RUnlock()

	return a.pluginVerificationTracker.unverifiedPluginsList
}

// StartPluginChecks starts checking unverified plugins
func (a *anonymous) StartPluginChecks() error {
	c := cron.New()
	_, err := c.AddFunc("@hourly", func() {
		mlog.Info(fmt.Sprintf("started updating validated plugins %s", time.Now()))

		plugins, err := a.checkPluginsVerificationStatus()
		if err != nil {
			mlog.Error(err.Error())
			return
		}

		a.pluginVerificationTracker.unverifiedPluginsLock.Lock()
		a.pluginVerificationTracker.unverifiedPluginsList = plugins
		a.pluginVerificationTracker.unverifiedPluginsLock.Unlock()
	})
	if err != nil {
		return errors.Wrap(err, "Error while starting plugin verification checker")
	}

	c.Start()

	return nil
}

func (a *anonymous) checkPluginsVerificationStatus() ([]PluginIdentifier, error) {

	activePlugins, err := a.PluginAPI.GetActivePlugins()
	if err != nil {
		return []PluginIdentifier{}, errors.Wrap(err, "Error while checking unverified plugins")
	}

	var plugins []PluginIdentifier

	for _, plugin := range activePlugins {
		if _, ok := a.pluginVerificationTracker.verifiedPlugins[plugin]; !ok {
			plugins = append(plugins, plugin)
		}
	}

	return plugins, nil
}

func (a *anonymous) initializeValidatedPackages() {
	a.pluginVerificationTracker = &pluginVerificationTracker{
		unverifiedPluginsList: []PluginIdentifier{},
		unverifiedPluginsLock: &sync.RWMutex{},

		verifiedPlugins: map[PluginIdentifier]bool{
			{
				ID:      a.GetConfiguration().PluginID,
				Version: a.GetConfiguration().PluginVersion,
			}: true,
		},
	}
}
