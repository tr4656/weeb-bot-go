package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commandSyntax  = regexp.MustCompile(`^\s*!([A-Za-z]+)((?: +[^ ]+)+)?\s*$`)
	repeatedSpaces = regexp.MustCompile(" +")
)

// Simple discord bot
// Provided command handlers, the bot will pass commands of a matching alias to the handler functions
type Bot struct {
	// Bot name - used for logging and for commands to recognise reference to the bot itself
	Name string
	// Help message on the bot level
	botHelpMessage string
	// Discord Session instance - null until Init() is called
	session *discordgo.Session

	// Slice of all Discord event functions to add as handlers to the Session on Init()
	customEventHandlers []interface{}
	// Map of all aliases to their corresponding handlers
	commands map[string]Command

	// Whether to log error messages to discord or not
	// TODO: Implement usage
	logEnabled bool
}

// Create a new Bot instance
func New(name string, botHelpMessage string, logEnabled bool) *Bot {
	bot := Bot{
		Name:       name,
		commands:   make(map[string]Command),
		logEnabled: logEnabled,
	}
	return &bot
}

// Initialise the Bot instance
// Loads in all commands and adds all event handlers, then opens the session
// Returns errors from discordgo
func (b *Bot) Init(token string) error {
	var err error
	b.session, err = discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	b.initEventHandlers()

	err = b.session.Open()
	if err != nil {
		return err
	}

	return nil
}

// Initialise all event handlers
// Add default handlers for Ready and CreateMessage, then all assigned custom handlers
// Will print a message if an invalid handler is present in in customEventHandlers
func (b *Bot) initEventHandlers() {
	b.session.AddHandlerOnce(b.handleReady)
	b.session.AddHandler(b.handleMessage)

	for _, handler := range b.customEventHandlers {
		b.session.AddHandler(handler)
	}
}

// Generates a help message from commands loaded in
func (b *Bot) GenerateHelpMessage() string {
	helpMessage := b.botHelpMessage

	helpMessage += "\n\n"

	for _, command := range b.commands {
		helpMessage += command.ProvideHelpMessage() + "\n"
	}

	return helpMessage
}

// Add a command to the commands map
func (b *Bot) AddCommand(command Command) {
	aliases := command.ProvideAliases()
	for _, alias := range aliases {
		b.commands[strings.ToLower(alias)] = command
	}
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

	// If successfully parsed, run appropriate handler with parsed args
	cmdArgs := b.parseCommand(e.Message)
	if cmdArgs.CommandName != "" {
		handler := b.commands[cmdArgs.CommandName]
		err := handler.Handle(b.session, cmdArgs)

		if err != nil {
			// TODO: Log better
			fmt.Printf("Error handling command '%s': %s\n", cmdArgs.CommandName, err)
		}
	}
}

// Parse arguments in the provided message
// Returns populated arguments if successfully parsed, otherwise returns empty arguments
func (b *Bot) parseCommand(message *discordgo.Message) CommandArguments {
	// Ensure message matches commnd syntax
	m := commandSyntax.FindStringSubmatch(message.Content)
	if m == nil {
		return CommandArguments{}
	}
	// Ensure command alias in message is present in commands lookup
	if _, exists := b.commands[strings.ToLower(m[1])]; !exists {
		return CommandArguments{}
	}

	// If arguments are present, clean up spaces and split arguments into a slice
	var args []string
	if m[2] != "" {
		argString := m[2]
		argString = repeatedSpaces.ReplaceAllLiteralString(argString, " ")
		argString = strings.TrimSpace(argString)
		args = strings.Split(argString, " ")
	}

	return CommandArguments{
		Message:     message,
		CommandName: strings.ToLower(m[1]), // Command alias is returned as lowercase string for consistency
		Arguments:   args,
	}
}
