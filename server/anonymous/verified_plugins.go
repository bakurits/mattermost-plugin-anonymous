package anonymous

import "github.com/pkg/errors"

// PluginIdentifier unique plugin identifier
type PluginIdentifier struct {
	ID      string
	Version string
}

func (a *anonymous) UnverifiedPlugins() ([]PluginIdentifier, error) {

	activePlugins, err := a.PluginAPI.GetActivePlugins()
	if err != nil {
		return []PluginIdentifier{}, errors.Wrap(err, "Error while checking unverified plugins")
	}

	var plugins []PluginIdentifier

	for _, plugin := range activePlugins {
		if _, ok := a.VerifiedPlugins[plugin]; !ok {
			plugins = append(plugins, plugin)
		}
	}

	return plugins, nil
}

func (a *anonymous) initializeValidatedPackages() {

	a.VerifiedPlugins = map[PluginIdentifier]bool{
		{
			ID:      "",
			Version: "",
		}: true,
	}
}
