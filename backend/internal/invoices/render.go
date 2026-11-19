package invoices

// Faktur siap cetak (A4). Dua bahasa: Khmer + Inggris (dokumen pajak di Kamboja).

import (
	"bytes"
	"html/template"
	"time"
)

var page = template.Must(template.New("inv").Funcs(template.FuncMap{
	"d":     func(t time.Time) string { return t.In(phnomPenh).Format("02 Jan 2006") },
	"month": func(p string) string { t, _ := time.Parse("2006-01", p); return t.Format("January 2006") },
}).Parse(`<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.Number}} · Phsar</title>
<style>
@page{size:A4;margin:16mm}
*{box-sizing:border-box}body{font-family:"Inter","Noto Sans Khmer",system-ui,sans-serif;color:#18213a;font-size:13px;margin:0;background:#f4f4f6}
.sheet{background:#fff;max-width:210mm;margin:16px auto;padding:18mm 16mm;box-shadow:0 2px 12px rgba(0,0,0,.08)}
header{display:flex;justify-content:space-between;align-items:flex-start;gap:24px;border-bottom:3px solid #e8175d;padding-bottom:14px}
.brand{font-size:28px;font-weight:800;color:#e8175d;letter-spacing:-.5px}.brand small{display:block;font-size:11px;color:#697386;font-weight:600;letter-spacing:0}
h1{font-size:20px;margin:0;text-align:right}h1 small{display:block;font-size:12px;color:#697386;font-weight:500}
.meta{text-align:right;margin-top:6px;line-height:1.6}.void{display:inline-block;margin-top:6px;background:#fdecec;color:#b42318;font-weight:800;padding:3px 10px;border-radius:6px}
.parties{display:grid;grid-template-columns:1fr 1fr;gap:24px;margin:18px 0}.parties h3{font-size:11px;text-transform:uppercase;color:#697386;margin:0 0 6px;letter-spacing:.5px}
.parties p{margin:0;line-height:1.55}
table{width:100%;border-collapse:collapse;margin-top:6px}th{text-align:left;font-size:11px;color:#697386;text-transform:uppercase;border-bottom:1px solid #d9dce3;padding:6px 4px}
td{padding:6px 4px;border-bottom:1px solid #eef0f4}td.n,th.n{text-align:right;font-variant-numeric:tabular-nums}tr.ret td{color:#b42318}
.totals{margin-left:auto;width:300px;margin-top:14px}.totals div{display:flex;justify-content:space-between;padding:4px 0}.totals .grand{border-top:2px solid #18213a;font-weight:800;font-size:15px;margin-top:4px;padding-top:8px}
.note{margin-top:22px;color:#697386;font-size:11px;line-height:1.6}
.bar{max-width:210mm;margin:16px auto 0;display:flex;justify-content:flex-end}.bar button{font:inherit;font-weight:700;padding:8px 18px;border-radius:999px;border:0;background:#e8175d;color:#fff;cursor:pointer}
@media print{body{background:#fff}.sheet{box-shadow:none;margin:0;padding:0;max-width:none}.bar{display:none}}
</style></head><body>
<div class="bar"><button type="button" onclick="window.print()">Print / Save PDF</button></div>
<div class="sheet">
<header><div class="brand">Phsar<small>{{.Issuer.LegalName}}</small></div>
<div><h1>វិក្កយបត្រអាករ · Tax invoice<small>Commission for {{month .Period}}</small></h1>
<div class="meta">No. <b>{{.Number}}</b><br>Date: {{d .IssuedAt}}{{if .ReplacesNo}}<br>Replaces {{.ReplacesNo}}{{end}}</div>
{{if eq .Status "void"}}<div class="void">VOID{{if .VoidNote}} — {{.VoidNote}}{{end}}</div>{{end}}</div></header>
<div class="parties">
<div><h3>From / អ្នកលក់សេវា</h3><p><b>{{.Issuer.LegalName}}</b><br>{{.Issuer.Address}}{{if .Issuer.VATTIN}}<br>VAT TIN: {{.Issuer.VATTIN}}{{end}}{{if .Issuer.Email}}<br>{{.Issuer.Email}}{{end}}{{if .Issuer.Phone}}<br>{{.Issuer.Phone}}{{end}}</p></div>
<div><h3>Bill to / អតិថិជន</h3><p><b>{{.Seller.LegalName}}</b>{{if ne .Seller.LegalName .Seller.Name}}<br>Store: {{.Seller.Name}}{{end}}{{if .Seller.Address}}<br>{{.Seller.Address}}{{end}}{{if .Seller.VATTIN}}<br>VAT TIN: {{.Seller.VATTIN}}{{else if .Seller.TIN}}<br>TIN: {{.Seller.TIN}}{{end}}{{if .Seller.Phone}}<br>{{.Seller.Phone}}{{end}}</p></div>
</div>
<table><thead><tr><th>Date</th><th>Reference</th><th>Description</th><th class="n">Sales</th><th class="n">Commission</th></tr></thead><tbody>
{{range .Lines}}<tr{{if eq .Kind "return"}} class="ret"{{end}}><td>{{d .Date}}</td><td>{{.Ref}}</td><td>{{if eq .Kind "return"}}Commission returned · refund of {{.Order}}{{else}}Marketplace commission{{end}}</td><td class="n">${{.Subtotal}}</td><td class="n">${{.Commission}}</td></tr>
{{end}}</tbody></table>
<div class="totals">
<div><span>Orders</span><span>{{.OrderCount}}</span></div>
<div><span>Sales total</span><span>${{.SalesTotal}}</span></div>
<div><span>Net commission</span><span>${{.Commission}}</span></div>
<div><span>Amount before VAT</span><span>${{.Subtotal}}</span></div>
<div><span>VAT {{.VATRate}}%</span><span>${{.VATAmount}}</span></div>
<div class="grand"><span>Total (USD)</span><span>${{.Total}}</span></div>
</div>
<p class="note">{{if .VATInclusive}}Commission amounts include VAT. The commission was already deducted from your payouts — this invoice is for your records; nothing more to pay.{{else}}VAT is charged on top of the commission, which was already deducted from your payouts.{{end}}
Commission returned on refunded returns is shown as a negative line.</p>
</div></body></html>`))

// HTML menghasilkan halaman faktur siap cetak.
func HTML(v *Invoice) ([]byte, error) {
	var b bytes.Buffer
	err := page.Execute(&b, v)
	return b.Bytes(), err
}
