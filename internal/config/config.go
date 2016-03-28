package config

import (
	"time"

	"github.com/vrischmann/flagutil"
)

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

type Config struct {
	Address             flagutil.NetworkAddresses
	Secret              string
	PersonalAccessToken string
	PollFrequency       time.Duration
	Webhook             struct {
		Endpoint         string
		ValidOwnerLogins []string
	}

	RepositoriesPath string

	Postgres Postgres `envconfig:"optional"`
}
