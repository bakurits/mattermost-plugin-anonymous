package anonymous_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"

	"github.com/bakurits/mattermost-plugin-anonymous/server/crypto"
	mockPlugin "github.com/bakurits/mattermost-plugin-anonymous/server/plugin/mock"
	"github.com/bakurits/mattermost-plugin-anonymous/server/store"
	mockStore "github.com/bakurits/mattermost-plugin-anonymous/server/store/mock"
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/test"
	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/stretchr/testify/assert"
)

type userIDMatcher string

func (s userIDMatcher) Matches(data interface{}) bool {

	user, ok := data.(*store.User)
	if !ok {
		return false
	}

	pattern := string(s)
	return strings.Contains(user.MattermostUserID, pattern)

}

func (s userIDMatcher) String() string {
	return fmt.Sprintf("should match with strings containing (%s)", string(s))
}

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
	return fmt.Sprintf("should match with strings containing ()")
}

func Test_anonymous_GetPublicKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	tassert := assert.New(t)

	storeMock := mockStore.NewMockStore(ctrl)

	storeMock.EXPECT().LoadUser(stringLikeMatcher("user_in")).Return(&store.User{
		MattermostUserID: "user_in",
		PublicKey:        crypto.PublicKey([]byte{1, 1, 1}),
	}, nil)

	storeMock.EXPECT().LoadUser(gomock.Not(stringLikeMatcher("user_in"))).Return(nil, errors.New("some error"))

	pluginMock := mockPlugin.NewMockPlugin(ctrl)
	pluginMock.EXPECT().GetConfiguration().Return(&config.Config{
		PluginID:      "",
		PluginVersion: "",
	}).AnyTimes()

	defer ctrl.Finish()

	type fields struct {
		Config                 anonymous.Config
		actingMattermostUserID string
		PluginContext          plugin.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    crypto.PublicKey
		wantErr bool
	}{
		{
			name: "basic test",
			fields: fields{
				Config: anonymous.Config{
					Dependencies: &anonymous.Dependencies{
						PluginAPI: pluginMock,
						Store:     storeMock,
					},
				},
				actingMattermostUserID: "user_in",
				PluginContext:          plugin.Context{},
			},
			want:    crypto.PublicKey([]byte{1, 1, 1}),
			wantErr: false,
		},
		{
			name: "test empty",
			fields: fields{
				Config: anonymous.Config{
					Dependencies: &anonymous.Dependencies{
						PluginAPI: pluginMock,
						Store:     storeMock,
					},
				},
				actingMattermostUserID: "user_not_in",
				PluginContext:          plugin.Context{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := anonymous.New(tt.fields.Config)
			got, err := a.GetPublicKey(tt.fields.actingMattermostUserID)
			test.CheckErr(tassert, tt.wantErr, err)
			tassert.Equal(tt.want, got)
		})
	}
}

func Test_anonymous_StorePublicKey(t *testing.T) {

	ctrl := gomock.NewController(t)
	tassert := assert.New(t)

	storeMock := mockStore.NewMockStore(ctrl)

	storeMock.EXPECT().StoreUser(userIDMatcher("user_not_in")).Return(errors.New("some error"))
	storeMock.EXPECT().StoreUser(gomock.Not(userIDMatcher("user_not_in"))).Return(nil)

	pluginMock := mockPlugin.NewMockPlugin(ctrl)
	pluginMock.EXPECT().GetConfiguration().Return(&config.Config{
		PluginID:      "",
		PluginVersion: "",
	}).AnyTimes()

	defer ctrl.Finish()

	type fields struct {
		Config                 anonymous.Config
		actingMattermostUserID string
		PluginContext          plugin.Context
	}
	type args struct {
		publicKey []byte
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
				Config: anonymous.Config{
					Dependencies: &anonymous.Dependencies{
						PluginAPI: pluginMock,
						Store:     storeMock,
					},
				},
				actingMattermostUserID: "user_in",
				PluginContext:          plugin.Context{},
			},
			args: args{
				publicKey: []byte{1, 1, 1},
			},
			wantErr: false,
		},
		{
			name: "test empty",
			fields: fields{
				Config: anonymous.Config{
					Dependencies: &anonymous.Dependencies{
						PluginAPI: pluginMock,
						Store:     storeMock,
					},
				},
				actingMattermostUserID: "user_not_in",
				PluginContext:          plugin.Context{},
			},
			args: args{
				publicKey: []byte{1, 1, 1},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := anonymous.New(tt.fields.Config)
			err := a.StorePublicKey(tt.fields.actingMattermostUserID, tt.args.publicKey)
			test.CheckErr(tassert, tt.wantErr, err)
		})
	}
}

