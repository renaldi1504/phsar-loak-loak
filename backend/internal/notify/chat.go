package notify

// chat: pesan baru dalam chat pembeli ↔ toko (dibatasi maks. 1 notifikasi per percakapan per 10 menit).

const (
	ChatNew   = "chat.new"   // ke anggota toko. Params: Person = nama pembeli, Note = cuplikan
	ChatReply = "chat.reply" // ke pembeli (EnqueueUser). Params: Store = nama toko, Note = cuplikan
)

func init() {
	texts[ChatNew] = map[string]tmpl{
		"en": {"New chat message 💬", "{person}: {note}"},
		"km": {"សារជជែកថ្មី 💬", "{person}៖ {note}"},
		"zh": {"新聊天消息 💬", "{person}：{note}"},
	}
	texts[ChatReply] = map[string]tmpl{
		"en": {"{store} replied 💬", "{note}"},
		"km": {"{store} បានឆ្លើយតប 💬", "{note}"},
		"zh": {"{store} 回复了您 💬", "{note}"},
	}
}
