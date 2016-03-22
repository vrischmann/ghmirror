package main

import "io"

// DataStore is the interface used to store metadata about repositories.
type DataStore interface {
	io.Closer
	Lock()
	Unlock()
	Repositories() (Repositories, error)
	GetByID(id int64) (*Repository, error)
	HasRepository(id int64) (bool, error)
	AddRepository(repo *Repository) error
}
