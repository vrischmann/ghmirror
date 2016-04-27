package main

import (
	"bytes"
	"io"
)

type rewindableReader struct {
	*bytes.Reader
}

func (b *rewindableReader) Close() error {
	return nil
}

func newRewindableReader(body []byte) io.ReadCloser {
	return &rewindableReader{
		Reader: bytes.NewReader(body),
	}
}

func rewind(r io.Reader) {
	if b, ok := r.(*rewindableReader); ok {
		b.Seek(0, 0)
	}
}
