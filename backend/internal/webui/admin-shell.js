'use strict';
// Kerangka admin (menu samping berkelompok, kepala halaman), dashboard, dan audit log — mengikuti "Mock-up Admin Phsar".
addStrings({
  en: {signInTitle: 'Sign in to Phsar Admin', nav: 'Admin menu', env: 'Internal', soon: 'Soon', phase2: 'Phase 2', mfaOn: 'MFA on', menu: 'Menu',
    g_overview: 'Overview', g_ops: 'Operations', g_fin: 'Finance', g_growth: 'Growth', g_sys: 'System',
    i_dashboard: 'Dashboard', i_sellers: 'Sellers & KYC', i_moderation: 'Product moderation', i_reports: 'Reports & strikes', i_orders: 'Orders & shipping', i_returns: 'Returns & disputes',
    i_support: 'Support tickets', i_customers: 'Customers', i_payments: 'Payments & webhooks', i_reconciliation: 'Bank reconciliation', i_payouts: 'Seller payouts', i_sweep: 'Commission sweep',
    i_refunds: 'Refunds', i_invoices: 'Invoices & tax', i_commission: 'Commission rules', i_ads: 'Ads', i_promos: 'Vouchers & flash sales', i_catalog: 'Categories & brands',
    i_content: 'Content & SEO', i_shipping: 'Shipping & FX rates', i_admins: 'Admins & roles', i_audit: 'Audit log', i_jobs: 'Scheduled jobs', i_ordertiming: 'Order time limits', i_legal: 'Legal documents',
    step: 'Step {n} of 2', note1: 'For the Phsar team only. Sellers and customers use the app.', note2: 'Sessions end after 8 hours or 30 minutes without activity. After 5 wrong codes the session is closed.',
    dashTitle: 'Dashboard', kpiLabel: 'Key figures', k_gmv: 'GMV, 7 days', k_gmvSub: '{n} paid orders', k_comm: 'Commission, 7 days', k_commSub: 'before gateway fees',
    k_sellers: 'Active sellers', k_sellersSub: 'of {n} registered', k_transfer: 'Payouts to transfer', k_transferSub: '{n} stores waiting',
    todo: 'Needs action', todoSub: 'Oldest first', open: 'Open', activity: 'Recent admin activity', noActivity: 'No admin activity yet.', allClear: 'All queues are clear',
    allClearSub: 'Nothing is waiting for your role right now.', alertFailed: '{n} failed deliveries', alertFailedSub: 'Call the buyer or return the parcel to the store.',
    alertRefund: 'Refunds waiting more than 24 h', alertRefundSub: 'Refund in ABA PayWay and record the reference.', oldest: 'Oldest {t}', nothing: 'Nothing waiting',
    q_stores: 'Stores awaiting verification', q_products: 'Products awaiting review', q_refunds: 'Refunds to make', q_disputes: 'Return disputes to decide', q_return_refunds: 'Approved returns to refund', q_payouts: 'Payouts to transfer',
    q_deliveries: 'Parcels waiting for a courier', q_failed: 'Failed deliveries', q_pickups: 'New pickup points to verify',
    r_super_admin: 'Super admin', r_moderator: 'Moderator', r_cs: 'Customer service', r_finance: 'Finance', r_ads_reviewer: 'Ads reviewer', r_system: 'System',
    ago_m: '{n} min', ago_h: '{n} h', ago_d: '{n} d',
    auditTitle: 'Audit log', auditDesc: 'Every data change and every sensitive read. Append-only: records can’t be edited or deleted.',
    f_actor: 'Actor', f_table: 'Table', f_action: 'Action', all: 'All', at_admin: 'Admins', at_seller: 'Sellers', at_buyer: 'Buyers', at_courier: 'Couriers', at_system: 'System',
    c_time: 'Time', c_actor: 'Actor', c_action: 'Action', c_record: 'Record', recordDetail: 'Record detail', pickRecord: 'Choose a record to see the details.',
    device: 'Device', field: 'Field', before: 'Before', after: 'After', readOnly: 'Read only · no data changed', sensitive: 'Sensitive read',
    secretNote: 'OTP codes, tokens and encrypted values are never shown in the log.', loadMore: 'Load more',
    a_insert: 'Created', a_update: 'Updated', a_soft_delete: 'Deleted', a_delete: 'Deleted permanently', a_restore: 'Restored',
    'a_admin.login': 'Signed in to admin', 'a_admin.login_backup_code': 'Signed in with a backup code', 'a_admin.mfa_enabled': 'Turned on MFA', 'a_admin.mfa_reset': 'Reset MFA of an admin',
    'a_admin.backup_codes_renewed': 'Created new backup codes', 'a_kyc_document.view': 'Viewed an ID document', 'a_payout_account.view': 'Viewed bank account numbers',
    tabUsers: 'Admin users', tabRoles: 'Role permissions', perm: 'Permission', rule1: 'Every admin signs in with a phone number and an authenticator app (MFA).',
    rule2: 'Keep at least two super admins, so one can reset MFA if another loses their phone.', rule3: 'Deactivate admins who leave the team — their sessions end immediately.'},
  km: {signInTitle: 'ចូលទៅ Phsar Admin', nav: 'ម៉ឺនុយអ្នកគ្រប់គ្រង', env: 'ផ្ទៃក្នុង', soon: 'ឆាប់ៗ', phase2: 'ដំណាក់កាលទី 2', mfaOn: 'MFA បើក', menu: 'ម៉ឺនុយ',
    g_overview: 'ទិដ្ឋភាពទូទៅ', g_ops: 'ប្រតិបត្តិការ', g_fin: 'ហិរញ្ញវត្ថុ', g_growth: 'កំណើន', g_sys: 'ប្រព័ន្ធ',
    i_dashboard: 'ផ្ទាំងគ្រប់គ្រង', i_sellers: 'អ្នកលក់ និង KYC', i_moderation: 'ត្រួតពិនិត្យផលិតផល', i_reports: 'របាយការណ៍ និងការព្រមាន', i_orders: 'ការបញ្ជាទិញ និងការដឹកជញ្ជូន', i_returns: 'ការប្រគល់វិញ និងវិវាទ',
    i_support: 'សំបុត្រជំនួយ', i_customers: 'អតិថិជន', i_payments: 'ការទូទាត់ និង webhook', i_reconciliation: 'ការផ្ទៀងផ្ទាត់ធនាគារ', i_payouts: 'ការទូទាត់ឱ្យអ្នកលក់', i_sweep: 'ការផ្ទេរកម្រៃជើងសារ',
    i_refunds: 'ការសងប្រាក់វិញ', i_invoices: 'វិក្កយបត្រ និងពន្ធ', i_commission: 'ច្បាប់កម្រៃជើងសារ', i_ads: 'ការផ្សាយពាណិជ្ជកម្ម', i_promos: 'ប័ណ្ណបញ្ចុះតម្លៃ និង flash sale', i_catalog: 'ប្រភេទ និងម៉ាក',
    i_content: 'មាតិកា និង SEO', i_shipping: 'ថ្លៃដឹក និងអត្រាប្តូរប្រាក់', i_admins: 'អ្នកគ្រប់គ្រង និងតួនាទី', i_audit: 'កំណត់ហេតុសវនកម្ម', i_jobs: 'ការងារកំណត់ពេល', i_ordertiming: 'ពេលកំណត់ការបញ្ជាទិញ', i_legal: 'ឯកសារច្បាប់',
    step: 'ជំហាន {n} នៃ 2', note1: 'សម្រាប់តែក្រុម Phsar ប៉ុណ្ណោះ។ អ្នកលក់ និងអតិថិជនប្រើកម្មវិធី។', note2: 'វគ្គបញ្ចប់បន្ទាប់ពី 8 ម៉ោង ឬ 30 នាទីដោយមិនប្រើ។ លេខកូដខុស 5 ដង វគ្គត្រូវបិទ។',
    dashTitle: 'ផ្ទាំងគ្រប់គ្រង', kpiLabel: 'តួលេខសំខាន់', k_gmv: 'GMV ៧ ថ្ងៃ', k_gmvSub: 'ការបញ្ជាទិញបានបង់ {n}', k_comm: 'កម្រៃជើងសារ ៧ ថ្ងៃ', k_commSub: 'មុនកាត់ថ្លៃ gateway',
    k_sellers: 'អ្នកលក់សកម្ម', k_sellersSub: 'ក្នុងចំណោម {n} ដែលបានចុះឈ្មោះ', k_transfer: 'ការទូទាត់ត្រូវផ្ទេរ', k_transferSub: 'ហាង {n} កំពុងរង់ចាំ',
    todo: 'ត្រូវចាត់វិធានការ', todoSub: 'ចាស់ជាងគេមុន', open: 'បើក', activity: 'សកម្មភាពអ្នកគ្រប់គ្រងថ្មីៗ', noActivity: 'មិនទាន់មានសកម្មភាពទេ។', allClear: 'គ្មានការងាររង់ចាំ',
    allClearSub: 'មិនមានអ្វីរង់ចាំតួនាទីរបស់អ្នកទេ។', alertFailed: 'ការដឹកជញ្ជូនបរាជ័យ {n}', alertFailedSub: 'ទូរស័ព្ទទៅអ្នកទិញ ឬប្រគល់កញ្ចប់ទៅហាងវិញ។',
    alertRefund: 'ការសងប្រាក់រង់ចាំលើស 24 ម៉ោង', alertRefundSub: 'សងប្រាក់ក្នុង ABA PayWay ហើយកត់លេខយោង។', oldest: 'ចាស់បំផុត {t}', nothing: 'គ្មានអ្វីរង់ចាំ',
    q_stores: 'ហាងរង់ចាំការផ្ទៀងផ្ទាត់', q_products: 'ផលិតផលរង់ចាំការត្រួតពិនិត្យ', q_refunds: 'ការសងប្រាក់ត្រូវធ្វើ', q_disputes: 'វិវាទការប្ដូរត្រូវសម្រេច', q_return_refunds: 'ការប្ដូរបានយល់ព្រម ត្រូវសងប្រាក់', q_payouts: 'ការទូទាត់ត្រូវផ្ទេរ',
    q_deliveries: 'កញ្ចប់រង់ចាំអ្នកដឹក', q_failed: 'ការដឹកជញ្ជូនបរាជ័យ', q_pickups: 'ទីតាំងទទួលថ្មីត្រូវផ្ទៀងផ្ទាត់',
    r_super_admin: 'អ្នកគ្រប់គ្រងកំពូល', r_moderator: 'អ្នកត្រួតពិនិត្យ', r_cs: 'សេវាអតិថិជន', r_finance: 'ហិរញ្ញវត្ថុ', r_ads_reviewer: 'អ្នកពិនិត្យការផ្សាយ', r_system: 'ប្រព័ន្ធ',
    ago_m: '{n} នាទី', ago_h: '{n} ម៉ោង', ago_d: '{n} ថ្ងៃ',
    auditTitle: 'កំណត់ហេតុសវនកម្ម', auditDesc: 'រាល់ការផ្លាស់ប្តូរទិន្នន័យ និងការមើលទិន្នន័យរសើប។ កំណត់ត្រាអាចបន្ថែមបានតែប៉ុណ្ណោះ មិនអាចកែ ឬលុបបានទេ។',
    f_actor: 'អ្នកធ្វើ', f_table: 'តារាង', f_action: 'សកម្មភាព', all: 'ទាំងអស់', at_admin: 'អ្នកគ្រប់គ្រង', at_seller: 'អ្នកលក់', at_buyer: 'អ្នកទិញ', at_courier: 'អ្នកដឹក', at_system: 'ប្រព័ន្ធ',
    c_time: 'ពេលវេលា', c_actor: 'អ្នកធ្វើ', c_action: 'សកម្មភាព', c_record: 'កំណត់ត្រា', recordDetail: 'ព័ត៌មានលម្អិតកំណត់ត្រា', pickRecord: 'ជ្រើសកំណត់ត្រាដើម្បីមើលលម្អិត។',
    device: 'ឧបករណ៍', field: 'វាល', before: 'មុន', after: 'ក្រោយ', readOnly: 'មើលតែប៉ុណ្ណោះ · គ្មានទិន្នន័យផ្លាស់ប្តូរ', sensitive: 'ការមើលទិន្នន័យរសើប',
    secretNote: 'លេខកូដ OTP token និងតម្លៃដែលបានអ៊ិនគ្រីប មិនដែលបង្ហាញក្នុងកំណត់ហេតុទេ។', loadMore: 'មើលបន្ថែម',
    a_insert: 'បានបង្កើត', a_update: 'បានកែប្រែ', a_soft_delete: 'បានលុប', a_delete: 'បានលុបជាអចិន្ត្រៃយ៍', a_restore: 'បានស្តារ',
    'a_admin.login': 'បានចូលអ្នកគ្រប់គ្រង', 'a_admin.login_backup_code': 'បានចូលដោយលេខកូដបម្រុង', 'a_admin.mfa_enabled': 'បានបើក MFA', 'a_admin.mfa_reset': 'បានកំណត់ MFA អ្នកគ្រប់គ្រងឡើងវិញ',
    'a_admin.backup_codes_renewed': 'បានបង្កើតលេខកូដបម្រុងថ្មី', 'a_kyc_document.view': 'បានមើលឯកសារអត្តសញ្ញាណ', 'a_payout_account.view': 'បានមើលលេខគណនីធនាគារ',
    tabUsers: 'អ្នកគ្រប់គ្រង', tabRoles: 'សិទ្ធិតាមតួនាទី', perm: 'សិទ្ធិ', rule1: 'អ្នកគ្រប់គ្រងទាំងអស់ចូលដោយលេខទូរស័ព្ទ និងកម្មវិធី authenticator (MFA)។',
    rule2: 'រក្សាអ្នកគ្រប់គ្រងកំពូលយ៉ាងហោចណាស់ពីរនាក់ ដើម្បីកំណត់ MFA ឡើងវិញពេលម្នាក់បាត់ទូរស័ព្ទ។', rule3: 'បិទអ្នកគ្រប់គ្រងដែលចាកចេញពីក្រុម — វគ្គរបស់គាត់បញ្ចប់ភ្លាមៗ។'},
  zh: {signInTitle: '登录 Phsar 管理后台', nav: '管理菜单', env: '内部', soon: '即将推出', phase2: '第二阶段', mfaOn: '已开启 MFA', menu: '菜单',
    g_overview: '概览', g_ops: '运营', g_fin: '财务', g_growth: '增长', g_sys: '系统',
    i_dashboard: '仪表盘', i_sellers: '卖家与 KYC', i_moderation: '商品审核', i_reports: '举报与违规', i_orders: '订单与配送', i_returns: '退货与纠纷',
    i_support: '客服工单', i_customers: '客户', i_payments: '支付与 Webhook', i_reconciliation: '银行对账', i_payouts: '卖家结算', i_sweep: '佣金划转',
    i_refunds: '退款', i_invoices: '发票与税务', i_commission: '佣金规则', i_ads: '广告', i_promos: '优惠券与限时抢购', i_catalog: '类目与品牌',
    i_content: '内容与 SEO', i_shipping: '运费与汇率', i_admins: '管理员与角色', i_audit: '审计日志', i_jobs: '定时任务', i_ordertiming: '订单时限', i_legal: '法律文件',
    step: '第 {n} 步，共 2 步', note1: '仅限 Phsar 团队。卖家和客户请使用 App。', note2: '会话 8 小时后或 30 分钟无操作后结束。验证码错误 5 次后会话关闭。',
    dashTitle: '仪表盘', kpiLabel: '关键数据', k_gmv: '7 天 GMV', k_gmvSub: '{n} 个已付款订单', k_comm: '7 天佣金', k_commSub: '未扣支付网关费',
    k_sellers: '活跃卖家', k_sellersSub: '共 {n} 家已注册', k_transfer: '待转账结算', k_transferSub: '{n} 家店铺等待',
    todo: '待处理', todoSub: '最早的在前', open: '打开', activity: '最近管理员操作', noActivity: '暂无管理员操作。', allClear: '没有待处理事项',
    allClearSub: '目前没有需要你的角色处理的事项。', alertFailed: '{n} 个配送失败', alertFailedSub: '联系买家或将包裹退回店铺。',
    alertRefund: '有退款等待超过 24 小时', alertRefundSub: '请在 ABA PayWay 退款并记录参考号。', oldest: '最早 {t}', nothing: '无待处理',
    q_stores: '待审核店铺', q_products: '待审核商品', q_refunds: '待退款', q_disputes: '待裁决的退货纠纷', q_return_refunds: '已同意待退款的退货', q_payouts: '待转账结算',
    q_deliveries: '等待快递员的包裹', q_failed: '配送失败', q_pickups: '待核实的新取件点',
    r_super_admin: '超级管理员', r_moderator: '审核员', r_cs: '客服', r_finance: '财务', r_ads_reviewer: '广告审核员', r_system: '系统',
    ago_m: '{n} 分钟', ago_h: '{n} 小时', ago_d: '{n} 天',
    auditTitle: '审计日志', auditDesc: '所有数据变更和敏感数据读取。只能追加：记录不能编辑或删除。',
    f_actor: '操作者', f_table: '数据表', f_action: '操作', all: '全部', at_admin: '管理员', at_seller: '卖家', at_buyer: '买家', at_courier: '快递员', at_system: '系统',
    c_time: '时间', c_actor: '操作者', c_action: '操作', c_record: '记录', recordDetail: '记录详情', pickRecord: '选择一条记录查看详情。',
    device: '设备', field: '字段', before: '修改前', after: '修改后', readOnly: '只读 · 未修改数据', sensitive: '敏感读取',
    secretNote: '验证码、令牌和加密值永远不会显示在日志中。', loadMore: '加载更多',
    a_insert: '创建', a_update: '修改', a_soft_delete: '删除', a_delete: '永久删除', a_restore: '恢复',
    'a_admin.login': '登录管理后台', 'a_admin.login_backup_code': '使用备用码登录', 'a_admin.mfa_enabled': '开启 MFA', 'a_admin.mfa_reset': '重置管理员 MFA',
    'a_admin.backup_codes_renewed': '生成新备用码', 'a_kyc_document.view': '查看身份证件', 'a_payout_account.view': '查看银行账号',
    tabUsers: '管理员', tabRoles: '角色权限', perm: '权限', rule1: '所有管理员使用手机号和验证器应用（MFA）登录。',
    rule2: '至少保留两名超级管理员，以便有人丢失手机时可以重置 MFA。', rule3: '停用离开团队的管理员——其会话会立即结束。'},
});

