package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/FryDay/zorkbot/config"
	"github.com/FryDay/zorkbot/zorkbot"
	"github.com/spf13/viper"
)

const defaultConfig = `[bot]
nick = "mybot"
password = ""
server = "irc.freenode.net"
port = 7000
channel = "#foo"`

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	confpath := path.Join(usr.HomeDir, ".config")
	if _, err := os.Stat(confpath); os.IsNotExist(err) {
		os.Mkdir(confpath, 0700)
	}

	v := viper.New()
	v.SetConfigName("zorkbot")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	v.AddConfigPath(confpath)
	if _, ok := v.ReadInConfig().(viper.ConfigFileNotFoundError); ok {
		// Create config file
		newConf := []byte(defaultConfig)
		if err := ioutil.WriteFile(path.Join(confpath, "zorkbot.toml"), newConf, 0600); err != nil {
			log.Fatal(err)
		}
	}
	var conf config.Config
	if err = v.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	zorkBot, err := zorkbot.NewBot(&conf)
	if err != nil {
		log.Fatal(err)
	}

	zorkBot.Run()
}
