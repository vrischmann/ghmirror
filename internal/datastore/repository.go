package datastore

import (
	"io"

	"github.com/vrischmann/ghmirror/internal"
)

// Repository is used to get and update metadata about repositories.
type Repository interface {
	io.Closer

	GetAll() (internal.Repositories, error)
	GetByID(id int64) (*internal.Repository, error)
	Has(id int64) (bool, error)
	Add(repo *internal.Repository) error
}
