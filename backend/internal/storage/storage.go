// Package storage menyimpan file upload. Sekarang di disk lokal; nanti diganti Cloudflare R2
// dengan interface yang sama (key tetap sama, jadi data di database tidak perlu diubah).
package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Storage interface {
	Put(ctx context.Context, key string, r io.Reader) (int64, error)
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	// PublicURL mengembalikan URL untuk key public/…; key private/… tidak punya URL publik.
	PublicURL(key string) string
}

var ErrBadKey = errors.New("key file tidak valid")

type Local struct {
	Dir     string // mis. ./data/uploads
	BaseURL string // mis. http://localhost:8080/media
}

func (l Local) path(key string) (string, error) {
	clean := filepath.ToSlash(filepath.Clean(key))
	if clean != key || strings.HasPrefix(clean, "/") || strings.Contains(clean, "..") ||
		!(strings.HasPrefix(clean, "public/") || strings.HasPrefix(clean, "private/")) {
		return "", ErrBadKey
	}
	return filepath.Join(l.Dir, filepath.FromSlash(clean)), nil
}

func (l Local) Put(_ context.Context, key string, r io.Reader) (int64, error) {
	p, err := l.path(key)
	if err != nil {
		return 0, err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
		return 0, err
	}
	tmp := p + ".part"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o640)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(f, r)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		_ = os.Remove(tmp)
		return 0, err
	}
	return n, os.Rename(tmp, p)
}

func (l Local) Open(_ context.Context, key string) (io.ReadCloser, error) {
	p, err := l.path(key)
	if err != nil {
		return nil, err
	}
	return os.Open(p)
}

func (l Local) PublicURL(key string) string {
	if key == "" || !strings.HasPrefix(key, "public/") {
		return ""
	}
	return strings.TrimRight(l.BaseURL, "/") + "/" + key
}