// Menu: href = halaman yang sudah ada; tanpa href = "Segera". perm = izin yang dibutuhkan.
const NAV = [
  ['overview', [['dashboard', '#/', null]]],
  ['ops', [['sellers', '#/stores', 'seller.verify', 'stores'], ['moderation', '#/products', 'product.moderate', 'products'], ['reports', '#/reports', 'report.handle', 'reports'],
    ['orders', '#/delivery', 'delivery.dispatch', 'deliveries+failed'], ['returns', '#/returns', 'order.view', 'disputes+return_refunds'], ['support', '#/support', 'support.handle', 'tickets'], ['chats', '#/chats', 'support.handle'] /* chat */, ['customers', '#/customers', 'customer.view']]],
  ['fin', [['payments', '#/payments', 'order.view'], ['reconciliation', '#/reconciliation', 'reconciliation.run'], ['payouts', '#/payouts', 'payout.approve', 'payouts'],
    ['sweep', '#/sweep', 'payout.approve'], ['refunds', '#/refunds', 'refund.decide', 'refunds'], ['invoices', '#/invoices', 'commission.manage'], ['commission', '#/commission', 'commission.manage']]],
  ['growth', [['ads', '#/ads', 'ads.review'], ['promos', '#/promos', 'commission.manage'], ['catalog', '#/catalog', 'catalog.manage'], ['content', '#/content', '*'], ['shipping', '#/shipping', 'commission.manage']]],
  ['sys', [['admins', '#/admins', 'admin.manage'], ['integrations', '#/integrations/ai', 'integration.configure', 'integrations'], ['audit', '#/audit', 'audit.view'], ['jobs', '#/jobs', '*'], ['ordertiming', '#/order-timing', 'order.rules'], ['legal', '#/legal', 'legal.manage']]],
];
const ROUTE_ITEM = {'#/': 'dashboard', '#/stores': 'sellers', '#/store': 'sellers', '#/products': 'moderation', '#/product': 'moderation', '#/delivery': 'orders',
  '#/payouts': 'payouts', '#/refunds': 'refunds', '#/admins': 'admins', '#/audit': 'audit', '#/gateway': 'integrations', '#/logins': 'integrations', '#/integrations': 'integrations', '#/commission': 'commission', '#/legal': 'legal', '#/returns': 'returns', '#/reports': 'reports', '#/support': 'support', '#/chats': 'chats', '#/reconciliation': 'reconciliation', '#/sweep': 'sweep', '#/invoices': 'invoices', '#/ads': 'ads', '#/promos': 'promos', '#/catalog': 'catalog', '#/content': 'content', '#/shipping': 'shipping', '#/jobs': 'jobs', '#/customers': 'customers', '#/order/': 'customers', '#/payments': 'payments', '#/order-timing': 'ordertiming'};
