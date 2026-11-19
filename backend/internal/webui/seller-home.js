'use strict';
// Seller Center v2: navigasi 5 bagian (Hari ini · Pesanan · Produk · Uang · Lainnya),
// layar "Hari ini" (daftar tugas), "Mulai jualan" (setup toko langkah demi langkah) dan "Lainnya" (menu).
addStrings({
  en: {durH: '{n} h', durM: '{n} min', durHM: '{h} h {m} min', openHrs: 'open hours', tk_packDue: 'pack by',
    closedNow: 'Your shop is closed now — order time limits are paused until {d}.', closedNowNoTime: 'Your shop is closed now — order time limits are paused.',
    mi_hours: 'Opening hours', mi_hoursD: 'Order time limits only run while you are open',
    hrTitle: 'Opening hours', hrWhy: 'Phsar only counts your order time limits while your shop is open. Example: an order paid at night waits until you open in the morning.',
    hrWhyOff: 'Order time limits currently run all day (Phsar setting). Your hours are still shown to Phsar support.', hrRules: 'Accept new orders within {a}. Mark accepted orders ready for pickup within {p}.', hrRulesNoPack: 'Accept new orders within {a}.',
    hrUsingDefault: 'You are using Phsar’s standard hours: {s}', hrUseMine: 'Set my own hours', hrUseDefault: 'Use Phsar’s standard hours', hrSave: 'Save opening hours', hrSaved: 'Opening hours saved', hrOpenNow: 'Open now', hrClosedNow: 'Closed now',
    su_hours: 'Opening hours', su_hoursN: 'Order time limits pause while you are closed', v2nav: 'Seller sections', nav_home: 'Today', nav_orders: 'Orders', nav_products: 'Products', nav_payouts: 'Money', nav_more: 'More',
    h_orders: 'Orders', h_products: 'Products', h_payouts: 'My money', h_more: 'More', backMore: '← More',
    hello: 'Hello, {n}', todayTitle: 'Today', doNow: 'Do these now', allDoneNow: 'All done for now 🎉 New tasks will show up here.',
    tk_new: 'Accept {n} new order(s)', tk_newW: 'Buyers are waiting. Orders not accepted within {a} cancel automatically.', tk_newB: 'Accept',
    tk_pack: 'Pack {n} order(s)', tk_packW: 'Pack the items, stick the Phsar label on, then mark them ready for pickup. Late orders cancel automatically.', tk_packB: 'Pack',
    tk_ret: 'Answer {n} return request(s)', tk_retW: 'Accept or explain. If you don’t answer in time, Phsar decides.', tk_retB: 'Answer',
    tk_fix: 'Fix {n} product(s)', tk_fixW: 'Phsar asked for changes before these products can be shown.', tk_fixB: 'Fix',
    tk_flash: 'Flash sale invitation: {n} product(s)', tk_flashW: 'Phsar invited your products to a flash sale. Accept or decline.', tk_flashB: 'View',
    tk_tick: 'Phsar replied to {n} help request(s)', tk_tickW: 'Read the reply and answer if needed.', tk_tickB: 'Read',
    tk_bank: 'Add your bank account', tk_bankW: 'Without it Phsar can’t send you money.', tk_bankB: 'Add',
    tk_ready: '{n} order(s) waiting for the courier', tk_readyW: 'Nothing to do — the Phsar courier will pick them up.', tk_readyB: 'View',
    leftT: '{t} left', lateT: 'Time is up', startsT: 'starts {d}',
    salesToday: 'Sales today', ordersN: '{n} orders', readyMoney: 'Ready to be paid', seeMoney: 'My money →',
    tipTitle: 'Tip of the day',
    tips: ['Reply to buyers fast — shops that answer quickly are shown higher in search.', 'Products with a clear video sell much more often. Show the item in your hand.', 'Accept orders quickly. Late or rejected orders lower your shop score.'],
    staffHome: 'You are helping this shop as staff. You can handle orders, returns and products.',
    su_title: 'Open your shop in 5 steps', su_sub: 'You can stop and continue any time.', nextStep: 'Next step', allSteps: 'All steps',
    su_profile: 'Shop name, logo & description', su_profileN: 'So buyers recognise you', su_profileW: 'Add a clear logo and a short description of what you sell. You can change it later.', su_profileC: 'Add logo & description',
    su_kyc: 'Identity check', su_kycN: 'ID card or passport photo', su_kycW: 'Only the Phsar team sees your document. It keeps buyers and sellers safe.', su_kycC: 'Upload my ID',
    su_pickup: 'Pickup point', su_pickupN: 'Where the courier collects your packages', su_pickupW: 'Put the pin exactly at your door so the Phsar courier can find you.', su_pickupC: 'Set my pickup point',
    su_submit: 'Send for review', su_submitN: 'Phsar checks everything in 1–2 working days', su_submitW: 'Everything is filled in. Send your shop to Phsar for checking.', su_submitC: 'Send for review',
    su_review: 'Approved by Phsar', su_reviewN: 'Then you can start selling',
    su_waitT: 'Phsar is checking your shop', su_waitW: 'This usually takes 1–2 working days. We will send you a notification. Meanwhile you can prepare products.', su_waitC: 'Prepare a product',
    su_min: '~{n} min', su_done: 'Done', su_todo: 'To do', su_wait: 'Waiting',
    su_extra: 'You can also prepare now', su_bank: 'Bank account', su_bankN: 'So Phsar can pay you', su_prod: 'First product (draft)', su_prodN: 'It goes live after your shop is approved',
    su_help: 'Stuck? Ask Phsar support', su_suspended: 'This shop is suspended. Contact Phsar support for help.', su_closed: 'This shop is closed.',
    staffWait: 'This shop is not open yet. The owner is still setting it up.',
    mg_sell: 'Sell more', mg_buyers: 'Buyers', mg_shop: 'My shop', mg_help: 'Help & account',
    mi_flash: 'Flash sale', mi_flashD: 'Invitations and your flash sale products', mi_ads: 'Ads', mi_adsD: 'Show your products at the top',
    mi_returns: 'Returns', mi_returnsD: 'Accept or explain return requests',
    mi_profile: 'Shop profile', mi_profileD: 'Name, logo, banner, description', mi_kyc: 'Identity (KYC)', mi_kycD: 'Your ID check status',
    mi_pickup: 'Pickup points', mi_pickupD: 'Where the courier collects', mi_team: 'Team', mi_teamD: 'Let staff help with orders',
    mi_help: 'Help', mi_helpD: 'Ask Phsar support, see your requests', mi_stores: 'My shops', mi_storesD: 'Your shop and shops you help as staff',
    mi_logout: 'Sign out', mi_logoutD: 'Sign out of Seller Center', isNew: 'New', langTitle: 'Language'},
  km: {durH: '{n} ម៉ោង', durM: '{n} នាទី', durHM: '{h} ម៉ោង {m} នាទី', openHrs: 'ម៉ោងបើក', tk_packDue: 'វេចខ្ចប់មុន',
    closedNow: 'ហាងរបស់អ្នកកំពុងបិទ — ពេលកំណត់ការបញ្ជាទិញត្រូវបានផ្អាករហូតដល់ {d}។', closedNowNoTime: 'ហាងរបស់អ្នកកំពុងបិទ — ពេលកំណត់ការបញ្ជាទិញត្រូវបានផ្អាក។',
    mi_hours: 'ម៉ោងបើក', mi_hoursD: 'ពេលកំណត់ការបញ្ជាទិញរាប់តែពេលហាងបើក',
    hrTitle: 'ម៉ោងបើក', hrWhy: 'Phsar រាប់ពេលកំណត់ការបញ្ជាទិញតែពេលហាងបើកប៉ុណ្ណោះ។ ឧទាហរណ៍៖ ការបញ្ជាទិញដែលបង់ពេលយប់ រង់ចាំរហូតដល់អ្នកបើកពេលព្រឹក។',
    hrWhyOff: 'បច្ចុប្បន្នពេលកំណត់រាប់ពេញមួយថ្ងៃ (ការកំណត់ Phsar)។ ម៉ោងរបស់អ្នកនៅតែបង្ហាញដល់ផ្នែកជំនួយ Phsar។', hrRules: 'ទទួលការបញ្ជាទិញថ្មីក្នុង {a}។ សម្គាល់ការបញ្ជាទិញដែលបានទទួលថារួចរាល់ឱ្យមកយកក្នុង {p}។', hrRulesNoPack: 'ទទួលការបញ្ជាទិញថ្មីក្នុង {a}។',
    hrUsingDefault: 'អ្នកកំពុងប្រើម៉ោងស្តង់ដារ Phsar៖ {s}', hrUseMine: 'កំណត់ម៉ោងផ្ទាល់ខ្លួន', hrUseDefault: 'ប្រើម៉ោងស្តង់ដារ Phsar', hrSave: 'រក្សាទុកម៉ោងបើក', hrSaved: 'បានរក្សាទុកម៉ោងបើក', hrOpenNow: 'កំពុងបើក', hrClosedNow: 'កំពុងបិទ',
    su_hours: 'ម៉ោងបើក', su_hoursN: 'ពេលកំណត់ការបញ្ជាទិញផ្អាកពេលហាងបិទ', v2nav: 'ផ្នែកអ្នកលក់', nav_home: 'ថ្ងៃនេះ', nav_orders: 'ការបញ្ជាទិញ', nav_products: 'ផលិតផល', nav_payouts: 'លុយ', nav_more: 'ច្រើនទៀត',
    h_orders: 'ការបញ្ជាទិញ', h_products: 'ផលិតផល', h_payouts: 'លុយរបស់ខ្ញុំ', h_more: 'ច្រើនទៀត', backMore: '← ច្រើនទៀត',
    hello: 'សួស្តី {n}', todayTitle: 'ថ្ងៃនេះ', doNow: 'ធ្វើទាំងនេះឥឡូវ', allDoneNow: 'រួចរាល់ទាំងអស់ 🎉 កិច្ចការថ្មីនឹងបង្ហាញនៅទីនេះ។',
    tk_new: 'ទទួលការបញ្ជាទិញថ្មី {n}', tk_newW: 'អ្នកទិញកំពុងរង់ចាំ។ ការបញ្ជាទិញដែលមិនទទួលក្នុង {a} នឹងបោះបង់ដោយស្វ័យប្រវត្តិ។', tk_newB: 'ទទួល',
    tk_pack: 'វេចខ្ចប់ការបញ្ជាទិញ {n}', tk_packW: 'វេចខ្ចប់ទំនិញ បិទស្លាក Phsar រួចសម្គាល់ថារួចរាល់ឱ្យមកយក។ ការបញ្ជាទិញយឺតនឹងបោះបង់ដោយស្វ័យប្រវត្តិ។', tk_packB: 'វេចខ្ចប់',
    tk_ret: 'ឆ្លើយសំណើប្តូរ/សង {n}', tk_retW: 'ទទួល ឬពន្យល់។ បើមិនឆ្លើយទាន់ពេល Phsar នឹងសម្រេច។', tk_retB: 'ឆ្លើយ',
    tk_fix: 'កែផលិតផល {n}', tk_fixW: 'Phsar ស្នើឱ្យកែ មុនពេលផលិតផលទាំងនេះអាចបង្ហាញ។', tk_fixB: 'កែ',
    tk_flash: 'ការអញ្ជើញ Flash sale៖ ផលិតផល {n}', tk_flashW: 'Phsar អញ្ជើញផលិតផលរបស់អ្នកចូល Flash sale។ ទទួល ឬបដិសេធ។', tk_flashB: 'មើល',
    tk_tick: 'Phsar បានឆ្លើយសំណើជំនួយ {n}', tk_tickW: 'អានចម្លើយ ហើយឆ្លើយតបបើចាំបាច់។', tk_tickB: 'អាន',
    tk_bank: 'បន្ថែមគណនីធនាគារ', tk_bankW: 'បើគ្មាន Phsar មិនអាចផ្ញើលុយជូនអ្នកបានទេ។', tk_bankB: 'បន្ថែម',
    tk_ready: 'ការបញ្ជាទិញ {n} កំពុងរង់ចាំអ្នកដឹក', tk_readyW: 'មិនចាំបាច់ធ្វើអ្វីទេ — អ្នកដឹក Phsar នឹងមកយក។', tk_readyB: 'មើល',
    leftT: 'នៅសល់ {t}', lateT: 'ផុតពេលហើយ', startsT: 'ចាប់ផ្តើម {d}',
    salesToday: 'ការលក់ថ្ងៃនេះ', ordersN: 'ការបញ្ជាទិញ {n}', readyMoney: 'រួចរាល់សម្រាប់ទូទាត់', seeMoney: 'លុយរបស់ខ្ញុំ →',
    tipTitle: 'គន្លឹះប្រចាំថ្ងៃ',
    tips: ['ឆ្លើយអ្នកទិញឱ្យលឿន — ហាងដែលឆ្លើយលឿនត្រូវបានបង្ហាញខ្ពស់ក្នុងការស្វែងរក។', 'ផលិតផលដែលមានវីដេអូច្បាស់លក់ដាច់ជាងច្រើន។ បង្ហាញទំនិញនៅក្នុងដៃ។', 'ទទួលការបញ្ជាទិញឱ្យលឿន។ ការយឺត ឬបដិសេធធ្វើឱ្យពិន្ទុហាងធ្លាក់។'],
    staffHome: 'អ្នកកំពុងជួយហាងនេះជាបុគ្គលិក។ អ្នកអាចគ្រប់គ្រងការបញ្ជាទិញ ការប្តូរ/សង និងផលិតផល។',
    su_title: 'បើកហាងក្នុង 5 ជំហាន', su_sub: 'អាចឈប់ ហើយបន្តពេលណាក៏បាន។', nextStep: 'ជំហានបន្ទាប់', allSteps: 'ជំហានទាំងអស់',
    su_profile: 'ឈ្មោះហាង ឡូហ្គោ និងការពិពណ៌នា', su_profileN: 'ដើម្បីឱ្យអ្នកទិញស្គាល់អ្នក', su_profileW: 'បន្ថែមឡូហ្គោច្បាស់ និងការពិពណ៌នាខ្លីអំពីអ្វីដែលអ្នកលក់។ អាចប្តូរពេលក្រោយ។', su_profileC: 'បន្ថែមឡូហ្គោ និងការពិពណ៌នា',
    su_kyc: 'ផ្ទៀងផ្ទាត់អត្តសញ្ញាណ', su_kycN: 'រូបថតអត្តសញ្ញាណប័ណ្ណ ឬលិខិតឆ្លងដែន', su_kycW: 'មានតែក្រុម Phsar ដែលមើលឃើញឯកសារ។ ដើម្បីសុវត្ថិភាពអ្នកទិញ និងអ្នកលក់។', su_kycC: 'ផ្ទុកអត្តសញ្ញាណប័ណ្ណ',
    su_pickup: 'ទីតាំងទទួលកញ្ចប់', su_pickupN: 'កន្លែងអ្នកដឹកមកយកកញ្ចប់', su_pickupW: 'ដាក់ម្ជុលត្រង់ទ្វារហាង ដើម្បីឱ្យអ្នកដឹក Phsar រកឃើញ។', su_pickupC: 'កំណត់ទីតាំងទទួល',
    su_submit: 'ផ្ញើឱ្យត្រួតពិនិត្យ', su_submitN: 'Phsar ពិនិត្យក្នុង 1–2 ថ្ងៃធ្វើការ', su_submitW: 'បានបំពេញគ្រប់យ៉ាង។ ផ្ញើហាងទៅ Phsar ដើម្បីពិនិត្យ។', su_submitC: 'ផ្ញើឱ្យត្រួតពិនិត្យ',
    su_review: 'Phsar បានអនុម័ត', su_reviewN: 'បន្ទាប់មកអ្នកអាចចាប់ផ្តើមលក់',
    su_waitT: 'Phsar កំពុងពិនិត្យហាងរបស់អ្នក', su_waitW: 'ជាធម្មតា 1–2 ថ្ងៃធ្វើការ។ យើងនឹងផ្ញើការជូនដំណឹង។ ចន្លោះពេលនេះ អ្នកអាចរៀបចំផលិតផល។', su_waitC: 'រៀបចំផលិតផល',
    su_min: '~{n} នាទី', su_done: 'រួចរាល់', su_todo: 'ត្រូវធ្វើ', su_wait: 'កំពុងរង់ចាំ',
    su_extra: 'អាចរៀបចំឥឡូវនេះផងដែរ', su_bank: 'គណនីធនាគារ', su_bankN: 'ដើម្បីឱ្យ Phsar ទូទាត់ជូនអ្នក', su_prod: 'ផលិតផលដំបូង (សេចក្តីព្រាង)', su_prodN: 'បង្ហាញបន្ទាប់ពីហាងត្រូវបានអនុម័ត',
    su_help: 'ជាប់គាំង? សួរផ្នែកជំនួយ Phsar', su_suspended: 'ហាងនេះត្រូវបានផ្អាក។ សូមទាក់ទងផ្នែកជំនួយ Phsar។', su_closed: 'ហាងនេះបានបិទ។',
    staffWait: 'ហាងនេះមិនទាន់បើកទេ។ ម្ចាស់ហាងកំពុងរៀបចំ។',
    mg_sell: 'លក់បានច្រើន', mg_buyers: 'អ្នកទិញ', mg_shop: 'ហាងខ្ញុំ', mg_help: 'ជំនួយ និងគណនី',
    mi_flash: 'Flash sale', mi_flashD: 'ការអញ្ជើញ និងផលិតផល Flash sale របស់អ្នក', mi_ads: 'ការផ្សាយពាណិជ្ជកម្ម', mi_adsD: 'បង្ហាញផលិតផលនៅខាងលើ',
    mi_returns: 'ការប្តូរ/សង', mi_returnsD: 'ទទួល ឬពន្យល់សំណើប្តូរ/សង',
    mi_profile: 'ព័ត៌មានហាង', mi_profileD: 'ឈ្មោះ ឡូហ្គោ បដា ការពិពណ៌នា', mi_kyc: 'អត្តសញ្ញាណ (KYC)', mi_kycD: 'ស្ថានភាពផ្ទៀងផ្ទាត់',
    mi_pickup: 'ទីតាំងទទួល', mi_pickupD: 'កន្លែងអ្នកដឹកមកយក', mi_team: 'ក្រុម', mi_teamD: 'ឱ្យបុគ្គលិកជួយការបញ្ជាទិញ',
    mi_help: 'ជំនួយ', mi_helpD: 'សួរផ្នែកជំនួយ Phsar មើលសំណើរបស់អ្នក', mi_stores: 'ហាងរបស់ខ្ញុំ', mi_storesD: 'ហាងរបស់អ្នក និងហាងដែលអ្នកជួយជាបុគ្គលិក',
    mi_logout: 'ចាកចេញ', mi_logoutD: 'ចាកចេញពីមជ្ឈមណ្ឌលអ្នកលក់', isNew: 'ថ្មី', langTitle: 'ភាសា'},
  zh: {durH: '{n} 小时', durM: '{n} 分钟', durHM: '{h} 小时 {m} 分钟', openHrs: '营业时间', tk_packDue: '打包截止',
    closedNow: '店铺目前休息中——订单时限暂停，{d} 恢复计时。', closedNowNoTime: '店铺目前休息中——订单时限已暂停。',
    mi_hours: '营业时间', mi_hoursD: '订单时限只在营业时间内计算',
    hrTitle: '营业时间', hrWhy: 'Phsar 只在店铺营业时计算订单时限。例如：夜间付款的订单会等到你早上开门后才开始计时。',
    hrWhyOff: '目前订单时限全天计算（Phsar 设置）。你的营业时间仍会显示给 Phsar 客服。', hrRules: '新订单请在 {a} 内接单。已接订单请在 {p} 内标记为待取件。', hrRulesNoPack: '新订单请在 {a} 内接单。',
    hrUsingDefault: '当前使用 Phsar 标准营业时间：{s}', hrUseMine: '设置我的营业时间', hrUseDefault: '使用 Phsar 标准时间', hrSave: '保存营业时间', hrSaved: '营业时间已保存', hrOpenNow: '营业中', hrClosedNow: '休息中',
    su_hours: '营业时间', su_hoursN: '休息时订单时限暂停', v2nav: '卖家栏目', nav_home: '今天', nav_orders: '订单', nav_products: '商品', nav_payouts: '收入', nav_more: '更多',
    h_orders: '订单', h_products: '商品', h_payouts: '我的收入', h_more: '更多', backMore: '← 更多',
    hello: '你好，{n}', todayTitle: '今天', doNow: '现在要做', allDoneNow: '全部完成 🎉 新任务会显示在这里。',
    tk_new: '接受 {n} 笔新订单', tk_newW: '买家在等待。{a} 内未接受的订单会自动取消。', tk_newB: '接单',
    tk_pack: '打包 {n} 笔订单', tk_packW: '打包商品，贴上 Phsar 面单，然后标记为待取件。超时订单会自动取消。', tk_packB: '打包',
    tk_ret: '处理 {n} 个退货申请', tk_retW: '同意或说明原因。逾期未处理将由 Phsar 决定。', tk_retB: '处理',
    tk_fix: '修改 {n} 个商品', tk_fixW: 'Phsar 要求修改后这些商品才能上架。', tk_fixB: '修改',
    tk_flash: '限时抢购邀请：{n} 个商品', tk_flashW: 'Phsar 邀请你的商品参加限时抢购。请接受或拒绝。', tk_flashB: '查看',
    tk_tick: 'Phsar 回复了 {n} 个求助', tk_tickW: '查看回复，如有需要请回复。', tk_tickB: '查看',
    tk_bank: '添加银行账户', tk_bankW: '没有银行账户，Phsar 无法给你打款。', tk_bankB: '添加',
    tk_ready: '{n} 笔订单等待快递员取件', tk_readyW: '无需操作，Phsar 快递员会上门取件。', tk_readyB: '查看',
    leftT: '剩 {t}', lateT: '已超时', startsT: '{d} 开始',
    salesToday: '今日销售', ordersN: '{n} 笔订单', readyMoney: '待结算', seeMoney: '我的收入 →',
    tipTitle: '今日小贴士',
    tips: ['尽快回复买家——回复快的店铺在搜索中排名更高。', '有清晰视频的商品成交率高得多。把商品拿在手里展示。', '尽快接单。超时或拒单会降低店铺评分。'],
    staffHome: '你是这家店的员工，可以处理订单、退货和商品。',
    su_title: '5 步开店', su_sub: '可随时暂停，之后继续。', nextStep: '下一步', allSteps: '全部步骤',
    su_profile: '店名、标志和介绍', su_profileN: '让买家认出你', su_profileW: '上传清晰的标志，并简单介绍你卖什么。之后可以修改。', su_profileC: '添加标志和介绍',
    su_kyc: '身份验证', su_kycN: '身份证或护照照片', su_kycW: '证件仅 Phsar 团队可见，保障买卖双方安全。', su_kycC: '上传证件',
    su_pickup: '取件点', su_pickupN: '快递员取件的地方', su_pickupW: '把定位放在店门口，方便 Phsar 快递员找到你。', su_pickupC: '设置取件点',
    su_submit: '提交审核', su_submitN: 'Phsar 会在 1–2 个工作日内审核', su_submitW: '资料已填写完整，提交给 Phsar 审核。', su_submitC: '提交审核',
    su_review: 'Phsar 审核通过', su_reviewN: '之后就可以开始销售',
    su_waitT: 'Phsar 正在审核你的店铺', su_waitW: '通常需要 1–2 个工作日，结果会通知你。现在可以先准备商品。', su_waitC: '准备商品',
    su_min: '约 {n} 分钟', su_done: '已完成', su_todo: '待完成', su_wait: '等待中',
    su_extra: '现在也可以准备', su_bank: '银行账户', su_bankN: '用于 Phsar 给你打款', su_prod: '第一个商品（草稿）', su_prodN: '店铺通过审核后上架',
    su_help: '遇到问题？联系 Phsar 客服', su_suspended: '该店铺已被暂停，请联系 Phsar 客服。', su_closed: '该店铺已关闭。',
    staffWait: '该店铺尚未开业，店主仍在设置中。',
    mg_sell: '多卖货', mg_buyers: '买家', mg_shop: '我的店铺', mg_help: '帮助和账户',
    mi_flash: '限时抢购', mi_flashD: '邀请和你的抢购商品', mi_ads: '广告', mi_adsD: '让商品排在前面',
    mi_returns: '退货', mi_returnsD: '同意或说明退货申请',
    mi_profile: '店铺资料', mi_profileD: '店名、标志、横幅、介绍', mi_kyc: '身份验证', mi_kycD: '证件审核状态',
    mi_pickup: '取件点', mi_pickupD: '快递员取件地址', mi_team: '团队', mi_teamD: '让员工协助处理订单',
    mi_help: '帮助', mi_helpD: '联系 Phsar 客服，查看你的求助', mi_stores: '我的店铺', mi_storesD: '你的店铺及你作为员工协助的店铺',
    mi_logout: '退出', mi_logoutD: '退出卖家中心', isNew: '新', langTitle: '语言'},
});

