package bot

import (
	"github.com/bwmarrin/discordgo"
)

// Arguments provided to CommandFunctions
type CommandArguments struct {
	Message   *discordgo.Message
	Command   string
	Arguments []string
}

// Function to handle a particular command on Bot
type CommandFunction func(b *Bot, cmd *CommandArguments)

// Interface that provides a handlable command
type CommandProvider interface {
	// Provides and alias and a command
	ProvideCommand() (aliases []string, handler CommandFunction)
}