const SH = {badges: {}, at: 0, open: false};

function hasPerm(p) { const ps = (A.me && A.me.permissions) || []; return !p || ps.includes('*') || ps.includes(p); }
function activeItem() { const h = location.hash || '#/'; const k = Object.keys(ROUTE_ITEM).filter(r => r === h || (r !== '#/' && h.startsWith(r))).sort((a, b) => b.length - a.length)[0]; return ROUTE_ITEM[k || '#/']; }

async function loadBadges(force) {
  if (!force && Date.now() - SH.at < 20000) return;
  try { const d = await api('/admin/api/dashboard'); SH.dash = d; SH.badges = Object.fromEntries(d.queues.map(q => [q.key, q.count])); SH.at = Date.now(); } catch {}
}

// Status buka/tutup kategori menu, diingat per browser. Kategori halaman aktif selalu terbuka.
const NAV_KEY = 'phsar-admin-nav';
function navState() {
  if (!SH.groups) { try { SH.groups = JSON.parse(localStorage.getItem(NAV_KEY) || '{}') || {}; } catch { SH.groups = {}; } }
  return SH.groups;
}
function toggleGroup(g, btn) {
  const open = btn.getAttribute('aria-expanded') !== 'true';
  navState()[g] = open;
  try { localStorage.setItem(NAV_KEY, JSON.stringify(SH.groups)); } catch { /* mode privat */ }
  btn.setAttribute('aria-expanded', String(open));
  btn.closest('.ngroup').classList.toggle('closed', !open);
}

