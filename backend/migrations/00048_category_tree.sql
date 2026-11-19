-- +goose Up
-- Struktur kategori baru: 12 kategori utama, masing-masing 3–6 sub-kategori (2 tingkat).
-- Seller wajib memilih sub-kategori. ID lama (fashion, accessories, electronics, beauty, kitchen) tetap dipakai
-- sebagai kategori utama supaya relasi (komisi, koleksi, iklan) tidak putus; nama & slug-nya diperbarui.
-- Bekas/terbatas di sub-kategori mengikuti induk (dicek di kode: allows_used = anak AND induk, restricted = anak OR induk).

UPDATE categories SET icon = 'shirt', allows_used = true, restricted = false, sort_order = 10, deleted_at = NULL, parent_id = NULL WHERE id = 'fashion';
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('fashion', 'en', 'Women''s Fashion', 'womens-fashion', 'womens-fashion')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('fashion', 'km', 'សម្លៀកបំពាក់នារី', 'womens-fashion', 'womens-fashion')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('fashion', 'zh', '女装', 'womens-fashion', 'womens-fashion')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('womens-tops', 'fashion', true, false, 10);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-tops', 'en', 'Tops', 'womens-tops', 'womens-fashion/womens-tops');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-tops', 'km', 'អាវ', 'womens-tops', 'womens-fashion/womens-tops');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-tops', 'zh', '上衣', 'womens-tops', 'womens-fashion/womens-tops');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('womens-dresses', 'fashion', true, false, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-dresses', 'en', 'Dresses & Skirts', 'womens-dresses', 'womens-fashion/womens-dresses');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-dresses', 'km', 'រ៉ូប និងសំពត់', 'womens-dresses', 'womens-fashion/womens-dresses');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-dresses', 'zh', '连衣裙和半身裙', 'womens-dresses', 'womens-fashion/womens-dresses');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('womens-bottoms', 'fashion', true, false, 30);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-bottoms', 'en', 'Pants & Shorts', 'womens-bottoms', 'womens-fashion/womens-bottoms');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-bottoms', 'km', 'ខោ', 'womens-bottoms', 'womens-fashion/womens-bottoms');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-bottoms', 'zh', '裤子', 'womens-bottoms', 'womens-fashion/womens-bottoms');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('womens-underwear', 'fashion', false, false, 40);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-underwear', 'en', 'Underwear & Sleepwear', 'womens-underwear', 'womens-fashion/womens-underwear');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-underwear', 'km', 'ខោអាវទ្រនាប់ និងខោអាវគេង', 'womens-underwear', 'womens-fashion/womens-underwear');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-underwear', 'zh', '内衣和睡衣', 'womens-underwear', 'womens-fashion/womens-underwear');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('womens-shoes', 'fashion', true, false, 50);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-shoes', 'en', 'Women''s Shoes', 'womens-shoes', 'womens-fashion/womens-shoes');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-shoes', 'km', 'ស្បែកជើងនារី', 'womens-shoes', 'womens-fashion/womens-shoes');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-shoes', 'zh', '女鞋', 'womens-shoes', 'womens-fashion/womens-shoes');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('womens-bags', 'fashion', true, false, 60);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-bags', 'en', 'Women''s Bags', 'womens-bags', 'womens-fashion/womens-bags');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-bags', 'km', 'កាបូបនារី', 'womens-bags', 'womens-fashion/womens-bags');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('womens-bags', 'zh', '女包', 'womens-bags', 'womens-fashion/womens-bags');
INSERT INTO categories (id, icon, allows_used, restricted, sort_order) VALUES ('mens-fashion', 'tshirt', true, false, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-fashion', 'en', 'Men''s Fashion', 'mens-fashion', 'mens-fashion')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-fashion', 'km', 'សម្លៀកបំពាក់បុរស', 'mens-fashion', 'mens-fashion')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-fashion', 'zh', '男装', 'mens-fashion', 'mens-fashion')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('mens-tops', 'mens-fashion', true, false, 10);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-tops', 'en', 'T-shirts & Shirts', 'mens-tops', 'mens-fashion/mens-tops');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-tops', 'km', 'អាវយឺត និងអាវសឺមី', 'mens-tops', 'mens-fashion/mens-tops');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-tops', 'zh', 'T恤和衬衫', 'mens-tops', 'mens-fashion/mens-tops');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('mens-bottoms', 'mens-fashion', true, false, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-bottoms', 'en', 'Men''s Pants & Shorts', 'mens-bottoms', 'mens-fashion/mens-bottoms');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-bottoms', 'km', 'ខោបុរស', 'mens-bottoms', 'mens-fashion/mens-bottoms');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-bottoms', 'zh', '男裤', 'mens-bottoms', 'mens-fashion/mens-bottoms');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('mens-underwear', 'mens-fashion', false, false, 30);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-underwear', 'en', 'Men''s Underwear', 'mens-underwear', 'mens-fashion/mens-underwear');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-underwear', 'km', 'ខោអាវទ្រនាប់បុរស', 'mens-underwear', 'mens-fashion/mens-underwear');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-underwear', 'zh', '男士内衣', 'mens-underwear', 'mens-fashion/mens-underwear');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('mens-shoes', 'mens-fashion', true, false, 40);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-shoes', 'en', 'Men''s Shoes', 'mens-shoes', 'mens-fashion/mens-shoes');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-shoes', 'km', 'ស្បែកជើងបុរស', 'mens-shoes', 'mens-fashion/mens-shoes');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-shoes', 'zh', '男鞋', 'mens-shoes', 'mens-fashion/mens-shoes');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('mens-bags', 'mens-fashion', true, false, 50);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-bags', 'en', 'Bags & Wallets', 'mens-bags', 'mens-fashion/mens-bags');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-bags', 'km', 'កាបូប និងកាបូបលុយ', 'mens-bags', 'mens-fashion/mens-bags');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mens-bags', 'zh', '男包和钱包', 'mens-bags', 'mens-fashion/mens-bags');
UPDATE categories SET icon = 'watch', allows_used = true, restricted = false, sort_order = 30, deleted_at = NULL, parent_id = NULL WHERE id = 'accessories';
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('accessories', 'en', 'Watches & Accessories', 'watches-accessories', 'watches-accessories')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('accessories', 'km', 'នាឡិកា និងគ្រឿងលម្អ', 'watches-accessories', 'watches-accessories')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('accessories', 'zh', '手表和配饰', 'watches-accessories', 'watches-accessories')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('watches', 'accessories', true, false, 10);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('watches', 'en', 'Watches', 'watches', 'watches-accessories/watches');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('watches', 'km', 'នាឡិកាដៃ', 'watches', 'watches-accessories/watches');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('watches', 'zh', '手表', 'watches', 'watches-accessories/watches');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('jewelry', 'accessories', true, true, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('jewelry', 'en', 'Jewelry', 'jewelry', 'watches-accessories/jewelry');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('jewelry', 'km', 'គ្រឿងអលង្ការ', 'jewelry', 'watches-accessories/jewelry');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('jewelry', 'zh', '珠宝首饰', 'jewelry', 'watches-accessories/jewelry');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('eyewear', 'accessories', true, false, 30);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('eyewear', 'en', 'Glasses & Sunglasses', 'eyewear', 'watches-accessories/eyewear');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('eyewear', 'km', 'វ៉ែនតា', 'eyewear', 'watches-accessories/eyewear');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('eyewear', 'zh', '眼镜和太阳镜', 'eyewear', 'watches-accessories/eyewear');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('hats', 'accessories', true, false, 40);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('hats', 'en', 'Hats & Caps', 'hats', 'watches-accessories/hats');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('hats', 'km', 'មួក', 'hats', 'watches-accessories/hats');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('hats', 'zh', '帽子', 'hats', 'watches-accessories/hats');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('belts', 'accessories', true, false, 50);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('belts', 'en', 'Belts', 'belts', 'watches-accessories/belts');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('belts', 'km', 'ខ្សែក្រវ៉ាត់', 'belts', 'watches-accessories/belts');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('belts', 'zh', '腰带', 'belts', 'watches-accessories/belts');
INSERT INTO categories (id, icon, allows_used, restricted, sort_order) VALUES ('phones', 'phone', true, false, 40);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('phones', 'en', 'Phones & Tablets', 'phones-tablets', 'phones-tablets')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('phones', 'km', 'ទូរសព្ទ និងថេប្លេត', 'phones-tablets', 'phones-tablets')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('phones', 'zh', '手机和平板', 'phones-tablets', 'phones-tablets')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('mobile-phones', 'phones', true, false, 10);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mobile-phones', 'en', 'Mobile Phones', 'mobile-phones', 'phones-tablets/mobile-phones');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mobile-phones', 'km', 'ទូរសព្ទដៃ', 'mobile-phones', 'phones-tablets/mobile-phones');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mobile-phones', 'zh', '手机', 'mobile-phones', 'phones-tablets/mobile-phones');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('tablets', 'phones', true, false, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('tablets', 'en', 'Tablets', 'tablets', 'phones-tablets/tablets');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('tablets', 'km', 'ថេប្លេត', 'tablets', 'phones-tablets/tablets');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('tablets', 'zh', '平板电脑', 'tablets', 'phones-tablets/tablets');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('phone-cases', 'phones', true, false, 30);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('phone-cases', 'en', 'Cases & Screen Protectors', 'phone-cases', 'phones-tablets/phone-cases');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('phone-cases', 'km', 'ស្រោម និងកញ្ចក់ការពារ', 'phone-cases', 'phones-tablets/phone-cases');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('phone-cases', 'zh', '手机壳和贴膜', 'phone-cases', 'phones-tablets/phone-cases');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('chargers-cables', 'phones', true, false, 40);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('chargers-cables', 'en', 'Chargers & Cables', 'chargers-cables', 'phones-tablets/chargers-cables');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('chargers-cables', 'km', 'ឆ្នាំងសាក និងខ្សែសាក', 'chargers-cables', 'phones-tablets/chargers-cables');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('chargers-cables', 'zh', '充电器和数据线', 'chargers-cables', 'phones-tablets/chargers-cables');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('power-banks', 'phones', true, false, 50);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('power-banks', 'en', 'Power Banks', 'power-banks', 'phones-tablets/power-banks');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('power-banks', 'km', 'ថ្មសាកបម្រុង', 'power-banks', 'phones-tablets/power-banks');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('power-banks', 'zh', '充电宝', 'power-banks', 'phones-tablets/power-banks');
UPDATE categories SET icon = 'laptop', allows_used = true, restricted = false, sort_order = 50, deleted_at = NULL, parent_id = NULL WHERE id = 'electronics';
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('electronics', 'en', 'Electronics & Computers', 'electronics', 'electronics')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('electronics', 'km', 'អេឡិចត្រូនិក និងកុំព្យូទ័រ', 'electronics', 'electronics')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('electronics', 'zh', '电子和电脑', 'electronics', 'electronics')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('laptops', 'electronics', true, false, 10);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('laptops', 'en', 'Laptops', 'laptops', 'electronics/laptops');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('laptops', 'km', 'កុំព្យូទ័រយួរដៃ', 'laptops', 'electronics/laptops');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('laptops', 'zh', '笔记本电脑', 'laptops', 'electronics/laptops');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('computer-accessories', 'electronics', true, false, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('computer-accessories', 'en', 'Computer Accessories', 'computer-accessories', 'electronics/computer-accessories');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('computer-accessories', 'km', 'គ្រឿងកុំព្យូទ័រ', 'computer-accessories', 'electronics/computer-accessories');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('computer-accessories', 'zh', '电脑配件', 'computer-accessories', 'electronics/computer-accessories');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('audio', 'electronics', true, false, 30);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('audio', 'en', 'Audio', 'audio', 'electronics/audio');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('audio', 'km', 'ឧបករណ៍សំឡេង', 'audio', 'electronics/audio');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('audio', 'zh', '音频设备', 'audio', 'electronics/audio');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('cameras', 'electronics', true, false, 40);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('cameras', 'en', 'Cameras', 'cameras', 'electronics/cameras');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('cameras', 'km', 'កាមេរ៉ា', 'cameras', 'electronics/cameras');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('cameras', 'zh', '相机', 'cameras', 'electronics/cameras');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('gaming', 'electronics', true, false, 50);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('gaming', 'en', 'Gaming', 'gaming', 'electronics/gaming');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('gaming', 'km', 'ហ្គេម', 'gaming', 'electronics/gaming');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('gaming', 'zh', '游戏', 'gaming', 'electronics/gaming');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('tvs', 'electronics', true, false, 60);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('tvs', 'en', 'TV & Video', 'tvs', 'electronics/tvs');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('tvs', 'km', 'ទូរទស្សន៍', 'tvs', 'electronics/tvs');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('tvs', 'zh', '电视和影音', 'tvs', 'electronics/tvs');
UPDATE categories SET icon = 'beauty', allows_used = false, restricted = true, sort_order = 60, deleted_at = NULL, parent_id = NULL WHERE id = 'beauty';
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('beauty', 'en', 'Beauty', 'beauty', 'beauty')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('beauty', 'km', 'សម្រស់', 'beauty', 'beauty')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('beauty', 'zh', '美妆', 'beauty', 'beauty')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('skincare', 'beauty', false, true, 10);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('skincare', 'en', 'Skincare', 'skincare', 'beauty/skincare');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('skincare', 'km', 'ថែរក្សាស្បែក', 'skincare', 'beauty/skincare');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('skincare', 'zh', '护肤', 'skincare', 'beauty/skincare');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('makeup', 'beauty', false, true, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('makeup', 'en', 'Makeup', 'makeup', 'beauty/makeup');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('makeup', 'km', 'គ្រឿងសម្អាង', 'makeup', 'beauty/makeup');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('makeup', 'zh', '彩妆', 'makeup', 'beauty/makeup');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('haircare', 'beauty', false, true, 30);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('haircare', 'en', 'Hair Care', 'haircare', 'beauty/haircare');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('haircare', 'km', 'ថែរក្សាសក់', 'haircare', 'beauty/haircare');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('haircare', 'zh', '护发', 'haircare', 'beauty/haircare');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('fragrance', 'beauty', false, true, 40);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('fragrance', 'en', 'Fragrance', 'fragrance', 'beauty/fragrance');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('fragrance', 'km', 'ទឹកអប់', 'fragrance', 'beauty/fragrance');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('fragrance', 'zh', '香水', 'fragrance', 'beauty/fragrance');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('body-care', 'beauty', false, true, 50);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('body-care', 'en', 'Body Care', 'body-care', 'beauty/body-care');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('body-care', 'km', 'ថែរក្សារាងកាយ', 'body-care', 'beauty/body-care');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('body-care', 'zh', '身体护理', 'body-care', 'beauty/body-care');
INSERT INTO categories (id, icon, allows_used, restricted, sort_order) VALUES ('health', 'health', false, true, 70);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('health', 'en', 'Health', 'health', 'health')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('health', 'km', 'សុខភាព', 'health', 'health')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('health', 'zh', '健康', 'health', 'health')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('supplements', 'health', false, true, 10);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('supplements', 'en', 'Supplements', 'supplements', 'health/supplements');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('supplements', 'km', 'អាហារបំប៉ន', 'supplements', 'health/supplements');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('supplements', 'zh', '保健品', 'supplements', 'health/supplements');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('health-devices', 'health', false, true, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('health-devices', 'en', 'Health Devices', 'health-devices', 'health/health-devices');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('health-devices', 'km', 'ឧបករណ៍សុខភាព', 'health-devices', 'health/health-devices');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('health-devices', 'zh', '健康器械', 'health-devices', 'health/health-devices');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('first-aid', 'health', false, true, 30);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('first-aid', 'en', 'Masks & First Aid', 'first-aid', 'health/first-aid');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('first-aid', 'km', 'ម៉ាស់ និងជំនួយបឋម', 'first-aid', 'health/first-aid');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('first-aid', 'zh', '口罩和急救用品', 'first-aid', 'health/first-aid');
INSERT INTO categories (id, icon, allows_used, restricted, sort_order) VALUES ('mother-baby', 'baby', true, false, 80);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mother-baby', 'en', 'Mother & Baby', 'mother-baby', 'mother-baby')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mother-baby', 'km', 'ម្ដាយ និងទារក', 'mother-baby', 'mother-baby')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('mother-baby', 'zh', '母婴', 'mother-baby', 'mother-baby')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('kids-clothing', 'mother-baby', true, false, 10);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('kids-clothing', 'en', 'Kids'' Clothing', 'kids-clothing', 'mother-baby/kids-clothing');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('kids-clothing', 'km', 'សម្លៀកបំពាក់កុមារ', 'kids-clothing', 'mother-baby/kids-clothing');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('kids-clothing', 'zh', '童装', 'kids-clothing', 'mother-baby/kids-clothing');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('diapers', 'mother-baby', false, false, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('diapers', 'en', 'Diapers', 'diapers', 'mother-baby/diapers');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('diapers', 'km', 'កន្ទបទារក', 'diapers', 'mother-baby/diapers');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('diapers', 'zh', '纸尿裤', 'diapers', 'mother-baby/diapers');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('baby-food', 'mother-baby', false, true, 30);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('baby-food', 'en', 'Baby Milk & Food', 'baby-food', 'mother-baby/baby-food');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('baby-food', 'km', 'ទឹកដោះគោ និងអាហារទារក', 'baby-food', 'mother-baby/baby-food');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('baby-food', 'zh', '奶粉和辅食', 'baby-food', 'mother-baby/baby-food');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('baby-gear', 'mother-baby', true, false, 40);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('baby-gear', 'en', 'Baby Gear', 'baby-gear', 'mother-baby/baby-gear');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('baby-gear', 'km', 'សម្ភារៈទារក', 'baby-gear', 'mother-baby/baby-gear');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('baby-gear', 'zh', '婴儿用品', 'baby-gear', 'mother-baby/baby-gear');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('kids-toys', 'mother-baby', true, false, 50);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('kids-toys', 'en', 'Kids'' Toys', 'kids-toys', 'mother-baby/kids-toys');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('kids-toys', 'km', 'ប្រដាប់ក្មេងលេង', 'kids-toys', 'mother-baby/kids-toys');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('kids-toys', 'zh', '儿童玩具', 'kids-toys', 'mother-baby/kids-toys');
UPDATE categories SET icon = 'cooker', allows_used = true, restricted = false, sort_order = 90, deleted_at = NULL, parent_id = NULL WHERE id = 'kitchen';
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('kitchen', 'en', 'Home & Kitchen', 'home-kitchen', 'home-kitchen')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('kitchen', 'km', 'ផ្ទះ និងផ្ទះបាយ', 'home-kitchen', 'home-kitchen')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('kitchen', 'zh', '家居厨房', 'home-kitchen', 'home-kitchen')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('kitchenware', 'kitchen', true, false, 10);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('kitchenware', 'en', 'Kitchenware', 'kitchenware', 'home-kitchen/kitchenware');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('kitchenware', 'km', 'ប្រដាប់ផ្ទះបាយ', 'kitchenware', 'home-kitchen/kitchenware');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('kitchenware', 'zh', '厨具', 'kitchenware', 'home-kitchen/kitchenware');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('home-appliances', 'kitchen', true, false, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('home-appliances', 'en', 'Home Appliances', 'home-appliances', 'home-kitchen/home-appliances');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('home-appliances', 'km', 'គ្រឿងអគ្គិសនីក្នុងផ្ទះ', 'home-appliances', 'home-kitchen/home-appliances');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('home-appliances', 'zh', '家用电器', 'home-appliances', 'home-kitchen/home-appliances');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('furniture', 'kitchen', true, false, 30);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('furniture', 'en', 'Furniture', 'furniture', 'home-kitchen/furniture');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('furniture', 'km', 'គ្រឿងសង្ហារិម', 'furniture', 'home-kitchen/furniture');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('furniture', 'zh', '家具', 'furniture', 'home-kitchen/furniture');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('home-decor', 'kitchen', true, false, 40);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('home-decor', 'en', 'Home Decor', 'home-decor', 'home-kitchen/home-decor');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('home-decor', 'km', 'គ្រឿងតុបតែងផ្ទះ', 'home-decor', 'home-kitchen/home-decor');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('home-decor', 'zh', '家居装饰', 'home-decor', 'home-kitchen/home-decor');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('cleaning', 'kitchen', true, false, 50);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('cleaning', 'en', 'Cleaning & Laundry', 'cleaning', 'home-kitchen/cleaning');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('cleaning', 'km', 'សម្ភារៈសម្អាត', 'cleaning', 'home-kitchen/cleaning');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('cleaning', 'zh', '清洁用品', 'cleaning', 'home-kitchen/cleaning');
INSERT INTO categories (id, icon, allows_used, restricted, sort_order) VALUES ('food', 'groceries', false, true, 100);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('food', 'en', 'Food & Drinks', 'food-drinks', 'food-drinks')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('food', 'km', 'អាហារ និងភេសជ្ជៈ', 'food-drinks', 'food-drinks')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('food', 'zh', '食品饮料', 'food-drinks', 'food-drinks')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('snacks', 'food', false, true, 10);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('snacks', 'en', 'Snacks', 'snacks', 'food-drinks/snacks');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('snacks', 'km', 'នំចំណី', 'snacks', 'food-drinks/snacks');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('snacks', 'zh', '零食', 'snacks', 'food-drinks/snacks');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('coffee-tea', 'food', false, true, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('coffee-tea', 'en', 'Coffee & Tea', 'coffee-tea', 'food-drinks/coffee-tea');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('coffee-tea', 'km', 'កាហ្វេ និងតែ', 'coffee-tea', 'food-drinks/coffee-tea');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('coffee-tea', 'zh', '咖啡和茶', 'coffee-tea', 'food-drinks/coffee-tea');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('cooking-ingredients', 'food', false, true, 30);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('cooking-ingredients', 'en', 'Spices & Ingredients', 'cooking-ingredients', 'food-drinks/cooking-ingredients');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('cooking-ingredients', 'km', 'គ្រឿងទេស និងគ្រឿងផ្សំ', 'cooking-ingredients', 'food-drinks/cooking-ingredients');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('cooking-ingredients', 'zh', '调味料和食材', 'cooking-ingredients', 'food-drinks/cooking-ingredients');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('packaged-food', 'food', false, true, 40);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('packaged-food', 'en', 'Packaged Food', 'packaged-food', 'food-drinks/packaged-food');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('packaged-food', 'km', 'អាហារវេចខ្ចប់', 'packaged-food', 'food-drinks/packaged-food');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('packaged-food', 'zh', '包装食品', 'packaged-food', 'food-drinks/packaged-food');
INSERT INTO categories (id, icon, allows_used, restricted, sort_order) VALUES ('sports-auto', 'ball', true, false, 110);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('sports-auto', 'en', 'Sports & Automotive', 'sports-automotive', 'sports-automotive')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('sports-auto', 'km', 'កីឡា និងយានយន្ត', 'sports-automotive', 'sports-automotive')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('sports-auto', 'zh', '运动和汽车', 'sports-automotive', 'sports-automotive')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('sportswear', 'sports-auto', true, false, 10);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('sportswear', 'en', 'Sportswear', 'sportswear', 'sports-automotive/sportswear');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('sportswear', 'km', 'សម្លៀកបំពាក់កីឡា', 'sportswear', 'sports-automotive/sportswear');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('sportswear', 'zh', '运动服', 'sportswear', 'sports-automotive/sportswear');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('fitness', 'sports-auto', true, false, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('fitness', 'en', 'Fitness', 'fitness', 'sports-automotive/fitness');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('fitness', 'km', 'ហាត់ប្រាណ', 'fitness', 'sports-automotive/fitness');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('fitness', 'zh', '健身', 'fitness', 'sports-automotive/fitness');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('outdoor', 'sports-auto', true, false, 30);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('outdoor', 'en', 'Outdoor & Beach', 'outdoor', 'sports-automotive/outdoor');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('outdoor', 'km', 'ក្រៅផ្ទះ និងឆ្នេរ', 'outdoor', 'sports-automotive/outdoor');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('outdoor', 'zh', '户外和海滩', 'outdoor', 'sports-automotive/outdoor');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('bicycles', 'sports-auto', true, false, 40);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('bicycles', 'en', 'Bicycles', 'bicycles', 'sports-automotive/bicycles');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('bicycles', 'km', 'កង់', 'bicycles', 'sports-automotive/bicycles');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('bicycles', 'zh', '自行车', 'bicycles', 'sports-automotive/bicycles');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('motorbike-accessories', 'sports-auto', true, false, 50);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('motorbike-accessories', 'en', 'Motorbike Accessories', 'motorbike-accessories', 'sports-automotive/motorbike-accessories');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('motorbike-accessories', 'km', 'គ្រឿងម៉ូតូ', 'motorbike-accessories', 'sports-automotive/motorbike-accessories');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('motorbike-accessories', 'zh', '摩托车配件', 'motorbike-accessories', 'sports-automotive/motorbike-accessories');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('helmets', 'sports-auto', true, false, 60);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('helmets', 'en', 'Helmets', 'helmets', 'sports-automotive/helmets');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('helmets', 'km', 'មួកសុវត្ថិភាព', 'helmets', 'sports-automotive/helmets');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('helmets', 'zh', '头盔', 'helmets', 'sports-automotive/helmets');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('car-accessories', 'sports-auto', true, false, 70);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('car-accessories', 'en', 'Car Accessories', 'car-accessories', 'sports-automotive/car-accessories');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('car-accessories', 'km', 'គ្រឿងរថយន្ត', 'car-accessories', 'sports-automotive/car-accessories');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('car-accessories', 'zh', '汽车配件', 'car-accessories', 'sports-automotive/car-accessories');
INSERT INTO categories (id, icon, allows_used, restricted, sort_order) VALUES ('hobbies', 'book', true, false, 120);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('hobbies', 'en', 'Hobbies & More', 'hobbies-more', 'hobbies-more')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('hobbies', 'km', 'ចំណូលចិត្ត និងផ្សេងៗ', 'hobbies-more', 'hobbies-more')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('hobbies', 'zh', '兴趣和其他', 'hobbies-more', 'hobbies-more')
  ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path;
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('books-stationery', 'hobbies', true, false, 10);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('books-stationery', 'en', 'Books & Stationery', 'books-stationery', 'hobbies-more/books-stationery');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('books-stationery', 'km', 'សៀវភៅ និងសម្ភារៈសិក្សា', 'books-stationery', 'hobbies-more/books-stationery');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('books-stationery', 'zh', '图书文具', 'books-stationery', 'hobbies-more/books-stationery');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('toys-collectibles', 'hobbies', true, false, 20);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('toys-collectibles', 'en', 'Toys & Collectibles', 'toys-collectibles', 'hobbies-more/toys-collectibles');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('toys-collectibles', 'km', 'ប្រដាប់លេង និងរបស់ប្រមូល', 'toys-collectibles', 'hobbies-more/toys-collectibles');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('toys-collectibles', 'zh', '玩具和收藏品', 'toys-collectibles', 'hobbies-more/toys-collectibles');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('musical-instruments', 'hobbies', true, false, 30);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('musical-instruments', 'en', 'Musical Instruments', 'musical-instruments', 'hobbies-more/musical-instruments');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('musical-instruments', 'km', 'ឧបករណ៍តន្ត្រី', 'musical-instruments', 'hobbies-more/musical-instruments');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('musical-instruments', 'zh', '乐器', 'musical-instruments', 'hobbies-more/musical-instruments');
INSERT INTO categories (id, parent_id, allows_used, restricted, sort_order) VALUES ('pet-supplies', 'hobbies', true, false, 40);
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('pet-supplies', 'en', 'Pet Supplies', 'pet-supplies', 'hobbies-more/pet-supplies');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('pet-supplies', 'km', 'សម្ភារៈសត្វចិញ្ចឹម', 'pet-supplies', 'hobbies-more/pet-supplies');
INSERT INTO category_translations (category_id, lang, name, slug, path) VALUES ('pet-supplies', 'zh', '宠物用品', 'pet-supplies', 'hobbies-more/pet-supplies');

-- Produk lama dipindah ke sub-kategori yang paling cocok berdasarkan nama (bahasa apa pun).
-- Yang tidak cocok (mis. kaos tanpa keterangan pria/wanita) tetap di kategori utama; seller diminta memilih sub-kategori saat mengedit.
CREATE TEMP TABLE cat_remap (prio int, from_id text, pattern text, to_id text) ON COMMIT DROP;
INSERT INTO cat_remap VALUES
  (0, 'electronics', 'power ?bank', 'power-banks'),
  (1, 'electronics', '\m(cable|charger|charging|adapter)\M', 'chargers-cables'),
  (2, 'electronics', '\m(phone case|case|screen protector|tempered glass)\M', 'phone-cases'),
  (3, 'electronics', '\m(speaker|earbuds?|earphones?|headphones?|headset|soundbar)\M', 'audio'),
  (4, 'electronics', '\m(laptop|notebook|macbook|thinkpad|chromebook)\M', 'laptops'),
  (5, 'electronics', '\m(camera|webcam|gopro|selfie|tripod|gimbal)\M', 'cameras'),
  (6, 'electronics', '\m(mouse|keyboard|usb hub|monitor|ssd|flash drive|router)\M', 'computer-accessories'),
  (7, 'electronics', '\m(tablet|ipad)\M', 'tablets'),
  (8, 'electronics', '\m(phone|iphone|smartphone|galaxy|redmi|xiaomi|oppo|vivo|realme|huawei|pixel)\M', 'mobile-phones'),
  (9, 'electronics', '\m(game|gaming|controller|console|playstation|nintendo)\M', 'gaming'),
  (10, 'electronics', '\m(tv|television|projector)\M', 'tvs'),
  (11, 'electronics', '\m(smart ?watch|watch)\M', 'watches'),
  (12, 'beauty', '\m(perfume|fragrance|cologne|eau de)\M', 'fragrance'),
  (13, 'beauty', '\m(shampoo|conditioner|hair)\M', 'haircare'),
  (14, 'beauty', '\m(lipstick|lip|mascara|foundation|powder|blush|eyeliner|eyeshadow|cushion|tint)\M', 'makeup'),
  (15, 'beauty', '\m(serum|toner|sunscreen|cleanser|moisturi[sz]er|face|facial|mask|cream)\M', 'skincare'),
  (16, 'beauty', '\m(lotion|soap|body|shower)\M', 'body-care'),
  (17, 'fashion', '\m(men|mens|men''s|male|boy)\M.*\m(shoes?|sneakers?|sandals?|slippers?)\M|\m(shoes?|sneakers?|sandals?|slippers?)\M.*\m(men|mens|men''s)\M', 'mens-shoes'),
  (18, 'fashion', '\m(wallet)\M', 'mens-bags'),
  (19, 'fashion', '\m(men|mens|men''s|male)\M.*\m(pants|jeans|shorts|trousers)\M', 'mens-bottoms'),
  (20, 'fashion', '\m(men|mens|men''s|male)\M', 'mens-tops'),
  (21, 'fashion', '\m(dress|skirt)\M', 'womens-dresses'),
  (22, 'fashion', '\m(heels?|women|womens|women''s|ladies)\M.*\m(shoes?|sneakers?|sandals?|slippers?)\M|\m(high heels?|pumps)\M', 'womens-shoes'),
  (23, 'fashion', '\m(tote|clutch|handbag|crossbody|purse)\M', 'womens-bags'),
  (24, 'fashion', '\m(leggings)\M|\m(women|womens|women''s|ladies)\M.*\m(pants|jeans|shorts|trousers)\M', 'womens-bottoms'),
  (25, 'fashion', '\m(bra|underwear|pajamas?|sleepwear|lingerie)\M', 'womens-underwear'),
  (26, 'fashion', '\m(blouse|crop top)\M|\m(women|womens|women''s|ladies)\M.*\m(shirt|t-shirt|tee|top|hoodie|jacket|sweater)\M', 'womens-tops'),
  (27, 'accessories', '\m(watch)\M', 'watches'),
  (28, 'accessories', '\m(ring|necklace|bracelet|earrings?|pendant|jewel(le)?ry)\M', 'jewelry'),
  (29, 'accessories', '\m(sunglasses|glasses|eyewear)\M', 'eyewear'),
  (30, 'accessories', '\m(hat|cap|beanie)\M', 'hats'),
  (31, 'accessories', '\m(belt)\M', 'belts'),
  (32, 'accessories', '\m(wallet)\M', 'mens-bags'),
  (33, 'accessories', '\m(tote|clutch|handbag|crossbody|purse)\M', 'womens-bags'),
  (34, 'accessories', '\m(phone case|case)\M', 'phone-cases'),
  (35, 'kitchen', '\m(rice cooker|blender|kettle|air fryer|fryer|fan|vacuum|iron|microwave|toaster)\M', 'home-appliances'),
  (36, 'kitchen', '\m(pan|pot|wok|knife|knives|mug|cup|plate|bowl|spatula|cutting board|chopsticks|glass|bottle|lunch box|container)\M', 'kitchenware'),
  (37, 'kitchen', '\m(chair|table|shelf|sofa|desk|cabinet|stool|rack)\M', 'furniture'),
  (38, 'kitchen', '\m(lamp|vase|decor|candle|frame|pillow|cushion|curtain|rug|mirror)\M', 'home-decor'),
  (39, 'kitchen', '\m(mop|broom|detergent|cleaner|cleaning|laundry|towel)\M', 'cleaning');
UPDATE products p SET category_id = m.to_id
FROM (
  SELECT DISTINCT ON (p2.id) p2.id AS product_id, r.to_id
  FROM products p2
  JOIN cat_remap r ON r.from_id = p2.category_id
  WHERE EXISTS (SELECT 1 FROM product_translations t WHERE t.product_id = p2.id AND t.name ~* r.pattern)
  ORDER BY p2.id, r.prio
) m
WHERE p.id = m.product_id;

-- +goose Down
UPDATE products p SET category_id = c.parent_id FROM categories c WHERE c.id = p.category_id AND c.parent_id IS NOT NULL;
UPDATE products SET category_id = 'electronics' WHERE category_id = 'phones';
UPDATE products SET category_id = 'fashion' WHERE category_id = 'mens-fashion';
UPDATE products SET category_id = 'beauty' WHERE category_id = 'health';
UPDATE products SET category_id = 'kitchen' WHERE category_id IN ('mother-baby', 'food', 'sports-auto', 'hobbies');
UPDATE content_collections SET category_id = NULL WHERE category_id IN (SELECT id FROM categories WHERE parent_id IS NOT NULL OR id IN ('phones','mens-fashion','health','mother-baby','food','sports-auto','hobbies'));
DELETE FROM category_attributes WHERE category_id IN (SELECT id FROM categories WHERE parent_id IS NOT NULL OR id IN ('phones','mens-fashion','health','mother-baby','food','sports-auto','hobbies'));
DELETE FROM brand_categories WHERE category_id IN (SELECT id FROM categories WHERE parent_id IS NOT NULL OR id IN ('phones','mens-fashion','health','mother-baby','food','sports-auto','hobbies'));
DELETE FROM commission_rules WHERE category_id IN (SELECT id FROM categories WHERE parent_id IS NOT NULL OR id IN ('phones','mens-fashion','health','mother-baby','food','sports-auto','hobbies'));
DELETE FROM category_translations WHERE category_id IN (SELECT id FROM categories WHERE parent_id IS NOT NULL OR id IN ('phones','mens-fashion','health','mother-baby','food','sports-auto','hobbies'));
DELETE FROM categories WHERE parent_id IS NOT NULL;
DELETE FROM categories WHERE id IN ('phones','mens-fashion','health','mother-baby','food','sports-auto','hobbies');
UPDATE categories SET icon = 'shirt' WHERE id = 'fashion';
UPDATE categories SET icon = 'star' WHERE id = 'accessories';
UPDATE categories SET icon = 'headphones' WHERE id = 'electronics';
UPDATE category_translations SET name = CASE lang WHEN 'en' THEN 'Fashion' WHEN 'km' THEN 'សម្លៀកបំពាក់' ELSE '时尚' END, slug = 'fashion', path = 'fashion' WHERE category_id = 'fashion';
UPDATE category_translations SET name = CASE lang WHEN 'en' THEN 'Accessories' WHEN 'km' THEN 'គ្រឿងបន្ថែម' ELSE '配饰' END, slug = 'accessories', path = 'accessories' WHERE category_id = 'accessories';
UPDATE category_translations SET name = CASE lang WHEN 'en' THEN 'Electronics' WHEN 'km' THEN 'អេឡិចត្រូនិក' ELSE '电子产品' END WHERE category_id = 'electronics';
