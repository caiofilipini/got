package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/caiofilipini/got/bot"
	"github.com/caiofilipini/got/command"
	"github.com/caiofilipini/got/irc"
)

var (
	server      *string
	port        *int
	channel     *string
	user        *string
	passwd      *string
	logFilePath *string
)

func init() {
	server = flag.String("s", "irc.freenode.org", "IRC server host")
	port = flag.Int("p", 6667, "IRC server port")
	user = flag.String("u", "gotgotgot", "bot username")
	channel = flag.String("c", "", "channel to connect")
	passwd = flag.String("k", "", "channel secret key")
	logFilePath = flag.String("l", "", "log file location; if empty, stdout will be used")

	flag.Parse()

	if *channel == "" {
		log.Println("No channel specified, aborting!")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func setupLogging() *os.File {
	if *logFilePath != "" {
		file, err := os.OpenFile(*logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)
		return file
	}

	return nil
}

func main() {
	logFile := setupLogging()
	if logFile != nil {
		defer logFile.Close()
	}

	conn := irc.NewIRC(*server, *port, *channel)
	defer conn.Close()

	bot := bot.NewBot(conn, *user, *passwd)
	defer bot.Shutdown()

	// Register commands
	bot.Register(command.Swear())
	bot.Register(command.Greet())
	bot.Register(command.Image())
	bot.Register(command.GIF())
	bot.Register(command.Video())
	bot.Register(command.XKCD())
	bot.Register(command.BeerOClock())
	bot.Register(command.Weather())
	bot.Register(command.Luca()) // tribute to lucapette

	bot.Start()
	go bot.Listen()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for _ = range signals {
		log.Println("KTHXBAI.")
		os.Exit(0)
	}
}
