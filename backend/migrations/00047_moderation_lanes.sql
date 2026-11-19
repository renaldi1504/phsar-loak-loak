-- +goose Up
-- Moderasi 3 jalur: produk langsung tayang (otomatis) kecuali ada alasan untuk ditinjau manual.
--   auto   : lolos semua cek → approved, tetap masuk daftar "Baru tayang" untuk dicek admin belakangan
--   manual : toko masa percobaan, kategori terbatas, kata "review", harga janggal, strike, laporan pembeli
--   block  : kata dengan aksi "block" → ditolak saat dikirim
-- Angka aturan (masa percobaan, ambang diskon, dll.) di platform_settings['moderation_rules'].

ALTER TABLE categories ADD COLUMN restricted boolean NOT NULL DEFAULT false;  -- produknya selalu ditinjau manual
UPDATE categories SET restricted = true WHERE id = 'beauty';

ALTER TABLE products ADD COLUMN review_lane varchar(8) CHECK (review_lane IN ('auto', 'manual'));
ALTER TABLE products ADD COLUMN post_review_pending boolean NOT NULL DEFAULT false;  -- tayang otomatis, belum dilihat admin
ALTER TABLE products ADD COLUMN approved_min_price numeric(12,2);                    -- harga saat terakhir disetujui (cek harga anjlok)
UPDATE products SET approved_min_price = min_price WHERE moderation_status = 'approved';
CREATE INDEX products_post_review ON products (published_at) WHERE post_review_pending AND deleted_at IS NULL;

-- Keputusan otomatis dicatat juga (reviewer_id kosong = sistem).
ALTER TABLE product_moderation_reviews ALTER COLUMN reviewer_id DROP NOT NULL;

