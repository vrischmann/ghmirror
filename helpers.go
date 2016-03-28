package main

import (
	"log"
	"os"

	"github.com/vrischmann/ghmirror/internal"
)

func UpdateRepository(r *internal.Repository) error {
	_, err := os.Stat(r.LocalPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		log.Printf("git clone from %s to %s", r.CloneURL, r.LocalPath)
		return gitClone(r.CloneURL, r.LocalPath)
	}

	log.Printf("git pull in %s", r.LocalPath)

	return gitPull(r.LocalPath)
}
