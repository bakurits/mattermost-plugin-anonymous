package store

import (
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/store"
)

type UserStore interface {
	LoadUser(mattermostUserId string) (*User, error)
	StoreUser(user *User) error
	DeleteUser(mattermostUserId string) error
}

// User specific data
type User struct {
	MattermostUserID string `json:"mattermost_user_id"`
	PublicKey        []byte `json:"public_key"`
}

func (s *pluginStore) LoadUser(mattermostUserId string) (*User, error) {
	user := &User{}
	err := store.LoadJSON(s.userStore, mattermostUserId, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *pluginStore) StoreUser(user *User) error {
	err := store.StoreJSON(s.userStore, user.MattermostUserID, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *pluginStore) DeleteUser(mattermostUserID string) error {
	err := s.userStore.Delete(mattermostUserID)
	if err != nil {
		return err
	}
	return nil
}
