// Package irc provides a basic implementation of the IRC protocol.
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

// IRC represents an connection to a channel.
type IRC struct {
	// The IRC server to connect to.
	server string

	// The server port to connect to.
	port int

	// The IRC channel to connect to.
	Channel string

	// The connection to the IRC server.
	conn net.Conn

	// The channel where to send PING messages.
	ping chan string

	// The channel where to send messages that should
	// be sent back to the server.
	out chan string

	// A map where the key is a regexp pattern to be matched against,
	// and the value is a channel where to send messages that match
	// the specified pattern.
	subscriptions map[*regexp.Regexp]chan string
}

// New connects to the specified server:port and returns
// an IRC value for interacting with the server.
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

// Close closes the underlying IRC connection.
func (irc IRC) Close() {
	irc.conn.Close()

	close(irc.ping)
	close(irc.out)
	for _, c := range irc.subscriptions {
		close(c)
	}
}

// SendMessages sends the given list of messages over the wire
// to the connected channel.
func (irc IRC) SendMessages(messages ...string) {
	for _, msg := range messages {
		irc.out <- fmt.Sprintf("PRIVMSG %s :%s", irc.Channel, msg)
	}
}

// Join joins the configured channel with the given
// user credentials.
func (irc IRC) Join(user string, passwd string) {
	irc.out <- fmt.Sprintf("NICK %s", user)
	irc.out <- fmt.Sprintf("USER %s 0.0.0.0 0.0.0.0 :%s", user, user)
	irc.out <- fmt.Sprintf("JOIN %s %s", irc.Channel, passwd)
}

// Subscribe configures a message subscription pattern that,
// when matched, causes the message to be sent to the specified
// channel.
func (irc IRC) Subscribe(pattern *regexp.Regexp, channel chan string) {
	irc.subscriptions[pattern] = channel
}

// handleRead reads all messages sent to the IRC channel.
// If it's a "PING" message, forwards it to the ping channel;
// otherwise, looks for a subscription that matches the message
// and forwards it to the registered channel.
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

// handleWrite reads messages from the out channel
// and sends them over the wire.
func (irc IRC) handleWrite() {
	for msg := range irc.out {
		irc.send(msg)
	}
}

// handlePing reads messages from the ping channel
// and sends the "PONG" response to the server originating
// the "PING" request.
func (irc IRC) handlePing() {
	for ping := range irc.ping {
		server := strings.Split(ping, ":")[1]

		irc.out <- fmt.Sprintf("PONG %s", server)
		log.Printf("[IRC] PONG sent to %s\n", server)
	}
}

// send is responsible for writing the bytes over the wire.
func (irc IRC) send(msg string) {
	_, err := irc.conn.Write([]byte(fmt.Sprintf("%s\r\n", msg)))
	if err != nil {
		log.Fatal(err)
	}
}

// connect dials to the configured server and returns
// the connection.
func connect(server string, port int) net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[IRC] Connected to %s (%s).\n", server, conn.RemoteAddr())
	return conn
}

// recoverable checks if the given error is temporary and could
// be recovered from.
func recoverable(err error) bool {
	if e, netError := err.(net.Error); netError && e.Temporary() {
		return true
	} else if err == io.EOF {
		return true
	}
	return false
}
