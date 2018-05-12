package zorkbot

import (
	"crypto/tls"
	"fmt"
	"os"
	"regexp"
	"strings"

	"bitbucket.org/zombiezen/gonorth/north"
	"github.com/thoj/go-ircevent"
)

// Bot ...
type Bot struct {
	Nick         string
	Room         string
	Server       *Server
	machine      *north.Machine
	joined       chan bool
	inputBuffer  *inputBuffer
	outputBuffer *outputBuffer
	connection   *irc.Connection
}

// NewBot ...
func NewBot(nick, room, serverURL string, serverPort int64) (*Bot, error) {
	var err error
	newConnection := irc.IRC(nick, nick)
	// newConnection.Debug = true
	newConnection.UseTLS = true
	newConnection.TLSConfig = &tls.Config{}

	newBot := &Bot{
		Nick:         nick,
		Room:         room,
		joined:       make(chan bool, 1),
		inputBuffer:  newInputBuffer(),
		outputBuffer: newOutputBuffer(),
		connection:   newConnection,
	}

	newBot.Server, err = NewServer(serverURL, serverPort)
	if err != nil {
		return nil, err
	}

	newBot.connection.AddCallback("001", func(e *irc.Event) { newBot.connection.Join(room) })
	newBot.connection.AddCallback("JOIN", newBot.join)
	newBot.connection.AddCallback("PRIVMSG", newBot.mention)
	if err = newBot.connection.Connect(newBot.Server.String()); err != nil {
		return nil, err
	}

	return newBot, nil
}

// OpenStory ...
func (b *Bot) OpenStory(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	b.machine, err = north.NewMachine(f, b)
	return err
}

// Run ...
func (b *Bot) Run() {
	go b.connection.Loop()
	<-b.joined // Wait until we join the room
	go b.watchOutput()
	if err := b.machine.Run(); err != nil {
		panic(err)
	}
}

func (b *Bot) watchOutput() {
	for {
		msg := <-b.outputBuffer.msg
		if msg == "" {
			continue
		}
		b.connection.Privmsg(b.Room, msg)
	}
}

// Input ...
func (b *Bot) Input(n int) ([]rune, error) {
	r := make([]rune, 0, n)
	for {
		rr, _, err := b.inputBuffer.ReadRune()
		if err != nil {
			return r, err
		} else if rr == '\n' {
			break
		}
		if len(r) < n {
			r = append(r, rr)
		}
	}
	return r, nil
}

// Output ...
func (b *Bot) Output(window int, s string) error {
	if window != 0 {
		return nil
	}
	_, err := fmt.Print(s)
	b.outputBuffer.WriteString(s)
	return err
}

// ReadRune ...
func (b *Bot) ReadRune() (rune, int, error) {
	return b.inputBuffer.ReadRune()
}

// Save ...
func (b *Bot) Save(m *north.Machine) error {
	return nil
}

// Restore ...
func (b *Bot) Restore(m *north.Machine) error {
	return nil
}

func (b *Bot) join(e *irc.Event) {
	b.joined <- true
}

//TODO: Cleanup
func (b *Bot) mention(e *irc.Event) {
	lower := strings.ToLower(e.Message())

	if strings.HasPrefix(lower, b.Nick) {
		var regex = regexp.MustCompile(`^[^:, ]+[:, ]`)

		message := strings.TrimSpace(strings.Replace(e.Message(), regex.FindString(e.Message()), "", 1))
		lower = strings.ToLower(message)
		if strings.HasPrefix(lower, b.Nick) {
			for {
				message = strings.Replace(message, b.Nick, "", 1)
				lower = strings.ToLower(message)
				if !strings.HasPrefix(lower, b.Nick) {
					break
				}
			}
		}

		fmt.Println(lower)
		b.inputBuffer.WriteString(lower + "\n")
	}
}