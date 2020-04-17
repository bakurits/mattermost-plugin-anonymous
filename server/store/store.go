package store

import (
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/store"
)

// Store encapsulates all store APIs
type Store interface {
	UserStore
}

type pluginStore struct {
	userStore store.KVStore
}

// NewPluginStore creates Store object from plugin.API
func NewPluginStore(api store.StoreAPI) Store {
	return &pluginStore{
		userStore: store.NewPluginStore(api),
	}
}

//// NewWithStores creates Store object from stores
//func NewWithStores(userStore store.KVStore) Store {
//	return &pluginStore{
//		userStore: userStore,
//	}
//}
