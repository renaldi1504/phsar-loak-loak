// Package chat: chat pembeli ↔ toko.
//
// Satu percakapan per (pembeli, toko). Pesan berisi teks (maks. 2000 karakter) dan/atau satu foto (upload privat
// chat_photo), dan boleh membawa kartu produk / pesanan di atasnya. Di sisi toko semua anggota (pemilik & staf) bisa
// membaca dan membalas. Teks disimpan apa adanya; semua tampilan wajib meng-escape (tidak ada HTML dari pengguna).
//
// Notifikasi pesan baru: ke toko lewat notify.Enqueue (lonceng + Telegram), ke pembeli lewat notify.EnqueueUser —
// maksimal satu per percakapan per 10 menit untuk tiap sisi supaya tidak spam.
package chat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/notify"
	"phsar/internal/platform/db"
	"phsar/internal/storage"
)

const (
	MaxBody     = 2000 // karakter per pesan
	PerMinute   = 20   // pesan per pengguna per menit
	PerDay      = 500  // pesan per pengguna per 24 jam
	NewPerDay   = 30   // percakapan baru per pembeli per 24 jam
	NotifyGap   = 10 * time.Minute
	previewSize = 140
	pageSize    = 200 // pesan per muatan
	PhotoKind   = "chat_photo"
)

var (
	ErrNotFound    = errors.New("percakapan tidak ditemukan")
	ErrRateLimited = errors.New("terlalu banyak pesan, tunggu sebentar")
)

// FieldError: input tidak valid (422).
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

type Service struct {
	Pool  *pgxpool.Pool
	Files storage.Storage
	Log   *slog.Logger
}

// Scope: siapa yang melihat. Pembeli: percakapan miliknya. Toko: semua percakapan toko (anggota sudah dicek router).
type Scope struct {
	UserID  string
	StoreID string // kosong = pembeli
	Lang    string
}

func (sc Scope) side() string {
	if sc.StoreID != "" {
		return "seller"
	}
	return "buyer"
}

func (sc Scope) other() string {
	if sc.StoreID != "" {
		return "buyer"
	}
	return "seller"
}

func (sc Scope) lang() string {
	switch sc.Lang {
	case "en", "km", "zh":
		return sc.Lang
	}
	return "en"
}

// where: batas akses untuk alias c (argumen mulai dari $n).
func (sc Scope) where(n int) (string, []any) {
	if sc.StoreID != "" {
		return fmt.Sprintf("c.store_id::text = $%d", n), []any{sc.StoreID}
	}
	return fmt.Sprintf("c.buyer_id = $%d", n), []any{sc.UserID}
}

// Card: kartu produk / pesanan di atas pesan.
type Card struct {
	Type      string `json:"type"` // product · order
	ProductID string `json:"product_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Slug      string `json:"slug,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	Price     string `json:"price,omitempty"`
	OrderNo   string `json:"order_no,omitempty"`
	Status    string `json:"status,omitempty"`
	Total     string `json:"total,omitempty"`
}

