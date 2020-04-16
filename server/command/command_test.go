package command_test

import (
	"strings"
	"testing"

	"github.com/bakurits/mattermost-plugin-anonymous/server/command"

	"github.com/bakurits/mattermost-plugin-anonymous/server/utils/test"

	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	mockAnonymous "github.com/bakurits/mattermost-plugin-anonymous/server/anonymous/mock"
	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/assert"
)

func Test_command_overwrite(t *testing.T) {

	ctrl := gomock.NewController(t)
	tassert := assert.New(t)
	a := mockAnonymous.NewMockAnonymous(ctrl)
	a.EXPECT().StorePublicKey(gomock.Any()).Return(nil)
	a.EXPECT().SendEphemeralPost(gomock.Any(), gomock.Any()).Return(nil)
	b := mockAnonymous.NewMockAnonymous(ctrl)
	b.EXPECT().StorePublicKey(gomock.Any()).Return(&model.AppError{Message: "something went wrong while storing the key"})
	defer ctrl.Finish()

	type fields struct {
		anonymous anonymous.Anonymous
	}
	type args struct {
		commandArgs *model.CommandArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "overwrite without args test",
			fields: fields{
				anonymous: a,
			},
			args: args{
				commandArgs: &model.CommandArgs{
					Command: "/anonymous keypair --overwrite",
				},
			},
			wantErr: true,
		},
		{
			name: "overwrite too many args test",
			fields: fields{
				anonymous: a,
			},
			args: args{
				commandArgs: &model.CommandArgs{
					Command: "/anonymous keypair --overwrite 1 2 3 4 5",
				},
			},
			wantErr: true,
		},
		{
			name: "overwrite with args test",
			fields: fields{
				anonymous: a,
			},
			args: args{
				commandArgs: &model.CommandArgs{
					Command: "/anonymous keypair --overwrite 2JUMmg86ptVNYZT1scuLUIFegDr513LHA6aqsJCLaZo=",
				},
			},
			wantErr: false,
		},
		{
			name: "overwrite with invalid key",
			fields: fields{
				anonymous: a,
			},
			args: args{
				commandArgs: &model.CommandArgs{
					Command: "/anonymous keypair --overwrite 1234567891235412415251",
				},
			},
			wantErr: true,
		},
		{
			name: "overwrite error test",
			fields: fields{
				anonymous: b,
			},
			args: args{
				commandArgs: &model.CommandArgs{
					Command: "/anonymous keypair --overwrite 2JUMmg86ptVNYZT1scuLUIFegDr513LHA6aqsJCLaZo=",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comm := command.NewHandler(tt.args.commandArgs, tt.fields.anonymous)
			_, err := comm.Handle(strings.Fields(tt.args.commandArgs.Command)...)
			test.CheckErr(tassert, tt.wantErr, err)
		})
	}
}
func Test_command_other(t *testing.T) {

	ctrl := gomock.NewController(t)
	tassert := assert.New(t)
	a := mockAnonymous.NewMockAnonymous(ctrl)
	a.EXPECT().SendEphemeralPost(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(3)
	defer ctrl.Finish()

	type fields struct {
		anonymous anonymous.Anonymous
	}
	type args struct {
		commandArgs *model.CommandArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "help test",
			fields: fields{
				anonymous: a,
			},
			args: args{
				commandArgs: &model.CommandArgs{
					Command: "help",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong command test 1",
			fields: fields{
				anonymous: a,
			},
			args: args{
				commandArgs: &model.CommandArgs{
					Command: "wrong command",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong command test 2",
			fields: fields{
				anonymous: a,
			},
			args: args{
				commandArgs: &model.CommandArgs{
					Command: "/anonymous keypair --wrongAction",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comm := command.NewHandler(tt.args.commandArgs, tt.fields.anonymous)
			_, err := comm.Handle(strings.Fields(tt.args.commandArgs.Command)...)
			test.CheckErr(tassert, tt.wantErr, err)
		})
	}
}
