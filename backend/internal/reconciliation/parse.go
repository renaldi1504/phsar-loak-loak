package reconciliation

// Membaca ekspor transaksi ABA PayWay (CSV). Format ekspor bisa berubah, jadi kolom dikenali dari nama header,
// bukan posisinya: BOM, baris judul sebelum header, pemisah , ; atau tab, tanda kutip, pemisah ribuan, "$", (1.00) negatif.

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const (
	MaxBytes = 5 << 20
	MaxRows  = 20000
)

// Line: satu transaksi dari file.
type Line struct {
	LineNo    int
	TranID    string
	Reference string
	Status    string
	Kind      string // payment · refund
	Amount    decimal.Decimal
	HasAmount bool
	Currency  string
	Fee       *decimal.Decimal
	At        *time.Time
	Raw       map[string]string
}

var (
	ErrNoHeader = errors.New("kolom ID transaksi (tran_id / Transaction ID) tidak ditemukan di 10 baris pertama")
	ErrNotText  = errors.New("file bukan CSV teks")
	ErrTooMany  = fmt.Errorf("maksimal %d baris per file", MaxRows)
)

// alias header (huruf kecil, hanya huruf/angka) → kolom
var aliases = map[string]string{
	"tranid": "tran", "transactionid": "tran", "transid": "tran", "txnid": "tran", "merchanttranid": "tran",
	"merchanttransactionid": "tran", "paywaytranid": "tran",
	"originaltranid": "orig", "originaltransactionid": "orig", "reftranid": "orig",
	"amount": "amount", "totalamount": "amount", "transactionamount": "amount", "paidamount": "amount", "amountusd": "amount",
	"grossamount": "amount", "refundamount": "amount",
	"currency": "currency", "ccy": "currency",
	"status": "status", "paymentstatus": "status", "transactionstatus": "status",
	"date": "date", "datetime": "date", "transactiondate": "date", "transactiondatetime": "date", "createdat": "date",
	"paymentdate": "date", "transactiontime": "date", "trandate": "date", "createddate": "date",
	"time": "time",
	"apv":  "ref", "apvcode": "ref", "reference": "ref", "ref": "ref", "approvalcode": "ref", "refundreference": "ref",
	"refundid": "ref", "referenceno": "ref", "refundtranid": "ref",
	"type": "type", "transactiontype": "type", "trantype": "type",
	"fee": "fee", "fees": "fee", "mdr": "fee", "servicefee": "fee", "transactionfee": "fee", "commission": "fee",
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]`)

func normHeader(h string) string {
	return nonAlnum.ReplaceAllString(strings.ToLower(strings.TrimSpace(h)), "")
}

// Parse membaca seluruh file. Baris kosong dilewati; baris tanpa ID transaksi diabaikan.
func Parse(data []byte) ([]Line, error) {
	if len(data) > MaxBytes {
		return nil, fmt.Errorf("file lebih dari %d MB", MaxBytes>>20)
	}
	data = bytes.TrimPrefix(data, []byte("\xEF\xBB\xBF"))
	if bytes.IndexByte(data, 0) >= 0 {
		return nil, ErrNotText
	}
	// pemisah: yang paling sering muncul di baris-baris awal
	head := string(data[:min(len(data), 4096)])
	sep := ','
	best := strings.Count(head, ",")
	for _, c := range []rune{';', '\t'} {
		if n := strings.Count(head, string(c)); n > best {
			sep, best = c, n
		}
	}
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma, r.LazyQuotes, r.FieldsPerRecord, r.TrimLeadingSpace = sep, true, -1, true
	var cols map[string]int
	var header []string
	var out []Line
	line := 0
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("baris %d: %v", line+1, err)
		}
		line++
		if cols == nil {
			if line > 10 {
				return nil, ErrNoHeader
			}
			c := map[string]int{}
			for i, h := range rec {
				if k, ok := aliases[normHeader(h)]; ok {
					if _, dup := c[k]; !dup {
						c[k] = i
					}
				}
			}
			if _, ok := c["tran"]; ok {
				cols, header = c, rec
			}
			continue
		}
		get := func(k string) string {
			if i, ok := cols[k]; ok && i < len(rec) {
				return strings.TrimSpace(rec[i])
			}
			return ""
		}
		tran := get("tran")
		if tran == "" {
			continue
		}
		if len(out) >= MaxRows {
			return nil, ErrTooMany
		}
		l := Line{LineNo: line, TranID: cut(tran, 64), Reference: cut(get("ref"), 80), Status: cut(get("status"), 40),
			Currency: strings.ToUpper(cut(get("currency"), 3)), Kind: "payment", Raw: map[string]string{}}
		for i, h := range header {
			if i < len(rec) && strings.TrimSpace(h) != "" && rec[i] != "" {
				l.Raw[cut(strings.TrimSpace(h), 60)] = cut(rec[i], 200)
			}
		}
		if a, ok := parseAmount(get("amount")); ok {
			l.Amount, l.HasAmount = a, true
		}
		if f, ok := parseAmount(get("fee")); ok {
			f = f.Abs()
			l.Fee = &f
		}
		typ := strings.ToLower(get("type") + " " + l.Status)
		if strings.Contains(typ, "refund") || strings.Contains(typ, "reversal") || l.Amount.IsNegative() {
			l.Kind = "refund"
			l.Amount = l.Amount.Abs()
			// baris refund: kolom "original tran id" (kalau ada) adalah transaksi asal; ID baris menjadi referensi refund
			if o := get("orig"); o != "" {
				if l.Reference == "" {
					l.Reference = cut(tran, 80)
				}
				l.TranID = cut(o, 64)
			}
		}
		if t, ok := parseTime(strings.TrimSpace(get("date") + " " + get("time"))); ok {
			l.At = &t
		}
		out = append(out, l)
	}
	if cols == nil {
		return nil, ErrNoHeader
	}
	return out, nil
}

func cut(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}

var decimalComma = regexp.MustCompile(`^-?\d+,\d{1,2}$`)

// parseAmount: "1,234.50" · "$12.00" · "12.00 USD" · "(1.00)" · "-1.00" · "12,50" (koma desimal)
func parseAmount(s string) (decimal.Decimal, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return decimal.Zero, false
	}
	neg := false
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		neg, s = true, s[1:len(s)-1]
	}
	s = strings.NewReplacer("$", "", "USD", "", "usd", "", "KHR", "", " ", "", " ", "").Replace(s)
	switch {
	case strings.Contains(s, ",") && strings.Contains(s, "."):
		// yang muncul terakhir adalah pemisah desimal: 1,234.50 atau 1.234,50
		if strings.LastIndex(s, ",") > strings.LastIndex(s, ".") {
			s = strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", ".")
		} else {
			s = strings.ReplaceAll(s, ",", "")
		}
	case decimalComma.MatchString(s):
		s = strings.ReplaceAll(s, ",", ".")
	default:
		s = strings.ReplaceAll(s, ",", "")
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero, false
	}
	if neg {
		d = d.Neg()
	}
	return d, true
}

var Zone = func() *time.Location {
	if l, err := time.LoadLocation("Asia/Phnom_Penh"); err == nil {
		return l
	}
	return time.FixedZone("ICT", 7*3600)
}()

// Kamboja menulis tanggal hari/bulan/tahun.
var layouts = []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02 15:04", "2006-01-02",
	"02/01/2006 15:04:05", "02/01/2006 15:04", "02/01/2006", "02-01-2006 15:04:05", "02-01-2006 15:04", "02-01-2006",
	"02-Jan-2006 15:04:05", "02-Jan-2006 15:04", "02-Jan-2006", "02 Jan 2006 15:04:05", "02 Jan 2006 15:04", "02 Jan 2006",
	"Jan 02, 2006 03:04:05 PM", "Jan 02, 2006 03:04 PM", "Jan 2, 2006 3:04 PM", "Jan 02, 2006", "2006/01/02 15:04:05", "2006/01/02"}

func parseTime(s string) (time.Time, bool) {
	s = strings.Join(strings.Fields(s), " ")
	if s == "" {
		return time.Time{}, false
	}
	for _, l := range layouts {
		if t, err := time.ParseInLocation(l, s, Zone); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
