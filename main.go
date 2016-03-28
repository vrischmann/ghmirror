package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/vrischmann/envconfig"
	"github.com/vrischmann/ghmirror/internal/config"
)

var (
	version string
	commit  string

	conf config.Config
)

func main() {
	log.Printf("ghmirror %s-%s", version, commit)
	log.Printf("config: %+v", conf)

	err := envconfig.Init(&conf)
	if err != nil {
		log.Fatal(err)
	}

	poller, err := newPoller(&conf)
	if err != nil {
		log.Fatal(err)
	}
	poller.run()

	handler, err := newHandler(&conf)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/hook", handler)

	// TODO(vincent): replace negroni

	n := negroni.Classic()
	n.UseFunc(makeBodyRewindable)
	n.UseFunc(hookAuthentication)
	n.UseFunc(eventTypeValidation)
	n.UseHandler(mux)
	n.Run(string(conf.Address.StringSlice()[0]))
}
