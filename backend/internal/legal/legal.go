// Package legal: dokumen legal Phsar (Syarat pembeli, Perjanjian seller, Privasi, Retur & refund, Hapus data) dalam 3 bahasa.
// Dokumen berversi dan append-only: menerbitkan versi baru = baris baru, versi lama tetap tersimpan (bukti persetujuan).
// Draf awal (defaults/*.md) dimasukkan otomatis kalau belum ada; super admin mengeditnya di admin → Legal documents.
// Teks memakai Markdown sederhana + placeholder {{company}} {{address}} {{email}} {{phone}} {{telegram}} {{site}} {{base}}.
package legal

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/platform/db"
)

//go:embed defaults/*.md
var defaults embed.FS

// Types: urutan tampil + slug URL (/km/legal/<slug>).
var Types = []struct{ Type, Slug string }{
	{"buyer_terms", "terms"}, {"privacy_policy", "privacy"}, {"return_policy", "returns"},
	{"seller_agreement", "seller-terms"}, {"data_deletion", "data-deletion"},
}

var Langs = []string{"en", "km", "zh"}

func TypeBySlug(slug string) (string, bool) {
	for _, t := range Types {
		if t.Slug == slug {
			return t.Type, true
		}
	}
	return "", false
}

func validType(t string) bool {
	for _, x := range Types {
		if x.Type == t {
			return true
		}
	}
	return false
}

type Store struct{ Pool *pgxpool.Pool }

// EnsureDefaults memasukkan draf awal untuk dokumen/bahasa yang belum punya versi sama sekali.
func (s *Store) EnsureDefaults(ctx context.Context) error {
	for _, t := range Types {
		for _, l := range Langs {
			b, err := defaults.ReadFile("defaults/" + t.Type + "." + l + ".md")
			if err != nil {
				return err
			}
			body := string(b)
			if _, err := s.Pool.Exec(ctx, `INSERT INTO legal_documents (type, version, lang, body, content_hash, effective_at)
				SELECT $1::varchar, '1.0', $2::lang_code, $3, $4, now()
				WHERE NOT EXISTS (SELECT 1 FROM legal_documents WHERE type = $1::varchar AND lang = $2::lang_code)`, t.Type, l, body, hash(body)); err != nil {
				return err
			}
		}
	}
	return nil
}

func hash(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }

type Doc struct {
	Type        string    `json:"type"`
	Version     string    `json:"version"`
	Lang        string    `json:"lang"`
	Body        string    `json:"body"`
	EffectiveAt time.Time `json:"effective_at"`
}

// Latest: versi terbaru yang sudah berlaku. Kalau bahasa itu belum ada, pakai bahasa Inggris.
func (s *Store) Latest(ctx context.Context, typ, lang string) (Doc, error) {
	d := Doc{Type: typ}
	err := s.Pool.QueryRow(ctx, `SELECT version, lang::text, body, effective_at FROM legal_documents
		WHERE type = $1 AND lang IN ($2::lang_code, 'en') AND effective_at <= now()
		ORDER BY (lang = $2::lang_code) DESC, effective_at DESC LIMIT 1`, typ, lang).Scan(&d.Version, &d.Lang, &d.Body, &d.EffectiveAt)
	return d, err
}

type Version struct {
	Version     string    `json:"version"`
	EffectiveAt time.Time `json:"effective_at"`
	By          string    `json:"published_by"`
}

type AdminDoc struct {
	Type     string         `json:"type"`
	Slug     string         `json:"slug"`
	Latest   map[string]Doc `json:"latest"`
	Versions []Version      `json:"versions"`
}

// Admin: semua dokumen dengan versi terbaru per bahasa + riwayat versi.
func (s *Store) Admin(ctx context.Context) ([]AdminDoc, error) {
	out := []AdminDoc{}
	for _, t := range Types {
		d := AdminDoc{Type: t.Type, Slug: t.Slug, Latest: map[string]Doc{}, Versions: []Version{}}
		for _, l := range Langs {
			x, err := s.Latest(ctx, t.Type, l)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
			if x.Lang == l {
				d.Latest[l] = x
			}
		}
		rows, err := s.Pool.Query(ctx, `SELECT l.version, max(l.effective_at), coalesce(max(nullif(u.name, '')), max(u.phone), '')
			FROM legal_documents l LEFT JOIN users u ON u.id = l.published_by WHERE l.type = $1
			GROUP BY l.version ORDER BY max(l.effective_at) DESC LIMIT 20`, t.Type)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var v Version
			if err := rows.Scan(&v.Version, &v.EffectiveAt, &v.By); err != nil {
				rows.Close()
				return nil, err
			}
			d.Versions = append(d.Versions, v)
		}
		rows.Close()
		out = append(out, d)
	}
	return out, nil
}

type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

var versionRe = regexp.MustCompile(`^[0-9A-Za-z.\-]{1,16}$`)

