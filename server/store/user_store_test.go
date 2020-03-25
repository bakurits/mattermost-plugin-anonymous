package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/store"
	mock_store "github.com/bakurits/mattermost-plugin-anonymous/server/utils/store/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
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

func checkErr(tassert *assert.Assertions, wantErr bool, err error) {

	if wantErr {
		tassert.Error(err)
	} else {
		tassert.NoError(err)
	}
}

func Test_pluginStore_DeleteUser(t *testing.T) {

	ctrl := gomock.NewController(t)
	tassert := assert.New(t)
	m := mock_store.NewMockKVStore(ctrl)
	m.EXPECT().Delete(stringLikeMatcher("key_in")).Return(nil)
	m.EXPECT().Delete(gomock.Not(stringLikeMatcher("key_in"))).Return(errors.New("no data"))

	defer ctrl.Finish()

	type fields struct {
		userStore store.KVStore
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
			s := &pluginStore{
				userStore: tt.fields.userStore,
			}

			err := s.DeleteUser(tt.args.mattermostUserID)
			checkErr(tassert, tt.wantErr, err)
		})
	}
}

func Test_pluginStore_LoadUser(t *testing.T) {

	ctrl := gomock.NewController(t)
	tassert := assert.New(t)
	m := mock_store.NewMockKVStore(ctrl)
	dt, _ := json.Marshal(User{
		MattermostUserID: "key_in",
		PublicKey:        []byte{1},
	})
	m.EXPECT().Load(stringLikeMatcher("key_in")).Return(dt, nil)
	m.EXPECT().Load(stringLikeMatcher("json_error")).Return([]byte{1}, nil)
	m.EXPECT().Load(stringLikeMatcher("no_data")).Return(nil, errors.New("no data"))
	defer ctrl.Finish()

	type fields struct {
		userStore store.KVStore
	}
	type args struct {
		mattermostUserID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
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
			want: &User{
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
			s := &pluginStore{
				userStore: tt.fields.userStore,
			}
			got, err := s.LoadUser(tt.args.mattermostUserID)
			checkErr(tassert, tt.wantErr, err)

			tassert.Equal(tt.want, got, "returned users must be the same")
		})
	}
}

func Test_pluginStore_StoreUser(t *testing.T) {

	ctrl := gomock.NewController(t)
	tassert := assert.New(t)
	m := mock_store.NewMockKVStore(ctrl)

	m.EXPECT().Store(stringLikeMatcher("user_1"), gomock.Any()).Return(nil)
	m.EXPECT().Store(stringLikeMatcher("cant_store"), gomock.Any()).Return(errors.New("failed plugin KVSet"))

	defer ctrl.Finish()

	type fields struct {
		userStore store.KVStore
	}
	type args struct {
		user *User
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
				user: &User{
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
				user: &User{
					MattermostUserID: "cant_store",
					PublicKey:        nil,
				}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pluginStore{
				userStore: tt.fields.userStore,
			}

			err := s.StoreUser(tt.args.user)
			checkErr(tassert, tt.wantErr, err)
		})
	}
}
