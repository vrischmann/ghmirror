package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/codegangsta/negroni"
	"github.com/google/go-github/github"
	"github.com/vrischmann/envconfig"
)

var conf struct {
	Port                int
	Secret              string
	PersonalAccessToken string
	PollFrequency       time.Duration
	WebHook             struct {
		Endpoint         string
		ValidOwnerLogins []string
	}
	RepositoriesPath string
	DatabasePath     string
}

func main() {
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
		n.Run(fmt.Sprintf(":%d", conf.Port))
	}
}