// Publish menerbitkan versi baru untuk ketiga bahasa sekaligus (berlaku sekarang).
func (s *Store) Publish(ctx context.Context, actorID, typ, version string, bodies map[string]string) error {
	if !validType(typ) {
		return FieldError{"type", "unknown document"}
	}
	version = strings.TrimSpace(version)
	if !versionRe.MatchString(version) {
		return FieldError{"version", "e.g. 1.1 or 2026-11"}
	}
	for _, l := range Langs {
		b := strings.TrimSpace(bodies[l])
		if len(b) < 50 || len(b) > 60000 {
			return FieldError{"body_" + l, "text is missing or too long"}
		}
	}
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		var exists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM legal_documents WHERE type = $1 AND version = $2)`, typ, version).Scan(&exists); err != nil {
			return err
		}
		if exists {
			return FieldError{"version", "this version number is already used — use a new one"}
		}
		for _, l := range Langs {
			b := strings.TrimSpace(bodies[l]) + "\n"
			if _, err := tx.Exec(ctx, `INSERT INTO legal_documents (type, version, lang, body, content_hash, effective_at, published_by)
				VALUES ($1, $2, $3::lang_code, $4, $5, now(), $6)`, typ, version, l, b, hash(b), actorID); err != nil {
				return err
			}
		}
		return nil
	})
}

// Settings: identitas perusahaan & kontak yang mengisi placeholder di semua dokumen (dan halaman Kontak).
type Settings struct {
	Company  string `json:"company"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Telegram string `json:"telegram"`
	Facebook string `json:"facebook"`
	Hours    string `json:"hours"`
}

func (s *Store) Settings(ctx context.Context) (Settings, error) {
	st := Settings{Company: "Phsar", Address: "Sihanoukville, Cambodia"} // sisanya diisi super admin di admin → Legal documents
	var raw []byte
	err := s.Pool.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'legal'`).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return st, nil
	}
	if err != nil {
		return st, err
	}
	_ = json.Unmarshal(raw, &st)
	return st, nil
}

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func (s *Store) SaveSettings(ctx context.Context, actorID string, in Settings) error {
	trim := func(v string, n int) string {
		r := []rune(strings.TrimSpace(v))
		if len(r) > n {
			r = r[:n]
		}
		return string(r)
	}
	in = Settings{Company: trim(in.Company, 120), Address: trim(in.Address, 200), Email: trim(in.Email, 120), Phone: trim(in.Phone, 40),
		Telegram: trim(in.Telegram, 60), Facebook: trim(in.Facebook, 200), Hours: trim(in.Hours, 80)}
	if in.Company == "" {
		return FieldError{"company", "required"}
	}
	if in.Email != "" && !emailRe.MatchString(in.Email) {
		return FieldError{"email", "invalid email"}
	}
	if in.Facebook != "" && !strings.HasPrefix(in.Facebook, "https://") {
		return FieldError{"facebook", "must start with https://"}
	}
	v, _ := json.Marshal(in)
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO platform_settings (key, value, updated_by, updated_at) VALUES ('legal', $1, $2, now())
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, v, actorID)
		return err
	})
}

// ---- tampilan ----

var (
	boldRe = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	linkRe = regexp.MustCompile(`\[([^\]]+)\]\(([^)\s]+)\)`)
	olRe   = regexp.MustCompile(`^\d+\.\s+`)
)

// Render: Markdown sederhana → HTML aman (semua teks di-escape dulu; tautan hanya http(s) atau path situs sendiri).
func Render(body string, st Settings, base, site string) template.HTML {
	dash := func(v string) string {
		if v == "" {
			return "—"
		}
		return v
	}
	st.Email, st.Phone, st.Telegram, st.Address = dash(st.Email), dash(st.Phone), dash(st.Telegram), dash(st.Address)
	vars := map[string]string{"company": st.Company, "address": st.Address, "email": st.Email, "phone": st.Phone,
		"telegram": st.Telegram, "site": site, "base": base}
	var b strings.Builder
	list := ""
	closeList := func() {
		if list != "" {
			b.WriteString("</" + list + ">")
			list = ""
		}
	}
	para := []string{}
	flush := func() {
		if len(para) > 0 {
			b.WriteString("<p>" + strings.Join(para, " ") + "</p>")
			para = para[:0]
		}
	}
	for _, raw := range strings.Split(strings.ReplaceAll(body, "\r", ""), "\n") {
		line := inline(strings.TrimSpace(raw), vars)
		switch {
		case line == "":
			flush()
			closeList()
		case strings.HasPrefix(line, "### "):
			flush()
			closeList()
			b.WriteString("<h3>" + line[4:] + "</h3>")
		case strings.HasPrefix(line, "## "):
			flush()
			closeList()
			b.WriteString("<h2>" + line[3:] + "</h2>")
		case strings.HasPrefix(line, "- "):
			flush()
			if list != "ul" {
				closeList()
				b.WriteString("<ul>")
				list = "ul"
			}
			b.WriteString("<li>" + line[2:] + "</li>")
		case olRe.MatchString(line):
			flush()
			if list != "ol" {
				closeList()
				b.WriteString("<ol>")
				list = "ol"
			}
			b.WriteString("<li>" + olRe.ReplaceAllString(line, "") + "</li>")
		default:
			closeList()
			para = append(para, line)
		}
	}
	flush()
	closeList()
	return template.HTML(b.String()) //nolint:gosec — semua teks sudah di-escape di inline()
}

func inline(s string, vars map[string]string) string {
	s = html.EscapeString(s)
	for k, v := range vars {
		s = strings.ReplaceAll(s, "{{"+k+"}}", html.EscapeString(v))
	}
	s = boldRe.ReplaceAllString(s, "<strong>$1</strong>")
	s = linkRe.ReplaceAllStringFunc(s, func(m string) string {
		p := linkRe.FindStringSubmatch(m)
		u := html.UnescapeString(p[2])
		if !(strings.HasPrefix(u, "https://") || strings.HasPrefix(u, "http://") || (strings.HasPrefix(u, "/") && !strings.HasPrefix(u, "//"))) {
			return p[1]
		}
		return fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(u), p[1])
	})
	return s
}
