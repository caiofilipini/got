package irc

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

const (
	Action     = "!got"
	WelcomeMsg = "OHAI"
)

type Handler func(...string)

type Bot struct {
	irc      IRC
	user     string
	passwd   string
	commands map[string]Handler
	request  chan string
	action   *regexp.Regexp
}

func (bot Bot) Register(command string, handler Handler) {
	bot.commands[command] = handler
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
		log.Printf("Received request: %s\n", r)

		parts := strings.Fields(r)
		command, args := parts[0], parts[1:]
		if handler, registered := bot.commands[command]; registered {
			handler(args...)
		} else {
			log.Printf("WARNING: Unknown command \"%s\".\n", r)
		}
	}
}

func (bot Bot) Swear(args ...string) {
	bot.irc.Send("ANNAGG A MARONN!")
}

func (bot Bot) ActionRequested(msg string) bool {
	return bot.action.Match([]byte(msg))
}

func (bot Bot) Handle(msg string) {
	req := bot.action.FindStringSubmatch(msg)
	if len(req) > 1 {
		bot.request <- req[1]
	}
}

func NewBot(irc IRC, user, passwd string) Bot {
	return Bot{
		irc,
		user,
		passwd,
		make(map[string]Handler),
		make(chan string),
		regexp.MustCompile(fmt.Sprintf("PRIVMSG %s :%s (.*)", irc.channel, Action)),
	}
}
