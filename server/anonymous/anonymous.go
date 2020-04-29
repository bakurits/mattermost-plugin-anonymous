package anonymous

import (
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	"github.com/bakurits/mattermost-plugin-anonymous/server/crypto"
	"github.com/bakurits/mattermost-plugin-anonymous/server/store"
	utils_store "github.com/bakurits/mattermost-plugin-anonymous/server/utils/store"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

// Anonymous API for business logic
type Anonymous interface {
	PluginAPI
	store.Store

	StorePublicKey(userID string, publicKey crypto.PublicKey) error
	GetPublicKey(userID string) (crypto.PublicKey, error)

	UnverifiedPlugins() []PluginIdentifier
	StartPluginChecks()
}

// Dependencies contains all API dependencies
type Dependencies struct {
	PluginAPI
	store.Store
}

// Config Anonymous configuration
type Config struct {
	*Dependencies
}

// PluginAPI API form mattermost plugin
type PluginAPI interface {
	SendEphemeralPost(userID string, post *model.Post) *model.Post
	GetActivePlugins() ([]PluginIdentifier, error)
	GetConfiguration() *config.Config
	utils_store.API
}

type anonymous struct {
	Config
	pluginVerificationTracker *pluginVerificationTracker
}

// New returns new Anonymous API object
func New(apiConfig Config) Anonymous {
	return newAnonymous(apiConfig)
}

func newAnonymous(apiConfig Config) *anonymous {
	a := &anonymous{
		Config: apiConfig,
	}
	a.initializeValidatedPackages()
	return a
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