const V2 = {today: null};
const NAV_OF = {home: 'home', orders: 'orders', products: 'products', payouts: 'payouts'}; // selain ini → "Lainnya"
const V2_ICON = {
  home: '<path d="M9 11l3 3L22 4"/><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>',
  orders: '<path d="M21 8 12 3 3 8v8l9 5 9-5z"/><path d="M3 8l9 5 9-5"/><path d="M12 13v8"/>',
  products: '<rect x="3" y="3" width="7" height="7" rx="1.5"/><rect x="14" y="3" width="7" height="7" rx="1.5"/><rect x="3" y="14" width="7" height="7" rx="1.5"/><path d="M17.5 14v7M14 17.5h7"/>',
  payouts: '<rect x="2" y="6" width="20" height="13" rx="2.5"/><circle cx="12" cy="12.5" r="2.5"/><path d="M6 9.5v.01M18 15.5v.01"/>',
  more: '<circle cx="5" cy="12" r="1.6"/><circle cx="12" cy="12" r="1.6"/><circle cx="19" cy="12" r="1.6"/>',
};
const v2svg = k => `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">${V2_ICON[k]}</svg>`;
const v2Section = tab => NAV_OF[tab] || 'more';

// Pindah ke bagian/halaman lain di toko yang sedang dibuka (filter opsional untuk pesanan/produk/retur).
function v2Open(tab, filter) {
  if (filter && tab === 'orders') Object.assign(SO, {filter, q: '', page: 1});
  if (filter && tab === 'products') PF.filter = filter;
  if (tab === 'returns' && typeof SR !== 'undefined') SR.filter = 'open';
  if (tab === 'help') S.helpTicket = null;
  if (tab === 'chat') S.chatConv = null; // chat:
  S.tab = tab; S.editing = null;
  try { sessionStorage.setItem('phsar-tab-' + S.storeId, tab); } catch {}
  render();
  window.scrollTo(0, 0);
}

