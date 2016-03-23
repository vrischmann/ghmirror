package main

import (
	"database/sql"
	"fmt"
	"sync"
)

type postgresDataStore struct {
	db *sql.DB
	mu sync.Mutex
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

func (s *postgresDataStore) Lock() {
	s.mu.Lock()
}

func (s *postgresDataStore) Unlock() {
	s.mu.Unlock()
}

func (s *postgresDataStore) Repositories() (Repositories, error) {
	var res Repositories

	const q = `SELECT id, name, local_path, clone_url, hook_id FROM repository`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		id, hookID                int64
		name, localPath, cloneURL string
	)

	for rows.Next() {
		if err := rows.Scan(&id, &name, &localPath, &cloneURL, &hookID); err != nil {
			return nil, err
		}

		repo := &Repository{
			ID:        id,
			Name:      name,
			LocalPath: localPath,
			CloneURL:  cloneURL,
			HookID:    hookID,
		}

		res = append(res, repo)
	}

	return res, nil
}

func (s *postgresDataStore) GetByID(id int64) (*Repository, error) {
	const q = `SELECT name, local_path, clone_url, hook_id FROM repository
               WHERE id = $1`

	var (
		name, localPath, cloneURL string
		hookID                    int64
	)

	err := s.db.QueryRow(q).Scan(&name, &localPath, &cloneURL, &hookID)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	repo := &Repository{
		ID:        id,
		Name:      name,
		LocalPath: localPath,
		CloneURL:  cloneURL,
		HookID:    hookID,
	}

	return repo, nil
}

func (s *postgresDataStore) HasRepository(id int64) (bool, error) {
	repo, err := s.GetByID(id)
	return repo != nil, err
}

func (s *postgresDataStore) AddRepository(repo *Repository) error {
	const q = `INSERT INTO repository(name, local_path, clone_url, hook_id)
               VALUES ($1, $2, $3, $4)`

	// TODO(vincent): do we need the last inserted id for something ?
	_, err := s.db.Exec(q, repo.Name, repo.LocalPath, repo.CloneURL, repo.HookID)
	return err
}
