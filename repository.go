package main

import (
	"bufio"
	"encoding/binary"
	"io"
)

type Repository struct {
	Name      string
	ID        int64
	LocalPath string
	CloneURL  string
}

func NewRepository(id int64, name, localPath, cloneURL string) *Repository {
	return &Repository{
		ID:        id,
		Name:      name,
		LocalPath: localPath,
		CloneURL:  cloneURL,
	}
}

func (r *Repository) Pull() error {
	return gitPull(r.LocalPath)
}

func (r *Repository) Clone() error {
	return gitClone(r.CloneURL, r.LocalPath)
}

func (r *Repository) Read(p []byte) (n int, err error) {
	return -1, nil
}

func (r *Repository) Write(p []byte) (n int, err error) {
	return -1, nil
}

func (r *Repository) ReadFrom(rd io.Reader) (n int64, err error) {
	var i int64
	br := bufio.NewReader(rd)

	binary.Read(br, binary.BigEndian, &r.ID)
	n += 8

	if i, err = readString(br, &r.Name); err != nil {
		return
	}
	n += i

	if i, err = readString(br, &r.LocalPath); err != nil {
		return
	}
	n += i

	if i, err = readString(br, &r.CloneURL); err != nil {
		return
	}
	n += i

	return
}

func (r *Repository) WriteTo(w io.Writer) (n int64, err error) {
	bw := bufio.NewWriter(w)

	if err = binary.Write(bw, binary.BigEndian, r.ID); err != nil {
		return
	}
	if err = binary.Write(bw, binary.BigEndian, int32(len(r.Name))); err != nil {
		return
	}
	bw.WriteString(r.Name)
	if err = binary.Write(bw, binary.BigEndian, int32(len(r.LocalPath))); err != nil {
		return
	}
	bw.WriteString(r.LocalPath)
	if err = binary.Write(bw, binary.BigEndian, int32(len(r.CloneURL))); err != nil {
		return
	}
	bw.WriteString(r.CloneURL)

	return int64(bw.Buffered()), bw.Flush()
}

func EncodeID(id int64) []byte {
	buf := make([]byte, 8)
	binary.PutVarint(buf, id)
	return buf
}

type Repositories []*Repository

func readString(r io.Reader, s *string) (n int64, err error) {
	var length int32
	var i int

	if err = binary.Read(r, binary.BigEndian, &length); err != nil {
		return
	}

	buf := make([]byte, length)
	if i, err = io.ReadFull(r, buf); err != nil {
		return
	}
	n = int64(i)

	*s = string(buf)

	return
}
