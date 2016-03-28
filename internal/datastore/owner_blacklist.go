package datastore

import (
	"io"

	"github.com/vrischmann/ghmirror/internal"
)

type OwnerBlacklist interface {
	io.Closer

	Get() (internal.OwnersBlacklist, error)
}
