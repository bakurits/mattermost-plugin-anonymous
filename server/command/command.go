package command

import (
	"fmt"
	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	"github.com/mattermost/mattermost-server/v5/model"
	"strings"
)

const (
	helpTextHeader = "###### Mattermost Anonymous Plugin - Slash command Help\n"
	helpText       = `
* |/Anonymous help| - print this help message
* |/Anonymous keypair [action]| - do one of the following actions regarding encryption keypair
  * |action| is one of the following:
    * |--generate| - generates and stores new keypair for encryption
	* |--overwrite [private key]| - you enter new 32byte private key, the plugin stores it along with the updated public key
    * |--export| - exports your existing keypair
`
)

type Command interface {
	Handle(args ...string) (*model.CommandResponse, *model.AppError)
	Help(args ...string) (*model.CommandResponse, *model.AppError)
	executeKeyPairGenerate(args ...string) (*model.CommandResponse, *model.AppError)
	executeKeyOverwrite(args ...string) (*model.CommandResponse, *model.AppError)
	executeKeyExport(args ...string) (*model.CommandResponse, *model.AppError)
	postCommandResponse(text string)
	responsef(format string, args ...interface{}) *model.CommandResponse
	responseRedirect(redirectURL string) *model.CommandResponse
}

// command stores command specific information
type command struct {
	args      *model.CommandArgs
	anonymous anonymous.Anonymous
	handler   HandlerMap
}

// HandlerFunc command handler function type
type HandlerFunc func(args ...string) (*model.CommandResponse, *model.AppError)

// HandlerMap map of command handler functions
type HandlerMap struct {
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

func New(args *model.CommandArgs, a anonymous.Anonymous) Command {
	c := &command{
		args:      args,
		anonymous: a,
	}

	c.handler = HandlerMap{
		handlers: map[string]HandlerFunc{
			"help":                c.Help,
			"keypair/--generate":  c.executeKeyPairGenerate,
			"keypair/--overwrite": c.executeKeyOverwrite,
			"keypair/--export":    c.executeKeyExport,
		},
		defaultHandler: c.Help,
	}
	return c
}

func (c *command) Handle(args ...string) (*model.CommandResponse, *model.AppError) {
	ch := c.handler
	for n := len(args); n > 0; n-- {
		h := ch.handlers[strings.Join(args[:n], "/")]
		if h != nil {
			return h(args[n:]...)
		}
	}
	return ch.defaultHandler(args...)
}

func (c *command) executeKeyPairGenerate(args ...string) (*model.CommandResponse, *model.AppError) {
	return &model.CommandResponse{}, nil
}

func (c *command) executeKeyOverwrite(args ...string) (*model.CommandResponse, *model.AppError) {
	if len(args) > 1 {
		return &model.CommandResponse{}, &model.AppError{
			Message: "Too many arguments passed to e",
		}
	}
	if len(args) == 0 {
		return &model.CommandResponse{}, &model.AppError{
			Message: "public key not passed to function",
		}
	}
	pub := args[0]
	err := c.anonymous.StorePublicKey([]byte(pub))
	if err != nil {
		return &model.CommandResponse{}, &model.AppError{
			Message: err.Error(),
		}
	}
	return &model.CommandResponse{}, nil
}

func (c *command) executeKeyExport(args ...string) (*model.CommandResponse, *model.AppError) {
	return &model.CommandResponse{}, nil
}

func (c *command) Help(args ...string) (*model.CommandResponse, *model.AppError) {
	helpText := helpTextHeader + helpText
	c.postCommandResponse(helpText)
	return &model.CommandResponse{}, nil
}

func (c *command) postCommandResponse(text string) {
	post := &model.Post{
		ChannelId: c.args.ChannelId,
		Message:   text,
	}
	_ = c.anonymous.SendEphemeralPost(c.args.UserId, post)
}

func (c *command) responsef(format string, args ...interface{}) *model.CommandResponse {
	c.postCommandResponse(fmt.Sprintf(format, args...))
	return &model.CommandResponse{}
}

func (c *command) responseRedirect(redirectURL string) *model.CommandResponse {
	return &model.CommandResponse{
		GotoLocation: redirectURL,
	}
}

func GetSlashCommand() *model.Command {
	return &model.Command{
		Trigger:          "Anonymous",
		DisplayName:      "Anonymous",
		Description:      "End to end message encryption",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: keypair [--generate, --export, --overwrite]",
		AutoCompleteHint: "[command][subcommands]",
	}
}
