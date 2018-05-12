package main

import (
	"github.com/FryDay/zorkbot/zorkbot"
)

func main() {
	zorkBot, err := zorkbot.NewBot("zorkbot", "#frybot-test-room", "irc.freenode.net", 7000)
	if err != nil {
		panic(err)
	}
	if err = zorkBot.OpenStory("./stories/zork1.z5"); err != nil {
		panic(err)
	}

	zorkBot.Run()
}
