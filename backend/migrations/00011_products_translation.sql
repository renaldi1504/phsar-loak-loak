-- +goose Up
-- Produk seller + worker terjemahan AI.

-- Job yang gagal dicoba lagi nanti (backoff), bukan langsung.
ALTER TABLE translation_jobs ADD COLUMN run_after timestamptz;
CREATE INDEX translation_jobs_entity ON translation_jobs (entity_table, entity_id) WHERE status IN ('queued', 'running');

-- Kapan seller mengirim produk ke moderasi + hasil cek otomatis (kata terlarang, foto cacat, dll).
ALTER TABLE products ADD COLUMN submitted_at timestamptz;
ALTER TABLE products ADD COLUMN auto_flags jsonb NOT NULL DEFAULT '{}';

-- Daftar awal kata yang perlu dicek moderator (bukan blokir otomatis). Admin bisa menambah.
INSERT INTO prohibited_keywords (keyword, lang, category, action) VALUES
  ('gun', 'en', 'weapons', 'review'),
  ('pistol', 'en', 'weapons', 'review'),
  ('rifle', 'en', 'weapons', 'review'),
  ('ammunition', 'en', 'weapons', 'review'),
  ('grenade', 'en', 'weapons', 'review'),
  ('taser', 'en', 'weapons', 'review'),
  ('cannabis', 'en', 'drugs', 'review'),
  ('marijuana', 'en', 'drugs', 'review'),
  ('methamphetamine', 'en', 'drugs', 'review'),
  ('cocaine', 'en', 'drugs', 'review'),
  ('heroin', 'en', 'drugs', 'review'),
  ('ivory', 'en', 'protected_wildlife', 'review'),
  ('pangolin', 'en', 'protected_wildlife', 'review'),
  ('tiger skin', 'en', 'protected_wildlife', 'review'),
  ('replica', 'en', 'counterfeit', 'review'),
  ('counterfeit', 'en', 'counterfeit', 'review'),
  ('kw super', 'en', 'counterfeit', 'review')
ON CONFLICT (lang, keyword) DO NOTHING;

-- +goose Down
DELETE FROM prohibited_keywords WHERE action = 'review' AND lang = 'en';
ALTER TABLE products DROP COLUMN auto_flags;
ALTER TABLE products DROP COLUMN submitted_at;
DROP INDEX translation_jobs_entity;
ALTER TABLE translation_jobs DROP COLUMN run_after;
