package config

// Config ...
type Config struct {
	Bot BotConfig
}

// BotConfig ...
type BotConfig struct {
	Nick     string
	Password string
	Server   string
	Port     int64
	Channel  string
}