function sidebar() {
  const act = activeItem(), st = navState();
  const groups = NAV.map(([g, items]) => {
    const vis = items.filter(it => hasPerm(it[2]));
    if (!vis.length) return '';
    const count = it => it[3] ? it[3].split('+').reduce((s, k) => s + (SH.badges[k] || 0), 0) : 0;
    const links = vis.map(it => {
      const [key, href, , , p2] = it, n = count(it);
      if (!href) return `<span class="nitem soon" aria-disabled="true"><span>${t('i_' + key)}</span><span class="soonb">${t(p2 ? 'phase2' : 'soon')}</span></span>`;
      return `<a class="nitem" href="${href}" ${act === key ? 'aria-current="page"' : ''}><span>${t('i_' + key)}</span>${n ? `<span class="badge">${n}</span>` : ''}</a>`;
    }).join('');
    if (vis.length === 1) return `<div class="ngroup"><span class="glabel">${t('g_' + g)}</span>${links}</div>`; // mis. Overview: tanpa tombol
    const open = vis.some(it => it[0] === act) || st[g] === true;
    const sum = vis.reduce((s, it) => s + count(it), 0);
    return `<div class="ngroup${open ? '' : ' closed'}">
      <button type="button" class="glabel gtoggle" aria-expanded="${open}" aria-controls="ng-${g}" onclick="toggleGroup('${g}', this)">
        <span>${t('g_' + g)}</span>${sum ? `<span class="badge gbadge" aria-label="${sum}">${sum}</span>` : ''}<span class="gicon" aria-hidden="true"></span></button>
      <div class="gitems" id="ng-${g}">${links}</div></div>`;
  }).join('');
  const name = (A.me && (A.me.name || A.me.phone)) || 'Phsar Admin';
  const ini = name.replace(/^\+855/, '').split(/\s+/).map(w => w[0]).join('').slice(0, 2).toUpperCase() || 'PA';
  return `<nav class="side" aria-label="${t('nav')}"><a class="brandx" href="#/">${brandSVG}<span><b>Phsar</b><small>ADMIN · ${t('env').toUpperCase()}</small></span></a>
    ${groups}
    <div class="me"><div class="avatar" aria-hidden="true">${esc(ini)}</div><div class="who"><b>${esc(name)}</b><small>${esc(t('r_' + (A.me?.role || 'super_admin')))} · ${t('mfaOn')}</small></div>
      <button type="button" onclick="logout()">${t('logout')}</button></div></nav>`;
}

