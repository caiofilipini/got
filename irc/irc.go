package irc

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

type IRC struct {
	server        string
	port          int
	Channel       string
	conn          net.Conn
	ping          chan string
	out           chan string
	subscriptions map[*regexp.Regexp]chan string
}

func NewIRC(server string, port int, channel string) IRC {
	conn := connect(server, port)

	irc := IRC{
		server:        server,
		port:          port,
		Channel:       channel,
		conn:          conn,
		ping:          make(chan string),
		out:           make(chan string),
		subscriptions: make(map[*regexp.Regexp]chan string),
	}

	go irc.handleRead()
	go irc.handlePing()
	go irc.handleWrite()

	return irc
}

func (irc IRC) Close() {
	irc.conn.Close()

	close(irc.ping)
	close(irc.out)
	for _, c := range irc.subscriptions {
		close(c)
	}
}

func (irc IRC) SendMessages(messages ...string) {
	for _, msg := range messages {
		irc.out <- fmt.Sprintf("PRIVMSG %s :%s", irc.Channel, msg)
	}
}

func (irc IRC) Join(user string, passwd string) {
	irc.out <- fmt.Sprintf("NICK %s", user)
	irc.out <- fmt.Sprintf("USER %s 0.0.0.0 0.0.0.0 :%s", user, user)
	irc.out <- fmt.Sprintf("JOIN %s %s", irc.Channel, passwd)
}

func (irc IRC) Subscribe(pattern *regexp.Regexp, channel chan string) {
	irc.subscriptions[pattern] = channel
}

func (irc *IRC) handleRead() {
	buf := bufio.NewReaderSize(irc.conn, 512)

	for {
		msg, err := buf.ReadString('\n')
		if err != nil {
			if recoverable(err) {
				log.Printf("Error [%s] while reading message, reconnecting in 1s...\n", err)
				<-time.After(1 * time.Second)

				irc.conn = connect(irc.server, irc.port)

				continue
			} else {
				log.Fatalf("Unrecoverable error while reading message: %v\n", err)
			}
		}

		msg = msg[:len(msg)-2]
		if strings.Contains(msg, "PING") {
			irc.ping <- msg
		} else {
			for pattern, channel := range irc.subscriptions {
				if pattern.Match([]byte(msg)) {
					channel <- msg
				}
			}
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
		log.Fatal(err)
	}
}

func connect(server string, port int) net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[IRC] Connected to %s (%s).\n", server, conn.RemoteAddr())
	return conn
}

func recoverable(err error) bool {
	if e, netError := err.(net.Error); netError && e.Temporary() {
		return true
	} else if err == io.EOF {
		return true
	}
	return false
}
