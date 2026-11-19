package storefront

import (
	"encoding/json"
	"html/template"
)

// Teks untuk halaman keranjang / checkout / pesanan (dirender oleh shop.js di browser).
var shopStrs = map[string]map[string]string{
	"en": {
		"cartTitle": "Your cart", "cartEmpty": "Your cart is empty.", "continueShopping": "Continue shopping", "remove": "Remove",
		"subtotal": "Subtotal", "shipping": "Delivery", "free": "Free", "total": "Total", "checkout": "Checkout",
		"added": "Added to cart", "viewCart": "View cart", "chooseFirst": "Please choose {opt} first",
		"problem_unavailable": "No longer available", "problem_out_of_stock": "Out of stock", "problem_qty_reduced": "Only {n} left — quantity updated",
		"cartChanged": "Some items changed. Please check your cart before paying.", "removeProblems": "Remove unavailable items",
		"loginTitle": "Log in to place your order", "loginSub": "We confirm your phone number so the courier can reach you.",
		"tg": "Continue with Telegram", "or": "or", "phone": "Phone number", "sendCode": "Send code by SMS", "code": "6-digit code",
		"verify": "Log in", "codeSent": "Code sent.", "waitingTg": "Open Telegram and press Start in the bot…",
		"tgExpired":       "The Telegram link expired. Please try again.",
		"err_otp_invalid": "Wrong or expired code.", "err_otp_too_soon": "Please wait a minute before asking for a new code.",
		"err_phone_invalid": "This phone number is not valid.", "err_otp_rate_limited": "Too many codes requested. Please try later.",
		"err_telegram_disabled": "Telegram login is not set up yet. Use your phone number.",
		"deliverTo":             "Delivery address", "recipient": "Recipient name", "recipientPhone": "Phone for the courier",
		"searchPlace": "Search a place, street or building", "address": "Address", "addressPh": "House no., street, sangkat",
		"note": "Note for the courier", "notePh": "Gate colour, floor, landmark…", "pinHint": "Tap the map or drag the pin to the exact delivery spot.",
		"pinMissing": "Choose your delivery spot on the map.", "yourOrder": "Your order", "fromStore": "From {store}",
		"payKHQR": "Pay ${amount} with KHQR", "placing": "Creating your payment…", "khqrNote": "Pay with ABA Mobile or any Cambodian bank app (KHQR).", "fb": "Continue with Facebook", "guest": "Continue as guest", "guestNote": "No account needed. Your orders are saved on this phone/browser — log in next time to see them anywhere.", "guestUser": "Checking out as guest", "guestPaid": "You ordered as a guest. Keep your order number — log in next time to see your orders on any device.", "err_facebook_cancelled": "Facebook login was cancelled.", "err_facebook_failed": "Facebook login failed. Please try again.", "err_facebook_expired": "The Facebook login took too long. Please try again.", "err_facebook_disabled": "Facebook login is not available yet.", "err_sms_failed": "The SMS could not be sent right now. Try again or use Telegram.", "err_guest_rate_limited": "Too many guest checkouts from this network. Please log in or try again later.", "payMethod": "Payment method", "khqrSub": "Scan to pay with any banking app", "scanToPay": "Scan to Pay", "orWord": "or", "downloadQR": "Download QR", "uploadHint": "and upload to a mobile banking app that supports KHQR", "payABA": "Pay with ABA Mobile", "closedTitle": "Shopping opens soon", "closedText": "Your cart is saved on this phone. You can pay once Phsar opens.", "err_checkout_closed": "Phsar is opening soon. Shopping is not open yet.", "err_sms_disabled": "SMS login is not available yet. Please continue with Telegram.",
		"payTitle": "Scan to pay", "paySub": "Open ABA Mobile or any bank app with KHQR and scan this code.", "openABA": "Open ABA Mobile",
		"expiresIn": "Expires in {t}", "waitingPay": "Waiting for your payment…", "cancelPay": "Cancel order",
		"simulate": "Simulate payment (test mode)", "testMode": "Test mode — no real money is charged.",
		"paidTitle": "Payment received 🎉", "paidSub": "Your order was sent to the seller. A Phsar courier will deliver it to you for free.",
		"expiredTitle": "Payment time ran out", "expiredSub": "Your order was cancelled and nothing was charged. You can try again.",
		"cancelledTitle": "Order cancelled", "cancelledSub": "Nothing was charged.", "backToCart": "Back to cart", "seeOrders": "See my orders",
		"refundPending": "This order was cancelled. We are refunding ${amount} to your bank account — usually within 1–3 days.", "refundDone": "Refunded ${amount}.", "timeline": "Order history",
		"ev_paid": "Paid", "ev_accepted": "Seller accepted", "ev_ready_for_pickup": "Packed, waiting for the courier", "ev_picked_up": "Picked up by the courier", "ev_delivered": "Delivered",
		"ev_cancelled": "Cancelled", "ev_refund_pending": "Refund started", "ev_refunded": "Refunded",
		"myOrders": "My orders", "noOrders": "You have no orders yet.", "orderNo": "Order {no}", "items": "{n} items", "item1": "1 item", "placedOn": "Ordered {d}",
		"paidWith": "Paid with {m}", "deliverAddress": "Deliver to", "logout": "Log out", "loggedInAs": "Logged in as {p}",
		"status_pending_payment": "Waiting for payment", "status_paid": "Paid — waiting for the seller", "status_accepted": "Seller is preparing",
		"status_ready_for_pickup": "Ready for the courier", "status_picked_up": "On the way", "status_delivered": "Delivered",
		"status_completed": "Completed", "status_cancelled": "Cancelled", "status_refunded": "Refunded",
		"error": "Something went wrong. Please try again.", "network": "No connection. Check your internet and try again.",
		"required": "Please fill in this field.", "confirmCancel": "Cancel this order?",
		"useLocation": "📍 Use my current location", "locating": "Finding your location…", "locationDenied": "Couldn't get your location. Tap the map instead.", "pinSet": "Delivery spot set ✓",
	},
	"km": {
		"cartTitle": "កន្ត្រករបស់អ្នក", "cartEmpty": "កន្ត្រករបស់អ្នកទទេ។", "continueShopping": "បន្តទិញទំនិញ", "remove": "លុប",
		"subtotal": "សរុបរង", "shipping": "ការដឹកជញ្ជូន", "free": "ឥតគិតថ្លៃ", "total": "សរុប", "checkout": "បញ្ជាទិញ",
		"added": "បានដាក់ក្នុងកន្ត្រក", "viewCart": "មើលកន្ត្រក", "chooseFirst": "សូមជ្រើស{opt}ជាមុនសិន",
		"problem_unavailable": "លែងមានលក់", "problem_out_of_stock": "អស់ស្តុក", "problem_qty_reduced": "នៅសល់តែ {n} — បានកែចំនួន",
		"cartChanged": "ទំនិញខ្លះបានផ្លាស់ប្តូរ។ សូមពិនិត្យកន្ត្រកមុនបង់ប្រាក់។", "removeProblems": "លុបទំនិញដែលលែងមាន",
		"loginTitle": "ចូលគណនីដើម្បីបញ្ជាទិញ", "loginSub": "យើងផ្ទៀងផ្ទាត់លេខទូរស័ព្ទ ដើម្បីឱ្យអ្នកដឹកអាចទាក់ទងអ្នកបាន។",
		"tg": "បន្តជាមួយ Telegram", "or": "ឬ", "phone": "លេខទូរស័ព្ទ", "sendCode": "ផ្ញើលេខកូដតាម SMS", "code": "លេខកូដ 6 ខ្ទង់",
		"verify": "ចូល", "codeSent": "បានផ្ញើលេខកូដ។", "waitingTg": "បើក Telegram ហើយចុច Start ក្នុងបូត…",
		"tgExpired":       "តំណ Telegram ផុតកំណត់។ សូមព្យាយាមម្តងទៀត។",
		"err_otp_invalid": "លេខកូដខុស ឬផុតកំណត់។", "err_otp_too_soon": "សូមរង់ចាំមួយនាទីមុនស្នើលេខកូដថ្មី។",
		"err_phone_invalid": "លេខទូរស័ព្ទមិនត្រឹមត្រូវ។", "err_otp_rate_limited": "ស្នើលេខកូដច្រើនពេក។ សូមព្យាយាមពេលក្រោយ។",
		"err_telegram_disabled": "មិនទាន់កំណត់ការចូលតាម Telegram។ សូមប្រើលេខទូរស័ព្ទ។",
		"deliverTo":             "អាសយដ្ឋានដឹកជញ្ជូន", "recipient": "ឈ្មោះអ្នកទទួល", "recipientPhone": "លេខទូរស័ព្ទសម្រាប់អ្នកដឹក",
		"searchPlace": "ស្វែងរកទីកន្លែង ផ្លូវ ឬអគារ", "address": "អាសយដ្ឋាន", "addressPh": "លេខផ្ទះ ផ្លូវ សង្កាត់",
		"note": "កំណត់ចំណាំសម្រាប់អ្នកដឹក", "notePh": "ពណ៌ទ្វារ ជាន់ ចំណុចសម្គាល់…", "pinHint": "ចុចលើផែនទី ឬអូសម្ជុលទៅកន្លែងដឹកពិតប្រាកដ។",
		"pinMissing": "សូមជ្រើសកន្លែងដឹកលើផែនទី។", "yourOrder": "ការបញ្ជាទិញរបស់អ្នក", "fromStore": "ពីហាង {store}",
		"payKHQR": "បង់ ${amount} តាម KHQR", "placing": "កំពុងបង្កើតការបង់ប្រាក់…", "khqrNote": "បង់តាម ABA Mobile ឬកម្មវិធីធនាគារណាមួយក្នុងកម្ពុជា (KHQR)។", "fb": "បន្តជាមួយ Facebook", "guest": "បន្តជាភ្ញៀវ", "guestNote": "មិនចាំបាច់មានគណនី។ ការបញ្ជាទិញរបស់អ្នករក្សាទុកក្នុងទូរស័ព្ទ/កម្មវិធីរុករកនេះ — ចូលគណនីលើកក្រោយដើម្បីមើលគ្រប់ទីកន្លែង។", "guestUser": "កំពុងទិញជាភ្ញៀវ", "guestPaid": "អ្នកបានបញ្ជាទិញជាភ្ញៀវ។ សូមរក្សាលេខការបញ្ជាទិញ — ចូលគណនីលើកក្រោយដើម្បីមើលការបញ្ជាទិញលើឧបករណ៍ណាក៏បាន។", "err_facebook_cancelled": "បានបោះបង់ការចូលតាម Facebook។", "err_facebook_failed": "ការចូលតាម Facebook បរាជ័យ។ សូមព្យាយាមម្តងទៀត។", "err_facebook_expired": "ការចូលតាម Facebook យូរពេក។ សូមព្យាយាមម្តងទៀត។", "err_facebook_disabled": "មិនទាន់អាចចូលតាម Facebook បានទេ។", "err_sms_failed": "មិនអាចផ្ញើ SMS ឥឡូវនេះទេ។ សូមព្យាយាមម្តងទៀត ឬប្រើ Telegram។", "err_guest_rate_limited": "ការទិញជាភ្ញៀវច្រើនពេកពីបណ្តាញនេះ។ សូមចូលគណនី ឬព្យាយាមពេលក្រោយ។", "payMethod": "វិធីបង់ប្រាក់", "khqrSub": "ស្កេនដើម្បីបង់ជាមួយកម្មវិធីធនាគារណាមួយ", "scanToPay": "ស្កេនដើម្បីបង់", "orWord": "ឬ", "downloadQR": "ទាញយក QR", "uploadHint": "ហើយបង្ហោះទៅកម្មវិធីធនាគារដែលគាំទ្រ KHQR", "payABA": "បង់ជាមួយ ABA Mobile", "closedTitle": "ការទិញនឹងបើកឆាប់ៗ", "closedText": "កន្ត្រករបស់អ្នកត្រូវបានរក្សាទុកក្នុងទូរស័ព្ទនេះ។ អ្នកអាចបង់ប្រាក់នៅពេលផ្សារបើក។", "err_checkout_closed": "ផ្សារនឹងបើកឆាប់ៗនេះ។ មិនទាន់អាចទិញបានទេ។", "err_sms_disabled": "មិនទាន់អាចចូលដោយ SMS បានទេ។ សូមបន្តជាមួយ Telegram។",
		"payTitle": "ស្កេនដើម្បីបង់ប្រាក់", "paySub": "បើក ABA Mobile ឬកម្មវិធីធនាគារដែលគាំទ្រ KHQR ហើយស្កេនកូដនេះ។", "openABA": "បើក ABA Mobile",
		"expiresIn": "ផុតកំណត់ក្នុង {t}", "waitingPay": "កំពុងរង់ចាំការបង់ប្រាក់…", "cancelPay": "បោះបង់ការបញ្ជាទិញ",
		"simulate": "សាកល្បងបង់ប្រាក់ (របៀបតេស្ត)", "testMode": "របៀបតេស្ត — មិនកាត់លុយពិតទេ។",
		"paidTitle": "បានទទួលការបង់ប្រាក់ 🎉", "paidSub": "ការបញ្ជាទិញបានផ្ញើទៅអ្នកលក់ហើយ។ អ្នកដឹករបស់ Phsar នឹងដឹកជូនអ្នកដោយឥតគិតថ្លៃ។",
		"expiredTitle": "ផុតពេលបង់ប្រាក់", "expiredSub": "ការបញ្ជាទិញត្រូវបានបោះបង់ ហើយមិនបានកាត់លុយទេ។ អ្នកអាចព្យាយាមម្តងទៀត។",
		"cancelledTitle": "បានបោះបង់ការបញ្ជាទិញ", "cancelledSub": "មិនបានកាត់លុយទេ។", "backToCart": "ត្រឡប់ទៅកន្ត្រក", "seeOrders": "មើលការបញ្ជាទិញរបស់ខ្ញុំ",
		"refundPending": "ការបញ្ជាទិញនេះត្រូវបានបោះបង់។ យើងកំពុងសងប្រាក់ ${amount} ទៅគណនីធនាគាររបស់អ្នក — ជាធម្មតាក្នុង 1–3 ថ្ងៃ។", "refundDone": "បានសងប្រាក់ ${amount}។", "timeline": "ប្រវត្តិការបញ្ជាទិញ",
		"ev_paid": "បានបង់ប្រាក់", "ev_accepted": "អ្នកលក់បានទទួល", "ev_ready_for_pickup": "វេចខ្ចប់រួច រង់ចាំអ្នកដឹក", "ev_picked_up": "អ្នកដឹកបានយក", "ev_delivered": "បានដឹកដល់",
		"ev_cancelled": "បានបោះបង់", "ev_refund_pending": "ចាប់ផ្តើមសងប្រាក់", "ev_refunded": "បានសងប្រាក់",
		"myOrders": "ការបញ្ជាទិញរបស់ខ្ញុំ", "noOrders": "អ្នកមិនទាន់មានការបញ្ជាទិញទេ។", "orderNo": "ការបញ្ជាទិញ {no}", "items": "{n} មុខ", "item1": "1 មុខ", "placedOn": "បញ្ជាទិញ {d}",
		"paidWith": "បង់តាម {m}", "deliverAddress": "ដឹកទៅ", "logout": "ចាកចេញ", "loggedInAs": "បានចូលជា {p}",
		"status_pending_payment": "កំពុងរង់ចាំការបង់ប្រាក់", "status_paid": "បានបង់ — រង់ចាំអ្នកលក់", "status_accepted": "អ្នកលក់កំពុងរៀបចំ",
		"status_ready_for_pickup": "រួចរាល់សម្រាប់អ្នកដឹក", "status_picked_up": "កំពុងដឹក", "status_delivered": "បានដឹកដល់",
		"status_completed": "បានបញ្ចប់", "status_cancelled": "បានបោះបង់", "status_refunded": "បានសងប្រាក់វិញ",
		"error": "មានបញ្ហា។ សូមព្យាយាមម្តងទៀត។", "network": "គ្មានអ៊ីនធឺណិត។ សូមពិនិត្យហើយព្យាយាមម្តងទៀត។",
		"required": "សូមបំពេញចន្លោះនេះ។", "confirmCancel": "បោះបង់ការបញ្ជាទិញនេះ?",
		"useLocation": "📍 ប្រើទីតាំងបច្ចុប្បន្នរបស់ខ្ញុំ", "locating": "កំពុងស្វែងរកទីតាំង…", "locationDenied": "រកទីតាំងមិនបាន។ សូមចុចលើផែនទីជំនួស។", "pinSet": "បានកំណត់កន្លែងដឹក ✓",
	},
	"zh": {
		"cartTitle": "购物车", "cartEmpty": "购物车是空的。", "continueShopping": "继续购物", "remove": "删除",
		"subtotal": "小计", "shipping": "配送", "free": "免费", "total": "合计", "checkout": "去结算",
		"added": "已加入购物车", "viewCart": "查看购物车", "chooseFirst": "请先选择{opt}",
		"problem_unavailable": "已下架", "problem_out_of_stock": "已售罄", "problem_qty_reduced": "仅剩 {n} 件，数量已调整",
		"cartChanged": "部分商品有变动，请在付款前检查购物车。", "removeProblems": "删除无法购买的商品",
		"loginTitle": "登录后下单", "loginSub": "验证手机号，方便快递员联系您。",
		"tg": "使用 Telegram 继续", "or": "或", "phone": "手机号", "sendCode": "发送短信验证码", "code": "6 位验证码",
		"verify": "登录", "codeSent": "验证码已发送。", "waitingTg": "请打开 Telegram 并在机器人中点击 Start…",
		"tgExpired":       "Telegram 链接已过期，请重试。",
		"err_otp_invalid": "验证码错误或已过期。", "err_otp_too_soon": "请稍等一分钟再获取验证码。",
		"err_phone_invalid": "手机号无效。", "err_otp_rate_limited": "验证码请求过多，请稍后再试。",
		"err_telegram_disabled": "暂未开通 Telegram 登录，请使用手机号。",
		"deliverTo":             "收货地址", "recipient": "收件人", "recipientPhone": "快递员联系电话",
		"searchPlace": "搜索地点、街道或楼名", "address": "地址", "addressPh": "门牌号、街道、分区",
		"note": "给快递员的备注", "notePh": "大门颜色、楼层、地标…", "pinHint": "点击地图或拖动图钉到准确的送货位置。",
		"pinMissing": "请在地图上选择送货位置。", "yourOrder": "您的订单", "fromStore": "来自 {store}",
		"payKHQR": "用 KHQR 支付 ${amount}", "placing": "正在创建付款…", "khqrNote": "可用 ABA Mobile 或任何柬埔寨银行 App（KHQR）付款。", "fb": "使用 Facebook 继续", "guest": "以游客身份继续", "guestNote": "无需注册。订单保存在这台手机/浏览器上——下次登录即可在任何设备查看。", "guestUser": "正在以游客身份结账", "guestPaid": "您以游客身份下单。请保存订单号——下次登录即可在任何设备查看订单。", "err_facebook_cancelled": "已取消 Facebook 登录。", "err_facebook_failed": "Facebook 登录失败，请重试。", "err_facebook_expired": "Facebook 登录超时，请重试。", "err_facebook_disabled": "暂不支持 Facebook 登录。", "err_sms_failed": "短信暂时无法发送，请重试或使用 Telegram。", "err_guest_rate_limited": "此网络的游客下单次数过多，请登录或稍后再试。", "payMethod": "付款方式", "khqrSub": "使用任何银行 App 扫码付款", "scanToPay": "扫码付款", "orWord": "或", "downloadQR": "下载二维码", "uploadHint": "然后在支持 KHQR 的手机银行 App 中上传", "payABA": "使用 ABA Mobile 付款", "closedTitle": "购物即将开放", "closedText": "购物车已保存在这台手机上，Phsar 开业后即可付款。", "err_checkout_closed": "Phsar 即将开业，暂时还不能下单。", "err_sms_disabled": "暂不支持短信登录，请使用 Telegram 继续。",
		"payTitle": "扫码付款", "paySub": "打开 ABA Mobile 或支持 KHQR 的银行 App 扫描此码。", "openABA": "打开 ABA Mobile",
		"expiresIn": "{t} 后过期", "waitingPay": "等待付款中…", "cancelPay": "取消订单",
		"simulate": "模拟付款（测试模式）", "testMode": "测试模式——不会真实扣款。",
		"paidTitle": "付款成功 🎉", "paidSub": "订单已发送给卖家，Phsar 快递员将免费为您送达。",
		"expiredTitle": "付款已超时", "expiredSub": "订单已取消，未扣款。您可以重新下单。",
		"cancelledTitle": "订单已取消", "cancelledSub": "未扣款。", "backToCart": "返回购物车", "seeOrders": "查看我的订单",
		"refundPending": "该订单已取消，我们正在将 ${amount} 退回您的银行账户，通常 1–3 天内到账。", "refundDone": "已退款 ${amount}。", "timeline": "订单记录",
		"ev_paid": "已付款", "ev_accepted": "卖家已接单", "ev_ready_for_pickup": "已打包，等待快递员", "ev_picked_up": "快递员已取件", "ev_delivered": "已送达",
		"ev_cancelled": "已取消", "ev_refund_pending": "退款处理中", "ev_refunded": "已退款",
		"myOrders": "我的订单", "noOrders": "您还没有订单。", "orderNo": "订单 {no}", "items": "{n} 件商品", "item1": "1 件商品", "placedOn": "下单时间 {d}",
		"paidWith": "支付方式 {m}", "deliverAddress": "送至", "logout": "退出登录", "loggedInAs": "已登录 {p}",
		"status_pending_payment": "等待付款", "status_paid": "已付款，等待卖家处理", "status_accepted": "卖家备货中",
		"status_ready_for_pickup": "等待快递员取件", "status_picked_up": "配送中", "status_delivered": "已送达",
		"status_completed": "已完成", "status_cancelled": "已取消", "status_refunded": "已退款",
		"error": "出错了，请重试。", "network": "网络连接失败，请检查网络后重试。",
		"required": "请填写此项。", "confirmCancel": "确定取消此订单吗？",
		"useLocation": "📍 使用我的当前位置", "locating": "正在获取位置…", "locationDenied": "无法获取位置，请直接点击地图。", "pinSet": "已设置送货位置 ✓",
	},
}

// ShopT: teks belanja untuk template (bahasa halaman → Inggris).
func (p *Page) ShopT() map[string]string {
	m := map[string]string{}
	for k, v := range shopStrs["en"] {
		m[k] = v
	}
	for k, v := range shopStrs[p.Lang] {
		m[k] = v
	}
	return m
}

// ShopJS: teks untuk shop.js (bahasa halaman, kosong → Inggris).
func (p *Page) ShopJS() template.JS {
	m := map[string]string{}
	for k, v := range shopStrs["en"] {
		m[k] = v
	}
	for k, v := range shopStrs[p.Lang] {
		m[k] = v
	}
	b, _ := json.Marshal(m)
	return template.JS(b)
}

func init() {
	// Judul halaman (dipakai Page.T di template).
	for l, m := range shopStrs {
		for _, k := range []string{"cartTitle", "myOrders", "checkout"} {
			strs[l][k] = m[k]
		}
	}
}
