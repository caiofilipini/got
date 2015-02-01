package bot

import (
	"fmt"
	"log"
	"regexp"

	"github.com/caiofilipini/got/irc"
)

const (
	Action      = "!got"
	WelcomeMsg  = "OHAI"
	HelpCommand = `(?i)help\s*(.*)`
)

type Command interface {
	Run(string) []string
	Name() string
	Pattern() *regexp.Regexp
	Help() string
	Usage() []string
}

type Bot struct {
	irc            irc.IRC
	user           string
	passwd         string
	commands       []Command
	commandsByName map[string]Command
	action         *regexp.Regexp
	helpPattern    *regexp.Regexp
	request        chan string
	in             chan string
}

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

func (bot *Bot) Register(command Command) {
	bot.commands = append(bot.commands, command)
	bot.commandsByName[command.Name()] = command
}

func (bot Bot) Start() {
	bot.irc.Subscribe(bot.action, bot.in)
	bot.irc.Join(bot.user, bot.passwd)
	bot.irc.SendMessages(WelcomeMsg)
}

func (bot Bot) Listen() {
	go bot.handleRequests()

	for msg := range bot.in {
		if req := bot.action.FindStringSubmatch(msg); len(req) > 1 {
			bot.request <- req[1]
		}
	}
}

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