function storeShell(st, inner) {
  const owner = st.role === 'owner', d = V2.today || {};
  const secs = owner ? ['home', 'orders', 'products', 'payouts', 'more'] : ['home', 'orders', 'products', 'more'];
  const cur = v2Section(S.tab);
  const badge = {orders: (d.new_orders || 0) + (d.to_pack || 0) + (d.courier_returns || 0) /* courierreturn */, more: (d.returns || 0) + (owner ? d.flash_invites || 0 : 0) + (d.ticket_replies || 0) + (d.chat_unread || 0)}; // chat: + belum dibalas
  const links = secs.map(k => `<a href="#/store/${st.id}/${k}" ${cur === k ? 'aria-current="page"' : ''} onclick="event.preventDefault();v2Open('${k}')">
    <span class="ic">${v2svg(k)}${badge[k] ? `<span class="nb">${badge[k] > 99 ? '99+' : badge[k]}</span>` : ''}</span><span>${t('nav_' + k)}</span></a>`).join('');
  const ini = esc((st.name || '?').trim().slice(0, 2).toUpperCase());
  return `<header class="top v2top"><a class="storeid" href="#/store/${st.id}/home" onclick="event.preventDefault();v2Open('home')">
      ${st.logo_url ? `<img src="${esc(st.logo_url)}" alt="">` : `<span class="ini" aria-hidden="true">${ini}</span>`}
      <span><b>${esc(st.name)}</b><small>${esc(t('appName')).toUpperCase()}</small></span></a>
    <div class="top-actions">${bellButton()}${langSwitch()}</div></header>
    <div class="v2wrap"><nav class="v2nav" style="--n:${secs.length}" aria-label="${t('v2nav')}">${links}</nav><main class="v2main">${inner}</main></div>`;
}

