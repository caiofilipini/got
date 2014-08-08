package command

import (
	"fmt"

	"github.com/caiofilipini/got/irc"
)

func Greet(bot irc.Bot, query string) {
	bot.Send(fmt.Sprintf("ohai there, %s!", query))
}
