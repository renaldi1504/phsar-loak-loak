package qa

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"math"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"phsar/internal/httpx"
)

// fakeEmbedder meniru embedder/server.py: vektor = warna rata-rata gambar (cukup untuk menguji alur).
func fakeEmbedder() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/embed" {
			http.NotFound(w, r)
			return
		}
		img, _, err := image.Decode(io.LimitReader(r.Body, 12<<20))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"error":"bad_image"}`))
			return
		}
		b := img.Bounds()
		var sr, sg, sb, n float64
		for y := b.Min.Y; y < b.Max.Y; y += 4 {
			for x := b.Min.X; x < b.Max.X; x += 4 {
				r, g, bl, _ := img.At(x, y).RGBA()
				sr, sg, sb, n = sr+float64(r)/65535, sg+float64(g)/65535, sb+float64(bl)/65535, n+1
			}
		}
		v := []float64{sr/n - .5, sg/n - .5, sb/n - .5, .3}
		norm := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2] + v[3]*v[3])
		for i := range v {
			v[i] /= norm
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"model": "fake-color", "vector": v})
	}))
}

func solidJPEG(c color.RGBA) []byte {
	img := image.NewRGBA(image.Rect(0, 0, 320, 320))
	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = c.R, c.G, c.B, 255
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	return buf.Bytes()
}

func (e *env) photoSearch(t *testing.T, data []byte, ip string) (int, string, string) {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, _ := mw.CreateFormFile("image", "photo.jpg")
	_, _ = fw.Write(data)
	_ = mw.Close()
	req, _ := http.NewRequest(http.MethodPost, e.srv.URL+"/en/photo-search", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("CF-Connecting-IP", ip)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return res.StatusCode, res.Header.Get("X-Results"), string(body)
}

// ---------------------------------------------------------------------------
// 14. Cari dengan foto
// ---------------------------------------------------------------------------

func TestPhotoSearch(t *testing.T) {
	e := getEnv(t)
	ctx := context.Background()
	mod := newAdmin(t, e, "moderator")
	tok, sid := newActiveStore(t, e)
	mk := func(name string, c color.RGBA) string {
		v := e.upload(t, tok, "product_media", mp4Bytes).want(201)
		f := e.upload(t, tok, "product_media", solidJPEG(c)).want(201)
		p := e.call(t, "POST", "/v1/seller/stores/"+sid+"/products", map[string]any{
			"category_id": "womens-tops", "source_lang": "en", "name": name, "description": "Photo search test.", "condition": "new",
			"variants": []any{map[string]any{"price": "7.50", "stock_qty": 4}}, "submit": true,
			"media": []any{map[string]any{"key": v.str("key"), "kind": "video"}, map[string]any{"key": f.str("key"), "kind": "photo"}},
		}, tok).want(201)
		e.call(t, "POST", "/admin/api/products/"+p.str("id")+"/decision", map[string]string{"decision": "approve"}, mod).want(200)
		return p.str("id")
	}
	red := mk("Crimson canvas tote", color.RGBA{220, 20, 30, 255})
	mk("Ocean blue cap", color.RGBA{20, 60, 220, 255})
	svc := httpx.NewPhotoSearch(e.cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), e.pool)
	for i := 0; i < 20; i++ {
		n, err := svc.IndexPending(ctx, 50)
		if err != nil {
			t.Fatalf("IndexPending: %v", err)
		}
		if n == 0 {
			break
		}
	}
	check(t, e.scalar(t, `SELECT count(*) FROM product_image_vectors WHERE product_id = $1 AND vec IS NOT NULL`, red) == "1", "foto produk tidak diindeks")

	t.Run("foto merah → produk merah, bukan yang biru", func(t *testing.T) {
		code, n, body := e.photoSearch(t, solidJPEG(color.RGBA{200, 30, 40, 255}), "10.9.0.1")
		check(t, code == 200, "status %d: %s", code, body)
		check(t, n != "" && n != "0", "tidak ada hasil (X-Results=%q)", n)
		check(t, strings.Contains(body, "Crimson canvas tote"), "produk merah tidak ditemukan: %s", body)
		check(t, !strings.Contains(body, "Ocean blue cap"), "produk biru ikut muncul")
		check(t, strings.Contains(body, `data-quick-add=`), "kartu tanpa tombol keranjang")
	})
	t.Run("produk dihapus tidak muncul", func(t *testing.T) {
		e.sql(t, `UPDATE products SET deleted_at = now() WHERE id = $1`, red)
		code, n, body := e.photoSearch(t, solidJPEG(color.RGBA{200, 30, 40, 255}), "10.9.0.1")
		check(t, code == 200 && n == "0" && !strings.Contains(body, "Crimson"), "produk terhapus muncul: %d %s", code, n)
	})
	t.Run("bukan gambar / terlalu besar / kosong", func(t *testing.T) {
		code, _, _ := e.photoSearch(t, []byte("hello, this is text not an image"), "10.9.0.2")
		check(t, code == 415, "teks harus 415, dapat %d", code)
		code, _, _ = e.photoSearch(t, append(solidJPEG(color.RGBA{1, 2, 3, 255}), make([]byte, 11<<20)...), "10.9.0.2")
		check(t, code == 413, "file >10MB harus 413, dapat %d", code)
		code, _, _ = e.photoSearch(t, []byte{0xff, 0xd8, 0xff, 0xe0, 0, 0, 0, 0}, "10.9.0.2")
		check(t, code == 422, "jpeg rusak harus 422, dapat %d", code)
	})
	t.Run("batas 12 pencarian per menit per IP", func(t *testing.T) {
		img := solidJPEG(color.RGBA{9, 200, 9, 255})
		last := 0
		for i := 0; i < 13; i++ {
			last, _, _ = e.photoSearch(t, img, "10.9.0.3")
		}
		check(t, last == 429, "pencarian ke-13 harus 429, dapat %d", last)
		code, _, _ := e.photoSearch(t, img, "10.9.0.4")
		check(t, code == 200, "IP lain ikut dibatasi: %d", code)
	})
	t.Run("halaman memuat tombol cari foto", func(t *testing.T) {
		r := e.page(t, "/en/")
		body := string(r.Raw)
		check(t, strings.Contains(body, `data-photo`) && strings.Contains(body, `id="photoInput"`), "tombol / input foto tidak ada")
	})
}
