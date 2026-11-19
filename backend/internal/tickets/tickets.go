// Package tickets: tiket bantuan dari pembeli dan seller ke tim CS Phsar.
//
//	open (menunggu CS) ──CS membalas──► pending (menunggu pengguna) ──pengguna membalas──► open
//	└──diselesaikan──► solved ──pengguna membalas lagi──► open
//
// Catatan internal CS (internal = true) tidak pernah dikirim ke pengguna.
package tickets

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/notify"
	"phsar/internal/platform/db"
)

const (
	MaxPerDay   = 10   // tiket baru per pengguna per 24 jam
	MaxBody     = 4000 // karakter per pesan
	MaxSubject  = 120
	snippetSize = 140
)

var (
	ErrNotFound    = errors.New("tiket tidak ditemukan")
	ErrRateLimited = errors.New("terlalu banyak tiket hari ini, coba lagi besok")
)

// FieldError: input tidak valid (422).
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

var Categories = []string{"order", "payment", "return", "account", "seller", "other"}
var Statuses = []string{"open", "pending", "solved"}

func oneOf(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

type Service struct {
	Pool *pgxpool.Pool
	Log  *slog.Logger
}

// Scope: siapa yang melihat. Pembeli: tiket miliknya (source buyer). Seller: semua tiket toko (anggota toko).
type Scope struct {
	UserID  string
	StoreID string // kosong = pembeli
	Lang    string
}

func (sc Scope) source() string {
	if sc.StoreID != "" {
		return "seller"
	}
	return "buyer"
}

// where: batas akses untuk alias tabel t (argumen mulai dari $n).
func (sc Scope) where(n int) (string, []any) {
	if sc.StoreID != "" {
		return fmt.Sprintf("t.source = 'seller' AND t.store_id = $%d", n), []any{sc.StoreID}
	}
	return fmt.Sprintf("t.source = 'buyer' AND t.requester_id = $%d", n), []any{sc.UserID}
}

type Message struct {
	ID         int64     `json:"id"`
	AuthorType string    `json:"author_type"` // user · admin
	AuthorName string    `json:"author_name"`
	Internal   bool      `json:"internal,omitempty"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
}

type Ticket struct {
	ID                string     `json:"id"`
	No                string     `json:"ticket_no"`
	Source            string     `json:"source"`
	Category          string     `json:"category"`
	Subject           string     `json:"subject"`
	OrderNo           *string    `json:"order_no"`
	Status            string     `json:"status"`
	Waiting           bool       `json:"waiting"` // menunggu jawaban CS
	LastUserMessageAt time.Time  `json:"last_user_message_at"`
	LastAdminReplyAt  *time.Time `json:"last_admin_reply_at"`
	SolvedAt          *time.Time `json:"solved_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	Preview           string     `json:"preview"` // cuplikan pesan terakhir yang boleh dilihat
	MessageCount      int        `json:"message_count"`
	// Hanya admin
	RequesterID    string    `json:"requester_id,omitempty"`
	RequesterName  string    `json:"requester_name,omitempty"`
	RequesterPhone string    `json:"requester_phone,omitempty"`
	StoreID        *string   `json:"store_id,omitempty"`
	StoreName      *string   `json:"store_name,omitempty"`
	AssignedTo     *string   `json:"assigned_to,omitempty"`
	AssignedName   *string   `json:"assigned_name,omitempty"`
	Messages       []Message `json:"messages,omitempty"`
}

// selectTicket: kolom ringkasan. internalOK = admin (cuplikan & jumlah termasuk catatan internal).
func selectTicket(internalOK bool) string {
	vis := "AND NOT m.internal"
	if internalOK {
		vis = ""
	}
	return `SELECT t.id, t.ticket_no, t.source, t.category, t.subject, o.order_no, t.status,
	  t.status = 'open' AND (t.last_admin_reply_at IS NULL OR t.last_user_message_at > t.last_admin_reply_at),
	  t.last_user_message_at, t.last_admin_reply_at, t.solved_at, t.created_at, t.updated_at,
	  coalesce((SELECT left(m.body, ` + fmt.Sprint(snippetSize) + `) FROM ticket_messages m WHERE m.ticket_id = t.id ` + vis + ` ORDER BY m.id DESC LIMIT 1), ''),
	  (SELECT count(*) FROM ticket_messages m WHERE m.ticket_id = t.id ` + vis + `),
	  t.requester_id, coalesce(u.name, ''), coalesce(u.phone, ''), t.store_id, st.name, t.assigned_to, coalesce(nullif(a.name, ''), a.phone)
	FROM tickets t JOIN users u ON u.id = t.requester_id LEFT JOIN orders o ON o.id = t.order_id
	LEFT JOIN stores st ON st.id = t.store_id LEFT JOIN users a ON a.id = t.assigned_to`
}

func scanTickets(rows pgx.Rows) ([]Ticket, error) {
	defer rows.Close()
	out := []Ticket{}
	for rows.Next() {
		var t Ticket
		if err := rows.Scan(&t.ID, &t.No, &t.Source, &t.Category, &t.Subject, &t.OrderNo, &t.Status, &t.Waiting,
			&t.LastUserMessageAt, &t.LastAdminReplyAt, &t.SolvedAt, &t.CreatedAt, &t.UpdatedAt, &t.Preview, &t.MessageCount,
			&t.RequesterID, &t.RequesterName, &t.RequesterPhone, &t.StoreID, &t.StoreName, &t.AssignedTo, &t.AssignedName); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// forUser menyembunyikan data internal admin.
func forUser(t *Ticket) *Ticket {
	t.RequesterID, t.RequesterName, t.RequesterPhone, t.AssignedTo, t.AssignedName = "", "", "", nil, nil
	if t.Source == "buyer" {
		t.StoreID, t.StoreName = nil, nil
	}
	return t
}

func (s *Service) messages(ctx context.Context, ticketID string, internalOK bool) ([]Message, error) {
	q := `SELECT m.id, m.author_type, coalesce(nullif(u.name, ''), u.phone, ''), m.internal, m.body, m.created_at
		FROM ticket_messages m JOIN users u ON u.id = m.author_id WHERE m.ticket_id = $1`
	if !internalOK {
		q += " AND NOT m.internal"
	}
	rows, err := s.Pool.Query(ctx, q+" ORDER BY m.id", ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Message{}
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.AuthorType, &m.AuthorName, &m.Internal, &m.Body, &m.CreatedAt); err != nil {
			return nil, err
		}
		if !internalOK && m.AuthorType == "admin" {
			m.AuthorName = "Phsar Support" // nama staf CS tidak ditampilkan ke pengguna
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func cleanText(v string, max int) (string, int) {
	v = strings.TrimSpace(strings.ReplaceAll(v, "\r\n", "\n"))
	return v, utf8.RuneCountInString(v)
}

// ---------------------------------------------------------------------------
// Pengguna (pembeli & seller)
// ---------------------------------------------------------------------------

func (s *Service) List(ctx context.Context, sc Scope) ([]Ticket, error) {
	w, args := sc.where(1)
	rows, err := s.Pool.Query(ctx, selectTicket(false)+" WHERE "+w+" ORDER BY greatest(t.last_user_message_at, coalesce(t.last_admin_reply_at, t.created_at)) DESC LIMIT 100", args...)
	if err != nil {
		return nil, err
	}
	l, err := scanTickets(rows)
	for i := range l {
		forUser(&l[i])
	}
	return l, err
}

func (s *Service) Get(ctx context.Context, sc Scope, id string) (*Ticket, error) {
	w, args := sc.where(2)
	rows, err := s.Pool.Query(ctx, selectTicket(false)+" WHERE t.id::text = $1 AND "+w, append([]any{id}, args...)...)
	if err != nil {
		return nil, err
	}
	l, err := scanTickets(rows)
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, ErrNotFound
	}
	t := forUser(&l[0])
	if t.Messages, err = s.messages(ctx, t.ID, false); err != nil {
		return nil, err
	}
	return t, nil
}

type CreateInput struct {
	Category string `json:"category"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
	OrderNo  string `json:"order_no"`
}

func (s *Service) Create(ctx context.Context, sc Scope, in CreateInput) (*Ticket, error) {
	subject, sn := cleanText(in.Subject, MaxSubject)
	body, bn := cleanText(in.Message, MaxBody)
	if !oneOf(Categories, in.Category) {
		return nil, FieldError{"category", "choose a topic"}
	}
	if sn < 3 || sn > MaxSubject {
		return nil, FieldError{"subject", "3–120 characters"}
	}
	if bn < 10 || bn > MaxBody {
		return nil, FieldError{"message", "describe your question (10–4000 characters)"}
	}
	var id string
	err := db.WithActor(ctx, s.Pool, sc.UserID, map[bool]string{true: "seller", false: "buyer"}[sc.StoreID != ""], "", func(tx pgx.Tx) error {
		var recent int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM tickets WHERE requester_id = $1 AND created_at > now() - interval '24 hours'`, sc.UserID).Scan(&recent); err != nil {
			return err
		}
		if recent >= MaxPerDay {
			return ErrRateLimited
		}
		var orderID *string
		if no := strings.ToUpper(strings.TrimSpace(in.OrderNo)); no != "" {
			var oid string
			q := `SELECT id FROM orders WHERE order_no = $1 AND buyer_id = $2`
			arg := sc.UserID
			if sc.StoreID != "" {
				q, arg = `SELECT id FROM orders WHERE order_no = $1 AND store_id = $2`, sc.StoreID
			}
			if err := tx.QueryRow(ctx, q, no, arg).Scan(&oid); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return FieldError{"order_no", "order not found in your account"}
				}
				return err
			}
			orderID = &oid
		}
		var store *string
		if sc.StoreID != "" {
			store = &sc.StoreID
		}
		no := "TK" + time.Now().In(phnomPenh).Format("060102") + "-" + randCode(5)
		if err := tx.QueryRow(ctx, `INSERT INTO tickets (ticket_no, source, requester_id, store_id, category, subject, order_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, no, sc.source(), sc.UserID, store, in.Category, subject, orderID).Scan(&id); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `INSERT INTO ticket_messages (ticket_id, author_id, author_type, body) VALUES ($1, $2, 'user', $3)`, id, sc.UserID, body)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, sc, id)
}

// Reply: pengguna membalas → tiket kembali "open" (menunggu CS), juga kalau sudah selesai.
func (s *Service) Reply(ctx context.Context, sc Scope, id, body string) (*Ticket, error) {
	body, n := cleanText(body, MaxBody)
	if n < 1 || n > MaxBody {
		return nil, FieldError{"body", "1–4000 characters"}
	}
	w, args := sc.where(2)
	err := db.WithActor(ctx, s.Pool, sc.UserID, map[bool]string{true: "seller", false: "buyer"}[sc.StoreID != ""], "", func(tx pgx.Tx) error {
		var tid string
		err := tx.QueryRow(ctx, `UPDATE tickets t SET status = 'open', solved_at = NULL, last_user_message_at = now()
			WHERE t.id::text = $1 AND `+w+` RETURNING t.id`, append([]any{id}, args...)...).Scan(&tid)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO ticket_messages (ticket_id, author_id, author_type, body) VALUES ($1, $2, 'user', $3)`, tid, sc.UserID, body)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, sc, id)
}