// Layar fokus (editor produk): kepala toko tanpa navigasi bawah, supaya seller tidak terganggu.
function focusShell(st, inner) {
  const ini = esc((st.name || '?').trim().slice(0, 2).toUpperCase());
  return `<header class="top v2top"><a class="storeid" href="#/store/${st.id}/home">
      ${st.logo_url ? `<img src="${esc(st.logo_url)}" alt="">` : `<span class="ini" aria-hidden="true">${ini}</span>`}
      <span><b>${esc(st.name)}</b><small>${esc(t('appName')).toUpperCase()}</small></span></a>
    <div class="top-actions">${bellButton()}${langSwitch()}</div></header><main class="v2main focus">${inner}</main>`;
}

// Judul halaman: bagian utama punya judul sendiri; halaman dari "Lainnya" punya tombol kembali.
function v2Header(tab) {
  if (tab === 'home') return '';
  if (NAV_OF[tab] || tab === 'more') return `<div class="v2head"><h1>${t('h_' + tab)}</h1></div>`;
  const titles = {hours: 'mi_hours', returns: 'mi_returns', flash: 'mi_flash', ads: 'mi_ads', profile: 'mi_profile', kyc: 'mi_kyc', pickup: 'mi_pickup', team: 'mi_team', help: 'mi_help', faq: 'mi_faq', chat: 'mi_chat'}; // askbot: faq · chat:
  return `<a class="backlink" href="#/store/${S.storeId}/more" onclick="event.preventDefault();v2Open('more')">${t('backMore')}</a>
    <div class="v2head"><h1>${t(titles[tab] || 'h_more')}</h1></div>`;
}

