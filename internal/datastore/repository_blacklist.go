package datastore

import (
	"io"

	"github.com/vrischmann/ghmirror/internal"
)

type RepositoryBlacklist interface {
	io.Closer

	Get() (internal.RepositoriesBlacklist, error)
}