// shell: halaman login (latar navy) atau tata letak admin dengan menu samping.
function shell(inner) {
  const signedIn = P.token && P.mfa && P.mfa.verified && !MFA.codes;
  if (!signedIn) {
    const step = P.token ? 2 : 1;
    inner = inner.replace(/<section class="card login([^"]*)">/, (m) => `${m}<div class="stepline"><span>${t('step', {n: step})}</span>${langSwitch()}</div>`);
    return `<div class="auth"><span class="brandx">${brandSVG.replace('width="28" height="32"', 'width="32" height="36"')}<span><b>Phsar</b><small>ADMIN</small></span></span>
      ${inner}<ul class="notes"><li>${t('note1')}</li><li>${t('note2')}</li></ul>
      ${P.token ? `<button type="button" class="signout" onclick="logout()">${t('logout')}</button>` : ''}</div>`;
  }
  return `<div class="adm${SH.open ? ' open' : ''}" onclick="if(event.target===this){SH.open=false;this.classList.remove('open')}">
    <div class="topbar"><a class="brandx" href="#/">${brandSVG.replace(/"bm"|#bm\)/g, m => m.replace('bm', 'bmt'))}<span><b>Phsar</b><small>ADMIN</small></span></a>
      <button type="button" onclick="SH.open=true;this.closest('.adm').classList.add('open')" aria-label="${t('menu')}">☰ ${t('menu')}</button></div>
    ${sidebar()}<main class="amain">${inner}</main></div>`;
}

