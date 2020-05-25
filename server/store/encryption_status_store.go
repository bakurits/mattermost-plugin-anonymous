package store

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/pkg/errors"
)

const (
	// EncryptionStatusStoreKeyPrefix prefix for encryption status data key is kvsotre
	EncryptionStatusStoreKeyPrefix = "encryption_status_"

	// EncryptionDisabled indicates that encryption is disabled for channel
	EncryptionDisabled byte = 0
	// EncryptionEnabled indicates that encryption is enebled for channel
	EncryptionEnabled byte = 1
)

// EncryptionStatusStore API for encryption statuses in KVStore
type EncryptionStatusStore interface {
	IsEncryptionEnabled(channelID, userID string) bool
	SetEncryptionStatus(channelID, userID string, status bool) error
}

// IsEncryptionEnabled checks if encryption is enabled for channel and user
func (s *pluginStore) IsEncryptionEnabled(channelID, userID string) bool {
	data, err := s.encryptionStatusStore.Load(fmt.Sprintf("%s%s:%s", EncryptionStatusStoreKeyPrefix, channelID, userID))
	if err != nil || len(data) == 0 {
		mlog.Err(err)
		return false
	}
	return data[0] == EncryptionEnabled
}

// SetEncryptionStatus changes encryption state
func (s *pluginStore) SetEncryptionStatus(channelID, userID string, status bool) error {
	var enableIndicator byte
	if status {
		enableIndicator = EncryptionEnabled
	} else {
		enableIndicator = EncryptionDisabled
	}

	err := s.encryptionStatusStore.Store(fmt.Sprintf("%s%s:%s", EncryptionStatusStoreKeyPrefix, channelID, userID), []byte{enableIndicator})
	if err != nil {
		return errors.Wrap(err, "error while storing encryption data")
	}
	return nil
}