// Lama waktu yang mudah dibaca: 90 → "1 h 30 min" (+ "jam buka" kalau hitungan hanya saat toko buka).
function v2Dur(m, withOpen) {
  const h = Math.floor(m / 60), r = m % 60;
  const s = h && r ? t('durHM', {h, m: r}) : h ? t('durH', {n: h}) : t('durM', {n: r});
  return withOpen && V2.today?.timing?.business_hours ? `${s} (${t('openHrs')})` : s;
}
function closedBanner() {
  const tm = V2.today?.timing;
  if (!tm || !tm.business_hours || tm.open_now) return '';
  return `<div class="banner info">🌙 ${tm.opens_at ? t('closedNow', {d: fmtTime(tm.opens_at)}) : t('closedNowNoTime')}</div>`;
}

function dueLeft(iso) {
  if (!iso) return '';
  const l = leftText(iso);
  return l ? t('leftT', {t: l}) : t('lateT');
}

// ---- Hari ini ----
async function tabHome(id) {
  const st = S.store, owner = st.role === 'owner';
  if (st.status !== 'active') {
    if (owner) return drawSetup(id);
    $('#tab').innerHTML = `<div class="banner info">${t('staffWait')}</div>`;
    return;
  }
  const d = V2.today || {};
  const pay = owner ? await api(`/v1/seller/stores/${id}/payouts`).catch(() => null) : null;
  const P_ = ['#d4063f', '#ffe4ec', '#a3002f'], O_ = ['#c25e00', '#fff0dc', '#8a4300'], V_ = ['#5b3cc4', '#ece7ff', '#45309a'], N_ = ['#102241', '#e6eaf2', '#102241'], G_ = ['#0b6b3c', '#e3f4ea', '#0b6b3c'];
  const tasks = [];
  const tm = d.timing || {accept_minutes: 120};
  const add = (n, key, icon, pal, go, due) => tasks.push({title: t('tk_' + key, {n}), why: t('tk_' + key + 'W', {a: v2Dur(tm.accept_minutes, true)}), btn: t('tk_' + key + 'B'), icon, pal, go, due});
  if (d.new_orders) add(d.new_orders, 'new', '📦', P_, "v2Open('orders','new')", dueLeft(d.new_due));
  if (d.to_pack) add(d.to_pack, 'pack', '🎁', O_, "v2Open('orders','processing')", dueLeft(d.pack_due));
  if (d.returns) add(d.returns, 'ret', '↩️', V_, "v2Open('returns')", dueLeft(d.returns_due));
  if (d.courier_returns && typeof crSellerOpen === 'function') add(d.courier_returns, 'crr', '📦', O_, 'crSellerOpen()', dueLeft(d.courier_returns_due)); // courierreturn
  if (d.product_fixes) add(d.product_fixes, 'fix', '✏️', P_, "v2Open('products','changes')", '');
  if (owner && d.flash_invites) add(d.flash_invites, 'flash', '⚡', N_, "v2Open('flash')", d.flash_starts ? t('startsT', {d: fmtTime(d.flash_starts)}) : '');
  if (d.chat_unread) add(d.chat_unread, 'chat', '💬', P_, "v2Open('chat')", ''); // chat: pembeli menunggu balasan
  if (d.ticket_replies) add(d.ticket_replies, 'tick', '💬', N_, "v2Open('help')", '');
  if (owner && pay && !pay.account) add(1, 'bank', '🏦', O_, "v2Open('payouts')", '');
  if (d.ready) add(d.ready, 'ready', '🛵', G_, "v2Open('orders','ready')", '');
  const card = k => `<article class="task" style="--tc:${k.pal[0]};--tb:${k.pal[1]};--td:${k.pal[2]}">
      <div class="tt"><span class="ti" aria-hidden="true">${k.icon}</span><div><h3>${esc(k.title)}</h3><p>${esc(k.why)}</p></div></div>
      <div class="tf">${k.due ? `<span class="due">⏱ ${esc(k.due)}</span>` : ''}<button type="button" class="primary" onclick="${k.go}">${esc(k.btn)}</button></div></article>`;
  const tips = t('tips');
  const tip = Array.isArray(tips) ? tips[new Date().getDate() % tips.length] : '';
  $('#tab').innerHTML = `
    <div class="v2head"><div><p class="sub">${t('hello', {n: esc(st.name)})}</p><h1>${t('todayTitle')}</h1></div></div>
    ${owner ? '' : `<div class="banner info">${t('staffHome')}</div>`}
    ${closedBanner()}
    <div class="spread"><h2>${t('doNow')}</h2></div>
    ${tasks.length ? tasks.map(card).join('') : `<div class="alldone">${t('allDoneNow')}</div>`}
    <div class="minis">
      <div class="mini"><span>${t('salesToday')}</span><b>${pmoney(d.sales_today || 0)}</b><small style="color:var(--muted)">${t('ordersN', {n: d.orders_today || 0})}</small></div>
      ${pay ? `<a class="mini" href="#/store/${id}/payouts" onclick="event.preventDefault();v2Open('payouts')"><span>${t('readyMoney')}</span><b>${pmoney(pay.available)}</b><small>${t('seeMoney')}</small></a>` : ''}
    </div>
    ${tip ? `<aside class="tip"><span aria-hidden="true">💡</span><div><b>${t('tipTitle')}</b>${esc(tip)}</div></aside>` : ''}`;
}

