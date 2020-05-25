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

// GetUsersInChannel gets paginated user list for channel
func (p *plugin) GetUsersInChannel(channelID, sortBy string, page, perPage int) ([]*model.User, error) {
	return p.API.GetUsersInChannel(channelID, sortBy, page, perPage)
}

// PublishWebSocketEvent sends broadcast
func (p *plugin) PublishWebSocketEvent(event string, payload map[string]interface{}, broadcast *model.WebsocketBroadcast) {
	p.API.PublishWebSocketEvent(event, payload, broadcast)
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

// KVGet retrieves a value based on the key, unique per plugin. Returns nil for non-existent keys.
func (p *plugin) KVGet(key string) ([]byte, *model.AppError) {
	return p.API.KVGet(key)
}

// KVSet stores a key-value pair, unique per plugin.
func (p *plugin) KVSet(key string, value []byte) *model.AppError {
	return p.API.KVSet(key, value)
}

// KVDelete removes a key-value pair, unique per plugin. Returns nil for non-existent keys.
func (p *plugin) KVDelete(key string) *model.AppError {
	return p.API.KVDelete(key)
}
