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
	SSLMode  string
}

type Config struct {
	ListenAddress       flagutil.NetworkAddresses
	Secret              string
	PersonalAccessToken string
	PollFrequency       time.Duration
	Webhook             struct {
		Endpoint string
	}
	RepositoriesPath string
	Postgres         Postgres
}
