package reconciliation

import "testing"

func TestParseCommaExport(t *testing.T) {
	csv := "Transaction ID,Transaction Date,Amount,Currency,Status,Transaction Type,Original Tran ID\r\n" +
		"T1,2026-11-16 09:30:00,\"1,234.50\",usd,APPROVED,Purchase,\r\n" +
		"R9,16/11/2026 10:00,$3.00,USD,SUCCESS,Refund,T1\r\n" +
		",,,,,,\r\n"
	l, err := Parse([]byte(csv))
	if err != nil || len(l) != 2 {
		t.Fatalf("%v %d", err, len(l))
	}
	if l[0].Amount.String() != "1234.5" || l[0].Currency != "USD" || l[0].At == nil || l[0].Kind != "payment" {
		t.Fatalf("baris 1: %+v", l[0])
	}
	if l[1].Kind != "refund" || l[1].TranID != "T1" || l[1].Reference != "R9" || l[1].Amount.String() != "3" {
		t.Fatalf("baris refund: %+v", l[1])
	}
	for in, want := range map[string]string{"(2.00)": "-2", "1.234,50": "1234.5", "12,50": "12.5", "1,000": "1000", "7 USD": "7"} {
		if d, ok := parseAmount(in); !ok || d.String() != want {
			t.Errorf("%s → %s", in, d)
		}
	}
	if _, err := Parse([]byte("a,b\n1,2\n")); err != ErrNoHeader {
		t.Fatalf("tanpa kolom ID harus ErrNoHeader: %v", err)
	}
}
