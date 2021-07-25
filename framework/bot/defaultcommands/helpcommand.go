package bot

import (
	"strings"

	"github.com/tr4656/weeb-bot-go/framework/bot"
)

// Command to sent a help message
type HelpCommand struct {
	// Help message string to send on command
	helpMessage string
}

func (h *HelpCommand) ProvideCommand() ([]string, bot.CommandFunction) {
	return []string{"h", "help"}, h.handleHelp
}

func (h *HelpCommand) handleHelp(b *bot.Bot, cmd *bot.CommandArguments) {
	// Ensure the sole argument to the command is the bot name
	if len(cmd.Arguments) == 0 ||
		!strings.EqualFold(cmd.Arguments[0], b.Name) {
		return
	}

	b.Session.ChannelMessageSend(cmd.Message.ChannelID, h.helpMessage)
}