function pageHead(title, desc, actions) {
  return `<header class="phead"><div><h1>${esc(title)}</h1>${desc ? `<p>${desc}</p>` : ''}</div><div class="acts">${bellButton()}${langSwitch()}${actions || ''}</div></header>`;
}
// Halaman lama memanggil adminNav(); menu kini di samping.
function adminNav() { return ''; }

function ago(d) {
  if (!d) return '';
  const m = Math.max(1, Math.round((Date.now() - new Date(d)) / 60000));
  return m < 60 ? t('ago_m', {n: m}) : m < 2880 ? t('ago_h', {n: Math.round(m / 60)}) : t('ago_d', {n: Math.round(m / 1440)});
}
function actionLabel(a, table) {
  const k = 'a_' + a; const s = t(k);
  return s !== k ? (['insert', 'update', 'soft_delete', 'delete', 'restore'].includes(a) ? `${s} · ${table}` : s) : `${a} · ${table}`;
}
const QROUTE = {tickets: ['#/support', 'cs'], reports: ['#/reports', 'moderator'], disputes: ['#/returns', 'cs'], return_refunds: ['#/returns?s=approved', 'finance'], stores: ['#/stores', 'moderator'], products: ['#/products', 'moderator'], refunds: ['#/refunds', 'finance'], payouts: ['#/payouts', 'finance'],
  deliveries: ['#/delivery', 'cs'], failed: ['#/delivery', 'cs'], pickups: ['#/stores', 'moderator'], integrations: ['#/integrations/usage', 'super_admin']};
