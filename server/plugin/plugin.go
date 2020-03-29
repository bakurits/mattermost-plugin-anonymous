package plugin

import (
	mattermostPlugin "github.com/mattermost/mattermost-server/v5/plugin"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	"github.com/bakurits/mattermost-plugin-anonymous/server/api"
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

	// configurationLock synchronizes access to the configuration.
	configurationLock *sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	config *config.Config
}

// NewWithConfig creates new plugin object from configuration
func NewWithConfig(conf *config.Config) Plugin {
	return &plugin{
		configurationLock: &sync.RWMutex{},
		config:            conf,
		httpHandler:       api.NewHTTPHandler(),
	}

}

// OnActivate called when plugin is activated
func (p *plugin) OnActivate() error {
	rand.Seed(time.Now().UnixNano())
	err := p.API.RegisterCommand(command.GetSlashCommand())
	if err != nil {
		return errors.WithMessage(err, "OnActivate: failed to register command")
	}
	return nil
}

// ExecuteCommand hook is called when slash command is submitted
func (p *plugin) ExecuteCommand(c *mattermostPlugin.Context, commandArgs *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	mattermostUserID := commandArgs.UserId
	if len(mattermostUserID) == 0 {
		return &model.CommandResponse{}, &model.AppError{Message: "Not authorized"}
	}

	an := anonymous.New(p.newAnonymousConfig(), mattermostUserID, *c)
	comm := command.New(commandArgs, an)
	args := strings.Fields(commandArgs.Command)
	var commandResponse *model.CommandResponse
	var err error
	if len(args) == 0 || args[0] != "/Anonymous" {
		commandResponse, err = comm.Help()
	} else {
		commandResponse, err = comm.Handle(args[1:]...)
	}

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

// getConfiguration retrieves the active Config under lock, making it safe to use
// concurrently. The active Config may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *plugin) getConfiguration() *config.Config {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.config == nil {
		return &config.Config{}
	}

	return p.config
}

// setConfiguration replaces the active Config under lock.
//
// Do not call setConfiguration while holding the configurationLock, as sync.Mutex is not
// re-entrant. In particular, avoid using the plugin API entirely, as this may in turn trigger a
// hook back into the plugin. If that hook attempts to acquire this lock, a deadlock may occur.
//
// This method panics if setConfiguration is called with the existing Config. This almost
// certainly means that the Config was modified without being cloned and may result in
// an unsafe access.
func (p *plugin) setConfiguration(configuration *config.Config) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.config == configuration {
		// Ignore assignment if the Config struct is empty. Go will optimize the
		// allocation for same to point at the same memory address, breaking the check
		// above.
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}

		panic("setConfiguration called with the existing Config")
	}

	p.config = configuration
}

// OnConfigurationChange is invoked when Config changes may have been made.
func (p *plugin) OnConfigurationChange() error {
	var configuration = new(config.Config)

	// Load the public Config fields from the Mattermost server Config.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin Config")
	}

	p.setConfiguration(configuration)

	return nil
}

func (p *plugin) ServeHTTP(pc *mattermostPlugin.Context, w http.ResponseWriter, req *http.Request) {
	mattermostUserID := req.Header.Get("Mattermost-User-ID")
	if len(mattermostUserID) == 0 {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
	}

	apiConf := p.newAnonymousConfig()

	ctx := req.Context()
	ctx = config.Context(ctx, p.config)
	ctx = anonymous.Context(ctx, anonymous.New(apiConf, mattermostUserID, *pc))
	p.httpHandler.ServeHTTP(w, req.WithContext(ctx))
}

func (p *plugin) newAnonymousConfig() anonymous.Config {
	conf := p.getConfiguration()
	pluginStore := store.NewPluginStore(p.API)

	return anonymous.Config{
		Config: conf,
		Dependencies: &anonymous.Dependencies{
			Store:     pluginStore,
			PluginAPI: p,
		},
	}
}
