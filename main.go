package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/vrischmann/envconfig"
)

type bufferedBody struct {
	*bytes.Reader
}

func (b *bufferedBody) Close() error {
	return nil
}

func newBufferedBody(body []byte) io.ReadCloser {
	return &bufferedBody{
		Reader: bytes.NewReader(body),
	}
}

type Config struct {
	Port             int
	Secret           string
	RepositoriesPath string
	DatabasePath     string
}

var (
	conf Config
	ds   DataStore
)

func rewind(r io.Reader) {
	if b, ok := r.(*bufferedBody); ok {
		b.Seek(0, 0)
	}
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
	rewind(r.Body)

	var hb hookBody
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&hb); err != nil {
		log.Printf("error while decoding json. err=%v", err)
		writeInternalServerError(w)
		return
	}

	ok, err := ds.HasRepository(hb.Repository.Id)
	if err != nil {
		log.Printf("error while checking for repository in the datastore. err=%v", err)
		writeInternalServerError(w)
		return
	}

	if !ok {
		localPath := filepath.Join(conf.RepositoriesPath, hb.Repository.FullName)
		repo = NewRepository(
			hb.Repository.Id,
			hb.Repository.Name,
			localPath,
			hb.Repository.CloneURL,
		)

		if err := ds.AddRepository(repo); err != nil {
			log.Printf("error while adding repository to the datastore. err=%v", err)
			writeInternalServerError(w)
			return
		}

		if err := repo.Clone(); err != nil {
			log.Printf("error while cloning repository. err=%v", err)
			writeInternalServerError(w)
			return
		}
	} else {
		repo, err := ds.GetByID(hb.Repository.Id)
		if err != nil {
			log.Printf("error while getting repository from the datastore. err=%v", err)
			writeInternalServerError(w)
			return
		}

		if err := repo.Pull(); err != nil {
			log.Printf("error while pulling repository. err=%v", err)
			writeInternalServerError(w)
			return
		}
	}

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

func bufferizeBody(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error while reading body. err=%v", err)
		writeInternalServerError(w)
		return
	}
	r.Body.Close()

	r.Body = newBufferedBody(body)

	next(w, r)
}

func hookAuthentication(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rewind(r.Body)

	sign := r.Header.Get("X-Hub-Signature")
	if !strings.HasPrefix(sign, "sha1=") {
		writeForbidden(w)
		return
	}

	messageMAC, err := hex.DecodeString(strings.Split(sign, "=")[1])
	if err != nil {
		log.Printf("error while decoding message MAC. err=%v", err)
		writeForbidden(w)
		return
	}

	mac := hmac.New(sha1.New, []byte(conf.Secret))
	io.Copy(mac, r.Body)
	expectedMAC := mac.Sum(nil)

	if !hmac.Equal(expectedMAC, messageMAC) {
		writeForbidden(w)
		return
	}

	next(w, r)
}

func eventTypeValidation(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	et := r.Header.Get("X-GitHub-Event")
	if et != "push" {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK")
		return
	}

	next(w, r)
}

func main() {
	err := envconfig.Init(&conf)
	if err != nil {
		log.Fatalln(err)
	}

	if ds, err = newBoltDataStore(conf.DatabasePath); err != nil {
		log.Fatalln(err)
	}
	defer ds.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/hook", hookHandler)

	n := negroni.Classic()
	n.UseFunc(bufferizeBody)
	n.UseFunc(hookAuthentication)
	n.UseFunc(eventTypeValidation)
	n.UseHandler(mux)
	n.Run(fmt.Sprintf(":%d", conf.Port))
}
