package store

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/pkg/errors"
)

// API that store uses for interactions with KVStore
type API interface {
	KVGet(key string) ([]byte, *model.AppError)
	KVSet(key string, value []byte) *model.AppError
	KVDelete(key string) *model.AppError
}

// KVStore abstraction for plugin.API.KVStore
type KVStore interface {
	Load(key string) ([]byte, error)
	Store(key string, data []byte) error
	Delete(key string) error
}

type pluginStore struct {
	api API
}

// NewPluginStore creates KVStore from plugin.API
func NewPluginStore(api API) KVStore {
	return &pluginStore{
		api: api,
	}
}

func (s *pluginStore) Load(key string) ([]byte, error) {
	data, appErr := s.api.KVGet(key)
	if appErr != nil {
		return nil, errors.WithMessage(appErr, "failed plugin KVGet")
	}
	if data == nil {
		return nil, errors.New("Error while loading data from KVStore")
	}
	return data, nil
}

func (s *pluginStore) Store(key string, data []byte) error {
	appErr := s.api.KVSet(key, data)
	if appErr != nil {
		return errors.Wrapf(appErr, "Error while storing data with KVStore with key : %q", key)
	}
	return nil
}

func (s *pluginStore) Delete(key string) error {
	appErr := s.api.KVDelete(key)
	if appErr != nil {
		return errors.Wrapf(appErr, "Error while deleting data from KVStore with key : %q", key)
	}
	return nil
}

// LoadJSON load json data from KVStore
func LoadJSON(s KVStore, key string, v interface{}) (returnErr error) {
	data, err := s.Load(key)
	if err != nil {
		return errors.Wrap(err, "Error while loading json")
	}
	return json.Unmarshal(data, v)
}

// SetJSON sets json data in KVStore
func SetJSON(s KVStore, key string, v interface{}) (returnErr error) {
	data, err := json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "Error while storing json")
	}
	return s.Store(key, data)
}
