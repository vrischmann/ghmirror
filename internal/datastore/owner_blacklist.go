package datastore

import "io"

type OwnerBlacklist interface {
	io.Closer

	IsBlacklisted(name string) (bool, error)
}
