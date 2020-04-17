package anonymous

import (
	"reflect"
	"sync"

	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	"github.com/bakurits/mattermost-plugin-anonymous/server/crypto"
	"github.com/bakurits/mattermost-plugin-anonymous/server/store"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

// Anonymous API for business logic
type Anonymous interface {
	PluginAPI
	store.Store

	StorePublicKey(userID string, publicKey crypto.PublicKey) error
	GetPublicKey(userID string) (crypto.PublicKey, error)

	GetConfiguration() *config.Config
	SetConfiguration(configuration *config.Config)
}

// Dependencies contains all API dependencies
type Dependencies struct {
	PluginAPI
	store.Store
}

// Config Anonymous configuration
type Config struct {
	*config.Config
	*Dependencies
}

// PluginAPI API form mattermost plugin
type PluginAPI interface {
	SendEphemeralPost(userID string, post *model.Post) *model.Post
}

type anonymous struct {
	Config
	// configurationLock synchronizes access to the configuration.
	configurationLock *sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	config *config.Config
}

// New returns new Anonymous API object
func New(apiConfig Config) Anonymous {
	return &anonymous{
		Config: apiConfig,
	}
}

//StorePublicKey store public key in plugin's KeyValue Store
func (a *anonymous) StorePublicKey(userID string, publicKey crypto.PublicKey) error {
	return a.StoreUser(&store.User{
		MattermostUserID: userID,
		PublicKey:        publicKey,
	})
}

//GetPublicKey get public key from plugin's KeyValue Store
func (a *anonymous) GetPublicKey(userID string) (crypto.PublicKey, error) {
	user, err := a.LoadUser(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while loading a user %s", userID)
	}
	return user.PublicKey, nil
}

// getConfiguration retrieves the active Config under lock, making it safe to use
// concurrently. The active Config may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (a *anonymous) GetConfiguration() *config.Config {
	a.configurationLock.RLock()
	defer a.configurationLock.RUnlock()

	if a.config == nil {
		return &config.Config{}
	}

	return a.config
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
func (a *anonymous) SetConfiguration(configuration *config.Config) {
	a.configurationLock.Lock()
	defer a.configurationLock.Unlock()

	if configuration != nil && a.config == configuration {
		// Ignore assignment if the Config struct is empty. Go will optimize the
		// allocation for same to point at the same memory address, breaking the check
		// above.
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}

		panic("setConfiguration called with the existing Config")
	}

	a.config = configuration
}
