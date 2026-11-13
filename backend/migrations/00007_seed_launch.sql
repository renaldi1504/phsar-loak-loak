-- +goose Up
-- Data awal untuk launch: Sihanoukville, ongkir gratis, komisi 8%, peran admin, kategori dari website lama.

INSERT INTO shipping_zones (code, fee_per_store, is_estimate) VALUES ('sihanoukville', 0.00, false);

-- Batas wilayah: KOTAK KASAR di sekitar kota Sihanoukville (lng 103.44–103.62, lat 10.55–10.72).
-- Ganti dengan poligon yang digambar di peta sebelum launch (admin: Pengiriman & wilayah).
INSERT INTO service_areas (code, name, boundary, shipping_zone_code, is_active) VALUES (
  'sihanoukville',
  '{"en":"Sihanoukville","km":"ក្រុងព្រះសីហនុ","zh":"西哈努克城"}',
  '{"type":"Polygon","coordinates":[[[103.44,10.55],[103.62,10.55],[103.62,10.72],[103.44,10.72],[103.44,10.55]]]}',
  'sihanoukville', true);

INSERT INTO commission_rules (scope, rate_percent) VALUES ('global', 8.00);

-- Kurs tampilan riel (website lama memakai $1 = 4,000). Job fx_fetch akan menambah kurs harian.
INSERT INTO exchange_rates (rate_date, riel_per_usd, source) VALUES (current_date, 4000, 'manual');

INSERT INTO admin_roles (code, name_i18n) VALUES
  ('super_admin',  '{"en":"Super admin","km":"អ្នកគ្រប់គ្រងកំពូល","zh":"超级管理员"}'),
  ('moderator',    '{"en":"Moderator","km":"អ្នកត្រួតពិនិត្យ","zh":"审核员"}'),
  ('cs',           '{"en":"Customer service","km":"សេវាអតិថិជន","zh":"客服"}'),
  ('finance',      '{"en":"Finance","km":"ហិរញ្ញវត្ថុ","zh":"财务"}'),
  ('ads_reviewer', '{"en":"Ads reviewer","km":"អ្នកពិនិត្យការផ្សាយ","zh":"广告审核员"}');

INSERT INTO role_permissions (role_code, permission) VALUES
  ('super_admin', '*'),
  ('moderator', 'seller.verify'), ('moderator', 'pickup_point.verify'), ('moderator', 'product.moderate'),
  ('moderator', 'translation.review'), ('moderator', 'report.handle'), ('moderator', 'catalog.manage'),
  ('cs', 'order.view'), ('cs', 'return.decide'), ('cs', 'support.handle'), ('cs', 'customer.view'), ('cs', 'delivery.dispatch'),
  ('finance', 'order.view'), ('finance', 'payout.approve'), ('finance', 'reconciliation.run'), ('finance', 'refund.decide'),
  ('finance', 'commission.manage'), ('finance', 'audit.view'),
  ('ads_reviewer', 'ads.review');

-- Kategori utama, sama dengan website lama (catalog.js / app.js).
INSERT INTO categories (id, icon, allows_used, sort_order) VALUES
  ('electronics', 'headphones', true,  10),
  ('beauty',      'beauty',     false, 20),
  ('fashion',     'shirt',      true,  30),
  ('accessories', 'star',       true,  40),
  ('kitchen',     'cooker',     true,  50);

INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES
  ('electronics', 'en', 'Electronics',    'electronics',    'electronics'),
  ('electronics', 'km', 'អេឡិចត្រូនិក',  'electronics',    'electronics'),
  ('electronics', 'zh', '电子产品',        'electronics',    'electronics'),
  ('beauty',      'en', 'Beauty',         'beauty',         'beauty'),
  ('beauty',      'km', 'សម្រស់',          'beauty',         'beauty'),
  ('beauty',      'zh', '美妆',            'beauty',         'beauty'),
  ('fashion',     'en', 'Fashion',        'fashion',        'fashion'),
  ('fashion',     'km', 'សម្លៀកបំពាក់',    'fashion',        'fashion'),
  ('fashion',     'zh', '时尚',            'fashion',        'fashion'),
  ('accessories', 'en', 'Accessories',    'accessories',    'accessories'),
  ('accessories', 'km', 'គ្រឿងបន្ថែម',      'accessories',    'accessories'),
  ('accessories', 'zh', '配饰',            'accessories',    'accessories'),
  ('kitchen',     'en', 'Home & Kitchen', 'home-kitchen',   'home-kitchen'),
  ('kitchen',     'km', 'ផ្ទះបាយ',          'home-kitchen',   'home-kitchen'),
  ('kitchen',     'zh', '家居厨房',         'home-kitchen',   'home-kitchen');

-- +goose Down
DELETE FROM category_translations;
DELETE FROM categories;
DELETE FROM role_permissions;
DELETE FROM admin_roles;
DELETE FROM exchange_rates;
DELETE FROM commission_rules;
DELETE FROM service_areas;
DELETE FROM shipping_zones;
