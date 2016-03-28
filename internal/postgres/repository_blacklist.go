package postgres

import (
	"database/sql"

	"github.com/vrischmann/ghmirror/internal"
	"github.com/vrischmann/ghmirror/internal/config"
	"github.com/vrischmann/ghmirror/internal/datastore"
)

type repositoryBlacklistStore struct {
	db *sql.DB
}

func NewRepositoryBlacklistStore(conf *config.Postgres) (datastore.RepositoryBlacklist, error) {
	return nil, nil
}

func (s *repositoryBlacklistStore) Close() error { return s.db.Close() }

func (s *repositoryBlacklistStore) Get() (internal.RepositoriesBlacklist, error) {
	return nil, nil
}

var _ datastore.RepositoryBlacklist = (*repositoryBlacklistStore)(nil)
