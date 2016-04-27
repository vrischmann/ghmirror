package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// makeBodyRewindable turns a request's Body into a rewind-able body.
func makeBodyRewindable(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error while reading body. err=%v", err)
		writeInternalServerError(w)
		return
	}
	r.Body.Close()

	r.Body = newRewindableReader(body)

	next(w, r)
}

// hookAuthentication checks that the webhook event is authenticated.
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

// eventTypeValidation checks that the webhook event is of the correct type.
func eventTypeValidation(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	et := r.Header.Get("X-GitHub-Event")
	if et != "push" {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK")
		return
	}

	next(w, r)
}
