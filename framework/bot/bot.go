package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var commandSyntax = regexp.MustCompile(`^\s*!([A-Za-z]+)((?: +[^ ]+)+)?\s*$`)
var repeatedSpaces = regexp.MustCompile(" +")

// Simple discord bot
// Provided command handlers, the bot will pass commands of a matching alias to the handler functions
type Bot struct {
	// Bot name - used for logging and for commands to recognise reference to the bot itself
	Name string
	// Discord Session instance - null until Init() is called
	Session *discordgo.Session

	// Slice of all CommandProvider interfaces to use on Init()
	commandProviders []CommandProvider
	// Slice of all Discord event functions to add as handlers to the Session on Init()
	customEventHandlers []interface{}
	// Map of all aliases to their corresponding handler functions
	commands map[string]CommandFunction

	// Whether to log error messages to discord or not
	// TODO: Implement usage
	logEnabled bool
}

// Create a new Bot instance
func New(name string, logEnabled bool) Bot {
	bot := Bot{
		Name:       name,
		commands:   make(map[string]CommandFunction),
		logEnabled: logEnabled,
	}
	return bot
}

// Initialise the Bot instance
// Loads in all commands and adds all event handlers, then opens the session
// Returns errors from discordgo
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

// Initialise all commands
// Gets all commands from assigned providers and generates a map lookup to handler functions
func (b *Bot) initCommands() {
	for _, provider := range b.commandProviders {
		aliases, handler := provider.ProvideCommand()
		for _, alias := range aliases {
			b.commands[strings.ToLower(alias)] = handler
		}
	}
}

// Initialise all event handlers
// Add default handlers for Ready and CreateMessage, then all assigned custom handlers
// Will print a message if an invalid handler is present in in customEventHandlers
func (b *Bot) initEventHandlers() {
	b.Session.AddHandlerOnce(b.handleReady)
	b.Session.AddHandler(b.handleMessage)

	for _, handler := range b.customEventHandlers {
		b.Session.AddHandler(handler)
	}
}

// Add a command in the form of a CommandProvider
func (b *Bot) AddCommand(provider CommandProvider) {
	b.commandProviders = append(b.commandProviders, provider)
}

// Add an event handler in the form of a Discord event handling function
func (b *Bot) AddEventHandler(handler interface{}) {
	b.customEventHandlers = append(b.customEventHandlers, handler)
}

// Handle the Ready even from Discord
func (b *Bot) handleReady(s *discordgo.Session, e *discordgo.Ready) {
	fmt.Println("Discord ready")
}

// Handle the Message even from Discord
func (b *Bot) handleMessage(s *discordgo.Session, e *discordgo.MessageCreate) {
	// Ignore bot messages
	if e.Author.Bot {
		return
	}

	// If successfully parsed, run appropriate handler function with parsed args
	cmdArgs := b.parseCommand(e.Message)
	if cmdArgs != nil {
		b.commands[cmdArgs.Command](b, cmdArgs)
	}
}

// Parse arguments in the provided message
// Returns arguments if successfully parsed, otherwise returns nil
func (b *Bot) parseCommand(message *discordgo.Message) *CommandArguments {
	// Ensure message matches commnd syntax
	m := commandSyntax.FindStringSubmatch(message.Content)
	if m == nil {
		return nil
	}
	// Ensure command alias in message is present in commands lookup
	if _, exists := b.commands[strings.ToLower(m[1])]; !exists {
		return nil
	}

	// If arguments are present, clean up spaces and split arguments into a slice
	var args []string
	if m[2] != "" {
		argString := m[2]
		argString = repeatedSpaces.ReplaceAllLiteralString(argString, " ")
		argString = strings.TrimSpace(argString)
		args = strings.Split(argString, " ")
	}

	return &CommandArguments{
		Message:   message,
		Command:   strings.ToLower(m[1]), // Command alias is returned as lowercase string for consistency
		Arguments: args,
	}
}
