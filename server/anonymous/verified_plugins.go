package anonymous

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/pkg/errors"
)

// PluginIdentifier unique plugin identifier
type PluginIdentifier struct {
	ID      string
	Version string
}

// UnverifiedPlugins returns list of unverified plugins
func (a *anonymous) UnverifiedPlugins() []PluginIdentifier {
	a.unverifiedPluginsLock.RLock()
	defer a.unverifiedPluginsLock.RUnlock()

	return a.unverifiedPluginsList
}

// StartPluginChecks starts checking unverified plugins
func (a *anonymous) StartPluginChecks() {
	go func() {
		for now := range time.Tick(time.Hour) {
			mlog.Info(fmt.Sprintf("started updating validated plugins %s", now.String()))

			plugins, err := a.unverifiedPlugins()
			if err != nil {
				mlog.Error(err.Error())
				return
			}

			a.unverifiedPluginsLock.Lock()
			a.unverifiedPluginsList = plugins
			a.unverifiedPluginsLock.Unlock()
		}
	}()
}

func (a *anonymous) unverifiedPlugins() ([]PluginIdentifier, error) {

	activePlugins, err := a.PluginAPI.GetActivePlugins()
	if err != nil {
		return []PluginIdentifier{}, errors.Wrap(err, "Error while checking unverified plugins")
	}

	var plugins []PluginIdentifier

	for _, plugin := range activePlugins {
		if _, ok := a.verifiedPlugins[plugin]; !ok {
			plugins = append(plugins, plugin)
		}
	}

	return plugins, nil
}

func (a *anonymous) initializeValidatedPackages() {

	a.verifiedPlugins = map[PluginIdentifier]bool{
		{
			ID:      "",
			Version: "",
		}: true,
	}
}