const amoney = v => '$' + Number(v || 0).toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2});

async function renderDashboard() {
  await loadBadges(true);
  const d = SH.dash || {kpis: {}, queues: [], activity: []}, k = d.kpis;
  const by = Object.fromEntries(d.queues.map(q => [q.key, q]));
  const alerts = [];
  if (by.failed && by.failed.count) alerts.push(`<div class="banner warn"><strong>${t('alertFailed', {n: by.failed.count})}</strong><span style="flex:1">${t('alertFailedSub')}</span><a href="#/delivery" style="color:inherit;font-weight:600">${t('open')}</a></div>`);
  if (by.refunds && by.refunds.count && Date.now() - new Date(by.refunds.oldest) > 864e5) alerts.push(`<div class="banner warn"><strong>${t('alertRefund')}</strong><span style="flex:1">${t('alertRefundSub')}</span><a href="#/refunds" style="color:inherit;font-weight:600">${t('open')}</a></div>`);
  (d.integration_alerts || []).forEach(a => alerts.push(usageAlertBanner(a, true))); // saldo / error AI & SMS
  const waiting = d.queues.filter(q => q.count);
  if (!alerts.length && !waiting.length) alerts.push(`<div class="banner ok"><strong>${t('allClear')}</strong><span>${t('allClearSub')}</span></div>`);
  const queues = [...d.queues].sort((a, b) => (b.count > 0) - (a.count > 0) || new Date(a.oldest || 0) - new Date(b.oldest || 0)).map(q => {
    const [href, role] = QROUTE[q.key] || ['#/', 'super_admin'];
    return `<div class="qrow"><span class="n${q.count ? '' : ' zero'}">${q.count}</span><div class="lbl"><b>${t('q_' + q.key)}</b><small>${q.count ? (q.oldest ? t('oldest', {t: ago(q.oldest)}) : '') : t('nothing')}</small></div>
      <span class="rolechip">${t('r_' + role)}</span><a class="btn" href="${href}">${t('open')}</a></div>`;
  }).join('');
  const acts = d.activity.length ? d.activity.map(a => `<div class="arow"><time>${new Date(a.at).toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'})}</time>
    <div><b>${esc(a.actor || '—')}${a.role ? ' · ' + esc(t('r_' + a.role)) : ''}</b><span>${esc(actionLabel(a.action, a.table))}</span></div></div>`).join('') : `<div class="arow"><span>${t('noActivity')}</span></div>`;
  $('#app').innerHTML = shell(`${pageHead(t('dashTitle'), new Date().toLocaleDateString(P.lang === 'km' ? 'km-KH' : P.lang === 'zh' ? 'zh-CN' : 'en-GB', {weekday: 'short', day: 'numeric', month: 'short', year: 'numeric'}))}
    <section class="alerts" aria-label="Alerts">${alerts.join('')}</section>
    <section class="kpis" aria-label="${t('kpiLabel')}">
      <div class="kpi"><span>${t('k_gmv')}</span><b>${amoney(k.gmv_7d)}</b><span>${t('k_gmvSub', {n: k.orders_7d || 0})}</span></div>
      <div class="kpi"><span>${t('k_comm')}</span><b>${amoney(k.commission_7d)}</b><span>${t('k_commSub')}</span></div>
      <div class="kpi"><span>${t('k_sellers')}</span><b>${k.active_sellers || 0}</b><span>${t('k_sellersSub', {n: k.registered_sellers || 0})}</span></div>
      <div class="kpi"><span>${t('k_transfer')}</span><b>${amoney(k.to_transfer)}</b><span>${t('k_transferSub', {n: k.to_transfer_count || 0})}</span></div></section>
    <div class="cols"><section class="panel"><div class="ph"><h2>${t('todo')}</h2><span class="muted">${t('todoSub')}</span></div>${queues || `<div class="qrow"><span class="muted">${t('nothing')}</span></div>`}</section>
      <section class="panel"><div class="ph"><h2>${t('activity')}</h2>${d.can_audit ? `<a href="#/audit">${t('i_audit')}</a>` : ''}</div>${acts}</section></div>`);
}

