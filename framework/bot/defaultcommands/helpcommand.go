package defaultcommands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tr4656/weeb-bot-go/framework/bot"
)

// Command to sent a help message
type HelpCommand struct {
	// Name of the owning bot
	BotName string
	// Help message string to send on command
	HelpMessage string
}

func (h *HelpCommand) ProvideAliases() []string {
	return []string{"h", "help"}
}

// HelpCommand does not have a help message itself
func (h *HelpCommand) ProvideHelpMessage() string {
	return ""
}

func (h *HelpCommand) Handle(session *discordgo.Session, cmd bot.CommandArguments) {
	// Ensure the sole argument to the command is the bot name
	if len(cmd.Arguments) == 0 ||
		!strings.EqualFold(cmd.Arguments[0], h.BotName) {
		return
	}

	session.ChannelMessageSend(cmd.Message.ChannelID, h.HelpMessage)
}
