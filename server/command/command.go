package command

import (
	"fmt"
	"strings"

	"github.com/bakurits/mattermost-plugin-anonymous/server/anonymous"
	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	helpTextHeader = "###### Mattermost Anonymous Plugin - Slash command Help\n"
	helpText       = `
* |/anonymous help| - print this help message
* |/anonymous keypair [action]| - do one of the following actions regarding encryption keypair
  * |action| is one of the following:
    * |--generate| - generates and stores new keypair for encryption
	* |--overwrite [private key]| - you enter new 32byte private key, the plugin stores it along with the updated public key
    * |--export| - exports your existing keypair
`
)

// Command returns API for interacting with plugin commands
type Command interface {
	Handle(args ...string) (*model.CommandResponse, error)
	Help(args ...string) (*model.CommandResponse, error)
}

// command stores command specific information
type command struct {
	args      *model.CommandArgs
	anonymous anonymous.Anonymous
	handler   HandlerMap
}

// HandlerFunc command handler function type
type HandlerFunc func(args ...string) (*model.CommandResponse, error)

// HandlerMap map of command handler functions
type HandlerMap struct {
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

func newCommand(args *model.CommandArgs, a anonymous.Anonymous) *command {
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

// New returns new Command with given dependencies
func New(args *model.CommandArgs, a anonymous.Anonymous) Command {
	return newCommand(args, a)
}

func (c *command) Handle(args ...string) (*model.CommandResponse, error) {
	ch := c.handler
	for n := len(args); n > 0; n-- {
		h := ch.handlers[strings.Join(args[:n], "/")]
		if h != nil {
			return h(args[n:]...)
		}
	}
	return ch.defaultHandler(args...)
}

func (c *command) executeKeyPairGenerate(args ...string) (*model.CommandResponse, error) {
	return &model.CommandResponse{}, nil
}

func (c *command) executeKeyOverwrite(args ...string) (*model.CommandResponse, error) {
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

func (c *command) executeKeyExport(args ...string) (*model.CommandResponse, error) {
	return &model.CommandResponse{}, nil
}

func (c *command) Help(args ...string) (*model.CommandResponse, error) {
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

// GetSlashCommand returns command to register
func GetSlashCommand() *model.Command {
	return &model.Command{
		Trigger:          "anonymous",
		DisplayName:      "anonymous",
		Description:      "End to end message encryption",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: keypair [--generate, --export, --overwrite]",
		AutoCompleteHint: "[command][subcommands]",
	}
}
