package bot

import (
	"github.com/bwmarrin/discordgo"
)

// Arguments provided to CommandFunctions
type CommandArguments struct {
	// Discord Message object corresponding to the command
	Message *discordgo.Message
	// CommandName alias used to call the command
	CommandName string
	// Slice of arguments passed in space separated to the command
	Arguments []string
}

// Interface that provides a handlable command
type Command interface {
	// Provides a list of aliases and a command
	ProvideAliases() (aliases []string)
	ProvideHelpMessage() (helpMessage string)
	Handle(session *discordgo.Session, cmd CommandArguments) error
}
