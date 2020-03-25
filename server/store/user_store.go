package store

import (
	"errors"
	"fmt"
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/store"
)

const userStoreKeyPrefix = "user_"

// UserStore API for user KVStore
type UserStore interface {
	LoadUser(mattermostUserID string) (*User, error)
	StoreUser(user *User) error
	DeleteUser(mattermostUserID string) error
}

// User stores user specific data
type User struct {
	MattermostUserID string `json:"mattermost_user_id"`
	PublicKey        []byte `json:"public_key"`
}

func (s *pluginStore) LoadUser(mattermostUserID string) (*User, error) {
	user := &User{}
	err := store.LoadJSON(s.userStore, fmt.Sprintf("%s%s", userStoreKeyPrefix, mattermostUserID), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *pluginStore) StoreUser(user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	err := store.SetJSON(s.userStore, fmt.Sprintf("%s%s", userStoreKeyPrefix, user.MattermostUserID), user)
	if err != nil {
		return err
	}
	return nil
}

func (s *pluginStore) DeleteUser(mattermostUserID string) error {
	err := s.userStore.Delete(fmt.Sprintf("%s%s", userStoreKeyPrefix, mattermostUserID))
	if err != nil {
		return err
	}
	return nil
}
