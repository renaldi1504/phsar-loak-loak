// Package secretbox: enkripsi AES-256-GCM untuk data sensitif di database (mis. nomor rekening seller).
package secretbox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

type Box struct{ aead cipher.AEAD }

// New: key = 32 byte dalam base64 (DATA_ENCRYPTION_KEY).
func New(b64 string) (*Box, error) {
	key, err := base64.StdEncoding.DecodeString(b64)
	if err != nil || len(key) != 32 {
		return nil, fmt.Errorf("DATA_ENCRYPTION_KEY harus 32 byte base64 (buat: openssl rand -base64 32)")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Box{aead}, nil
}

// Seal: nonce || ciphertext. aad mengikat data ke konteksnya (mis. id toko) supaya tidak bisa dipindah ke baris lain.
func (b *Box) Seal(plain, aad []byte) ([]byte, error) {
	nonce := make([]byte, b.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return b.aead.Seal(nonce, nonce, plain, aad), nil
}

func (b *Box) Open(sealed, aad []byte) ([]byte, error) {
	n := b.aead.NonceSize()
	if len(sealed) < n {
		return nil, errors.New("secretbox: data terlalu pendek")
	}
	return b.aead.Open(nil, sealed[:n], sealed[n:], aad)
}
