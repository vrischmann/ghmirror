package main

import (
	"database/sql"
	"fmt"
)

type postgresDataStore struct {
	db *sql.DB
}

func newPostgresDataStore(conf *postgresConf) (DataStore, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=verify-full", conf.Host, conf.Port, conf.User, conf.Dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &postgresDataStore{db: db}, nil
}

func (s *postgresDataStore) Close() error {
	return s.db.Close()
}

func (s *postgresDataStore) Lock()   {}
func (s *postgresDataStore) Unlock() {}

func (s *postgresDataStore) Repositories() (res Repositories, err error) {

	return
}

func (s *postgresDataStore) GetByID(id int64) (*Repository, error) { return nil, nil }
func (s *postgresDataStore) HasRepository(id int64) (bool, error)  { return false, nil }
func (s *postgresDataStore) AddRepository(repo *Repository) error  { return nil }
