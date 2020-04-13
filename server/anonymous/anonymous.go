package anonymous

import (
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	"github.com/bakurits/mattermost-plugin-anonymous/server/crypto"
	"github.com/bakurits/mattermost-plugin-anonymous/server/store"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

// Anonymous API for business logic
type Anonymous interface {
	PluginAPI
	store.Store

	StorePublicKey(publicKey crypto.PublicKey) error
	GetPublicKey(userID string) (crypto.PublicKey, error)
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
	GetActivePlugins() ([]PluginIdentifier, error)
}

type anonymous struct {
	Config
	actingMattermostUserID string
	PluginContext          plugin.Context
	VerifiedPlugins        map[PluginIdentifier]bool
}

// New returns new Anonymous API object
func New(apiConfig Config, mattermostUserID string, ctx plugin.Context) Anonymous {
	return newAnonymous(apiConfig, mattermostUserID, ctx)
}

func newAnonymous(apiConfig Config, mattermostUserID string, ctx plugin.Context) *anonymous {
	a := &anonymous{
		Config:                 apiConfig,
		actingMattermostUserID: mattermostUserID,
		PluginContext:          ctx,
	}

	a.initializeValidatedPackages()

	return a
}

//StorePublicKey store public key in plugin's KeyValue Store
func (a *anonymous) StorePublicKey(publicKey crypto.PublicKey) error {
	return a.StoreUser(&store.User{
		MattermostUserID: a.actingMattermostUserID,
		PublicKey:        publicKey,
	})
}

//GetPublicKey get public key from plugin's KeyValue Store
func (a *anonymous) GetPublicKey(userID string) (crypto.PublicKey, error) {
	user, err := a.LoadUser(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while loading a user %s", a.actingMattermostUserID)
	}
	return user.PublicKey, nil
}