type Message struct {
	ID         int64     `json:"id"`
	Sender     string    `json:"sender"`                // buyer · seller
	SenderName string    `json:"sender_name,omitempty"` // nama staf toko (hanya untuk toko & admin)
	Body       string    `json:"body"`
	HasPhoto   bool      `json:"has_photo"`
	Card       *Card     `json:"card,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	productID  *string
	orderID    *string
}

type Conversation struct {
	ID            string     `json:"id"`
	StoreID       string     `json:"store_id"`
	StoreName     string     `json:"store_name"`
	StoreSlug     string     `json:"store_slug"`
	StoreLogo     string     `json:"store_logo,omitempty"`
	BuyerName     string     `json:"buyer_name,omitempty"` // untuk toko & admin (tanpa nomor HP)
	LastMessageAt *time.Time `json:"last_message_at"`
	LastSender    *string    `json:"last_sender"`
	Preview       string     `json:"preview"`
	Unread        int        `json:"unread"`               // belum dibaca oleh yang melihat
	Waiting       bool       `json:"waiting"`              // pesan terakhir dari pembeli (belum dibalas toko)
	ProductID     *string    `json:"product_id,omitempty"` // kartu produk terakhir
	Messages      []Message  `json:"messages,omitempty"`
	// Hanya admin
	BuyerID      string `json:"buyer_id,omitempty"`
	BuyerPhone   string `json:"buyer_phone,omitempty"`
	BuyerUnread  *int   `json:"buyer_unread,omitempty"`
	SellerUnread *int   `json:"seller_unread,omitempty"`
}

// selectConv: kolom ringkasan; unreadCol = kolom belum dibaca untuk sisi yang melihat.
func selectConv(unreadCol string) string {
	return `SELECT c.id, c.store_id, st.name, st.slug, coalesce(st.logo_key, ''), coalesce(u.name, ''), c.last_message_at, c.last_sender,
	  c.last_preview, ` + unreadCol + `, coalesce(c.last_sender = 'buyer', false), c.product_id, c.buyer_id, coalesce(u.phone, ''), c.buyer_unread, c.seller_unread
	FROM chat_conversations c JOIN stores st ON st.id = c.store_id JOIN users u ON u.id = c.buyer_id`
}

func (s *Service) scanConvs(rows pgx.Rows) ([]Conversation, error) {
	defer rows.Close()
	out := []Conversation{}
	for rows.Next() {
		var c Conversation
		var bu, su int
		if err := rows.Scan(&c.ID, &c.StoreID, &c.StoreName, &c.StoreSlug, &c.StoreLogo, &c.BuyerName, &c.LastMessageAt, &c.LastSender,
			&c.Preview, &c.Unread, &c.Waiting, &c.ProductID, &c.BuyerID, &c.BuyerPhone, &bu, &su); err != nil {
			return nil, err
		}
		c.StoreLogo = s.url(c.StoreLogo)
		c.BuyerUnread, c.SellerUnread = &bu, &su
		out = append(out, c)
	}
	return out, rows.Err()
}

// forUser menyembunyikan data khusus admin; pembeli juga tidak perlu nama dirinya sendiri.
func forUser(c *Conversation, sc Scope) {
	c.BuyerID, c.BuyerPhone, c.BuyerUnread, c.SellerUnread = "", "", nil, nil
	if sc.StoreID == "" {
		c.BuyerName = ""
	}
}

func (s *Service) url(keyOrURL string) string {
	if keyOrURL == "" || strings.HasPrefix(keyOrURL, "http") || strings.HasPrefix(keyOrURL, "/") || s.Files == nil {
		return keyOrURL
	}
	return s.Files.PublicURL(keyOrURL)
}

func unreadCol(sc Scope) string { return "c." + sc.side() + "_unread" }

// ---------------------------------------------------------------------------
// Daftar & isi percakapan
// ---------------------------------------------------------------------------

// List: percakapan yang sudah berisi pesan, terbaru di atas. Mengembalikan juga total belum dibaca & menunggu balasan.
func (s *Service) List(ctx context.Context, sc Scope) ([]Conversation, int, int, error) {
	w, args := sc.where(1)
	rows, err := s.Pool.Query(ctx, selectConv(unreadCol(sc))+" WHERE "+w+" AND c.last_message_at IS NOT NULL ORDER BY c.last_message_at DESC LIMIT 100", args...)
	if err != nil {
		return nil, 0, 0, err
	}
	l, err := s.scanConvs(rows)
	if err != nil {
		return nil, 0, 0, err
	}
	unread, waiting := 0, 0
	for i := range l {
		forUser(&l[i], sc)
		unread += l[i].Unread
		if l[i].Waiting {
			waiting++
		}
	}
	return l, unread, waiting, nil
}

// Get: percakapan + pesan. after > 0 = hanya pesan setelah id itu (polling). Menandai sudah dibaca untuk sisi yang melihat.
func (s *Service) Get(ctx context.Context, sc Scope, id string, after int64) (*Conversation, error) {
	w, args := sc.where(2)
	rows, err := s.Pool.Query(ctx, selectConv(unreadCol(sc))+" WHERE c.id::text = $1 AND "+w, append([]any{id}, args...)...)
	if err != nil {
		return nil, err
	}
	l, err := s.scanConvs(rows)
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, ErrNotFound
	}
	c := &l[0]
	if c.Messages, err = s.messages(ctx, c.ID, after, sc.lang(), sc.StoreID != ""); err != nil {
		return nil, err
	}
	if c.Unread > 0 {
		if _, err := s.Pool.Exec(ctx, `UPDATE chat_conversations SET `+sc.side()+`_unread = 0 WHERE id = $1`, c.ID); err != nil {
			return nil, err
		}
		c.Unread = 0
	}
	forUser(c, sc)
	return c, nil
}

// messages: maksimal pageSize pesan (yang terbaru), urut lama → baru. names = tampilkan nama staf/pembeli.
func (s *Service) messages(ctx context.Context, convID string, after int64, lang string, names bool) ([]Message, error) {
	rows, err := s.Pool.Query(ctx, `SELECT * FROM (
		SELECT m.id, m.sender_type, coalesce(nullif(u.name, ''), ''), m.body, m.photo_key IS NOT NULL, m.created_at, m.product_id::text, m.order_id::text
		FROM chat_messages m JOIN users u ON u.id = m.sender_id
		WHERE m.conversation_id = $1 AND m.id > $2 ORDER BY m.id DESC LIMIT $3) x ORDER BY id`, convID, after, pageSize)
	if err != nil {
		return nil, err
	}
	out, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Message, error) {
		var m Message
		err := r.Scan(&m.ID, &m.Sender, &m.SenderName, &m.Body, &m.HasPhoto, &m.CreatedAt, &m.productID, &m.orderID)
		if !names {
			m.SenderName = ""
		}
		return m, err
	})
	if err != nil {
		return nil, err
	}
	if out == nil {
		out = []Message{}
	}
	return out, s.attachCards(ctx, out, lang)
}

// attachCards mengisi kartu produk / pesanan untuk pesan yang membawanya.
func (s *Service) attachCards(ctx context.Context, ms []Message, lang string) error {
	var pids, oids []string
	for _, m := range ms {
		if m.productID != nil {
			pids = append(pids, *m.productID)
		}
		if m.orderID != nil {
			oids = append(oids, *m.orderID)
		}
	}
	prods, err := s.productCards(ctx, pids, lang)
	if err != nil {
		return err
	}
	orders, err := s.orderCards(ctx, oids)
	if err != nil {
		return err
	}
	for i := range ms {
		if ms[i].productID != nil {
			ms[i].Card = prods[*ms[i].productID]
		} else if ms[i].orderID != nil {
			ms[i].Card = orders[*ms[i].orderID]
		}
	}
	return nil
}

func (s *Service) productCards(ctx context.Context, ids []string, lang string) (map[string]*Card, error) {
	out := map[string]*Card{}
	if len(ids) == 0 {
		return out, nil
	}
	rows, err := s.Pool.Query(ctx, `SELECT p.id::text, coalesce(t.name, ts.name), coalesce(t.slug, ts.slug), coalesce(p.min_price::text, ''),
		  coalesce((SELECT coalesce(m.image_url, m.storage_key) FROM product_media m WHERE m.product_id = p.id AND m.deleted_at IS NULL
		            AND m.kind = 'photo' ORDER BY m.sort_order LIMIT 1), '')
		FROM products p JOIN product_translations ts ON ts.product_id = p.id AND ts.lang = p.source_lang
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = $2::text::lang_code
		WHERE p.id = ANY($1::uuid[])`, ids, lang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		c := &Card{Type: "product"}
		if err := rows.Scan(&c.ProductID, &c.Name, &c.Slug, &c.Price, &c.ImageURL); err != nil {
			return nil, err
		}
		c.ImageURL = s.url(c.ImageURL)
		out[c.ProductID] = c
	}
	return out, rows.Err()
}

func (s *Service) orderCards(ctx context.Context, ids []string) (map[string]*Card, error) {
	out := map[string]*Card{}
	if len(ids) == 0 {
		return out, nil
	}
	rows, err := s.Pool.Query(ctx, `SELECT id::text, order_no, status::text, total::text FROM orders WHERE id = ANY($1::uuid[])`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		c := &Card{Type: "order"}
		if err := rows.Scan(&id, &c.OrderNo, &c.Status, &c.Total); err != nil {
			return nil, err
		}
		out[id] = c
	}
	return out, rows.Err()
}

// ---------------------------------------------------------------------------
// Pembeli membuka percakapan (dari halaman produk / toko / pesanan)
// ---------------------------------------------------------------------------

type OpenInput struct {
	StoreID   string `json:"store_id"`
	ProductID string `json:"product_id"`
	OrderNo   string `json:"order_no"`
}

// Opened: percakapan + kartu yang akan dilampirkan ke pesan pertama pembeli (belum disimpan).
type Opened struct {
	*Conversation
	Card *Card `json:"card,omitempty"`
}

// Open mencari atau membuat percakapan pembeli dengan toko. Toko ditentukan dari store_id, produk, atau pesanan.
func (s *Service) Open(ctx context.Context, sc Scope, in OpenInput) (*Opened, bool, error) {
	storeID := strings.TrimSpace(in.StoreID)
	var card *Card
	viaOrder := false
	if no := strings.ToUpper(strings.TrimSpace(in.OrderNo)); no != "" {
		var oid, sid string
		err := s.Pool.QueryRow(ctx, `SELECT id::text, store_id::text FROM orders WHERE order_no = $1 AND buyer_id = $2`, no, sc.UserID).Scan(&oid, &sid)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, FieldError{"order_no", "order not found in your account"}
		}
		if err != nil {
			return nil, false, err
		}
		if storeID != "" && storeID != sid {
			return nil, false, FieldError{"store_id", "order is from another shop"}
		}
		storeID, viaOrder = sid, true
		cs, err := s.orderCards(ctx, []string{oid})
		if err != nil {
			return nil, false, err
		}
		card = cs[oid]
	}
	if pid := strings.TrimSpace(in.ProductID); pid != "" {
		sid, err := s.productStore(ctx, pid)
		if err != nil {
			return nil, false, err
		}
		if storeID != "" && storeID != sid {
			return nil, false, FieldError{"product_id", "product is from another shop"}
		}
		storeID = sid
		if card == nil {
			cs, err := s.productCards(ctx, []string{pid}, sc.lang())
			if err != nil {
				return nil, false, err
			}
			card = cs[pid]
		}
	}
	if storeID == "" {
		return nil, false, FieldError{"store_id", "choose a shop"}
	}
	var ok, member bool
	if err := s.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM stores WHERE id::text = $1 AND deleted_at IS NULL AND (status = 'active' OR $3)),
		EXISTS (SELECT 1 FROM store_members WHERE store_id::text = $1 AND user_id = $2)`, storeID, sc.UserID, viaOrder).Scan(&ok, &member); err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, FieldError{"store_id", "shop not found"}
	}
	if member {
		return nil, false, FieldError{"store_id", "this is your own shop"}
	}
	var id string
	created := false
	err := db.WithActor(ctx, s.Pool, sc.UserID, "buyer", "", func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `SELECT id::text FROM chat_conversations WHERE buyer_id = $1 AND store_id::text = $2`, sc.UserID, storeID).Scan(&id)
		if err == nil || !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		var recent int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM chat_conversations WHERE buyer_id = $1 AND created_at > now() - interval '24 hours'`, sc.UserID).Scan(&recent); err != nil {
			return err
		}
		if recent >= NewPerDay {
			return ErrRateLimited
		}
		created = true
		return tx.QueryRow(ctx, `INSERT INTO chat_conversations (buyer_id, store_id) VALUES ($1, $2::text::uuid)
			ON CONFLICT (buyer_id, store_id) DO UPDATE SET buyer_id = EXCLUDED.buyer_id RETURNING id::text`, sc.UserID, storeID).Scan(&id)
	})
	if err != nil {
		return nil, false, err
	}
	c, err := s.Get(ctx, sc, id, 0)
	if err != nil {
		return nil, false, err
	}
	return &Opened{Conversation: c, Card: card}, created, nil
}

// productStore: toko pemilik produk (produk yang belum dihapus).
func (s *Service) productStore(ctx context.Context, pid string) (string, error) {
	var sid string
	err := s.Pool.QueryRow(ctx, `SELECT store_id::text FROM products WHERE id::text = $1 AND deleted_at IS NULL`, pid).Scan(&sid)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", FieldError{"product_id", "product not found"}
	}
	return sid, err
}

// ---------------------------------------------------------------------------
// Kirim pesan
// ---------------------------------------------------------------------------

type SendInput struct {
	Body      string `json:"body"`
	PhotoKey  string `json:"photo_key"`
	ProductID string `json:"product_id"` // kartu produk (produk toko ini)
	OrderNo   string `json:"order_no"`   // kartu pesanan (pesanan pembeli ini di toko ini)
}

// cleanBody: rapikan baris baru, buang karakter kontrol (kecuali baris baru & tab), potong spasi di ujung.
func cleanBody(v string) (string, int) {
	v = strings.ReplaceAll(v, "\r\n", "\n")
	v = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' {
			return r
		}
		if unicode.IsControl(r) || r == utf8.RuneError {
			return -1
		}
		return r
	}, v)
	v = strings.TrimSpace(v)
	return v, utf8.RuneCountInString(v)
}

func snippet(body string, photo bool) string {
	s := []rune(strings.Join(strings.Fields(body), " "))
	if len(s) > previewSize {
		s = append(s[:previewSize-1], '…')
	}
	if len(s) == 0 && photo {
		return "📷"
	}
	return string(s)
}

func (s *Service) Send(ctx context.Context, sc Scope, id string, in SendInput) (*Message, error) {
	body, n := cleanBody(in.Body)
	key := strings.TrimSpace(in.PhotoKey)
	if n > MaxBody {
		return nil, FieldError{"body", fmt.Sprintf("at most %d characters", MaxBody)}
	}
	if n == 0 && key == "" {
		return nil, FieldError{"body", "write a message or add a photo"}
	}
	if key != "" {
		if err := storage.CheckOwner(ctx, s.Pool, key, sc.UserID, PhotoKind); err != nil {
			if errors.Is(err, storage.ErrNotOwner) {
				return nil, FieldError{"photo_key", "photo not found"}
			}
			return nil, err
		}
	}
	w, args := sc.where(2)
	var msgID int64
	err := db.WithActor(ctx, s.Pool, sc.UserID, sc.side(), "", func(tx pgx.Tx) error {
		var convID, storeID, buyerID, storeName, buyerName, buyerLang string
		var otherNotified *time.Time
		err := tx.QueryRow(ctx, `SELECT c.id::text, c.store_id::text, c.buyer_id::text, st.name, coalesce(nullif(u.name, ''), ''), coalesce(u.language::text, 'km'),
			  c.`+sc.other()+`_notified_at
			FROM chat_conversations c JOIN stores st ON st.id = c.store_id JOIN users u ON u.id = c.buyer_id
			WHERE c.id::text = $1 AND `+w+` FOR UPDATE OF c`, append([]any{id}, args...)...).
			Scan(&convID, &storeID, &buyerID, &storeName, &buyerName, &buyerLang, &otherNotified)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		var perMin, perDay int
		if err := tx.QueryRow(ctx, `SELECT count(*) FILTER (WHERE created_at > now() - interval '1 minute'), count(*)
			FROM chat_messages WHERE sender_id = $1 AND created_at > now() - interval '24 hours'`, sc.UserID).Scan(&perMin, &perDay); err != nil {
			return err
		}
		if perMin >= PerMinute || perDay >= PerDay {
			return ErrRateLimited
		}
		var productID, orderID *string
		if pid := strings.TrimSpace(in.ProductID); pid != "" {
			var ok bool
			if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM products WHERE id::text = $1 AND store_id::text = $2 AND deleted_at IS NULL)`, pid, storeID).Scan(&ok); err != nil {
				return err
			}
			if !ok {
				return FieldError{"product_id", "product not found in this shop"}
			}
			productID = &pid
		}
		if no := strings.ToUpper(strings.TrimSpace(in.OrderNo)); no != "" {
			var oid string
			err := tx.QueryRow(ctx, `SELECT id::text FROM orders WHERE order_no = $1 AND buyer_id::text = $2 AND store_id::text = $3`, no, buyerID, storeID).Scan(&oid)
			if errors.Is(err, pgx.ErrNoRows) {
				return FieldError{"order_no", "order not found"}
			}
			if err != nil {
				return err
			}
			orderID = &oid
		}
		var photo *string
		if key != "" {
			photo = &key
		}
		if err := tx.QueryRow(ctx, `INSERT INTO chat_messages (conversation_id, sender_type, sender_id, body, photo_key, product_id, order_id)
			VALUES ($1, $2, $3, $4, $5, $6::text::uuid, $7::text::uuid) RETURNING id`, convID, sc.side(), sc.UserID, body, photo, productID, orderID).Scan(&msgID); err != nil {
			return err
		}
		snip := snippet(body, key != "")
		notifyNow := otherNotified == nil || time.Since(*otherNotified) >= NotifyGap
		if _, err := tx.Exec(ctx, `UPDATE chat_conversations SET last_message_at = now(), last_sender = $2, last_preview = $3,
			  `+sc.other()+`_unread = `+sc.other()+`_unread + 1, `+sc.side()+`_unread = 0,
			  product_id = coalesce($4::text::uuid, product_id), order_id = coalesce($5::text::uuid, order_id),
			  `+sc.other()+`_notified_at = CASE WHEN $6 THEN now() ELSE `+sc.other()+`_notified_at END
			WHERE id = $1`, convID, sc.side(), snip, productID, orderID, notifyNow); err != nil {
			return err
		}
		if !notifyNow {
			return nil
		}
		if sc.StoreID == "" {
			who := buyerName
			if who == "" {
				who = "Buyer"
			}
			return notify.Enqueue(ctx, tx, storeID, notify.ChatNew, notify.Params{Person: who, Note: snip}, "#/store/"+storeID+"/chat/"+convID)
		}
		return notify.EnqueueUser(ctx, tx, buyerID, notify.ChatReply, notify.Params{Store: storeName, Note: snip}, "/"+buyerLang+"/chat#"+convID)
	})
	if err != nil {
		return nil, err
	}
	ms, err := s.messages(ctx, id, msgID-1, sc.lang(), sc.StoreID != "")
	if err != nil {
		return nil, err
	}
	if len(ms) == 0 {
		return nil, ErrNotFound
	}
	return &ms[0], nil
}

