package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

type handler struct {
	ds DataStore
}

func newHandler(ds DataStore) *handler {
	return &handler{ds: ds}
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

	// Obtain the datastore lock
	h.ds.Lock()
	defer h.ds.Unlock()

	ok, err := h.ds.HasRepository(hb.Repository.Id)
	if err != nil {
		log.Printf("error while checking for repository in the datastore. err=%v", err)
		writeInternalServerError(w)
		return
	}

	var repo *Repository
	{
		if !ok {
			log.Printf("repository %d does not exist yet, adding it", hb.Repository.Id)

			localPath := filepath.Join(conf.RepositoriesPath, hb.Repository.FullName)
			repo = NewRepository(
				hb.Repository.Id,
				hb.Repository.Name,
				localPath,
				hb.Repository.CloneURL,
			)

			if err := h.ds.AddRepository(repo); err != nil {
				log.Printf("error while adding repository to the datastore. err=%v", err)
				writeInternalServerError(w)
				return
			}
		} else {
			repo, err = h.ds.GetByID(hb.Repository.Id)
			if err != nil {
				log.Printf("error while getting repository from the datastore. err=%v", err)
				writeInternalServerError(w)
				return
			}
		}
	}

	log.Printf("updating repo %d, %s", repo.ID, hb.Repository.FullName)

	if err := repo.Update(); err != nil {
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