// ---- Mulai jualan (toko belum aktif) ----
async function drawSetup(id) {
  const st = S.store, c = st.checklist || {};
  if (st.status === 'suspended' || st.status === 'closed') {
    $('#tab').innerHTML = `<div class="banner bad">${t(st.status === 'suspended' ? 'su_suspended' : 'su_closed')}</div>
      <div><button type="button" class="primary" onclick="v2Open('help')">${t('su_help')}</button></div>`;
    return;
  }
  const sent = st.status === 'pending_review';
  const kycOk = c.kyc === 'pending' || c.kyc === 'verified';
  const steps = [
    {k: 'profile', done: c.logo && c.about, min: 2, go: "v2Open('profile')"},
    {k: 'kyc', done: kycOk, min: 3, go: "v2Open('kyc')", note: c.kyc && c.kyc !== 'missing' ? t('kyc_' + c.kyc) : ''},
    {k: 'pickup', done: c.pickup_points > 0, min: 3, go: "v2Open('pickup')"},
    {k: 'submit', done: sent, min: 1, go: 'setupSubmit()'},
    {k: 'review', done: false, wait: sent},
  ];
  const done = steps.filter(s => s.done).length;
  const next = steps.find(s => !s.done);
  let nextHTML;
  if (next.k === 'review') {
    nextHTML = `<section class="nextcard"><span class="lbl">${t('su_wait')}</span><h2>${t('su_waitT')}</h2><p>${t('su_waitW')}</p>
      <a class="btn primary" href="#/store/${id}/product/new">${t('su_waitC')}</a></section>`;
  } else {
    const can = next.k !== 'submit' || c.can_submit;
    nextHTML = `<section class="nextcard"><span class="lbl">${t('nextStep')}</span><h2>${t('su_' + next.k)}</h2><p>${t('su_' + next.k + 'W')}</p>
      ${next.min ? `<p>⏱ ${t('su_min', {n: next.min})}</p>` : ''}
      <button type="button" class="primary" id="suGo" ${can ? '' : 'disabled'}>${t('su_' + next.k + 'C')}</button></section>`;
  }
  const row = (s, i) => {
    const isNext = s === next;
    const chip = s.done ? `<span class="chip ok">${t('su_done')}</span>` : s.wait ? `<span class="chip warn">${t('su_wait')}</span>` : `<span class="chip ${isNext ? 'pink' : ''}">${t('su_todo')}</span>`;
    const inner = `<span class="n" aria-hidden="true">${s.done ? '✓' : i + 1}</span><span class="sx"><b>${t('su_' + s.k)}</b><span>${esc(s.note || t('su_' + s.k + 'N'))}</span></span>${chip}`;
    return s.go && s.k !== 'submit' ? `<a class="step ${s.done ? 'done' : isNext ? 'next' : ''}" href="#" onclick="event.preventDefault();${s.go}">${inner}</a>`
      : `<div class="step ${s.done ? 'done' : isNext ? 'next' : ''}">${inner}</div>`;
  };
  $('#tab').innerHTML = `
    <div class="v2head"><div class="ring" style="--p:${done * 20}%" role="img" aria-label="${done}/5"><span>${done}/5</span></div>
      <div><h1>${t('su_title')}</h1><p class="sub">${t('su_sub')}</p></div></div>
    ${st.review_note && !sent ? `<div class="banner bad"><b>${t('needsFix')}</b> ${esc(st.review_note)}</div>` : ''}
    ${nextHTML}
    <section class="card steps"><h2>${t('allSteps')}</h2>${steps.map(row).join('')}</section>
    <section class="card steps"><h2>${t('su_extra')}</h2>
      <a class="step" href="#" onclick="event.preventDefault();v2Open('payouts')"><span class="n" aria-hidden="true">🏦</span><span class="sx"><b>${t('su_bank')}</b><span>${t('su_bankN')}</span></span></a>
      <a class="step" href="#" onclick="event.preventDefault();v2Open('hours')"><span class="n" aria-hidden="true">🕗</span><span class="sx"><b>${t('su_hours')}</b><span>${t('su_hoursN')}</span></span></a>
      <a class="step" href="#/store/${id}/product/new"><span class="n" aria-hidden="true">📦</span><span class="sx"><b>${t('su_prod')}</b><span>${t('su_prodN')}</span></span></a></section>
    <p class="muted" style="text-align:center"><a href="#" onclick="event.preventDefault();v2Open('help')">${t('su_help')}</a></p>`;
  if ($('#suGo')) $('#suGo').onclick = () => { if (next.k === 'submit') setupSubmit(); else v2Open(next.k); };
}

