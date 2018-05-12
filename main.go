package main

import (
	"bufio"

	"github.com/FryDay/zorkbot/bot"
)

var in *bufio.Reader

func main() {
	zorkBot, err := zorkbot.NewBot("zorkbot", "#frybot-test-room", "irc.freenode.net", 7000)
	if err != nil {
		panic(err)
	}
	if err = zorkBot.OpenStory("./stories/zork1.z5"); err != nil {
		panic(err)
	}

	zorkBot.Run()

	// in = bufio.NewReader(os.Stdin)
	// m, err := openStory("./stories/zork1.z5")
	// if err != nil {
	// 	panic(err)
	// }
	// for {
	// 	err = m.Run()
	// 	switch err {
	// 	case io.EOF, north.ErrQuit:
	// 		os.Exit(0)
	// 	case north.ErrRestart:
	// 		m, err = openStory(flag.Arg(0))
	// 		if err != nil {
	// 			fmt.Fprintln(os.Stderr, err)
	// 			os.Exit(1)
	// 		}
	// 	default:
	// 		fmt.Fprintln(os.Stderr, "** Internal Error:", err)
	// 		os.Exit(1)
	// 	}
	// }
}

// func openStory(path string) (*north.Machine, error) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	return north.NewMachine(f, new(terminalUI))
// }

// type terminalUI struct{}

// func (t *terminalUI) Input(n int) ([]rune, error) {
// 	r := make([]rune, 0, n)
// 	for {
// 		rr, _, err := in.ReadRune()
// 		if err != nil {
// 			return r, err
// 		} else if rr == '\n' {
// 			break
// 		}
// 		if len(r) < n {
// 			r = append(r, rr)
// 		}
// 	}
// 	return r, nil
// }

// func (t *terminalUI) Output(window int, s string) error {
// 	if window != 0 {
// 		return nil
// 	}
// 	_, err := fmt.Print(s)
// 	return err
// }

// func (t *terminalUI) ReadRune() (rune, int, error) {
// 	return in.ReadRune()
// }

// func (t *terminalUI) Save(m *north.Machine) error {
// 	return nil
// }

// func (t *terminalUI) Restore(m *north.Machine) error {
// 	return nil
// }
