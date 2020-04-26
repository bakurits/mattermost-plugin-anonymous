package plugin

import "github.com/mattermost/mattermost-server/v5/model"

// SendEphemeralPost responds user request with message
func (p *plugin) SendEphemeralPost(userID string, post *model.Post) *model.Post {
	return p.API.SendEphemeralPost(userID, post)
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
