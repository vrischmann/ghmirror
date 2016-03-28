package postgres

import (
	"database/sql"

	"github.com/vrischmann/ghmirror/internal"
	"github.com/vrischmann/ghmirror/internal/config"
	"github.com/vrischmann/ghmirror/internal/datastore"
)

type repositoryStore struct {
	db *sql.DB
}

func NewRepositoryStore(conf *config.Postgres) (datastore.Repository, error) {
	s := new(repositoryStore)

	var err error
	s.db, err = makeDB(conf)

	return s, err
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

	err := s.db.QueryRow(q, id).Scan(&name, &localPath, &cloneURL, &hookID)
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

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// TODO(vincent): do we need the last inserted id for something ?
	_, err = tx.Exec(q, repo.Name, repo.LocalPath, repo.CloneURL, repo.HookID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

var _ datastore.Repository = (*repositoryStore)(nil)
