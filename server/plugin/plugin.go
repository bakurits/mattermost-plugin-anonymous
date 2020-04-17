package plugin

import (
	"github.com/bakurits/mattermost-plugin-anonymous/server/api"
	"math/rand"
	"net/http"
	"strings"
	"time"

	mattermostPlugin "github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	"github.com/bakurits/mattermost-plugin-anonymous/server/command"
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	"github.com/bakurits/mattermost-plugin-anonymous/server/store"
	"github.com/pkg/errors"
)

// Plugin is interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin interface {
	anonymous.PluginAPI
	OnActivate() error
	OnConfigurationChange() error
	ServeHTTP(pc *mattermostPlugin.Context, w http.ResponseWriter, r *http.Request)
}

type plugin struct {
	mattermostPlugin.MattermostPlugin

	httpHandler http.Handler

	an anonymous.Anonymous
}

// NewWithConfig creates new plugin object from configuration
func NewWithConfig(conf *config.Config) Plugin {
	p := &plugin{}
	pluginStore := store.NewPluginStore(p.API)
	p.an = anonymous.New(anonymous.Config{
		Config: conf,
		Dependencies: &anonymous.Dependencies{
			Store:     pluginStore,
			PluginAPI: p,
		},
	})
	p.httpHandler = api.NewHTTPHandler(p.an)
	return p
}

// NewWithStore creates new plugin object from configuration and store object
func NewWithStore(mockStore store.Store, conf *config.Config) Plugin {
	p := &plugin{}

	p.an = anonymous.New(anonymous.Config{
		Config: conf,
		Dependencies: &anonymous.Dependencies{
			Store:     mockStore,
			PluginAPI: p,
		},
	})
	p.httpHandler = api.NewHTTPHandler(p.an)
	return p
}

// OnActivate called when plugin is activated
func (p *plugin) OnActivate() error {
	rand.Seed(time.Now().UnixNano())
	err := p.API.RegisterCommand(command.GetSlashCommand())
	if err != nil {
		return errors.Wrap(err, "OnActivate: failed to register command")
	}
	return nil
}

// ExecuteCommand hook is called when slash command is submitted
func (p *plugin) ExecuteCommand(_ *mattermostPlugin.Context, commandArgs *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	mattermostUserID := commandArgs.UserId
	if len(mattermostUserID) == 0 {
		return &model.CommandResponse{}, &model.AppError{Message: "Not authorized"}
	}

	commandHandler := command.NewHandler(commandArgs, p.an)
	args := strings.Fields(commandArgs.Command)

	commandResponse, err := commandHandler.Handle(args...)
	if err == nil {
		return commandResponse, nil
	}

	if appError, ok := err.(*model.AppError); ok {
		return commandResponse, appError
	}

	return commandResponse, &model.AppError{
		Message: err.Error(),
	}

}

// OnConfigurationChange is invoked when Config changes may have been made.
func (p *plugin) OnConfigurationChange() error {
	var configuration = new(config.Config)

	// Load the public Config fields from the Mattermost server Config.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin Config")
	}

	p.an.SetConfiguration(configuration)

	return nil
}

func (p *plugin) ServeHTTP(_ *mattermostPlugin.Context, w http.ResponseWriter, req *http.Request) {
	p.httpHandler.ServeHTTP(w, req)
}
