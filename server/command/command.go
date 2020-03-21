package command

import (
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	anonymous "github.com/bakurits/mattermost-plugin-anonymous/server/plugin"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type Command struct {
	Context   *plugin.Context
	Args      *model.CommandArgs
	ChannelID string
	Config    *config.Config
	plugin    *anonymous.Plugin

	subCommand string
}
