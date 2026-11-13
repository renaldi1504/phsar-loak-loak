-- +goose Up
-- Moderasi produk, terjemahan AI, komisi, dokumen legal, kurs.

CREATE TABLE prohibited_keywords (
  id       integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  keyword  varchar(80) NOT NULL,
  lang     lang_code NOT NULL,
  category varchar(40) NOT NULL,                   -- weapons · drugs · protected_wildlife · counterfeit · …
  action   varchar(8) NOT NULL CHECK (action IN ('block','review')),
  UNIQUE (lang, keyword)
);

-- Riwayat keputusan moderator; append-only.
CREATE TABLE product_moderation_reviews (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id  uuid NOT NULL REFERENCES products(id),
  decision    moderation_status NOT NULL CHECK (decision <> 'pending_review'),
  reason_code varchar(32) CHECK (reason_code IN ('defect_photo_missing','unclear_condition','prohibited','counterfeit','wrong_category','price','other')),
  note        text,                                 -- terlihat oleh seller
  auto_flags  jsonb NOT NULL DEFAULT '{}',          -- snapshot cek otomatis saat diputuskan
  reviewer_id uuid NOT NULL REFERENCES users(id),
  created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX product_moderation_reviews_product_idx ON product_moderation_reviews (product_id, created_at DESC);

CREATE TABLE media_moderation_results (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  media_id       uuid NOT NULL REFERENCES product_media(id),
  model          varchar(64) NOT NULL,
  labels         jsonb NOT NULL,                    -- {"nudity":0.02,"weapon":0.01}
  max_score      numeric(4,3) NOT NULL,
  auto_decision  varchar(8) NOT NULL CHECK (auto_decision IN ('allow','review','block')),
  final_decision varchar(8) CHECK (final_decision IN ('allow','block')),  -- moderasi otomatis tidak memutuskan sendiri
  reviewed_by    uuid REFERENCES users(id),
  reviewed_at    timestamptz,
  created_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX media_moderation_results_media_idx ON media_moderation_results (media_id);

-- Antrean terjemahan AI: seller menulis 1 bahasa, sistem mengisi 2 bahasa lainnya.
CREATE TABLE translation_jobs (
  id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  entity_table varchar(64) NOT NULL,                -- product_translations · store_translations · …
  entity_id    text NOT NULL,
  source_lang  lang_code NOT NULL,
  target_lang  lang_code NOT NULL CHECK (target_lang <> source_lang),
  model        varchar(64),
  status       varchar(16) NOT NULL DEFAULT 'queued' CHECK (status IN ('queued','running','done','failed','skipped')),
  attempts     integer NOT NULL DEFAULT 0,
  error        text,
  priority     smallint NOT NULL DEFAULT 0,
  created_at   timestamptz NOT NULL DEFAULT now(),
  done_at      timestamptz
);
CREATE INDEX translation_jobs_queue ON translation_jobs (priority DESC, created_at) WHERE status = 'queued';

-- Komisi: aturan paling spesifik menang (store > category > global). Launch: global 8%.
CREATE TABLE commission_rules (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  scope        varchar(16) NOT NULL CHECK (scope IN ('global','category','store')),
  category_id  varchar(60) REFERENCES categories(id),
  store_id     uuid REFERENCES stores(id),
  rate_percent numeric(5,2) NOT NULL CHECK (rate_percent >= 0 AND rate_percent <= 100),
  fixed_fee    numeric(12,2) NOT NULL DEFAULT 0.00 CHECK (fixed_fee >= 0),
  min_fee      numeric(12,2),
  max_fee      numeric(12,2),
  valid_from   timestamptz NOT NULL DEFAULT now(),
  valid_to     timestamptz,
  created_by   uuid REFERENCES users(id),
  created_at   timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT commission_scope_target CHECK (
    (scope = 'global'   AND category_id IS NULL     AND store_id IS NULL) OR
    (scope = 'category' AND category_id IS NOT NULL AND store_id IS NULL) OR
    (scope = 'store'    AND store_id IS NOT NULL)
  ),
  CONSTRAINT commission_valid_range CHECK (valid_to IS NULL OR valid_to > valid_from)
);
CREATE INDEX commission_rules_lookup ON commission_rules (scope, category_id, store_id, valid_from DESC);
CREATE TRIGGER commission_rules_audit AFTER INSERT OR UPDATE OR DELETE ON commission_rules FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Dokumen legal berversi + bukti persetujuan (append-only).
CREATE TABLE legal_documents (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  type         varchar(32) NOT NULL CHECK (type IN ('buyer_terms','seller_agreement','privacy_policy','ads_terms')),
  version      varchar(16) NOT NULL,
  lang         lang_code NOT NULL,
  body         text NOT NULL,
  content_hash char(64) NOT NULL,
  effective_at timestamptz NOT NULL,
  created_at   timestamptz NOT NULL DEFAULT now(),
  UNIQUE (type, version, lang)
);

CREATE TABLE terms_acceptances (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     uuid NOT NULL REFERENCES users(id),
  document_id uuid NOT NULL REFERENCES legal_documents(id),
  store_id    uuid REFERENCES stores(id),           -- untuk seller_agreement
  accepted_at timestamptz NOT NULL DEFAULT now(),
  ip          inet,
  user_agent  text
);
CREATE INDEX terms_acceptances_user_idx ON terms_acceptances (user_id, accepted_at DESC);

CREATE TABLE exchange_rates (
  id           integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  rate_date    date NOT NULL,
  riel_per_usd numeric(10,4) NOT NULL CHECK (riel_per_usd > 0),
  source       varchar(8) NOT NULL CHECK (source IN ('nbc','payway','manual')),
  fetched_at   timestamptz NOT NULL DEFAULT now(),
  UNIQUE (rate_date, source)
);

-- +goose Down
DROP TABLE exchange_rates;
DROP TABLE terms_acceptances;
DROP TABLE legal_documents;
DROP TABLE commission_rules;
DROP TABLE translation_jobs;
DROP TABLE media_moderation_results;
DROP TABLE product_moderation_reviews;
DROP TABLE prohibited_keywords;
