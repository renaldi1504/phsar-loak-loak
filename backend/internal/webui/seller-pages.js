'use strict';
// sellerpages: blok tampilan bersama untuk halaman dari menu "Lainnya" (retur, flash sale, iklan, profil, KYC,
// titik jemput, tim, bantuan) — gaya v2: langkah 1-2-3, kosong-state ramah, hitung mundur tenggat.
addStrings({
  en: {spHow: 'How it works', spBack: 'Go back', spToday: 'Back to Today', spLeftAns: '{t} left to answer', spLeft: '{t} left',
    spRt1: 'See the request', spRt1W: 'Photos and reason', spRt2: 'Accept or explain', spRt2W: 'Within 48 hours', spRt3: 'Done', spRt3W: 'Phsar handles the refund',
    spRtHelp: 'Answer each request within 48 hours. If you don’t answer, Phsar decides for you.', spRtBuyerSays: 'Buyer says', spRtRejectBtn: 'Reject and explain',
    spRtNoneT: 'No returns to answer', spRtNoneW: 'When a buyer asks to return something, it shows up here and on your Today screen.', spRtSeeAll: 'See past returns',
    spRtNoneAllW: 'You have no returns yet. Keep it up!',
    spRtNext_approved: 'Phsar sends the refund to the buyer. Nothing to do.', spRtNext_disputed: 'Phsar is checking the photos and your reason. We will notify you.',
    spRtNext_rejected: 'Phsar agreed with you — no refund.', spRtNext_refunded: 'Done — the buyer got the money back.', spRtNext_cancelled: 'The buyer cancelled this request. Nothing to do.',
    spRtNext_late: 'Time is up — Phsar is deciding now.',
    spFl1: 'Pick a session', spFl1W: 'Check the dates', spFl2: 'Set price & quota', spFl2W: 'Lower price, how many', spFl3: 'Phsar checks', spFl3W: 'Then it goes live',
    spFlInvW: 'Phsar picked these products for you. Accept to join.', spFlNoneW: 'New sessions open often. We will notify you when one opens.',
    spFlRegLeft: 'Register within {t}', spFlClose: 'Close', spFlNoMine: 'Products you register show up here with their status.', spFlFindHint: 'Tick the products you want to add, then set the flash price and quota.',
    spFlPicked: 'Your products ({n})',
    spAd1: 'Choose a product', spAd1W: 'It must be live', spAd2: 'Set your budget', spAd2W: 'Per click and per day', spAd3: 'Phsar reviews', spAd3W: 'Within 1 working day',
    spAdBal: 'Ad balance', spAdBalW: 'Clicks are paid from your Phsar balance.', spAdChange: 'Change price or budget', spAdYours: 'Your campaigns',
    spAdNoneT: 'Get more buyers with ads', spAdNoneW: 'Your product shows at the top of search. You pay only when someone clicks.',
    spAdAll: 'All time: {v} views · {c} clicks · {s} spent · CTR {r}',
    spPfLook: 'How buyers see your shop', spPfName: 'Name and description', spPfSocialW: 'Optional — buyers can follow you here.', spPfLogo: 'Change logo',
    spKy_missing: 'Not sent yet', spKy_missingW: 'Fill in the 3 steps below. It takes about 3 minutes.',
    spKy_pending: 'We are checking', spKy_pendingW: 'Phsar usually checks within 1–2 working days. You can still make changes.',
    spKy_verified: 'Verified', spKy_verifiedW: 'Your identity is confirmed. Nothing else to do.',
    spKy_rejected: 'Please fix and send again', spKy_rejectedW: 'Read the note below, correct it and save again.',
    spKyFileW: 'All 4 corners visible, no glare. JPG, PNG or PDF.', spKyMore: 'Tax details (optional)',
    spPkNoneT: 'Where should the courier pick up?', spPkS1: 'Put the pin on your door', spPkS2: 'Address and contact', spPkS3: 'When can the courier come?', spPkS4: 'Extra (optional)',
    spPkDelQ: 'Delete this pickup point?',
    spTm1: 'Create a link', spTm1W: 'Tap the button below', spTm2: 'Send it', spTm2W: 'e.g. in Telegram', spTm3: 'They join', spTm3W: 'One tap on the link',
    spTmPeople: 'People in this store',
    spHpWrite: 'Write to Phsar', spHpNoneW: 'Your questions and our answers will show here.', spHpReplyPh: 'Write your message…', spHpMine: 'Your requests'},
  km: {spHow: 'របៀបដំណើរការ', spBack: 'ត្រឡប់ក្រោយ', spToday: 'ត្រឡប់ទៅថ្ងៃនេះ', spLeftAns: 'នៅសល់ {t} ដើម្បីឆ្លើយ', spLeft: 'នៅសល់ {t}',
    spRt1: 'មើលសំណើ', spRt1W: 'រូបថត និងមូលហេតុ', spRt2: 'ទទួល ឬពន្យល់', spRt2W: 'ក្នុង 48 ម៉ោង', spRt3: 'រួចរាល់', spRt3W: 'Phsar ចាត់ការសងប្រាក់',
    spRtHelp: 'សូមឆ្លើយសំណើនីមួយៗក្នុង 48 ម៉ោង។ បើមិនឆ្លើយ Phsar នឹងសម្រេចជំនួសអ្នក។', spRtBuyerSays: 'អ្នកទិញនិយាយថា', spRtRejectBtn: 'បដិសេធ និងពន្យល់',
    spRtNoneT: 'គ្មានការប្តូរវិញដែលត្រូវឆ្លើយទេ', spRtNoneW: 'ពេលអ្នកទិញស្នើប្តូរវិញ វានឹងបង្ហាញនៅទីនេះ និងនៅទំព័រថ្ងៃនេះ។', spRtSeeAll: 'មើលការប្តូរវិញពីមុន',
    spRtNoneAllW: 'អ្នកមិនទាន់មានការប្តូរវិញទេ។ ល្អណាស់!',
    spRtNext_approved: 'Phsar កំពុងសងប្រាក់ទៅអ្នកទិញ។ អ្នកមិនចាំបាច់ធ្វើអ្វីទេ។', spRtNext_disputed: 'Phsar កំពុងពិនិត្យរូបថត និងមូលហេតុរបស់អ្នក។ យើងនឹងជូនដំណឹង។',
    spRtNext_rejected: 'Phsar យល់ស្របនឹងអ្នក — មិនសងប្រាក់ទេ។', spRtNext_refunded: 'រួចរាល់ — អ្នកទិញបានទទួលប្រាក់វិញ។', spRtNext_cancelled: 'អ្នកទិញបានបោះបង់សំណើនេះ។ មិនចាំបាច់ធ្វើអ្វីទេ។',
    spRtNext_late: 'ផុតពេលហើយ — Phsar កំពុងសម្រេច។',
    spFl1: 'ជ្រើសវគ្គ', spFl1W: 'មើលកាលបរិច្ឆេទ', spFl2: 'ដាក់តម្លៃ និងកូតា', spFl2W: 'តម្លៃទាប ចំនួនប៉ុន្មាន', spFl3: 'Phsar ពិនិត្យ', spFl3W: 'បន្ទាប់មកចាប់ផ្តើមលក់',
    spFlInvW: 'Phsar បានជ្រើសផលិតផលទាំងនេះឱ្យអ្នក។ ចុចទទួលដើម្បីចូលរួម។', spFlNoneW: 'វគ្គថ្មីបើកជាញឹកញាប់។ យើងនឹងជូនដំណឹងពេលមានវគ្គបើក។',
    spFlRegLeft: 'ចុះឈ្មោះក្នុង {t}', spFlClose: 'បិទ', spFlNoMine: 'ផលិតផលដែលអ្នកចុះឈ្មោះនឹងបង្ហាញនៅទីនេះ ជាមួយស្ថានភាព។', spFlFindHint: 'ធីកផលិតផលដែលចង់បន្ថែម បន្ទាប់មកដាក់តម្លៃ flash និងកូតា។',
    spFlPicked: 'ផលិតផលរបស់អ្នក ({n})',
    spAd1: 'ជ្រើសផលិតផល', spAd1W: 'ត្រូវតែកំពុងលក់', spAd2: 'កំណត់ថវិកា', spAd2W: 'ក្នុងមួយចុច និងមួយថ្ងៃ', spAd3: 'Phsar ពិនិត្យ', spAd3W: 'ក្នុង 1 ថ្ងៃធ្វើការ',
    spAdBal: 'សមតុល្យសម្រាប់ផ្សាយ', spAdBalW: 'ការចុចត្រូវបង់ពីសមតុល្យ Phsar របស់អ្នក។', spAdChange: 'ប្តូរតម្លៃ ឬថវិកា', spAdYours: 'យុទ្ធនាការរបស់អ្នក',
    spAdNoneT: 'ទទួលអ្នកទិញបន្ថែមជាមួយការផ្សាយ', spAdNoneW: 'ផលិតផលរបស់អ្នកបង្ហាញនៅខាងលើការស្វែងរក។ អ្នកបង់តែពេលមានគេចុច។',
    spAdAll: 'សរុប: {v} ដងមើល · {c} ចុច · ចំណាយ {s} · CTR {r}',
    spPfLook: 'អ្នកទិញឃើញហាងរបស់អ្នកបែបនេះ', spPfName: 'ឈ្មោះ និងការពិពណ៌នា', spPfSocialW: 'មិនចាំបាច់ — អ្នកទិញអាចតាមដានអ្នកនៅទីនេះ។', spPfLogo: 'ប្តូរឡូហ្គោ',
    spKy_missing: 'មិនទាន់ផ្ញើ', spKy_missingW: 'បំពេញ 3 ជំហានខាងក្រោម។ ប្រហែល 3 នាទី។',
    spKy_pending: 'យើងកំពុងពិនិត្យ', spKy_pendingW: 'Phsar ជាធម្មតាពិនិត្យក្នុង 1–2 ថ្ងៃធ្វើការ។ អ្នកនៅតែអាចកែបាន។',
    spKy_verified: 'បានផ្ទៀងផ្ទាត់', spKy_verifiedW: 'អត្តសញ្ញាណរបស់អ្នកត្រូវបានបញ្ជាក់។ មិនចាំបាច់ធ្វើអ្វីទៀតទេ។',
    spKy_rejected: 'សូមកែ ហើយផ្ញើម្តងទៀត', spKy_rejectedW: 'អានកំណត់ចំណាំខាងក្រោម កែវា ហើយរក្សាទុកម្តងទៀត។',
    spKyFileW: 'ឃើញជ្រុងទាំង 4 មិនចាំងពន្លឺ។ JPG, PNG ឬ PDF។', spKyMore: 'ព័ត៌មានពន្ធ (មិនចាំបាច់)',
    spPkNoneT: 'តើអ្នកដឹកគួរមកទទួលនៅឯណា?', spPkS1: 'ដាក់ម្ជុលត្រង់ទ្វាររបស់អ្នក', spPkS2: 'អាសយដ្ឋាន និងទំនាក់ទំនង', spPkS3: 'តើអ្នកដឹកអាចមកពេលណា?', spPkS4: 'បន្ថែម (មិនចាំបាច់)',
    spPkDelQ: 'លុបទីតាំងទទួលនេះ?',
    spTm1: 'បង្កើតតំណ', spTm1W: 'ចុចប៊ូតុងខាងក្រោម', spTm2: 'ផ្ញើវា', spTm2W: 'ឧ. ក្នុង Telegram', spTm3: 'ពួកគេចូលរួម', spTm3W: 'ចុចតំណម្តង',
    spTmPeople: 'មនុស្សក្នុងហាងនេះ',
    spHpWrite: 'សរសេរទៅ Phsar', spHpNoneW: 'សំណួររបស់អ្នក និងចម្លើយរបស់យើងនឹងបង្ហាញនៅទីនេះ។', spHpReplyPh: 'សរសេរសាររបស់អ្នក…', spHpMine: 'សំណើរបស់អ្នក'},
  zh: {spHow: '怎么做', spBack: '返回', spToday: '回到今日', spLeftAns: '还剩 {t} 回复', spLeft: '还剩 {t}',
    spRt1: '查看申请', spRt1W: '照片和原因', spRt2: '同意或说明', spRt2W: '48 小时内', spRt3: '完成', spRt3W: 'Phsar 处理退款',
    spRtHelp: '每个申请请在 48 小时内回复。如果不回复，Phsar 会替你决定。', spRtBuyerSays: '买家说', spRtRejectBtn: '拒绝并说明',
    spRtNoneT: '没有需要回复的退货', spRtNoneW: '买家申请退货时，会显示在这里和“今日”页面。', spRtSeeAll: '查看以往退货',
    spRtNoneAllW: '你还没有任何退货。继续保持！',
    spRtNext_approved: 'Phsar 正在给买家退款，你无需操作。', spRtNext_disputed: 'Phsar 正在查看照片和你的理由，结果会通知你。',
    spRtNext_rejected: 'Phsar 同意你的意见——不退款。', spRtNext_refunded: '已完成——买家已收到退款。', spRtNext_cancelled: '买家已取消此申请，你无需操作。',
    spRtNext_late: '已超时——Phsar 正在决定。',
    spFl1: '选择场次', spFl1W: '看好日期', spFl2: '设置价格和数量', spFl2W: '更低价格，卖多少', spFl3: 'Phsar 审核', spFl3W: '然后开始售卖',
    spFlInvW: 'Phsar 为你挑选了这些商品，接受即可参加。', spFlNoneW: '新场次经常开放，开放时我们会通知你。',
    spFlRegLeft: '请在 {t} 内报名', spFlClose: '收起', spFlNoMine: '你报名的商品会和状态一起显示在这里。', spFlFindHint: '勾选要加入的商品，然后设置秒杀价和数量。',
    spFlPicked: '你的商品（{n}）',
    spAd1: '选择商品', spAd1W: '必须在售', spAd2: '设置预算', spAd2W: '每次点击和每天', spAd3: 'Phsar 审核', spAd3W: '1 个工作日内',
    spAdBal: '广告余额', spAdBalW: '点击费用从你的 Phsar 余额中扣除。', spAdChange: '修改出价或预算', spAdYours: '你的推广',
    spAdNoneT: '用广告吸引更多买家', spAdNoneW: '你的商品显示在搜索结果顶部，只有被点击才付费。',
    spAdAll: '累计：{v} 次浏览 · {c} 次点击 · 花费 {s} · 点击率 {r}',
    spPfLook: '买家看到的店铺', spPfName: '名称和简介', spPfSocialW: '可选——买家可以在这里关注你。', spPfLogo: '更换标志',
    spKy_missing: '尚未提交', spKy_missingW: '完成下面 3 步，大约 3 分钟。',
    spKy_pending: '审核中', spKy_pendingW: 'Phsar 通常在 1–2 个工作日内审核。你仍可修改。',
    spKy_verified: '已认证', spKy_verifiedW: '你的身份已确认，无需其他操作。',
    spKy_rejected: '请修改后重新提交', spKy_rejectedW: '请阅读下方说明，修改后再次保存。',
    spKyFileW: '四个角都清楚，无反光。JPG、PNG 或 PDF。', spKyMore: '税务信息（可选）',
    spPkNoneT: '快递员应在哪里取件？', spPkS1: '把图钉放在你的店门口', spPkS2: '地址和联系方式', spPkS3: '快递员什么时候可以来？', spPkS4: '其他（可选）',
    spPkDelQ: '删除这个取件点？',
    spTm1: '创建链接', spTm1W: '点击下方按钮', spTm2: '发送链接', spTm2W: '例如通过 Telegram', spTm3: '对方加入', spTm3W: '点一下链接即可',
    spTmPeople: '店铺成员',
    spHpWrite: '联系 Phsar', spHpNoneW: '你的问题和我们的回复会显示在这里。', spHpReplyPh: '输入你的消息…', spHpMine: '你的请求'},
});