// Solve: pengguna menandai masalahnya sudah beres.
func (s *Service) Solve(ctx context.Context, sc Scope, id string) (*Ticket, error) {
	w, args := sc.where(2)
	err := db.WithActor(ctx, s.Pool, sc.UserID, map[bool]string{true: "seller", false: "buyer"}[sc.StoreID != ""], "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `UPDATE tickets t SET status = 'solved', solved_at = coalesce(solved_at, now()) WHERE t.id::text = $1 AND `+w, append([]any{id}, args...)...)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, sc, id)
}

// ---------------------------------------------------------------------------
// Admin (CS)
// ---------------------------------------------------------------------------

type Counts struct {
	Open       int `json:"open"`
	Waiting    int `json:"waiting"`
	Pending    int `json:"pending"`
	Solved     int `json:"solved"`
	Mine       int `json:"mine"`
	Unassigned int `json:"unassigned"`
}

// AdminList: filter open · pending · solved · mine · unassigned · all; q = nomor tiket / pesanan / HP / subjek.
func (s *Service) AdminList(ctx context.Context, adminID, filter, q string) ([]Ticket, Counts, error) {
	var c Counts
	if err := s.Pool.QueryRow(ctx, `SELECT count(*) FILTER (WHERE status = 'open'),
		count(*) FILTER (WHERE status = 'open' AND (last_admin_reply_at IS NULL OR last_user_message_at > last_admin_reply_at)),
		count(*) FILTER (WHERE status = 'pending'), count(*) FILTER (WHERE status = 'solved'),
		count(*) FILTER (WHERE status <> 'solved' AND assigned_to = $1), count(*) FILTER (WHERE status <> 'solved' AND assigned_to IS NULL)
		FROM tickets`, adminID).Scan(&c.Open, &c.Waiting, &c.Pending, &c.Solved, &c.Mine, &c.Unassigned); err != nil {
		return nil, c, err
	}
	where, args := []string{}, []any{}
	switch filter {
	case "", "open":
		where = append(where, "t.status = 'open'")
	case "pending", "solved":
		where = append(where, "t.status = '"+filter+"'")
	case "mine":
		args = append(args, adminID)
		where = append(where, fmt.Sprintf("t.status <> 'solved' AND t.assigned_to = $%d", len(args)))
	case "unassigned":
		where = append(where, "t.status <> 'solved' AND t.assigned_to IS NULL")
	case "all":
	default:
		return nil, c, FieldError{"filter", "unknown filter"}
	}
	if q = strings.TrimSpace(q); q != "" {
		args = append(args, q)
		p := fmt.Sprintf("$%d", len(args))
		conds := []string{"t.ticket_no = upper(" + p + ")", "o.order_no = upper(" + p + ")", "t.subject ILIKE '%' || " + p + " || '%'", "st.name ILIKE '%' || " + p + " || '%'"}
		if d := strings.TrimPrefix(strings.TrimPrefix(digits(q), "855"), "0"); len(d) >= 4 {
			args = append(args, d)
			conds = append(conds, fmt.Sprintf("regexp_replace(coalesce(u.phone, ''), '\\D', '', 'g') LIKE '%%' || $%d || '%%'", len(args)))
		}
		where = append(where, "("+strings.Join(conds, " OR ")+")")
	}
	w := "TRUE"
	if len(where) > 0 {
		w = strings.Join(where, " AND ")
	}
	// Yang menunggu CS paling lama di atas; selesai: yang terbaru di atas.
	order := "(t.status = 'open' AND (t.last_admin_reply_at IS NULL OR t.last_user_message_at > t.last_admin_reply_at)) DESC, t.last_user_message_at ASC"
	if filter == "solved" || filter == "all" {
		order = "t.updated_at DESC"
	}
	rows, err := s.Pool.Query(ctx, selectTicket(true)+" WHERE "+w+" ORDER BY "+order+" LIMIT 200", args...)
	if err != nil {
		return nil, c, err
	}
	l, err := scanTickets(rows)
	return l, c, err
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

func (s *Service) AdminGet(ctx context.Context, id string) (*Ticket, error) {
	rows, err := s.Pool.Query(ctx, selectTicket(true)+" WHERE t.id::text = $1", id)
	if err != nil {
		return nil, err
	}
	l, err := scanTickets(rows)
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, ErrNotFound
	}
	t := &l[0]
	if t.Messages, err = s.messages(ctx, t.ID, true); err != nil {
		return nil, err
	}
	return t, nil
}

type AdminReply struct {
	Body     string `json:"body"`
	Internal bool   `json:"internal"`
	Status   string `json:"status"` // kosong: balasan → pending, catatan → tidak berubah
}

// AdminReplyTo: balasan (dikirim ke pengguna + notifikasi) atau catatan internal.
func (s *Service) AdminReplyTo(ctx context.Context, adminID, id string, in AdminReply) (*Ticket, error) {
	body, n := cleanText(in.Body, MaxBody)
	if n < 1 || n > MaxBody {
		return nil, FieldError{"body", "1–4000 characters"}
	}
	if in.Status != "" && !oneOf(Statuses, in.Status) {
		return nil, FieldError{"status", "open, pending or solved"}
	}
	status := in.Status
	if status == "" && !in.Internal {
		status = "pending"
	}
	err := db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		var no, source, requester, lang string
		var store *string
		err := tx.QueryRow(ctx, `SELECT t.ticket_no, t.source, t.requester_id, t.store_id, coalesce(u.language::text, 'km') FROM tickets t
			JOIN users u ON u.id = t.requester_id WHERE t.id::text = $1 FOR UPDATE OF t`, id).Scan(&no, &source, &requester, &store, &lang)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `INSERT INTO ticket_messages (ticket_id, author_id, author_type, internal, body) VALUES ($1, $2, 'admin', $3, $4)`,
			id, adminID, in.Internal, body); err != nil {
			return err
		}
		set := "assigned_to = coalesce(assigned_to, $2)"
		if !in.Internal {
			set += ", last_admin_reply_at = now()"
		}
		args := []any{id, adminID}
		if status != "" {
			args = append(args, status)
			set += ", status = $3::text, solved_at = CASE WHEN $3::text = 'solved' THEN coalesce(solved_at, now()) ELSE NULL END"
		}
		if _, err := tx.Exec(ctx, `UPDATE tickets SET `+set+` WHERE id::text = $1`, args...); err != nil {
			return err
		}
		if in.Internal {
			return nil
		}
		return s.notifyUser(ctx, tx, requester, store, source, lang, id, no, body)
	})
	if err != nil {
		return nil, err
	}
	return s.AdminGet(ctx, id)
}

