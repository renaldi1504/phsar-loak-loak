package storefront

// chat: teks chat pembeli ↔ toko (/km/chat + tombol "Chat penjual") — digabung ke shopStrs (chat.js lewat shop.js)
// dan strs (judul halaman & tombol di halaman produk / toko).
var chatStrs = map[string]map[string]string{
	"en": {
		"chatTitle": "Chats", "chatSeller": "Chat with seller", "chatNone": "No chats yet. Open a product or a shop and tap “Chat with seller”.",
		"chatBack": "All chats", "chatPh": "Write a message…", "chatSend": "Send", "chatPhoto": "Add photo", "chatYou": "You",
		"chatAbout": "You are asking about:", "chatOrder": "Order {no}", "chatProduct": "Product", "chatStore": "Visit shop",
		"chatHint": "Keep payments on Phsar. Never send money to a seller directly.", "chatTooLong": "At most 2000 characters.",
		"chatUploading": "Uploading photo…", "chatPhotoAlt": "Photo from chat", "chatEmpty": "Say hello and ask your question.",
		"err_chat_rate_limited": "You are sending messages too fast. Please wait a moment.", "err_guest_not_allowed": "Please log in with your phone number to chat.",
	},
	"km": {
		"chatTitle": "ការជជែក", "chatSeller": "ជជែកជាមួយអ្នកលក់", "chatNone": "មិនទាន់មានការជជែកទេ។ បើកផលិតផល ឬហាង ហើយចុច «ជជែកជាមួយអ្នកលក់»។",
		"chatBack": "ការជជែកទាំងអស់", "chatPh": "សរសេរសារ…", "chatSend": "ផ្ញើ", "chatPhoto": "បន្ថែមរូបថត", "chatYou": "អ្នក",
		"chatAbout": "អ្នកកំពុងសួរអំពី៖", "chatOrder": "ការបញ្ជាទិញ {no}", "chatProduct": "ផលិតផល", "chatStore": "ទៅកាន់ហាង",
		"chatHint": "សូមបង់ប្រាក់តាម Phsar ជានិច្ច។ កុំផ្ញើប្រាក់ទៅអ្នកលក់ដោយផ្ទាល់។", "chatTooLong": "អតិបរមា 2000 តួអក្សរ។",
		"chatUploading": "កំពុងផ្ទុករូបថត…", "chatPhotoAlt": "រូបថតពីការជជែក", "chatEmpty": "សួស្ដី ហើយសួរសំណួររបស់អ្នក។",
		"err_chat_rate_limited": "អ្នកផ្ញើសារលឿនពេក។ សូមរង់ចាំបន្តិច។", "err_guest_not_allowed": "សូមចូលដោយលេខទូរស័ព្ទដើម្បីជជែក។",
	},
	"zh": {
		"chatTitle": "聊天", "chatSeller": "联系卖家", "chatNone": "暂无聊天。打开商品或店铺，点击“联系卖家”。",
		"chatBack": "全部聊天", "chatPh": "输入消息…", "chatSend": "发送", "chatPhoto": "添加图片", "chatYou": "你",
		"chatAbout": "你正在咨询：", "chatOrder": "订单 {no}", "chatProduct": "商品", "chatStore": "进入店铺",
		"chatHint": "请始终通过 Phsar 付款，切勿直接向卖家转账。", "chatTooLong": "最多 2000 个字符。",
		"chatUploading": "正在上传图片…", "chatPhotoAlt": "聊天图片", "chatEmpty": "打个招呼，提出你的问题吧。",
		"err_chat_rate_limited": "发送太快了，请稍等片刻。", "err_guest_not_allowed": "请使用手机号登录后再聊天。",
	},
}

func init() {
	for l, m := range chatStrs {
		for k, v := range m {
			shopStrs[l][k] = v
		}
		strs[l]["chatTitle"], strs[l]["chatSeller"] = m["chatTitle"], m["chatSeller"]
	}
}
