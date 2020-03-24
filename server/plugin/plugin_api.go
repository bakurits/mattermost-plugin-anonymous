package plugin

import "github.com/mattermost/mattermost-server/v5/model"

func (p *Plugin) SendEphemeralPost(userId string, post *model.Post) *model.Post {
	return p.API.SendEphemeralPost(userId, post)
}
