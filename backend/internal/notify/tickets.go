package notify

// support: notifikasi tiket bantuan (balasan CS ke pembeli / seller yang membuka tiket).

// TicketReply: CS Phsar membalas tiket. Params: Order = nomor tiket, Note = cuplikan balasan.
const TicketReply = "ticket.reply"

func init() {
	texts[TicketReply] = map[string]tmpl{
		"en": {"Phsar Support replied 💬", "Ticket {order}: {note}"},
		"km": {"ជំនួយ Phsar បានឆ្លើយតប 💬", "សំបុត្រ {order}៖ {note}"},
		"zh": {"Phsar 客服已回复 💬", "工单 {order}：{note}"},
	}
}
