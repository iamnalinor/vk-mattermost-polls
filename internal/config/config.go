package config

import "os"

type Config struct {
	TarantoolAddress  string
	TarantoolUser     string
	TarantoolPassword string
	MattermostUrl     string
	MattermostToken   string
}

func LoadConfig() Config {
	return Config{
		TarantoolAddress:  os.Getenv("TARANTOOL_ADDRESS"),
		TarantoolUser:     os.Getenv("TARANTOOL_USER"),
		TarantoolPassword: os.Getenv("TARANTOOL_PASSWORD"),
		MattermostUrl:     os.Getenv("MATTERMOST_URL"),
		MattermostToken:   os.Getenv("MATTERMOST_TOKEN"),
	}
}