// Langkah 1-2-3: [[judul, keterangan], …]; now = indeks langkah yang sedang dikerjakan (-1 = hanya penjelasan).
function spSteps(steps, now = -1) {
  return `<ol class="sp-steps" aria-label="${t('spHow')}">${steps.map(([a, b], i) => `<li class="${i === now ? 'now' : i < now ? 'done' : ''}">
    <span class="n" aria-hidden="true">${i < now ? '✓' : i + 1}</span><b>${a}</b>${b ? `<small>${b}</small>` : ''}</li>`).join('')}</ol>`;
}

// Kosong-state ramah: ikon besar, judul, penjelasan, tombol ajakan (HTML) opsional.
function spEmpty(icon, title, text, action = '') {
  return `<section class="sp-empty"><span class="sp-eic" aria-hidden="true">${icon}</span><h2>${title}</h2>${text ? `<p>${text}</p>` : ''}${action}</section>`;
}

// Sisa menit sampai iso (0 kalau lewat).
function spMinLeft(iso) { return Math.max(0, Math.floor((new Date(iso) - Date.now()) / 60000)); }

// Hitung mundur tenggat dengan bar: total = lama tenggat penuh (menit), mis. 48 jam untuk retur.
function spDue(iso, total, key = 'spLeftAns', lateKey = 'rtLate') {
  return `<div class="sp-due" data-spdue="${esc(iso)}" data-sptotal="${total}" data-spkey="${key}" data-splate="${lateKey}">${spDueInner(iso, total, key, lateKey)}</div>`;
}
function spDueInner(iso, total, key, lateKey) {
  const m = spMinLeft(iso), pct = Math.max(0, Math.min(100, Math.round(100 * m / total)));
  const lvl = !m ? 'late' : m < total / 4 ? 'hot' : '';
  return `<span class="sp-dt ${lvl}">⏱ ${m ? t(key, {t: v2Dur(m)}) : t(lateKey)}</span><span class="sp-bar ${lvl}"><i style="width:${pct}%"></i></span>`;
}
// Perbarui semua hitung mundur di layar (dipanggil timer halaman).
function spDueTick() {
  document.querySelectorAll('[data-spdue]').forEach(el => { el.innerHTML = spDueInner(el.dataset.spdue, +el.dataset.sptotal, el.dataset.spkey, el.dataset.splate); });
}
