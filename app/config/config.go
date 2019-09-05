package config

import (
	"flag"
	"os"
)

var (
	Client string
	User   string
)

func init() {
	if clientEnv := os.Getenv("CLIENT_KEY"); clientEnv != "" {
		Client = clientEnv
	}

	if userEnv := os.Getenv("USER_KEY"); userEnv != "" {
		User = userEnv
	}

	flag.StringVar(&Client, "client", Client, "Gracenote Client ID")
	flag.StringVar(&User, "user", User, "Gracenote User ID")
}
