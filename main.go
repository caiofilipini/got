package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/caiofilipini/got/irc"
)

var (
	server  *string
	port    *int
	channel *string
	user    *string
	passwd  *string
)

func init() {
	server = flag.String("s", "irc.freenode.org", "IRC server host")
	port = flag.Int("p", 6667, "IRC server port")
	user = flag.String("u", "gotgotgot", "bot username")
	channel = flag.String("c", "", "channel to connect")
	passwd = flag.String("k", "", "channel secret key")

	flag.Parse()

	if *channel == "" {
		log.Println("No channel specified, aborting!")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	conn := irc.NewIRC(*server, *port, *channel)
	defer conn.Close()

	bot := irc.NewBot(conn, *user, *passwd)

	// Register commands
	bot.Register("swear", bot.Swear)

	bot.Start()
	go bot.Listen()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for _ = range signals {
		log.Println("KTHXBAI.")
		os.Exit(0)
	}
}
