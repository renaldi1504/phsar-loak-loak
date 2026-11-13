-- +goose Up
-- Fondasi: extension, enum bersama, trigger updated_at, dan audit log.
CREATE EXTENSION IF NOT EXISTS pgcrypto;  -- gen_random_uuid(), enkripsi rekening
CREATE EXTENSION IF NOT EXISTS citext;    -- email tidak peka huruf besar/kecil
-- pgvector (cari pakai foto) diaktifkan di migrasi Tahap 2 bersama product_embeddings.

CREATE TYPE lang_code AS ENUM ('en', 'km', 'zh');
CREATE TYPE actor_type AS ENUM ('buyer', 'seller', 'admin', 'courier', 'system');

-- +goose StatementBegin
CREATE FUNCTION set_updated_at() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  NEW.updated_at := now();
  RETURN NEW;
END $$;
-- +goose StatementEnd

-- Audit log: hanya bisa ditambah. Dipartisi per bulan.
CREATE TABLE audit_logs (
  id             bigint GENERATED ALWAYS AS IDENTITY,
  occurred_at    timestamptz NOT NULL DEFAULT now(),
  actor_type     actor_type  NOT NULL DEFAULT 'system',
  actor_user_id  uuid,
  store_id       uuid,
  action         varchar(64) NOT NULL,        -- insert · update · soft_delete · restore · delete · atau aksi bisnis
  entity_table   varchar(64) NOT NULL,
  entity_id      text,
  old_data       jsonb,
  new_data       jsonb,
  changed_fields text[],
  request_id     text,
  ip             inet,
  user_agent     text,
  PRIMARY KEY (id, occurred_at)
) PARTITION BY RANGE (occurred_at);

CREATE TABLE audit_logs_default PARTITION OF audit_logs DEFAULT;
CREATE INDEX audit_logs_entity_idx ON audit_logs (entity_table, entity_id, occurred_at DESC);
CREATE INDEX audit_logs_store_idx  ON audit_logs (store_id, occurred_at DESC) WHERE store_id IS NOT NULL;
CREATE INDEX audit_logs_actor_idx  ON audit_logs (actor_user_id, occurred_at DESC) WHERE actor_user_id IS NOT NULL;

-- Partisi bulanan dibuat oleh job terjadwal; fungsi ini juga dipanggil di sini untuk 14 bulan pertama.
-- +goose StatementBegin
CREATE FUNCTION ensure_audit_partition(month_start date) RETURNS void LANGUAGE plpgsql AS $$
DECLARE
  part text := 'audit_logs_' || to_char(month_start, 'YYYY_MM');
BEGIN
  IF to_regclass(part) IS NULL THEN
    EXECUTE format('CREATE TABLE %I PARTITION OF audit_logs FOR VALUES FROM (%L) TO (%L)',
                   part, month_start, (month_start + interval '1 month')::date);
  END IF;
END $$;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
DECLARE m date := date_trunc('month', now())::date;
BEGIN
  FOR i IN 0..13 LOOP
    PERFORM ensure_audit_partition((m + make_interval(months => i))::date);
  END LOOP;
END $$;
-- +goose StatementEnd

-- Larang UPDATE/DELETE di audit_logs.
-- +goose StatementBegin
CREATE FUNCTION audit_logs_immutable() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  RAISE EXCEPTION 'audit_logs is append-only';
END $$;
-- +goose StatementEnd
CREATE TRIGGER audit_logs_no_update BEFORE UPDATE OR DELETE ON audit_logs
  FOR EACH ROW EXECUTE FUNCTION audit_logs_immutable();

-- Trigger generik: catat perubahan baris. Siapa pelakunya dibaca dari
-- SET LOCAL app.actor_id / app.actor_type / app.request_id (lihat db.WithActor di Go).
-- +goose StatementBegin
CREATE FUNCTION audit_row_change() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
  old_j jsonb := CASE WHEN TG_OP IN ('UPDATE','DELETE') THEN to_jsonb(OLD) END;
  new_j jsonb := CASE WHEN TG_OP IN ('INSERT','UPDATE') THEN to_jsonb(NEW) END;
  changed text[];
  act text := lower(TG_OP);
  actor text := nullif(current_setting('app.actor_id', true), '');
  atype text := coalesce(nullif(current_setting('app.actor_type', true), ''), 'system');
  rec jsonb := coalesce(new_j, old_j);
  sid uuid;
BEGIN
  IF TG_OP = 'UPDATE' THEN
    SELECT array_agg(k ORDER BY k) INTO changed
    FROM jsonb_object_keys(new_j) k
    WHERE k <> 'updated_at' AND new_j -> k IS DISTINCT FROM old_j -> k;
    IF changed IS NULL THEN
      RETURN NEW;  -- tidak ada yang berubah
    END IF;
    IF 'deleted_at' = ANY(changed) THEN
      act := CASE WHEN new_j ->> 'deleted_at' IS NULL THEN 'restore' ELSE 'soft_delete' END;
    END IF;
    SELECT jsonb_object_agg(k, old_j -> k) INTO old_j FROM unnest(changed) k;
    SELECT jsonb_object_agg(k, new_j -> k) INTO new_j FROM unnest(changed) k;
  END IF;

  IF TG_TABLE_NAME = 'stores' THEN
    sid := (rec ->> 'id')::uuid;
  ELSIF rec ? 'store_id' THEN
    sid := (rec ->> 'store_id')::uuid;
  END IF;

  INSERT INTO audit_logs (actor_type, actor_user_id, store_id, action, entity_table, entity_id,
                          old_data, new_data, changed_fields, request_id)
  VALUES (atype::actor_type, actor::uuid, sid, act, TG_TABLE_NAME,
          coalesce(rec ->> 'id', rec ->> 'code', rec ->> 'user_id'),
          old_j, new_j, changed, nullif(current_setting('app.request_id', true), ''));
  RETURN coalesce(NEW, OLD);
END $$;
-- +goose StatementEnd

-- Jejak job terjadwal (release dana, payout, rekonsiliasi, partisi audit, kurs).
CREATE TABLE job_runs (
  id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  job           varchar(64) NOT NULL,
  started_at    timestamptz NOT NULL DEFAULT now(),
  finished_at   timestamptz,
  status        varchar(16) NOT NULL DEFAULT 'running' CHECK (status IN ('running','success','failed')),
  rows_affected integer,
  error         text,
  lock_key      varchar(64)
);
CREATE INDEX job_runs_job_idx ON job_runs (job, started_at DESC);

-- +goose Down
DROP TABLE job_runs;
DROP TABLE audit_logs;
DROP FUNCTION audit_row_change();
DROP FUNCTION audit_logs_immutable();
DROP FUNCTION ensure_audit_partition(date);
DROP FUNCTION set_updated_at();
DROP TYPE actor_type;
DROP TYPE lang_code;
