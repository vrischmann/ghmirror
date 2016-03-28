package postgres

import (
	"database/sql"

	"github.com/vrischmann/ghmirror/internal"
	"github.com/vrischmann/ghmirror/internal/config"
	"github.com/vrischmann/ghmirror/internal/datastore"
)

type ownerBlacklistStore struct {
	db *sql.DB
}

func NewOwnerBlacklistStore(conf *config.Postgres) (datastore.OwnerBlacklist, error) {
	return nil, nil
}

func (s *ownerBlacklistStore) Close() error { return s.db.Close() }

func (s *ownerBlacklistStore) Get() (internal.OwnersBlacklist, error) {
	return nil, nil
}

var _ datastore.OwnerBlacklist = (*ownerBlacklistStore)(nil)
