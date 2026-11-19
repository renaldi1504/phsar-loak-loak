package storefront

// support: teks halaman tiket bantuan (/km/support) — digabung ke shopStrs (shop.js) dan strs (judul halaman).
var ticketStrs = map[string]map[string]string{
	"en": {
		"supportTitle": "Help & support", "supportOpen": "Send us a message", "tkSub": "Ask Phsar about an order, payment, return or your account. We usually reply within one working day.",
		"tkNew": "New request", "tkNone": "You have no support requests yet.", "tkTopic": "Topic", "tkSubject": "Subject", "tkSubjectPh": "Short summary, e.g. Payment not confirmed",
		"tkOrder": "Order (optional)", "tkNoOrder": "Not about an order", "tkMessage": "Message", "tkMessagePh": "Tell us what happened. Include anything that helps us check quickly.",
		"tkSend": "Send", "tkReply": "Reply", "tkReplyPh": "Write a reply…", "tkSolve": "My problem is solved", "tkSolved": "Marked as solved — thank you!",
		"tkSent": "Message sent. We will reply here and notify you.", "tkBack": "All requests", "tkYou": "You", "tkReopen": "Write again to reopen this request.",
		"tkCat_order": "Order", "tkCat_payment": "Payment", "tkCat_return": "Return / refund", "tkCat_account": "Account", "tkCat_seller": "Selling on Phsar", "tkCat_other": "Other",
		"tkSt_open": "Waiting for Phsar", "tkSt_pending": "Phsar replied", "tkSt_solved": "Solved", "err_rate_limited": "Too many requests today. Please try again tomorrow.",
	},
	"km": {
		"supportTitle": "ជំនួយ", "supportOpen": "ផ្ញើសារមកយើង", "tkSub": "សួរ Phsar អំពីការបញ្ជាទិញ ការបង់ប្រាក់ ការប្ដូរ ឬគណនីរបស់អ្នក។ ជាធម្មតាយើងឆ្លើយក្នុងមួយថ្ងៃធ្វើការ។",
		"tkNew": "សំណើថ្មី", "tkNone": "អ្នកមិនទាន់មានសំណើជំនួយទេ។", "tkTopic": "ប្រធានបទ", "tkSubject": "ចំណងជើង", "tkSubjectPh": "សង្ខេបខ្លី ឧ. ការបង់ប្រាក់មិនទាន់បញ្ជាក់",
		"tkOrder": "ការបញ្ជាទិញ (ស្រេចចិត្ត)", "tkNoOrder": "មិនទាក់ទងនឹងការបញ្ជាទិញ", "tkMessage": "សារ", "tkMessagePh": "ប្រាប់យើងពីអ្វីដែលបានកើតឡើង។ បញ្ចូលព័ត៌មានដែលជួយយើងពិនិត្យបានលឿន។",
		"tkSend": "ផ្ញើ", "tkReply": "ឆ្លើយតប", "tkReplyPh": "សរសេរការឆ្លើយតប…", "tkSolve": "បញ្ហារបស់ខ្ញុំត្រូវបានដោះស្រាយ", "tkSolved": "បានសម្គាល់ថាដោះស្រាយរួច — សូមអរគុណ!",
		"tkSent": "បានផ្ញើសារ។ យើងនឹងឆ្លើយនៅទីនេះ ហើយជូនដំណឹងអ្នក។", "tkBack": "សំណើទាំងអស់", "tkYou": "អ្នក", "tkReopen": "សរសេរម្តងទៀតដើម្បីបើកសំណើនេះឡើងវិញ។",
		"tkCat_order": "ការបញ្ជាទិញ", "tkCat_payment": "ការបង់ប្រាក់", "tkCat_return": "ការប្ដូរ / សងប្រាក់", "tkCat_account": "គណនី", "tkCat_seller": "ការលក់នៅលើ Phsar", "tkCat_other": "ផ្សេងៗ",
		"tkSt_open": "រង់ចាំ Phsar", "tkSt_pending": "Phsar បានឆ្លើយ", "tkSt_solved": "បានដោះស្រាយ", "err_rate_limited": "សំណើច្រើនពេកថ្ងៃនេះ។ សូមព្យាយាមម្តងទៀតថ្ងៃស្អែក។",
	},
	"zh": {
		"supportTitle": "帮助与支持", "supportOpen": "给我们留言", "tkSub": "向 Phsar 咨询订单、付款、退货或账户问题。我们通常在一个工作日内回复。",
		"tkNew": "新建请求", "tkNone": "您还没有支持请求。", "tkTopic": "主题", "tkSubject": "标题", "tkSubjectPh": "简短概括，例如：付款未确认",
		"tkOrder": "订单（可选）", "tkNoOrder": "与订单无关", "tkMessage": "内容", "tkMessagePh": "请描述发生了什么，并提供有助于我们快速核实的信息。",
		"tkSend": "发送", "tkReply": "回复", "tkReplyPh": "写下回复…", "tkSolve": "我的问题已解决", "tkSolved": "已标记为解决——谢谢！",
		"tkSent": "消息已发送，我们会在这里回复并通知您。", "tkBack": "全部请求", "tkYou": "您", "tkReopen": "再次留言即可重新打开此请求。",
		"tkCat_order": "订单", "tkCat_payment": "付款", "tkCat_return": "退货 / 退款", "tkCat_account": "账户", "tkCat_seller": "在 Phsar 开店", "tkCat_other": "其他",
		"tkSt_open": "等待 Phsar 回复", "tkSt_pending": "Phsar 已回复", "tkSt_solved": "已解决", "err_rate_limited": "今天的请求过多，请明天再试。",
	},
}

func init() {
	for l, m := range ticketStrs {
		for k, v := range m {
			shopStrs[l][k] = v
		}
		strs[l]["supportTitle"], strs[l]["supportOpen"] = m["supportTitle"], m["supportOpen"]
	}
}