func Test_anonymous_SetEncryptionStatusForChannel(t *testing.T) {

	ctrl := gomock.NewController(t)
	is := assert.New(t)

	storeMock := mockStore.NewMockStore(ctrl)

	storeMock.EXPECT().SetEncryptionStatus("general", "storing_err_user", gomock.Any()).Return(errors.New("some error")).AnyTimes()
	storeMock.EXPECT().SetEncryptionStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	pluginMock := mockPlugin.NewMockPlugin(ctrl)
	pluginMock.EXPECT().GetUsersInChannel("general", gomock.Any(), 0, gomock.Any()).Return([]*model.User{
		{
			Id: "some other",
		},
		{
			Id: "some_other",
		},
	}, nil).AnyTimes()
	pluginMock.EXPECT().GetUsersInChannel("general", gomock.Any(), 1, gomock.Any()).Return([]*model.User{
		{
			Id: "in_general",
		},
		{
			Id: "storing_err_user",
		},
	}, nil).AnyTimes()
	pluginMock.EXPECT().GetUsersInChannel(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*model.User{}, errors.New("some error")).AnyTimes()

	pluginMock.EXPECT().GetConfiguration().Return(&config.Config{
		PluginID:      "",
		PluginVersion: "",
	}).AnyTimes()

	defer ctrl.Finish()

	type fields struct {
		Config anonymous.Config
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
				Config: anonymous.Config{
					Dependencies: &anonymous.Dependencies{
						PluginAPI: pluginMock,
						Store:     storeMock,
					},
				},
			},
			args: args{
				userID:    "storing_err_user",
				channelID: "general",
				status:    false,
			},
			wantErr: true,
		},
		{
			name: "test user not in channel",
			fields: fields{
				Config: anonymous.Config{
					Dependencies: &anonymous.Dependencies{
						PluginAPI: pluginMock,
						Store:     storeMock,
					},
				},
			},
			args: args{
				userID:    "not_in_general",
				channelID: "general",
				status:    false,
			},
			wantErr: true,
		},
		{
			name: "success change",
			fields: fields{
				Config: anonymous.Config{
					Dependencies: &anonymous.Dependencies{
						PluginAPI: pluginMock,
						Store:     storeMock,
					},
				},
			},
			args: args{
				userID:    "in_general",
				channelID: "general",
				status:    false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := anonymous.New(tt.fields.Config)
			err := a.SetEncryptionStatusForChannel(tt.args.channelID, tt.args.userID, tt.args.status)
			test.CheckErr(is, tt.wantErr, err)
		})
	}
}

func Test_anonymous_IsEncryptionEnabledForChannel(t *testing.T) {

	ctrl := gomock.NewController(t)
	is := assert.New(t)

	storeMock := mockStore.NewMockStore(ctrl)

	storeMock.EXPECT().IsEncryptionEnabled("general", "in_general").Return(true).AnyTimes()
	storeMock.EXPECT().IsEncryptionEnabled(gomock.Any(), gomock.Any()).Return(false).AnyTimes()

	pluginMock := mockPlugin.NewMockPlugin(ctrl)
	pluginMock.EXPECT().GetConfiguration().Return(&config.Config{
		PluginID:      "",
		PluginVersion: "",
	}).AnyTimes()

	defer ctrl.Finish()

	type fields struct {
		Config anonymous.Config
	}
	type args struct {
		userID    string
		channelID string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		isEnabled bool
	}{
		{
			name: "enabled",
			fields: fields{
				Config: anonymous.Config{
					Dependencies: &anonymous.Dependencies{
						PluginAPI: pluginMock,
						Store:     storeMock,
					},
				},
			},
			args: args{
				userID:    "in_general",
				channelID: "general",
			},
			isEnabled: true,
		},
		{
			name: "disabled",
			fields: fields{
				Config: anonymous.Config{
					Dependencies: &anonymous.Dependencies{
						PluginAPI: pluginMock,
						Store:     storeMock,
					},
				},
			},
			args: args{
				userID:    "not_in_general",
				channelID: "general",
			},
			isEnabled: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := anonymous.New(tt.fields.Config)
			isEnabled := a.IsEncryptionEnabledForChannel(tt.args.channelID, tt.args.userID)
			is.Equal(tt.isEnabled, isEnabled)
		})
	}
}
