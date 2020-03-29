package anonymous

import (
	"errors"
	"fmt"
	"github.com/bakurits/mattermost-plugin-anonymous/server/config"
	mockPlugin "github.com/bakurits/mattermost-plugin-anonymous/server/plugin/mock"
	"github.com/bakurits/mattermost-plugin-anonymous/server/store"
	mockStore "github.com/bakurits/mattermost-plugin-anonymous/server/store/mock"
	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/test"
	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
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
	return fmt.Sprintf("should match with strings containging (%s)", string(s))
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
	return fmt.Sprintf("should match with strings containging ()")
}

func Test_anonymous_GetPublicKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	tassert := assert.New(t)

	storeMock := mockStore.NewMockStore(ctrl)

	storeMock.EXPECT().LoadUser(stringLikeMatcher("user_in")).Return(&store.User{
		MattermostUserID: "user_in",
		PublicKey:        []byte{1, 1, 1},
	}, nil)

	storeMock.EXPECT().LoadUser(gomock.Not(stringLikeMatcher("user_in"))).Return(nil, errors.New("some error"))

	pluginMock := mockPlugin.NewMockPlugin(ctrl)
	defer ctrl.Finish()

	type fields struct {
		Config                 Config
		actingMattermostUserID string
		PluginContext          plugin.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "basic test",
			fields: fields{
				Config: Config{
					Config: &config.Config{},
					Dependencies: &Dependencies{
						PluginAPI: pluginMock,
						Store:     storeMock,
					},
				},
				actingMattermostUserID: "user_in",
				PluginContext:          plugin.Context{},
			},
			want:    []byte{1, 1, 1},
			wantErr: false,
		},
		{
			name: "test empty",
			fields: fields{
				Config: Config{
					Config: &config.Config{},
					Dependencies: &Dependencies{
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
			a := New(tt.fields.Config, tt.fields.actingMattermostUserID, tt.fields.PluginContext)
			got, err := a.GetPublicKey()
			test.CheckErr(tassert, tt.wantErr, err)
			tassert.Equal(got, tt.want)
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
	defer ctrl.Finish()

	type fields struct {
		Config                 Config
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
				Config: Config{
					Config: &config.Config{},
					Dependencies: &Dependencies{
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
				Config: Config{
					Config: &config.Config{},
					Dependencies: &Dependencies{
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
			a := New(tt.fields.Config, tt.fields.actingMattermostUserID, tt.fields.PluginContext)
			err := a.StorePublicKey(tt.args.publicKey)
			test.CheckErr(tassert, tt.wantErr, err)
		})
	}
}
