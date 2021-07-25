package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var commandSyntax = regexp.MustCompile("^\\s*!([A-Za-z]+)((?: +[^ ]+)+)?\\s*$")
var repeatedSpaces = regexp.MustCompile(" +")

// Simple discord bot
// Provided command handlers, the bot will pass commands of a matching alias to the handler functions
type Bot struct {
	Name    string
	HelpMsg *string
	Session *discordgo.Session

	commandProviders    []CommandProvider
	customEventHandlers []interface{}
	commands            map[string]CommandFunction

	logEnabled bool
}

func New(name string, helpMsg *string, logEnabled bool) Bot {
	bot := Bot{
		Name:       name,
		HelpMsg:    helpMsg,
		commands:   make(map[string]CommandFunction),
		logEnabled: logEnabled,
	}
	return bot
}

func (b *Bot) Init(token string) error {
	var err error
	b.Session, err = discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	b.initCommands()
	b.initEventHandlers()

	err = b.Session.Open()
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) initCommands() {
	for _, provider := range b.commandProviders {
		aliases, handler := provider.ProvideCommand()
		for _, alias := range aliases {
			b.commands[alias] = handler
		}
	}
}

func (b *Bot) initEventHandlers() {
	b.Session.AddHandlerOnce(b.handleReady)
	b.Session.AddHandler(b.handleMessage)

	for _, handler := range b.customEventHandlers {
		b.Session.AddHandler(handler)
	}
}

func (b *Bot) AddCommand(provider CommandProvider) {
	b.commandProviders = append(b.commandProviders, provider)
}

func (b *Bot) AddEventHandler(handler interface{}) {
	b.customEventHandlers = append(b.customEventHandlers, handler)
}

func (b *Bot) handleReady(s *discordgo.Session, e *discordgo.Ready) {
	fmt.Println("Discord ready")
}

func (b *Bot) handleMessage(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.Bot {
		return
	}

	cmdArgs := b.parseCommand(e.Message)
	if cmdArgs != nil {
		b.commands[cmdArgs.Command](b, cmdArgs)
	}
}

func (b *Bot) parseCommand(message *discordgo.Message) *CommandArguments {
	m := commandSyntax.FindStringSubmatch(message.Content)
	if m == nil {
		return nil
	}
	if _, exists := b.commands[m[1]]; !exists {
		return nil
	}

	var args []string
	if m[2] != "" {
		argString := m[2]
		argString = repeatedSpaces.ReplaceAllLiteralString(argString, " ")
		argString = strings.TrimSpace(argString)
		args = strings.Split(argString, " ")
	}

	return &CommandArguments{
		Message:   message,
		Command:   m[1],
		Arguments: args,
	}
}