async function setupSubmit() {
  try { await api(`/v1/seller/stores/${S.storeId}/submit`, {method: 'POST'}); render(); } catch (e) { toast(errMsg(e), 'error'); }
}

// ---- Lainnya ----
function tabMore(id) {
  const owner = S.store.role === 'owner', d = V2.today || {};
  const item = (tab, icon, bg, key, badge) => `<a class="mitem" href="#/store/${id}/${tab}" onclick="event.preventDefault();v2Open('${tab}')" style="--mb:${bg}">
      <span class="mi" aria-hidden="true">${icon}</span><span class="mx"><b>${t('mi_' + key)}${badge ? `<span class="mbadge ${badge === 'new' ? 'new' : ''}">${badge === 'new' ? t('isNew') : badge}</span>` : ''}</b><span>${t('mi_' + key + 'D')}</span></span>
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true" style="color:var(--muted)"><path d="m9 6 6 6-6 6"/></svg></a>`;
  const group = (title, items) => items.length ? `<section class="mgroup"><h2>${t(title)}</h2><div class="mlist">${items.join('')}</div></section>` : '';
  $('#tab').innerHTML = [
    owner ? group('mg_sell', [item('flash', '⚡', '#ffe4ec', 'flash', d.flash_invites || ''), item('ads', '📣', '#fff6d9', 'ads', '')]) : '',
    group('mg_buyers', [item('chat', '💬', '#ffe4ec', 'chat', d.chat_unread || ''), item('returns', '↩️', '#ece7ff', 'returns', d.returns || ''), item('faq', '💡', '#e3f4ea', 'faq')]), // chat: · askbot: FAQ
    owner ? group('mg_shop', [item('profile', '🏪', '#e3f4ea', 'profile'), item('kyc', '🪪', '#e6eaf2', 'kyc'), item('hours', '🕗', '#fff6d9', 'hours'), item('pickup', '📍', '#fff0dc', 'pickup'), item('team', '👥', '#e6eaf2', 'team')]) : '',
    `<section class="mgroup"><h2>${t('mg_help')}</h2><div class="mlist">${item('help', '❓', '#f1eef3', 'help', d.ticket_replies || '')}
      <a class="mitem" href="#/" style="--mb:#f1eef3"><span class="mi" aria-hidden="true">🏬</span><span class="mx"><b>${t('mi_stores')}</b><span>${t('mi_storesD')}</span></span></a>
      <a class="mitem" href="#" onclick="event.preventDefault();logout()" style="--mb:#fee2e2"><span class="mi" aria-hidden="true">🚪</span><span class="mx"><b>${t('mi_logout')}</b><span>${t('mi_logoutD')}</span></span></a></div></section>
    <section class="mgroup langrow"><div>${langSwitch()}</div></section>`,
  ].join('');
}

