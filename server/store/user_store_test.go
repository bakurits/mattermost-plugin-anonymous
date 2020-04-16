package store_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/bakurits/mattermost-plugin-anonymous/server/store"
	utilStore "github.com/bakurits/mattermost-plugin-anonymous/server/utils/store"
	mockStore "github.com/bakurits/mattermost-plugin-anonymous/server/utils/store/mock"
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type stringLikeMatcher string

func (s stringLikeMatcher) Matches(data interface{}) bool {

	text, ok := data.(string)
	if !ok {
		return false
	}

	pattern := string(s)
	return strings.Contains(text, pattern)

}

func (s stringLikeMatcher) String() string {
	return fmt.Sprintf("should match with strings containging ()")
}

func Test_pluginStore_DeleteUser(t *testing.T) {

	ctrl := gomock.NewController(t)
	tassert := assert.New(t)
	m := mockStore.NewMockKVStore(ctrl)
	m.EXPECT().Delete(stringLikeMatcher("key_in")).Return(nil)
	m.EXPECT().Delete(gomock.Not(stringLikeMatcher("key_in"))).Return(errors.New("no data"))

	defer ctrl.Finish()

	type fields struct {
		userStore utilStore.KVStore
	}
	type args struct {
		mattermostUserID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "key present test",
			fields: fields{
				userStore: m,
			},
			args: args{
				mattermostUserID: "key_in",
			},
			wantErr: false,
		},
		{
			name: "key not present test",
			fields: fields{
				userStore: m,
			},
			args: args{
				mattermostUserID: "key_not_in",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewPluginsStore(tt.fields.userStore)

			err := s.DeleteUser(tt.args.mattermostUserID)
			test.CheckErr(tassert, tt.wantErr, err)
		})
	}
}

func Test_pluginStore_LoadUser(t *testing.T) {

	ctrl := gomock.NewController(t)
	tassert := assert.New(t)
	m := mockStore.NewMockKVStore(ctrl)
	dt, _ := json.Marshal(store.User{
		MattermostUserID: "key_in",
		PublicKey:        []byte{1},
	})
	m.EXPECT().Load(stringLikeMatcher("key_in")).Return(dt, nil)
	m.EXPECT().Load(stringLikeMatcher("json_error")).Return([]byte{1}, nil)
	m.EXPECT().Load(stringLikeMatcher("no_data")).Return(nil, errors.New("no data"))
	defer ctrl.Finish()

	type fields struct {
		userStore utilStore.KVStore
	}
	type args struct {
		mattermostUserID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *store.User
		wantErr bool
	}{
		{
			name: "key present test",
			fields: fields{
				userStore: m,
			},
			args: args{
				mattermostUserID: "key_in",
			},
			want: &store.User{
				MattermostUserID: "key_in",
				PublicKey:        []byte{1},
			},
			wantErr: false,
		},
		{
			name: "json error",
			fields: fields{
				userStore: m,
			},
			args: args{
				mattermostUserID: "json_error",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no data test",
			fields: fields{
				userStore: m,
			},
			args: args{
				mattermostUserID: "no_data",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewPluginsStore(tt.fields.userStore)

			got, err := s.LoadUser(tt.args.mattermostUserID)
			test.CheckErr(tassert, tt.wantErr, err)

			tassert.Equal(tt.want, got, "returned users must be the same")
		})
	}
}

func Test_pluginStore_StoreUser(t *testing.T) {

	ctrl := gomock.NewController(t)
	tassert := assert.New(t)
	m := mockStore.NewMockKVStore(ctrl)

	m.EXPECT().Store(stringLikeMatcher("user_1"), gomock.Any()).Return(nil)
	m.EXPECT().Store(stringLikeMatcher("cant_store"), gomock.Any()).Return(errors.New("failed plugin KVSet"))

	defer ctrl.Finish()

	type fields struct {
		userStore utilStore.KVStore
	}
	type args struct {
		user *store.User
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
				userStore: m,
			},
			args: args{
				user: &store.User{
					MattermostUserID: "1",
					PublicKey:        nil,
				},
			},
			wantErr: false,
		},
		{
			name: "storing empty data",
			fields: fields{
				userStore: m,
			},
			args:    args{user: nil},
			wantErr: true,
		},
		{
			name: "storing bad data",
			fields: fields{
				userStore: m,
			},
			args: args{
				user: &store.User{
					MattermostUserID: "cant_store",
					PublicKey:        nil,
				}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewPluginsStore(tt.fields.userStore)

			err := s.StoreUser(tt.args.user)
			test.CheckErr(tassert, tt.wantErr, err)
		})
	}
}
