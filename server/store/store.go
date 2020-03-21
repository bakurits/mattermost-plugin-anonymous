package store

import (
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/store"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type Store interface {
	UserStore
}

type pluginStore struct {
	userStore store.KVStore
}

func NewPluginStore(api plugin.API) Store {
	return &pluginStore{
		userStore: store.NewPluginStore(api),
	}
}
