package storefront

// Teks retur untuk shop.js (digabung ke shopStrs).
var returnStrs = map[string]map[string]string{
	"en": {
		"back": "Back", "rtTitle": "Return or refund", "rtAsk": "Request a return", "rtUntil": "Problems with your order? You can ask until {d}.",
		"rtWhich": "Which items?", "rtReason": "What happened?", "rtDesc": "Describe the problem", "rtDescPh": "For example: the screen has a crack on the left side",
		"rtPhotos": "Photos (strongly recommended for damaged or wrong items)", "rtAddPhoto": "Add photos", "rtSend": "Send request", "rtCancelReq": "Cancel request",
		"rtConfirmCancel": "Cancel this return request?", "rtPolicy": "Read the return policy", "rtNo": "Return {no}", "rtAmount": "Requested {a}",
		"rtReason_changed_mind": "Changed my mind", "rtReason_wrong_item": "Wrong item", "rtReason_damaged": "Damaged or broken",
		"rtReason_not_as_described": "Not as shown in the video", "rtReason_missing_parts": "Missing parts", "rtReason_other": "Other problem",
		"rtNoChangeMind": "Change of mind is no longer possible for these items — only problems with the item.",
		"rtSt_requested": "Waiting for the seller", "rtSt_disputed": "Phsar is reviewing", "rtSt_approved": "Approved", "rtSt_rejected": "Not approved",
		"rtSt_refunded": "Refunded", "rtSt_cancelled": "Cancelled",
		"rtInfo_requested": "The seller has until {d} to answer. If they do not, Phsar decides.",
		"rtInfo_disputed":  "Phsar is reviewing your request and usually decides within 2 working days.",
		"rtInfo_approved":  "We will refund {a} to the account that paid, usually within 7 working days.",
		"rtPickup":         "A Phsar courier will contact you to collect the item for free. Keep the item and its packaging until then.",
		"rtInfo_refunded":  "{a} was refunded on {d}.", "rtInfo_cancelled": "You cancelled this request.",
		"rtSellerSaid": "Seller", "rtPhsarSaid": "Phsar", "rtSent": "Return request sent", "rtTooMany": "At most 6 photos", "rtUploading": "Uploading photos…",
		"ev_return_requested": "Return requested", "ev_return_approved": "Return approved", "ev_return_disputed": "Return sent to Phsar review",
		"ev_return_rejected": "Return not approved", "ev_return_refunded": "Return refunded", "ev_return_cancelled": "Return cancelled", "ev_completed": "Completed",
	},
	"km": {
		"back": "ត្រឡប់", "rtTitle": "ប្ដូរ ឬសងប្រាក់វិញ", "rtAsk": "ស្នើប្ដូរទំនិញ", "rtUntil": "មានបញ្ហាជាមួយការបញ្ជាទិញ? អ្នកអាចស្នើបានរហូតដល់ {d}។",
		"rtWhich": "ទំនិញណាខ្លះ?", "rtReason": "មានអ្វីកើតឡើង?", "rtDesc": "ពិពណ៌នាបញ្ហា", "rtDescPh": "ឧទាហរណ៍៖ អេក្រង់មានស្នាមប្រេះនៅខាងឆ្វេង",
		"rtPhotos": "រូបថត (ណែនាំយ៉ាងខ្លាំងសម្រាប់ទំនិញខូច ឬខុស)", "rtAddPhoto": "បន្ថែមរូបថត", "rtSend": "ផ្ញើសំណើ", "rtCancelReq": "បោះបង់សំណើ",
		"rtConfirmCancel": "បោះបង់សំណើប្ដូរនេះ?", "rtPolicy": "អានគោលការណ៍ប្ដូរទំនិញ", "rtNo": "ការប្ដូរ {no}", "rtAmount": "ស្នើ {a}",
		"rtReason_changed_mind": "ប្ដូរចិត្ត", "rtReason_wrong_item": "ទំនិញខុស", "rtReason_damaged": "ខូច ឬបែក",
		"rtReason_not_as_described": "មិនដូចក្នុងវីដេអូ", "rtReason_missing_parts": "ខ្វះគ្រឿង", "rtReason_other": "បញ្ហាផ្សេងទៀត",
		"rtNoChangeMind": "មិនអាចប្ដូរដោយសារប្ដូរចិត្តបានទៀតទេ — មានតែបញ្ហាជាមួយទំនិញប៉ុណ្ណោះ។",
		"rtSt_requested": "រង់ចាំអ្នកលក់", "rtSt_disputed": "Phsar កំពុងពិនិត្យ", "rtSt_approved": "បានយល់ព្រម", "rtSt_rejected": "មិនយល់ព្រម",
		"rtSt_refunded": "បានសងប្រាក់វិញ", "rtSt_cancelled": "បានបោះបង់",
		"rtInfo_requested": "អ្នកលក់មានពេលដល់ {d} ដើម្បីឆ្លើយ។ បើមិនឆ្លើយ Phsar នឹងសម្រេច។",
		"rtInfo_disputed":  "Phsar កំពុងពិនិត្យសំណើរបស់អ្នក ហើយជាធម្មតាសម្រេចក្នុង 2 ថ្ងៃធ្វើការ។",
		"rtInfo_approved":  "យើងនឹងសងប្រាក់ {a} ទៅគណនីដែលបានបង់ ជាធម្មតាក្នុង 7 ថ្ងៃធ្វើការ។",
		"rtPickup":         "អ្នកដឹក Phsar នឹងទាក់ទងអ្នកដើម្បីមកយកទំនិញដោយឥតគិតថ្លៃ។ សូមរក្សាទំនិញ និងកញ្ចប់រហូតដល់ពេលនោះ។",
		"rtInfo_refunded":  "បានសងប្រាក់ {a} នៅ {d}។", "rtInfo_cancelled": "អ្នកបានបោះបង់សំណើនេះ។",
		"rtSellerSaid": "អ្នកលក់", "rtPhsarSaid": "Phsar", "rtSent": "បានផ្ញើសំណើប្ដូរទំនិញ", "rtTooMany": "រូបថតច្រើនបំផុត 6", "rtUploading": "កំពុងផ្ទុករូបថត…",
		"ev_return_requested": "បានស្នើប្ដូរ", "ev_return_approved": "ការប្ដូរបានយល់ព្រម", "ev_return_disputed": "បញ្ជូនទៅ Phsar ពិនិត្យ",
		"ev_return_rejected": "ការប្ដូរមិនយល់ព្រម", "ev_return_refunded": "បានសងប្រាក់សម្រាប់ការប្ដូរ", "ev_return_cancelled": "បានបោះបង់ការប្ដូរ", "ev_completed": "បានបញ្ចប់",
	},
	"zh": {
		"back": "返回", "rtTitle": "退货或退款", "rtAsk": "申请退货", "rtUntil": "订单有问题？您可以在 {d} 前申请。",
		"rtWhich": "哪些商品？", "rtReason": "发生了什么？", "rtDesc": "描述问题", "rtDescPh": "例如：屏幕左侧有裂痕",
		"rtPhotos": "照片（商品损坏或发错时强烈建议上传）", "rtAddPhoto": "添加照片", "rtSend": "提交申请", "rtCancelReq": "取消申请",
		"rtConfirmCancel": "确定取消此退货申请吗？", "rtPolicy": "查看退货政策", "rtNo": "退货 {no}", "rtAmount": "申请 {a}",
		"rtReason_changed_mind": "不想要了", "rtReason_wrong_item": "发错商品", "rtReason_damaged": "损坏或破裂",
		"rtReason_not_as_described": "与视频不符", "rtReason_missing_parts": "缺少配件", "rtReason_other": "其他问题",
		"rtNoChangeMind": "这些商品已不能无理由退货——只能因商品问题申请。",
		"rtSt_requested": "等待卖家处理", "rtSt_disputed": "Phsar 审核中", "rtSt_approved": "已同意", "rtSt_rejected": "未同意",
		"rtSt_refunded": "已退款", "rtSt_cancelled": "已取消",
		"rtInfo_requested": "卖家需在 {d} 前回复，否则由 Phsar 决定。",
		"rtInfo_disputed":  "Phsar 正在审核您的申请，通常在 2 个工作日内决定。",
		"rtInfo_approved":  "我们将把 {a} 退回到付款账户，通常在 7 个工作日内。",
		"rtPickup":         "Phsar 快递员会联系您免费上门取件，在此之前请保留商品和包装。",
		"rtInfo_refunded":  "已于 {d} 退款 {a}。", "rtInfo_cancelled": "您已取消此申请。",
		"rtSellerSaid": "卖家", "rtPhsarSaid": "Phsar", "rtSent": "退货申请已提交", "rtTooMany": "最多 6 张照片", "rtUploading": "正在上传照片…",
		"ev_return_requested": "已申请退货", "ev_return_approved": "退货已同意", "ev_return_disputed": "已转交 Phsar 审核",
		"ev_return_rejected": "退货未同意", "ev_return_refunded": "退货已退款", "ev_return_cancelled": "退货已取消", "ev_completed": "已完成",
	},
}

func init() {
	for l, m := range returnStrs {
		for k, v := range m {
			shopStrs[l][k] = v
		}
	}
	for l := range shopStrs {
		shopStrs[l]["err_return_state"] = map[string]string{"en": "This return was already handled. Refresh the page.", "km": "ការប្ដូរនេះត្រូវបានដោះស្រាយរួចហើយ។ សូមផ្ទុកទំព័រឡើងវិញ។", "zh": "此退货已处理，请刷新页面。"}[l]
		shopStrs[l]["err_return_not_allowed"] = map[string]string{"en": "This order cannot be returned.", "km": "ការបញ្ជាទិញនេះមិនអាចប្ដូរបានទេ។", "zh": "此订单无法退货。"}[l]
	}
}
