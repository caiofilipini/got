package bot

import (
	"fmt"
	"log"
	"regexp"

	"github.com/caiofilipini/got/irc"
)

const (
	Action     = "!got"
	WelcomeMsg = "OHAI"
)

type Command interface {
	Run(string) []string
	Pattern() *regexp.Regexp
}

type Bot struct {
	irc      irc.IRC
	user     string
	passwd   string
	commands []Command
	action   *regexp.Regexp
	request  chan string
	in       chan string
}

func NewBot(irc irc.IRC, user, passwd string) Bot {
	return Bot{
		irc:      irc,
		user:     user,
		passwd:   passwd,
		commands: make([]Command, 0),
		action:   regexp.MustCompile(fmt.Sprintf("PRIVMSG %s :%s (.*)", irc.Channel, Action)),
		request:  make(chan string),
		in:       make(chan string),
	}
}

func (bot *Bot) Register(command Command) {
	bot.commands = append(bot.commands, command)
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

func (bot Bot) recognise(request string) (Command, string, error) {
	for _, c := range bot.commands {
		if match := c.Pattern().FindStringSubmatch(request); len(match) > 0 {
			return c, match[len(match)-1], nil
		}
	}
	return nil, "", fmt.Errorf("Don't know how to handle \"%s\"", request)
}

func (bot Bot) handleRequests() {
	for r := range bot.request {
		info(fmt.Sprintf("Received request: %s", r))

		if command, query, err := bot.recognise(r); err == nil {
			messages := command.Run(query)
			bot.irc.SendMessages(messages...)
		} else {
			info(fmt.Sprintf("WARNING: %s", err.Error()))
		}
	}
}

func info(msg string) {
	log.Printf("[Bot] %s\n", msg)
}
