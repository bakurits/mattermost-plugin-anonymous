package plugin

import "github.com/mattermost/mattermost-server/v5/model"

// SendEphemeralPost responds user request with message
func (p *plugin) SendEphemeralPost(userID string, post *model.Post) *model.Post {
	return p.API.SendEphemeralPost(userID, post)
}
