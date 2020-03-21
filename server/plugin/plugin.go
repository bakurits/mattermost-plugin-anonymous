package plugin

import (
	"math/rand"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	"github.com/bakurits/mattermost-plugin-anonymous/server/api"
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	"github.com/bakurits/mattermost-plugin-anonymous/server/store"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	httpHandler *api.Handler

	// configurationLock synchronizes access to the configuration.
	configurationLock *sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	config *config.Config
}

// NewWithConfig creates new plugin object from configuration
func NewWithConfig(conf *config.Config) *Plugin {
	return &Plugin{
		configurationLock: &sync.RWMutex{},
		config:            conf,
		httpHandler:       api.NewHTTPHandler(),
	}

}

// OnActivate called when plugin is activated
func (p *Plugin) OnActivate() error {
	rand.Seed(time.Now().UnixNano())
	return nil
}

// getConfiguration retrieves the active Config under lock, making it safe to use
// concurrently. The active Config may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *Plugin) getConfiguration() *config.Config {
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
func (p *Plugin) setConfiguration(configuration *config.Config) {
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
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(config.Config)

	// Load the public Config fields from the Mattermost server Config.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin Config")
	}

	p.setConfiguration(configuration)

	return nil
}

func (p *Plugin) ServeHTTP(pc *plugin.Context, w http.ResponseWriter, req *http.Request) {
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

func (p *Plugin) newAnonymousConfig() anonymous.Config {
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
