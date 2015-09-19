package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/codegangsta/negroni"
	"github.com/google/go-github/github"
	"github.com/vrischmann/envconfig"
)

type networkAddress string

func (n *networkAddress) Unmarshal(s string) error {
	_, _, err := net.SplitHostPort(s)
	if err != nil {
		return err
	}

	*n = networkAddress(s)

	return nil
}

var conf struct {
	Address             networkAddress
	Secret              string
	PersonalAccessToken string
	PollFrequency       time.Duration
	Webhook             struct {
		Endpoint         string
		ValidOwnerLogins []string
	}
	RepositoriesPath string
	DatabasePath     string
}

var (
	version string
	commit  string
)

func main() {
	log.Printf("ghmirror %s-%s", version, commit)

	var ds DataStore
	var gh *github.Client
	{
		err := envconfig.Init(&conf)
		if err != nil {
			log.Fatalln(err)
		}

		ds, err = newBoltDataStore(conf.DatabasePath)
		if err != nil {
			log.Fatalln(err)
		}
		defer ds.Close()

		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.PersonalAccessToken})
		tc := oauth2.NewClient(oauth2.NoContext, ts)

		gh = github.NewClient(tc)
	}

	{
		poller := poller{
			ds:   ds,
			gh:   gh,
			freq: conf.PollFrequency,
		}
		go poller.run()
	}

	{
		handler := newHandler(ds)

		mux := http.NewServeMux()
		mux.Handle("/hook", handler)

		n := negroni.Classic()
		n.UseFunc(makeBodyRewindable)
		n.UseFunc(hookAuthentication)
		n.UseFunc(eventTypeValidation)
		n.UseHandler(mux)
		n.Run(string(conf.Address))
	}
}
