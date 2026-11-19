-- +goose Up
-- Jadwal latar (worker): riwayat run yang berarti + detak terakhir per job. Jeda per job ada di platform_settings['jobs'].
ALTER TABLE job_runs
  ADD COLUMN duration_ms  integer,
  ADD COLUMN trigger      varchar(8) NOT NULL DEFAULT 'auto' CHECK (trigger IN ('auto','manual')),
  ADD COLUMN triggered_by uuid REFERENCES users(id);
CREATE INDEX job_runs_started_idx ON job_runs (started_at);

-- Detak terakhir: worker yang jalan tiap beberapa detik tidak mencatat setiap putaran kosong di job_runs,
-- cukup memperbarui baris ini (paling sering tiap 30 detik).
CREATE TABLE job_state (
  job             varchar(64) PRIMARY KEY,
  last_tick_at    timestamptz,
  last_success_at timestamptz,
  last_error      text,
  last_error_at   timestamptz,
  updated_at      timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE job_state;
DROP INDEX job_runs_started_idx;
ALTER TABLE job_runs DROP COLUMN triggered_by, DROP COLUMN trigger, DROP COLUMN duration_ms;