-- Harga minimum wajar per kata kunci di nama produk (mis. "iphone 15" < $250 → ditinjau).
CREATE TABLE price_floor_rules (
  id         integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  keyword    varchar(80) NOT NULL,
  min_price  numeric(12,2) NOT NULL CHECK (min_price > 0),
  note       varchar(200),
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX price_floor_rules_keyword ON price_floor_rules (lower(keyword));
CREATE TRIGGER price_floor_rules_audit AFTER INSERT OR UPDATE OR DELETE ON price_floor_rules FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Angka awal = perkiraan kasar harga bekas termurah yang masih masuk akal. Silakan disesuaikan admin.
INSERT INTO price_floor_rules (keyword, min_price, note) VALUES
  ('iphone 17', 450, NULL), ('iphone 16', 350, NULL), ('iphone 15', 250, NULL), ('iphone 14', 180, NULL), ('iphone 13', 140, NULL),
  ('galaxy s25', 350, NULL), ('galaxy s24', 250, NULL), ('galaxy z fold', 300, NULL), ('galaxy z flip', 200, NULL),
  ('macbook', 200, NULL), ('ipad pro', 200, NULL), ('ipad', 80, NULL), ('apple watch', 60, NULL),
  ('airpods pro', 40, NULL), ('airpods', 25, NULL),
  ('playstation 5', 200, NULL), ('ps5', 200, NULL), ('nintendo switch', 100, NULL), ('dyson', 100, NULL), ('gopro', 80, NULL),
  ('rolex', 1000, 'luxury watch'), ('patek philippe', 2000, 'luxury watch'), ('audemars piguet', 2000, 'luxury watch'),
  ('louis vuitton', 150, 'luxury brand'), ('gucci', 100, 'luxury brand'), ('chanel', 150, 'luxury brand'), ('hermes', 300, 'luxury brand'),
  ('dior', 100, 'luxury brand'), ('prada', 100, 'luxury brand'), ('balenciaga', 100, 'luxury brand')
ON CONFLICT DO NOTHING;

-- Barang yang jelas dilarang → blokir langsung.
UPDATE prohibited_keywords SET action = 'block'
  WHERE lang = 'en' AND keyword IN ('ammunition', 'grenade', 'methamphetamine', 'cocaine', 'heroin');

INSERT INTO prohibited_keywords (keyword, lang, category, action) VALUES
  -- dilarang
  ('crystal meth', 'en', 'drugs', 'block'), ('ecstasy pill', 'en', 'drugs', 'block'), ('ketamine', 'en', 'drugs', 'block'),
  ('e-cigarette', 'en', 'vape', 'block'), ('vape', 'en', 'vape', 'block'), ('e-liquid', 'en', 'vape', 'block'), ('pod system', 'en', 'vape', 'block'),
  ('bullet', 'en', 'weapons', 'review'), ('explosive', 'en', 'weapons', 'block'), ('fireworks', 'en', 'explosives', 'review'),
  ('firecracker', 'en', 'explosives', 'block'),
  ('cigarette', 'en', 'tobacco_alcohol', 'review'), ('tobacco', 'en', 'tobacco_alcohol', 'review'), ('whisky', 'en', 'tobacco_alcohol', 'review'),
  ('vodka', 'en', 'tobacco_alcohol', 'review'),
  ('fake id', 'en', 'documents', 'block'), ('counterfeit money', 'en', 'counterfeit', 'block'),
  ('spy camera', 'en', 'surveillance', 'review'), ('hidden camera', 'en', 'surveillance', 'review'),
  ('live animal', 'en', 'animals', 'review'),
  ('antique', 'en', 'antiquities', 'review'), ('artifact', 'en', 'antiquities', 'review'),
  ('prescription', 'en', 'medicine', 'review'), ('antibiotic', 'en', 'medicine', 'review'), ('sildenafil', 'en', 'medicine', 'review'),
  -- terbatas (ditinjau manual)
  ('supplement', 'en', 'health', 'review'), ('vitamin', 'en', 'health', 'review'), ('collagen', 'en', 'health', 'review'),
  ('slimming', 'en', 'health', 'review'), ('diet pill', 'en', 'health', 'review'), ('herbal medicine', 'en', 'health', 'review'),
  ('whitening', 'en', 'cosmetics', 'review'), ('bleaching cream', 'en', 'cosmetics', 'review'), ('mercury', 'en', 'cosmetics', 'review'),
  ('infant formula', 'en', 'baby', 'review'), ('baby formula', 'en', 'baby', 'review'), ('car seat', 'en', 'baby', 'review'),
  ('gift card', 'en', 'digital', 'review'), ('voucher', 'en', 'digital', 'review'), ('top up', 'en', 'digital', 'review'),
  ('game account', 'en', 'digital', 'review'), ('netflix account', 'en', 'digital', 'review'),
  ('24k', 'en', 'jewelry', 'review'), ('gold bar', 'en', 'jewelry', 'review'), ('diamond', 'en', 'jewelry', 'review'),
  ('machete', 'en', 'blades', 'review'), ('pocket knife', 'en', 'blades', 'review'), ('pepper spray', 'en', 'weapons', 'review'),
  ('pesticide', 'en', 'chemicals', 'review'), ('insecticide', 'en', 'chemicals', 'review'), ('car battery', 'en', 'chemicals', 'review'),
  ('buddha statue', 'en', 'religious', 'review'), ('amulet', 'en', 'religious', 'review'),
  ('first copy', 'en', 'counterfeit', 'review'), ('mirror quality', 'en', 'counterfeit', 'review'), ('1:1', 'en', 'counterfeit', 'review'),
  ('original quality', 'en', 'counterfeit', 'review'), ('kw', 'en', 'counterfeit', 'review'),
  -- 中文
  ('电子烟', 'zh', 'vape', 'block'), ('烟油', 'zh', 'vape', 'block'), ('冰毒', 'zh', 'drugs', 'block'), ('毒品', 'zh', 'drugs', 'block'),
  ('子弹', 'zh', 'weapons', 'block'), ('炸药', 'zh', 'weapons', 'block'), ('枪', 'zh', 'weapons', 'review'), ('象牙', 'zh', 'protected_wildlife', 'review'),
  ('香烟', 'zh', 'tobacco_alcohol', 'review'), ('白酒', 'zh', 'tobacco_alcohol', 'review'), ('啤酒', 'zh', 'tobacco_alcohol', 'review'), ('红酒', 'zh', 'tobacco_alcohol', 'review'),
  ('保健品', 'zh', 'health', 'review'), ('减肥', 'zh', 'health', 'review'), ('美白', 'zh', 'cosmetics', 'review'), ('奶粉', 'zh', 'baby', 'review'),
  ('礼品卡', 'zh', 'digital', 'review'), ('充值', 'zh', 'digital', 'review'), ('黄金', 'zh', 'jewelry', 'review'), ('钻石', 'zh', 'jewelry', 'review'),
  ('农药', 'zh', 'chemicals', 'review'), ('古董', 'zh', 'antiquities', 'review'), ('佛像', 'zh', 'religious', 'review'),
  ('高仿', 'zh', 'counterfeit', 'review'), ('复刻', 'zh', 'counterfeit', 'review'), ('精仿', 'zh', 'counterfeit', 'review'),
  -- ខ្មែរ (ពាក្យវែង ដើម្បីកុំឲ្យប៉ះពាក្យផ្សេង)
  ('បារីអេឡិចត្រូនិក', 'km', 'vape', 'block'), ('គ្រឿងញៀន', 'km', 'drugs', 'block'), ('គ្រាប់កាំភ្លើង', 'km', 'weapons', 'block'),
  ('កាំភ្លើង', 'km', 'weapons', 'review'), ('វត្ថុបុរាណ', 'km', 'antiquities', 'review')
ON CONFLICT (lang, keyword) DO NOTHING;

-- +goose Down
DELETE FROM prohibited_keywords WHERE keyword IN (
  'crystal meth','ecstasy pill','ketamine','e-cigarette','vape','e-liquid','pod system','bullet','explosive','fireworks','firecracker',
  'cigarette','tobacco','whisky','vodka','fake id','counterfeit money','spy camera','hidden camera','live animal',
  'antique','artifact','prescription','antibiotic','sildenafil','supplement','vitamin','collagen','slimming',
  'diet pill','herbal medicine','whitening','bleaching cream','mercury','infant formula','baby formula','car seat','gift card','voucher',
  'top up','game account','netflix account','24k','gold bar','diamond','machete','pocket knife','pepper spray','pesticide','insecticide',
  'car battery','buddha statue','amulet','first copy','mirror quality','1:1','original quality','kw',
  '电子烟','烟油','冰毒','毒品','子弹','炸药','枪','象牙','香烟','白酒','啤酒','红酒','保健品','减肥','美白','奶粉','礼品卡','充值','黄金','钻石',
  '农药','古董','佛像','高仿','复刻','精仿','បារីអេឡិចត្រូនិក','គ្រឿងញៀន','គ្រាប់កាំភ្លើង','កាំភ្លើង','វត្ថុបុរាណ');
UPDATE prohibited_keywords SET action = 'review'
  WHERE lang = 'en' AND keyword IN ('ammunition', 'grenade', 'methamphetamine', 'cocaine', 'heroin');
DROP TABLE price_floor_rules;
DELETE FROM product_moderation_reviews WHERE reviewer_id IS NULL;
ALTER TABLE product_moderation_reviews ALTER COLUMN reviewer_id SET NOT NULL;
DROP INDEX products_post_review;
ALTER TABLE products DROP COLUMN approved_min_price;
ALTER TABLE products DROP COLUMN post_review_pending;
ALTER TABLE products DROP COLUMN review_lane;
ALTER TABLE categories DROP COLUMN restricted;
