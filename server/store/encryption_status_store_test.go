package store_test

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/bakurits/mattermost-plugin-anonymous/server/store"
	utilsStore "github.com/bakurits/mattermost-plugin-anonymous/server/utils/store"
	mockStore "github.com/bakurits/mattermost-plugin-anonymous/server/utils/store/mock"
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/test"
	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/assert"
)

func Test_pluginStore_SetEncryptionStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	is := assert.New(t)
	m := mockStore.NewMockAPI(ctrl)
	m.EXPECT().KVGet(stringLikeMatcher(store.EncryptionStatusStoreKeyPrefix+"userIn")).Return(stringsToGob([]string{"userIn"}), nil).AnyTimes()
	m.EXPECT().KVGet(stringLikeMatcher(store.EncryptionStatusStoreKeyPrefix+"multiple_userIn")).Return(stringsToGob([]string{"1", "2", "userIn", "3"}), nil).AnyTimes()
	m.EXPECT().KVGet(stringLikeMatcher(store.EncryptionStatusStoreKeyPrefix+"empty")).Return(stringsToGob([]string{""}), nil).AnyTimes()
	m.EXPECT().KVGet(gomock.Any()).Return([]byte{}, &model.AppError{}).AnyTimes()

	m.EXPECT().KVSet(stringLikeMatcher(store.EncryptionStatusStoreKeyPrefix+"cant_store"), gomock.Any()).Return(&model.AppError{}).AnyTimes()
	m.EXPECT().KVSet(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	defer ctrl.Finish()

	type fields struct {
		storeAPI utilsStore.API
	}
	type args struct {
		userID    string
		channelID string
		status    bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "basic test",
			fields: fields{
				storeAPI: m,
			},
			args: args{
				userID:    "userIn",
				channelID: "userIn",
				status:    true,
			},
			wantErr: false,
		},
		{
			name: "storing bad data",
			fields: fields{
				storeAPI: m,
			},
			args: args{
				userID:    "notIn",
				channelID: "cant_store",
				status:    true,
			},
			wantErr: true,
		},
		{
			name: "turning off",
			fields: fields{
				storeAPI: m,
			},
			args: args{
				userID:    "userIn",
				channelID: "multiple_userIn",
				status:    false,
			},
			wantErr: false,
		},
		{
			name: "turning off with when already is off",
			fields: fields{
				storeAPI: m,
			},
			args: args{
				userID:    "4",
				channelID: "multiple_userIn",
				status:    false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewPluginStore(tt.fields.storeAPI)

			err := s.SetEncryptionStatus(tt.args.channelID, tt.args.userID, tt.args.status)
			test.CheckErr(is, tt.wantErr, err)
		})
	}
}

func stringsToGob(users []string) []byte {
	var data bytes.Buffer
	_ = gob.NewEncoder(&data).Encode(users)
	return data.Bytes()
}

func Test_pluginStore_IsEncryptionEnabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	is := assert.New(t)
	m := mockStore.NewMockAPI(ctrl)

	m.EXPECT().KVGet(stringLikeMatcher(store.EncryptionStatusStoreKeyPrefix+"channelIn")).Return(stringsToGob([]string{"userIn"}), nil).AnyTimes()
	m.EXPECT().KVGet(stringLikeMatcher(store.EncryptionStatusStoreKeyPrefix+"channelInDisabled:userInDisabled")).Return([]byte{}, nil).AnyTimes()
	m.EXPECT().KVGet(gomock.Any()).Return([]byte{}, &model.AppError{}).AnyTimes()

	defer ctrl.Finish()

	type fields struct {
		storeAPI utilsStore.API
	}
	type args struct {
		userID    string
		channelID string
		status    bool
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		isEnabled bool
	}{
		{
			name: "is enabled test",
			fields: fields{
				storeAPI: m,
			},
			args: args{
				userID:    "userIn",
				channelID: "channelIn",
				status:    false,
			},
			isEnabled: true,
		},
		{
			name: "is disabled test",
			fields: fields{
				storeAPI: m,
			},
			args: args{
				userID:    "userInDisabled",
				channelID: "channelInDisabled",
				status:    false,
			},
			isEnabled: false,
		},
		{
			name: "not stored test",
			fields: fields{
				storeAPI: m,
			},
			args: args{
				userID:    "not in",
				channelID: "not in",
				status:    false,
			},
			isEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewPluginStore(tt.fields.storeAPI)

			isEnabled := s.IsEncryptionEnabled(tt.args.channelID, tt.args.userID)
			is.Equal(tt.isEnabled, isEnabled)
		})
	}
}
