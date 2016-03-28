package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/vrischmann/ghmirror/internal/config"
	"github.com/vrischmann/ghmirror/internal/datastore"
)

type ownerBlacklistStore struct {
	db *sql.DB
}

func NewOwnerBlacklistStore(conf *config.Postgres) (datastore.OwnerBlacklist, error) {
	s := new(ownerBlacklistStore)

	var err error
	s.db, err = makeDB(conf)

	return s, err
}

func (s *ownerBlacklistStore) Close() error { return s.db.Close() }

func (s *ownerBlacklistStore) IsBlacklisted(name string) (bool, error) {
	const q = `SELECT 1 FROM owner_blacklist
               WHERE name = $1`

	var i int

	err := s.db.QueryRow(q, name).Scan(&i)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

var _ datastore.OwnerBlacklist = (*ownerBlacklistStore)(nil)
