package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/vrischmann/ghmirror/internal"
	"github.com/vrischmann/ghmirror/internal/config"
	"github.com/vrischmann/ghmirror/internal/datastore"
	"github.com/vrischmann/ghmirror/internal/postgres"
)

type handler struct {
	conf *config.Config

	rs  datastore.Repository
	obs datastore.OwnerBlacklist
	rbs datastore.RepositoryBlacklist
}

func newHandler(conf *config.Config) (*handler, error) {
	h := &handler{conf: conf}

	var err error

	h.rs, err = postgres.NewRepositoryStore(&conf.Postgres)
	if err != nil {
		return nil, fmt.Errorf("unable to create repository store. err=%v", err)
	}

	h.obs, err = postgres.NewOwnerBlacklistStore(&conf.Postgres)
	if err != nil {
		return nil, fmt.Errorf("unable to create owner blacklist store. err=%v", err)
	}

	h.rbs, err = postgres.NewRepositoryBlacklistStore(&conf.Postgres)
	if err != nil {
		return nil, fmt.Errorf("unable to create repository blacklist store. err=%v", err)
	}

	return h, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rewind(r.Body)

	var hb hookBody
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&hb); err != nil {
		log.Printf("error while decoding json. err=%v", err)
		writeInternalServerError(w)
		return
	}

	// TODO(vincent): transactions

	ok, err := h.rs.Has(hb.Repository.ID)
	if err != nil {
		log.Printf("error while checking for repository in the datastore. err=%v", err)
		writeInternalServerError(w)
		return
	}

	var repo *internal.Repository
	if !ok {
		log.Printf("repository %d does not exist yet, adding it", hb.Repository.ID)

		localPath := filepath.Join(conf.RepositoriesPath, hb.Repository.FullName)
		repo = internal.NewRepository(
			hb.Repository.ID,
			hb.Repository.Name,
			localPath,
			hb.Repository.CloneURL,
		)

		if err := h.rs.Add(repo); err != nil {
			log.Printf("error while adding repository to the datastore. err=%v", err)
			writeInternalServerError(w)
			return
		}
	} else {
		repo, err = h.rs.GetByID(hb.Repository.ID)
		if err != nil {
			log.Printf("error while getting repository from the datastore. err=%v", err)
			writeInternalServerError(w)
			return
		}
	}

	log.Printf("updating repo %d, %s", repo.ID, hb.Repository.FullName)

	if err := UpdateRepository(repo); err != nil {
		log.Printf("error while cloning repository. err=%v", err)
		writeInternalServerError(w)
		return
	}

	log.Printf("repo %d, %s updated", repo.ID, hb.Repository.FullName)

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

func writeForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	io.WriteString(w, "Forbidden")
}

func writeInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, "Oh Noes !")
}