// ---- Audit log ----
const AU = {actor: '', table: '', action: '', rows: [], sel: null, tables: [], actions: [], more: false};
async function renderAudit(append) {
  const qs = new URLSearchParams({actor: AU.actor, table: AU.table, action: AU.action});
  if (append && AU.rows.length) qs.set('before', AU.rows[AU.rows.length - 1].id);
  const r = await api('/admin/api/audit?' + qs);
  if (!append) { AU.rows = []; AU.tables = r.tables || AU.tables; AU.actions = r.actions || AU.actions; AU.sel = null; }
  AU.rows = AU.rows.concat(r.records); AU.more = r.records.length === 50;
  drawAudit();
}
function drawAudit() {
  const sel = (id, label, cur, opts) => `<label>${label}<select id="${id}"><option value="">${t('all')}</option>${opts.map(([v, l]) => `<option value="${esc(v)}" ${v === cur ? 'selected' : ''}>${esc(l)}</option>`).join('')}</select></label>`;
  const isView = a => /\.view$/.test(a);
  const rows = AU.rows.map((x, i) => `<tr data-i="${i}" aria-selected="${AU.sel === i}"><td class="muted">${new Date(x.at).toLocaleString()}</td>
    <td><b>${esc(x.actor || t('at_' + x.actor_type))}</b><br><span class="muted small">${esc(x.role ? t('r_' + x.role) : t('at_' + x.actor_type))}</span></td>
    <td>${esc(actionLabel(x.action, x.table))}${isView(x.action) ? `<span class="sens">${t('sensitive')}</span>` : ''}<span class="acode">${esc(x.action)}</span></td>
    <td class="muted">${esc(x.table)}${x.entity ? ' · ' + esc(x.entity.slice(0, 8)) : ''}</td></tr>`).join('');
  const c = AU.sel != null ? AU.rows[AU.sel] : null;
  let detail = `<p class="muted">${t('pickRecord')}</p>`;
  if (c) {
    const keys = c.changed && c.changed.length ? c.changed : Object.keys(c.new || c.old || {}).filter(k => !['created_at', 'updated_at'].includes(k));
    const fmt = v => v == null ? '—' : typeof v === 'object' ? JSON.stringify(v) : String(v);
    const changes = keys.length && (c.old || c.new) ? `<div class="tbl-wrap"><table style="min-width:0"><thead><tr><th>${t('field')}</th><th>${t('before')}</th><th>${t('after')}</th></tr></thead><tbody>
      ${keys.slice(0, 40).map(k => `<tr><td><code>${esc(k)}</code></td><td class="muted">${esc(fmt(c.old && c.old[k]))}</td><td>${esc(fmt(c.new && c.new[k]))}</td></tr>`).join('')}</tbody></table></div>` : '';
    detail = `<div><h2>${esc(actionLabel(c.action, c.table))}</h2><p class="muted">${new Date(c.at).toLocaleString()} · ${esc(c.actor || t('at_' + c.actor_type))}</p></div>
      <dl class="dl"><dt>${t('c_record')}</dt><dd>${esc(c.table)} · ${esc(c.entity || '—')}</dd><dt>Request ID</dt><dd>${esc(c.request_id || '—')}</dd>
        <dt>IP</dt><dd>${esc(c.ip || '—')}</dd><dt>${t('device')}</dt><dd>${esc(c.user_agent || '—')}</dd></dl>
      ${changes}${isView(c.action) || !changes ? `<div class="banner info">${t('readOnly')}</div>` : ''}`;
  }
  $('#app').innerHTML = shell(`${pageHead(t('auditTitle'), t('auditDesc'))}
    <form class="filters" onsubmit="return false">${sel('auActor', t('f_actor'), AU.actor, ['admin', 'seller', 'buyer', 'courier', 'system'].map(a => [a, t('at_' + a)]))}
      ${sel('auTable', t('f_table'), AU.table, AU.tables.map(x => [x, x]))}${sel('auAction', t('f_action'), AU.action, AU.actions.map(x => [x, actionLabel(x, '').replace(/ · $/, '')]))}</form>
    <div class="audit-split"><section class="panel"><div class="tbl-wrap"><table><thead><tr><th>${t('c_time')}</th><th>${t('c_actor')}</th><th>${t('c_action')}</th><th>${t('c_record')}</th></tr></thead>
      <tbody>${rows || `<tr><td colspan="4" class="muted">${t('empty')}</td></tr>`}</tbody></table></div>
      ${AU.more ? `<p class="foot"><button type="button" id="auMore">${t('loadMore')}</button></p>` : ''}</section>
      <section class="panel"><div class="ph"><h2>${t('recordDetail')}</h2></div><div class="card" style="border:0;border-radius:0">${detail}<p class="muted small">${t('secretNote')}</p></div></section></div>`);
  ['auActor', 'auTable', 'auAction'].forEach(id => $('#' + id).onchange = () => { AU.actor = $('#auActor').value; AU.table = $('#auTable').value; AU.action = $('#auAction').value; renderAudit().catch(e => toast(errMsg(e), 'error')); });
  document.querySelectorAll('tr[data-i]').forEach(tr => tr.onclick = () => { AU.sel = +tr.dataset.i; drawAudit(); });
  if ($('#auMore')) $('#auMore').onclick = () => renderAudit(true).catch(e => toast(errMsg(e), 'error'));
}
