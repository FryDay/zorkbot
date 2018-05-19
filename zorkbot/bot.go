package zorkbot

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"bitbucket.org/zombiezen/gonorth/north"
	"github.com/FryDay/zorkbot/config"
	"github.com/thoj/go-ircevent"
)

// Bot ...
type Bot struct {
	Nick         string
	Channel      string
	Password     string
	Server       *Server
	nickRegex    *regexp.Regexp
	machine      *north.Machine
	joined       chan bool
	inputBuffer  *inputBuffer
	outputBuffer *outputBuffer
	connection   *irc.Connection
}

// NewBot ...
func NewBot(conf *config.Config) (*Bot, error) {
	var err error
	newConnection := irc.IRC(conf.Bot.Nick, conf.Bot.Nick)
	// newConnection.Debug = true
	newConnection.UseTLS = true
	newConnection.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	newBot := &Bot{
		Nick:         conf.Bot.Nick,
		Channel:      conf.Bot.Channel,
		Password:     conf.Bot.Password,
		joined:       make(chan bool, 1),
		inputBuffer:  newInputBuffer(),
		outputBuffer: newOutputBuffer(),
		connection:   newConnection,
	}

	if newBot.nickRegex, err = regexp.Compile(`^[^:, ]+[:, ]`); err != nil {
		return nil, err
	}

	newBot.Server, err = NewServer(conf.Bot.Server, conf.Bot.Port)
	if err != nil {
		return nil, err
	}

	newBot.connection.AddCallback("001", func(e *irc.Event) { newBot.connection.Join(conf.Bot.Channel) })
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
	err := b.OpenStory("./stories/zork1.z5")
	for {
		err = b.machine.Run()
		switch err {
		case io.EOF, north.ErrQuit:
			os.Exit(0)
		case north.ErrRestart:
			err = b.OpenStory("./stories/zork1.z5")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		default:
			fmt.Fprintln(os.Stderr, "** Internal Error:", err)
			os.Exit(1)
		}
	}
}

func (b *Bot) watchOutput() {
	for {
		msg := <-b.outputBuffer.msg
		if msg == "" {
			continue
		}
		b.connection.Privmsg(b.Channel, msg)
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
	if b.Password != "" {
		b.connection.Privmsgf("NickServ", "identify %s", b.Password)
	}
	b.joined <- true
}

func (b *Bot) login(e *irc.Event) {

}

func (b *Bot) mention(e *irc.Event) {
	lower := strings.ToLower(e.Message())

	if strings.HasPrefix(lower, b.Nick) {
		message := strings.TrimSpace(strings.Replace(e.Message(), b.nickRegex.FindString(e.Message()), "", 1))
		lower = strings.ToLower(message)

		fmt.Println(lower)
		b.inputBuffer.WriteString(lower + "\n")
	}
}
