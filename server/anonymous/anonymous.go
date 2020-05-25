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

	StorePublicKey(userID string, publicKey crypto.PublicKey) error
	GetPublicKey(userID string) (crypto.PublicKey, error)

	IsEncryptionEnabledForChannel(channelID, userID string) bool
	SetEncryptionStatusForChannel(channelID, userID string, status bool) error

	UnverifiedPlugins() []PluginIdentifier
	StartPluginChecks() error
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

	GetUsersInChannel(channelID, sortBy string, page, perPage int) ([]*model.User, error)
	PublishWebSocketEvent(event string, payload map[string]interface{}, broadcast *model.WebsocketBroadcast)
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

//IsEncryptionEnabledForChannel checks if encryption is enabled for channel
func (a *anonymous) IsEncryptionEnabledForChannel(channelID, userID string) bool {
	return a.IsEncryptionEnabled(channelID, userID)
}

//SetEncryptionStatusForChannel sets new connection status
func (a *anonymous) SetEncryptionStatusForChannel(channelID, userID string, status bool) error {
	if !a.isUserInChannel(channelID, userID) {
		return errors.New("can't find user in channel")
	}

	err := a.SetEncryptionStatus(channelID, userID, status)
	if err != nil {
		return errors.Wrap(err, "error while setting connection status")
	}
	return nil
}

func (a *anonymous) isUserInChannel(channelID, userID string) bool {
	for page := 0; ; page = page + 1 {
		users, err := a.PluginAPI.GetUsersInChannel(channelID, "username", page, 100)
		if err != nil || len(users) == 0 {
			return false
		}
		for _, user := range users {
			if user.Id == userID {
				return true
			}
		}
	}
}
