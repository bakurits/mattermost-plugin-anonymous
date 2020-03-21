package store

import (
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/store"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Store encapsulates all store APIs
type Store interface {
	UserStore
}

type pluginStore struct {
	userStore store.KVStore
}

// NewPluginStore creates Store object from plugin.API
func NewPluginStore(api plugin.API) Store {
	return &pluginStore{
		userStore: store.NewPluginStore(api),
	}
}
