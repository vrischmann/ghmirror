package main

import (
	"bytes"
	"io"
	"sync"

	"github.com/boltdb/bolt"
)

const (
	repositoriesBucket_ = "repositories"
)

var (
	repositoriesBucket = []byte(repositoriesBucket_)
)

// DataStore is the interface used to store metadata about repositories.
type DataStore interface {
	io.Closer
	Lock()
	Unlock()
	Repositories() (Repositories, error)
	GetByID(id int64) (*Repository, error)
	HasRepository(id int64) (bool, error)
	AddRepository(repo *Repository) error
}

type boltDataStore struct {
	db *bolt.DB
	mu sync.Mutex
}

func newBoltDataStore(path string) (DataStore, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &boltDataStore{db: db}, nil
}

func (s *boltDataStore) Close() error {
	return s.db.Close()
}

func (s *boltDataStore) Lock() {
	s.mu.Lock()
}

func (s *boltDataStore) Unlock() {
	s.mu.Unlock()
}

func (s *boltDataStore) Repositories() (Repositories, error) {
	return nil, nil
}

func (s *boltDataStore) GetByID(id int64) (*Repository, error) {
	var repo *Repository

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(repositoriesBucket)
		if bucket == nil {
			return nil
		}

		data := bucket.Get(EncodeID(id))
		if data == nil {
			return nil
		}

		var tmp Repository

		_, err := tmp.ReadFrom(bytes.NewReader(data))
		if err != nil {
			return err
		}

		repo = &tmp

		return nil
	})

	return repo, err
}

func (s *boltDataStore) HasRepository(id int64) (bool, error) {
	res := false
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(repositoriesBucket)
		if bucket == nil {
			return nil
		}

		data := bucket.Get(EncodeID(id))
		if data != nil {
			res = true
		}

		return nil
	})
	return res, err
}

func (s *boltDataStore) AddRepository(repo *Repository) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(repositoriesBucket)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		_, err = io.Copy(&buf, repo)
		if err != nil {
			return err
		}

		return bucket.Put(EncodeID(repo.ID), buf.Bytes())
	})
}