// notifyUser: pembeli → link website; seller → lonceng Seller Center toko itu (tab Bantuan).
func (s *Service) notifyUser(ctx context.Context, tx pgx.Tx, userID string, store *string, source, lang, id, no, body string) error {
	snip := []rune(strings.Join(strings.Fields(body), " "))
	if len(snip) > snippetSize {
		snip = append(snip[:snippetSize-1], '…')
	}
	p := notify.Params{Order: no, Note: string(snip)}
	if source == "seller" && store != nil {
		params, _ := json.Marshal(p)
		_, err := tx.Exec(ctx, `INSERT INTO notifications (user_id, store_id, kind, params, link) VALUES ($1, $2, $3, $4, $5)`,
			userID, *store, notify.TicketReply, params, "#/store/"+*store+"/help/"+id)
		return err
	}
	return notify.EnqueueUser(ctx, tx, userID, notify.TicketReply, p, "/"+lang+"/support#"+id)
}

func (s *Service) Assign(ctx context.Context, adminID, id string, toMe bool) (*Ticket, error) {
	var to *string
	if toMe {
		to = &adminID
	}
	err := db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `UPDATE tickets SET assigned_to = $2 WHERE id::text = $1`, id, to)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.AdminGet(ctx, id)
}

func (s *Service) SetStatus(ctx context.Context, adminID, id, status string) (*Ticket, error) {
	if !oneOf(Statuses, status) {
		return nil, FieldError{"status", "open, pending or solved"}
	}
	err := db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `UPDATE tickets SET status = $2::text, solved_at = CASE WHEN $2::text = 'solved' THEN coalesce(solved_at, now()) ELSE NULL END
			WHERE id::text = $1`, id, status)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.AdminGet(ctx, id)
}

// ---------------------------------------------------------------------------

const codeAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

func randCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		k, _ := rand.Int(rand.Reader, big.NewInt(int64(len(codeAlphabet))))
		b[i] = codeAlphabet[k.Int64()]
	}
	return string(b)
}

var phnomPenh = func() *time.Location {
	if l, err := time.LoadLocation("Asia/Phnom_Penh"); err == nil {
		return l
	}
	return time.FixedZone("ICT", 7*3600)
}()
