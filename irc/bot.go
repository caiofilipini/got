package irc

import (
	"fmt"
	"log"
	"regexp"
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
	irc      IRC
	user     string
	passwd   string
	commands []Command
	request  chan string
	action   *regexp.Regexp
}

func (bot *Bot) Register(command Command) {
	bot.commands = append(bot.commands, command)
}

func (bot Bot) Recognise(request string) (Command, string, error) {
	for _, c := range bot.commands {
		if match := c.Pattern().FindStringSubmatch(request); len(match) > 0 {
			return c, match[len(match)-1], nil
		}
	}
	return nil, "", fmt.Errorf("Don't know how to handle \"%s\"", request)
}

func (bot Bot) Start() {
	go bot.irc.handleRead(bot)
	go bot.irc.handlePing()
	go bot.irc.handleWrite()

	bot.irc.out <- fmt.Sprintf("NICK %s", bot.user)
	bot.irc.out <- fmt.Sprintf("USER %s 0.0.0.0 0.0.0.0 :%s", bot.user, bot.user)
	bot.irc.out <- fmt.Sprintf("JOIN %s %s", bot.irc.channel, bot.passwd)

	bot.irc.Send(WelcomeMsg)
}

func (bot Bot) Listen() {
	for r := range bot.request {
		info(fmt.Sprintf("Received request: %s", r))

		if command, query, err := bot.Recognise(r); err == nil {
			messages := command.Run(query)
			bot.irc.Send(messages...)
		} else {
			info(fmt.Sprintf("WARNING: %s", err.Error()))
		}
	}
}

func (bot Bot) ActionRequested(msg string) bool {
	return bot.action.Match([]byte(msg))
}

func (bot Bot) Handle(msg string) {
	if req := bot.action.FindStringSubmatch(msg); len(req) > 1 {
		bot.request <- req[1]
	}
}

func NewBot(irc IRC, user, passwd string) Bot {
	return Bot{
		irc,
		user,
		passwd,
		make([]Command, 0),
		make(chan string),
		regexp.MustCompile(fmt.Sprintf("PRIVMSG %s :%s (.*)", irc.channel, Action)),
	}
}

func info(msg string) {
	log.Printf("[Bot] %s\n", msg)
}
