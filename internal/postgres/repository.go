package postgres

import (
	"database/sql"
	"fmt"

	"github.com/vrischmann/ghmirror/internal"
	"github.com/vrischmann/ghmirror/internal/config"
	"github.com/vrischmann/ghmirror/internal/datastore"
)

type repositoryStore struct {
	db *sql.DB
}

func NewRepositoryStore(conf *config.Postgres) (datastore.Repository, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=verify-full", conf.Host, conf.Port, conf.User, conf.Dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &repositoryStore{db: db}, nil
}

func (s *repositoryStore) Close() error { return s.db.Close() }

func (s *repositoryStore) GetAll() (internal.Repositories, error) {
	var res internal.Repositories

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

		repo := &internal.Repository{
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

func (s *repositoryStore) GetByID(id int64) (*internal.Repository, error) {
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

	repo := &internal.Repository{
		ID:        id,
		Name:      name,
		LocalPath: localPath,
		CloneURL:  cloneURL,
		HookID:    hookID,
	}

	return repo, nil
}

func (s *repositoryStore) Has(id int64) (bool, error) {
	repo, err := s.GetByID(id)
	return repo != nil, err
}

func (s *repositoryStore) Add(repo *internal.Repository) error {
	const q = `INSERT INTO repository(name, local_path, clone_url, hook_id)
               VALUES ($1, $2, $3, $4)`

	// TODO(vincent): do we need the last inserted id for something ?
	_, err := s.db.Exec(q, repo.Name, repo.LocalPath, repo.CloneURL, repo.HookID)
	return err
}

var _ datastore.Repository = (*repositoryStore)(nil)
