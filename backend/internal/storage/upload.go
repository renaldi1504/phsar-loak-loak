package storage

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Aturan per jenis upload.
type rule struct {
	private  bool
	maxBytes int64
	types    map[string]string // content-type hasil deteksi → ekstensi
	folder   string
	typeMax  map[string]int64 // batas khusus per content-type (mis. video lebih besar dari foto)
}

var images = map[string]string{"image/jpeg": "jpg", "image/png": "png", "image/webp": "webp"}

var rules = map[string]rule{
	"store_logo":   {false, 5 << 20, images, "public/stores", nil},
	"store_banner": {false, 8 << 20, images, "public/stores", nil},
	"pickup_photo": {false, 8 << 20, images, "public/pickup", nil},
	"kyc_document": {true, 10 << 20, map[string]string{"image/jpeg": "jpg", "image/png": "png", "image/webp": "webp", "application/pdf": "pdf"}, "private/kyc", nil},
	// Foto 8 MB; video maks. 100 MB (± 3 menit dari HP). Transcode/kompres menyusul.
	"product_media": {false, MaxVideoBytes, map[string]string{"image/jpeg": "jpg", "image/png": "png", "image/webp": "webp",
		"video/mp4": "mp4", "video/quicktime": "mov", "video/webm": "webm"}, "public/products",
		map[string]int64{"image/jpeg": 8 << 20, "image/png": 8 << 20, "image/webp": 8 << 20}},
}

// MaxVideoBytes: batas upload terbesar (dipakai juga untuk membatasi body request).
const MaxVideoBytes = 100 << 20

// IsVideo / IsImage untuk content-type hasil deteksi.
func IsVideo(ct string) bool { return len(ct) > 6 && ct[:6] == "video/" }
func IsImage(ct string) bool { return len(ct) > 6 && ct[:6] == "image/" }

// sniff: http.DetectContentType tidak mengenali .mov dari iPhone (brand "qt  "), jadi cek kotak ftyp sendiri.
func sniff(head []byte) string {
	if len(head) >= 12 && string(head[4:8]) == "ftyp" {
		if string(head[8:12]) == "qt  " {
			return "video/quicktime"
		}
		return "video/mp4"
	}
	return http.DetectContentType(head)
}

var (
	ErrKind     = errors.New("jenis upload tidak dikenal")
	ErrTooLarge = errors.New("file terlalu besar")
	ErrType     = errors.New("jenis file tidak diizinkan")
	ErrNotOwner = errors.New("file tidak ditemukan atau bukan milik kamu")
)

type Uploaded struct {
	Key         string `json:"key"`
	URL         string `json:"url,omitempty"` // kosong untuk file privat (KYC)
	ContentType string `json:"content_type"`
	Bytes       int64  `json:"bytes"`
}

type Uploader struct {
	Store Storage
	Pool  *pgxpool.Pool
	// CheckVideo (opsional): dipanggil setelah video tersimpan, mis. cek durasi dengan ffprobe.
	// Error → upload ditolak (dan tidak dicatat di tabel uploads).
	CheckVideo func(ctx context.Context, key string) error
}

// Save memeriksa isi file (bukan sekadar nama/ekstensi), menyimpan, dan mencatat pemiliknya.
func (u Uploader) Save(ctx context.Context, ownerID, kind string, r io.Reader) (Uploaded, error) {
	rl, ok := rules[kind]
	if !ok {
		return Uploaded{}, ErrKind
	}
	head := make([]byte, 512)
	n, err := io.ReadFull(r, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return Uploaded{}, err
	}
	head = head[:n]
	ctype := sniff(head)
	ext, ok := rl.types[ctype]
	if !ok {
		return Uploaded{}, ErrType
	}
	maxBytes := rl.maxBytes
	if m, ok := rl.typeMax[ctype]; ok {
		maxBytes = m
	}
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	key := fmt.Sprintf("%s/%s.%s", rl.folder, hex.EncodeToString(b), ext)

	limited := io.LimitReader(io.MultiReader(bytes.NewReader(head), r), maxBytes+1)
	size, err := u.Store.Put(ctx, key, limited)
	if err != nil {
		return Uploaded{}, err
	}
	if size > maxBytes {
		return Uploaded{}, ErrTooLarge // file sisa ditinggal; dibersihkan job harian (uploads tanpa baris)
	}
	if IsVideo(ctype) && u.CheckVideo != nil {
		if err := u.CheckVideo(ctx, key); err != nil {
			return Uploaded{}, err
		}
	}
	if _, err := u.Pool.Exec(ctx, `INSERT INTO uploads (key, owner_user_id, kind, content_type, bytes, is_private) VALUES ($1,$2,$3,$4,$5,$6)`,
		key, ownerID, kind, ctype, size, rl.private); err != nil {
		return Uploaded{}, err
	}
	return Uploaded{Key: key, URL: u.Store.PublicURL(key), ContentType: ctype, Bytes: size}, nil
}

// CheckOwner memastikan key milik user dan jenisnya sesuai. key kosong = tidak diubah.
func CheckOwner(ctx context.Context, pool *pgxpool.Pool, key, ownerID, kind string) error {
	if key == "" {
		return nil
	}
	var ok bool
	if err := pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM uploads WHERE key = $1 AND owner_user_id = $2 AND kind = $3)`,
		key, ownerID, kind).Scan(&ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotOwner
	}
	return nil
}
