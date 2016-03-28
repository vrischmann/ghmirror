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
	err := envconfig.Init(&conf)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ghmirror %s-%s", version, commit)
	log.Printf("listen address: %v", conf.ListenAddress)
	log.Printf("postgres conf: %+v", conf.Postgres)

	// TODO(vincent): maybe it's better to keep one sql.DB

	poller, err := newPoller(&conf)
	if err != nil {
		log.Fatal(err)
	}
	go poller.run()

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
	n.Run(string(conf.ListenAddress.StringSlice()[0]))
}
