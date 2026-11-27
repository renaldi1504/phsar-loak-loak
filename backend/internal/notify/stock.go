package notify

// Stok menipis: dibuat trigger database (migrasi 00052, fungsi phsar_low_stock) untuk semua anggota toko saat
// stok tersedia sebuah varian turun ke/di bawah ambang (ambang varian, kalau tidak ada ambang toko). Sekali per
// penurunan: penanda direset saat stok naik lagi di atas ambang. Grup "store", bunyi "ding" (bawaan).
// Params: product (nama), note (SKU), order (stok tersedia).

// StockLow: jenis notifikasi stok menipis.
const StockLow = "stock.low"

func init() {
	texts[StockLow] = map[string]tmpl{
		"en": {"Stock running low", "{product} ({note}) — {order} left."},
		"km": {"ស្តុកជិតអស់", "{product} ({note}) — នៅសល់ {order}។"},
		"zh": {"库存不足", "{product}（{note}）——仅剩 {order} 件。"},
	}
}
