package bot

import (
	"github.com/bwmarrin/discordgo"
)

// Arguments provided to CommandFunctions
type CommandArguments struct {
	// Discord Message object corresponding to the command
	Message *discordgo.Message
	// Command alias used to call the command
	Command string
	// Slice of arguments passed in space separated to the command
	Arguments []string
}

// Function to handle a particular command on Bot
type CommandFunction func(b *Bot, cmd *CommandArguments)

// Interface that provides a handlable command
type CommandProvider interface {
	// Provides a list of aliases and a command
	ProvideCommand() (aliases []string, handler CommandFunction)
}
