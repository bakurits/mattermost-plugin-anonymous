package command

import (
	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Command stores command specific information
type Command struct {
	Context   *plugin.Context
	Args      *model.CommandArgs
	anonymous *anonymous.Anonymous
}
