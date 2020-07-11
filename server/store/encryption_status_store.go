package store

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/pkg/errors"
)

const (
	// EncryptionStatusStoreKeyPrefix prefix for encryption status data key is kvsotre
	EncryptionStatusStoreKeyPrefix = "es_"
)

// EncryptionStatusStore API for encryption statuses in KVStore
type EncryptionStatusStore interface {
	IsEncryptionEnabled(channelID, userID string) bool
	SetEncryptionStatus(channelID, userID string, status bool) error
}

// IsEncryptionEnabled checks if encryption is enabled for channel and user
func (s *pluginStore) IsEncryptionEnabled(channelID, userID string) bool {
	data, err := s.encryptionStatusStore.Load(fmt.Sprintf("%s%s", EncryptionStatusStoreKeyPrefix, channelID))
	if err != nil {
		return false
	}
	var users []string
	_ = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&users)

	for _, user := range users {
		if user == userID {
			return true
		}
	}
	return false
}

// SetEncryptionStatus changes encryption state
func (s *pluginStore) SetEncryptionStatus(channelID, userID string, status bool) error {
	data, err := s.encryptionStatusStore.Load(fmt.Sprintf("%s%s", EncryptionStatusStoreKeyPrefix, channelID))
	var users []string
	if err == nil {
		_ = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&users)
	}

	ind := -1
	for idx, user := range users {
		if user == userID {
			ind = idx
			break
		}
	}

	if status {
		if ind != -1 {
			return nil
		}
		users = append(users, userID)
	} else {
		if ind == -1 {
			return nil
		}
		users = append(users[:ind], users[ind+1:]...)
	}
	var newData bytes.Buffer
	_ = gob.NewEncoder(&newData).Encode(users)
	err = s.encryptionStatusStore.Store(fmt.Sprintf("%s%s", EncryptionStatusStoreKeyPrefix, channelID), newData.Bytes())
	if err != nil {
		return errors.Wrap(err, "error while storing encryption data")
	}
	return nil
}
