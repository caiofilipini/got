package irc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type IRC struct {
	server  string
	port    int
	channel string
	conn    net.Conn
	ping    chan string
	out     chan string
}

func (irc IRC) handleRead(bot Bot) {
	buf := bufio.NewReaderSize(irc.conn, 512)

	for {
		msg, err := buf.ReadString('\n')
		if err != nil {
			panic(err)
		}

		msg = msg[:len(msg)-2]
		if strings.Contains(msg, "PING") {
			irc.ping <- msg
		} else if bot.ActionRequested(msg) {
			bot.Handle(msg)
		}
	}
}

func (irc IRC) handleWrite() {
	for msg := range irc.out {
		irc.send(msg)
	}
}

func (irc IRC) handlePing() {
	for ping := range irc.ping {
		server := strings.Split(ping, ":")[1]

		irc.out <- fmt.Sprintf("PONG %s", server)
		log.Printf("[IRC] PONG sent to %s\n", server)
	}
}

func (irc IRC) send(msg string) {
	_, err := irc.conn.Write([]byte(fmt.Sprintf("%s\r\n", msg)))
	if err != nil {
		panic(err)
	}
}

func (irc IRC) Close() {
	irc.conn.Close()
}

func (irc IRC) Send(messages ...string) {
	for _, msg := range messages {
		irc.out <- fmt.Sprintf("PRIVMSG %s :%s", irc.channel, msg)
	}
}

func NewIRC(server string, port int, channel string) IRC {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		panic(err)
	}
	log.Printf("[IRC] Connected to %s (%s).\n", server, conn.RemoteAddr())

	return IRC{
		server,
		port,
		channel,
		conn,
		make(chan string),
		make(chan string),
	}
}
