package store

import (
	"fmt"

	"github.com/bakurits/mattermost-plugin-anonymous/server/crypto"
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/store"
	"github.com/pkg/errors"
)

// UserStoreKeyPrefix prefix for user data key is kvsotre
const UserStoreKeyPrefix = "user_"

// UserStore API for user KVStore
type UserStore interface {
	LoadUser(mattermostUserID string) (*User, error)
	StoreUser(user *User) error
	DeleteUser(mattermostUserID string) error
}

// User stores user specific data
type User struct {
	MattermostUserID string           `json:"mattermost_user_id"`
	PublicKey        crypto.PublicKey `json:"public_key"`
}

func (s *pluginStore) LoadUser(mattermostUserID string) (*User, error) {
	user := &User{}
	err := store.LoadJSON(s.userStore, fmt.Sprintf("%s%s", UserStoreKeyPrefix, mattermostUserID), user)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while loading a user with id : %s", mattermostUserID)
	}
	return user, nil
}

func (s *pluginStore) StoreUser(user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	err := store.SetJSON(s.userStore, fmt.Sprintf("%s%s", UserStoreKeyPrefix, user.MattermostUserID), user)
	if err != nil {
		return errors.Wrap(err, "Error while storing user")
	}
	return nil
}

func (s *pluginStore) DeleteUser(mattermostUserID string) error {
	err := s.userStore.Delete(fmt.Sprintf("%s%s", UserStoreKeyPrefix, mattermostUserID))
	if err != nil {
		return errors.Wrapf(err, "Error while deleting a user with id : %s", mattermostUserID)
	}
	return nil
}
