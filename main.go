package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/codegangsta/negroni"
	"github.com/google/go-github/github"
	"github.com/vrischmann/envconfig"
	"github.com/vrischmann/flagutil"
)

type dataStoreType int

func (t *dataStoreType) Unmarshal(s string) error {
	switch strings.ToLower(s) {
	case "bolt":
		*t = boltDataStoreType
	case "postgres":
		*t = postgresDataStoreType
	default:
		return fmt.Errorf("unknown data store type '%s'", s)
	}

	return nil
}

const (
	boltDataStoreType dataStoreType = iota + 1
	postgresDataStoreType
)

type boltConf struct {
	DatabasePath string
}

type postgresConf struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

var conf struct {
	Address             flagutil.NetworkAddresses
	Secret              string
	PersonalAccessToken string
	PollFrequency       time.Duration
	Webhook             struct {
		Endpoint         string
		ValidOwnerLogins []string
	}

	RepositoriesPath string

	DataStoreType dataStoreType

	Bolt     boltConf     `envconfig:"optional"`
	Postgres postgresConf `envconfig:"optional"`
}

var (
	version string
	commit  string
)

func main() {
	log.Printf("ghmirror %s-%s", version, commit)

	var (
		ds DataStore
		gh *github.Client
	)

	err := envconfig.Init(&conf)
	if err != nil {
		log.Fatal(err)
	}

	switch conf.DataStoreType {
	case boltDataStoreType:
		ds, err = newBoltDataStore(&conf.Bolt)
		if err != nil {
			log.Fatal(err)
		}
		defer ds.Close()

	case postgresDataStoreType:
		ds, err = newPostgresDataStore(&conf.Postgres)
		if err != nil {
			log.Fatal(err)
		}
		defer ds.Close()
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.PersonalAccessToken})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	gh = github.NewClient(tc)

	poller := poller{
		ds:   ds,
		gh:   gh,
		freq: conf.PollFrequency,
	}
	go poller.run()

	handler := newHandler(ds)

	mux := http.NewServeMux()
	mux.Handle("/hook", handler)

	n := negroni.Classic()
	n.UseFunc(makeBodyRewindable)
	n.UseFunc(hookAuthentication)
	n.UseFunc(eventTypeValidation)
	n.UseHandler(mux)
	n.Run(string(conf.Address.StringSlice()[0]))
}
