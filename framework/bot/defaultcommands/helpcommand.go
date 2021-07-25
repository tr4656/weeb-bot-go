package bot

import (
	"strings"

	"github.com/tr4656/weeb-bot-go/framework/bot"
)

type HelpCommand struct{}

func (h *HelpCommand) ProvideCommand() ([]string, bot.CommandFunction) {
	return []string{"h", "help"}, handleHelp
}

func handleHelp(b *bot.Bot, cmd *bot.CommandArguments) {
	if len(cmd.Arguments) == 0 ||
		strings.EqualFold(cmd.Arguments[0], b.Name) {
		return
	}

	b.Session.ChannelMessageSend(cmd.Message.ChannelID, *b.HelpMsg)
}
