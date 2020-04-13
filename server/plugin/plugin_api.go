package plugin

import (
	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

// SendEphemeralPost responds user request with message
func (p *plugin) SendEphemeralPost(userID string, post *model.Post) *model.Post {
	return p.API.SendEphemeralPost(userID, post)
}

// GetActivePlugins returns list of installed plugins which are active
func (p *plugin) GetActivePlugins() ([]anonymous.PluginIdentifier, error) {
	pluginManifests, err := p.API.GetPlugins()
	if err != nil {
		return []anonymous.PluginIdentifier{}, errors.Wrap(err, "Error while retrieving plugins list")
	}

	var activePlugins []anonymous.PluginIdentifier

	for _, pluginManifest := range pluginManifests {
		activePlugins = append(activePlugins, anonymous.PluginIdentifier{
			ID:      pluginManifest.Id,
			Version: pluginManifest.Version,
		})
	}

	return activePlugins, nil

}