// ---- Jam buka toko ----
async function tabHours(id) {
  const v = await api(`/v1/seller/stores/${id}/hours`);
  const own = !!v.hours;
  const draw = mine => {
    $('#tab').innerHTML = `
      <section class="card"><div class="spread"><h2>${t('hrTitle')}</h2><span class="chip ${v.open_now ? 'ok' : ''}">${v.open_now ? t('hrOpenNow') : t('hrClosedNow')}</span></div>
        <p>${v.business_hours ? t('hrWhy') : t('hrWhyOff')}</p>
        <p class="ohelp">⏱ ${v.pack_enabled ? t('hrRules', {a: v2Dur(v.accept_minutes, v.business_hours), p: v2Dur(v.pack_minutes, v.business_hours)}) : t('hrRulesNoPack', {a: v2Dur(v.accept_minutes, v.business_hours)})}</p></section>
      <section class="card">
        ${mine ? `${hoursRows(v.hours || v.default_hours, 'hrMine')}
            <div class="row"><button type="button" class="primary" id="hrSave">${t('hrSave')}</button>${own ? `<button type="button" id="hrDef">${t('hrUseDefault')}</button>` : `<button type="button" id="hrCancel">${t('cancel')}</button>`}</div>`
          : `<p>${t('hrUsingDefault', {s: esc(hoursSummary(v.default_hours))})}</p><div class="row"><button type="button" class="primary" id="hrMineBtn">${t('hrUseMine')}</button></div>`}
      </section>`;
    if (mine) {
      bindHoursRows('hrMine');
      $('#hrSave').onclick = () => save(readHoursRows('hrMine'));
      if ($('#hrDef')) $('#hrDef').onclick = () => save(null, true);
      if ($('#hrCancel')) $('#hrCancel').onclick = () => draw(false);
    } else $('#hrMineBtn').onclick = () => draw(true);
  };
  const save = async (week, reset) => {
    if (!week && !reset) { toast(t('hrNeedOne'), 'warning'); return; }
    try { await api(`/v1/seller/stores/${id}/hours`, {method: 'PUT', body: {hours: week}}); toast(t('hrSaved')); render(); } catch (e) { toast(errMsg(e), 'error'); }
  };
  draw(own);
}