// Photo membuka foto pesan mid dalam percakapan yang boleh dilihat (sc nil = admin).
func (s *Service) Photo(ctx context.Context, sc *Scope, convID string, mid int64) (io.ReadCloser, string, error) {
	q := `SELECT m.photo_key, u.content_type FROM chat_messages m JOIN chat_conversations c ON c.id = m.conversation_id
		JOIN uploads u ON u.key = m.photo_key WHERE m.id = $1 AND c.id::text = $2`
	args := []any{mid, convID}
	if sc != nil {
		w, a := sc.where(3)
		q += " AND " + w
		args = append(args, a...)
	}
	var key, ctype string
	if err := s.Pool.QueryRow(ctx, q, args...).Scan(&key, &ctype); err != nil {
		return nil, "", ErrNotFound
	}
	f, err := s.Files.Open(ctx, key)
	if err != nil {
		return nil, "", ErrNotFound
	}
	return f, ctype, nil
}

// ---------------------------------------------------------------------------
// Admin (sengketa) — hanya baca, setiap pembukaan percakapan dicatat di audit log
// ---------------------------------------------------------------------------

// AdminList: q = nomor pesanan, HP pembeli, atau nama toko. Tanpa q: percakapan terbaru.
func (s *Service) AdminList(ctx context.Context, q string) ([]Conversation, error) {
	where, args := []string{"c.last_message_at IS NOT NULL"}, []any{}
	if q = strings.TrimSpace(q); q != "" {
		args = append(args, q)
		conds := []string{"st.name ILIKE '%' || $1 || '%'", "u.name ILIKE '%' || $1 || '%'",
			`EXISTS (SELECT 1 FROM orders o WHERE o.order_no = upper($1) AND o.buyer_id = c.buyer_id AND o.store_id = c.store_id)`}
		if d := strings.TrimPrefix(strings.TrimPrefix(digits(q), "855"), "0"); len(d) >= 4 {
			args = append(args, d)
			conds = append(conds, "regexp_replace(coalesce(u.phone, ''), '\\D', '', 'g') LIKE '%' || $2 || '%'")
		}
		where = append(where, "("+strings.Join(conds, " OR ")+")")
	}
	rows, err := s.Pool.Query(ctx, selectConv("0")+" WHERE "+strings.Join(where, " AND ")+" ORDER BY c.last_message_at DESC LIMIT 100", args...)
	if err != nil {
		return nil, err
	}
	return s.scanConvs(rows)
}

func digits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// AdminGet: isi lengkap percakapan (tidak mengubah tanda dibaca) + catatan audit "chat.view".
func (s *Service) AdminGet(ctx context.Context, adminID, id string) (*Conversation, error) {
	rows, err := s.Pool.Query(ctx, selectConv("0")+" WHERE c.id::text = $1", id)
	if err != nil {
		return nil, err
	}
	l, err := s.scanConvs(rows)
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, ErrNotFound
	}
	c := &l[0]
	if c.Messages, err = s.messages(ctx, c.ID, 0, "en", true); err != nil {
		return nil, err
	}
	if _, err := s.Pool.Exec(ctx, `INSERT INTO audit_logs (actor_type, actor_user_id, action, entity_table, entity_id)
		VALUES ('admin', $1, 'chat.view', 'chat_conversations', $2)`, adminID, c.ID); err != nil {
		return nil, err
	}
	return c, nil
}
