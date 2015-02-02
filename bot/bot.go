// Package bot provides the implementation of the IRC bot.
package bot

import (
	"fmt"
	"log"
	"regexp"

	"github.com/caiofilipini/got/irc"
)

const (
	// Action is the command that triggers the bot.
	Action = "!got"

	// WelcomeMsg is the message to be printed out when the bot is online.
	WelcomeMsg = "OHAI"

	// HelpCommand is the pattern for the help command.
	HelpCommand = `(?i)help\s*(.*)`
)

// Command is the interface that registered commands need to
// implement in order to receive messages.
type Command interface {
	// Name returns the command name.
	Name() string

	// Pattern returns the pattern to be matched against
	// in order to check if this command should be triggered.
	Pattern() *regexp.Regexp

	// Help returns the help message for this command.
	Help() string

	// Usage returns details about how to use this command.
	Usage() []string

	// Run receives a query and returns a list of messages
	// to be sent in response.
	Run(string) []string
}

// Bot represents a running instance of the bot.
type Bot struct {
	// The IRC connection.
	irc irc.IRC

	// The user representing the bot in the IRC channel.
	user string

	// The password for the IRC channel (if applicable).
	passwd string

	// The list of registered commands.
	commands []Command

	// A map where the key is the command name, and the
	// value is the command itself.
	commandsByName map[string]Command

	// The regexp pattern that matches the action to trigger
	// the bot.
	action *regexp.Regexp

	// The regexp pattern that matches the help command.
	helpPattern *regexp.Regexp

	// The channel where filtered requests are sent.
	request chan string

	// The channel where messages that match the configured
	// action are sent.
	in chan string
}

// NewBot creates and return a value representing
// a connected bot.
func NewBot(irc irc.IRC, user, passwd string) Bot {
	return Bot{
		irc:            irc,
		user:           user,
		passwd:         passwd,
		commands:       make([]Command, 0),
		commandsByName: make(map[string]Command),
		action:         regexp.MustCompile(fmt.Sprintf("PRIVMSG %s :%s (.*)", irc.Channel, Action)),
		helpPattern:    regexp.MustCompile(HelpCommand),
		request:        make(chan string),
		in:             make(chan string),
	}
}

// Register registers the given command.
func (bot *Bot) Register(command Command) {
	bot.commands = append(bot.commands, command)
	bot.commandsByName[command.Name()] = command
}

// Start joins the channel, sends a welcome message and
// subscribes to messages that match the configured action.
func (bot Bot) Start() {
	bot.irc.Subscribe(bot.action, bot.in)
	bot.irc.Join(bot.user, bot.passwd)
	bot.irc.SendMessages(WelcomeMsg)
}

// Listen starts a background process to listen to
// incoming requests.
func (bot Bot) Listen() {
	go bot.handleRequests()

	for msg := range bot.in {
		if req := bot.action.FindStringSubmatch(msg); len(req) > 1 {
			bot.request <- req[1]
		}
	}
}

// Shutdown closes the incoming request channels.
func (bot Bot) Shutdown() {
	close(bot.in)
	close(bot.request)
}

func (bot Bot) recognise(request string) (Command, string, error) {
	for _, c := range bot.commands {
		if match := c.Pattern().FindStringSubmatch(request); len(match) > 0 {
			return c, match[len(match)-1], nil
		}
	}
	return nil, "", fmt.Errorf("Don't know how to handle \"%s\"", request)
}

func (bot Bot) showHelp(command string) {
	var helpMessages []string

	if command != "" {
		if c, found := bot.commandsByName[command]; found {
			helpMessages = append(helpMessages, formatHelp(c.Usage()...)...)
		} else {
			helpMessages = append(helpMessages, "unknown command: "+command)
		}
	} else {
		for _, c := range bot.commands {
			helpMessages = append(helpMessages, formatHelp(c.Help())...)
		}
		helpHelp := []string{
			"help – displays this message",
			"help <command> – displays usage for the given command",
		}
		helpMessages = append(helpMessages, formatHelp(helpHelp...)...)
	}
	bot.irc.SendMessages(helpMessages...)
}

func (bot Bot) handleRequests() {
	for r := range bot.request {
		info(fmt.Sprintf("Received request: %s", r))

		if match := bot.helpPattern.FindStringSubmatch(r); len(match) > 0 {
			command := match[len(match)-1]
			bot.showHelp(command)
		} else if command, query, err := bot.recognise(r); err == nil {
			messages := command.Run(query)
			bot.irc.SendMessages(messages...)
		} else {
			info(fmt.Sprintf("WARNING: %s", err.Error()))
		}
	}
}

func formatHelp(messages ...string) []string {
	formatted := make([]string, len(messages))
	for i, m := range messages {
		formatted[i] = fmt.Sprintf("%s %s", Action, m)
	}
	return formatted
}

func info(msg string) {
	log.Printf("[Bot] %s\n", msg)
}
