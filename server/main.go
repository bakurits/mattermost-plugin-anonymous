package main

import (
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	anonymous "github.com/bakurits/mattermost-plugin-anonymous/server/plugin"
	mattermost "github.com/mattermost/mattermost-server/v5/plugin"
)

func main() {
	mattermost.ClientMain(
		anonymous.NewWithConfig(
			&config.Config{
				PluginID:      manifest.Id,
				PluginVersion: manifest.Version,
			}))
}
